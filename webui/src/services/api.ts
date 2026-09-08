import { APIError } from '../types/api'

const BASE_URL = ''

export interface RequestOptions extends RequestInit {
  params?: Record<string, string | number | boolean | undefined>
}

export async function request<T>(endpoint: string, options: RequestOptions = {}): Promise<T> {
  const { params, headers, ...customConfig } = options

  let url = `${BASE_URL}${endpoint}`
  if (params) {
    const searchParams = new URLSearchParams()
    for (const [key, value] of Object.entries(params)) {
      if (value !== undefined) {
        searchParams.append(key, String(value))
      }
    }
    const queryString = searchParams.toString()
    if (queryString) {
      url += (url.includes('?') ? '&' : '?') + queryString
    }
  }

  const isFormData = customConfig.body instanceof FormData

  const config: RequestInit = {
    ...customConfig,
    headers: {
      ...(!isFormData ? { 'Content-Type': 'application/json' } : {}),
      Accept: 'application/json',
      ...headers,
    },
  }

  const response = await fetch(url, config)

  if (!response.ok) {
    let errorMsg = `HTTP Error ${response.status}: ${response.statusText}`
    let details: Record<string, unknown> | undefined

    try {
      const errorJson = await response.json()
      if (errorJson.message) {
        errorMsg = errorJson.message
      } else if (errorJson.error) {
        errorMsg = errorJson.error
      }
      details = errorJson
    } catch {
      // Body is not JSON, try text
      const text = await response.text().catch(() => '')
      if (text) errorMsg = text
    }

    throw new APIError(errorMsg, response.status, details)
  }

  // Handle 202 Accepted or 204 No Content with empty body
  if (response.status === 204 || response.status === 202) {
    const text = await response.text()
    if (!text) {
      return {} as T
    }
    try {
      return JSON.parse(text) as T
    } catch {
      return {} as T
    }
  }

  const contentType = response.headers.get('content-type')
  if (contentType && contentType.includes('application/json')) {
    return (await response.json()) as T
  }

  return (await response.text()) as unknown as T
}
