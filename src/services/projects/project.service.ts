import { prisma } from "lib/prisma.js";
import { ApiError } from "errors/ApiError.js";
import { StatusCodes } from "http-status-codes";
import type {
  CreateProjectInput,
  UpdateProjectInput,
  AddMemberInput,
  UpdateMemberRoleInput,
} from "types/projects/project.schema.js";
import type {
  ProjectResponse,
  ProjectMemberResponse,
  ProjectSummaryResponse,
} from "types/projects/project.dto.js";

const mapToProjectMemberResponse = (
  prismaMember: any,
): ProjectMemberResponse => {
  const { users: user, ...memberData } = prismaMember;

  const userSummary = {
    id: user.id,
    first_name: user.first_name,
    last_name: user.last_name,
    login_email: user.login_email,
  };

  return {
    ...memberData,
    user: userSummary,
  } as ProjectMemberResponse;
};

const mapToProjectResponse = (project: any): ProjectResponse => {
  const { users: manager, project_members, ...projectData } = project;

  const managerSummary = {
    id: manager.id,
    first_name: manager.first_name,
    last_name: manager.last_name,
    login_email: manager.login_email,
  };

  const members = project_members.map(mapToProjectMemberResponse);

  return {
    ...projectData,
    manager: managerSummary,

    business: project.businesses
      ? {
          id: project.businesses.id,
          name: project.businesses.name,
          business_type: project.businesses.business_type,
        }
      : null,
    members: members,
  } as ProjectResponse;
};

export class ProjectService {
  /**
   * Creates a new project record and automatically sets the managed_by_user_id
   * as the initial 'manager' role member.
   *
   * @param {CreateProjectInput} data The data for the new project.
   * @returns {Promise<ProjectResponse>} The newly created project with manager, business, and members included.
   * @throws {ApiError} 409 Conflict if a project with the same unique fields already exists.
   * @throws {ApiError} 500 Internal Server Error on database failure.
   */
  public async createProject(
    data: CreateProjectInput,
  ): Promise<ProjectResponse> {
    try {
      const project = await prisma.projects.create({
        data: {
          ...data,
          project_members: {
            create: {
              user_id: data.managed_by_user_id,
              role: "manager",
            },
          },
        },
        include: {
          users: true,
          businesses: true,
          project_members: { include: { users: true } },
        },
      });

      return mapToProjectResponse(project);
    } catch (error: any) {
      if (error.code === "P2002") {
        throw new ApiError(
          StatusCodes.CONFLICT,
          "A project with this name already exists.",
        );
      }
      throw new ApiError(
        StatusCodes.INTERNAL_SERVER_ERROR,
        "Failed to create project.",
      );
    }
  }

  /**
   * Retrieves a project by its unique ID, including detailed relationship data
   * for the manager, associated business, and project members.
   *
   * @param {number} projectId The unique ID of the project.
   * @returns {Promise<ProjectResponse | null>} The project details or null if not found.
   */
  public async getProjectById(
    projectId: number,
  ): Promise<ProjectResponse | null> {
    const project = await prisma.projects.findUnique({
      where: { id: projectId },
      include: {
        users: true,
        businesses: true,
        project_members: { include: { users: true } },
      },
    });

    if (!project) return null;

    return mapToProjectResponse(project);
  }

  /**
   * Updates core project details based on the provided partial data.
   *
   * @param {number} projectId The unique ID of the project to update.
   * @param {UpdateProjectInput} data The fields to be updated.
   * @returns {Promise<ProjectResponse>} The updated project object.
   * @throws {ApiError} 404 Not Found if the project does not exist.
   * @throws {ApiError} 500 Internal Server Error on database failure.
   */
  public async updateProject(
    projectId: number,
    data: UpdateProjectInput,
  ): Promise<ProjectResponse> {
    try {
      const updatedProject = await prisma.projects.update({
        where: { id: projectId },
        data: data,
        include: {
          users: true,
          businesses: true,
          project_members: { include: { users: true } },
        },
      });

      return mapToProjectResponse(updatedProject);
    } catch (error: any) {
      if (error.code === "P2025") {
        throw new ApiError(StatusCodes.NOT_FOUND, "Project not found.");
      }
      throw new ApiError(
        StatusCodes.INTERNAL_SERVER_ERROR,
        "Failed to update project.",
      );
    }
  }

  /**
   * Permanently deletes a project by ID, which should cascade to all
   * related project members and skills in the database.
   *
   * @param {number} projectId The unique ID of the project to delete.
   * @returns {Promise<boolean>} True if the project was deleted, false if the project was not found.
   * @throws {ApiError} 500 Internal Server Error on database failure.
   */
  public async deleteProject(projectId: number): Promise<boolean> {
    try {
      await prisma.projects.delete({
        where: { id: projectId },
      });
      return true;
    } catch (error: any) {
      if (error.code === "P2025") {
        return false;
      }
      throw new ApiError(
        StatusCodes.INTERNAL_SERVER_ERROR,
        "Failed to delete project.",
      );
    }
  }

