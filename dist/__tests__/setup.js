"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const test_utils_1 = require("./test-utils");
beforeAll(async () => {
    await test_utils_1.testUtils.startServer();
    await new Promise(resolve => setTimeout(resolve, 1000));
    await test_utils_1.testUtils.cleanupTestData();
}, 30000);
afterAll(async () => {
    await test_utils_1.testUtils.cleanupTestData();
    await test_utils_1.testUtils.stopServer();
}, 30000);
