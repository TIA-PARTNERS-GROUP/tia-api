import request from "supertest";
import { StatusCodes } from "http-status-codes";
import { ProjectService } from "@/services/projects/project.service";
import { HttpErrors } from "@/errors/ApiError";
import {
  CreateProjectInput,
  UpdateProjectInput,
  AddMemberInput,
  UpdateMemberRoleInput,
} from "@/types/projects/project.schema";

// --- MOCK SETUP ---
jest.mock("@/services/projects/project.service");
const mockProjectService = ProjectService as jest.MockedClass<
  typeof ProjectService
>;

// Placeholder for the Express application instance
const app = {};

const VALID_TOKEN = "Bearer valid-jwt-token";
const MOCK_MEMBER_RESPONSE = {
  project_id: 1,
  user_id: 5,
  role: "contributor",
  joined_at: new Date(),
  user: {
    id: 5,
    first_name: "Bob",
    last_name: "Member",
    login_email: "bob@test.com",
  },
};
const MOCK_PROJECT_RESPONSE = {
  id: 1,
  name: "New App Project",
  project_status: "planning",
  managed_by_user_id: 1,
  members: [MOCK_MEMBER_RESPONSE],
  manager: MOCK_MEMBER_RESPONSE.user,
  created_at: new Date(),
  updated_at: new Date(),
};
const MOCK_CREATE_INPUT: CreateProjectInput = {
  managed_by_user_id: 1,
  name: "New Project",
  description: "Initial project plan",
  project_status: "planning",
};
const MOCK_ADD_MEMBER: AddMemberInput = { user_id: 10, role: "reviewer" };

describe("ProjectController", () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  // --- POST /projects (Create Project) ---
  describe("POST /projects", () => {
    it("should return 201 and the created project", async () => {
      mockProjectService.prototype.createProject.mockResolvedValue(
        MOCK_PROJECT_RESPONSE as any,
      );

      const response = await request(app)
        .post("/projects")
        .set("Authorization", VALID_TOKEN)
        .send(MOCK_CREATE_INPUT);

      expect(response.statusCode).toBe(StatusCodes.CREATED);
      expect(response.body.name).toBe("New App Project");
    });
  });

  // --- GET /projects/{projectId} ---
  describe("GET /projects/:projectId", () => {
    it("should return 200 and the project details", async () => {
      mockProjectService.prototype.getProjectById.mockResolvedValue(
        MOCK_PROJECT_RESPONSE as any,
      );

      const response = await request(app)
        .get("/projects/1")
        .set("Authorization", VALID_TOKEN);

      expect(response.statusCode).toBe(StatusCodes.OK);
      expect(response.body.id).toBe(1);
    });

    it("should return 404 Not Found", async () => {
      mockProjectService.prototype.getProjectById.mockResolvedValue(null);

      const response = await request(app)
        .get("/projects/999")
        .set("Authorization", VALID_TOKEN);

      expect(response.statusCode).toBe(StatusCodes.NOT_FOUND);
    });
  });

  // --- PUT /projects/{projectId} (Update Project) ---
  describe("PUT /projects/:projectId", () => {
    const updateInput: UpdateProjectInput = { project_status: "active" };

    it("should return 200 and the updated project", async () => {
      mockProjectService.prototype.updateProject.mockResolvedValue({
        ...MOCK_PROJECT_RESPONSE,
        ...updateInput,
      } as any);

      const response = await request(app)
        .put("/projects/1")
        .set("Authorization", VALID_TOKEN)
        .send(updateInput);

      expect(response.statusCode).toBe(StatusCodes.OK);
      expect(response.body.project_status).toBe("active");
    });
  });

  // --- DELETE /projects/{projectId} ---
  describe("DELETE /projects/:projectId", () => {
    it("should return 204 No Content on success", async () => {
      mockProjectService.prototype.deleteProject.mockResolvedValue(true);

      const response = await request(app)
        .delete("/projects/1")
        .set("Authorization", VALID_TOKEN);

      expect(response.statusCode).toBe(StatusCodes.NO_CONTENT);
    });

    it("should return 404 Not Found if deletion fails", async () => {
      mockProjectService.prototype.deleteProject.mockResolvedValue(false);

      const response = await request(app)
        .delete("/projects/999")
        .set("Authorization", VALID_TOKEN);

      expect(response.statusCode).toBe(StatusCodes.NOT_FOUND);
    });
  });

  // --- POST /projects/{projectId}/members (Add Member) ---
  describe("POST /projects/:projectId/members", () => {
    it("should return 201 and the new member", async () => {
      mockProjectService.prototype.addMember.mockResolvedValue({
        ...MOCK_MEMBER_RESPONSE,
        user_id: 10,
      } as any);

      const response = await request(app)
        .post("/projects/1/members")
        .set("Authorization", VALID_TOKEN)
        .send(MOCK_ADD_MEMBER);

      expect(response.statusCode).toBe(StatusCodes.CREATED);
      expect(response.body.user_id).toBe(10);
    });

    it("should return 409 Conflict if user is already a member", async () => {
      mockProjectService.prototype.addMember.mockRejectedValue(
        HttpErrors.Conflict("User is already a member"),
      );

      const response = await request(app)
        .post("/projects/1/members")
        .set("Authorization", VALID_TOKEN)
        .send(MOCK_ADD_MEMBER);

      expect(response.statusCode).toBe(StatusCodes.CONFLICT);
    });
  });

  // --- GET /projects/{projectId}/members ---
  describe("GET /projects/:projectId/members", () => {
    it("should return 200 and a list of members", async () => {
      mockProjectService.prototype.getMembersByProjectId.mockResolvedValue([
        MOCK_MEMBER_RESPONSE,
      ] as any);

      const response = await request(app)
        .get("/projects/1/members")
        .set("Authorization", VALID_TOKEN);

      expect(response.statusCode).toBe(StatusCodes.OK);
      expect(Array.isArray(response.body)).toBe(true);
    });
  });

  // --- GET /projects/{projectId}/members/{userId} ---
  describe("GET /projects/:projectId/members/:userId", () => {
    it("should return 200 and the member details", async () => {
      mockProjectService.prototype.getMemberByKeys.mockResolvedValue(
        MOCK_MEMBER_RESPONSE as any,
      );

      const response = await request(app)
        .get("/projects/1/members/5")
        .set("Authorization", VALID_TOKEN);

      expect(response.statusCode).toBe(StatusCodes.OK);
      expect(response.body.user_id).toBe(5);
    });

    it("should return 404 Not Found", async () => {
      mockProjectService.prototype.getMemberByKeys.mockResolvedValue(null);

      const response = await request(app)
        .get("/projects/1/members/999")
        .set("Authorization", VALID_TOKEN);

      expect(response.statusCode).toBe(StatusCodes.NOT_FOUND);
    });
  });

  // --- DELETE /projects/{projectId}/members/{userId} ---
  describe("DELETE /projects/:projectId/members/:userId", () => {
    it("should return 204 No Content", async () => {
      mockProjectService.prototype.removeMember.mockResolvedValue(true);

      const response = await request(app)
        .delete("/projects/1/members/5")
        .set("Authorization", VALID_TOKEN);

      expect(response.statusCode).toBe(StatusCodes.NO_CONTENT);
    });

    it("should return 404 Not Found if member does not exist", async () => {
      mockProjectService.prototype.removeMember.mockResolvedValue(false);

      const response = await request(app)
        .delete("/projects/1/members/999")
        .set("Authorization", VALID_TOKEN);

      expect(response.statusCode).toBe(StatusCodes.NOT_FOUND);
    });
  });
});
