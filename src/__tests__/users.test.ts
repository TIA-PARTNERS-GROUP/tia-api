import request from 'supertest';
import { testUtils } from './test-utils';

describe('Users API', () => {
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
    testUser = await testUtils.createTestUser();
    authToken = await testUtils.getAuthToken(testUser);
  });

  afterAll(async () => {
    await testUtils.cleanupTestData();
  });

  describe('GET /users', () => {
    it('should return all users', async () => {
      const response = await request(baseUrl)
        .get('/users')
        .set('Authorization', `Bearer ${authToken}`)
        .expect(200);

      expect(Array.isArray(response.body)).toBe(true);
    });

    it('should require authentication', async () => {
      const response = await request(baseUrl)
        .get('/users')
        .expect(401);

      expect(response.body).toHaveProperty('message');
    });
  });

  describe('GET /users/{id}', () => {
    it('should return user by ID', async () => {
      const response = await request(baseUrl)
        .get(`/users/${testUser.id}`)
        .set('Authorization', `Bearer ${authToken}`)
        .expect(200);

      expect(response.body).toHaveProperty('id', testUser.id);
      expect(response.body).toHaveProperty('login_email', testUser.login_email);
    });

    it('should return 404 for non-existent user', async () => {
      const response = await request(baseUrl)
        .get('/users/99999')
        .set('Authorization', `Bearer ${authToken}`)
        .expect(404);

      expect(response.body).toHaveProperty('message');
    });
  });

  describe('POST /users', () => {
    it('should create a new user', async () => {
      const newUser = {
        first_name: 'New',
        last_name: 'User',
        login_email: 'newuser@example.com',
        password: 'StrongPass123!',
        contact_email: 'newcontact@example.com'
      };

      const response = await request(baseUrl)
        .post('/users')
        .set('Authorization', `Bearer ${authToken}`)
        .send(newUser)
        .expect(201);

      expect(response.body).toHaveProperty('id');
      expect(response.body.first_name).toBe(newUser.first_name);
      expect(response.body.login_email).toBe(newUser.login_email);
    });

    it('should validate required fields', async () => {
      const invalidUser = {
        first_name: 'Invalid'
       
      };

      const response = await request(baseUrl)
        .post('/users')
        .set('Authorization', `Bearer ${authToken}`)
        .send(invalidUser)
        .expect(422);

      expect(response.body).toHaveProperty('details');
    });
  });

  describe('PUT /users/{id}', () => {
    it('should update user', async () => {
      const updateData = {
        first_name: 'Updated',
        contact_email: 'updated@example.com'
      };

      const response = await request(baseUrl)
        .put(`/users/${testUser.id}`)
        .set('Authorization', `Bearer ${authToken}`)
        .send(updateData)
        .expect(200);

      expect(response.body.first_name).toBe(updateData.first_name);
      expect(response.body.contact_email).toBe(updateData.contact_email);
    });
  });

  describe('DELETE /users/{id}', () => {
    it('should delete user', async () => {
      const response = await request(baseUrl)
        .delete(`/users/${testUser.id}`)
        .set('Authorization', `Bearer ${authToken}`)
        .expect(204);

     
      expect(response.body).toEqual({});
    });
  });
});
