import {
  createContext,
  useContext,
  useEffect,
  useMemo,
  useRef,
  useState,
  type ReactNode,
} from "react";
import {
  ApiError,
  api,
  clearToken,
  getCookie,
  getToken,
  setToken,
} from "../lib/api";

type Role = "user" | "admin";

type User = {
  id: number;
  email: string;
  role: Role;
  token_ver: number;
  email_verified: boolean;
};

type SignupResponse = {
  user: User;
};

type LoginResponse = {
  access_token: string;
  user: User;
};

type RefreshResponse = {
  access_token: string;
  user: User;
};

type MeResponse = {
  user: User;
};

type AuthCtx = {
  user: User | null;
  loading: boolean;
  signup: (email: string, password: string) => Promise<string>;
  verifyEmail: (token: string) => Promise<void>;
  resendVerify: (email: string) => Promise<string>;
  login: (email: string, password: string) => Promise<void>;
  refresh: () => Promise<boolean>;
  logout: () => Promise<void>;
  loadMe: () => Promise<void>;
};

const AuthContext = createContext<AuthCtx | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const initStartedRefresh = useRef<boolean>(false);

  async function signup(email: string, password: string): Promise<string> {
    const res = await api<SignupResponse>("/auth/signup", {
      method: "POST",
      body: { email, password },
    });

    if (!res) {
      return "登録しました。メール確認をしてください。";
    }

    return `${res.user.email} を登録しました。メール確認をしてください。`;
  }

  async function verifyEmail(token: string): Promise<void> {
    await api<void>("/auth/verify-email", {
      method: "POST",
      body: { token },
    });
  }

  async function resendVerify(email: string): Promise<string> {
    await api<void>("/auth/resend-verify", {
      method: "POST",
      body: { email },
    });

    return `${email} 宛に確認メールを再送しました。`;
  }

  async function loadMe(): Promise<void> {
    const res = await api<MeResponse>("/me", {
      method: "GET",
      auth: true,
    });

    if (!res) {
      throw new Error("me response is empty");
    }

    setUser(res.user);
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
    setUser(res.user);
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
      setUser(res.user);
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
        csrf: true,
      });
    } finally {
      clearToken();
      setUser(null);
    }
  }

  useEffect(() => {
    if (initStartedRefresh.current) {
      return;
    }

    initStartedRefresh.current = true;

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

        // csrf_token cookie がある時だけ refresh を試す
        const csrf = getCookie("csrf_token");
        if (!csrf) {
          setUser(null);
          return;
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
      verifyEmail,
      resendVerify,
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
