import { create } from 'zustand'
import { persist } from 'zustand/middleware'

export type TerminalTheme = 'term-dark' | 'term-light'

interface SettingsState {
  theme: TerminalTheme
  soundEffects: boolean
  setTheme: (theme: TerminalTheme) => void
  toggleSoundEffects: () => void
}

export const useSettingsStore = create<SettingsState>()(
  persist(
    (set) => ({
      theme: 'term-dark',
      soundEffects: false,
      setTheme: (theme) => set({ theme: theme === 'term-light' ? 'term-light' : 'term-dark' }),
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

