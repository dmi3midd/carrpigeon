import React, { useState, useEffect } from 'react'
import styles from './Groups.module.css'
import { AsciiCard } from '../../components/AsciiCard/AsciiCard'
import { Badge } from '../../components/Badge/Badge'
import { TerminalTable, type Column } from '../../components/TerminalTable/TerminalTable'
import { Modal } from '../../components/Modal/Modal'
import { useGroupsStore } from '../../stores/useGroupsStore'
import { useReceiversStore } from '../../stores/useReceiversStore'
import { useUiStore } from '../../stores/useUiStore'
import { truncateId, formatDate } from '../../utils/formatters'
import type { Group } from '../../types/group'
import type { Receiver } from '../../types/receiver'

export const Groups: React.FC = () => {
  const {
    groups,
    isLoading,
    offset,
    limit,
    hasMore,
    selectedGroup,
    groupReceivers,
    isLoadingMembers,
    fetchGroups,
    createGroup,
    updateGroup,
    deleteGroup,
    setSelectedGroup,
    addReceiverToGroup,
    removeReceiverFromGroup,
    nextPage,
    prevPage,
  } = useGroupsStore()

  const { receivers, fetchReceivers } = useReceiversStore()
  const { addToast } = useUiStore()

  const [search, setSearch] = useState('')

  // Modals state
  const [isCreateOpen, setIsCreateOpen] = useState(false)
  const [editTarget, setEditTarget] = useState<Group | null>(null)
  const [deleteTarget, setDeleteTarget] = useState<Group | null>(null)

  // Form states
  const [name, setName] = useState('')
  const [description, setDescription] = useState('')
  const [selectedReceiverToAdd, setSelectedReceiverToAdd] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)

  useEffect(() => {
    fetchGroups(0)
    fetchReceivers(0)
  }, [fetchGroups, fetchReceivers])

  const handleOpenCreate = () => {
    setName('')
    setDescription('')
    setIsCreateOpen(true)
  }

  const handleOpenEdit = (g: Group) => {
    setEditTarget(g)
    setName(g.name)
    setDescription(g.description)
  }

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsSubmitting(true)
    try {
      const id = await createGroup({ name, description })
      addToast('success', `Created group '${name}' (ID: ${truncateId(id)})`)
      setIsCreateOpen(false)
    } catch (err) {
      addToast('error', err instanceof Error ? err.message : 'Failed to create group')
    } finally {
      setIsSubmitting(false)
    }
  }

  const handleUpdate = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!editTarget) return
    setIsSubmitting(true)
    try {
      await updateGroup(editTarget.id, { name, description })
      addToast('success', `Updated group '${name}'`)
      setEditTarget(null)
    } catch (err) {
      addToast('error', err instanceof Error ? err.message : 'Failed to update group')
    } finally {
      setIsSubmitting(false)
    }
  }

  const handleDelete = async () => {
    if (!deleteTarget) return
    setIsSubmitting(true)
    try {
      await deleteGroup(deleteTarget.id)
      addToast('success', `Deleted group '${deleteTarget.name}'`)
      setDeleteTarget(null)
    } catch (err) {
      addToast('error', err instanceof Error ? err.message : 'Failed to delete group')
    } finally {
      setIsSubmitting(false)
    }
  }

  const handleAddMember = async () => {
    if (!selectedGroup || !selectedReceiverToAdd) return
    setIsSubmitting(true)
    try {
      await addReceiverToGroup(selectedGroup.id, selectedReceiverToAdd)
      const rec = receivers.find((r) => r.id === selectedReceiverToAdd)
      addToast('success', `Added ${rec ? rec.name : 'receiver'} to ${selectedGroup.name}`)
      setSelectedReceiverToAdd('')
    } catch (err) {
      addToast('error', err instanceof Error ? err.message : 'Failed to add receiver to group')
    } finally {
      setIsSubmitting(false)
    }
  }

  const handleRemoveMember = async (receiver: Receiver) => {
    if (!selectedGroup) return
    try {
      await removeReceiverFromGroup(selectedGroup.id, receiver.id)
      addToast('info', `Removed ${receiver.name} from ${selectedGroup.name}`)
    } catch (err) {
      addToast('error', err instanceof Error ? err.message : 'Failed to remove member')
    }
  }

  const handleCopyId = (id: string) => {
    navigator.clipboard.writeText(id)
    addToast('info', `Copied Group ID: ${truncateId(id)}`)
  }

  const filteredGroups = groups.filter(
    (g) =>
      g.name.toLowerCase().includes(search.toLowerCase()) ||
      g.description.toLowerCase().includes(search.toLowerCase()) ||
      g.id.toLowerCase().includes(search.toLowerCase())
  )

  const groupColumns: Column<Group>[] = [
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
      key: 'description',
      header: 'DESCRIPTION',
      render: (item) => <span className={styles.descCell}>{item.description || '—'}</span>,
    },
    {
      key: 'receivers_count',
      header: 'MEMBERS',
      width: '110px',
      render: (item) => (
        <Badge variant={item.receivers_count > 0 ? 'cyan' : 'muted'}>
          {item.receivers_count} RECS
        </Badge>
      ),
    },
    {
      key: 'created_at',
      header: 'CREATED',
      width: '170px',
      render: (item) => formatDate(item.created_at),
    },
    {
      key: 'actions',
      header: 'ACTIONS',
      width: '210px',
      align: 'right',
      render: (item) => {
        const isSelected = selectedGroup?.id === item.id
        return (
          <div className={styles.rowActions}>
            <button
              type="button"
              className={`${styles.actionBtn} ${isSelected ? styles.selectedBtn : ''}`}
              onClick={() => setSelectedGroup(isSelected ? null : item)}
              title="View and manage members"
            >
              {isSelected ? '[close]' : '[members]'}
            </button>
            <button
              type="button"
              className={styles.actionBtn}
              onClick={() => handleOpenEdit(item)}
              title="Edit group metadata"
            >
              [edit]
            </button>
            <button
              type="button"
              className={`${styles.actionBtn} ${styles.deleteBtn}`}
              onClick={() => setDeleteTarget(item)}
              title="Delete group"
            >
              [del]
            </button>
          </div>
        )
      },
    },
  ]

  const memberColumns: Column<Receiver>[] = [
    {
      key: 'name',
      header: 'NAME',
      render: (item) => <strong>{item.name}</strong>,
    },
    {
      key: 'email',
      header: 'EMAIL',
      render: (item) => item.email,
    },
    {
      key: 'actions',
      header: 'ACTIONS',
      width: '90px',
      align: 'right',
      render: (item) => (
        <button
          type="button"
          className={`${styles.actionBtn} ${styles.deleteBtn}`}
          onClick={() => handleRemoveMember(item)}
          title="Remove receiver from this group"
        >
          [remove]
        </button>
      ),
    },
  ]

  return (
    <div className={styles.groupsPage}>
      {/* Groups Table Card */}
      <AsciiCard
        title="Receiver Groups / Distribution Lists"
        subtitle={`(${filteredGroups.length} groups)`}
        badge={<Badge variant="green">ACTIVE</Badge>}
        actions={
          <button
            type="button"
            className={styles.newBtn}
            onClick={handleOpenCreate}
          >
            [ + New Group ]
          </button>
        }
      >
        <div className={styles.toolbar}>
          <div className={styles.searchWrapper}>
            <span className={styles.promptArrow}>&gt;</span>
            <input
              type="text"
              placeholder="filter by group name, desc, or id..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className={styles.searchInput}
            />
          </div>
          <button
            type="button"
            onClick={() => fetchGroups()}
            disabled={isLoading}
            className={styles.refreshBtn}
          >
            {isLoading ? '[ refreshing... ]' : '[ refresh ]'}
          </button>
        </div>

        <TerminalTable
          columns={groupColumns}
          data={filteredGroups}
          isLoading={isLoading}
          keyExtractor={(item) => item.id}
          emptyMessage="No groups registered. Click [+ New Group] to create a list."
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

      {/* Selected Group Member Management Drawer */}
      {selectedGroup && (
        <AsciiCard
          title={`Group Members: ${selectedGroup.name}`}
          subtitle={`(ID: ${selectedGroup.id})`}
          badge={<Badge variant="cyan">{groupReceivers.length} ASSIGNED</Badge>}
          actions={
            <button
              type="button"
              className={styles.actionBtn}
              onClick={() => setSelectedGroup(null)}
            >
              [ Hide Members ]
            </button>
          }
        >
          {/* Add member toolbar */}
          <div className={styles.addMemberBar}>
            <select
              value={selectedReceiverToAdd}
              onChange={(e) => setSelectedReceiverToAdd(e.target.value)}
              className={styles.receiverSelect}
              aria-label="Select receiver to add"
            >
              <option value="">[ Select receiver to add to group... ]</option>
              {receivers
                .filter((r) => !groupReceivers.some((m) => m.id === r.id))
                .map((r) => (
                  <option key={r.id} value={r.id}>
                    {r.name} ({r.email})
                  </option>
                ))}
            </select>
            <button
              type="button"
              disabled={!selectedReceiverToAdd || isSubmitting}
              onClick={handleAddMember}
              className={styles.newBtn}
            >
              [ + Add Member ]
            </button>
          </div>

          <TerminalTable
            columns={memberColumns}
            data={groupReceivers}
            isLoading={isLoadingMembers}
            keyExtractor={(item) => item.id}
            emptyMessage="No members currently in this group. Use the selector above to add receivers."
          />
        </AsciiCard>
      )}

      {/* Create Group Modal */}
      <Modal
        isOpen={isCreateOpen}
        title="Create New Receiver Group"
        onClose={() => setIsCreateOpen(false)}
      >
        <form onSubmit={handleCreate} className={styles.modalForm}>
          <div className={styles.fieldGroup}>
            <label htmlFor="grp-name" className={styles.fieldLabel}>
              Group Name [name]:
            </label>
            <input
              id="grp-name"
              type="text"
              required
              maxLength={128}
              placeholder="e.g. Marketing Subscribers"
              value={name}
              onChange={(e) => setName(e.target.value)}
              autoFocus
            />
          </div>
          <div className={styles.fieldGroup}>
            <label htmlFor="grp-desc" className={styles.fieldLabel}>
              Description [description]:
            </label>
            <input
              id="grp-desc"
              type="text"
              required
              maxLength={256}
              placeholder="e.g. Newsletter and promotion recipients"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
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
              {isSubmitting ? '[ Saving... ]' : '[ > Save Group ]'}
            </button>
          </div>
        </form>
      </Modal>

      {/* Edit Group Modal */}
      <Modal
        isOpen={editTarget !== null}
        title={`Edit Group: ${editTarget?.name}`}
        onClose={() => setEditTarget(null)}
      >
        <form onSubmit={handleUpdate} className={styles.modalForm}>
          <div className={styles.fieldGroup}>
            <label htmlFor="edit-grp-name" className={styles.fieldLabel}>
              Group Name:
            </label>
            <input
              id="edit-grp-name"
              type="text"
              maxLength={128}
              value={name}
              onChange={(e) => setName(e.target.value)}
              autoFocus
            />
          </div>
          <div className={styles.fieldGroup}>
            <label htmlFor="edit-grp-desc" className={styles.fieldLabel}>
              Description:
            </label>
            <input
              id="edit-grp-desc"
              type="text"
              maxLength={256}
              value={description}
              onChange={(e) => setDescription(e.target.value)}
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
        title="Confirm Group Deletion"
        onClose={() => setDeleteTarget(null)}
      >
        <div className={styles.deletePrompt}>
          <p>
            Are you sure you want to delete group{' '}
            <strong>{deleteTarget?.name}</strong>?
          </p>
          <p className={styles.warningText}>
            [WARNING: This will delete the group record. Receivers inside this group will not be deleted.]
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
