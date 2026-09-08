import React from 'react'
import styles from './Dashboard.module.css'
import { CARRPIGEON_ASCII } from '../../assets/carrpigeon-ascii'
import { AsciiCard } from '../../components/AsciiCard/AsciiCard'
import { Badge } from '../../components/Badge/Badge'
import { useHealthStore } from '../../stores/useHealthStore'
import { useUiStore } from '../../stores/useUiStore'
import { useSendStore } from '../../stores/useSendStore'

export const Dashboard: React.FC = () => {
  const { data: health, isLoading: isHealthLoading, fetchHealth } = useHealthStore()
  const { setActiveTab } = useUiStore()
  const { logs } = useSendStore()

  const isUp = health?.status === 'up'

  return (
    <div className={styles.dashboard}>
      {/* ASCII Hero Banner */}
      <div className={styles.asciiBanner}>
        <pre>{CARRPIGEON_ASCII}</pre>
        <div className={styles.bannerSub}>
          [ Carrpigeon Notification Service — Email Delivery & Template Engine ]
        </div>
      </div>

      {/* Grid of System Status Cards */}
      <div className={styles.grid}>
        {/* PostgreSQL Database Health */}
        <AsciiCard
          title="System Health & DB Pool"
          badge={
            <Badge variant={isUp ? 'green' : 'red'}>
              {isUp ? 'STATUS: HEALTHY' : 'STATUS: DEGRADED'}
            </Badge>
          }
          actions={
            <button
              type="button"
              onClick={() => fetchHealth()}
              disabled={isHealthLoading}
              className={styles.refreshBtn}
            >
              {isHealthLoading ? '[ Ping... ]' : '[ Ping DB ]'}
            </button>
          }
        >
          <div className={styles.metricGrid}>
            <div className={styles.metricItem}>
              <span className={styles.metricLabel}>Postgres Status:</span>
              <span className={`${styles.metricVal} ${isUp ? styles.valGreen : styles.valRed}`}>
                {health?.status || 'unknown'} ({health?.message || 'checking...'})
              </span>
            </div>
            <div className={styles.metricItem}>
              <span className={styles.metricLabel}>Open Connections:</span>
              <span className={styles.metricVal}>{health?.open_connections || '0'}</span>
            </div>
            <div className={styles.metricItem}>
              <span className={styles.metricLabel}>In-Use Connections:</span>
              <span className={styles.metricVal}>{health?.in_use || '0'}</span>
            </div>
            <div className={styles.metricItem}>
              <span className={styles.metricLabel}>Idle Connections:</span>
              <span className={styles.metricVal}>{health?.idle || '0'}</span>
            </div>
            <div className={styles.metricItem}>
              <span className={styles.metricLabel}>Wait Count:</span>
              <span className={styles.metricVal}>{health?.wait_count || '0'}</span>
            </div>
            <div className={styles.metricItem}>
              <span className={styles.metricLabel}>Wait Duration:</span>
              <span className={styles.metricVal}>{health?.wait_duration || '0s'}</span>
            </div>
          </div>
        </AsciiCard>

        {/* Quick Actions */}
        <AsciiCard title="Quick Actions & Operations">
          <div className={styles.actionList}>
            <button
              type="button"
              className={styles.actionBtn}
              onClick={() => setActiveTab('send')}
            >
              <span className={styles.actionArrow}>&gt;</span>
              <div className={styles.actionText}>
                <strong>Dispatch Email</strong>
                <small>Send single, group, or templated emails</small>
              </div>
            </button>

            <button
              type="button"
              className={styles.actionBtn}
              onClick={() => setActiveTab('receivers')}
            >
              <span className={styles.actionArrow}>&gt;</span>
              <div className={styles.actionText}>
                <strong>Manage Receivers</strong>
                <small>View and add recipient email records</small>
              </div>
            </button>

            <button
              type="button"
              className={styles.actionBtn}
              onClick={() => setActiveTab('groups')}
            >
              <span className={styles.actionArrow}>&gt;</span>
              <div className={styles.actionText}>
                <strong>Receiver Groups</strong>
                <small>Create distribution lists & assign members</small>
              </div>
            </button>

            <button
              type="button"
              className={styles.actionBtn}
              onClick={() => setActiveTab('templates')}
            >
              <span className={styles.actionArrow}>&gt;</span>
              <div className={styles.actionText}>
                <strong>Templates Catalog</strong>
                <small>Upload and preview HTML/text templates</small>
              </div>
            </button>
          </div>
        </AsciiCard>
      </div>

      {/* Recent Dispatches Stream */}
      <AsciiCard
        title="Session Activity Log"
        subtitle={`(${logs.length} events)`}
        actions={
          <button
            type="button"
            className={styles.refreshBtn}
            onClick={() => setActiveTab('logs')}
          >
            [ Full Logs &gt; ]
          </button>
        }
      >
        {logs.length === 0 ? (
          <div className={styles.emptyLog}>
            [ No email dispatches recorded in this session. Go to Dispatch to send an email. ]
          </div>
        ) : (
          <div className={styles.logList}>
            {logs.slice(0, 5).map((log) => (
              <div key={log.id} className={styles.logRow}>
                <span className={styles.logTime}>[{log.timestamp}]</span>
                <Badge variant={log.status === 'queued' ? 'green' : 'red'}>
                  {log.status.toUpperCase()}
                </Badge>
                <span className={styles.logType}>[{log.type}]</span>
                <span className={styles.logTarget}>{log.target}</span>
                <span className={styles.logSubject}>"{log.subject}"</span>
                {log.error && <span className={styles.logError}>err: {log.error}</span>}
              </div>
            ))}
          </div>
        )}
      </AsciiCard>
    </div>
  )
}
