import { testUtils } from './test-utils';

beforeAll(async () => {
 
  await testUtils.startServer();

 
  await new Promise(resolve => setTimeout(resolve, 1000));

 
  await testUtils.cleanupTestData();
}, 30000);

afterAll(async () => {
 
  await testUtils.cleanupTestData();

 
  await testUtils.stopServer();
}, 30000);
