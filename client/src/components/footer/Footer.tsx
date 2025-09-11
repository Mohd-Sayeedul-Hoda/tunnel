import { Link } from '@tanstack/react-router'
import Container from '@/components/global/Container'
import { GiArrowDunk } from "react-icons/gi"

export default function Footer() {
  return (
    <footer className="bg-background border-t">
      <Container>
        <div className="flex flex-col space-y-6 py-8">
          <div className="flex flex-col md:flex-row md:items-center md:justify-between space-y-4 md:space-y-0">
            <div className="flex items-center space-x-2">
              <Link to="/" className="flex items-center space-x-2">
                <GiArrowDunk className="w-6 h-6" />
                <span className="text-xl font-bold">Tunnel</span>
              </Link>
            </div>

            <div className="flex flex-wrap gap-6">
              <Link to="/docs" className="text-sm text-muted-foreground hover:text-foreground transition-colors">
                Documentation
              </Link>
              <Link to="/about" className="text-sm text-muted-foreground hover:text-foreground transition-colors">
                About
              </Link>
              <Link to="/contact" className="text-sm text-muted-foreground hover:text-foreground transition-colors">
                Contact
              </Link>
              <Link to="/privacy" className="text-sm text-muted-foreground hover:text-foreground transition-colors">
                Privacy Policy
              </Link>
              <Link to="/terms" className="text-sm text-muted-foreground hover:text-foreground transition-colors">
                Terms of Service
              </Link>
            </div>
          </div>

          <div className="flex flex-col md:flex-row md:items-center md:justify-between space-y-4 md:space-y-0 pt-6 ">
            <div className="text-sm text-muted-foreground">
              © 2025 Tunnel.com. All rights reserved.
            </div>

            <div className="flex items-center space-x-4">
              <div className="flex items-center space-x-4">
                <a href="https://www.linkedin.com/in/mohd-sayeedul-hoda/" className="text-sm text-muted-foreground hover:text-foreground transition-colors">
                  Linkedin
                </a>
                <a href="https://x.com/HodaSayeed" className="text-sm text-muted-foreground hover:text-foreground transition-colors">
                  Twitter
                </a>
                <a href="https://github.com/Mohd-Sayeedul-Hoda" className="text-sm text-muted-foreground hover:text-foreground transition-colors">
                  GitHub
                </a>
              </div>
            </div>
          </div>
        </div>
      </Container>
    </footer>
  )
}
