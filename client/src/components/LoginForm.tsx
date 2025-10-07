import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Link } from "@tanstack/react-router";
import { ArrowLeftIcon } from "lucide-react";
import { useFormValidation } from "@/hooks/useFormValidation";
import { useState } from "react";
import { loginSchema, type LoginFormData } from "@/lib/validations";

interface LoginFormProps extends React.ComponentProps<"div"> {
  onNext: (data: LoginFormData) => void;
  isLoading?: boolean;
  error?: string | null;
}

export function LoginForm({ className, onNext, isLoading }: LoginFormProps) {
  const [formData, setFormData] = useState<LoginFormData>({
    email: "",
    password: "",
  });
  const { errors, touched, handleFieldChange, handleFieldBlur, validateForm } =
    useFormValidation({
      schema: loginSchema,
      validateOnBlur: true,
      validateOnChange: false,
    });

  const handleInputChange = (field: keyof LoginFormData, value: string) => {
    const newData = { ...formData, [field]: value };
    setFormData(newData);
    handleFieldChange(field, value, newData);
  };

  const handleBlur = (field: keyof LoginFormData) => {
    handleFieldBlur(field, formData[field], formData);
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    const validation = validateForm(formData);
    if (!validation.isValid) {
      console.log("Validation failed, errors:", validation.errors);
      return;
    }
    onNext({ email: formData.email, password: formData.password });
  };

  return (
    <div className={cn("flex flex-col gap-6", className)}>
      <Card>
        <CardHeader>
          <div>
            <CardTitle>Login to your account</CardTitle>
            <CardDescription>
              Enter your email and password below
            </CardDescription>
          </div>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit}>
            <div className="flex flex-col gap-6">
              <div className="grid gap-3">
                <Label htmlFor="email">Email</Label>
                <Input
                  id="email"
                  type="email"
                  placeholder="demo@demo.com"
                  value={formData.email}
                  onChange={(e) => handleInputChange("email", e.target.value)}
                  onBlur={() => handleBlur("email")}
                  className={cn(
                    "transition-colors",
                    errors.email &&
                    touched.email &&
                    "border-red-500 focus-visible:ring-red-500",
                  )}
                  required
                />
                {errors.email && touched["email"] && (
                  <p className="text-sm text-red-500 flex items-center gap-1">
                    <span className="text-red-500">•</span>
                    {errors.email}
                  </p>
                )}
              </div>
              <div className="grid gap-3">
                <div className="flex items-center">
                  <Label htmlFor="password">Password</Label>
                  <Link
                    to="/forgot-password"
                    className="ml-auto inline-block text-sm underline-offset-4 hover:underline"
                  >
                    Forgot password?
                  </Link>
                </div>
                <Input
                  id="password"
                  type="password"
                  value={formData.password}
                  onChange={(e) =>
                    handleInputChange("password", e.target.value)
                  }
                  onBlur={() => handleBlur("password")}
                  className={cn(
                    "transition-colors",
                    errors.password &&
                    "border-red-500 focus-visible:ring-red-500",
                  )}
                  required
                />
                {errors.password && (
                  <p className="text-sm text-red-500 flex items-center gap-1">
                    <span className="text-red-500">•</span>
                    {errors.password}
                  </p>
                )}
              </div>
              <div className="flex flex-col gap-3">
                <Button type="submit" className="w-full" disabled={isLoading}>
                  {isLoading ? "Logging in" : "Login"}
                </Button>
              </div>
            </div>
            <div className="text-center mt-4">
              <p className="text-sm text-muted-foreground">
                Don&apos;t have an account?{" "}
                <Button asChild variant="link" className="p-0 h-auto text-sm">
                  <Link to="/signup">Sign up</Link>
                </Button>
              </p>
            </div>
          </form>
        </CardContent>
      </Card>

      <div className="text-center">
        <Button variant="link" asChild>
          <Link to="/">
            <ArrowLeftIcon className="w-4 h-4" /> Back to home
          </Link>
        </Button>
      </div>
    </div>
  );
}
