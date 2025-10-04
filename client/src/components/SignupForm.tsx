import React, { useState } from "react";

import { signupSchema, type SignupFormData } from "@/lib/validations";
import { cn } from "@/lib/utils";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "./ui/card";
import { Label } from "./ui/label";
import { Input } from "./ui/input";
import { Button } from "./ui/button";
import { ArrowLeftIcon, EyeIcon, EyeOffIcon } from "lucide-react";
import { Link } from "@tanstack/react-router";
import { useFormValidation } from "@/hooks/useFormValidation";

interface SignupFormProps extends React.ComponentProps<"div"> {
  onNext?: (data: SignupFormData) => void;
  isLoading?: boolean;
  error?: string | null;
}

export function SignupForm({ className, onNext }: SignupFormProps) {
  const [formData, setFormData] = useState<SignupFormData>({
    name: "",
    email: "",
    password: "",
    confirmPassword: "",
  });

  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);

  const { errors, touched, handleFieldChange, handleFieldBlur, validateForm } =
    useFormValidation({
      schema: signupSchema,
      validateOnBlur: true,
      validateOnChange: false,
    });

  const handleInputChange = (field: keyof SignupFormData, value: string) => {
    const updateData = { ...formData, [field]: value };
    setFormData(updateData);
    handleFieldChange(field, formData[field], updateData);
  };

  const handleBlur = (field: keyof SignupFormData) => {
    handleFieldBlur(field, formData[field], formData);
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    const validation = validateForm(formData);
    if (!validation.isValid) {
      console.log("Validation failed, errors:", validation.errors);
      return;
    }
    console.log("Signup data:", formData);
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
          <form onSubmit={handleSubmit}>
            <div className="flex flex-col gap-6">
              <div className="grid gap-3">
                <Label htmlFor="name">Full Name</Label>
                <Input
                  id="name"
                  type="text"
                  placeholder="Bruce Wayne"
                  value={formData.name}
                  onChange={(e) => handleInputChange("name", e.target.value)}
                  onBlur={() => handleBlur("name")}
                  className={cn(
                    "transition-colors",
                    errors.name &&
                    touched.name &&
                    "border-red-500 focus-visible:ring-red-500",
                  )}
                  required
                />
                {errors.name && touched.name && (
                  <p className="text-sm text-red-500 flex items-center gap-1">
                    <span className="text-red-500">.</span>
                    {errors.name}
                  </p>
                )}
              </div>

              <div className="grid gap-3">
                <Label htmlFor="email">Email</Label>
                <Input
                  id="email"
                  type="email"
                  value={formData.email}
                  onChange={(e) => handleInputChange("email", e.target.value)}
                  onBlur={() => handleBlur("email")}
                  placeholder="demo@demo.com"
                  className={cn(
                    "transition-colors",
                    errors.email &&
                    touched.email &&
                    "border-red-500 focus-visible:ring-red-500",
                  )}
                  required
                />
                {errors.email && touched.email && (
                  <p className="flex items-center gap-1 text-sm text-red-500">
                    <span className="text-red-500">.</span>
                    {errors.email}
                  </p>
                )}
              </div>

              <div className="grid gap-3">
                <Label htmlFor="password">Password</Label>
                <div className="relative">
                  <Input
                    id="password"
                    placeholder="password"
                    type={showPassword ? "text" : "password"}
                    onChange={(e) =>
                      handleInputChange("password", e.target.value)
                    }
                    onBlur={() => handleBlur("password")}
                    className={cn(
                      "pr-10 transition-colors",
                      errors.password &&
                      touched.password &&
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
                {errors.password && touched.password && (
                  <p className="text-sm text-red-500 flex items-center gap-1">
                    <span className="text-red-500">.</span>
                    {errors.password}
                  </p>
                )}
              </div>

              <div className="grid gap-3">
                <Label htmlFor="confirmPassword">Confirm Password</Label>
                <div className="relative">
                  <Input
                    id="confirmPassword"
                    type={showConfirmPassword ? "text" : "password"}
                    placeholder="Confirm your password"
                    onBlur={() => {
                      handleBlur("confirmPassword");
                    }}
                    onChange={(e) => {
                      handleInputChange("confirmPassword", e.target.value);
                    }}
                    className={cn(
                      "pr-10 transition-colors",
                      errors.confirmPassword &&
                      touched.confirmPassword &&
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
                {errors.confirmPassword && touched.confirmPassword && (
                  <p className="text-sm text-red-500 flex items-center gap-1">
                    <span className="text-red-500">.</span>
                    {errors.confirmPassword}
                  </p>
                )}
              </div>

              <div className="flex flex-col gap-3">
                <Button type="submit" className="w-full">
                  Create Account
                </Button>
              </div>
              <div className="flex gap-1 justify-center items-center">
                <p className="text-sm">Aleady have an account?</p>
                <Button asChild variant="link" className="p-0 h-auto text-sm">
                  <Link to="/login">Sign In</Link>
                </Button>
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
