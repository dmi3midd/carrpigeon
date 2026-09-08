import React from 'react'
import styles from './Toggle.module.css'

interface ToggleProps {
  label: string
  checked: boolean
  onChange: (checked: boolean) => void
  disabled?: boolean
  id?: string
}

export const Toggle: React.FC<ToggleProps> = ({
  label,
  checked,
  onChange,
  disabled = false,
  id,
}) => {
  const toggleId = id || `toggle-${label.toLowerCase().replace(/\s+/g, '-')}`

  return (
    <div className={styles.container}>
      <label htmlFor={toggleId} className={styles.label}>
        {label}
      </label>
      <button
        id={toggleId}
        type="button"
        role="switch"
        aria-checked={checked}
        disabled={disabled}
        className={`${styles.toggleBtn} ${checked ? styles.checked : ''}`}
        onClick={() => onChange(!checked)}
      >
        {checked ? '[ ON ]' : '[ OFF ]'}
      </button>
    </div>
  )
}
