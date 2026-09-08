export interface SendSingleRequest {
  to: string
  subject: string
  body: string
}

export interface SendSingleWithTemplateRequest {
  to: string
  subject: string
  template_id: string
  data: Record<string, unknown>
}

export interface SendGroupRequest {
  group_id: string
  subject: string
  body: string
}

export interface SendGroupWithTemplateRequest {
  group_id: string
  subject: string
  template_id: string
  data: Record<string, unknown>
}

export interface DispatchLogItem {
  id: string
  type: 'single' | 'single_template' | 'group' | 'group_template'
  target: string
  subject: string
  templateId?: string
  status: 'queued' | 'error'
  timestamp: string
  error?: string
}
