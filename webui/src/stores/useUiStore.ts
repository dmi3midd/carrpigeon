import { create } from 'zustand'

export type NavTab = 'dashboard' | 'send' | 'receivers' | 'groups' | 'templates' | 'logs' | 'settings'

export interface ToastMessage {
  id: string
  type: 'info' | 'success' | 'warn' | 'error'
  text: string
  timestamp: string
  durationMs?: number
}

export interface ModalState {
  type:
    | 'createReceiver'
    | 'editReceiver'
    | 'deleteReceiver'
    | 'createGroup'
    | 'editGroup'
    | 'deleteGroup'
    | 'uploadTemplate'
    | 'viewTemplate'
    | 'deleteTemplate'
    | null
  data?: unknown
}

interface UiState {
  activeTab: NavTab
  toasts: ToastMessage[]
  modal: ModalState
  isSidebarOpen: boolean
  setActiveTab: (tab: NavTab) => void
  addToast: (type: ToastMessage['type'], text: string, durationMs?: number) => void
  removeToast: (id: string) => void
  openModal: (type: ModalState['type'], data?: unknown) => void
  closeModal: () => void
  toggleSidebar: () => void
}

export const useUiStore = create<UiState>((set) => ({
  activeTab: 'dashboard',
  toasts: [],
  modal: { type: null },
  isSidebarOpen: true,

  setActiveTab: (tab) => set({ activeTab: tab }),

  addToast: (type, text, durationMs = 3000) =>
    set((state) => {
      const id = `${Date.now()}-${Math.random().toString(36).slice(2, 7)}`
      const newToast: ToastMessage = {
        id,
        type,
        text,
        timestamp: new Date().toLocaleTimeString(),
        durationMs,
      }
      return { toasts: [...state.toasts.slice(-5), newToast] }
    }),

  removeToast: (id) =>
    set((state) => ({
      toasts: state.toasts.filter((t) => t.id !== id),
    })),

  openModal: (type, data) => set({ modal: { type, data } }),

  closeModal: () => set({ modal: { type: null } }),

  toggleSidebar: () => set((s) => ({ isSidebarOpen: !s.isSidebarOpen })),
}))
