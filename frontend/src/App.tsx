import { type ReactNode } from "react";
import {
  BrowserRouter,
  NavLink,
  Navigate,
  Route,
  Routes,
} from "react-router-dom";
import { AuthProvider, useAuth } from "./auth/auth";

import { LoginPage } from "./pages/login";
import { MePage } from "./pages/me";
import { ResendVerifyPage } from "./pages/resend-verify";
import { SignupPage } from "./pages/signup";
import { VerifyEmailPage } from "./pages/verify-email";
import { AdminPage } from "./pages/admin";
import TopPage from "./pages/top";
import { ItemsPage } from "./pages/items";

type GuardProps = {
  children: ReactNode;
};

function GuardLoading() {
  return (
    <div className="mx-auto max-w-[1280px] px-4 py-10 md:px-8">
      <div className="rounded-[28px] border border-[#e6d9ce] bg-white px-6 py-10 text-center text-base font-bold text-[#6f6259] shadow-[0_8px_24px_rgba(110,78,56,0.06)]">
        読み込み中です...
      </div>
    </div>
  );
}

function RequireAuth({ children }: GuardProps) {
  const { user, loading } = useAuth();

  if (loading) {
    return <GuardLoading />;
  }

  if (!user) {
    return <Navigate to="/login" replace />;
  }

  return <>{children}</>;
}

function RequireAdmin({ children }: GuardProps) {
  const { user, loading } = useAuth();

  if (loading) {
    return <GuardLoading />;
  }

  if (!user) {
    return <Navigate to="/login" replace />;
  }

  if (user.role !== "admin") {
    return <Navigate to="/me" replace />;
  }

  return <>{children}</>;
}

function HeaderNavItem({ to, label }: { to: string; label: string }) {
  return (
    <NavLink
      to={to}
      className={({ isActive }) =>
        [
          "inline-flex items-center rounded-full border px-4 py-2 text-sm font-bold transition",
          isActive
            ? "border-[#7b523a] bg-[#7b523a] text-white shadow-[0_8px_20px_rgba(123,82,58,0.22)]"
            : "border-[#d9c6b8] bg-white text-[#7b523a] hover:-translate-y-0.5 hover:bg-[#f8efe7]",
        ].join(" ")
      }
    >
      {label}
    </NavLink>
  );
}

function HeaderAuthActions() {
  const { user, logout } = useAuth();

  async function onLogout() {
    await logout();
  }

  if (!user) {
    return (
      <div className="flex flex-wrap items-center justify-end gap-2">
        <NavLink
          to="/login"
          className={({ isActive }) =>
            [
              "inline-flex items-center rounded-full border px-4 py-2 text-sm font-bold transition",
              isActive
                ? "border-[#7b523a] bg-[#7b523a] text-white"
                : "border-[#d9c6b8] bg-white text-[#7b523a] hover:bg-[#f8efe7]",
            ].join(" ")
          }
        >
          ログイン
        </NavLink>

        <NavLink
          to="/signup"
          className="inline-flex items-center rounded-full bg-[#4e342e] px-4 py-2 text-sm font-bold text-white transition hover:opacity-90"
        >
          新規登録
        </NavLink>

        <NavLink
          to="/resend-verify"
          className={({ isActive }) =>
            [
              "inline-flex items-center rounded-full border px-4 py-2 text-sm font-bold transition",
              isActive
                ? "border-[#c8b2a1] bg-[#efe4db] text-[#7b523a]"
                : "border-[#e1d2c7] bg-[#fbf7f2] text-[#8b6a58] hover:bg-[#f5ece4]",
            ].join(" ")
          }
        >
          確認メール再送
        </NavLink>
      </div>
    );
  }

  return (
    <div className="flex flex-wrap items-center justify-end gap-3">
      <div className="flex items-center gap-2 rounded-full border border-[#e2d3c7] bg-[#fbf7f2] px-4 py-2">
        <span className="rounded-full bg-[#f1e3d6] px-2.5 py-1 text-[11px] font-black tracking-[0.2em] text-[#7b523a] uppercase">
          {user.role}
        </span>

        <span className="max-w-[180px] truncate text-sm font-bold text-[#5f4a40]">
          {user.email}
        </span>
      </div>

      <button
        type="button"
        onClick={() => void onLogout()}
        className="inline-flex items-center rounded-full bg-[#4e342e] px-4 py-2 text-sm font-bold text-white transition hover:opacity-90"
      >
        ログアウト
      </button>
    </div>
  );
}

function AppHeader() {
  const { user } = useAuth();

  return (
    <header className="sticky top-0 z-50 border-b border-[#e6d9ce]/80 bg-[#fffaf5]/90 backdrop-blur">
      <div className="h-1 w-full bg-gradient-to-r from-[#6f4e37] via-[#b08968] to-[#e7d6c7]" />

      <div className="mx-auto max-w-[1440px] px-4 py-4 md:px-8">
        <div className="rounded-[28px] border border-[#e6d9ce] bg-white/95 px-4 py-4 shadow-[0_10px_28px_rgba(110,78,56,0.08)] md:px-6">
          <div className="flex flex-col gap-4 xl:flex-row xl:items-center xl:justify-between">
            <div className="flex items-center gap-4">
              <NavLink
                to="/"
                className="flex min-w-0 items-center gap-3 rounded-[22px] border border-[#eadfd5] bg-[#fcf7f2] px-3 py-3 transition hover:bg-[#f7efe8]"
              >
                <div className="flex h-12 w-12 items-center justify-center rounded-2xl bg-[#6f4e37] text-lg font-black tracking-[0.12em] text-white">
                  CS
                </div>

                <div className="min-w-0">
                  <h1 className="truncate text-lg font-black text-[#4e342e] md:text-xl">
                    coffee spa
                  </h1>
                  <p className="hidden text-sm font-semibold text-[#8a7b71] md:block">
                    豆・レシピ・セール・店舗情報を一つに。
                  </p>
                </div>
              </NavLink>
            </div>

            <div className="flex flex-col gap-3 xl:items-end">
              <nav className="flex flex-wrap items-center gap-2">
                <HeaderNavItem to="/" label="トップ" />
                {user ? <HeaderNavItem to="/me" label="マイページ" /> : null}
                {user?.role === "admin" ? (
                  <HeaderNavItem to="/admin" label="管理" />
                ) : null}
              </nav>

              <HeaderAuthActions />
            </div>
          </div>
        </div>
      </div>
    </header>
  );
}

function NotFoundPage() {
  return (
    <div className="mx-auto max-w-[1280px] px-4 py-10 md:px-8">
      <div className="rounded-[32px] border border-[#eadfd4] bg-white px-8 py-16 text-center shadow-[0_8px_24px_rgba(110,78,56,0.06)]">
        <p className="mb-3 text-sm font-black tracking-[0.28em] text-[#a1775b] uppercase">
          404
        </p>
        <h2 className="mb-4 text-3xl font-black text-[#4e342e]">
          page not found
        </h2>
        <p className="text-base font-semibold text-[#7a6f68]">
          指定されたページは見つかりませんでした。
        </p>
      </div>
    </div>
  );
}

function AppRoutes() {
  return (
    <div className="min-h-screen bg-[#f6f1eb]">
      <AppHeader />

      <main className="pb-10">
        <Routes>
          <Route path="/" element={<TopPage />} />
          <Route path="/items" element={<ItemsPage />} />
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
