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

//認証状態を配るProvider
export function AuthProvider({ children }: { children: ReactNode }) {
  //現在ログイン中のユーザー
  const [user, setUser] = useState<User | null>(null);

  //初期認証確認中かどうか
  const [loading, setLoading] = useState<boolean>(true);

  //開発中の二重実行や多重初期化を防ぐため
  const initStartedRef = useRef<boolean>(false);

  //サインアップ処理
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

  //メール確認処理
  async function verifyEmail(token: string): Promise<void> {
    await api<void>("/auth/verify-email", {
      method: "POST",
      body: { token },
    });
  }

  //確認メール再送処理
  async function resendVerify(email: string): Promise<string> {
    await api<void>("/auth/resend-verify", {
      method: "POST",
      body: { email },
    });

    return `${email} 宛に確認メールを再送しました。`;
  }

  //access tokenを使って現在ユーザー情報を取得
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

  //ログイン処理
  async function login(email: string, password: string): Promise<void> {
    const res = await api<LoginResponse>("/auth/login", {
      method: "POST",
      body: { email, password },
    });
    if (!res) {
      throw new Error("login response is empty");
    }

    //access tokenを保存
    setToken(res.access_token);

    //画面上の認証状態も即更新
    setUser(res.user);
  }

  //refresh cookie + csrf cookieを使ってaccess tokenを取り直す
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

      //新しいaccess tokenを保存
      setToken(res.access_token);

      //ユーザー情報も更新
      setUser(res.user);

      return true;
    } catch (err: unknown) {
      //refresh に失敗したらフロント側の認証状態はクリア
      clearToken();
      if (err instanceof ApiError) {
        if (err.status === 401 || err.status === 403) {
          return false;
        }
      }
      return false;
    }
  }

  //ログアウト処理
  async function logout(): Promise<void> {
    try {
      await api<void>("/auth/logout", {
        method: "POST",
        auth: true,
        csrf: true,
      });
    } finally {
      //ログアウト状態に戻す
      clearToken();
      setUser(null);
    }
  }

  useEffect(() => {
    //StrictModeで二重にinitが走らないように
    if (initStartedRef.current) {
      return;
    }

    initStartedRef.current = true;

    async function init() {
      try {
        //localStorageにaccess tokenがあるか確認
        const token = getToken();

        //accesstokenがあるなら/me確認
        if (token) {
          try {
            await loadMe();
            return;
          } catch {
            //壊れたtokenや期限切れなら一旦破棄
            clearToken();
            setUser(null);
          }
        }

        //未ログイン状態で無駄にrefreshを打たない
        //csrf_token cookie がある時だけrefreshを試す
        const csrfToken = getCookie("csrf_token");

        if (!csrfToken) {
          setUser(null);
          return;
        }

        await refresh();
      } finally {
        //認証初期化が終わったらloadingを落とす
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
