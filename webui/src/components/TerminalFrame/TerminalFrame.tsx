import React, { useState, useEffect } from 'react'
import styles from './TerminalFrame.module.css'
import { useUiStore, type NavTab } from '../../stores/useUiStore'
import { useSettingsStore, type TerminalTheme } from '../../stores/useSettingsStore'
import { useHealthStore } from '../../stores/useHealthStore'
import { Badge } from '../Badge/Badge'
import { ToastContainer } from '../Toast/ToastContainer'

interface TerminalFrameProps {
  children: React.ReactNode
}

export const TerminalFrame: React.FC<TerminalFrameProps> = ({ children }) => {
  const { activeTab, setActiveTab } = useUiStore()
  const { theme, setTheme } = useSettingsStore()
  const { data: health, error: healthError, fetchHealth } = useHealthStore()
  const [currentTime, setCurrentTime] = useState(new Date().toLocaleTimeString())

  useEffect(() => {
    // Apply theme attribute to document element
    document.documentElement.setAttribute('data-theme', theme)
  }, [theme])

  useEffect(() => {
    fetchHealth()
    const timer = setInterval(() => {
      setCurrentTime(new Date().toLocaleTimeString())
    }, 1000)
    return () => clearInterval(timer)
  }, [fetchHealth])

  const tabs: { key: NavTab; label: string }[] = [
    { key: 'dashboard', label: 'Dashboard' },
    { key: 'send', label: 'Dispatch' },
    { key: 'receivers', label: 'Receivers' },
    { key: 'groups', label: 'Groups' },
    { key: 'templates', label: 'Templates' },
    { key: 'logs', label: 'Logs' },
    { key: 'settings', label: 'Settings' },
  ]

  const isOnline = health?.status === 'up' && !healthError

  return (
    <div className={styles.terminalWindow}>
      {/* Titlebar */}
      <header className={styles.titlebar}>
        <div className={styles.controls}>
          <span className={`${styles.dot} ${styles.red}`} title="Close" />
          <span className={`${styles.dot} ${styles.yellow}`} title="Minimize" />
          <span className={`${styles.dot} ${styles.green}`} title="Zoom" />
        </div>

        <div className={styles.windowTitle}>
          <span className={styles.cliPrompt}>carrpigeon-cli</span>
          <span className={styles.separator}>—</span>
          <span className={styles.path}>~/carrpigeon/{activeTab}</span>
        </div>

        <div className={styles.statusGroup}>
          <Badge variant={isOnline ? 'green' : 'red'}>
            {isOnline ? 'SYS: ONLINE' : 'SYS: OFFLINE'}
          </Badge>
          <select
            value={theme}
            onChange={(e) => setTheme(e.target.value as TerminalTheme)}
            className={styles.themeSelect}
            aria-label="Select theme"
          >
            <option value="term-dark">Dark</option>
            <option value="term-light">Light</option>
          </select>
        </div>
      </header>

      {/* Navigation tabs */}
      <nav className={styles.tabBar} aria-label="Main Navigation">
        <div className={styles.tabList}>
          {tabs.map((tab) => {
            const isActive = activeTab === tab.key
            return (
              <button
                key={tab.key}
                type="button"
                className={`${styles.tabBtn} ${isActive ? styles.tabActive : ''}`}
                onClick={() => setActiveTab(tab.key)}
                aria-current={isActive ? 'page' : undefined}
              >
                <span className={styles.tabLabel}>{tab.label}</span>
              </button>
            )
          })}
        </div>
        <div className={styles.timeTicker}>{currentTime}</div>
      </nav>

      {/* Main Terminal Body */}
      <main className={styles.contentArea}>{children}</main>

      {/* Status Bar */}
      <footer className={styles.statusBar}>
        <div className={styles.statusLeft}>
          <span>ENV: prod</span>
          <span className={styles.delim}>|</span>
          <span>HTTP: localhost:2500</span>
          <span className={styles.delim}>|</span>
          <span>DB: {health ? `${health.in_use || 0} in use` : 'unknown'}</span>
        </div>
        <div className={styles.statusRight}>
          <span>CARRPIGEON NOTIFICATION SERVICE</span>
        </div>
      </footer>

      {/* Global Toast Alerts */}
      <ToastContainer />
    </div>
  )
}
