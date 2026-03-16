import "./App.css";
import { BrowserRouter, Link, Navigate, Route, Routes } from "react-router-dom";
import { AuthProvider, useAuth } from "./auth/auth";

import { LoginPage } from "./pages/login";
import { MePage } from "./pages/me";
import { ResendVerifyPage } from "./pages/resend-verify";
import { SignupPage } from "./pages/signup";
import { TopPage } from "./pages/top";
import { VerifyEmailPage } from "./pages/verify-email";
import { AdminPage } from "./pages/admin";

type GuardProps = {
  children: React.ReactNode;
};

function RequireAuth({ children }: GuardProps) {
  const { user, loading } = useAuth();

  if (loading) {
    return <div>loading...</div>;
  }

  if (!user) {
    return <Navigate to="/login" replace />;
  }

  return <>{children}</>;
}

function RequireAdmin({ children }: GuardProps) {
  const { user, loading } = useAuth();

  if (loading) {
    return <div>loading...</div>;
  }

  if (!user) {
    return <Navigate to="/login" replace />;
  }

  if (user.role !== "admin") {
    return <Navigate to="/me" replace />;
  }

  return <>{children}</>;
}

function NotFoundPage() {
  return <div>page not found</div>;
}

function AppRoutes() {
  const { user, logout } = useAuth();

  async function onLogout() {
    await logout();
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

          {user?.role === "admin" ? <Link to="/admin">admin</Link> : null}

          {user ? (
            <button type="button" onClick={() => void onLogout()}>
              logout
            </button>
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

          <Route
            path="/admin"
            element={
              <RequireAdmin>
                <AdminPage />
              </RequireAdmin>
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
