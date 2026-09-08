import { request } from './api'
import type {
  CreateGroupRequest,
  CreateGroupResponse,
  GetGroupResponse,
  Group,
  ListGroupReceiversResponse,
  ListGroupsResponse,
  UpdateGroupRequest,
  UpdateGroupResponse,
} from '../types/group'
import type { Receiver } from '../types/receiver'

export const groupsService = {
  list: async (limit = 10, offset = 0): Promise<Group[]> => {
    const res = await request<ListGroupsResponse>('/groups', {
      params: { limit, offset },
    })
    return res.groups || []
  },

  getById: async (id: string): Promise<Group> => {
    const res = await request<GetGroupResponse>(`/groups/${encodeURIComponent(id)}`)
    return res.group
  },

  create: async (data: CreateGroupRequest): Promise<CreateGroupResponse> => {
    return request<CreateGroupResponse>('/groups', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  },

  update: async (id: string, data: UpdateGroupRequest): Promise<UpdateGroupResponse> => {
    return request<UpdateGroupResponse>(`/groups/${encodeURIComponent(id)}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    })
  },

  delete: async (id: string): Promise<void> => {
    await request<void>(`/groups/${encodeURIComponent(id)}`, {
      method: 'DELETE',
    })
  },

  listReceivers: async (groupId: string, limit = 50, offset = 0): Promise<Receiver[]> => {
    const res = await request<ListGroupReceiversResponse>(
      `/groups/${encodeURIComponent(groupId)}/receivers`,
      { params: { limit, offset } }
    )
    return res.receivers || []
  },

  addReceiver: async (groupId: string, receiverId: string): Promise<void> => {
    await request<void>(
      `/groups/${encodeURIComponent(groupId)}/receivers/${encodeURIComponent(receiverId)}`,
      { method: 'PUT' }
    )
  },

  removeReceiver: async (groupId: string, receiverId: string): Promise<void> => {
    await request<void>(
      `/groups/${encodeURIComponent(groupId)}/receivers/${encodeURIComponent(receiverId)}`,
      { method: 'DELETE' }
    )
  },
}
