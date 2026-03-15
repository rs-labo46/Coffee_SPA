import { useState } from "react";
import { Link } from "react-router-dom";
import { useAuth } from "../auth/auth";
import { ApiError } from "../lib/api";

export function ResendVerifyPage() {
  const { resendVerify } = useAuth();

  const [email, setEmail] = useState<string>("");
  const [msg, setMsg] = useState<string>("");
  const [loading, setLoading] = useState<boolean>(false);

  async function onSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();

    setMsg("");
    setLoading(true);

    try {
      const message = await resendVerify(email);
      setMsg(
        `${message} backendログにverifyリンクが出るので、そのリンクを開いてください。`,
      );
    } catch (err: unknown) {
      if (err instanceof ApiError) {
        setMsg(err.message);
      } else {
        setMsg("確認メールの再送に失敗しました");
      }
    } finally {
      setLoading(false);
    }
  }

  return (
    <div>
      <h1>確認メール再送</h1>

      <form onSubmit={onSubmit}>
        <div>
          <input
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="email"
            type="email"
          />
        </div>

        <button type="submit" disabled={loading}>
          {loading ? "loading..." : "resend"}
        </button>
      </form>

      {msg ? <p>{msg}</p> : null}

      <p>
        <Link to="/login">ログインへ</Link>
      </p>
    </div>
  );
}
