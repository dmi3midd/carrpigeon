import { create } from 'zustand'
import { groupsService } from '../services/groups'
import type {
  CreateGroupRequest,
  Group,
  UpdateGroupRequest,
} from '../types/group'
import type { Receiver } from '../types/receiver'

interface GroupsState {
  groups: Group[]
  isLoading: boolean
  error: string | null
  limit: number
  offset: number
  hasMore: boolean
  selectedGroup: Group | null
  groupReceivers: Receiver[]
  isLoadingMembers: boolean
  searchQuery: string
  fetchGroups: (offset?: number) => Promise<void>
  createGroup: (data: CreateGroupRequest) => Promise<string>
  updateGroup: (id: string, data: UpdateGroupRequest) => Promise<void>
  deleteGroup: (id: string) => Promise<void>
  setSelectedGroup: (group: Group | null) => void
  fetchGroupReceivers: (groupId: string) => Promise<void>
  addReceiverToGroup: (groupId: string, receiverId: string) => Promise<void>
  removeReceiverFromGroup: (groupId: string, receiverId: string) => Promise<void>
  setSearchQuery: (query: string) => void
  nextPage: () => void
  prevPage: () => void
}

export const useGroupsStore = create<GroupsState>((set, get) => ({
  groups: [],
  isLoading: false,
  error: null,
  limit: 10,
  offset: 0,
  hasMore: false,
  selectedGroup: null,
  groupReceivers: [],
  isLoadingMembers: false,
  searchQuery: '',

  fetchGroups: async (customOffset?: number) => {
    const { limit, offset } = get()
    const targetOffset = customOffset !== undefined ? customOffset : offset
    set({ isLoading: true, error: null })
    try {
      const items = await groupsService.list(limit + 1, targetOffset)
      const hasMore = items.length > limit
      const groups = hasMore ? items.slice(0, limit) : items
      set({
        groups,
        isLoading: false,
        offset: targetOffset,
        hasMore,
        error: null,
      })
    } catch (err) {
      set({
        isLoading: false,
        error: err instanceof Error ? err.message : 'Failed to fetch groups',
      })
    }
  },

  createGroup: async (data) => {
    set({ isLoading: true, error: null })
    try {
      const res = await groupsService.create(data)
      await get().fetchGroups(0)
      return res.id
    } catch (err) {
      set({ isLoading: false })
      throw err
    }
  },

  updateGroup: async (id, data) => {
    set({ isLoading: true, error: null })
    try {
      await groupsService.update(id, data)
      await get().fetchGroups()
    } catch (err) {
      set({ isLoading: false })
      throw err
    }
  },

  deleteGroup: async (id) => {
    set({ isLoading: true, error: null })
    try {
      await groupsService.delete(id)
      if (get().selectedGroup?.id === id) {
        set({ selectedGroup: null, groupReceivers: [] })
      }
      await get().fetchGroups()
    } catch (err) {
      set({ isLoading: false })
      throw err
    }
  },

  setSelectedGroup: (selectedGroup) => {
    set({ selectedGroup })
    if (selectedGroup) {
      get().fetchGroupReceivers(selectedGroup.id)
    } else {
      set({ groupReceivers: [] })
    }
  },

  fetchGroupReceivers: async (groupId: string) => {
    set({ isLoadingMembers: true })
    try {
      const receivers = await groupsService.listReceivers(groupId, 100, 0)
      set({ groupReceivers: receivers, isLoadingMembers: false })
    } catch (err) {
      set({
        isLoadingMembers: false,
        error: err instanceof Error ? err.message : 'Failed to fetch group receivers',
      })
    }
  },

  addReceiverToGroup: async (groupId: string, receiverId: string) => {
    try {
      await groupsService.addReceiver(groupId, receiverId)
      await get().fetchGroupReceivers(groupId)
      await get().fetchGroups()
    } catch (err) {
      throw err
    }
  },

  removeReceiverFromGroup: async (groupId: string, receiverId: string) => {
    try {
      await groupsService.removeReceiver(groupId, receiverId)
      await get().fetchGroupReceivers(groupId)
      await get().fetchGroups()
    } catch (err) {
      throw err
    }
  },

  setSearchQuery: (searchQuery) => set({ searchQuery }),

  nextPage: () => {
    const { offset, limit, hasMore } = get()
    if (hasMore) {
      get().fetchGroups(offset + limit)
    }
  },

  prevPage: () => {
    const { offset, limit } = get()
    const newOffset = Math.max(0, offset - limit)
    get().fetchGroups(newOffset)
  },
}))
