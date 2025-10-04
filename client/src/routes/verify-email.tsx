import { createFileRoute } from "@tanstack/react-router";
import { VerifyEmailForm } from "@/components/VerifyEmailForm";

export const Route = createFileRoute("/verify-email")({
  component: VerifyEmailComponent,
});

function VerifyEmailComponent() {
  return (
    <div className="flex min-h-svh w-full items-center justify-center p-6 md:p-10">
      <div className="w-full max-w-sm">
        <VerifyEmailForm email="sayeedulhoda@gmail.com" />
      </div>
    </div>
  );
}
