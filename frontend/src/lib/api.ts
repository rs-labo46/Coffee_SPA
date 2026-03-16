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

const tokenKey = "access_token";

function getBaseUrl(): string {
  return import.meta.env.VITE_API_BASE_URL || "http://localhost:8080";
}

function joinUrl(base: string, path: string): string {
  const left = base.endsWith("/") ? base.slice(0, -1) : base;
  const right = path.startsWith("/") ? path : `/${path}`;
  return `${left}${right}`;
}

export function getToken(): string {
  return localStorage.getItem(tokenKey) || "";
}

export function setToken(token: string): void {
  localStorage.setItem(tokenKey, token);
}

export function clearToken(): void {
  localStorage.removeItem(tokenKey);
}

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

async function readJsonSafe<T>(res: Response): Promise<T | null> {
  const contentType = res.headers.get("Content-Type") || "";

  if (!contentType.includes("application/json")) {
    return null;
  }

  return (await res.json()) as T;
}

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
