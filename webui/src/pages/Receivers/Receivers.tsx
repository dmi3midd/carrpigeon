import React, { useState, useEffect } from 'react'
import styles from './Receivers.module.css'
import { AsciiCard } from '../../components/AsciiCard/AsciiCard'
import { Badge } from '../../components/Badge/Badge'
import { TerminalTable, type Column } from '../../components/TerminalTable/TerminalTable'
import { Modal } from '../../components/Modal/Modal'
import { useReceiversStore } from '../../stores/useReceiversStore'
import { useUiStore } from '../../stores/useUiStore'
import { truncateId, formatDate } from '../../utils/formatters'
import type { Receiver } from '../../types/receiver'

export const Receivers: React.FC = () => {
  const {
    receivers,
    isLoading,
    offset,
    limit,
    hasMore,
    fetchReceivers,
    createReceiver,
    updateReceiver,
    deleteReceiver,
    nextPage,
    prevPage,
  } = useReceiversStore()

  const { addToast } = useUiStore()

  const [search, setSearch] = useState('')

  // Modals state
  const [isCreateOpen, setIsCreateOpen] = useState(false)
  const [editTarget, setEditTarget] = useState<Receiver | null>(null)
  const [deleteTarget, setDeleteTarget] = useState<Receiver | null>(null)

  // Form states
  const [name, setName] = useState('')
  const [email, setEmail] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)

  useEffect(() => {
    fetchReceivers(0)
  }, [fetchReceivers])

  const handleOpenCreate = () => {
    setName('')
    setEmail('')
    setIsCreateOpen(true)
  }

  const handleOpenEdit = (receiver: Receiver) => {
    setEditTarget(receiver)
    setName(receiver.name)
    setEmail(receiver.email)
  }

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsSubmitting(true)
    try {
      const id = await createReceiver({ name, email })
      addToast('success', `Created receiver '${name}' (ID: ${truncateId(id)})`)
      setIsCreateOpen(false)
    } catch (err) {
      addToast('error', err instanceof Error ? err.message : 'Failed to create receiver')
    } finally {
      setIsSubmitting(false)
    }
  }

  const handleUpdate = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!editTarget) return
    setIsSubmitting(true)
    try {
      await updateReceiver(editTarget.id, { name, email })
      addToast('success', `Updated receiver '${name}'`)
      setEditTarget(null)
    } catch (err) {
      addToast('error', err instanceof Error ? err.message : 'Failed to update receiver')
    } finally {
      setIsSubmitting(false)
    }
  }

  const handleDelete = async () => {
    if (!deleteTarget) return
    setIsSubmitting(true)
    try {
      await deleteReceiver(deleteTarget.id)
      addToast('success', `Deleted receiver '${deleteTarget.name}'`)
      setDeleteTarget(null)
    } catch (err) {
      addToast('error', err instanceof Error ? err.message : 'Failed to delete receiver')
    } finally {
      setIsSubmitting(false)
    }
  }

  const handleCopyId = (id: string) => {
    navigator.clipboard.writeText(id)
    addToast('info', `Copied ID to clipboard: ${truncateId(id)}`)
  }

  const filteredReceivers = receivers.filter(
    (r) =>
      r.name.toLowerCase().includes(search.toLowerCase()) ||
      r.email.toLowerCase().includes(search.toLowerCase()) ||
      r.id.toLowerCase().includes(search.toLowerCase())
  )

  const columns: Column<Receiver>[] = [
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
          title={`Click to copy full ID: ${item.id}`}
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
      key: 'email',
      header: 'EMAIL',
      render: (item) => <span className={styles.emailCell}>{item.email}</span>,
    },
    {
      key: 'created_at',
      header: 'REGISTERED',
      width: '180px',
      render: (item) => formatDate(item.created_at),
    },
    {
      key: 'actions',
      header: 'ACTIONS',
      width: '160px',
      align: 'right',
      render: (item) => (
        <div className={styles.rowActions}>
          <button
            type="button"
            className={styles.actionBtn}
            onClick={() => handleOpenEdit(item)}
            title="Edit receiver details"
          >
            [edit]
          </button>
          <button
            type="button"
            className={`${styles.actionBtn} ${styles.deleteBtn}`}
            onClick={() => setDeleteTarget(item)}
            title="Delete receiver"
          >
            [del]
          </button>
        </div>
      ),
    },
  ]

  return (
    <div className={styles.receiversPage}>
      <AsciiCard
        title="Email Receivers Directory"
        subtitle={`(${filteredReceivers.length} listed)`}
        badge={<Badge variant="green">ACTIVE</Badge>}
        actions={
          <button
            type="button"
            className={styles.newBtn}
            onClick={handleOpenCreate}
          >
            [ + New Receiver ]
          </button>
        }
      >
        {/* Search filter and quick stats */}
        <div className={styles.toolbar}>
          <div className={styles.searchWrapper}>
            <span className={styles.promptArrow}>&gt;</span>
            <input
              type="text"
              placeholder="filter by name, email, or id..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className={styles.searchInput}
            />
          </div>
          <button
            type="button"
            onClick={() => fetchReceivers()}
            disabled={isLoading}
            className={styles.refreshBtn}
          >
            {isLoading ? '[ refreshing... ]' : '[ refresh ]'}
          </button>
        </div>

        {/* Monospace Table */}
        <TerminalTable
          columns={columns}
          data={filteredReceivers}
          isLoading={isLoading}
          keyExtractor={(item) => item.id}
          emptyMessage="No receivers registered yet. Click [+ New Receiver] to add one."
        />

        {/* Pagination controls */}
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

      {/* Create Receiver Modal */}
      <Modal
        isOpen={isCreateOpen}
        title="Create New Receiver"
        onClose={() => setIsCreateOpen(false)}
      >
        <form onSubmit={handleCreate} className={styles.modalForm}>
          <div className={styles.fieldGroup}>
            <label htmlFor="rec-name" className={styles.fieldLabel}>
              Full Name [name]:
            </label>
            <input
              id="rec-name"
              type="text"
              required
              maxLength={128}
              placeholder="e.g. Sarah Connor"
              value={name}
              onChange={(e) => setName(e.target.value)}
              autoFocus
            />
          </div>
          <div className={styles.fieldGroup}>
            <label htmlFor="rec-email" className={styles.fieldLabel}>
              Email Address [email]:
            </label>
            <input
              id="rec-email"
              type="email"
              required
              placeholder="e.g. sarah@cyberdyne.com"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
            />
          </div>
          <div className={styles.modalActions}>
            <button
              type="button"
              onClick={() => setIsCreateOpen(false)}
              className={styles.cancelBtn}
            >
              [ Cancel ]
            </button>
            <button
              type="submit"
              disabled={isSubmitting}
              className={styles.confirmBtn}
            >
              {isSubmitting ? '[ Saving... ]' : '[ > Save Receiver ]'}
            </button>
          </div>
        </form>
      </Modal>

      {/* Edit Receiver Modal */}
      <Modal
        isOpen={editTarget !== null}
        title={`Edit Receiver: ${editTarget?.name}`}
        onClose={() => setEditTarget(null)}
      >
        <form onSubmit={handleUpdate} className={styles.modalForm}>
          <div className={styles.fieldGroup}>
            <label htmlFor="edit-rec-name" className={styles.fieldLabel}>
              Full Name:
            </label>
            <input
              id="edit-rec-name"
              type="text"
              maxLength={128}
              value={name}
              onChange={(e) => setName(e.target.value)}
              autoFocus
            />
          </div>
          <div className={styles.fieldGroup}>
            <label htmlFor="edit-rec-email" className={styles.fieldLabel}>
              Email Address:
            </label>
            <input
              id="edit-rec-email"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
            />
          </div>
          <div className={styles.modalActions}>
            <button
              type="button"
              onClick={() => setEditTarget(null)}
              className={styles.cancelBtn}
            >
              [ Cancel ]
            </button>
            <button
              type="submit"
              disabled={isSubmitting}
              className={styles.confirmBtn}
            >
              {isSubmitting ? '[ Updating... ]' : '[ > Update Details ]'}
            </button>
          </div>
        </form>
      </Modal>

      {/* Delete Confirmation Modal */}
      <Modal
        isOpen={deleteTarget !== null}
        title="Confirm Deletion"
        onClose={() => setDeleteTarget(null)}
      >
        <div className={styles.deletePrompt}>
          <p>
            Are you sure you want to delete receiver{' '}
            <strong>{deleteTarget?.name}</strong> (<code>{deleteTarget?.email}</code>)?
          </p>
          <p className={styles.warningText}>
            [WARNING: This will remove this recipient from all associated groups.]
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
