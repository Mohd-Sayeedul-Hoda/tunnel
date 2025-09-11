import { ThemeProvider } from "./ThemeProvider";
import { Toaster } from "@/components/ui/sonner"

function Provider({ children }: { children: React.ReactNode }) {
  return (
    <>
      <ThemeProvider defaultTheme="light">
        {children}
      </ThemeProvider>
      <Toaster />
    </>
  )
}
export default Provider;
