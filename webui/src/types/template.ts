export interface Template {
  id: string
  name: string
  content: string
  is_html: boolean
  fields: string[] | null
  created_at: string
  updated_at: string
}

export interface TemplateMetadata {
  id: string
  name: string
  is_html: boolean
  fields: string[] | null
  created_at: string
  updated_at: string
}

export interface GetTemplateRawResponse {
  template: Template
}

export interface ListTemplateMetadataResponse {
  templates: TemplateMetadata[] | null
}

export interface CreateTemplateResponse {
  id: string
}

export interface UpdateTemplateResponse {
  id: string
}
