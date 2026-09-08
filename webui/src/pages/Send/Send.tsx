import React, { useState, useEffect } from 'react'
import styles from './Send.module.css'
import { AsciiCard } from '../../components/AsciiCard/AsciiCard'
import { Badge } from '../../components/Badge/Badge'
import { useSendStore } from '../../stores/useSendStore'
import { useReceiversStore } from '../../stores/useReceiversStore'
import { useGroupsStore } from '../../stores/useGroupsStore'
import { useTemplatesStore } from '../../stores/useTemplatesStore'
import { useUiStore } from '../../stores/useUiStore'

type SendMode = 'single' | 'single_template' | 'group' | 'group_template'

export const Send: React.FC = () => {
  const [mode, setMode] = useState<SendMode>('single')

  // Form fields
  const [to, setTo] = useState('')
  const [groupId, setGroupId] = useState('')
  const [subject, setSubject] = useState('')
  const [body, setBody] = useState('')
  const [templateId, setTemplateId] = useState('')
  const [templateDataJson, setTemplateDataJson] = useState('{\n  "Name": "John Doe"\n}')
  const [jsonError, setJsonError] = useState<string | null>(null)

  const { isSending, sendSingle, sendSingleWithTemplate, sendGroup, sendGroupWithTemplate } = useSendStore()
  const { receivers, fetchReceivers } = useReceiversStore()
  const { groups, fetchGroups } = useGroupsStore()
  const { templates, fetchTemplates } = useTemplatesStore()
  const { addToast, setActiveTab } = useUiStore()

  useEffect(() => {
    fetchReceivers(0)
    fetchGroups(0)
    fetchTemplates(0)
  }, [fetchReceivers, fetchGroups, fetchTemplates])

  // If a template is selected, initialize template data fields if available
  const handleTemplateChange = (id: string) => {
    setTemplateId(id)
    const tmpl = templates.find((t) => t.id === id)
    if (tmpl && tmpl.fields && tmpl.fields.length > 0) {
      const initialData: Record<string, string> = {}
      tmpl.fields.forEach((f) => {
        initialData[f] = `Value for ${f}`
      })
      setTemplateDataJson(JSON.stringify(initialData, null, 2))
    }
  }

  // Generate payload preview
  const getPayloadPreview = () => {
    let dataObj: Record<string, unknown> = {}
    try {
      if (mode.includes('template')) {
        dataObj = JSON.parse(templateDataJson)
      }
    } catch {
      dataObj = { error: 'Invalid JSON' }
    }

    if (mode === 'single') {
      return { to, subject, body }
    } else if (mode === 'single_template') {
      return { to, subject, template_id: templateId, data: dataObj }
    } else if (mode === 'group') {
      return { group_id: groupId, subject, body }
    } else {
      return { group_id: groupId, subject, template_id: templateId, data: dataObj }
    }
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setJsonError(null)

    try {
      if (mode === 'single') {
        if (!to) throw new Error('Recipient email is required')
        if (!subject) throw new Error('Subject is required')
        if (!body) throw new Error('Body is required')
        await sendSingle({ to, subject, body })
        addToast('success', `Email queued for ${to}`)
      } else if (mode === 'single_template') {
        if (!to) throw new Error('Recipient email is required')
        if (!subject) throw new Error('Subject is required')
        if (!templateId) throw new Error('Template ID is required')
        let data: Record<string, unknown> = {}
        try {
          data = JSON.parse(templateDataJson)
        } catch {
          setJsonError('Template Data is not valid JSON')
          throw new Error('Template Data is not valid JSON')
        }
        await sendSingleWithTemplate({ to, subject, template_id: templateId, data })
        addToast('success', `Templated email queued for ${to}`)
      } else if (mode === 'group') {
        if (!groupId) throw new Error('Group ID is required')
        if (!subject) throw new Error('Subject is required')
        if (!body) throw new Error('Body is required')
        await sendGroup({ group_id: groupId, subject, body })
        addToast('success', `Group email queued for group ${groupId}`)
      } else if (mode === 'group_template') {
        if (!groupId) throw new Error('Group ID is required')
        if (!subject) throw new Error('Subject is required')
        if (!templateId) throw new Error('Template ID is required')
        let data: Record<string, unknown> = {}
        try {
          data = JSON.parse(templateDataJson)
        } catch {
          setJsonError('Template Data is not valid JSON')
          throw new Error('Template Data is not valid JSON')
        }
        await sendGroupWithTemplate({ group_id: groupId, subject, template_id: templateId, data })
        addToast('success', `Group templated email queued for group ${groupId}`)
      }

      // Clear subject and body on success
      setSubject('')
      setBody('')
    } catch (err) {
      addToast('error', err instanceof Error ? err.message : 'Dispatch failed')
    }
  }

  return (
    <div className={styles.sendPage}>
      {/* Mode Selection Header */}
      <div className={styles.modeBar}>
        <span className={styles.modeLabel}>[ DISPATCH MODE ]:</span>
        <div className={styles.modeButtons}>
          <button
            type="button"
            className={`${styles.modeBtn} ${mode === 'single' ? styles.modeActive : ''}`}
            onClick={() => setMode('single')}
          >
            [1. Single Email]
          </button>
          <button
            type="button"
            className={`${styles.modeBtn} ${mode === 'single_template' ? styles.modeActive : ''}`}
            onClick={() => setMode('single_template')}
          >
            [2. Single + Template]
          </button>
          <button
            type="button"
            className={`${styles.modeBtn} ${mode === 'group' ? styles.modeActive : ''}`}
            onClick={() => setMode('group')}
          >
            [3. Group Email]
          </button>
          <button
            type="button"
            className={`${styles.modeBtn} ${mode === 'group_template' ? styles.modeActive : ''}`}
            onClick={() => setMode('group_template')}
          >
            [4. Group + Template]
          </button>
        </div>
      </div>

      <div className={styles.splitLayout}>
        {/* Form area */}
        <AsciiCard
          title={`Email Dispatch Form: ${mode.replace('_', ' ').toUpperCase()}`}
          badge={<Badge variant="cyan">{mode.toUpperCase()}</Badge>}
          className={styles.formCard}
        >
          <form onSubmit={handleSubmit} className={styles.form}>
            {/* Recipient Target */}
            {mode.startsWith('single') ? (
              <div className={styles.fieldGroup}>
                <label htmlFor="send-to" className={styles.fieldLabel}>
                  Recipient Email [to]:
                </label>
                <div className={styles.inputWithSelect}>
                  <input
                    id="send-to"
                    type="email"
                    required
                    placeholder="e.g. user@example.com"
                    value={to}
                    onChange={(e) => setTo(e.target.value)}
                  />
                  {receivers.length > 0 && (
                    <select
                      onChange={(e) => e.target.value && setTo(e.target.value)}
                      className={styles.quickSelect}
                      aria-label="Select from existing receivers"
                    >
                      <option value="">[ Pick Receiver ]</option>
                      {receivers.map((r) => (
                        <option key={r.id} value={r.email}>
                          {r.name} ({r.email})
                        </option>
                      ))}
                    </select>
                  )}
                </div>
              </div>
            ) : (
              <div className={styles.fieldGroup}>
                <label htmlFor="send-group" className={styles.fieldLabel}>
                  Target Group [group_id]:
                </label>
                <div className={styles.inputWithSelect}>
                  <input
                    id="send-group"
                    type="text"
                    required
                    placeholder="Group ID (xid)"
                    value={groupId}
                    onChange={(e) => setGroupId(e.target.value)}
                  />
                  {groups.length > 0 && (
                    <select
                      onChange={(e) => e.target.value && setGroupId(e.target.value)}
                      className={styles.quickSelect}
                      aria-label="Select from existing groups"
                    >
                      <option value="">[ Pick Group ]</option>
                      {groups.map((g) => (
                        <option key={g.id} value={g.id}>
                          {g.name} ({g.receivers_count} members)
                        </option>
                      ))}
                    </select>
                  )}
                </div>
              </div>
            )}

            {/* Template selector if template mode */}
            {mode.includes('template') && (
              <div className={styles.fieldGroup}>
                <label htmlFor="send-template" className={styles.fieldLabel}>
                  Template ID [template_id]:
                </label>
                <div className={styles.inputWithSelect}>
                  <input
                    id="send-template"
                    type="text"
                    required
                    placeholder="Template ID (xid)"
                    value={templateId}
                    onChange={(e) => handleTemplateChange(e.target.value)}
                  />
                  {templates.length > 0 && (
                    <select
                      onChange={(e) => e.target.value && handleTemplateChange(e.target.value)}
                      className={styles.quickSelect}
                      aria-label="Select registered template"
                    >
                      <option value="">[ Pick Template ]</option>
                      {templates.map((t) => (
                        <option key={t.id} value={t.id}>
                          {t.name} ({t.is_html ? 'HTML' : 'TXT'}
                          {t.fields?.length ? `, fields: ${t.fields.join(',')}` : ''})
                        </option>
                      ))}
                    </select>
                  )}
                </div>
              </div>
            )}

            {/* Subject */}
            <div className={styles.fieldGroup}>
              <label htmlFor="send-subject" className={styles.fieldLabel}>
                Subject:
              </label>
              <input
                id="send-subject"
                type="text"
                required
                maxLength={256}
                placeholder="Subject line..."
                value={subject}
                onChange={(e) => setSubject(e.target.value)}
              />
            </div>

            {/* Body or Template Data */}
            {!mode.includes('template') ? (
              <div className={styles.fieldGroup}>
                <label htmlFor="send-body" className={styles.fieldLabel}>
                  Message Body (Plain text, max 2048 chars):
                </label>
                <textarea
                  id="send-body"
                  rows={6}
                  required
                  maxLength={2048}
                  placeholder="Enter message content..."
                  value={body}
                  onChange={(e) => setBody(e.target.value)}
                />
              </div>
            ) : (
              <div className={styles.fieldGroup}>
                <div className={styles.labelRow}>
                  <label htmlFor="send-data" className={styles.fieldLabel}>
                    Template Variables [JSON data]:
                  </label>
                  {jsonError && <span className={styles.errLabel}>[ERR: {jsonError}]</span>}
                </div>
                <textarea
                  id="send-data"
                  rows={6}
                  required
                  placeholder={'{\n  "Name": "Jane Doe"\n}'}
                  value={templateDataJson}
                  onChange={(e) => {
                    setTemplateDataJson(e.target.value)
                    setJsonError(null)
                  }}
                  className={styles.jsonTextarea}
                />
              </div>
            )}

            {/* Submit Button */}
            <div className={styles.formActions}>
              <button
                type="submit"
                disabled={isSending}
                className={styles.submitBtn}
              >
                {isSending ? '[ TRANSMITTING... ]' : '[ > TRANSMIT NOTIFICATION ]'}
              </button>
            </div>
          </form>
        </AsciiCard>

        {/* Live Payload Preview & Guidance */}
        <div className={styles.sideCol}>
          <AsciiCard title="HTTP Request Payload Preview" subtitle="POST /send/...">
            <pre className={styles.previewBox}>
              {JSON.stringify(getPayloadPreview(), null, 2)}
            </pre>
          </AsciiCard>

          <AsciiCard title="Terminal Quick Reference">
            <div className={styles.refInfo}>
              <p>
                <strong>Endpoint:</strong>{' '}
                <code>
                  {mode === 'single' && 'POST /send/single'}
                  {mode === 'single_template' && 'POST /send/single/template'}
                  {mode === 'group' && 'POST /send/group'}
                  {mode === 'group_template' && 'POST /send/group/template'}
                </code>
              </p>
              <p>
                <strong>Response:</strong> <code>202 Accepted</code>. Messages are
                queued asynchronously for background processing by the Carrpigeon worker pool.
              </p>
              <div className={styles.linkRow}>
                <button
                  type="button"
                  className={styles.linkBtn}
                  onClick={() => setActiveTab('logs')}
                >
                  [ View Dispatch Logs &gt; ]
                </button>
              </div>
            </div>
          </AsciiCard>
        </div>
      </div>
    </div>
  )
}
