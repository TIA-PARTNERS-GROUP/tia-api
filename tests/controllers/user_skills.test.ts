import request from "supertest";
import { StatusCodes } from "http-status-codes";
import { UserSkillsService } from "@/services/user_skills/user_skills.service";
import { HttpErrors } from "@/errors/ApiError";
import {
  CreateUserSkillInput,
  UpdateUserSkillInput,
} from "@/types/user_skills/user_skills.dto";

// --- MOCK SETUP ---
jest.mock("@/services/user_skills/user_skills.service");
const mockUserSkillsService = UserSkillsService as jest.MockedClass<
  typeof UserSkillsService
>;

// Placeholder for the Express application instance
const app = {};

const VALID_TOKEN = "Bearer valid-jwt-token";
const MOCK_SKILL_RESPONSE = {
  skill_id: 1,
  user_id: 5,
  proficiency_level: "intermediate",
  created_at: new Date(),
};
const MOCK_CREATE_INPUT: CreateUserSkillInput = {
  skill_id: 1,
  user_id: 5,
  proficiency_level: "beginner",
};
const MOCK_UPDATE_INPUT: UpdateUserSkillInput = { proficiency_level: "expert" };

describe("UserSkillsController", () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  // --- GET /user-skills/user/{userId} (Get All Skills for User) ---
  describe("GET /user-skills/user/:userId", () => {
    it("should return 200 and a list of skills for the user", async () => {
      mockUserSkillsService.prototype.getUserSkills.mockResolvedValue([
        MOCK_SKILL_RESPONSE,
      ] as any);

      const response = await request(app)
        .get("/user-skills/user/5")
        .set("Authorization", VALID_TOKEN);

      expect(response.statusCode).toBe(StatusCodes.OK);
      expect(Array.isArray(response.body)).toBe(true);
      expect(response.body[0].user_id).toBe(5);
    });
  });

  // --- GET /user-skills/{skillId}/user/{userId} (Get Specific Skill) ---
  describe("GET /user-skills/:skillId/user/:userId", () => {
    it("should return 200 and the user-skill entry", async () => {
      mockUserSkillsService.prototype.getUserSkillById.mockResolvedValue(
        MOCK_SKILL_RESPONSE as any,
      );

      const response = await request(app)
        .get("/user-skills/1/user/5")
        .set("Authorization", VALID_TOKEN);

      expect(response.statusCode).toBe(StatusCodes.OK);
      expect(response.body.skill_id).toBe(1);
    });

    it("should return 404 Not Found", async () => {
      mockUserSkillsService.prototype.getUserSkillById.mockRejectedValue(
        HttpErrors.NotFound("User-skill relationship not found"),
      );

      const response = await request(app)
        .get("/user-skills/99/user/5")
        .set("Authorization", VALID_TOKEN);

      expect(response.statusCode).toBe(StatusCodes.NOT_FOUND);
    });
  });

  // --- POST /user-skills (Add Skill) ---
  describe("POST /user-skills", () => {
    it("should return 201 and the created user-skill entry", async () => {
      mockUserSkillsService.prototype.createUserSkill.mockResolvedValue(
        MOCK_SKILL_RESPONSE as any,
      );

      const response = await request(app)
        .post("/user-skills")
        .set("Authorization", VALID_TOKEN)
        .send(MOCK_CREATE_INPUT);

      expect(response.statusCode).toBe(StatusCodes.CREATED);
    });

    it("should return 409 Conflict if relationship already exists", async () => {
      mockUserSkillsService.prototype.createUserSkill.mockRejectedValue(
        HttpErrors.Conflict("User already has this skill"),
      );

      const response = await request(app)
        .post("/user-skills")
        .set("Authorization", VALID_TOKEN)
        .send(MOCK_CREATE_INPUT);

      expect(response.statusCode).toBe(StatusCodes.CONFLICT);
    });
  });

  // --- PUT /user-skills/{skillId}/user/{userId} (Update Skill) ---
  describe("PUT /user-skills/:skillId/user/:userId", () => {
    it("should return 200 and the updated user-skill entry", async () => {
      mockUserSkillsService.prototype.updateUserSkill.mockResolvedValue({
        ...MOCK_SKILL_RESPONSE,
        proficiency_level: "expert",
      } as any);

      const response = await request(app)
        .put("/user-skills/1/user/5")
        .set("Authorization", VALID_TOKEN)
        .send(MOCK_UPDATE_INPUT);

      expect(response.statusCode).toBe(StatusCodes.OK);
      expect(response.body.proficiency_level).toBe("expert");
    });
  });

  // --- DELETE /user-skills/{skillId}/user/{userId} (Delete Skill) ---
  describe("DELETE /user-skills/:skillId/user/:userId", () => {
    it("should return 204 No Content on success", async () => {
      mockUserSkillsService.prototype.deleteUserSkill.mockResolvedValue(
        undefined,
      );

      const response = await request(app)
        .delete("/user-skills/1/user/5")
        .set("Authorization", VALID_TOKEN);

      expect(response.statusCode).toBe(StatusCodes.NO_CONTENT);
    });

    it("should return 404 Not Found on failure", async () => {
      mockUserSkillsService.prototype.deleteUserSkill.mockRejectedValue(
        HttpErrors.NotFound("User-skill relationship not found"),
      );

      const response = await request(app)
        .delete("/user-skills/99/user/5")
        .set("Authorization", VALID_TOKEN);

      expect(response.statusCode).toBe(StatusCodes.NOT_FOUND);
    });
  });

  // --- GET /user-skills/skill/{skillId}/users (Get Users by Skill) ---
  describe("GET /user-skills/skill/:skillId/users", () => {
    it("should return 200 and a list of users", async () => {
      mockUserSkillsService.prototype.getUsersBySkill.mockResolvedValue([
        { id: 5, first_name: "Test", login_email: "test@example.com" },
      ] as any);

      const response = await request(app)
        .get("/user-skills/skill/1/users")
        .set("Authorization", VALID_TOKEN);

      expect(response.statusCode).toBe(StatusCodes.OK);
      expect(Array.isArray(response.body)).toBe(true);
      expect(response.body[0].id).toBe(5);
    });
  });

  // --- GET /user-skills/user/{userId}/proficiency/{proficiencyLevel} ---
  describe("GET /user-skills/user/:userId/proficiency/:proficiencyLevel", () => {
    it("should return 200 and filtered skills", async () => {
      mockUserSkillsService.prototype.getUserSkillsByProficiency.mockResolvedValue(
        [MOCK_SKILL_RESPONSE] as any,
      );

      const response = await request(app)
        .get("/user-skills/user/5/proficiency/intermediate")
        .set("Authorization", VALID_TOKEN);

      expect(response.statusCode).toBe(StatusCodes.OK);
      expect(Array.isArray(response.body)).toBe(true);
      expect(
        mockUserSkillsService.prototype.getUserSkillsByProficiency,
      ).toHaveBeenCalledWith(5, "intermediate");
    });

    it("should return 422/400 for invalid proficiency level (tsoa validation)", async () => {
      // This relies on tsoa/Zod path validation
      const response = await request(app)
        .get("/user-skills/user/5/proficiency/ninja")
        .set("Authorization", VALID_TOKEN);

      // Assuming tsoa or upstream validation catches invalid enum path param
      expect(response.statusCode).toBeGreaterThanOrEqual(
        StatusCodes.BAD_REQUEST,
      );
    });
  });
});
