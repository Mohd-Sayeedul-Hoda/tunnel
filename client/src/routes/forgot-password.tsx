import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { ForgotPasswordForm } from "@/components/ForgotPasswordForm";
import {
  useSendForgotPasswordEmail,
  useVerifyForgotPasswordOtp,
} from "@/hooks/use-auth";
import { toast } from "sonner";
import { useEffect, useState } from "react";
import { email } from "zod";

export const Route = createFileRoute("/forgot-password")({
  component: RouteComponent,
});

function RouteComponent() {
  const navigate = useNavigate();
  const sendForgotEmailMutation = useSendForgotPasswordEmail();
  const verifyForgotEmailOtpMutation = useVerifyForgotPasswordOtp();

  const [resendCooldown, setResendCooldown] = useState(0);

  useEffect(() => {
    if (resendCooldown > 0) {
      const timer = setTimeout(() => {
        setResendCooldown(resendCooldown - 1);
      }, 1000);
      return () => clearTimeout(timer);
    }
  }, [resendCooldown]);

  const onSendEmail = (email: string) => {
    if (email) {
      sendForgotEmailMutation.mutate(email, {
        onSuccess: () => {
          toast.success("Password Reset Successfully");
          navigate({
            to: "/login",
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
                  : "Failed to send password reset email. Please try again";
              toast.error(errorMsg);
            }
          } else {
            toast.error("Network error. Please check your connection.");
          }
        },
      });
    }
  };

  const onVerifyEmail = (data: { email: string; opt: string }) => {
    verifyForgotEmailOtpMutation.mutate(
      { email: data.email, otp: data.otp },
      {
        onSuccess: () => {
          toast.success("Email verified successfully. Please enter your new password");
          navigate({ to: "/login" });
        },
        onError: (error: any) {
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
                  : "Reset password failed. Please try again.";
              toast.error(errorMsg);
            }
          } else {
            toast.error("Network error. Please check your connection.");
          }
        }
      },
    );
  };

  return (
    <div className="flex min-h-svh w-full items-center justify-center p-6 md:p-10">
      <div className="w-full max-w-sm">
        <ForgotPasswordForm />
      </div>
    </div>
  );
}
