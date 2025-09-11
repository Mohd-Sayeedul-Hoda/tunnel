import { createFileRoute } from '@tanstack/react-router'
import { Link } from '@tanstack/react-router'
import { Button } from '@/components/ui/button'

export const Route = createFileRoute('/')({
  component: Index,
})

function Index() {
  return (
    <section className="flex flex-col items-center text-center py-20">
      <h1 className="text-4xl font-bold">
        Instant, Secure Tunnels to Your Localhost 🚀
      </h1>

      <p className="text-lg text-muted-foreground max-w-xl mt-2">
        Tunnel makes it effortless to share your local apps with the world.
        No complicated configs, just one command and you're online.
      </p>

      <div className="flex gap-2 mt-6">
        <Button>
          <Link to="/login">Get Started</Link>
        </Button>
        <Button variant="secondary">
          <Link to="/download-cli">Download CLI</Link>
        </Button>
      </div>

      <pre className="bg-muted p-4 rounded-xl text-sm mt-8 max-w-xl text-left shadow-sm">
        <code>
          curl -s https://get.tunnel.sh | bash{'\n'}
          tunnel auth &lt;API_KEY&gt;{'\n'}
          tunnel http 3000
        </code>
      </pre>

      <div className="flex gap-6 mt-6 text-sm text-muted-foreground">
        <span>✅ Free tier</span>
        <span>✅ Cross-platform</span>
        <span>✅ Secure NAT traversal</span>
      </div>
    </section>
  )
}
