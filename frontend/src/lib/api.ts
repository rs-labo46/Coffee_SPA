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

//access tokenを保存するキー
const tokenKey = "access_token";

//apiの接続先
function getBaseUrl(): string {
  return import.meta.env.VITE_API_BASE_URL || "http://localhost:8080";
}

function joinUrl(base: string, path: string): string {
  const left = base.endsWith("/") ? base.slice(0, -1) : base;
  const right = path.startsWith("/") ? path : `/${path}`;
  return `${left}${right}`;
}

//localStorageからaccess tokenを取得。

export function getToken(): string {
  return localStorage.getItem(tokenKey) || "";
}

//localStorageにaccess tokenを保存。
export function setToken(token: string): void {
  localStorage.setItem(tokenKey, token);
}

//logoutやrefresh失敗時にlocalStorageからaccess tokenを削除
export function clearToken(): void {
  localStorage.removeItem(tokenKey);
}

//Cookie から指定された csrftokenの名前の値を取得。
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

//fetch に渡すヘッダーを組み立て,Content-Type は body がある時だけ付ける。
function buildHeaders(opt?: ApiOption): Headers {
  const headers = new Headers();

  if (opt?.auth) {
    const token = getToken();

    if (token) {
      headers.set("Authorization", `Bearer ${token}`);
    }
  }

  if (opt?.csrf) {
    const csrf = getCookie("csrf_token");

    if (csrf) {
      headers.set("X-CSRF-Token", csrf);
    }
  }

  if (opt?.body !== undefined) {
    headers.set("Content-Type", "application/json");
  }

  return headers;
}

//レスポンスがJSONの時だけ安全にJSONとして読む。
async function readJsonSafe<T>(res: Response): Promise<T | null> {
  const contentType = res.headers.get("Content-Type") || "";

  if (!contentType.includes("application/json")) {
    return null;
  }

  return (await res.json()) as T;
}

//サーバのerror / message が取れれば使って、無ければ既定値を使う。
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
