import { Link } from '@tanstack/react-router'
import { Button } from "@/components/ui/button"
import Container from '@/components/global/Container'
import Logo from './Logo'

export default function NavBar() {

  return (
    <nav className="bg-transparent backdrop-blur ">
      <Container>
        <div className="flex h-16 items-center justify-between">
          <div className="flex items-center">
            <Logo />
          </div>
          <div className="flex items-center space-x-2">
            <Button variant="ghost" asChild>
              <Link to="/signup">Sign Up</Link>
            </Button>
            <Button asChild>
              <Link to="/login">Login</Link>
            </Button>
          </div>
        </div>
      </Container>
    </nav>
  )
}

