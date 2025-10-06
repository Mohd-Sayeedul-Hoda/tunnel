import { createFileRoute } from "@tanstack/react-router";
import { VerifyEmailForm } from "@/components/VerifyEmailForm";

export const Route = createFileRoute("/verify-email")({
  component: VerifyEmailComponent,
  validateSearch: (search: Record<string, unknown>) => {
    return {
      email: search.email as string,
    };
  },
});

function VerifyEmailComponent() {
  const { email } = Route.useSearch();
  return (
    <div className="flex min-h-svh w-full items-center justify-center p-6 md:p-10">
      <div className="w-full max-w-sm">
        <VerifyEmailForm email={email} />
      </div>
    </div>
  );
}
