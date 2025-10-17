import request from "supertest";
import { StatusCodes } from "http-status-codes";
// NOTE: Using relative paths for imports that aren't path aliases
import { AuthService } from "../../src/services/auth/auth.service";
import { HttpErrors } from "../../src/errors/ApiError";
// NOTE: Assuming your Express app instance is exported as 'app' from src/index.ts
import { app } from "../../src/index";

// --- MOCK SETUP ---
// 1. Mock the entire service module using the alias path (critical for Jest)
jest.mock("@/services/auth/auth.service");

// 2. Type the imported service as a Mocked static module (correct way for static methods)
//    This allows us to call .mockResolvedValue directly on the static method names.
const mockAuthService = AuthService as jest.Mocked<typeof AuthService>;

// Mock Auth Data (rest of the data remains the same)
const VALID_TOKEN = "Bearer valid-jwt-token";
const MOCK_USER = {
  id: 1,
  first_name: "Auth",
  last_name: "User",
  login_email: "auth@test.com",
  active: true,
  email_verified: true,
};
const MOCK_SESSION = {
  id: 101,
  ip_address: "127.0.0.1",
  created_at: new Date(),
  expires_at: new Date(Date.now() + 3600000),
};
const MOCK_LOGIN_RESPONSE = {
  user: MOCK_USER,
  token: "new-jwt",
  session_id: 101,
  expires_at: new Date(),
  token_type: "Bearer",
};
const MOCK_LOGIN_INPUT = {
  login_email: "test@example.com",
  password: "Password123!",
};
const MOCK_REQUEST_AUTH = {
  user: { ...MOCK_USER, session: MOCK_SESSION },
  id: MOCK_USER.id,
};

describe("AuthController", () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  // --- POST /auth/login ---
  describe("POST /auth/login", () => {
    it("should return 200 and login response on success", async () => {
      // FIX: Call mockResolvedValue directly on the static method
      mockAuthService.login.mockResolvedValue(MOCK_LOGIN_RESPONSE as any);

      const response = await request(app)
        .post("/auth/login")
        .send(MOCK_LOGIN_INPUT);

      expect(response.statusCode).toBe(StatusCodes.OK);
      expect(response.body.token).toBe("new-jwt");
      expect(mockAuthService.login).toHaveBeenCalled();
    });

    it("should return 401 Unauthorized for invalid credentials", async () => {
      // FIX: Call mockRejectedValue directly on the static method
      mockAuthService.login.mockRejectedValue(
        HttpErrors.Unauthorized("Invalid credentials"),
      );

      const response = await request(app)
        .post("/auth/login")
        .send(MOCK_LOGIN_INPUT);

      expect(response.statusCode).toBe(StatusCodes.UNAUTHORIZED);
    });
  });

  // --- POST /auth/logout ---
  describe("POST /auth/logout", () => {
    it("should return 200 and revoke current session", async () => {
      // FIX: Call mockResolvedValue directly on the static method
      mockAuthService.logout.mockResolvedValue(true);

      const response = await request(app)
        .post("/auth/logout")
        .set("Authorization", VALID_TOKEN)
        .send()
        // Simulate authentication middleware injecting user/session
        .query({ user: MOCK_REQUEST_AUTH.user, id: MOCK_USER.id });

      expect(response.statusCode).toBe(StatusCodes.OK);
      expect(response.body.sessions_ended).toBe(1);
    });

    it("should return 401 if session data is missing", async () => {
      const response = await request(app)
        .post("/auth/logout")
        .set("Authorization", "Bearer invalid-token")
        .send()
        .query({ user: null });

      expect(response.statusCode).toBe(StatusCodes.UNAUTHORIZED);
    });
  });

  // --- GET /auth/validate ---
  describe("GET /auth/validate", () => {
    it("should return 200 and validation result if token is valid", async () => {
      const response = await request(app)
        .get("/auth/validate")
        .set("Authorization", VALID_TOKEN)
        .query({ user: MOCK_REQUEST_AUTH.user, id: MOCK_USER.id });

      expect(response.statusCode).toBe(StatusCodes.OK);
      expect(response.body.valid).toBe(true);
      expect(response.body.user.id).toBe(MOCK_USER.id);
    });

    it("should return 401 if token is invalid or user data is missing", async () => {
      const response = await request(app)
        .get("/auth/validate")
        .set("Authorization", "Bearer expired-token")
        .query({ user: null });

      expect(response.statusCode).toBe(StatusCodes.UNAUTHORIZED);
    });
  });
});
