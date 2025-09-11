import { createFileRoute, useNavigate, useSearch } from '@tanstack/react-router'
import { VerifyEmailForm } from '@/components/VerifyEmailForm'
import { useState } from 'react'

export const Route = createFileRoute('/verify-email')({
  component: VerifyEmailComponent,
  validateSearch: (search: Record<string, unknown>) => {
    return {
      email: (search.email as string) || '',
    }
  },
})

function VerifyEmailComponent() {
  const { email } = useSearch({ from: '/verify-email' })
  const navigate = useNavigate()
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const handleVerifyOtp = async (otp: string) => {
    setIsLoading(true)
    setError(null)
    
    try {
      // TODO: Replace with actual API call
      // const response = await apiClient.verifyEmail({ email, otp })
      
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 1000))
      
      // TODO: Handle successful verification
      console.log('Email verified successfully')
      
      // Redirect to login or dashboard
      navigate({ to: '/login' })
      
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Verification failed')
    } finally {
      setIsLoading(false)
    }
  }

  const handleResendOtp = async () => {
    setIsLoading(true)
    setError(null)
    
    try {
      // TODO: Replace with actual API call
      // await apiClient.resendVerificationEmail(email)
      
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 500))
      
      console.log('Verification email resent to:', email)
      
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to resend email')
    } finally {
      setIsLoading(false)
    }
  }

  if (!email) {
    return (
      <div className="flex min-h-svh w-full items-center justify-center p-6 md:p-10">
        <div className="w-full max-w-sm text-center">
          <p className="text-muted-foreground">No email provided for verification.</p>
        </div>
      </div>
    )
  }

  return (
    <div className="flex min-h-svh w-full items-center justify-center p-6 md:p-10">
      <div className="w-full max-w-sm">
        <VerifyEmailForm
          email={email}
          onVerify={handleVerifyOtp}
          onResend={handleResendOtp}
          isLoading={isLoading}
          error={error || undefined}
        />
      </div>
    </div>
  )
}