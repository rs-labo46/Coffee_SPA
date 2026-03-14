import { useState } from "react";
import { Link } from "react-router-dom";
import { useAuth } from "../auth/auth";
import { ApiError } from "../lib/api";

export function SignupPage() {
  const { signup } = useAuth();

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [msg, setMsg] = useState("");
  const [loading, setLoading] = useState(false);

  async function onSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();

    setMsg("");
    setLoading(true);

    try {
      const message = await signup(email, password);
      setMsg(message);
    } catch (err: unknown) {
      if (err instanceof ApiError) {
        setMsg(err.message);
      } else {
        setMsg("登録に失敗しました");
      }
    } finally {
      setLoading(false);
    }
  }

  return (
    <div>
      <h1>サインアップ</h1>

      <form onSubmit={onSubmit}>
        <div>
          <input
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="email"
            type="email"
          />
        </div>

        <div>
          <input
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            placeholder="password"
            type="password"
          />
        </div>

        <button type="submit" disabled={loading}>
          {loading ? "loading..." : "signup"}
        </button>
      </form>

      {msg ? <p>{msg}</p> : null}

      <p>
        <Link to="/login">ログインへ</Link>
      </p>
    </div>
  );
}
