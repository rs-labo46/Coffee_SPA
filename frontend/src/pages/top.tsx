import { useEffect, useMemo, useState } from "react";
import { Link } from "react-router-dom";
import { api, toErrorMessage } from "../lib/api";
import {
  cardImage,
  hasRef,
  kindBadgeLabel,
  kindTitleLabel,
  type ItemKind,
} from "../lib/item";
import { formatDisplayDate } from "../lib/date";

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

function sectionDesc(kind: ItemKind): string {
  switch (kind) {
    case "news":
      return "豆、器具、焙煎、業界トレンドの更新をまとめています。";
    case "recipe":
      return "自宅で再現しやすい抽出レシピや淹れ方のコツを掲載しています。";
    case "deal":
      return "クーポン、セール、期間限定キャンペーンを確認できます。";
    case "shop":
      return "新店舗や気になるコーヒーショップの情報をまとめています。";
    default:
      return "";
  }
}

function sectionListPath(kind: ItemKind): string {
  return `/items?kind=${kind}`;
}

function ItemCard({ item }: { item: Item }) {
  const ref = hasRef(item.url);

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
            {kindBadgeLabel(item.kind)}
          </span>

          <span className="text-sm font-semibold text-[#8d8178]">
            {formatDisplayDate(item.published_at)}
          </span>
        </div>

        <h3 className="mb-3 line-clamp-2 text-xl font-extrabold leading-tight text-[#4e342e]">
          {item.title}
        </h3>

        <p className="mb-5 line-clamp-3 text-sm font-semibold leading-7 text-[#6d625b]">
          {item.summary ?? "概要は未登録です。"}
        </p>

        <div className="mt-auto">
          {ref ? (
            <span className="inline-flex items-center rounded-full border border-[#d9c6b8] bg-[#fcf7f2] px-4 py-2 text-sm font-bold text-[#7b523a]">
              参考記事あり
            </span>
          ) : (
            <span className="inline-flex items-center rounded-full border border-[#eadfd5] bg-[#faf5f0] px-4 py-2 text-sm font-bold text-[#9b8b7f]">
              記事カード
            </span>
          )}
        </div>
      </div>
    </article>
  );
}

function SectionBlock({ kind, items }: { kind: ItemKind; items: Item[] }) {
  const recent = items.slice(0, 3);

  return (
    <section
      id={kind}
      className="rounded-[36px] border border-[#e6d9ce] bg-[#fffdfa] px-6 py-7 shadow-[0_8px_24px_rgba(110,78,56,0.06)] md:px-8 md:py-8"
    >
      <div className="mb-6 flex flex-col gap-4 border-b border-[#eadfd5] pb-6 md:flex-row md:items-end md:justify-between">
        <div>
          <div className="mb-3 inline-flex rounded-full bg-[#f3e8de] px-4 py-2 text-xs font-bold tracking-[0.32em] text-[#7b523a]">
            {kindBadgeLabel(kind)}
          </div>

          <h2 className="mb-2 text-[26px] font-black text-[#4e342e] md:text-[30px]">
            {kindTitleLabel(kind)}
          </h2>

          <p className="max-w-3xl text-sm font-semibold leading-7 text-[#766b63]">
            {sectionDesc(kind)}
          </p>
        </div>

        <Link
          to={sectionListPath(kind)}
          className="inline-flex items-center rounded-full border border-[#d9c6b8] bg-white px-4 py-2 text-sm font-bold text-[#7b523a] transition hover:bg-[#f7efe8]"
        >
          一覧へ
        </Link>
      </div>

      <div className="grid gap-5 lg:grid-cols-[1.1fr_1.9fr]">
        <div className="rounded-[28px] border border-[#eadfd5] bg-[#fcf8f4] px-5 py-5">
          <p className="mb-3 text-xs font-black tracking-[0.24em] text-[#a1775b] uppercase">
            latest 3
          </p>

          <div className="grid gap-3">
            {recent.map((item, idx) => (
              <div
                key={item.id}
                className="rounded-2xl border border-[#eadfd5] bg-white px-4 py-4"
              >
                <div className="mb-2 flex items-center gap-3">
                  <span className="inline-flex h-7 w-7 items-center justify-center rounded-full bg-[#6f4e37] text-xs font-black text-white">
                    {idx + 1}
                  </span>

                  <span className="text-xs font-bold tracking-[0.18em] text-[#9b7a66] uppercase">
                    {formatDisplayDate(item.published_at)}
                  </span>
                </div>

                <p className="mb-2 line-clamp-2 text-sm font-black leading-6 text-[#4e342e]">
                  {item.title}
                </p>

                <p className="text-xs font-bold text-[#8a7b71]">
                  {hasRef(item.url) ? "参考記事あり" : "記事カード"}
                </p>
              </div>
            ))}
          </div>
        </div>

        <div className="grid gap-5 md:grid-cols-2 xl:grid-cols-3">
          {recent.map((item) => (
            <ItemCard key={item.id} item={item} />
          ))}
        </div>
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
        const res = await api<TopRes>("/items/top?limit=3", {
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
      <div className="mx-auto flex max-w-[1280px] flex-col gap-8">
        <section className="rounded-[36px] border border-[#e6d9ce] bg-gradient-to-br from-[#fffdfa] via-[#fbf5ef] to-[#f3e6d9] px-7 py-8 shadow-[0_8px_24px_rgba(110,78,56,0.06)] md:px-10">
          <div className="grid gap-8 lg:grid-cols-[1.4fr_0.9fr]">
            <div>
              <p className="mb-3 text-sm font-black tracking-[0.32em] text-[#a1775b] uppercase">
                coffee portal
              </p>

              <h2 className="mb-4 text-3xl font-black leading-tight text-[#4e342e] md:text-5xl">
                豆・抽出・セール・店舗の最新情報をお届け
              </h2>

              <div className="mt-6 flex flex-wrap gap-3">
                {quickLinks.map((link) => (
                  <a
                    key={link.label}
                    href={link.to}
                    className="rounded-full border border-[#d9c6b8] bg-white px-5 py-2.5 text-sm font-bold text-[#7b523a] transition hover:bg-[#f7efe8]"
                  >
                    {link.label}
                  </a>
                ))}
              </div>
            </div>
          </div>
        </section>

        <SectionBlock kind="news" items={data.news} />
        <SectionBlock kind="recipe" items={data.recipe} />
        <SectionBlock kind="deal" items={data.deal} />
        <SectionBlock kind="shop" items={data.shop} />
      </div>
    </main>
  );
}
