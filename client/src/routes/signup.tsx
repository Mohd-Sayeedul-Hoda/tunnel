import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { SignupForm } from '@/components/SignupForm'
import { VerifyEmailForm } from '@/components/VerifyEmailForm'
import { useSignup } from '@/contexts/SignupContext'
import { useState } from 'react'

export const Route = createFileRoute('/signup')({
  component: SignupComponent,
})

function SignupComponent() {
  const { signupData, setSignupData, currentStep } = useSignup()
  const navigate = useNavigate()
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const handleSignupNext = async (data: { email: string; password: string; name: string }) => {
    setIsLoading(true)
    setError(null)
    
    try {
      // TODO: Replace with actual API call
      // const response = await apiClient.signup(data)
      
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 1000))
      
      // Store signup data and move to verification step
      setSignupData(data)
      
      // TODO: Send verification email
      console.log('Sending verification email to:', data.email)
      
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Signup failed')
    } finally {
      setIsLoading(false)
    }
  }

  const handleVerifyOtp = async (otp: string) => {
    if (!signupData) return
    
    setIsLoading(true)
    setError(null)
    
    try {
      // TODO: Replace with actual API call
      // const response = await apiClient.verifyEmail({ email: signupData.email, otp })
      
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
    if (!signupData) return
    
    setIsLoading(true)
    setError(null)
    
    try {
      // TODO: Replace with actual API call
      // await apiClient.resendVerificationEmail(signupData.email)
      
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 500))
      
      console.log('Verification email resent to:', signupData.email)
      
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to resend email')
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="flex min-h-svh w-full items-center justify-center p-6 md:p-10">
      <div className="w-full max-w-sm">
        {currentStep === 'signup' ? (
          <SignupForm onNext={handleSignupNext} />
        ) : signupData ? (
          <VerifyEmailForm
            email={signupData.email}
            onVerify={handleVerifyOtp}
            onResend={handleResendOtp}
            isLoading={isLoading}
            error={error || undefined}
          />
        ) : (
          <SignupForm onNext={handleSignupNext} />
        )}
      </div>
    </div>
  )
}
