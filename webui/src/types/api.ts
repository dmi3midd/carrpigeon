export interface APIErrorResponse {
  error?: string
  message?: string
  code?: number
  details?: Record<string, unknown>
}

export class APIError extends Error {
  status: number
  details?: Record<string, unknown>

  constructor(message: string, status: number, details?: Record<string, unknown>) {
    super(message)
    this.name = 'APIError'
    this.status = status
    this.details = details
  }
}

export interface PaginationParams {
  limit?: number
  offset?: number
}
