import { prisma } from "@lib/prisma.js";
import { StatusCodes } from "http-status-codes";
import { ApiError } from "@errors/ApiError.js";
import { JWTUtils } from "@utils/jwt.utils.js";
import { PasswordUtils } from "@utils/password.utils.js";
import type { LoginResponse } from "types/auth/auth.dto.js";
import type { LoginInput } from "types/auth/auth.validation.js";

import * as cron from "node-cron";

export class AuthService {
  private static cleanupJob: cron.ScheduledTask;
  /**
   * Authenticate user and create session
   */
  static async login(
    loginData: LoginInput,
    ipAddress?: string,
    userAgent?: string,
  ): Promise<LoginResponse> {
    const { login_email, password } = loginData;
    const user = await prisma.users.findUnique({
      where: { login_email },
      select: {
        id: true,
        first_name: true,
        last_name: true,
        login_email: true,
        password_hash: true,
        contact_email: true,
        email_verified: true,
        active: true,
        created_at: true,
      },
    });

    if (!user) {
      throw new ApiError(StatusCodes.UNAUTHORIZED, "Invalid email or password");
    }

    if (!user.active) {
      throw new ApiError(StatusCodes.UNAUTHORIZED, "Account is deactivated");
    }

    if (!user.password_hash) {
      throw new ApiError(
        StatusCodes.UNAUTHORIZED,
        "No password set for this account",
      );
    }

    const isPasswordValid = await PasswordUtils.verifyPassword(
      password,
      user.password_hash,
    );
    if (!isPasswordValid) {
      throw new ApiError(StatusCodes.UNAUTHORIZED, "Invalid email or password");
    }

    const placeholderToken = "pending_" + Date.now();
    const sessionData: any = {
      user_id: user.id,
      token_hash: JWTUtils.hashToken(placeholderToken),
      expires_at: new Date(Date.now() + 7 * 24 * 60 * 60 * 1000),
    };

    if (ipAddress) sessionData.ip_address = ipAddress;
    if (userAgent) sessionData.user_agent = userAgent;

    const session = await prisma.user_sessions.create({
      data: sessionData,
    });

    const realToken = JWTUtils.generateToken({
      userId: user.id,
      email: user.login_email,
      sessionId: session.id,
    });

    await prisma.user_sessions.update({
      where: { id: session.id },
      data: {
        token_hash: JWTUtils.hashToken(realToken),
        expires_at: JWTUtils.getTokenExpiry(realToken),
      },
    });

    const { password_hash, ...userWithoutPassword } = user;

    return {
      user: userWithoutPassword,
      token: realToken,
      session_id: session.id,
      expires_at: session.expires_at,
      token_type: "Bearer",
    };
  }

  /**
   * Logout user by revoking session
   */
  static async logout(sessionId: number, userId: number): Promise<boolean> {
    const session = await prisma.user_sessions.updateMany({
      where: {
        id: sessionId,
        user_id: userId,
        revoked_at: null,
      },
      data: {
        revoked_at: new Date(),
      },
    });

    return session.count > 0;
  }

  /**
   * Logout all sessions for user
   */
  static async logoutAll(userId: number): Promise<number> {
    const result = await prisma.user_sessions.updateMany({
      where: {
        user_id: userId,
        revoked_at: null,
        expires_at: { gt: new Date() },
      },
      data: {
        revoked_at: new Date(),
      },
    });

    return result.count;
  }

  /**
   * Validate JWT token and return user
   */
  static async validateToken(
    token: string,
  ): Promise<{ user: any; session: any }> {
    try {
      JWTUtils.verifyToken(token);
      const tokenHash = JWTUtils.hashToken(token);

      const session = await prisma.user_sessions.findFirst({
        where: {
          token_hash: tokenHash,
          revoked_at: null,
          expires_at: { gt: new Date() },
        },
        include: {
          users: {
            select: {
              id: true,
              first_name: true,
              last_name: true,
              login_email: true,
              contact_email: true,
              email_verified: true,
              active: true,
            },
          },
        },
      });

      if (!session) {
        throw new ApiError(
          StatusCodes.UNAUTHORIZED,
          "Invalid or expired session",
        );
      }

      if (!session.users.active) {
        throw new ApiError(StatusCodes.UNAUTHORIZED, "Account is deactivated");
      }

      return {
        user: session.users,
        session: {
          id: session.id,
          ip_address: session.ip_address,
          created_at: session.created_at,
          expires_at: session.expires_at,
        },
      };
    } catch (error) {
      if (error instanceof ApiError) {
        throw error;
      }
      throw new ApiError(
        StatusCodes.UNAUTHORIZED,
        "Invalid authentication token",
      );
    }
  }

  /**
   * Get active sessions for user
   */
  static async getUserSessions(userId: number): Promise<any[]> {
    const sessions = await prisma.user_sessions.findMany({
      where: {
        user_id: userId,
        revoked_at: null,
        expires_at: { gt: new Date() },
      },
      select: {
        id: true,
        ip_address: true,
        user_agent: true,
        created_at: true,
        expires_at: true,
      },
      orderBy: {
        created_at: "desc",
      },
    });

    return sessions;
  }

  /**
   * Initialize session cleanup job
   */
  static initializeSessionCleanup(): void {
    this.cleanupJob = cron.schedule(
      "0 2 * * *",
      async () => {
        try {
          console.log("Starting session cleanup job...");
          const deletedCount = await this.cleanupExpiredSessions();
          console.log(
            `Session cleanup completed: ${deletedCount} sessions removed`,
          );
        } catch (error) {
          console.error("Session cleanup failed:", error);
        }
      },
      {},
    );

    this.cleanupJob.stop();
  }

  /**
   * Start the cleanup job
   */
  static startSessionCleanup(): void {
    if (this.cleanupJob) {
      this.cleanupJob.start();
      console.log("Session cleanup job started");
    }
  }

  /**
   * Stop the cleanup job
   */
  static stopSessionCleanup(): void {
    if (this.cleanupJob) {
      this.cleanupJob.stop();
      console.log("Session cleanup job stopped");
    }
  }

  /**
   * Clean up expired sessions
   */
  static async cleanupExpiredSessions(): Promise<number> {
    const result = await prisma.user_sessions.deleteMany({
      where: {
        OR: [{ expires_at: { lt: new Date() } }, { revoked_at: { not: null } }],
      },
    });

    return result.count;
  }

  /**
   * Manual cleanup (for testing or admin purposes)
   */
  static async manualCleanup(): Promise<{ deleted: number }> {
    const deleted = await this.cleanupExpiredSessions();
    return { deleted };
  }
}
