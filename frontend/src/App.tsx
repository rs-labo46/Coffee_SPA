import "./App.css";
import { BrowserRouter, Link, Navigate, Route, Routes } from "react-router-dom";
import { AuthProvider, useAuth } from "./auth/auth";
import { LoginPage } from "./pages/login";
import { MePage } from "./pages/me";
import { ResendVerifyPage } from "./pages/resend-verify";
import { SignupPage } from "./pages/signup";
import { TopPage } from "./pages/top";
import { VerifyEmailPage } from "./pages/verify-email";

type RequireAuthProps = {
  children: React.ReactNode;
};

function RequireAuth({ children }: RequireAuthProps) {
  const { user, loading } = useAuth();

  if (loading) {
    return <div>loading...</div>;
  }

  if (!user) {
    return <Navigate to="/login" replace />;
  }

  return <>{children}</>;
}

function AppRoutes() {
  const { user, logout } = useAuth();

  async function onLogout() {
    await logout();
  }
  function NotFoundPage() {
    return <div>page not found</div>;
  }

  return (
    <div className="app">
      <header className="head">
        <nav className="nav">
          <Link to="/">top</Link>
          <Link to="/me">me</Link>
          {!user ? <Link to="/login">login</Link> : null}
          {!user ? <Link to="/signup">signup</Link> : null}
          {!user ? <Link to="/resend-verify">resend verify</Link> : null}
          {user ? (
            <button onClick={() => void onLogout()}>logout</button>
          ) : null}
        </nav>
      </header>

      <main>
        <Routes>
          <Route path="/" element={<TopPage />} />
          <Route path="/login" element={<LoginPage />} />
          <Route path="/signup" element={<SignupPage />} />
          <Route path="/verify-email" element={<VerifyEmailPage />} />
          <Route path="/resend-verify" element={<ResendVerifyPage />} />
          <Route
            path="/me"
            element={
              <RequireAuth>
                <MePage />
              </RequireAuth>
            }
          />
          <Route path="*" element={<NotFoundPage />} />
        </Routes>
      </main>
    </div>
  );
}

export default function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <AppRoutes />
      </BrowserRouter>
    </AuthProvider>
  );
}
