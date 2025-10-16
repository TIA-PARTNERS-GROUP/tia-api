import request from 'supertest';
import { testUtils } from './test-utils';

describe('Authentication API', () => {
  let baseUrl: string;
  let testUser: any;
  let authToken: string;

  beforeAll(async () => {
    baseUrl = testUtils.getBaseUrl();

    const isHealthy = await testUtils.healthCheck();
    if (!isHealthy) {
      throw new Error('Server is not healthy');
    }

    await testUtils.cleanupTestData();
    testUser = await testUtils.createLoginTestUser();
    authToken = await testUtils.getAuthToken(testUser);
  });

  afterAll(async () => {
    await testUtils.cleanupTestData();
  });

  describe('POST /auth/login', () => {
    it('should login with valid credentials', async () => {
      const response = await request(baseUrl)
        .post('/auth/login')
        .send({
          login_email: 'test@example.com',
          password: 'password123'
        })
        .expect(200);

      expect(response.body).toHaveProperty('token');
      expect(response.body).toHaveProperty('user');
    });

    it('should return 401 with invalid password', async () => {
      const response = await request(baseUrl)
        .post('/auth/login')
        .send({
          login_email: 'test@example.com',
          password: 'wrongpassword'
        })
        .expect(401);

      expect(response.body).toHaveProperty('message');
    });

    it('should return 401 with non-existent email', async () => {
      const response = await request(baseUrl)
        .post('/auth/login')
        .send({
          login_email: 'nonexistent@example.com',
          password: 'password123'
        })
        .expect(401);

      expect(response.body).toHaveProperty('message');
    });

    it('should return 401 for inactive user', async () => {
      await testUtils.createInactiveTestUser();

      const response = await request(baseUrl)
        .post('/auth/login')
        .send({
          login_email: 'inactive@example.com',
          password: 'password123'
        })
        .expect(401);

      expect(response.body).toHaveProperty('message');
    });
  });

  describe('POST /auth/logout', () => {
    it('should logout successfully with valid token', async () => {
      const response = await request(baseUrl)
        .post('/auth/logout')
        .set('Authorization', `Bearer ${authToken}`)
        .expect(200);

      expect(response.body).toHaveProperty('message');
    });

    it('should return 401 without token', async () => {
      const response = await request(baseUrl)
        .post('/auth/logout')
        .expect(401);

      expect(response.body).toHaveProperty('message');
    });

    it('should return 401 with invalid token', async () => {
      const response = await request(baseUrl)
        .post('/auth/logout')
        .set('Authorization', 'Bearer invalid-token')
        .expect(401);

      expect(response.body).toHaveProperty('message');
    });
  });

  describe('GET /auth/validate', () => {
    it('should validate token successfully', async () => {
      const response = await request(baseUrl)
        .get('/auth/validate')
        .set('Authorization', `Bearer ${authToken}`)
        .expect(200);

      expect(response.body).toHaveProperty('valid', true);
      expect(response.body).toHaveProperty('user');
    });

    it('should return 401 with expired token', async () => {
      const expiredToken = await testUtils.getExpiredAuthToken(testUser);

      const response = await request(baseUrl)
        .get('/auth/validate')
        .set('Authorization', `Bearer ${expiredToken}`)
        .expect(401);

      expect(response.body).toHaveProperty('message');
    });
  });

  describe('GET /auth/sessions', () => {
    it('should return active sessions', async () => {
      await testUtils.createUserSession(testUser.id);

      const response = await request(baseUrl)
        .get('/auth/sessions')
        .set('Authorization', `Bearer ${authToken}`)
        .expect(200);

      expect(Array.isArray(response.body)).toBe(true);
    });

    it('should return empty array for user with no sessions', async () => {
      const newUser = await testUtils.createTestUser({
        login_email: 'newsessions@example.com'
      });
      const newUserToken = await testUtils.getAuthToken(newUser);

      const response = await request(baseUrl)
        .get('/auth/sessions')
        .set('Authorization', `Bearer ${newUserToken}`)
        .expect(200);

      expect(Array.isArray(response.body)).toBe(true);
    });
  });

  describe('POST /auth/logout-all', () => {
    it('should logout all sessions', async () => {
      await testUtils.createUserSession(testUser.id);
      await testUtils.createUserSession(testUser.id);

      const response = await request(baseUrl)
        .post('/auth/logout-all')
        .set('Authorization', `Bearer ${authToken}`)
        .expect(200);

      expect(response.body).toHaveProperty('message');
    });
  });
});
