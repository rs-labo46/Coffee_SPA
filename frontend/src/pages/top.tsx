import { useEffect, useState } from "react";
import { ApiError, api } from "../lib/api";

type ItemKind = "news" | "recipe" | "deal" | "shop";

type Item = {
  id: number;
  title: string;
  summary: string;
  url: string;
  image_url: string;
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
          setMsg("データが空です");
          return;
        }

        setData(res);
      } catch (err: unknown) {
        if (err instanceof ApiError) {
          setMsg(err.message);
        } else {
          setMsg("一覧の取得に失敗しました");
        }
      } finally {
        setLoading(false);
      }
    }

    void run();
  }, []);

  if (loading) {
    return <div>loading...</div>;
  }

  if (msg) {
    return <div>{msg}</div>;
  }

  if (!data) {
    return <div>no data</div>;
  }

  return (
    <div>
      <h1>top</h1>

      <section>
        <h2>news</h2>
        {data.news.map((item) => (
          <article key={item.id}>
            <a href={item.url} target="_blank" rel="noreferrer">
              {item.title}
            </a>
            <p>{item.summary}</p>
          </article>
        ))}
      </section>

      <section>
        <h2>recipe</h2>
        {data.recipe.map((item) => (
          <article key={item.id}>
            <a href={item.url} target="_blank" rel="noreferrer">
              {item.title}
            </a>
            <p>{item.summary}</p>
          </article>
        ))}
      </section>

      <section>
        <h2>deal</h2>
        {data.deal.map((item) => (
          <article key={item.id}>
            <a href={item.url} target="_blank" rel="noreferrer">
              {item.title}
            </a>
            <p>{item.summary}</p>
          </article>
        ))}
      </section>

      <section>
        <h2>shop</h2>
        {data.shop.map((item) => (
          <article key={item.id}>
            <a href={item.url} target="_blank" rel="noreferrer">
              {item.title}
            </a>
            <p>{item.summary}</p>
          </article>
        ))}
      </section>
    </div>
  );
}
