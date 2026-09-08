import React, { useEffect } from 'react'
import styles from './Toast.module.css'
import { useUiStore, type ToastMessage } from '../../stores/useUiStore'

interface ToastItemProps {
  toast: ToastMessage
  onDismiss: (id: string) => void
}

const ToastItem: React.FC<ToastItemProps> = ({ toast, onDismiss }) => {
  useEffect(() => {
    const duration = toast.durationMs ?? 3000
    const timer = setTimeout(() => {
      onDismiss(toast.id)
    }, duration)

    return () => clearTimeout(timer)
  }, [toast.id, toast.durationMs, onDismiss])

  return (
    <div
      className={`${styles.toast} ${styles[toast.type]}`}
      onClick={() => onDismiss(toast.id)}
      title="Click to dismiss"
    >
      <div className={styles.header}>
        <span className={styles.typeLabel}>[{toast.type.toUpperCase()}]</span>
        <span className={styles.time}>{toast.timestamp}</span>
        <button
          type="button"
          className={styles.closeBtn}
          onClick={(e) => {
            e.stopPropagation()
            onDismiss(toast.id)
          }}
          aria-label="Dismiss alert"
        >
          x
        </button>
      </div>
      <div className={styles.text}>{toast.text}</div>
    </div>
  )
}

export const ToastContainer: React.FC = () => {
  const { toasts, removeToast } = useUiStore()

  if (toasts.length === 0) return null

  return (
    <div className={styles.container} role="region" aria-label="Terminal Alerts">
      {toasts.map((toast: ToastMessage) => (
        <ToastItem key={toast.id} toast={toast} onDismiss={removeToast} />
      ))}
    </div>
  )
}

