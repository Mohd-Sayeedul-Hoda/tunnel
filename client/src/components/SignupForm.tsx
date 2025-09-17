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
import { ArrowLeftIcon, EyeIcon, EyeOffIcon } from "lucide-react"
import { useState } from "react"
import { signupFormSchema, type SignupFormData } from "@/lib/validations"
import { ZodError } from "zod"

interface SignupFormProps extends React.ComponentProps<"div"> {
  onNext: (data: { email: string; password: string; name: string }) => void
}

export function SignupForm({
  className,
  onNext,
  ...props
}: SignupFormProps) {
  const [formData, setFormData] = useState<SignupFormData>({
    name: "",
    email: "",
    password: "",
    confirmPassword: ""
  })
  const [showPassword, setShowPassword] = useState(false)
  const [showConfirmPassword, setShowConfirmPassword] = useState(false)
  const [errors, setErrors] = useState<Record<string, string>>({})
  const [touchedFields, setTouchedFields] = useState<Record<string, boolean>>({})
  const [isSubmitting, setIsSubmitting] = useState(false)

  const validateForm = (data: SignupFormData) => {
    try {
      signupFormSchema.parse(data)
      setErrors({})
      return true
    } catch (error) {
      if (error instanceof ZodError) {
        const newErrors: Record<string, string> = {}
        error.issues.forEach((err) => {
          if (err.path[0]) {
            newErrors[err.path[0] as string] = err.message
          }
        })
        setErrors(newErrors)
      }
      return false
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsSubmitting(true)

    // Mark all fields as touched
    const allFields: (keyof SignupFormData)[] = ['name', 'email', 'password', 'confirmPassword']
    const newTouchedFields = allFields.reduce((acc, field) => {
      acc[field] = true
      return acc
    }, {} as Record<keyof SignupFormData, boolean>)

    setTouchedFields(prev => ({ ...prev, ...newTouchedFields }))

    if (validateForm(formData)) {
      try {
        await onNext({
          email: formData.email,
          password: formData.password,
          name: formData.name
        })
      } catch (error) {
        console.error('Signup error:', error)
      }
    }

    setIsSubmitting(false)
  }

  const handleInputChange = (field: keyof SignupFormData, value: string) => {
    setFormData(prev => ({ ...prev, [field]: value }))

    if (!touchedFields[field]) {
      setTouchedFields(prev => ({ ...prev, [field]: true }))
    }

    if (errors[field]) {
      setErrors(prev => ({ ...prev, [field]: "" }))
    }

    if (touchedFields[field] || isSubmitting) {
      const updatedData = { ...formData, [field]: value }
      if (field === 'password' || field === 'confirmPassword') {
        validateForm(updatedData)
      } else {
        try {
          signupFormSchema.pick({ [field]: true }).parse({ [field]: value })
          if (errors[field]) {
            setErrors(prev => ({ ...prev, [field]: "" }))
          }
        } catch (error) {
          if (error instanceof ZodError) {
            const fieldError = error.issues.find(err => err.path[0] === field)
            if (fieldError) {
              setErrors(prev => ({ ...prev, [field]: fieldError.message }))
            }
          }
        }
      }
    }
  }

  return (
    <div className={cn("flex flex-col gap-6", className)} {...props}>
      <Card>
        <CardHeader>
          <div>
            <CardTitle>Create your account</CardTitle>
            <CardDescription>
              Enter your details to get started
            </CardDescription>
          </div>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit}>
            <div className="flex flex-col gap-6">
              <div className="grid gap-3">
                <Label htmlFor="name">Full Name</Label>
                <Input
                  id="name"
                  type="text"
                  placeholder="John Doe"
                  value={formData.name}
                  onChange={(e) => handleInputChange("name", e.target.value)}
                  onBlur={() => setTouchedFields(prev => ({ ...prev, name: true }))}
                  className={cn(
                    "transition-colors",
                    errors.name && (touchedFields.name || isSubmitting) && "border-red-500 focus-visible:ring-red-500"
                  )}
                  required
                />
                {errors.name && (touchedFields.name || isSubmitting) && (
                  <p className="text-sm text-red-500 flex items-center gap-1">
                    <span className="text-red-500">•</span>
                    {errors.name}
                  </p>
                )}
              </div>

              <div className="grid gap-3">
                <Label htmlFor="email">Email</Label>
                <Input
                  id="email"
                  type="email"
                  placeholder="john@example.com"
                  value={formData.email}
                  onChange={(e) => handleInputChange("email", e.target.value)}
                  onBlur={() => setTouchedFields(prev => ({ ...prev, email: true }))}
                  className={cn(
                    "transition-colors",
                    errors.email && (touchedFields.email || isSubmitting) && "border-red-500 focus-visible:ring-red-500"
                  )}
                  required
                />
                {errors.email && (touchedFields.email || isSubmitting) && (
                  <p className="text-sm text-red-500 flex items-center gap-1">
                    <span className="text-red-500">•</span>
                    {errors.email}
                  </p>
                )}
              </div>

              <div className="grid gap-3">
                <Label htmlFor="password">Password</Label>
                <div className="relative">
                  <Input
                    id="password"
                    type={showPassword ? "text" : "password"}
                    placeholder="Enter your password"
                    value={formData.password}
                    onChange={(e) => handleInputChange("password", e.target.value)}
                    onBlur={() => setTouchedFields(prev => ({ ...prev, password: true }))}
                    className={cn(
                      "pr-10 transition-colors",
                      errors.password && (touchedFields.password || isSubmitting) && "border-red-500 focus-visible:ring-red-500"
                    )}
                    required
                  />
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    className="absolute right-0 top-0 h-full px-3 py-2 hover:bg-transparent"
                    onClick={() => setShowPassword(!showPassword)}
                  >
                    {showPassword ? (
                      <EyeOffIcon className="h-4 w-4" />
                    ) : (
                      <EyeIcon className="h-4 w-4" />
                    )}
                  </Button>
                </div>
                {errors.password && (touchedFields.password || isSubmitting) && (
                  <p className="text-sm text-red-500 flex items-center gap-1">
                    <span className="text-red-500">•</span>
                    {errors.password}
                  </p>
                )}
                {!errors.password && formData.password && (
                  <div className="text-xs text-muted-foreground">
                    Password must contain: uppercase, lowercase, number, and special character
                  </div>
                )}
              </div>

              <div className="grid gap-3">
                <Label htmlFor="confirmPassword">Confirm Password</Label>
                <div className="relative">
                  <Input
                    id="confirmPassword"
                    type={showConfirmPassword ? "text" : "password"}
                    placeholder="Confirm your password"
                    value={formData.confirmPassword}
                    onChange={(e) => handleInputChange("confirmPassword", e.target.value)}
                    onBlur={() => setTouchedFields(prev => ({ ...prev, confirmPassword: true }))}
                    className={cn(
                      "pr-10 transition-colors",
                      errors.confirmPassword && (touchedFields.confirmPassword || isSubmitting) && "border-red-500 focus-visible:ring-red-500"
                    )}
                    required
                  />
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    className="absolute right-0 top-0 h-full px-3 py-2 hover:bg-transparent"
                    onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                  >
                    {showConfirmPassword ? (
                      <EyeOffIcon className="h-4 w-4" />
                    ) : (
                      <EyeIcon className="h-4 w-4" />
                    )}
                  </Button>
                </div>
                {errors.confirmPassword && (touchedFields.confirmPassword || isSubmitting) && (
                  <p className="text-sm text-red-500 flex items-center gap-1">
                    <span className="text-red-500">•</span>
                    {errors.confirmPassword}
                  </p>
                )}
              </div>

              <div className="flex flex-col gap-3">
                <Button
                  type="submit"
                  className="w-full"
                  disabled={isSubmitting}
                >
                  {isSubmitting ? "Creating Account..." : "Create Account"}
                </Button>
              </div>
            </div>
            <div className="mt-4 text-center text-sm">
              Already have an account?{" "}
              <Link to="/login" className="underline underline-offset-4">
                Sign in
              </Link>
            </div>
          </form>
        </CardContent>
      </Card>

      <div className="text-center">
        <Link
          to="/"
          className="inline-flex items-center gap-2 text-sm text-muted-foreground hover:text-foreground transition-colors"
        >
          <ArrowLeftIcon className="h-4 w-4" />
          Back to home
        </Link>
      </div>
    </div>
  )
}
