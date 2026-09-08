import type { Receiver } from './receiver'

export interface Group {
  id: string
  name: string
  description: string
  receivers_count: number
  created_at: string
  updated_at: string
}

export interface GetGroupResponse {
  group: Group
}

export interface ListGroupsResponse {
  groups: Group[] | null
}

export interface CreateGroupRequest {
  name: string
  description: string
}

export interface CreateGroupResponse {
  id: string
}

export interface UpdateGroupRequest {
  name?: string
  description?: string
}

export interface UpdateGroupResponse {
  id: string
}

export interface ListGroupReceiversResponse {
  receivers: Receiver[] | null
}
