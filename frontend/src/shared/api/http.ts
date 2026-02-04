export type HttpMethod = "GET" | "POST" | "PUT" | "DELETE";

export interface HttpOptions<TBody = unknown> {
  method?: HttpMethod;
  body?: TBody;
  token?: string;
}

const baseUrl = import.meta.env.VITE_API_BASE_URL ?? "http://localhost:54321";

export async function request<TResponse, TBody = unknown>(
  path: string,
  options: HttpOptions<TBody> = {},
): Promise<TResponse> {
  const headers: HeadersInit = {
    "Content-Type": "application/json",
  };

  if (options.token) {
    headers.Authorization = `Bearer ${options.token}`;
  }

  const response = await fetch(`${baseUrl}${path}`, {
    method: options.method ?? "GET",
    headers,
    body: options.body ? JSON.stringify(options.body) : undefined,
  });

  if (!response.ok) {
    const message = await response.text();
    throw new Error(message || `Request failed (${response.status})`);
  }

  const text = await response.text();
  if (!text) {
    return undefined as TResponse;
  }

  return JSON.parse(text) as TResponse;
}
