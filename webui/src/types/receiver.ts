export interface Receiver {
  id: string
  name: string
  email: string
  created_at: string
  updated_at: string
}

export interface GetReceiverByIdResponse {
  receiver: Receiver
}

export interface ListReceiversResponse {
  receivers: Receiver[] | null
}

export interface CreateReceiverRequest {
  name: string
  email: string
}

export interface CreateReceiverResponse {
  id: string
}

export interface UpdateReceiverRequest {
  name?: string
  email?: string
}

export interface UpdateReceiverResponse {
  id: string
}
