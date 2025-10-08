import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { LoginForm } from "@/components/LoginForm";
import { useAuthLogin } from "@/hooks/use-auth";
import { toast } from "sonner";
import type { AxiosError } from "axios";
import type { LoginFormData } from "@/lib/validations";

export const Route = createFileRoute("/login")({
  component: LoginComponent,
});

function LoginComponent() {
  const navigate = useNavigate();
  const loginMutation = useAuthLogin();

  const onSubmit = (formData: LoginFormData) => {
    loginMutation.mutate(formData, {
      onSuccess: () => {
        toast.success("Logged in successfully");
        navigate({
          to: "/dashboard",
        });
      },
      onError: (error: Error) => {
        const axiosError = error as AxiosError;
        if (axiosError.response) {
          const { status, data } = axiosError.response;

          if (status === 422) {
            if (typeof data.error === 'object' && data.error !== null) {
              Object.entries(data.error).forEach(([field, errorMsg]) => {
                const fieldName = field.replace(/_/g, ' ').replace(/\b\w/g, l => l.toUpperCase());
                const message = Array.isArray(errorMsg) ? errorMsg[0] : errorMsg;
                toast.error(`${fieldName}: ${message}`);
              });
            } else {
              const errorMessages = Object.values(data.error) as string[];
              errorMessages.forEach((errorMsg) => {
                toast.error(errorMsg);
              });
            }
          } else {
            const errorMsg =
              typeof data.error === "string"
                ? data.error
                : "Login failed. Please try again.";
            toast.error(errorMsg);
          }
        } else {
          toast.error("Network error. Please check your connection.");
        }
      },
    });
  };

  return (
    <div className="flex min-h-svh w-full items-center justify-center p-6 md:p-10">
      <div className="w-full max-w-sm">
        <LoginForm onNext={onSubmit} isLoading={loginMutation.isPending} />
      </div>
    </div>
  );
}
