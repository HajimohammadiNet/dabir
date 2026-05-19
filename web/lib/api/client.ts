export type ApiError = {
  code: string;
  message: string;
};

export type ApiResponse<T> = {
  success: boolean;
  data?: T;
  error?: ApiError;
};

const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080/api/v1";

export class ApiClientError extends Error {
  code: string;
  status: number;

  constructor(message: string, code: string, status: number) {
    super(message);
    this.name = "ApiClientError";
    this.code = code;
    this.status = status;
  }
}

type RequestOptions = {
  token?: string | null;
  body?: unknown;
  headers?: Record<string, string>;
};

async function request<T>(
  path: string,
  method: string,
  options: RequestOptions = {}
): Promise<T> {
  const headers: Record<string, string> = {
    Accept: "application/json",
    ...options.headers,
  };

  let body: BodyInit | undefined;

  if (options.body instanceof FormData) {
    body = options.body;
  } else if (options.body !== undefined) {
    headers["Content-Type"] = "application/json";
    body = JSON.stringify(options.body);
  }

  if (options.token) {
    headers.Authorization = `Bearer ${options.token}`;
  }

  const res = await fetch(`${API_BASE_URL}${path}`, {
    method,
    headers,
    body,
    cache: "no-store",
  });

  const payload = (await res.json()) as ApiResponse<T>;

  if (!res.ok || !payload.success) {
    throw new ApiClientError(
      payload.error?.message || "Request failed",
      payload.error?.code || "REQUEST_FAILED",
      res.status
    );
  }

  return payload.data as T;
}

export const apiClient = {
  get: <T>(path: string, token?: string | null) =>
    request<T>(path, "GET", { token }),

  post: <T>(path: string, body?: unknown, token?: string | null) =>
    request<T>(path, "POST", { body, token }),

  patch: <T>(path: string, body?: unknown, token?: string | null) =>
    request<T>(path, "PATCH", { body, token }),

  delete: <T>(path: string, token?: string | null) =>
    request<T>(path, "DELETE", { token }),
};