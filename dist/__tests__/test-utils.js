"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.testUtils = void 0;
const prisma_1 = require("../lib/prisma");
const jsonwebtoken_1 = __importDefault(require("jsonwebtoken"));
const bcrypt_1 = __importDefault(require("bcrypt"));
const crypto_1 = __importDefault(require("crypto"));
const index_1 = require("../index"); // Import your Express app
class TestUtils {
    server;
    baseUrl;
    constructor() {
        this.baseUrl = `http://localhost:${process.env.PORT || 8000}`;
    }
    /**
     * Start test server
     */
    async startServer() {
        return new Promise((resolve) => {
            const port = process.env.TEST_PORT || 8001; // Use different port for tests
            this.server = index_1.app.listen(port, () => {
                this.baseUrl = `http://localhost:${port}`;
                console.log(`Test server running on ${this.baseUrl}`);
                resolve();
            });
        });
    }
    /**
     * Stop test server
     */
    async stopServer() {
        return new Promise((resolve, reject) => {
            if (this.server) {
                this.server.close((err) => {
                    if (err) {
                        reject(err);
                    }
                    else {
                        console.log('Test server stopped');
                        resolve();
                    }
                });
            }
            else {
                resolve();
            }
        });
    }
    /**
     * Get base URL for tests
     */
    getBaseUrl() {
        return this.baseUrl;
    }
    /**
     * Clean up all test data
     */
    async cleanupTestData() {
        try {
            await prisma_1.prisma.user_sessions.deleteMany({});
            await prisma_1.prisma.users.deleteMany({
                where: {
                    OR: [
                        { login_email: { contains: 'test@test.com' } },
                        { login_email: { contains: 'test' } },
                        { login_email: 'test@example.com' },
                        { login_email: 'inactive@example.com' },
                        { login_email: 'newuser@example.com' },
                        { login_email: { contains: 'newsessions@example.com' } }
                    ]
                }
            });
        }
        catch (error) {
            console.error('Cleanup error:', error);
        }
    }
    /**
     * Create a test user with unique email
     */
    async createTestUser(userData) {
        const timestamp = Date.now();
        const random = Math.random().toString(36).substring(7);
        return await prisma_1.prisma.users.create({
            data: {
                first_name: 'Test',
                last_name: 'User',
                login_email: `test${timestamp}${random}@test.com`,
                contact_email: `test${timestamp}${random}@test.com`,
                password_hash: 'hashed_password',
                active: true,
                ...userData
            }
        });
    }
    /**
     * Create a test user with properly hashed password
     */
    async createTestUserWithHashedPassword(userData) {
        const timestamp = Date.now();
        const random = Math.random().toString(36).substring(7);
        const password = userData?.password || 'testpassword123';
        const password_hash = await bcrypt_1.default.hash(password, 10);
        return await prisma_1.prisma.users.create({
            data: {
                first_name: 'Test',
                last_name: 'User',
                login_email: `test${timestamp}${random}@test.com`,
                contact_email: `test${timestamp}${random}@test.com`,
                password_hash,
                active: true,
                ...userData,
                password: undefined // Remove plain text password if provided
            }
        });
    }
    /**
     * Create test user with specific email for login tests
     */
    async createLoginTestUser() {
        const password = 'password123';
        const password_hash = await bcrypt_1.default.hash(password, 10);
        return await prisma_1.prisma.users.create({
            data: {
                first_name: 'Login',
                last_name: 'Test',
                login_email: 'test@example.com',
                contact_email: 'test@example.com',
                password_hash,
                active: true
            }
        });
    }
    /**
     * Create inactive test user
     */
    async createInactiveTestUser() {
        const password = 'password123';
        const password_hash = await bcrypt_1.default.hash(password, 10);
        return await prisma_1.prisma.users.create({
            data: {
                first_name: 'Inactive',
                last_name: 'User',
                login_email: 'inactive@example.com',
                contact_email: 'inactive@example.com',
                password_hash,
                active: false
            }
        });
    }
    /**
     * Generate JWT token for a user
     */
    async getAuthToken(user) {
        return jsonwebtoken_1.default.sign({
            userId: user.id,
            email: user.login_email
        }, process.env.JWT_SECRET || 'test-secret', { expiresIn: '1h' });
    }
    /**
     * Generate expired JWT token for testing
     */
    async getExpiredAuthToken(user) {
        return jsonwebtoken_1.default.sign({
            userId: user.id,
            email: user.login_email
        }, process.env.JWT_SECRET || 'test-secret', { expiresIn: '-1h' } // Already expired
        );
    }
    /**
     * Create test user and return with token
     */
    async createTestUserWithToken(userData) {
        const user = await this.createTestUser(userData);
        const token = await this.getAuthToken(user);
        return { user, token };
    }
    /**
     * Create a user session
     */
    async createUserSession(userId, sessionData) {
        const token = `session-token-${Date.now()}-${Math.random().toString(36).substring(7)}`;
        const token_hash = crypto_1.default.createHash('sha256').update(token).digest('hex');
        return await prisma_1.prisma.user_sessions.create({
            data: {
                user_id: userId,
                token,
                token_hash,
                expires_at: new Date(Date.now() + 24 * 60 * 60 * 1000), // 24 hours from now
                ...sessionData
            }
        });
    }
    /**
     * Health check - verify server is responding
     */
    async healthCheck() {
        try {
            const response = await fetch(`${this.baseUrl}/health`);
            return response.status === 200;
        }
        catch (error) {
            return false;
        }
    }
}
exports.testUtils = new TestUtils();
