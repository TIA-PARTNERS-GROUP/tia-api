import request from 'supertest';
import { testUtils } from './test-utils';

describe('Rate Limiting', () => {
  let baseUrl: string;

  beforeAll(async () => {
    baseUrl = testUtils.getBaseUrl();

   
    const isHealthy = await testUtils.healthCheck();
    if (!isHealthy) {
      throw new Error('Server is not healthy');
    }

    await testUtils.cleanupTestData();
    await testUtils.createLoginTestUser();
  });

  afterAll(async () => {
    await testUtils.cleanupTestData();
  });

  it('should rate limit login attempts', async () => {
    const attempts = 6;

    for (let i = 0; i < attempts; i++) {
      const response = await request(baseUrl)
        .post('/auth/login')
        .send({
          login_email: 'test@example.com',
          password: 'wrongpassword'
        });

      if (i < 5) {
        expect(response.status).toBe(401);
      } else {
        expect(response.status).toBe(429);
      }
    }
  }, 30000);
});
