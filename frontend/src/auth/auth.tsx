import {
  createContext,
  useContext,
  useEffect,
  useMemo,
  useState,
  type ReactNode,
} from "react";
import { ApiError, api, clearToken, getToken, setToken } from "../lib/api";

type Role = "user" | "admin";

type User = {
  id: number;
  email: string;
  role: Role;
  email_verified: boolean;
};
type SignupResponse = {
  message?: string;
};

type LoginResponse = {
  access_token: string;
};

type RefreshResponse = {
  access_token: string;
};

type AuthCtx = {
  user: User | null;
  loading: boolean;
  signup: (email: string, password: string) => Promise<string>;
  login: (email: string, password: string) => Promise<void>;
  refresh: () => Promise<boolean>;
  logout: () => Promise<void>;
  loadMe: () => Promise<void>;
};

const AuthContext = createContext<AuthCtx | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState<boolean>(true);

  async function signup(email: string, password: string): Promise<string> {
    const res = await api<SignupResponse>("/auth/signup", {
      method: "POST",
      body: { email, password },
    });

    if (!res) {
      return "登録しました";
    }

    return res.message || "登録しました。メール確認をしてください。";
  }

  async function loadMe(): Promise<void> {
    const res = await api<User>("/me", {
      method: "GET",
      auth: true,
    });

    if (!res) {
      throw new Error("me response is empty");
    }

    setUser(res);
  }

  async function login(email: string, password: string): Promise<void> {
    const res = await api<LoginResponse>("/auth/login", {
      method: "POST",
      body: { email, password },
    });

    if (!res) {
      throw new Error("login response is empty");
    }

    setToken(res.access_token);
    await loadMe();
  }

  async function refresh(): Promise<boolean> {
    try {
      const res = await api<RefreshResponse>("/auth/refresh", {
        method: "POST",
        csrf: true,
      });

      if (!res) {
        clearToken();
        setUser(null);
        return false;
      }

      setToken(res.access_token);
      await loadMe();
      return true;
    } catch (err: unknown) {
      clearToken();
      setUser(null);

      if (err instanceof ApiError) {
        if (err.status === 401 || err.status === 403) {
          return false;
        }
      }

      return false;
    }
  }

  async function logout(): Promise<void> {
    try {
      await api<void>("/auth/logout", {
        method: "POST",
        auth: true,
      });
    } finally {
      clearToken();
      setUser(null);
    }
  }

  useEffect(() => {
    async function init() {
      try {
        const token = getToken();

        if (token) {
          try {
            await loadMe();
            return;
          } catch {
            clearToken();
            setUser(null);
          }
        }

        await refresh();
      } finally {
        setLoading(false);
      }
    }

    void init();
  }, []);

  const value = useMemo<AuthCtx>(() => {
    return {
      user,
      loading,
      signup,
      login,
      refresh,
      logout,
      loadMe,
    };
  }, [user, loading]);

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth(): AuthCtx {
  const ctx = useContext(AuthContext);

  if (!ctx) {
    throw new Error("auth context not found");
  }

  return ctx;
}
