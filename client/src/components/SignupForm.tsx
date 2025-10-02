import React, { useCallback, useState } from "react";

import { signupSchema, type SignupFormData } from "@/lib/validations";
import { ZodError } from "zod";
import { cn } from "@/lib/utils";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "./ui/card";
import { Label } from "@radix-ui/react-label";
import { Input } from "./ui/input";
import { Button } from "./ui/button";
import { ArrowLeftIcon, EyeIcon, EyeOffIcon } from "lucide-react";
import { Link } from "@tanstack/react-router";

interface SignupFormProps extends React.ComponentProps<"div"> {
  onNext?: (data: SignupFormData) => void;
  isLoading?: boolean;
  error?: string | null;
}

const useFormValidation = (schema: typeof signupSchema) => {
  const validateField = useCallback(
    (field: keyof SignupFormData, value: string, allData?: SignupFormData) => {
      try {
        if (field === "confirmPassword" && allData) {
          schema.parse(allData);
        } else {
          schema.pick({ [field]: true }).parse({ [field]: value });
          return null;
        }
      } catch (error) {
        if (error instanceof ZodError) {
          const fieldError = error.issues.find((err) => err.path[0] === field);
          return fieldError?.message || null;
        }
        return null;
      }
    },
    [schema],
  );

  const validateForm = useCallback(
    (data: SignupFormData) => {
      try {
        schema.parse(data);
        return { isValid: true, errors: {} };
      } catch (error) {
        if (error instanceof ZodError) {
          const errors: Record<string, string> = {};
          error.issues.forEach((err) => {
            if (err.path[0]) {
              errors[err.path[0] as string] = err.message;
            }
          });
          return { isValid: false, errors };
        }
        return { isValid: false, errors: {} };
      }
    },
    [schema],
  );
  return { validateField, validateForm };
};

export function SignupForm({ className, onNext }: SignupFormProps) {
  const [formData, setFormData] = useState<SignupFormData>({
    name: "",
    email: "",
    password: "",
    confirmPassword: "",
  });

  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [touchedFields, setTouchedFields] = useState<Record<string, boolean>>(
    {},
  );

  const { validateField } = useFormValidation(signupSchema);

  const handleInputChange = (field: keyof SignupFormData, value: string) => {
    const updateData = { ...formData, [field]: value };
    setFormData(updateData);

    if (!touchedFields[field]) {
      setTouchedFields((prev) => ({ ...prev, [field]: true }));
    }

    if (errors[field]) {
      setErrors((prev) => ({ ...prev, [field]: "" }));
    }

    if (touchedFields[field]) {
      const fieldError = validateField(field, value, updateData);
      if (fieldError) {
        setErrors((prev) => ({ ...prev, [field]: fieldError }));
      }
    }
  };

  const handleFieldBlur = (field: keyof SignupFormData) => {
    setTouchedFields((prev) => ({ ...prev, [field]: true }));

    const fieldError = validateField(field, formData[field], formData);
    if (fieldError) {
      setErrors((prev) => ({ ...prev, [field]: fieldError }));
    }
  };

  return (
    <div className={cn("flex flex-col gap-6", className)}>
      <Card>
        <CardHeader>
          <div>
            <CardTitle>Create your account</CardTitle>
            <CardDescription>Enter you details to get started</CardDescription>
          </div>
        </CardHeader>
        <CardContent>
          <form>
            <div className="flex flex-col gap-6">
              <div className="grid gap-3">
                <Label htmlFor="name">Full Name</Label>
                <Input
                  id="name"
                  type="text"
                  placeholder="Bruce Wayne"
                  value={formData.name}
                  onChange={(e) => handleInputChange("name", e.target.value)}
                  onBlur={() => handleFieldBlur("name")}
                  className={cn(
                    "transition-colors",
                    errors.name &&
                    touchedFields.name &&
                    "border-red-500 focus-visible:ring-red-500",
                  )}
                  required
                />
                {errors.name && touchedFields.name && (
                  <p className="text-sm text-red-500 flex items-center gap-1">
                    <span className="text-red-500">.</span>
                    {errors.name}
                  </p>
                )}
              </div>

              <div className="grid gap-3">
                <Label id="email">Email</Label>
                <Input
                  id="email"
                  type="email"
                  onChange={(e) => handleInputChange("email", e.target.value)}
                  onBlur={() => handleFieldBlur("email")}
                  placeholder="enterprise@wayne.com"
                  className={cn(
                    "transition-colors",
                    errors.email &&
                    touchedFields.email &&
                    "border-red-500 focus-visible:ring-red-500",
                  )}
                  required
                />
                {errors.email && touchedFields.email && (
                  <p className="flex items-center gap-1 text-sm text-red-500">
                    <span className="text-red-500">.</span>
                    {errors.email}
                  </p>
                )}
              </div>

              <div className="grid gap-3">
                <Label id="password">Password</Label>
                <div className="relative">
                  <Input
                    id="password"
                    placeholder="password"
                    type={showPassword ? "text" : "password"}
                    onChange={(e) =>
                      handleInputChange("password", e.target.value)
                    }
                    onBlur={() => handleFieldBlur("password")}
                    className={cn(
                      "pr-10 transition-colors",
                      errors.password &&
                      touchedFields.password &&
                      "border-red-500 focus-visible:ring-red-500",
                    )}
                    required
                  />
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    className="absolute right-0 top-0 h-full px-3 py-2 hover:bg-transparent"
                    onClick={() => setShowPassword(!showPassword)}
                  >
                    {showPassword ? (
                      <EyeOffIcon className="h-4 w-4" />
                    ) : (
                      <EyeIcon className="w-4 h-4" />
                    )}
                  </Button>
                </div>
                {errors.password && touchedFields.password && (
                  <p className="text-sm text-red-500 flex items-center gap-1">
                    <span className="text-red-500">.</span>
                    {errors.password}
                  </p>
                )}
              </div>

              <div className="grid gap-3">
                <Label id="confirmPassword">Confirm Password</Label>
                <div className="relative">
                  <Input
                    id="confirmPassword"
                    type={showConfirmPassword ? "text" : "password"}
                    placeholder="Confirm your password"
                    onBlur={() => {
                      handleFieldBlur("confirmPassword");
                    }}
                    onChange={(e) => {
                      handleInputChange("confirmPassword", e.target.value);
                    }}
                    className={cn(
                      "pr-10 transition-colors",
                      errors.confirmPassword &&
                      touchedFields.confirmPassword &&
                      "border-red-500 focus-visible:ring-red-500",
                    )}
                    required
                  ></Input>
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                    className="absolute right-0 top-0 h-full px-3 py-2 hover:bg-transparent"
                  >
                    {showConfirmPassword ? (
                      <EyeOffIcon className="w-4 h-4" />
                    ) : (
                      <EyeIcon className="w-4 h-4" />
                    )}
                  </Button>
                </div>
              </div>

              <div className="flex flex-col gap-3">
                <Button type="submit" className="w-full">
                  Create Account
                </Button>
              </div>
              <div className="flex gap-1 justify-center items-center">
                <p className="text-sm">Aleady have an account?</p>
                <Link to="/login" className="text-sm underline">
                  Sign In
                </Link>
              </div>
            </div>
          </form>
        </CardContent>
      </Card>

      <div className="flex items-center justify-center">
        <Button variant="link" asChild>
          <Link to="/">
            <ArrowLeftIcon className="w-4 h-4 inline" /> Back to home
          </Link>
        </Button>
      </div>
    </div>
  );
}
