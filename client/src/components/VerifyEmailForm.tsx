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
import { ArrowLeftIcon, MailIcon } from "lucide-react";
import { useState, useEffect } from "react";
import {
  emailValidationSchema,
  type EmailValidationData,
} from "@/lib/validations";
import { useFormValidation } from "@/hooks/useFormValidation";

interface VerifyEmailFormProps extends React.ComponentProps<"div"> {
  email: string;
  onVerify: (otp: string) => void;
  onResend: () => void;
  isLoading?: boolean;
  error?: string;
}

export function VerifyEmailForm({
  email,
  className,
  ...props
}: VerifyEmailFormProps) {
  const [formData, setFormData] = useState<EmailValidationData>({
    email: email || "",
    otp: "",
  });
  const [resendCooldown, setResendCooldown] = useState(0);

  useEffect(() => {
    if (resendCooldown > 0) {
      const timer = setTimeout(() => {
        setResendCooldown(resendCooldown - 1);
      }, 1000);
      return () => clearTimeout(timer);
    }
  }, [resendCooldown]);

  const isEmailProvided = email && email.trim() !== "";

  const { errors, touched, handleFieldChange, handleFieldBlur, validateForm } =
    useFormValidation({
      schema: emailValidationSchema,
      validateOnChange: true,
      validateOnBlur: true,
    });

  const handleInput = (field: keyof EmailValidationData, value: string) => {
    const newData = { ...formData, [field]: value };
    setFormData(newData);
    handleFieldChange(field, value, newData);
  };

  const handleBlur = (field: keyof EmailValidationData) => {
    handleFieldBlur(field, formData[field], formData);
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    const validation = validateForm(formData);
    if (!validation.isValid) {
      console.log("Validation failed, errors:", validation.errors);
      return;
    }
    console.log("EmailVerification data:", formData);
  };

  const handleResend = () => {
    console.log("resend the mail");
  };

  // Show error if email is not provided
  if (!isEmailProvided) {
    return (
      <div className={cn("flex flex-col gap-6", className)} {...props}>
        <Card>
          <CardHeader>
            <div className="text-center">
              <div className="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-red-100 dark:bg-red-900/20">
                <MailIcon className="h-6 w-6 text-red-500" />
              </div>
              <CardTitle className="text-red-600 dark:text-red-400">
                Email Not Set
              </CardTitle>
              <CardDescription>
                No email address provided for verification
              </CardDescription>
            </div>
          </CardHeader>
          <CardContent>
            <div className="text-center">
              <p className="text-sm text-muted-foreground mb-4">
                Please go back to the signup page and provide a valid email
                address.
              </p>
              <Button variant="link" asChild>
                <Link to="/signup">
                  <ArrowLeftIcon className="w-4 h-4" /> Back to signup
                </Link>
              </Button>
            </div>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className={cn("flex flex-col gap-6", className)} {...props}>
      <Card>
        <CardHeader>
          <div className="text-center">
            <div className="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-primary/10">
              <MailIcon className="h-6 w-6 text-primary" />
            </div>
            <CardTitle>Verify your email</CardTitle>
            <CardDescription>
              We've sent a 6-digit verification code to
            </CardDescription>
            <p className="font-medium text-foreground">{email}</p>
          </div>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit}>
            <div className="flex flex-col gap-6">
              <div className="grid gap-3">
                <Label id="email">Email</Label>
                <Input
                  id="email"
                  type="email"
                  placeholder="demo@demo.com"
                  value={formData.email}
                  onChange={(e) => handleInput("email", e.target.value)}
                  disabled
                  onBlur={() => {
                    handleBlur("email");
                  }}
                  className={cn(
                    "transition-colors",
                    errors.email &&
                    touched.email &&
                    "border-red-500 focus-visible:ring-red-500",
                  )}
                  required
                />
                {errors.email ||
                  (touched.email && (
                    <p className="text-sm text-red-500 text-center">
                      <span className="text-red-500">.</span>
                      {errors.email}
                    </p>
                  ))}
              </div>

              <div className="grid gap-3">
                <Label htmlFor="otp">Otp</Label>
                <Input
                  id="otp"
                  value={formData.otp}
                  type="text"
                  maxLength={6}
                  placeholder="000000"
                  onChange={(e) => {
                    handleInput("otp", e.target.value);
                  }}
                  onBlur={() => handleBlur("otp")}
                  className={cn(
                    "transition-colors",
                    errors.otp &&
                    touched.otp &&
                    "border-red-500 focus-visible:ring-red-500",
                  )}
                  required
                />
                {errors.otp && touched.otp && (
                  <p className="text-sm text-red-500 flex items-center">
                    <span className="text-red-500">.</span>
                    {errors.otp}
                  </p>
                )}
              </div>

              <div className="flex flex-col gap-3">
                <Button type="submit" className="w-full">
                  Verify Email
                </Button>
              </div>

              <div className="text-center">
                <p className="text-sm text-muted-foreground">
                  Didn't receive the code?{" "}
                  <Button
                    type="button"
                    variant="link"
                    className="p-0 h-auto text-sm"
                    onClick={() => handleResend()}
                  >
                    {resendCooldown > 0
                      ? `Resend in ${resendCooldown}s`
                      : "Resend code"}
                  </Button>
                </p>
              </div>
            </div>
          </form>
        </CardContent>
      </Card>

      <div className="text-center">
        <Button variant="link" asChild>
          <Link to="/signup">
            <ArrowLeftIcon className="w-4 h-4" /> Back to signup
          </Link>
        </Button>
      </div>
    </div>
  );
}
