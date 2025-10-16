"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const test_utils_1 = require("./test-utils");
beforeAll(async () => {
    // Start the server before all tests
    await test_utils_1.testUtils.startServer();
    // Wait a bit for server to be ready
    await new Promise(resolve => setTimeout(resolve, 1000));
    // Cleanup before tests
    await test_utils_1.testUtils.cleanupTestData();
}, 30000); // Increase timeout for server startup
afterAll(async () => {
    // Cleanup after all tests
    await test_utils_1.testUtils.cleanupTestData();
    // Stop the server after all tests
    await test_utils_1.testUtils.stopServer();
}, 30000);
