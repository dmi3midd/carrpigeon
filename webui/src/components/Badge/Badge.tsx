import React from 'react'
import styles from './Badge.module.css'

export type BadgeVariant = 'green' | 'amber' | 'red' | 'cyan' | 'muted'

interface BadgeProps {
  variant?: BadgeVariant
  children: React.ReactNode
  showBrackets?: boolean
  className?: string
}

export const Badge: React.FC<BadgeProps> = ({
  variant = 'muted',
  children,
  showBrackets = true,
  className = '',
}) => {
  return (
    <span className={`${styles.badge} ${styles[variant]} ${className}`}>
      {showBrackets ? `[${children}]` : children}
    </span>
  )
}
