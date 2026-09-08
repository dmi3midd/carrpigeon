export interface HealthResponse {
  status: string
  message: string
  open_connections: string
  in_use: string
  idle: string
  wait_count: string
  wait_duration: string
  max_idle_closed: string
  max_lifetime_closed: string
}
