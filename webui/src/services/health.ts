import { request } from './api'
import type { HealthResponse } from '../types/health'

export const healthService = {
  getHealth: async (): Promise<HealthResponse> => {
    return request<HealthResponse>('/health')
  },
}
