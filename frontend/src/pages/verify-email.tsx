import { useEffect, useState } from "react";
import { Link, useSearchParams } from "react-router-dom";
import { useAuth } from "../auth/auth";
import { ApiError } from "../lib/api";

export function VerifyEmailPage() {
  const { verifyEmail } = useAuth();
  const [params] = useSearchParams();

  const [msg, setMsg] = useState<string>("確認中...");
  const [done, setDone] = useState<boolean>(false);

  const token = params.get("token") || "";
  useEffect(() => {
    async function run() {
      if (!token) {
        setMsg("tokenがありません");
        return;
      }

      try {
        await verifyEmail(token);
        setMsg("メール確認が完了しました。ログインしてください。");
        setDone(true);
      } catch (err: unknown) {
        if (err instanceof ApiError) {
          setMsg(err.message);
        } else {
          setMsg("メール確認に失敗しました");
        }
      }
    }

    void run();
  }, [token]);

  return (
    <div>
      <h1>メール確認</h1>
      <p>{msg}</p>

      {done ? (
        <p>
          <Link to="/login">ログインへ</Link>
        </p>
      ) : (
        <p>
          <Link to="/resend-verify">確認メールを再送する</Link>
        </p>
      )}
    </div>
  );
}
