import { create } from 'zustand'
import { templatesService } from '../services/templates'
import type { Template, TemplateMetadata } from '../types/template'

interface TemplatesState {
  templates: TemplateMetadata[]
  isLoading: boolean
  error: string | null
  limit: number
  offset: number
  hasMore: boolean
  selectedTemplate: Template | null
  isLoadingRaw: boolean
  searchQuery: string
  fetchTemplates: (offset?: number) => Promise<void>
  fetchRawTemplate: (id: string) => Promise<Template>
  createTemplate: (name: string, file: File) => Promise<string>
  updateTemplate: (id: string, name: string, file: File) => Promise<void>
  deleteTemplate: (id: string) => Promise<void>
  setSelectedTemplate: (template: Template | null) => void
  setSearchQuery: (query: string) => void
  nextPage: () => void
  prevPage: () => void
}

export const useTemplatesStore = create<TemplatesState>((set, get) => ({
  templates: [],
  isLoading: false,
  error: null,
  limit: 10,
  offset: 0,
  hasMore: false,
  selectedTemplate: null,
  isLoadingRaw: false,
  searchQuery: '',

  fetchTemplates: async (customOffset?: number) => {
    const { limit, offset } = get()
    const targetOffset = customOffset !== undefined ? customOffset : offset
    set({ isLoading: true, error: null })
    try {
      const items = await templatesService.list(limit + 1, targetOffset)
      const hasMore = items.length > limit
      const templates = hasMore ? items.slice(0, limit) : items
      set({
        templates,
        isLoading: false,
        offset: targetOffset,
        hasMore,
        error: null,
      })
    } catch (err) {
      set({
        isLoading: false,
        error: err instanceof Error ? err.message : 'Failed to fetch templates',
      })
    }
  },

  fetchRawTemplate: async (id: string) => {
    set({ isLoadingRaw: true })
    try {
      const tmpl = await templatesService.getRaw(id)
      set({ selectedTemplate: tmpl, isLoadingRaw: false })
      return tmpl
    } catch (err) {
      set({
        isLoadingRaw: false,
        error: err instanceof Error ? err.message : 'Failed to fetch raw template',
      })
      throw err
    }
  },

  createTemplate: async (name, file) => {
    set({ isLoading: true, error: null })
    try {
      const res = await templatesService.create(name, file)
      await get().fetchTemplates(0)
      return res.id
    } catch (err) {
      set({ isLoading: false })
      throw err
    }
  },

  updateTemplate: async (id, name, file) => {
    set({ isLoading: true, error: null })
    try {
      await templatesService.update(id, name, file)
      await get().fetchTemplates()
      if (get().selectedTemplate?.id === id) {
        await get().fetchRawTemplate(id)
      }
    } catch (err) {
      set({ isLoading: false })
      throw err
    }
  },

  deleteTemplate: async (id) => {
    set({ isLoading: true, error: null })
    try {
      await templatesService.delete(id)
      if (get().selectedTemplate?.id === id) {
        set({ selectedTemplate: null })
      }
      await get().fetchTemplates()
    } catch (err) {
      set({ isLoading: false })
      throw err
    }
  },

  setSelectedTemplate: (selectedTemplate) => set({ selectedTemplate }),
  setSearchQuery: (searchQuery) => set({ searchQuery }),

  nextPage: () => {
    const { offset, limit, hasMore } = get()
    if (hasMore) {
      get().fetchTemplates(offset + limit)
    }
  },

  prevPage: () => {
    const { offset, limit } = get()
    const newOffset = Math.max(0, offset - limit)
    get().fetchTemplates(newOffset)
  },
}))
