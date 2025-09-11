import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Link } from "@tanstack/react-router"
import { ArrowLeftIcon, MailIcon, RefreshCwIcon } from "lucide-react"
import { useState, useEffect } from "react"
import { otpSchema } from "@/lib/validations"
import { ZodError } from "zod"

interface VerifyEmailFormProps extends React.ComponentProps<"div"> {
  email: string
  onVerify: (otp: string) => void
  onResend: () => void
  isLoading?: boolean
  error?: string
}

export function VerifyEmailForm({
  className,
  email,
  onVerify,
  onResend,
  isLoading = false,
  error,
  ...props
}: VerifyEmailFormProps) {
  const [otp, setOtp] = useState("")
  const [resendCooldown, setResendCooldown] = useState(0)
  const [otpError, setOtpError] = useState("")

  // Resend cooldown timer
  useEffect(() => {
    if (resendCooldown > 0) {
      const timer = setTimeout(() => {
        setResendCooldown(resendCooldown - 1)
      }, 1000)
      return () => clearTimeout(timer)
    }
  }, [resendCooldown])

  const handleOtpChange = (value: string) => {
    // Only allow numbers and limit to 6 digits
    const numericValue = value.replace(/\D/g, "").slice(0, 6)
    setOtp(numericValue)
    setOtpError("")
  }

  const validateOtp = (otpValue: string) => {
    try {
      otpSchema.parse({ otp: otpValue })
      setOtpError("")
      return true
    } catch (error) {
      if (error instanceof ZodError) {
        const otpError = error.issues.find(err => err.path[0] === 'otp')
        if (otpError) {
          setOtpError(otpError.message)
        }
      }
      return false
    }
  }

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (validateOtp(otp)) {
      onVerify(otp)
    }
  }

  const handleResend = () => {
    if (resendCooldown === 0) {
      onResend()
      setResendCooldown(60) // 60 seconds cooldown
    }
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    // Auto-submit when 6 digits are entered
    if (otp.length === 6 && e.key === "Enter") {
      handleSubmit(e)
    }
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
                <Label htmlFor="otp">Verification Code</Label>
                <Input
                  id="otp"
                  type="text"
                  inputMode="numeric"
                  placeholder="000000"
                  value={otp}
                  onChange={(e) => handleOtpChange(e.target.value)}
                  onKeyDown={handleKeyDown}
                  className={cn(
                    "text-center text-2xl tracking-widest",
                    (otpError || error) ? "border-red-500" : ""
                  )}
                  maxLength={6}
                  required
                />
                {(otpError || error) && (
                  <p className="text-sm text-red-500 text-center">
                    {otpError || error}
                  </p>
                )}
              </div>

              <div className="flex flex-col gap-3">
                <Button
                  type="submit"
                  className="w-full"
                  disabled={otp.length !== 6 || isLoading}
                >
                  {isLoading ? (
                    <>
                      <RefreshCwIcon className="mr-2 h-4 w-4 animate-spin" />
                      Verifying...
                    </>
                  ) : (
                    "Verify Email"
                  )}
                </Button>
              </div>

              <div className="text-center">
                <p className="text-sm text-muted-foreground">
                  Didn't receive the code?{" "}
                  <Button
                    type="button"
                    variant="link"
                    className="p-0 h-auto text-sm"
                    onClick={handleResend}
                    disabled={resendCooldown > 0 || isLoading}
                  >
                    {resendCooldown > 0 ? (
                      `Resend in ${resendCooldown}s`
                    ) : (
                      "Resend code"
                    )}
                  </Button>
                </p>
              </div>
            </div>
          </form>
        </CardContent>
      </Card>

      <div className="text-center">
        <Link
          to="/signup"
          className="inline-flex items-center gap-2 text-sm text-muted-foreground hover:text-foreground transition-colors"
        >
          <ArrowLeftIcon className="h-4 w-4" />
          Back to signup
        </Link>
      </div>
    </div>
  )
}
