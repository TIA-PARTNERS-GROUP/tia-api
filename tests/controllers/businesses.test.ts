import request from "supertest";
import { StatusCodes } from "http-status-codes";
import { BusinessesService } from "@/services/businesses/businesses.service";
import { HttpErrors } from "@/errors/ApiError";
import {
  CreateBusinessInput,
  UpdateBusinessInput,
} from "@/types/businesses/businesses.validation";

// --- MOCK SETUP ---
jest.mock("@/services/businesses/businesses.service");
const mockBusinessesService = BusinessesService as jest.MockedClass<
  typeof BusinessesService
>;

// Placeholder for the Express application instance
const app = {};

const VALID_TOKEN = "Bearer valid-jwt-token";
const MOCK_BUSINESS_RESPONSE = {
  id: 1,
  name: "Tech Corp",
  business_type: "Technology",
  operator_user_id: 5,
  active: 1,
  created_at: new Date(),
  updated_at: new Date(),
};
const MOCK_CREATE_INPUT: CreateBusinessInput = {
  operator_user_id: 5,
  name: "New Business",
  business_type: "Consulting",
  business_category: "B2B",
  business_phase: "Startup",
};
const MOCK_UPDATE_INPUT: UpdateBusinessInput = { tagline: "New Tagline" };
const MOCK_STATS_RESPONSE = { total_projects: 10, total_connections: 5 };

describe("BusinessesController", () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  // --- GET /businesses (Filtered List) ---
  describe("GET /businesses", () => {
    it("should return 200 and a list of businesses", async () => {
      mockBusinessesService.prototype.getBusinesses.mockResolvedValue([
        MOCK_BUSINESS_RESPONSE,
      ]);

      const response = await request(app)
        .get("/businesses?business_type=Technology&active=true")
        .set("Authorization", VALID_TOKEN);

      expect(response.statusCode).toBe(StatusCodes.OK);
      expect(Array.isArray(response.body)).toBe(true);
      expect(
        mockBusinessesService.prototype.getBusinesses,
      ).toHaveBeenCalledWith(
        expect.objectContaining({
          business_type: "Technology",
          active: true,
        }),
      );
    });
  });

  // --- GET /businesses/{businessId} ---
  describe("GET /businesses/:businessId", () => {
    it("should return 200 and the business details", async () => {
      mockBusinessesService.prototype.getBusinessById.mockResolvedValue(
        MOCK_BUSINESS_RESPONSE as any,
      );

      const response = await request(app)
        .get("/businesses/1")
        .set("Authorization", VALID_TOKEN);

      expect(response.statusCode).toBe(StatusCodes.OK);
    });

    it("should return 404 Not Found", async () => {
      mockBusinessesService.prototype.getBusinessById.mockRejectedValue(
        HttpErrors.NotFound("Business not found"),
      );

      const response = await request(app)
        .get("/businesses/999")
        .set("Authorization", VALID_TOKEN);

      expect(response.statusCode).toBe(StatusCodes.NOT_FOUND);
    });
  });

  // --- POST /businesses ---
  describe("POST /businesses", () => {
    it("should return 201 and the created business", async () => {
      mockBusinessesService.prototype.createBusiness.mockResolvedValue(
        MOCK_BUSINESS_RESPONSE as any,
      );

      const response = await request(app)
        .post("/businesses")
        .set("Authorization", VALID_TOKEN)
        .send(MOCK_CREATE_INPUT);

      expect(response.statusCode).toBe(StatusCodes.CREATED);
    });

    it("should return 422 for invalid input (Zod/tsoa validation)", async () => {
      // Missing required field 'name'
      const invalidInput = { ...MOCK_CREATE_INPUT, name: "" };

      const response = await request(app)
        .post("/businesses")
        .set("Authorization", VALID_TOKEN)
        .send(invalidInput);

      // Expecting tsoa/Express validation middleware to catch this
      expect(response.statusCode).toBe(StatusCodes.UNPROCESSABLE_ENTITY);
    });
  });

  // --- PUT /businesses/{businessId} ---
  describe("PUT /businesses/:businessId", () => {
    it("should return 200 and the updated business", async () => {
      mockBusinessesService.prototype.updateBusiness.mockResolvedValue({
        ...MOCK_BUSINESS_RESPONSE,
        ...MOCK_UPDATE_INPUT,
      } as any);

      const response = await request(app)
        .put("/businesses/1")
        .set("Authorization", VALID_TOKEN)
        .send(MOCK_UPDATE_INPUT);

      expect(response.statusCode).toBe(StatusCodes.OK);
      expect(response.body.tagline).toBe("New Tagline");
    });
  });

  // --- DELETE /businesses/{businessId} ---
  describe("DELETE /businesses/:businessId", () => {
    it("should return 204 No Content", async () => {
      mockBusinessesService.prototype.deleteBusiness.mockResolvedValue(
        undefined,
      );

      const response = await request(app)
        .delete("/businesses/1")
        .set("Authorization", VALID_TOKEN);

      expect(response.statusCode).toBe(StatusCodes.NO_CONTENT);
    });
  });

  // --- GET /businesses/{businessId}/stats ---
  describe("GET /businesses/:businessId/stats", () => {
    it("should return 200 and business statistics", async () => {
      mockBusinessesService.prototype.getBusinessStats.mockResolvedValue(
        MOCK_STATS_RESPONSE as any,
      );

      const response = await request(app)
        .get("/businesses/1/stats")
        .set("Authorization", VALID_TOKEN);

      expect(response.statusCode).toBe(StatusCodes.OK);
      expect(response.body.total_projects).toBe(10);
    });
  });
});
