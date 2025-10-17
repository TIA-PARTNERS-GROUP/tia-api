import request from "supertest";
import { StatusCodes } from "http-status-codes";
import * as userService from "@/services/user/user.services";
import { HttpErrors } from "@/errors/ApiError";
import { UserCreateInput, UserUpdateInput } from "@/types/user/user.validation";

// --- MOCK SETUP ---
jest.mock("@/services/user/user.services");
const mockUserService = userService as jest.Mocked<typeof userService>;

// Placeholder for the Express application instance
const app = {};

const VALID_TOKEN = "Bearer valid-jwt-token";
const MOCK_USER_RESPONSE = {
  id: 1,
  first_name: "Test",
  last_name: "User",
  login_email: "test@example.com",
  active: true,
  created_at: new Date(),
  updated_at: new Date(),
};
const MOCK_CREATE_INPUT: UserCreateInput = {
  first_name: "New",
  last_name: "Guy",
  login_email: "new@example.com",
  password: "SecurePassword123!",
};
const MOCK_UPDATE_INPUT: UserUpdateInput = { first_name: "UpdatedName" };

describe("UsersController", () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  // --- GET /users ---
  describe("GET /users", () => {
    it("should return 200 and a list of users", async () => {
      mockUserService.findAllUsers.mockResolvedValue([MOCK_USER_RESPONSE]);

      const response = await request(app)
        .get("/users")
        .set("Authorization", VALID_TOKEN);

      expect(response.statusCode).toBe(StatusCodes.OK);
      expect(Array.isArray(response.body)).toBe(true);
      expect(response.body.length).toBe(1);
    });
  });

  // --- GET /users/{id} ---
  describe("GET /users/:id", () => {
    it("should return 200 and the user details", async () => {
      mockUserService.findUserById.mockResolvedValue(MOCK_USER_RESPONSE);

      const response = await request(app)
        .get("/users/1")
        .set("Authorization", VALID_TOKEN);

      expect(response.statusCode).toBe(StatusCodes.OK);
      expect(response.body.id).toBe(1);
    });

    it("should return 404 Not Found if user does not exist", async () => {
      mockUserService.findUserById.mockResolvedValue(null);

      const response = await request(app)
        .get("/users/999")
        .set("Authorization", VALID_TOKEN);

      expect(response.statusCode).toBe(StatusCodes.NOT_FOUND);
    });
  });

  // --- POST /users ---
  describe("POST /users", () => {
    it("should return 201 and the created user", async () => {
      mockUserService.createUser.mockResolvedValue(MOCK_USER_RESPONSE as any);

      const response = await request(app)
        .post("/users")
        .send(MOCK_CREATE_INPUT);

      expect(response.statusCode).toBe(StatusCodes.CREATED);
      expect(mockUserService.createUser).toHaveBeenCalledWith(
        MOCK_CREATE_INPUT,
      );
    });

    it("should return 409 Conflict if email is already registered", async () => {
      mockUserService.createUser.mockRejectedValue(
        HttpErrors.Conflict("Email already exists"),
      );

      const response = await request(app)
        .post("/users")
        .send(MOCK_CREATE_INPUT);

      expect(response.statusCode).toBe(StatusCodes.CONFLICT);
    });
  });

  // --- PUT /users/{id} ---
  describe("PUT /users/:id", () => {
    it("should return 200 and the updated user", async () => {
      mockUserService.updateUser.mockResolvedValue({
        ...MOCK_USER_RESPONSE,
        ...MOCK_UPDATE_INPUT,
      } as any);

      const response = await request(app)
        .put("/users/1")
        .set("Authorization", VALID_TOKEN)
        .send(MOCK_UPDATE_INPUT);

      expect(response.statusCode).toBe(StatusCodes.OK);
      expect(response.body.first_name).toBe("UpdatedName");
    });

    it("should return 404 Not Found if user ID is invalid", async () => {
      // NOTE: Assuming the service throws a 404 ApiError on not found
      mockUserService.updateUser.mockRejectedValue(
        HttpErrors.NotFound("User not found"),
      );

      const response = await request(app)
        .put("/users/999")
        .set("Authorization", VALID_TOKEN)
        .send(MOCK_UPDATE_INPUT);

      expect(response.statusCode).toBe(StatusCodes.NOT_FOUND);
    });
  });

  // --- DELETE /users/{id} ---
  describe("DELETE /users/:id", () => {
    it("should return 204 No Content on successful deletion", async () => {
      mockUserService.deleteUser.mockResolvedValue(true as any); // Service returns success/void

      const response = await request(app)
        .delete("/users/1")
        .set("Authorization", VALID_TOKEN);

      expect(response.statusCode).toBe(StatusCodes.NO_CONTENT);
      expect(mockUserService.deleteUser).toHaveBeenCalledWith(1);
    });

    it("should return 404 Not Found if user does not exist", async () => {
      mockUserService.deleteUser.mockRejectedValue(
        HttpErrors.NotFound("User not found"),
      );

      const response = await request(app)
        .delete("/users/999")
        .set("Authorization", VALID_TOKEN);

      expect(response.statusCode).toBe(StatusCodes.NOT_FOUND);
    });
  });
});
