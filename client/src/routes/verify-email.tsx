import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { VerifyEmailForm } from "@/components/VerifyEmailForm";
import { useSendVerficationEmail, useVerfiyEmailOtp } from "@/hooks/use-auth";
import { useEffect, useState } from "react";
import { toast } from "sonner";

export const Route = createFileRoute("/verify-email")({
  component: VerifyEmailComponent,
  validateSearch: (search: Record<string, unknown>) => {
    return {
      email: search.email as string,
    };
  },
});

function VerifyEmailComponent() {
  const navigate = useNavigate();

  const [resendCooldown, setResendCooldown] = useState(0);

  const { email } = Route.useSearch();
  const sendEmailMutation = useSendVerficationEmail();
  const verifyEmailMutation = useVerfiyEmailOtp();

  useEffect(() => {
    if (resendCooldown > 0) {
      const timer = setTimeout(() => {
        setResendCooldown(resendCooldown - 1);
      }, 1000);
      return () => clearTimeout(timer);
    }
  }, [resendCooldown]);

  const handleSendEmail = () => {
    if (email) {
      sendEmailMutation.mutate(
        { email: email },
        {
          onSuccess: () => {
            setResendCooldown(60);
            toast.success("Verification email sent successfully!");
          },
          onError: (error: any) => {
            if (error.response) {
              const { status, data } = error.response;
              const errorMsg =
                typeof data.error === "string"
                  ? data.error
                  : "Failed to send verification email. Please try again";
              toast.error(errorMsg);
            } else {
              toast.error("Network error. Please check your connection.");
            }
          },
        },
      );
    }
  };

  const handleVerifyEmail = (otp: string) => {
    if (email) {
      verifyEmailMutation.mutate(
        { email: email, otp: otp },
        {
          onSuccess: () => {
            toast.success("Email verfication succesfull");
            navigate({ to: "/dashboard" });
          },
          onError: (error: any) => {
            if (error.response) {
              const { status, data } = error.response;

              if (status === 422) {
                const errorMessages = Object.values(data.error) as string[];
                errorMessages.forEach((errorMsg) => {
                  toast.error(errorMsg);
                });
              } else {
                const errorMsg =
                  typeof data.error === "string"
                    ? data.error
                    : "Email verfication failed. Please try again.";
                toast.error(errorMsg);
              }
            } else {
              toast.error("Network error. Please check your connection.");
            }
          },
        },
      );
    }
  };

  const handleResend = () => {
    if (resendCooldown === 0) {
      handleSendEmail();
    }
  };

  useEffect(() => {
    if (email && !sendEmailMutation.isPending && !sendEmailMutation.isSuccess) {
      handleSendEmail();
    }
  }, [email]);

  return (
    <div className="flex min-h-svh w-full items-center justify-center p-6 md:p-10">
      <div className="w-full max-w-sm">
        <VerifyEmailForm
          email={email}
          onResend={handleResend}
          onVerify={handleVerifyEmail}
          resendCooldown={resendCooldown}
          isSendingEmail={sendEmailMutation.isPending}
          isVerfyingEmail={verifyEmailMutation.isPending}
        />
      </div>
    </div>
  );
}
