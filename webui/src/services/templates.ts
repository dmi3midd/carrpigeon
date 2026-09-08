import { request } from './api'
import type {
  CreateTemplateResponse,
  GetTemplateRawResponse,
  ListTemplateMetadataResponse,
  Template,
  TemplateMetadata,
  UpdateTemplateResponse,
} from '../types/template'

export const templatesService = {
  list: async (limit = 10, offset = 0): Promise<TemplateMetadata[]> => {
    const res = await request<ListTemplateMetadataResponse>('/templates', {
      params: { limit, offset },
    })
    return res.templates || []
  },

  getRaw: async (id: string): Promise<Template> => {
    const res = await request<GetTemplateRawResponse>(`/templates/${encodeURIComponent(id)}`)
    return res.template
  },

  create: async (name: string, file: File): Promise<CreateTemplateResponse> => {
    const formData = new FormData()
    formData.append('name', name)
    formData.append('file', file)

    return request<CreateTemplateResponse>('/templates', {
      method: 'POST',
      body: formData,
    })
  },

  update: async (id: string, name: string, file: File): Promise<UpdateTemplateResponse> => {
    const formData = new FormData()
    formData.append('name', name)
    formData.append('file', file)

    return request<UpdateTemplateResponse>(`/templates/${encodeURIComponent(id)}`, {
      method: 'PUT',
      body: formData,
    })
  },

  delete: async (id: string): Promise<void> => {
    await request<void>(`/templates/${encodeURIComponent(id)}`, {
      method: 'DELETE',
    })
  },
}
