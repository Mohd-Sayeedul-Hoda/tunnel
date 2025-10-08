import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { SignupForm } from "@/components/SignupForm";
import { useAuthSignup } from "@/hooks/use-auth";
import { type SignupFormData } from "@/lib/validations";
import { toast } from "sonner";

export const Route = createFileRoute("/signup")({
  component: SignupComponent,
});

function SignupComponent() {
  const navigate = useNavigate();
  const mutation = useAuthSignup();

  const onSubmit = (formData: SignupFormData) => {
    mutation.mutate(formData, {
      onSuccess: () => {
        toast.success(
          "Account created successfully! Please verify your email.",
        );
        navigate({
          to: "/verify-email",
          search: { email: formData.email },
        });
      },
      onError: (error: any) => {
        if (error.response) {
          const { status, data } = error.response;

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
                : "Signup failed. Please try again.";
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
        <SignupForm onSubmit={onSubmit} isLoading={mutation.isPending} />
      </div>
    </div>
  );
}
