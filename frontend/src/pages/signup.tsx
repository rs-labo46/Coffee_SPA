import { useState } from "react";
import { Link } from "react-router-dom";
import { useAuth } from "../auth/auth";
import { ApiError } from "../lib/api";

export function SignupPage() {
  const { signup } = useAuth();
  const [ok, setOk] = useState(false);
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
      setMsg(
        `${message} backendログにverifyリンクが出るので、そのリンクを開いてください。`,
      );
    } catch (err: unknown) {
      if (err instanceof ApiError) {
        setMsg(err.message);
      } else {
        setMsg("登録に失敗しました。");
      }
    } finally {
      setLoading(false);
    }
  }

  return (
    <main className="min-h-[calc(100vh-120px)] bg-[#f6f1eb] px-4 py-8 md:px-8 md:py-10">
      <div className="mx-auto grid max-w-[1280px] gap-6 lg:grid-cols-[1.02fr_0.98fr]">
        <section className="overflow-hidden rounded-[36px] border border-[#e6d9ce] bg-[#fffdfa] shadow-[0_10px_28px_rgba(110,78,56,0.08)]">
          <div className="border-b border-[#eadfd5] bg-gradient-to-br from-[#4e342e] via-[#6f4e37] to-[#b08a6b] px-8 py-10 text-white md:px-10">
            <p className="mb-3 text-sm font-black tracking-[0.32em] text-white/80 uppercase">
              signup
            </p>

            <h1 className="mb-4 text-4xl font-black leading-tight md:text-5xl">
              新規登録して
              <br />
              Coffee SPAを使い始める
            </h1>
          </div>

          <div className="grid gap-4 px-8 py-8 md:px-10 md:py-10">
            <div className="rounded-[28px] border border-[#eadfd5] bg-white px-6 py-6">
              <p className="mb-2 text-xs font-black tracking-[0.24em] text-[#a1775b] uppercase">
                step 1
              </p>
              <h2 className="mb-2 text-xl font-black text-[#4e342e]">
                メールとパスワードを登録
              </h2>
              <p className="text-sm font-semibold leading-7 text-[#766b63]">
                まずは基本情報を登録して、認証フローを開始します。
              </p>
            </div>

            <div className="rounded-[28px] border border-[#eadfd5] bg-white px-6 py-6">
              <p className="mb-2 text-xs font-black tracking-[0.24em] text-[#a1775b] uppercase">
                step 2
              </p>
              <h2 className="mb-2 text-xl font-black text-[#4e342e]">
                メール確認を完了
              </h2>
              <p className="text-sm font-semibold leading-7 text-[#766b63]">
                backendログに出るverifyリンクを開いて、確認を完了させます。
              </p>
            </div>

            <div className="rounded-[28px] border border-[#eadfd5] bg-white px-6 py-6">
              <p className="mb-2 text-xs font-black tracking-[0.24em] text-[#a1775b] uppercase">
                step 3
              </p>
              <h2 className="mb-2 text-xl font-black text-[#4e342e]">
                ログインして利用開始
              </h2>
              <p className="text-sm font-semibold leading-7 text-[#766b63]">
                ログイン後はマイページに入り、アカウント情報や管理導線を確認できます。
              </p>
            </div>
          </div>
        </section>

        <section className="rounded-[36px] border border-[#e6d9ce] bg-white px-6 py-8 shadow-[0_10px_28px_rgba(110,78,56,0.08)] md:px-8 md:py-10">
          <div className="mb-8">
            <p className="mb-3 text-sm font-black tracking-[0.28em] text-[#a1775b] uppercase">
              create account
            </p>

            <h2 className="mb-3 text-3xl font-black text-[#4e342e]">
              アカウント登録
            </h2>

            <p className="text-base font-semibold leading-8 text-[#766b63]">
              登録後は確認メールの完了が必要です。そこで詰まらないよう、この画面で次の手順まで案内します。
            </p>
          </div>

          {msg ? (
            <div
              className={[
                "mb-6 rounded-[24px] px-5 py-4 text-sm font-bold leading-7",
                ok
                  ? "border border-[#cfe4cc] bg-[#f3fff1] text-[#42613e]"
                  : "border border-[#e6c7bd] bg-[#fff3ef] text-[#8a4b3a]",
              ].join(" ")}
            >
              {msg}
            </div>
          ) : null}

          <form onSubmit={onSubmit} className="grid gap-5">
            <label className="grid gap-2">
              <span className="text-sm font-black tracking-[0.08em] text-[#5f4a40]">
                メールアドレス
              </span>

              <input
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="you@example.com"
                type="email"
                autoComplete="email"
                className="w-full rounded-2xl border border-[#dccabc] bg-[#fffdfa] px-4 py-4 text-base font-semibold text-[#4e342e] outline-none transition placeholder:text-[#b09d90] focus:border-[#9c7257] focus:ring-4 focus:ring-[#ead8ca]"
              />
            </label>

            <label className="grid gap-2">
              <span className="text-sm font-black tracking-[0.08em] text-[#5f4a40]">
                パスワード
              </span>

              <input
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="new password"
                type="password"
                autoComplete="new-password"
                className="w-full rounded-2xl border border-[#dccabc] bg-[#fffdfa] px-4 py-4 text-base font-semibold text-[#4e342e] outline-none transition placeholder:text-[#b09d90] focus:border-[#9c7257] focus:ring-4 focus:ring-[#ead8ca]"
              />
            </label>

            <button
              type="submit"
              disabled={loading}
              className="mt-2 inline-flex items-center justify-center rounded-2xl bg-[#4e342e] px-5 py-4 text-base font-black text-white transition hover:opacity-90 disabled:cursor-not-allowed disabled:opacity-60"
            >
              {loading ? "登録中..." : "アカウントを作成する"}
            </button>
          </form>

          <div className="mt-8 grid gap-3 rounded-[26px] border border-[#eadfd5] bg-[#fcf8f4] px-5 py-5">
            <p className="text-sm font-black tracking-[0.24em] text-[#a1775b] uppercase">
              next action
            </p>

            <p className="text-sm font-semibold leading-7 text-[#766b63]">
              すでに登録済みならログインへ。確認メールをもう一度送りたい場合は再送ページへ進んでください。
            </p>

            <div className="flex flex-wrap gap-3">
              <Link
                to="/login"
                className="inline-flex rounded-full border border-[#d9c6b8] bg-white px-4 py-2.5 text-sm font-bold text-[#7b523a] transition hover:bg-[#f7efe8]"
              >
                ログインへ
              </Link>

              <Link
                to="/resend-verify"
                className="inline-flex rounded-full border border-[#d9c6b8] bg-white px-4 py-2.5 text-sm font-bold text-[#7b523a] transition hover:bg-[#f7efe8]"
              >
                確認メールを再送
              </Link>

              <Link
                to="/"
                className="inline-flex rounded-full border border-[#d9c6b8] bg-white px-4 py-2.5 text-sm font-bold text-[#7b523a] transition hover:bg-[#f7efe8]"
              >
                topへ戻る
              </Link>
            </div>
          </div>
        </section>
      </div>
    </main>
  );
}
