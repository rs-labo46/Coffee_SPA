import { useState } from "react";
import { Link, Navigate, useNavigate } from "react-router-dom";
import { useAuth } from "../auth/auth";
import { ApiError } from "../lib/api";

export function LoginPage() {
  const nav = useNavigate();
  const { login, user } = useAuth();

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [msg, setMsg] = useState("");
  const [loading, setLoading] = useState(false);

  if (user) {
    return <Navigate to="/me" replace />;
  }

  async function onSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();

    setMsg("");
    setLoading(true);

    try {
      await login(email, password);
      nav("/me");
    } catch (err: unknown) {
      if (err instanceof ApiError) {
        if (err.status === 401) {
          setMsg(
            "ログインに失敗しました。メール確認未完了、または認証情報が正しくありません。",
          );
        } else {
          setMsg(err.message);
        }
      } else {
        setMsg("ログインに失敗しました。");
      }
    } finally {
      setLoading(false);
    }
  }

  return (
    <div>
      <h1>login</h1>

      <form onSubmit={onSubmit}>
        <div>
          <input
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="email"
            type="email"
            autoComplete="email"
          />
        </div>

        <div>
          <input
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            placeholder="password"
            type="password"
            autoComplete="current-password"
          />
        </div>

        <button type="submit" disabled={loading}>
          {loading ? "loading..." : "login"}
        </button>
      </form>

      {msg ? <p>{msg}</p> : null}

      <p>
        <Link to="/signup">サインアップへ</Link>
      </p>

      <p>
        <Link to="/resend-verify">確認メールを再送する</Link>
      </p>
    </div>
  );
}
