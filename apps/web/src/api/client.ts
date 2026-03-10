export type ApiError = {
  error: string
}

const apiBaseUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

async function request<T>(path: string, init?: RequestInit, token?: string): Promise<T> {
  const headers = new Headers(init?.headers)
  if (!headers.has('Content-Type') && init?.body) {
    headers.set('Content-Type', 'application/json')
  }
  if (token) {
    headers.set('Authorization', `Bearer ${token}`)
  }

  const response = await fetch(`${apiBaseUrl}${path}`, {
    ...init,
    headers,
  })

  if (!response.ok) {
    const payload = (await response.json().catch(() => ({ error: 'Request failed' }))) as ApiError
    throw new Error(payload.error || 'Request failed')
  }

  return (await response.json()) as T
}

export const apiClient = {
  get: <T>(path: string, token?: string) => request<T>(path, undefined, token),
  post: <T>(path: string, body: unknown, token?: string) =>
    request<T>(path, { method: 'POST', body: JSON.stringify(body) }, token),
  patch: <T>(path: string, body: unknown, token?: string) =>
    request<T>(path, { method: 'PATCH', body: JSON.stringify(body) }, token),
}
