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

  const data = await response.json();

  return { ok: response.ok, status: response.status, data };
}