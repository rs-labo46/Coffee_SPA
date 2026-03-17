import { useEffect, useMemo, useState } from "react";
import { Link } from "react-router-dom";
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

function toErrorMessage(err: unknown, fallback: string): string {
  if (err instanceof ApiError) {
    return err.message;
  }

  return fallback;
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
      return "";
  }
}

function sectionTitle(kind: ItemKind): string {
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

function sectionDesc(kind: ItemKind): string {
  switch (kind) {
    case "news":
      return "豆、器具、トレンドなどの更新情報をまとめています。";
    case "recipe":
      return "自宅で試せる抽出レシピやコツを確認できます。";
    case "deal":
      return "お得なセール・キャンペーン情報を掲載しています。";
    case "shop":
      return "新店舗や注目ショップの情報をまとめています。";
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

function ItemCard({ item }: { item: Item }) {
  return (
    <Link
      to={`/items/${item.id}`}
      className="group flex h-full flex-col overflow-hidden rounded-[28px] border border-[#e5d7cb] bg-white shadow-[0_6px_20px_rgba(93,64,55,0.08)] transition duration-200 hover:-translate-y-1 hover:shadow-[0_12px_32px_rgba(93,64,55,0.14)]"
    >
      <div className="aspect-[16/10] overflow-hidden">
        <img
          src={cardImage(item.image_url)}
          alt={item.title}
          className="h-full w-full object-cover transition duration-300 group-hover:scale-[1.03]"
        />
      </div>

      <div className="flex flex-1 flex-col px-7 py-6">
        <div className="mb-4 flex items-center gap-3">
          <span className="rounded-full bg-[#f4ebe3] px-4 py-2 text-xs font-bold tracking-[0.28em] text-[#7b523a]">
            {kindLabel(item.kind)}
          </span>
          <span className="text-sm font-semibold text-[#8d8178]">
            {fmtDate(item.published_at)}
          </span>
        </div>

        <h3 className="mb-4 line-clamp-2 text-[22px] font-extrabold leading-tight text-[#4e342e]">
          {item.title}
        </h3>

        <p className="line-clamp-3 text-base font-semibold leading-8 text-[#6d625b]">
          {item.summary ?? "詳細は記事ページで確認できます。"}
        </p>
      </div>
    </Link>
  );
}

function SectionBlock({ kind, items }: { kind: ItemKind; items: Item[] }) {
  return (
    <section className="rounded-[36px] border border-[#e6d9ce] bg-[#fffdfa] px-7 py-7 shadow-[0_8px_24px_rgba(110,78,56,0.06)] md:px-10 md:py-8">
      <div className="mb-8 text-center">
        <div className="mb-5 inline-flex rounded-full bg-[#f3e8de] px-6 py-3 text-sm font-bold tracking-[0.32em] text-[#7b523a]">
          {kindLabel(kind)}
        </div>

        <h2 className="mb-4 text-[28px] font-black text-[#4e342e] md:text-[32px]">
          {sectionTitle(kind)}
        </h2>

        <p className="mx-auto max-w-3xl text-base font-semibold leading-8 text-[#766b63]">
          {sectionDesc(kind)}
        </p>
      </div>

      <div className="mb-8 border-t border-[#eadfd5]" />

      <div className="grid gap-6 md:grid-cols-2 xl:grid-cols-4">
        {items.map((item) => (
          <ItemCard key={item.id} item={item} />
        ))}
      </div>
    </section>
  );
}

export default function TopPage() {
  const [data, setData] = useState<TopRes | null>(null);
  const [loading, setLoading] = useState(true);
  const [msg, setMsg] = useState("");

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

  const quickLinks = useMemo(
    () => [
      { label: "admin", to: "/admin" },
      { label: "ニュース", to: "#news" },
      { label: "レシピ", to: "#recipe" },
      { label: "セール", to: "#deal" },
      { label: "店舗", to: "#shop" },
    ],
    [],
  );

  if (loading) {
    return (
      <main className="min-h-screen bg-[#f6f1eb] px-4 py-10 md:px-8">
        <div className="mx-auto max-w-[1280px] rounded-[32px] border border-[#eadfd4] bg-white px-8 py-16 text-center text-lg font-bold text-[#6f6259] shadow-[0_8px_24px_rgba(110,78,56,0.06)]">
          読み込み中です...
        </div>
      </main>
    );
  }

  if (msg) {
    return (
      <main className="min-h-screen bg-[#f6f1eb] px-4 py-10 md:px-8">
        <div className="mx-auto max-w-[1280px] rounded-[32px] border border-[#eadfd4] bg-white px-8 py-16 text-center text-lg font-bold text-[#8a4b3a] shadow-[0_8px_24px_rgba(110,78,56,0.06)]">
          {msg}
        </div>
      </main>
    );
  }

  if (!data) {
    return null;
  }

  return (
    <main className="min-h-screen bg-[#f6f1eb] px-4 py-8 md:px-8 md:py-10">
      <div className="mx-auto flex max-w-[1280px] flex-col gap-10">
        <section className="rounded-[36px] border border-[#e6d9ce] bg-[#fffdfa] px-7 py-6 shadow-[0_8px_24px_rgba(110,78,56,0.06)] md:px-10">
          <div className="flex flex-wrap items-center gap-4">
            {quickLinks.map((link) => (
              <a
                key={link.label}
                href={link.to}
                className="rounded-full border border-[#d9c6b8] bg-white px-6 py-3 text-xl font-bold text-[#7b523a] transition hover:bg-[#f7efe8]"
              >
                {link.label}
              </a>
            ))}
          </div>
        </section>

        <div id="news">
          <SectionBlock kind="news" items={data.news} />
        </div>

        <div id="recipe">
          <SectionBlock kind="recipe" items={data.recipe} />
        </div>

        <div id="deal">
          <SectionBlock kind="deal" items={data.deal} />
        </div>

        <div id="shop">
          <SectionBlock kind="shop" items={data.shop} />
        </div>
      </div>
    </main>
  );
}
