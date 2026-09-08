import { request } from './api'
import type {
  CreateReceiverRequest,
  CreateReceiverResponse,
  GetReceiverByIdResponse,
  ListReceiversResponse,
  Receiver,
  UpdateReceiverRequest,
  UpdateReceiverResponse,
} from '../types/receiver'

export const receiversService = {
  list: async (limit = 10, offset = 0): Promise<Receiver[]> => {
    const res = await request<ListReceiversResponse>('/receivers', {
      params: { limit, offset },
    })
    return res.receivers || []
  },

  getById: async (id: string): Promise<Receiver> => {
    const res = await request<GetReceiverByIdResponse>(`/receivers/${encodeURIComponent(id)}`)
    return res.receiver
  },

  create: async (data: CreateReceiverRequest): Promise<CreateReceiverResponse> => {
    return request<CreateReceiverResponse>('/receivers', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  },

  update: async (id: string, data: UpdateReceiverRequest): Promise<UpdateReceiverResponse> => {
    return request<UpdateReceiverResponse>(`/receivers/${encodeURIComponent(id)}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    })
  },

  delete: async (id: string): Promise<void> => {
    await request<void>(`/receivers/${encodeURIComponent(id)}`, {
      method: 'DELETE',
    })
  },
}
