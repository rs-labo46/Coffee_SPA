import { Link, useParams } from "react-router-dom";
import { useEffect, useMemo, useState } from "react";
import { api, toErrorMessage } from "../lib/api";

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

type ItemDetailRes = {
  item: Item;
};

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

function kindTitle(kind: ItemKind): string {
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

function imageSrc(url: string | null): string {
  if (url && url.trim() !== "") {
    return url;
  }

  return "https://images.unsplash.com/photo-1509042239860-f550ce710b93?auto=format&fit=crop&w=1600&q=80";
}

export default function ItemDetailPage() {
  const { id } = useParams();
  const [item, setItem] = useState<Item | null>(null);
  const [loading, setLoading] = useState(true);
  const [msg, setMsg] = useState("");

  const itemID = useMemo(() => Number(id), [id]);

  useEffect(() => {
    async function run() {
      try {
        if (!Number.isInteger(itemID) || itemID <= 0) {
          setMsg("記事IDが不正です。");
          return;
        }

        const res = await api<ItemDetailRes>(`/items/${itemID}`, {
          method: "GET",
        });

        if (!res || !res.item) {
          setMsg("記事データの取得に失敗しました。");
          return;
        }

        setItem(res.item);
      } catch (err: unknown) {
        setMsg(toErrorMessage(err, "記事詳細の取得に失敗しました。"));
      } finally {
        setLoading(false);
      }
    }

    void run();
  }, [itemID]);

  if (loading) {
    return (
      <main className="min-h-screen bg-[#f6f1eb] px-4 py-10 md:px-8">
        <div className="mx-auto max-w-[1100px] rounded-[32px] border border-[#eadfd4] bg-white px-8 py-16 text-center text-lg font-bold text-[#6f6259] shadow-[0_8px_24px_rgba(110,78,56,0.06)]">
          読み込み中です...
        </div>
      </main>
    );
  }

  if (msg) {
    return (
      <main className="min-h-screen bg-[#f6f1eb] px-4 py-10 md:px-8">
        <div className="mx-auto max-w-[1100px] rounded-[32px] border border-[#eadfd4] bg-white px-8 py-16 shadow-[0_8px_24px_rgba(110,78,56,0.06)]">
          <p className="mb-6 text-lg font-bold text-[#8a4b3a]">{msg}</p>

          <Link
            to="/"
            className="inline-flex rounded-full border border-[#d9c6b8] bg-white px-6 py-3 text-base font-bold text-[#7b523a] transition hover:bg-[#f7efe8]"
          >
            トップへ戻る
          </Link>
        </div>
      </main>
    );
  }

  if (!item) {
    return null;
  }

  return (
    <main className="min-h-screen bg-[#f6f1eb] px-4 py-10 md:px-8">
      <div className="mx-auto flex max-w-[1100px] flex-col gap-8">
        <div>
          <Link
            to="/"
            className="inline-flex rounded-full border border-[#d9c6b8] bg-white px-6 py-3 text-base font-bold text-[#7b523a] transition hover:bg-[#f7efe8]"
          >
            ← トップへ戻る
          </Link>
        </div>

        <article className="overflow-hidden rounded-[36px] border border-[#e6d9ce] bg-[#fffdfa] shadow-[0_8px_24px_rgba(110,78,56,0.06)]">
          <div className="aspect-[16/8] overflow-hidden">
            <img
              src={imageSrc(item.image_url)}
              alt={item.title}
              className="h-full w-full object-cover"
            />
          </div>

          <div className="px-7 py-8 md:px-12 md:py-12">
            <div className="mb-5 flex flex-wrap items-center gap-3">
              <span className="rounded-full bg-[#f4ebe3] px-4 py-2 text-xs font-bold tracking-[0.28em] text-[#7b523a]">
                {kindLabel(item.kind)}
              </span>

              <span className="text-sm font-semibold text-[#8d8178]">
                {fmtDate(item.published_at)}
              </span>

              <span className="text-sm font-semibold text-[#8d8178]">
                カテゴリ: {kindTitle(item.kind)}
              </span>
            </div>

            <h1 className="mb-6 text-[34px] font-black leading-tight text-[#4e342e] md:text-[42px]">
              {item.title}
            </h1>

            <p className="mb-8 text-lg font-semibold leading-9 text-[#6d625b]">
              {item.summary ?? "概要は未登録です。"}
            </p>

            {item.url ? (
              <div className="mt-8">
                <a
                  href={item.url}
                  target="_blank"
                  rel="noreferrer"
                  className="inline-flex rounded-full bg-[#6f4e37] px-6 py-3 text-base font-bold text-white transition hover:opacity-90"
                >
                  元URLを開く
                </a>
              </div>
            ) : null}
          </div>
        </article>
      </div>
    </main>
  );
}
