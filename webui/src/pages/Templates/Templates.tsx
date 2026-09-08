import React, { useState, useEffect } from 'react'
import styles from './Templates.module.css'
import { AsciiCard } from '../../components/AsciiCard/AsciiCard'
import { Badge } from '../../components/Badge/Badge'
import { TerminalTable, type Column } from '../../components/TerminalTable/TerminalTable'
import { Modal } from '../../components/Modal/Modal'
import { useTemplatesStore } from '../../stores/useTemplatesStore'
import { useUiStore } from '../../stores/useUiStore'
import { truncateId, formatDate } from '../../utils/formatters'
import type { Template, TemplateMetadata } from '../../types/template'

export const Templates: React.FC = () => {
  const {
    templates,
    isLoading,
    offset,
    limit,
    hasMore,
    fetchTemplates,
    fetchRawTemplate,
    createTemplate,
    deleteTemplate,
    nextPage,
    prevPage,
  } = useTemplatesStore()

  const { addToast, setActiveTab } = useUiStore()

  const [search, setSearch] = useState('')

  // Modals state
  const [isUploadOpen, setIsUploadOpen] = useState(false)
  const [viewTarget, setViewTarget] = useState<Template | null>(null)
  const [deleteTarget, setDeleteTarget] = useState<TemplateMetadata | null>(null)

  // Upload Form states
  const [name, setName] = useState('')
  const [file, setFile] = useState<File | null>(null)
  const [isSubmitting, setIsSubmitting] = useState(false)

  useEffect(() => {
    fetchTemplates(0)
  }, [fetchTemplates])

  const handleOpenUpload = () => {
    setName('')
    setFile(null)
    setIsUploadOpen(true)
  }

  const handleUpload = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!file) {
      addToast('error', 'Please select a template file (.html or .txt)')
      return
    }
    setIsSubmitting(true)
    try {
      const id = await createTemplate(name, file)
      addToast('success', `Uploaded template '${name}' (ID: ${truncateId(id)})`)
      setIsUploadOpen(false)
    } catch (err) {
      addToast('error', err instanceof Error ? err.message : 'Upload failed')
    } finally {
      setIsSubmitting(false)
    }
  }

  const handleOpenView = async (item: TemplateMetadata) => {
    try {
      const tmpl = await fetchRawTemplate(item.id)
      setViewTarget(tmpl)
    } catch (err) {
      addToast('error', err instanceof Error ? err.message : 'Failed to fetch raw template')
    }
  }

  const handleDelete = async () => {
    if (!deleteTarget) return
    setIsSubmitting(true)
    try {
      await deleteTemplate(deleteTarget.id)
      addToast('success', `Deleted template '${deleteTarget.name}'`)
      setDeleteTarget(null)
    } catch (err) {
      addToast('error', err instanceof Error ? err.message : 'Failed to delete template')
    } finally {
      setIsSubmitting(false)
    }
  }

  const handleCopyId = (id: string) => {
    navigator.clipboard.writeText(id)
    addToast('info', `Copied Template ID: ${truncateId(id)}`)
  }

  const filteredTemplates = templates.filter(
    (t) =>
      t.name.toLowerCase().includes(search.toLowerCase()) ||
      t.id.toLowerCase().includes(search.toLowerCase()) ||
      (t.fields && t.fields.some((f) => f.toLowerCase().includes(search.toLowerCase())))
  )

  const columns: Column<TemplateMetadata>[] = [
    {
      key: 'id',
      header: 'ID',
      width: '180px',
      render: (item) => (
        <span
          className={styles.idCell}
          onClick={(e) => {
            e.stopPropagation()
            handleCopyId(item.id)
          }}
          title={`Click to copy: ${item.id}`}
        >
          {truncateId(item.id, 10)} <span className={styles.copyIcon}>[copy]</span>
        </span>
      ),
    },
    {
      key: 'name',
      header: 'NAME',
      render: (item) => <strong>{item.name}</strong>,
    },
    {
      key: 'is_html',
      header: 'TYPE',
      width: '90px',
      render: (item) => (
        <Badge variant={item.is_html ? 'cyan' : 'amber'}>
          {item.is_html ? 'HTML' : 'TXT'}
        </Badge>
      ),
    },
    {
      key: 'fields',
      header: 'DETECTED FIELDS',
      render: (item) => {
        if (!item.fields || item.fields.length === 0) {
          return <span className={styles.mutedText}>[ none ]</span>
        }
        return (
          <div className={styles.fieldsList}>
            {item.fields.map((f) => (
              <span key={f} className={styles.fieldPill}>
                {`{{.${f}}}`}
              </span>
            ))}
          </div>
        )
      },
    },
    {
      key: 'updated_at',
      header: 'UPDATED',
      width: '170px',
      render: (item) => formatDate(item.updated_at),
    },
    {
      key: 'actions',
      header: 'ACTIONS',
      width: '170px',
      align: 'right',
      render: (item) => (
        <div className={styles.rowActions}>
          <button
            type="button"
            className={styles.actionBtn}
            onClick={() => handleOpenView(item)}
            title="Inspect raw template & preview"
          >
            [view]
          </button>
          <button
            type="button"
            className={`${styles.actionBtn} ${styles.deleteBtn}`}
            onClick={() => setDeleteTarget(item)}
            title="Delete template"
          >
            [del]
          </button>
        </div>
      ),
    },
  ]

  return (
    <div className={styles.templatesPage}>
      <AsciiCard
        title="Templates Catalog"
        subtitle={`(${filteredTemplates.length} templates)`}
        badge={<Badge variant="green">ENGINE READY</Badge>}
        actions={
          <button
            type="button"
            className={styles.newBtn}
            onClick={handleOpenUpload}
          >
            [ + Upload Template ]
          </button>
        }
      >
        <div className={styles.toolbar}>
          <div className={styles.searchWrapper}>
            <span className={styles.promptArrow}>&gt;</span>
            <input
              type="text"
              placeholder="filter by template name, id, or fields..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className={styles.searchInput}
            />
          </div>
          <button
            type="button"
            onClick={() => fetchTemplates()}
            disabled={isLoading}
            className={styles.refreshBtn}
          >
            {isLoading ? '[ refreshing... ]' : '[ refresh ]'}
          </button>
        </div>

        <TerminalTable
          columns={columns}
          data={filteredTemplates}
          isLoading={isLoading}
          keyExtractor={(item) => item.id}
          emptyMessage="No templates found. Upload an HTML or TXT file with {{.Field}} placeholders."
        />

        <div className={styles.pagination}>
          <div className={styles.pageInfo}>
            Offset: {offset} | Page size: {limit}
          </div>
          <div className={styles.pageBtns}>
            <button
              type="button"
              disabled={offset === 0 || isLoading}
              onClick={prevPage}
            >
              [ &lt; Prev ]
            </button>
            <button
              type="button"
              disabled={!hasMore || isLoading}
              onClick={nextPage}
            >
              [ Next &gt; ]
            </button>
          </div>
        </div>
      </AsciiCard>

      {/* Upload Template Modal */}
      <Modal
        isOpen={isUploadOpen}
        title="Upload Template File"
        onClose={() => setIsUploadOpen(false)}
      >
        <form onSubmit={handleUpload} className={styles.modalForm}>
          <div className={styles.fieldGroup}>
            <label htmlFor="tmpl-name" className={styles.fieldLabel}>
              Template Name [name]:
            </label>
            <input
              id="tmpl-name"
              type="text"
              required
              maxLength={128}
              placeholder="e.g. Welcome Email"
              value={name}
              onChange={(e) => setName(e.target.value)}
              autoFocus
            />
          </div>
          <div className={styles.fieldGroup}>
            <label htmlFor="tmpl-file" className={styles.fieldLabel}>
              Template File (.html or .txt) [file]:
            </label>
            <input
              id="tmpl-file"
              type="file"
              required
              accept=".html,.htm,.txt"
              onChange={(e) => setFile(e.target.files?.[0] || null)}
            />
            <small className={styles.fileHint}>
              Template files can define placeholders such as {'{{.Name}}'} or {'{{.PromoCode}}'}.
            </small>
          </div>
          <div className={styles.modalActions}>
            <button
              type="button"
              onClick={() => setIsUploadOpen(false)}
              className={styles.cancelBtn}
            >
              [ Cancel ]
            </button>
            <button
              type="submit"
              disabled={isSubmitting || !file}
              className={styles.confirmBtn}
            >
              {isSubmitting ? '[ Uploading... ]' : '[ > Upload Template ]'}
            </button>
          </div>
        </form>
      </Modal>

      {/* View/Preview Template Modal */}
      <Modal
        isOpen={viewTarget !== null}
        title={`Template Inspector: ${viewTarget?.name}`}
        onClose={() => setViewTarget(null)}
        maxWidth="750px"
      >
        {viewTarget && (
          <div className={styles.viewerContainer}>
            <div className={styles.viewerMeta}>
              <div>
                <strong>ID:</strong> <code>{viewTarget.id}</code>
              </div>
              <div>
                <strong>Format:</strong>{' '}
                <Badge variant={viewTarget.is_html ? 'cyan' : 'amber'}>
                  {viewTarget.is_html ? 'HTML' : 'TEXT'}
                </Badge>
              </div>
              <div>
                <strong>Fields:</strong>{' '}
                {viewTarget.fields && viewTarget.fields.length > 0 ? (
                  viewTarget.fields.map((f) => (
                    <span key={f} className={styles.fieldPill}>
                      {`{{.${f}}}`}
                    </span>
                  ))
                ) : (
                  <span className={styles.mutedText}>None detected</span>
                )}
              </div>
            </div>

            <div className={styles.codeView}>
              <div className={styles.codeHeader}>[ RAW SOURCE CODE ]</div>
              <pre className={styles.codeBlock}>{viewTarget.content}</pre>
            </div>

            {viewTarget.is_html && (
              <div className={styles.previewBox}>
                <div className={styles.codeHeader}>[ RENDERED HTML PREVIEW ]</div>
                <iframe
                  srcDoc={viewTarget.content}
                  title="Template Preview"
                  sandbox="allow-same-origin"
                  className={styles.previewIframe}
                />
              </div>
            )}

            <div className={styles.modalActions}>
              <button
                type="button"
                className={styles.confirmBtn}
                onClick={() => {
                  setViewTarget(null)
                  setActiveTab('send')
                }}
              >
                {"[ > Dispatch with this template ]"}
              </button>
              <button
                type="button"
                onClick={() => setViewTarget(null)}
                className={styles.cancelBtn}
              >
                [ Close ]
              </button>
            </div>
          </div>
        )}
      </Modal>

      {/* Delete Confirmation Modal */}
      <Modal
        isOpen={deleteTarget !== null}
        title="Confirm Template Deletion"
        onClose={() => setDeleteTarget(null)}
      >
        <div className={styles.deletePrompt}>
          <p>
            Are you sure you want to delete template{' '}
            <strong>{deleteTarget?.name}</strong> (<code>{deleteTarget?.id}</code>)?
          </p>
          <div className={styles.modalActions}>
            <button
              type="button"
              onClick={() => setDeleteTarget(null)}
              className={styles.cancelBtn}
            >
              [ Cancel ]
            </button>
            <button
              type="button"
              disabled={isSubmitting}
              onClick={handleDelete}
              className={styles.deleteConfirmBtn}
            >
              {isSubmitting ? '[ Deleting... ]' : '[ > Confirm Delete ]'}
            </button>
          </div>
        </div>
      </Modal>
    </div>
  )
}
