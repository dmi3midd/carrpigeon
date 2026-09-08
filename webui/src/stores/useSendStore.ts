import { create } from 'zustand'
import { sendService } from '../services/send'
import type {
  DispatchLogItem,
  SendGroupRequest,
  SendGroupWithTemplateRequest,
  SendSingleRequest,
  SendSingleWithTemplateRequest,
} from '../types/send'

interface SendState {
  isSending: boolean
  logs: DispatchLogItem[]
  lastError: string | null
  sendSingle: (data: SendSingleRequest) => Promise<void>
  sendSingleWithTemplate: (data: SendSingleWithTemplateRequest) => Promise<void>
  sendGroup: (data: SendGroupRequest) => Promise<void>
  sendGroupWithTemplate: (data: SendGroupWithTemplateRequest) => Promise<void>
  clearLogs: () => void
}

export const useSendStore = create<SendState>((set) => ({
  isSending: false,
  logs: [],
  lastError: null,

  sendSingle: async (data) => {
    set({ isSending: true, lastError: null })
    const logId = `${Date.now()}`
    try {
      await sendService.sendSingle(data)
      const newLog: DispatchLogItem = {
        id: logId,
        type: 'single',
        target: data.to,
        subject: data.subject,
        status: 'queued',
        timestamp: new Date().toLocaleTimeString(),
      }
      set((s) => ({ isSending: false, logs: [newLog, ...s.logs] }))
    } catch (err) {
      const errorMsg = err instanceof Error ? err.message : 'Send single failed'
      const newLog: DispatchLogItem = {
        id: logId,
        type: 'single',
        target: data.to,
        subject: data.subject,
        status: 'error',
        error: errorMsg,
        timestamp: new Date().toLocaleTimeString(),
      }
      set((s) => ({ isSending: false, lastError: errorMsg, logs: [newLog, ...s.logs] }))
      throw err
    }
  },

  sendSingleWithTemplate: async (data) => {
    set({ isSending: true, lastError: null })
    const logId = `${Date.now()}`
    try {
      await sendService.sendSingleWithTemplate(data)
      const newLog: DispatchLogItem = {
        id: logId,
        type: 'single_template',
        target: data.to,
        subject: data.subject,
        templateId: data.template_id,
        status: 'queued',
        timestamp: new Date().toLocaleTimeString(),
      }
      set((s) => ({ isSending: false, logs: [newLog, ...s.logs] }))
    } catch (err) {
      const errorMsg = err instanceof Error ? err.message : 'Send single template failed'
      const newLog: DispatchLogItem = {
        id: logId,
        type: 'single_template',
        target: data.to,
        subject: data.subject,
        templateId: data.template_id,
        status: 'error',
        error: errorMsg,
        timestamp: new Date().toLocaleTimeString(),
      }
      set((s) => ({ isSending: false, lastError: errorMsg, logs: [newLog, ...s.logs] }))
      throw err
    }
  },

  sendGroup: async (data) => {
    set({ isSending: true, lastError: null })
    const logId = `${Date.now()}`
    try {
      await sendService.sendGroup(data)
      const newLog: DispatchLogItem = {
        id: logId,
        type: 'group',
        target: `group:${data.group_id}`,
        subject: data.subject,
        status: 'queued',
        timestamp: new Date().toLocaleTimeString(),
      }
      set((s) => ({ isSending: false, logs: [newLog, ...s.logs] }))
    } catch (err) {
      const errorMsg = err instanceof Error ? err.message : 'Send group failed'
      const newLog: DispatchLogItem = {
        id: logId,
        type: 'group',
        target: `group:${data.group_id}`,
        subject: data.subject,
        status: 'error',
        error: errorMsg,
        timestamp: new Date().toLocaleTimeString(),
      }
      set((s) => ({ isSending: false, lastError: errorMsg, logs: [newLog, ...s.logs] }))
      throw err
    }
  },

  sendGroupWithTemplate: async (data) => {
    set({ isSending: true, lastError: null })
    const logId = `${Date.now()}`
    try {
      await sendService.sendGroupWithTemplate(data)
      const newLog: DispatchLogItem = {
        id: logId,
        type: 'group_template',
        target: `group:${data.group_id}`,
        subject: data.subject,
        templateId: data.template_id,
        status: 'queued',
        timestamp: new Date().toLocaleTimeString(),
      }
      set((s) => ({ isSending: false, logs: [newLog, ...s.logs] }))
    } catch (err) {
      const errorMsg = err instanceof Error ? err.message : 'Send group template failed'
      const newLog: DispatchLogItem = {
        id: logId,
        type: 'group_template',
        target: `group:${data.group_id}`,
        subject: data.subject,
        templateId: data.template_id,
        status: 'error',
        error: errorMsg,
        timestamp: new Date().toLocaleTimeString(),
      }
      set((s) => ({ isSending: false, lastError: errorMsg, logs: [newLog, ...s.logs] }))
      throw err
    }
  },

  clearLogs: () => set({ logs: [] }),
}))
