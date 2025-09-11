import { Button } from '@/components/ui/button'
import { Link } from '@tanstack/react-router'


export function NotFound() {
  return (
    <div className="flex flex-col items-center justify-center py-20 text-center mt-5">
      <h1 className="text-4xl font-bold mb-4">404 - Page Not Found</h1>
      <p className="text-lg text-muted-foreground mb-8">
        Oops! The page you're looking for doesn't exist.
      </p>
      <Button asChild variant='secondary'>
        <Link to="/" >
          Go back to home
        </Link>
      </Button>
    </div>
  )
}
