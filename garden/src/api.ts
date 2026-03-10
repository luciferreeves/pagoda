import { API_URL } from "./config";

interface APIOptions {
  method?: string;
  body?: unknown;
  token?: string | null;
}

interface APIResponse<T> {
  ok: boolean;
  status: number;
  data: T;
}

export async function api<T>(path: string, options: APIOptions = {}): Promise<APIResponse<T>> {
  const headers: Record<string, string> = {};

  if (options.body) {
    headers["Content-Type"] = "application/json";
  }

  if (options.token) {
    headers["Authorization"] = `Bearer ${options.token}`;
  }

  const response = await fetch(`${API_URL}${path}`, {
    method: options.method || "GET",
    headers,
    body: options.body ? JSON.stringify(options.body) : undefined,
  });

  const text = await response.text();
  const data = text ? JSON.parse(text) : null;

  return { ok: response.ok, status: response.status, data };
}

export async function uploadFile<T>(path: string, file: File, token: string): Promise<APIResponse<T>> {
  const form = new FormData();
  form.append("file", file);

  const response = await fetch(`${API_URL}${path}`, {
    method: "POST",
    headers: { "Authorization": `Bearer ${token}` },
    body: form,
  });

  const text = await response.text();
  const data = text ? JSON.parse(text) : null;

  return { ok: response.ok, status: response.status, data };
}