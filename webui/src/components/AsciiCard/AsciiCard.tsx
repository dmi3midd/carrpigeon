import React from 'react'
import styles from './AsciiCard.module.css'

interface AsciiCardProps {
  title: string
  subtitle?: string
  badge?: React.ReactNode
  actions?: React.ReactNode
  children: React.ReactNode
  className?: string
  noPadding?: boolean
}

export const AsciiCard: React.FC<AsciiCardProps> = ({
  title,
  subtitle,
  badge,
  actions,
  children,
  className = '',
  noPadding = false,
}) => {
  return (
    <div className={`${styles.card} ${className}`}>
      <div className={styles.header}>
        <div className={styles.titleGroup}>
          <span className={styles.bracket}>+--[</span>
          <h3 className={styles.title}>{title}</h3>
          <span className={styles.bracket}>]--</span>
          {subtitle && <span className={styles.subtitle}>{subtitle}</span>}
          {badge && <div className={styles.badgeWrapper}>{badge}</div>}
        </div>
        {actions && <div className={styles.actions}>{actions}</div>}
      </div>
      <div className={`${styles.body} ${noPadding ? styles.noPadding : ''}`}>
        {children}
      </div>
    </div>
  )
}
