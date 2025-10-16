import { testUtils } from './test-utils';

beforeAll(async () => {
  // Start the server before all tests
  await testUtils.startServer();

  // Wait a bit for server to be ready
  await new Promise(resolve => setTimeout(resolve, 1000));

  // Cleanup before tests
  await testUtils.cleanupTestData();
}, 30000); // Increase timeout for server startup

afterAll(async () => {
  // Cleanup after all tests
  await testUtils.cleanupTestData();

  // Stop the server after all tests
  await testUtils.stopServer();
}, 30000);
