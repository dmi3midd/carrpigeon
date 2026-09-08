import React from 'react'
import styles from './Settings.module.css'
import { AsciiCard } from '../../components/AsciiCard/AsciiCard'
import { Toggle } from '../../components/Toggle/Toggle'
import { Badge } from '../../components/Badge/Badge'
import { useSettingsStore, type TerminalTheme } from '../../stores/useSettingsStore'
import { useHealthStore } from '../../stores/useHealthStore'
import { useUiStore } from '../../stores/useUiStore'

export const Settings: React.FC = () => {
  const {
    theme,
    setTheme,
    autoRefresh,
    setAutoRefresh,
    refreshIntervalSeconds,
    setRefreshIntervalSeconds,
  } = useSettingsStore()

  const { data: health, fetchHealth } = useHealthStore()
  const { addToast } = useUiStore()

  const themes: { id: TerminalTheme; name: string; desc: string; sample: string }[] = [
    {
      id: 'term-dark',
      name: 'Dark Terminal (Default)',
      desc: 'High contrast dark background with crisp green and cyan accents',
      sample: '#0b0d10 / #7ed957',
    },
    {
      id: 'term-light',
      name: 'Light Terminal',
      desc: 'Clean paper-style light terminal theme for daytime readability',
      sample: '#f6f8fa / #1a7f37',
    },
  ]


  const handlePing = async () => {
    await fetchHealth()
    addToast('info', 'Pinged Carrpigeon API backend')
  }

  return (
    <div className={styles.settingsPage}>
      {/* Visual Terminal Themes */}
      <AsciiCard
        title="Terminal Color Theme"
      >
        <div className={styles.themeGrid}>
          {themes.map((t) => {
            const isSelected = theme === t.id
            return (
              <div
                key={t.id}
                className={`${styles.themeCard} ${isSelected ? styles.themeActive : ''}`}
                onClick={() => setTheme(t.id)}
              >
                <div className={styles.themeHeader}>
                  <strong className={styles.themeName}>{t.name}</strong>
                  {isSelected && <Badge variant="green">ACTIVE</Badge>}
                </div>
                <p className={styles.themeDesc}>{t.desc}</p>
                <div className={styles.themeSample}>
                  <code>{t.sample}</code>
                  <button
                    type="button"
                    className={styles.applyBtn}
                    onClick={(e) => {
                      e.stopPropagation()
                      setTheme(t.id)
                    }}
                  >
                    {isSelected ? '[ Applied ]' : '[ Apply Theme ]'}
                  </button>
                </div>
              </div>
            )
          })}
        </div>
      </AsciiCard>

      {/* Connectivity & Auto-Refresh */}
      <div className={styles.twoCol}>
        <AsciiCard title="Backend API & Documentation">
          <div className={styles.configList}>
            <div className={styles.configItem}>
              <span className={styles.configLabel}>API Base URL:</span>
              <span className={styles.configValue}>
                <code>http://localhost:2500</code> (relative <code>/</code>)
              </span>
            </div>
            <div className={styles.configItem}>
              <span className={styles.configLabel}>Swagger Documentation:</span>
              <span className={styles.configValue}>
                <a
                  href="/swagger/index.html"
                  target="_blank"
                  rel="noreferrer"
                  className={styles.swaggerLink}
                >
                  [ Open /swagger/index.html &gt; ]
                </a>
              </span>
            </div>
            <div className={styles.configItem}>
              <span className={styles.configLabel}>Backend Status:</span>
              <span className={styles.configValue}>
                <Badge variant={health?.status === 'up' ? 'green' : 'red'}>
                  {health?.status === 'up' ? 'ONLINE (200 OK)' : 'OFFLINE'}
                </Badge>
                <button
                  type="button"
                  onClick={handlePing}
                  className={styles.pingBtn}
                >
                  [ Ping ]
                </button>
              </span>
            </div>
          </div>
        </AsciiCard>

        <AsciiCard title="Auto-Refresh & Polling">
          <div className={styles.configList}>
            <Toggle
              label="Enable Background Polling"
              checked={autoRefresh}
              onChange={setAutoRefresh}
            />
            <div className={styles.configItem}>
              <label htmlFor="refresh-interval" className={styles.configLabel}>
                Refresh Interval:
              </label>
              <select
                id="refresh-interval"
                value={refreshIntervalSeconds}
                onChange={(e) => setRefreshIntervalSeconds(Number(e.target.value))}
                className={styles.intervalSelect}
              >
                <option value={5}>5 seconds</option>
                <option value={10}>10 seconds</option>
                <option value={30}>30 seconds</option>
                <option value={60}>60 seconds</option>
              </select>
            </div>
          </div>
        </AsciiCard>
      </div>
    </div>
  )
}
