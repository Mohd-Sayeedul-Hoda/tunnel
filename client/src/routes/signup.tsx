import { SignupForm } from "@/components/SignupForm";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/signup")({
  component: SignupComponent,
});

function SignupComponent() {
  return (
    <div className="flex min-h-svh w-full items-center justify-center p-6 md:p-10">
      <div className="w-full max-w-sm">
        <SignupForm onNext={() => console.log("hello word")} />
      </div>
    </div>
  );
}
