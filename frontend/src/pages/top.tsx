import { useEffect, useMemo, useState } from "react";
import { Link } from "react-router-dom";
import heroImg from "../assets/hero.png";
import { ApiError, api } from "../lib/api";

type ItemKind = "news" | "recipe" | "deal" | "shop";

type Item = {
  id: number;
  title: string;
  summary: string | null;
  url: string | null;
  image_url: string | null;
  kind: ItemKind;
  source_id: number;
  published_at: string;
  created_at: string;
};

type TopRes = {
  news: Item[];
  recipe: Item[];
  deal: Item[];
  shop: Item[];
};

type SectionProps = {
  title: string;
  sub: string;
  items: Item[];
  badge: string;
};

function toErrorMessage(err: unknown, fallback: string): string {
  if (err instanceof ApiError) {
    return err.message;
  }

  return fallback;
}

function formatDate(value: string): string {
  const d = new Date(value);

  if (Number.isNaN(d.getTime())) {
    return value;
  }

  return new Intl.DateTimeFormat("ja-JP", {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
  }).format(d);
}

function kindLabel(kind: ItemKind): string {
  switch (kind) {
    case "news":
      return "NEWS";
    case "recipe":
      return "RECIPE";
    case "deal":
      return "DEAL";
    case "shop":
      return "SHOP";
    default:
      return kind;
  }
}

function ItemCard({ item }: { item: Item }) {
  const body = (
    <>
      <div className="aspect-[16/10] w-full overflow-hidden bg-[#efe6de]">
        {item.image_url ? (
          <img
            src={item.image_url}
            alt={item.title}
            className="h-full w-full object-cover transition duration-300 group-hover:scale-[1.03]"
          />
        ) : (
          <div className="flex h-full items-center justify-center bg-gradient-to-br from-[#e8d7c8] to-[#d6b89c] text-sm font-semibold tracking-[0.2em] text-[#6a4731]">
            {kindLabel(item.kind)}
          </div>
        )}
      </div>

      <div className="flex min-h-[160px] flex-col items-center justify-center px-5 py-6 text-center">
        <div className="mb-3 flex items-center justify-center gap-3">
          <span className="rounded-full bg-[#f4ece5] px-3 py-1 text-[11px] font-bold tracking-[0.15em] text-[#6b452f]">
            {kindLabel(item.kind)}
          </span>
          <time className="text-xs text-stone-500">
            {formatDate(item.published_at)}
          </time>
        </div>

        <h3 className="mx-auto line-clamp-2 max-w-[18rem] text-lg font-bold leading-7 text-[#3f281a]">
          {item.title}
        </h3>

        {item.summary ? (
          <p className="mx-auto mt-3 line-clamp-2 max-w-[18rem] text-sm leading-6 text-stone-600">
            {item.summary}
          </p>
        ) : null}
      </div>
    </>
  );

  if (item.url) {
    return (
      <a
        href={item.url}
        target="_blank"
        rel="noreferrer"
        className="group block overflow-hidden rounded-3xl border border-[#e5d7cb] bg-white shadow-sm transition hover:-translate-y-1 hover:shadow-md"
      >
        {body}
      </a>
    );
  }

  return (
    <article className="overflow-hidden rounded-3xl border border-[#e5d7cb] bg-white shadow-sm opacity-90">
      {body}
    </article>
  );
}

function SectionBlock({ title, sub, items, badge }: SectionProps) {
  return (
    <section className="rounded-[2rem] border border-[#dccabd] bg-[#fffaf6] p-6 shadow-sm">
      <div className="mb-6 border-b border-[#eadfd6] pb-4 text-center">
        <div className="mb-3 inline-flex rounded-full bg-[#f1e6dd] px-4 py-2 text-[11px] font-bold tracking-[0.2em] text-[#7a5239]">
          {badge}
        </div>

        <h2 className="text-2xl font-bold tracking-wide text-[#452b1c]">
          {title}
        </h2>

        <p className="mt-3 text-sm leading-6 text-stone-600">{sub}</p>
      </div>

      {items.length === 0 ? (
        <div className="rounded-3xl border border-dashed border-[#d8c6b8] bg-white px-5 py-8 text-center text-sm text-stone-500">
          まだデータがありません。
        </div>
      ) : (
        <div className="grid gap-5 md:grid-cols-2 xl:grid-cols-4">
          {items.map((item) => (
            <ItemCard key={item.id} item={item} />
          ))}
        </div>
      )}
    </section>
  );
}

