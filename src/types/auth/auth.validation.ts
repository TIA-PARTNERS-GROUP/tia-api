import { z } from "zod";

export const loginSchema = z.object({
  login_email: z.string().email("Please provide a valid email address"),
  password: z.string().min(1, "Password is required"),
});

export const passwordStrengthSchema = z
  .string()
  .min(8, "Password must be at least 8 characters long")
  .regex(/^(?=.*[a-z])/, "Password must contain at least one lowercase letter")
  .regex(/^(?=.*[A-Z])/, "Password must contain at least one uppercase letter")
  .regex(/^(?=.*\d)/, "Password must contain at least one number")
  .regex(
    /^(?=.*[@$!%*?&])/,
    "Password must contain at least one special character (@$!%*?&)",
  );

export type LoginInput = z.infer<typeof loginSchema>;
