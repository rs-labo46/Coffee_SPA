import { useEffect, useMemo, useState } from "react";
import { Link } from "react-router-dom";
import { ApiError, api } from "../lib/api";
import { japanDateTime } from "../lib/date";

type ItemKind = "news" | "recipe" | "deal" | "shop";

type Source = {
  id: number;
  name: string;
  site_url: string | null;
  created_at: string;
};

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

type SourceListRes = {
  sources: Source[];
};

type SourceRes = {
  source: Source;
};

type ItemRes = {
  item: Item;
};
type SourceForm = {
  name: string;
  site_url: string;
};
type ItemForm = {
  title: string;
  summary: string;
  url: string;
  image_url: string;
  kind: ItemKind;
  source_id: string;
  published_at: string;
};

function newSourceForm(): SourceForm {
  return {
    name: "",
    site_url: "",
  };
}
function newItemForm(): ItemForm {
  return {
    title: "",
    summary: "",
    url: "",
    image_url: "",
    kind: "news",
    source_id: "",
    published_at: "",
  };
}
function toNullableText(value: string): string | null {
  const trimmed = value.trim();

  if (!trimmed) {
    return null;
  }

  return trimmed;
}
function toApiDateTime(value: string): string {
  if (!value) {
    return "";
  }
  return `${value}:00+09:00`;
}

function toErrorMessage(err: unknown, fallback: string): string {
  if (err instanceof ApiError) {
    return err.message;
  }

  return fallback;
}
function SectionTitle({ title, sub }: { title: string; sub: string }) {
  return (
    <div className="mb-4 border-b border-[#d9c7b8] pb-3">
      <h2 className="text-xl font-bold tracking-wide text-[#4f2f1d]">
        {title}
      </h2>
      <p className="mt-1 text-sm text-stone-600">{sub}</p>
    </div>
  );
}
function FieldLabel({
  htmlFor,
  label,
  required,
}: {
  htmlFor: string;
  label: string;
  required?: boolean;
}) {
  return (
    <label
      htmlFor={htmlFor}
      className="mb-1 block text-sm font-semibold text-stone-700"
    >
      {label}
      {required ? <span className="ml-1 text-red-700">*</span> : null}
    </label>
  );
}

