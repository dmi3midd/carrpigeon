import { create } from 'zustand'
import { healthService } from '../services/health'
import type { HealthResponse } from '../types/health'

interface HealthState {
  data: HealthResponse | null
  isLoading: boolean
  error: string | null
  lastChecked: string | null
  fetchHealth: () => Promise<void>
}

export const useHealthStore = create<HealthState>((set) => ({
  data: null,
  isLoading: false,
  error: null,
  lastChecked: null,

  fetchHealth: async () => {
    set({ isLoading: true, error: null })
    try {
      const data = await healthService.getHealth()
      set({
        data,
        isLoading: false,
        lastChecked: new Date().toLocaleTimeString(),
        error: null,
      })
    } catch (err) {
      set({
        isLoading: false,
        error: err instanceof Error ? err.message : 'Failed to reach service',
        lastChecked: new Date().toLocaleTimeString(),
      })
    }
  },
}))
