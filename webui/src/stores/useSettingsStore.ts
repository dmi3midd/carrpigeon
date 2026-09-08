import { create } from 'zustand'
import { persist } from 'zustand/middleware'

export type TerminalTheme = 'term-dark' | 'term-light'

interface SettingsState {
  theme: TerminalTheme
  autoRefresh: boolean
  refreshIntervalSeconds: number
  soundEffects: boolean
  setTheme: (theme: TerminalTheme) => void
  setAutoRefresh: (enabled: boolean) => void
  setRefreshIntervalSeconds: (seconds: number) => void
  toggleSoundEffects: () => void
}

export const useSettingsStore = create<SettingsState>()(
  persist(
    (set) => ({
      theme: 'term-dark',
      autoRefresh: true,
      refreshIntervalSeconds: 10,
      soundEffects: false,
      setTheme: (theme) => set({ theme: theme === 'term-light' ? 'term-light' : 'term-dark' }),
      setAutoRefresh: (autoRefresh) => set({ autoRefresh }),
      setRefreshIntervalSeconds: (refreshIntervalSeconds) => set({ refreshIntervalSeconds }),
      toggleSoundEffects: () => set((s) => ({ soundEffects: !s.soundEffects })),
    }),
    {
      name: 'carrpigeon-cli-settings',
      merge: (persistedState, currentState) => {
        const persisted = persistedState as Partial<SettingsState> | undefined
        return {
          ...currentState,
          ...persisted,
          theme: persisted?.theme === 'term-light' ? 'term-light' : 'term-dark',
        }
      },
    }
  )
)