  /**
   * Retrieves a list of all members for a specified project.
   *
   * @param {number} projectId The ID of the project.
   * @returns {Promise<ProjectMemberResponse[]>} A list of project member objects.
   * @throws {ApiError} 500 Internal Server Error on database failure.
   */
  public async getMembersByProjectId(
    projectId: number,
  ): Promise<ProjectMemberResponse[]> {
    try {
      const members = await prisma.project_members.findMany({
        where: { project_id: projectId },
        include: { users: true },
      });

      return members.map(mapToProjectMemberResponse);
    } catch (error) {
      throw new ApiError(
        StatusCodes.INTERNAL_SERVER_ERROR,
        "Failed to retrieve project members.",
      );
    }
  }

  /**
   * Retrieves a specific project member by their composite key (project ID and user ID).
   *
   * @param {number} projectId The ID of the project.
   * @param {number} userId The ID of the user (member).
   * @returns {Promise<ProjectMemberResponse | null>} The project member object or null if not found.
   */
  public async getMemberByKeys(
    projectId: number,
    userId: number,
  ): Promise<ProjectMemberResponse | null> {
    const member = await prisma.project_members.findUnique({
      where: {
        project_id_user_id: {
          project_id: projectId,
          user_id: userId,
        },
      },
      include: { users: true },
    });

    if (!member) return null;

    return mapToProjectMemberResponse(member);
  }

  /**
   * Adds a user to a project with a specified role.
   *
   * @param {number} projectId The ID of the project.
   * @param {AddMemberInput} data The user ID and role for the new member.
   * @returns {Promise<ProjectMemberResponse>} The newly created project member object.
   * @throws {ApiError} 404 Not Found if the project or user does not exist.
   * @throws {ApiError} 409 Conflict if the user is already a member.
   * @throws {ApiError} 500 Internal Server Error on database failure.
   */
  public async addMember(
    projectId: number,
    data: AddMemberInput,
  ): Promise<ProjectMemberResponse> {
    try {
      const member = await prisma.project_members.create({
        data: {
          project_id: projectId,
          user_id: data.user_id,
          role: data.role,
        },
        include: {
          users: true,
        },
      });

      return mapToProjectMemberResponse(member);
    } catch (error: any) {
      if (error.code === "P2003") {
        throw new ApiError(StatusCodes.NOT_FOUND, "Project or user not found.");
      }
      if (error.code === "P2002") {
        throw new ApiError(
          StatusCodes.CONFLICT,
          "User is already a member of this project.",
        );
      }
      throw new ApiError(
        StatusCodes.INTERNAL_SERVER_ERROR,
        "Failed to add project member.",
      );
    }
  }

  /**
   * Updates the role of an existing project member.
   *
   * @param {number} projectId The ID of the project.
   * @param {number} userId The ID of the user whose role is being updated.
   * @param {UpdateMemberRoleInput} data The new role.
   * @returns {Promise<ProjectMemberResponse>} The updated project member object.
   * @throws {ApiError} 404 Not Found if the member record does not exist.
   * @throws {ApiError} 500 Internal Server Error on database failure.
   */
  public async updateMemberRole(
    projectId: number,
    userId: number,
    data: UpdateMemberRoleInput,
  ): Promise<ProjectMemberResponse> {
    try {
      const updatedMember = await prisma.project_members.update({
        where: {
          project_id_user_id: {
            project_id: projectId,
            user_id: userId,
          },
        },
        data: {
          role: data.role,
        },
        include: {
          users: true,
        },
      });

      return mapToProjectMemberResponse(updatedMember);
    } catch (error: any) {
      if (error.code === "P2025") {
        throw new ApiError(StatusCodes.NOT_FOUND, "Project member not found.");
      }
      throw new ApiError(
        StatusCodes.INTERNAL_SERVER_ERROR,
        "Failed to update member role.",
      );
    }
  }

  /**
   * Removes a member from the project using the composite key.
   *
   * @param {number} projectId The ID of the project.
   * @param {number} userId The ID of the user (member) to remove.
   * @returns {Promise<boolean>} True if the member was successfully removed, false if not found.
   * @throws {ApiError} 500 Internal Server Error on database failure.
   */
  public async removeMember(
    projectId: number,
    userId: number,
  ): Promise<boolean> {
    try {
      await prisma.project_members.delete({
        where: {
          project_id_user_id: {
            project_id: projectId,
            user_id: userId,
          },
        },
      });
      return true;
    } catch (error: any) {
      if (error.code === "P2025") {
        return false;
      }
      throw new ApiError(
        StatusCodes.INTERNAL_SERVER_ERROR,
        "Failed to remove project member.",
      );
    }
  }
}
