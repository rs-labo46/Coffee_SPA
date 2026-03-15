type Json = string | number | boolean | null | Json[] | { [key: string]: Json };

type Method = "GET" | "POST" | "PATCH" | "PUT" | "DELETE";

type ApiOption = {
  method?: Method;
  body?: Json;
  auth?: boolean;
  csrf?: boolean;
};

type ApiErrBody = {
  error?: string;
  message?: string;
};

//フロント側で扱うAPIエラー
export class ApiError extends Error {
  status: number;
  code: string;

  constructor(status: number, code: string, message: string) {
    super(message);
    this.name = "ApiError";
    this.status = status;
    this.code = code;
  }
}

//localStorageにaccess tokenを保存するキー
const tokenKey = "access_token";

//APIのベースURLを返す
function getBaseUrl(): string {
  return import.meta.env.VITE_API_BASE_URL || "http://localhost:8080";
}

//ベースURLとパスをつなぐ
function joinUrl(base: string, path: string): string {
  const left = base.endsWith("/") ? base.slice(0, -1) : base;
  const right = path.startsWith("/") ? path : `/${path}`;
  return `${left}${right}`;
}

//localStorageからaccess tokenを読む
export function getToken(): string {
  return localStorage.getItem(tokenKey) || "";
}

//localStorageにaccess tokenを保存
export function setToken(token: string): void {
  localStorage.setItem(tokenKey, token);
}

//localStorageからaccess tokenを削除
export function clearToken(): void {
  localStorage.removeItem(tokenKey);
}

//Cookieから指定名の値を取り出す
export function getCookie(name: string): string {
  const items = document.cookie.split(";");

  for (const item of items) {
    const value = item.trim();

    if (value.startsWith(`${name}=`)) {
      return decodeURIComponent(value.slice(name.length + 1));
    }
  }

  return "";
}

//fetchに渡すヘッダー
function buildHeaders(opt?: ApiOption): Headers {
  const headers = new Headers();

  //認証が必要ならAuthorizationを付ける
  if (opt?.auth) {
    const token = getToken();

    if (token) {
      headers.set("Authorization", `Bearer ${token}`);
    }
  }

  //CSRFが必要ならcookie から読み出してheaderに付ける
  if (opt?.csrf) {
    const csrf = getCookie("csrf_token");

    if (csrf) {
      headers.set("X-CSRF-Token", csrf);
    }
  }

  //bodyがある時だけJSONを宣言
  if (opt?.body !== undefined) {
    headers.set("Content-Type", "application/json");
  }

  return headers;
}

//レスポンスがJSONの時だけ安全に読む
async function readJsonSafe<T>(res: Response): Promise<T | null> {
  const contentType = res.headers.get("Content-Type") || "";

  if (!contentType.includes("application/json")) {
    return null;
  }

  return (await res.json()) as T;
}

//エラーレスポンスをApiErrorに変換
async function throwApiError(res: Response): Promise<never> {
  const body = await readJsonSafe<ApiErrBody>(res);

  const code = body?.error || "request_failed";
  const message = body?.message || body?.error || `HTTP ${res.status}`;

  throw new ApiError(res.status, code, message);
}

export async function api<T>(
  path: string,
  opt?: ApiOption,
): Promise<T | undefined> {
  const init: RequestInit = {
    method: opt?.method || "GET",
    headers: buildHeaders(opt),
    credentials: "include",
  };

  //bodyがある時だけJSON文字列に
  if (opt?.body !== undefined) {
    init.body = JSON.stringify(opt.body);
  }

  const res = await fetch(joinUrl(getBaseUrl(), path), init);

  if (!res.ok) {
    await throwApiError(res);
  }

  if (res.status === 204) {
    return undefined;
  }

  const data = await readJsonSafe<T>(res);

  if (data === null) {
    throw new ApiError(500, "invalid_response", "response is not json");
  }

  return data;
}