export function TopPage() {
  const [data, setData] = useState<TopRes | null>(null);
  const [msg, setMsg] = useState("");
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function run() {
      try {
        const res = await api<TopRes>("/items/top?limit=4", {
          method: "GET",
        });

        if (!res) {
          setMsg("データが空です。");
          return;
        }

        setData(res);
      } catch (err: unknown) {
        setMsg(toErrorMessage(err, "一覧の取得に失敗しました。"));
      } finally {
        setLoading(false);
      }
    }

    void run();
  }, []);

  const pickup = useMemo<Item[]>(() => {
    if (!data) {
      return [];
    }

    return [data.news[0], data.recipe[0], data.deal[0], data.shop[0]].filter(
      (item): item is Item => Boolean(item),
    );
  }, [data]);

  if (loading) {
    return (
      <div className="min-h-screen bg-[#f7f3ee] px-4 py-10 text-stone-700">
        <div className="mx-auto max-w-6xl rounded-3xl border border-[#e2d4c7] bg-white px-6 py-10 shadow-sm">
          loading...
        </div>
      </div>
    );
  }

  if (msg) {
    return (
      <div className="min-h-screen bg-[#f7f3ee] px-4 py-10 text-stone-700">
        <div className="mx-auto max-w-6xl rounded-3xl border border-[#e2d4c7] bg-white px-6 py-10 shadow-sm">
          {msg}
        </div>
      </div>
    );
  }

  if (!data) {
    return (
      <div className="min-h-screen bg-[#f7f3ee] px-4 py-10 text-stone-700">
        <div className="mx-auto max-w-6xl rounded-3xl border border-[#e2d4c7] bg-white px-6 py-10 shadow-sm">
          no data
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-[#f7f3ee] text-stone-800">
      <div className="mx-auto max-w-6xl px-4 py-6 md:px-6 md:py-8">
        <header className="mb-8 overflow-hidden rounded-[2rem] border border-[#d9c8bb] bg-white shadow-sm">
          <div className="grid gap-0 lg:grid-cols-[1.15fr_0.85fr]">
            <div className="bg-gradient-to-br from-[#4d2f1f] via-[#6a4731] to-[#9a7454] px-6 py-8 text-white md:px-8 md:py-10">
              <p className="text-sm font-semibold uppercase tracking-[0.25em] text-white/80">
                Coffee SPA
              </p>
              <h1 className="mt-3 text-3xl font-bold leading-tight md:text-5xl">
                コーヒーの情報を、
                <br />
                ひと目で整理して見られるトップページ
              </h1>

              <p className="mt-4 max-w-2xl text-sm leading-7 text-white/90 md:text-base">
                ニュース、レシピ、セール、店舗情報をカテゴリごとに整理しました。
                まず何が更新されたかをすぐ把握できる、情報ポータル寄りのトップです。
              </p>
              <div className="mt-6 flex flex-wrap gap-3">
                <a
                  href="#pickup"
                  className="rounded-full bg-white px-5 py-3 text-sm font-semibold text-[#4d2f1f] transition hover:opacity-90"
                >
                  注目情報を見る
                </a>

                <Link
                  to="/login"
                  className="rounded-full border border-white/40 px-5 py-3 text-sm font-semibold text-white transition hover:bg-white/10"
                >
                  ログイン
                </Link>

                <Link
                  to="/signup"
                  className="rounded-full border border-white/40 px-5 py-3 text-sm font-semibold text-white transition hover:bg-white/10"
                >
                  新規登録
                </Link>
              </div>

              <div className="mt-8 grid gap-3 sm:grid-cols-3">
                <div className="rounded-2xl bg-white/10 px-4 py-4 backdrop-blur-sm">
                  <p className="text-xs tracking-[0.15em] text-white/70">
                    NEWS
                  </p>
                  <p className="mt-2 text-sm font-semibold">
                    新着情報を整理して確認
                  </p>
                </div>
                <div className="rounded-2xl bg-white/10 px-4 py-4 backdrop-blur-sm">
                  <p className="text-xs tracking-[0.15em] text-white/70">
                    RECIPE
                  </p>
                  <p className="mt-2 text-sm font-semibold">
                    家でも試しやすい抽出レシピ
                  </p>
                </div>
                <div className="rounded-2xl bg-white/10 px-4 py-4 backdrop-blur-sm">
                  <p className="text-xs tracking-[0.15em] text-white/70">
                    SHOP / DEAL
                  </p>
                  <p className="mt-2 text-sm font-semibold">
                    店舗情報とお得情報をまとめて表示
                  </p>
                </div>
              </div>
            </div>

            <div className="relative min-h-[260px] bg-[#eadccf]">
              <img
                src={heroImg}
                alt="コーヒーのヒーロー画像"
                className="h-full w-full object-cover"
              />
              <div className="absolute inset-0 bg-gradient-to-t from-[#2f1b10]/30 via-transparent to-transparent" />
            </div>
          </div>

          <div className="border-t border-[#eadfd6] bg-[#fffaf6] px-6 py-4">
            <div className="flex flex-wrap items-center gap-3 text-sm">
              <span className="font-semibold text-[#5c3925]">
                クイックリンク
              </span>
              <Link
                to="/me"
                className="rounded-full border border-[#d6c1af] bg-white px-4 py-2 font-medium text-[#6a4731] hover:bg-[#f7efe8]"
              >
                me
              </Link>
              <Link
                to="/admin"
                className="rounded-full border border-[#d6c1af] bg-white px-4 py-2 font-medium text-[#6a4731] hover:bg-[#f7efe8]"
              >
                admin
              </Link>
              <a
                href="#news"
                className="rounded-full border border-[#d6c1af] bg-white px-4 py-2 font-medium text-[#6a4731] hover:bg-[#f7efe8]"
              >
                ニュース
              </a>
              <a
                href="#recipe"
                className="rounded-full border border-[#d6c1af] bg-white px-4 py-2 font-medium text-[#6a4731] hover:bg-[#f7efe8]"
              >
                レシピ
              </a>
              <a
                href="#deal"
                className="rounded-full border border-[#d6c1af] bg-white px-4 py-2 font-medium text-[#6a4731] hover:bg-[#f7efe8]"
              >
                セール
              </a>
              <a
                href="#shop"
                className="rounded-full border border-[#d6c1af] bg-white px-4 py-2 font-medium text-[#6a4731] hover:bg-[#f7efe8]"
              >
                店舗
              </a>
            </div>
          </div>
        </header>

        <section
          id="pickup"
          className="mb-8 rounded-[2rem] border border-[#dccabd] bg-white px-6 py-6 shadow-sm"
        >
          <div className="mb-5 border-b border-[#eadfd6] pb-4 text-center">
            <div className="mb-3 inline-flex rounded-full bg-[#f1e6dd] px-4 py-2 text-[11px] font-bold tracking-[0.2em] text-[#7a5239]">
              PICKUP
            </div>
            <h2 className="text-2xl font-bold tracking-wide text-[#452b1c]">
              今日の注目トピック
            </h2>
            <p className="mt-3 text-sm leading-6 text-stone-600">
              各カテゴリの注目記事をまとめて確認できます。
            </p>
          </div>

          <div className="grid gap-5 md:grid-cols-2 xl:grid-cols-4">
            {pickup.map((item) => (
              <ItemCard key={item.id} item={item} />
            ))}
          </div>
        </section>

        <div className="space-y-8">
          <div id="news">
            <SectionBlock
              title="ニュース"
              sub="豆、器具、トレンドなどの更新情報をまとめています。"
              items={data.news}
              badge="NEWS"
            />
          </div>

          <div id="recipe">
            <SectionBlock
              title="レシピ"
              sub="家でも再現しやすい淹れ方や抽出手順を整理しています。"
              items={data.recipe}
              badge="RECIPE"
            />
          </div>

          <div id="deal">
            <SectionBlock
              title="セール"
              sub="クーポン、送料無料、期間限定割引などのお得情報です。"
              items={data.deal}
              badge="DEAL"
            />
          </div>

          <div id="shop">
            <SectionBlock
              title="店舗"
              sub="新店舗や立ち寄りやすいカフェ情報を確認できます。"
              items={data.shop}
              badge="SHOP"
            />
          </div>
        </div>

        <footer className="mt-8 rounded-[2rem] border border-[#dccabd] bg-[#fffaf6] px-6 py-6 shadow-sm">
          <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
            <div>
              <p className="text-lg font-bold text-[#452b1c]">Coffee SPA</p>
              <p className="mt-1 text-sm text-stone-600">
                コーヒー情報を見やすく整理するためのシングルページです。
              </p>
            </div>

            <div className="flex flex-wrap gap-3 text-sm">
              <Link
                to="/login"
                className="rounded-full border border-[#d6c1af] bg-white px-4 py-2 font-medium text-[#6a4731] hover:bg-[#f7efe8]"
              >
                login
              </Link>
              <Link
                to="/signup"
                className="rounded-full border border-[#d6c1af] bg-white px-4 py-2 font-medium text-[#6a4731] hover:bg-[#f7efe8]"
              >
                signup
              </Link>
              <Link
                to="/me"
                className="rounded-full border border-[#d6c1af] bg-white px-4 py-2 font-medium text-[#6a4731] hover:bg-[#f7efe8]"
              >
                me
              </Link>
              <Link
                to="/admin"
                className="rounded-full border border-[#d6c1af] bg-white px-4 py-2 font-medium text-[#6a4731] hover:bg-[#f7efe8]"
              >
                admin
              </Link>
            </div>
          </div>
        </footer>
      </div>
    </div>
  );
}
