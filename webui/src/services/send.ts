import { request } from './api'
import type {
  SendGroupRequest,
  SendGroupWithTemplateRequest,
  SendSingleRequest,
  SendSingleWithTemplateRequest,
} from '../types/send'

export const sendService = {
  sendSingle: async (data: SendSingleRequest): Promise<void> => {
    await request<void>('/send/single', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  },

  sendSingleWithTemplate: async (data: SendSingleWithTemplateRequest): Promise<void> => {
    await request<void>('/send/single/template', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  },

  sendGroup: async (data: SendGroupRequest): Promise<void> => {
    await request<void>('/send/group', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  },

  sendGroupWithTemplate: async (data: SendGroupWithTemplateRequest): Promise<void> => {
    await request<void>('/send/group/template', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  },
}
