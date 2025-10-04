import { cn } from "@/lib/utils";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "./ui/card";
import { Input } from "./ui/input";
import { Label } from "./ui/label";
import { Button } from "./ui/button";
import { Link } from "@tanstack/react-router";
import { ArrowLeftIcon } from "lucide-react";
import { useState } from "react";
import { type EmailFormData, emailSchema } from "@/lib/validations";
import { useFormValidation } from "@/hooks/useFormValidation";

interface ForgotPasswordProps extends React.ComponentProps<"div"> {
  onNext?: (data: EmailFormData) => void;
  isLoading?: boolean;
  error?: string | null;
}

export function ForgotPasswordForm({ className, onNext }: ForgotPasswordProps) {
  const [formData, setFormData] = useState<EmailFormData>({
    email: "",
  });
  const { errors, touched, handleFieldChange, handleFieldBlur, validateForm } =
    useFormValidation({
      schema: emailSchema,
      validateOnBlur: true,
      validateOnChange: true,
    });

  const handleInputChange = (field: keyof EmailFormData, value: string) => {
    const newData = { ...formData, [field]: value };
    setFormData(newData);
    handleFieldChange("email", value, newData);
  };

  const handleBlur = (field: keyof EmailFormData) => {
    handleFieldBlur(field, formData[field], formData);
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    const validation = validateForm(formData);
    if (!validation.isValid) {
      console.log("Validation failed, errors:", validation.errors);
      return;
    }
    console.log("email: ", formData);
  };

  return (
    <div className={cn("flex flex-col gap-6", className)}>
      <Card>
        <CardHeader>
          <CardTitle>Reset Password</CardTitle>
          <CardDescription>
            Enter you email to reset you password
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit}>
            <div className="flex flex-col gap-6">
              <div className="grid gap-3">
                <Label htmlFor="email">Email</Label>
                <Input
                  id="email"
                  placeholder="demo@demo.com"
                  value={formData.email}
                  type="email"
                  onChange={(e) => handleInputChange("email", e.target.value)}
                  onBlur={() => handleBlur("email")}
                  className={cn(
                    "transition-colors",
                    errors.email &&
                    "border-red-500 focus-visible:border-red-500",
                  )}
                  required
                />
                {errors.email && touched.email && (
                  <p className="text-sm text-red-500 flex items-center, gap-1">
                    {errors.email}
                  </p>
                )}
              </div>
              <Button type="submit">Reset Password</Button>
            </div>
          </form>
        </CardContent>
      </Card>

      <div className="flex items-center justify-center">
        <Button variant="link" type="submit" asChild>
          <Link to="/">
            <ArrowLeftIcon className="w-4 h-4 inline" /> Back to home
          </Link>
        </Button>
      </div>
    </div>
  );
}
