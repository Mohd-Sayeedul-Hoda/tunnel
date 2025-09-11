import { z } from 'zod'

// Signup form validation schema
export const signupSchema = z.object({
  name: z
    .string()
    .min(1, "Name is required")
    .min(2, "Name must be at least 2 characters")
    .max(50, "Name must be less than 50 characters")
    .regex(/^[a-zA-Z\s]+$/, "Name can only contain letters and spaces"),
  
  email: z
    .string()
    .min(1, "Email is required")
    .email("Please enter a valid email address")
    .max(100, "Email must be less than 100 characters"),
  
  password: z
    .string()
    .min(1, "Password is required")
    .min(8, "Password must be at least 8 characters")
    .max(100, "Password must be less than 100 characters")
    .regex(/^(?=.*[a-z])/, "Password must contain at least one lowercase letter")
    .regex(/^(?=.*[A-Z])/, "Password must contain at least one uppercase letter")
    .regex(/^(?=.*\d)/, "Password must contain at least one number")
    .regex(/^(?=.*[@$!%*?&])/, "Password must contain at least one special character (@$!%*?&)"),
  
  confirmPassword: z
    .string()
    .min(1, "Please confirm your password")
})

// Refine to check if passwords match
export const signupFormSchema = signupSchema.refine(
  (data) => data.password === data.confirmPassword,
  {
    message: "Passwords do not match",
    path: ["confirmPassword"], // This will show the error on the confirmPassword field
  }
)

// Login form validation schema
export const loginSchema = z.object({
  email: z
    .string()
    .min(1, "Email is required")
    .email("Please enter a valid email address"),
  
  password: z
    .string()
    .min(1, "Password is required")
})

// OTP verification schema
export const otpSchema = z.object({
  otp: z
    .string()
    .min(1, "Verification code is required")
    .length(6, "Verification code must be 6 digits")
    .regex(/^\d{6}$/, "Verification code must contain only numbers")
})

// Type exports for TypeScript
export type SignupFormData = z.infer<typeof signupFormSchema>
export type LoginFormData = z.infer<typeof loginSchema>
export type OtpFormData = z.infer<typeof otpSchema>

