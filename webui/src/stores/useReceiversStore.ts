import { create } from 'zustand'
import { receiversService } from '../services/receivers'
import type {
  CreateReceiverRequest,
  Receiver,
  UpdateReceiverRequest,
} from '../types/receiver'

interface ReceiversState {
  receivers: Receiver[]
  isLoading: boolean
  error: string | null
  limit: number
  offset: number
  hasMore: boolean
  selectedReceiver: Receiver | null
  searchQuery: string
  fetchReceivers: (offset?: number) => Promise<void>
  createReceiver: (data: CreateReceiverRequest) => Promise<string>
  updateReceiver: (id: string, data: UpdateReceiverRequest) => Promise<void>
  deleteReceiver: (id: string) => Promise<void>
  setSelectedReceiver: (receiver: Receiver | null) => void
  setSearchQuery: (query: string) => void
  nextPage: () => void
  prevPage: () => void
}

export const useReceiversStore = create<ReceiversState>((set, get) => ({
  receivers: [],
  isLoading: false,
  error: null,
  limit: 10,
  offset: 0,
  hasMore: false,
  selectedReceiver: null,
  searchQuery: '',

  fetchReceivers: async (customOffset?: number) => {
    const { limit, offset } = get()
    const targetOffset = customOffset !== undefined ? customOffset : offset
    set({ isLoading: true, error: null })
    try {
      const items = await receiversService.list(limit + 1, targetOffset)
      const hasMore = items.length > limit
      const receivers = hasMore ? items.slice(0, limit) : items
      set({
        receivers,
        isLoading: false,
        offset: targetOffset,
        hasMore,
        error: null,
      })
    } catch (err) {
      set({
        isLoading: false,
        error: err instanceof Error ? err.message : 'Failed to fetch receivers',
      })
    }
  },

  createReceiver: async (data) => {
    set({ isLoading: true, error: null })
    try {
      const res = await receiversService.create(data)
      await get().fetchReceivers(0)
      return res.id
    } catch (err) {
      set({ isLoading: false })
      throw err
    }
  },

  updateReceiver: async (id, data) => {
    set({ isLoading: true, error: null })
    try {
      await receiversService.update(id, data)
      await get().fetchReceivers()
    } catch (err) {
      set({ isLoading: false })
      throw err
    }
  },

  deleteReceiver: async (id) => {
    set({ isLoading: true, error: null })
    try {
      await receiversService.delete(id)
      await get().fetchReceivers()
    } catch (err) {
      set({ isLoading: false })
      throw err
    }
  },

  setSelectedReceiver: (selectedReceiver) => set({ selectedReceiver }),
  setSearchQuery: (searchQuery) => set({ searchQuery }),

  nextPage: () => {
    const { offset, limit, hasMore } = get()
    if (hasMore) {
      get().fetchReceivers(offset + limit)
    }
  },

  prevPage: () => {
    const { offset, limit } = get()
    const newOffset = Math.max(0, offset - limit)
    get().fetchReceivers(newOffset)
  },
}))
