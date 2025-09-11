import Container from '@/components/global/Container';
import { createRootRoute, Outlet, useMatchRoute } from '@tanstack/react-router'
import { TanStackRouterDevtools } from '@tanstack/react-router-devtools'
import Provider from '@/components/provider/Provider'
import NavBar from '@/components/navbar/NavBar'
import Footer from '@/components/footer/Footer';


export const Route = createRootRoute({

  component: RootLayout,
})

function RootLayout() {
  const matchRoute = useMatchRoute();
  const hideNavRoutes = ["/login", "/forgot-password", "/signup", "/verify-email"];
  const matchedNoNavRoutes = hideNavRoutes.some((route) => matchRoute({ to: route }));

  return (
    <Provider>
      <div className="min-h-screen w-full relative">
        <div
          className="absolute inset-0 z-0"
          style={{
            background: "radial-gradient(125% 125% at 50% 90%, var(--background) 40%, var(--primary) 100%)",
          }}
        />
        <div className="relative z-10 flex flex-col min-h-screen">
          {!matchedNoNavRoutes ? (
            <NavBar />
          ) : undefined}
          <main className="flex-1">
            <Container className={matchedNoNavRoutes ? '' : 'py-20'}>
              <Outlet />
            </Container>
          </main>
          {!matchedNoNavRoutes ? (
            <Footer />
          ) : undefined}
          <TanStackRouterDevtools />
        </div>
      </div>
    </Provider >
  )
}
