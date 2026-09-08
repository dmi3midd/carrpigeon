import React from 'react'
import styles from './TerminalTable.module.css'

export interface Column<T> {
  key: string
  header: string
  width?: string
  align?: 'left' | 'center' | 'right'
  render?: (item: T, index: number) => React.ReactNode
}

interface TerminalTableProps<T> {
  columns: Column<T>[]
  data: T[]
  isLoading?: boolean
  emptyMessage?: string
  keyExtractor: (item: T) => string
  onRowClick?: (item: T) => void
}

export function TerminalTable<T>({
  columns,
  data,
  isLoading = false,
  emptyMessage = 'No records found.',
  keyExtractor,
  onRowClick,
}: TerminalTableProps<T>) {
  return (
    <div className={styles.tableWrapper}>
      <table className={styles.table}>
        <thead>
          <tr>
            {columns.map((col) => (
              <th
                key={col.key}
                style={{
                  width: col.width,
                  textAlign: col.align || 'left',
                }}
              >
                [{col.header}]
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {isLoading ? (
            <tr>
              <td colSpan={columns.length} className={styles.statusCell}>
                <span className={styles.loadingPulse}>[ LOADING DATA... ]</span>
              </td>
            </tr>
          ) : data.length === 0 ? (
            <tr>
              <td colSpan={columns.length} className={styles.statusCell}>
                [ {emptyMessage} ]
              </td>
            </tr>
          ) : (
            data.map((item, idx) => (
              <tr
                key={keyExtractor(item)}
                onClick={() => onRowClick?.(item)}
                className={onRowClick ? styles.clickableRow : ''}
              >
                {columns.map((col) => {
                  const val = (item as Record<string, unknown>)[col.key]
                  return (
                    <td
                      key={col.key}
                      style={{ textAlign: col.align || 'left' }}
                    >
                      {col.render ? col.render(item, idx) : String(val ?? '—')}
                    </td>
                  )
                })}
              </tr>
            ))
          )}
        </tbody>
      </table>
    </div>
  )
}
