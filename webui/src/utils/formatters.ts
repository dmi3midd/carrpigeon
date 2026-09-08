export function truncateId(id: string, len = 8): string {
  if (!id) return ''
  if (id.length <= len + 3) return id
  return `${id.slice(0, len)}…`
}

export function formatDate(dateString: string | undefined): string {
  if (!dateString) return '—'
  try {
    const date = new Date(dateString)
    if (isNaN(date.getTime())) return dateString
    return date.toLocaleString(undefined, {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
      hour12: false,
    })
  } catch {
    return dateString
  }
}

export function formatTimeOnly(dateString: string | undefined): string {
  if (!dateString) return '—'
  try {
    const date = new Date(dateString)
    return date.toLocaleTimeString(undefined, { hour12: false })
  } catch {
    return dateString
  }
}
