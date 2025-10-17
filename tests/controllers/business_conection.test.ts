import request from "supertest";
import { StatusCodes } from "http-status-codes";
import { BusinessConnectionService } from "@/services/business_connection/business_connection.service";
import { ApiError, HttpErrors } from "@/errors/ApiError";
import {
  InitiateConnectionInput,
  UpdateConnectionStatusInput,
} from "@/types/buisness_connection/business_connection.schema";

// --- MOCK SETUP ---
jest.mock("@/services/business_connection/business_connection.service");
const mockService = BusinessConnectionService as jest.MockedClass<
  typeof BusinessConnectionService
>;

// Placeholder for the Express application instance
const app = {};

const VALID_TOKEN = "Bearer VALID_TEST_TOKEN";

// Mock Data
const MOCK_CONNECTION_RESPONSE = {
  id: 1,
  initiating_business_id: 101,
  receiving_business_id: 202,
  connection_type: "Partnership",
  status: "pending",
  initiated_by_user_id: 5,
  created_at: new Date(),
  initiating_business: {
    id: 101,
    name: "Alpha Inc.",
    business_type: "Technology",
  },
  receiving_business: {
    id: 202,
    name: "Beta Solutions",
    business_type: "Consulting",
  },
};

const MOCK_INITIATE_INPUT: InitiateConnectionInput = {
  initiating_business_id: 101,
  receiving_business_id: 202,
  connection_type: "Partnership",
  initiated_by_user_id: 5,
};

describe("BusinessConnectionController", () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  // --- POST /connections (Initiate Connection) ---
  describe("POST /connections", () => {
    it("should return 201 and the new connection response on success", async () => {
      mockService.prototype.initiateConnection.mockResolvedValue(
        MOCK_CONNECTION_RESPONSE as any,
      );

      const response = await request(app)
        .post("/connections")
        .set("Authorization", VALID_TOKEN)
        .send(MOCK_INITIATE_INPUT);

      expect(response.statusCode).toBe(StatusCodes.CREATED);
      expect(response.body.status).toBe("pending");
    });

    it("should return 409 Conflict if connection already exists", async () => {
      mockService.prototype.initiateConnection.mockRejectedValue(
        new ApiError(StatusCodes.CONFLICT, "Connection already exists"),
      );

      const response = await request(app)
        .post("/connections")
        .set("Authorization", VALID_TOKEN)
        .send(MOCK_INITIATE_INPUT);

      expect(response.statusCode).toBe(StatusCodes.CONFLICT);
    });
  });

  // --- GET /connections/record/{connectionId} (Get by ID) ---
  describe("GET /connections/record/:connectionId", () => {
    it("should return 200 and the connection data if found", async () => {
      mockService.prototype.getConnectionById.mockResolvedValue(
        MOCK_CONNECTION_RESPONSE as any,
      );

      const response = await request(app)
        .get("/connections/record/1")
        .set("Authorization", VALID_TOKEN);

      expect(response.statusCode).toBe(StatusCodes.OK);
      expect(response.body.id).toBe(1);
    });

    it("should return 404 Not Found if connection does not exist", async () => {
      mockService.prototype.getConnectionById.mockResolvedValue(null);

      const response = await request(app)
        .get("/connections/record/999")
        .set("Authorization", VALID_TOKEN);

      expect(response.statusCode).toBe(StatusCodes.NOT_FOUND);
    });
  });

  // --- PUT /connections/{connectionId}/status (Update Status) ---
  describe("PUT /connections/:connectionId/status", () => {
    const updateBody: UpdateConnectionStatusInput = { status: "active" };
    const mockUpdatedResponse = {
      ...MOCK_CONNECTION_RESPONSE,
      status: "active",
    };

    it("should return 200 and the updated connection on success", async () => {
      mockService.prototype.updateConnectionStatus.mockResolvedValue(
        mockUpdatedResponse as any,
      );

      const response = await request(app)
        .put("/connections/1/status")
        .set("Authorization", VALID_TOKEN)
        .send(updateBody);

      expect(response.statusCode).toBe(StatusCodes.OK);
      expect(response.body.status).toBe("active");
    });

    it("should return 404 Not Found if connection ID is invalid", async () => {
      mockService.prototype.updateConnectionStatus.mockRejectedValue(
        new ApiError(StatusCodes.NOT_FOUND, "Business connection not found."),
      );

      const response = await request(app)
        .put("/connections/999/status")
        .set("Authorization", VALID_TOKEN)
        .send(updateBody);

      expect(response.statusCode).toBe(StatusCodes.NOT_FOUND);
    });
  });

  // --- GET /connections/{businessId} (Get Connections for Business) ---
  describe("GET /connections/:businessId", () => {
    const mockSummaryList = [
      {
        id: 1,
        status: "active",
        connection_type: "Client",
        initiating_business_id: 100,
        receiving_business_id: 1,
        created_at: new Date(),
      },
    ];

    it("should return 200 and a list of connection summaries", async () => {
      mockService.prototype.getConnectionsForBusiness.mockResolvedValue(
        mockSummaryList as any,
      );

      const response = await request(app)
        .get("/connections/100")
        .set("Authorization", VALID_TOKEN);

      expect(response.statusCode).toBe(StatusCodes.OK);
      expect(Array.isArray(response.body)).toBe(true);
    });

    it("should return 400 Bad Request for invalid businessId (<= 0)", async () => {
      const response = await request(app)
        .get("/connections/0")
        .set("Authorization", VALID_TOKEN);

      expect(response.statusCode).toBe(StatusCodes.BAD_REQUEST);
      expect(response.body.message).toContain("Invalid Business ID");
    });
  });
});