export function AdminPage() {
  const [sources, setSources] = useState<Source[]>([]);
  const [loadingSources, setLoadingSources] = useState<boolean>(true);
  const [sourcesMsg, setSourcesMsg] = useState<string>("");
  const [sourceForm, setSourceForm] = useState<SourceForm>(newSourceForm());
  const [savingSource, setSavingSource] = useState<boolean>(false);
  const [sourceMsg, setSourceMsg] = useState<string>("");
  const [itemForm, setItemForm] = useState<ItemForm>(newItemForm());
  const [savingItem, setSavingItem] = useState<boolean>(false);
  const [itemMsg, setItemMsg] = useState<string>("");
  const [createdItem, setCreatedItem] = useState<Item | null>(null);
  const firstSourceId = useMemo<string>(() => {
    if (sources.length === 0) {
      return "";
    }

    return String(sources[0].id);
  }, [sources]);

  // source一覧を取得する関数。
  async function loadSources(): Promise<void> {
    setLoadingSources(true);
    setSourcesMsg("");

    try {
      const res = await api<SourceListRes>("/sources", {
        method: "GET",
      });

      if (!res) {
        setSources([]);
        setSourcesMsg("source一覧が空です。");
        return;
      }

      setSources(res.sources);

      //item formのsource_idが未選択なら、先頭を自動セットする。
      setItemForm((prev) => {
        if (prev.source_id) {
          return prev;
        }

        if (res.sources.length === 0) {
          return prev;
        }

        return {
          ...prev,
          source_id: String(res.sources[0].id),
        };
      });
    } catch (err: unknown) {
      setSources([]);
      setSourcesMsg(toErrorMessage(err, "source一覧の取得に失敗しました。"));
    } finally {
      setLoadingSources(false);
    }
  }

  // 初回表示時に一覧を取得する。
  useEffect(() => {
    void loadSources();
  }, []);

  useEffect(() => {
    if (!itemForm.source_id && firstSourceId) {
      setItemForm((prev) => ({
        ...prev,
        source_id: firstSourceId,
      }));
    }
  }, [firstSourceId, itemForm.source_id]);

  function onChangeSourceForm(e: React.ChangeEvent<HTMLInputElement>): void {
    const { name, value } = e.target;

    setSourceForm((prev) => ({
      ...prev,
      [name]: value,
    }));
  }

  //itemformの更新関数。
  function onChangeItemForm(
    e:
      | React.ChangeEvent<HTMLInputElement>
      | React.ChangeEvent<HTMLTextAreaElement>
      | React.ChangeEvent<HTMLSelectElement>,
  ): void {
    const { name, value } = e.target;

    setItemForm((prev) => ({
      ...prev,
      [name]: value,
    }));
  }
  async function onSubmitSource(
    e: React.FormEvent<HTMLFormElement>,
  ): Promise<void> {
    e.preventDefault();

    setSourceMsg("");
    setSavingSource(true);

    try {
      const body = {
        name: sourceForm.name.trim(),
        site_url: toNullableText(sourceForm.site_url),
      };

      const res = await api<SourceRes>("/sources", {
        method: "POST",
        auth: true,
        body,
      });

      if (res?.source) {
        setSourceMsg(`source「${res.source.name}」を作成しました。`);
      } else {
        setSourceMsg("sourceを作成しました。");
      }

      setSourceForm(newSourceForm());
      await loadSources();
    } catch (err: unknown) {
      setSourceMsg(toErrorMessage(err, "sourceの作成に失敗しました。"));
    } finally {
      setSavingSource(false);
    }
  }

  async function onSubmitItem(
    e: React.FormEvent<HTMLFormElement>,
  ): Promise<void> {
    e.preventDefault();

    setItemMsg("");
    setCreatedItem(null);

    if (sources.length === 0) {
      setItemMsg("先にsourceを1件以上登録してください。");
      return;
    }

    const sourceId = Number(itemForm.source_id);

    if (!sourceId) {
      setItemMsg("source を選択してください。");
      return;
    }

    if (!itemForm.published_at) {
      setItemMsg("公開日時を入力してください。");
      return;
    }

    setSavingItem(true);

    try {
      const body = {
        title: itemForm.title.trim(),
        summary: toNullableText(itemForm.summary),
        url: toNullableText(itemForm.url),
        image_url: toNullableText(itemForm.image_url),
        kind: itemForm.kind,
        source_id: sourceId,
        published_at: toApiDateTime(itemForm.published_at),
      };

      const res = await api<ItemRes>("/items", {
        method: "POST",
        auth: true,
        body,
      });

      if (res?.item) {
        setCreatedItem(res.item);
        setItemMsg(`item「${res.item.title}」を作成しました。`);
      } else {
        setItemMsg("itemを作成しました。");
      }

      setItemForm((prev) => ({
        ...newItemForm(),
        source_id: prev.source_id || firstSourceId,
        kind: "news",
      }));
    } catch (err: unknown) {
      setItemMsg(toErrorMessage(err, "itemの作成に失敗しました。"));
    } finally {
      setSavingItem(false);
    }
  }

  return (
    <div className="min-h-screen bg-[#f7f3ee] px-4 py-8 text-stone-800">
      <div className="mx-auto max-w-6xl">
        <header className="mb-8 overflow-hidden rounded-3xl border border-[#d9c7b8] bg-white shadow-sm">
          <div className="border-b border-[#e9ddd3] bg-gradient-to-r from-[#5a3825] via-[#7a5239] to-[#a67c52] px-6 py-8 text-white">
            <p className="text-sm font-semibold uppercase tracking-[0.2em]">
              Coffee SPA Admin
            </p>
            <h1 className="mt-2 text-3xl font-bold">管理画面</h1>
          </div>

          <div className="flex flex-wrap gap-3 px-6 py-4 text-sm">
            <Link
              to="/"
              className="rounded-full border border-[#ccb39e] bg-[#fffaf6] px-4 py-2 font-medium text-[#5a3825] transition hover:bg-[#f5ece5]"
            >
              公開トップへ
            </Link>

            <Link
              to="/me"
              className="rounded-full border border-[#ccb39e] bg-[#fffaf6] px-4 py-2 font-medium text-[#5a3825] transition hover:bg-[#f5ece5]"
            >
              me へ
            </Link>
          </div>
        </header>

        <div className="grid gap-6 lg:grid-cols-[1.1fr_1.4fr]">
          <section className="rounded-3xl border border-[#d9c7b8] bg-white p-6 shadow-sm">
            <SectionTitle title="Source 登録" sub="出典元を作成" />

            <form
              onSubmit={(e) => void onSubmitSource(e)}
              className="space-y-4"
            >
              <div>
                <FieldLabel htmlFor="source-name" label="Source名" required />
                <input
                  id="source-name"
                  name="name"
                  type="text"
                  value={sourceForm.name}
                  onChange={onChangeSourceForm}
                  placeholder="例: Coffee Daily"
                  className="w-full rounded-2xl border border-[#d8c8bc] bg-[#fffdfb] px-4 py-3 text-sm outline-none ring-0 transition focus:border-[#8b5e3c]"
                />
              </div>

              <div>
                <FieldLabel htmlFor="source-site-url" label="サイトURL" />
                <input
                  id="source-site-url"
                  name="site_url"
                  type="url"
                  value={sourceForm.site_url}
                  onChange={onChangeSourceForm}
                  placeholder="https://example.com"
                  className="w-full rounded-2xl border border-[#d8c8bc] bg-[#fffdfb] px-4 py-3 text-sm outline-none ring-0 transition focus:border-[#8b5e3c]"
                />
              </div>

              <button
                type="submit"
                disabled={savingSource}
                className="inline-flex rounded-2xl bg-[#5a3825] px-5 py-3 text-sm font-semibold text-white transition hover:opacity-90 disabled:cursor-not-allowed disabled:opacity-60"
              >
                {savingSource ? "作成中..." : "Sourceを作成"}
              </button>
            </form>

            {sourceMsg ? (
              <p className="mt-4 rounded-2xl bg-[#f8efe7] px-4 py-3 text-sm text-[#5a3825]">
                {sourceMsg}
              </p>
            ) : null}
          </section>

          <section className="rounded-3xl border border-[#d9c7b8] bg-white p-6 shadow-sm">
            <SectionTitle
              title="Item 登録"
              sub="ニュース、レシピ、セール、店舗情報を登録。"
            />

            <form onSubmit={(e) => void onSubmitItem(e)} className="space-y-4">
              <div className="grid gap-4 md:grid-cols-2">
                <div className="md:col-span-2">
                  <FieldLabel htmlFor="item-title" label="タイトル" required />
                  <input
                    id="item-title"
                    name="title"
                    type="text"
                    value={itemForm.title}
                    onChange={onChangeItemForm}
                    placeholder="例: 今週の新作コーヒー豆"
                    className="w-full rounded-2xl border border-[#d8c8bc] bg-[#fffdfb] px-4 py-3 text-sm outline-none transition focus:border-[#8b5e3c]"
                  />
                </div>

                <div className="md:col-span-2">
                  <FieldLabel htmlFor="item-summary" label="要約" />
                  <textarea
                    id="item-summary"
                    name="summary"
                    value={itemForm.summary}
                    onChange={onChangeItemForm}
                    placeholder="短い説明"
                    rows={4}
                    className="w-full rounded-2xl border border-[#d8c8bc] bg-[#fffdfb] px-4 py-3 text-sm outline-none transition focus:border-[#8b5e3c]"
                  />
                </div>

                <div>
                  <FieldLabel htmlFor="item-url" label="遷移先URL" />
                  <input
                    id="item-url"
                    name="url"
                    type="url"
                    value={itemForm.url}
                    onChange={onChangeItemForm}
                    placeholder="https://example.com/article"
                    className="w-full rounded-2xl border border-[#d8c8bc] bg-[#fffdfb] px-4 py-3 text-sm outline-none transition focus:border-[#8b5e3c]"
                  />
                </div>

                <div>
                  <FieldLabel htmlFor="item-image-url" label="画像URL" />
                  <input
                    id="item-image-url"
                    name="image_url"
                    type="url"
                    value={itemForm.image_url}
                    onChange={onChangeItemForm}
                    placeholder="https://example.com/image.jpg"
                    className="w-full rounded-2xl border border-[#d8c8bc] bg-[#fffdfb] px-4 py-3 text-sm outline-none transition focus:border-[#8b5e3c]"
                  />
                </div>

                <div>
                  <FieldLabel htmlFor="item-kind" label="種別" required />
                  <select
                    id="item-kind"
                    name="kind"
                    value={itemForm.kind}
                    onChange={onChangeItemForm}
                    className="w-full rounded-2xl border border-[#d8c8bc] bg-[#fffdfb] px-4 py-3 text-sm outline-none transition focus:border-[#8b5e3c]"
                  >
                    <option value="news">news</option>
                    <option value="recipe">recipe</option>
                    <option value="deal">deal</option>
                    <option value="shop">shop</option>
                  </select>
                </div>

                <div>
                  <FieldLabel
                    htmlFor="item-source-id"
                    label="Source"
                    required
                  />
                  <select
                    id="item-source-id"
                    name="source_id"
                    value={itemForm.source_id}
                    onChange={onChangeItemForm}
                    disabled={sources.length === 0}
                    className="w-full rounded-2xl border border-[#d8c8bc] bg-[#fffdfb] px-4 py-3 text-sm outline-none transition focus:border-[#8b5e3c] disabled:cursor-not-allowed disabled:bg-stone-100"
                  >
                    {sources.length === 0 ? (
                      <option value="">sourceを先に登録</option>
                    ) : (
                      sources.map((source) => (
                        <option key={source.id} value={String(source.id)}>
                          {source.name} / id:{source.id}
                        </option>
                      ))
                    )}
                  </select>
                </div>

                <div className="md:col-span-2">
                  <FieldLabel
                    htmlFor="item-published-at"
                    label="公開日時"
                    required
                  />
                  <input
                    id="item-published-at"
                    name="published_at"
                    type="datetime-local"
                    value={japanDateTime(itemForm.published_at)}
                    onChange={onChangeItemForm}
                    className="w-full rounded-2xl border border-[#d8c8bc] bg-[#fffdfb] px-4 py-3 text-sm outline-none transition focus:border-[#8b5e3c]"
                  />
                </div>
              </div>

              <button
                type="submit"
                disabled={savingItem || sources.length === 0}
                className="inline-flex rounded-2xl bg-[#8b5e3c] px-5 py-3 text-sm font-semibold text-white transition hover:opacity-90 disabled:cursor-not-allowed disabled:opacity-60"
              >
                {savingItem ? "作成中..." : "Itemを作成"}
              </button>
            </form>

            {itemMsg ? (
              <p className="mt-4 rounded-2xl bg-[#f8efe7] px-4 py-3 text-sm text-[#5a3825]">
                {itemMsg}
              </p>
            ) : null}

            {createdItem ? (
              <div className="mt-4 rounded-3xl border border-[#e3d6cc] bg-[#fffaf6] p-4">
                <p className="text-sm font-semibold text-[#5a3825]">
                  直近の作成結果
                </p>

                <dl className="mt-3 grid gap-2 text-sm text-stone-700">
                  <div className="grid grid-cols-[110px_1fr] gap-3">
                    <dt className="font-semibold">ID</dt>
                    <dd>{createdItem.id}</dd>
                  </div>

                  <div className="grid grid-cols-[110px_1fr] gap-3">
                    <dt className="font-semibold">タイトル</dt>
                    <dd>{createdItem.title}</dd>
                  </div>

                  <div className="grid grid-cols-[110px_1fr] gap-3">
                    <dt className="font-semibold">種別</dt>
                    <dd>{createdItem.kind}</dd>
                  </div>

                  <div className="grid grid-cols-[110px_1fr] gap-3">
                    <dt className="font-semibold">Source ID</dt>
                    <dd>{createdItem.source_id}</dd>
                  </div>

                  <div className="grid grid-cols-[110px_1fr] gap-3">
                    <dt className="font-semibold">公開日時</dt>
                    <dd>{createdItem.published_at}</dd>
                  </div>
                </dl>
              </div>
            ) : null}
          </section>
        </div>

        <section className="mt-6 rounded-3xl border border-[#d9c7b8] bg-white p-6 shadow-sm">
          <SectionTitle title="Source 一覧" sub="itemを登録します。" />

          {loadingSources ? (
            <p className="text-sm text-stone-600">loading...</p>
          ) : null}

          {!loadingSources && sourcesMsg ? (
            <p className="rounded-2xl bg-[#f8efe7] px-4 py-3 text-sm text-[#5a3825]">
              {sourcesMsg}
            </p>
          ) : null}

          {!loadingSources && !sourcesMsg && sources.length === 0 ? (
            <p className="text-sm text-stone-600">まだありません。</p>
          ) : null}

          {!loadingSources && sources.length > 0 ? (
            <div className="overflow-x-auto">
              <table className="min-w-full border-separate border-spacing-y-2 text-sm">
                <thead>
                  <tr className="text-left text-stone-600">
                    <th className="px-3 py-2">id</th>
                    <th className="px-3 py-2">name</th>
                    <th className="px-3 py-2">site_url</th>
                    <th className="px-3 py-2">created_at</th>
                  </tr>
                </thead>

                <tbody>
                  {sources.map((source) => (
                    <tr
                      key={source.id}
                      className="rounded-2xl bg-[#fffaf6] text-stone-800"
                    >
                      <td className="px-3 py-3">{source.id}</td>
                      <td className="px-3 py-3 font-medium">{source.name}</td>
                      <td className="px-3 py-3">
                        {source.site_url ? (
                          <a
                            href={source.site_url}
                            target="_blank"
                            rel="noreferrer"
                            className="break-all text-[#7a5239] underline"
                          >
                            {source.site_url}
                          </a>
                        ) : (
                          <span className="text-stone-400">-</span>
                        )}
                      </td>
                      <td className="px-3 py-3">
                        {japanDateTime(source.created_at)}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          ) : null}
        </section>
      </div>
    </div>
  );
}
