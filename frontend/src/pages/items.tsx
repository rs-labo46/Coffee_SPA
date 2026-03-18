import { useEffect, useMemo, useState } from "react";
import { Link, useSearchParams } from "react-router-dom";
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

type ItemListRes = {
  items: Item[];
};

function toErrorMessage(err: unknown, fallback: string): string {
  if (err instanceof ApiError) {
    return err.message;
  }
  return fallback;
}

function isItemKind(v: string | null): v is ItemKind {
  return v === "news" || v === "recipe" || v === "deal" || v === "shop";
}

function kindLabel(kind: ItemKind): string {
  switch (kind) {
    case "news":
      return "ニュース";
    case "recipe":
      return "レシピ";
    case "deal":
      return "セール";
    case "shop":
      return "店舗";
    default:
      return "";
  }
}

function kindDesc(kind: ItemKind): string {
  switch (kind) {
    case "news":
      return "上位3件を含めたニュース一覧です。";
    case "recipe":
      return "上位3件を含めたレシピ一覧です。";
    case "deal":
      return "上位3件を含めたセール一覧です。";
    case "shop":
      return "上位3件を含めた店舗一覧です。";
    default:
      return "";
  }
}

function fmtDate(v: string): string {
  const d = new Date(v);
  if (Number.isNaN(d.getTime())) {
    return "";
  }

  const y = d.getFullYear();
  const m = String(d.getMonth() + 1).padStart(2, "0");
  const day = String(d.getDate()).padStart(2, "0");
  return `${y}/${m}/${day}`;
}

function cardImage(url: string | null): string {
  if (url && url.trim() !== "") {
    return url;
  }
  return "https://images.unsplash.com/photo-1509042239860-f550ce710b93?auto=format&fit=crop&w=1200&q=80";
}

function hasRef(url: string | null): boolean {
  return !!url && url.trim() !== "";
}

function ItemCard({ item }: { item: Item }) {
  return (
    <article className="group flex h-full flex-col overflow-hidden rounded-[28px] border border-[#e5d7cb] bg-white shadow-[0_6px_20px_rgba(93,64,55,0.08)]">
      <div className="aspect-[16/10] overflow-hidden">
        <img
          src={cardImage(item.image_url)}
          alt={item.title}
          className="h-full w-full object-cover transition duration-300 group-hover:scale-[1.01]"
        />
      </div>

      <div className="flex flex-1 flex-col px-6 py-5">
        <div className="mb-3 flex items-center gap-3">
          <span className="rounded-full bg-[#f4ebe3] px-3 py-1.5 text-[11px] font-bold tracking-[0.28em] text-[#7b523a]">
            {item.kind.toUpperCase()}
          </span>

          <span className="text-sm font-semibold text-[#8d8178]">
            {fmtDate(item.published_at)}
          </span>
        </div>

        <h3 className="mb-3 line-clamp-2 text-xl font-extrabold leading-tight text-[#4e342e]">
          {item.title}
        </h3>

        <p className="mb-5 line-clamp-3 text-sm font-semibold leading-7 text-[#6d625b]">
          {item.summary ?? "概要は未登録です。"}
        </p>

        <div className="mt-auto">
          <span className="inline-flex items-center rounded-full border border-[#eadfd5] bg-[#faf5f0] px-4 py-2 text-sm font-bold text-[#9b8b7f]">
            {hasRef(item.url) ? "参考記事あり" : "記事カード"}
          </span>
        </div>
      </div>
    </article>
  );
}

export function ItemsPage() {
  const [params] = useSearchParams();
  const kindParam = params.get("kind");
  const kind: ItemKind = isItemKind(kindParam) ? kindParam : "news";

  const [items, setItems] = useState<Item[]>([]);
  const [loading, setLoading] = useState(true);
  const [msg, setMsg] = useState("");

  useEffect(() => {
    async function run() {
      setLoading(true);
      setMsg("");

      try {
        const res = await api<ItemListRes>(
          `/items?kind=${kind}&limit=10&offset=0`,
          { method: "GET" },
        );

        setItems(res?.items ?? []);
      } catch (err: unknown) {
        setMsg(toErrorMessage(err, "一覧の取得に失敗しました。"));
      } finally {
        setLoading(false);
      }
    }

    void run();
  }, [kind]);

  const tabs = useMemo(
    () => [
      { kind: "news" as const, label: "ニュース" },
      { kind: "recipe" as const, label: "レシピ" },
      { kind: "deal" as const, label: "セール" },
      { kind: "shop" as const, label: "店舗" },
    ],
    [],
  );

  if (loading) {
    return (
      <main className="min-h-screen bg-[#f6f1eb] px-4 py-8 md:px-8 md:py-10">
        <div className="mx-auto max-w-[1280px] rounded-[32px] border border-[#eadfd4] bg-white px-8 py-16 text-center text-lg font-bold text-[#6f6259] shadow-[0_8px_24px_rgba(110,78,56,0.06)]">
          読み込み中です...
        </div>
      </main>
    );
  }

  if (msg) {
    return (
      <main className="min-h-screen bg-[#f6f1eb] px-4 py-8 md:px-8 md:py-10">
        <div className="mx-auto max-w-[1280px] rounded-[32px] border border-[#eadfd4] bg-white px-8 py-16 text-center text-lg font-bold text-[#8a4b3a] shadow-[0_8px_24px_rgba(110,78,56,0.06)]">
          {msg}
        </div>
      </main>
    );
  }

  return (
    <main className="min-h-screen bg-[#f6f1eb] px-4 py-8 md:px-8 md:py-10">
      <div className="mx-auto flex max-w-[1280px] flex-col gap-8">
        <section className="rounded-[36px] border border-[#e6d9ce] bg-gradient-to-br from-[#fffdfa] via-[#fbf5ef] to-[#f3e6d9] px-7 py-8 shadow-[0_8px_24px_rgba(110,78,56,0.06)] md:px-10">
          <p className="mb-3 text-sm font-black tracking-[0.32em] text-[#a1775b] uppercase">
            item list
          </p>

          <h1 className="mb-3 text-3xl font-black text-[#4e342e] md:text-5xl">
            {kindLabel(kind)}一覧
          </h1>

          <div className="mt-6 flex flex-wrap gap-3">
            {tabs.map((tab) => (
              <Link
                key={tab.kind}
                to={`/items?kind=${tab.kind}`}
                className={[
                  "rounded-full border px-5 py-2.5 text-sm font-bold transition",
                  tab.kind === kind
                    ? "border-[#7b523a] bg-[#7b523a] text-white"
                    : "border-[#d9c6b8] bg-white text-[#7b523a] hover:bg-[#f7efe8]",
                ].join(" ")}
              >
                {tab.label}
              </Link>
            ))}

            <Link
              to="/"
              className="rounded-full border border-[#d9c6b8] bg-white px-5 py-2.5 text-sm font-bold text-[#7b523a] transition hover:bg-[#f7efe8]"
            >
              トップへ戻る
            </Link>
          </div>
        </section>

        {items.length === 0 ? (
          <section className="rounded-[32px] border border-[#eadfd4] bg-white px-8 py-16 text-center shadow-[0_8px_24px_rgba(110,78,56,0.06)]">
            <p className="text-base font-bold text-[#7a6f68]">
              表示できるデータがありません。
            </p>
          </section>
        ) : (
          <section className="grid gap-5 md:grid-cols-2 xl:grid-cols-3">
            {items.map((item) => (
              <ItemCard key={item.id} item={item} />
            ))}
          </section>
        )}
      </div>
    </main>
  );
}
