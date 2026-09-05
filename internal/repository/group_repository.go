package repository

import (
	"carrpigeo/internal/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

var (
	ErrNoGroup = errors.New("group not found")
)

type GroupRepository interface {
	// GetByID returns a group by its ID
	// Returns ErrNoGroup if no group is found
	GetByID(ctx context.Context, id string) (*domain.Group, error)
	// GetByName returns a group by its name
	// Returns ErrNoGroup if no group is found
	GetByName(ctx context.Context, name string) (*domain.Group, error)
	// List returns a list of groups
	List(ctx context.Context, limit, offset int) ([]domain.Group, error)
	// Create creates a new group
	Create(ctx context.Context, group *domain.Group) error
	// Update updates an existing group
	Update(ctx context.Context, group *domain.Group) error
	// Delete deletes a group by its ID
	Delete(ctx context.Context, id string) error

	// IsReceiverInGroup checks if a receiver is in a group
	IsReceiverInGroup(ctx context.Context, groupID, receiverID string) (bool, error)
	// AddReceiver adds a receiver to a group
	AddReceiver(ctx context.Context, groupID, receiverID string) error
	// RemoveReceiver removes a receiver from a group
	RemoveReceiver(ctx context.Context, groupID, receiverID string) error
	// ListReceivers returns a list of receivers in a group
	ListReceivers(ctx context.Context, groupID string, limit, offset int) ([]domain.Receiver, error)
}

type groupRepository struct {
	DB *sqlx.DB
}

func NewGroupRepository(db *sqlx.DB) GroupRepository {
	return &groupRepository{
		DB: db,
	}
}

func (r *groupRepository) GetByID(ctx context.Context, id string) (*domain.Group, error) {
	op := "GroupRepository.GetByID"
	query := `
		SELECT g.id, g.name, g.description, 
		       COALESCE(COUNT(gr.receiver_id), 0) AS receivers_count, 
		       g.created_at, g.updated_at
		FROM groups g
		LEFT JOIN groups_receivers gr ON g.id = gr.group_id
		WHERE g.id = $1
		GROUP BY g.id
	`
	executor := ExtractTx(ctx, r.DB)
	var group domain.Group
	err := sqlx.GetContext(ctx, executor, &group, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrNoGroup)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &group, nil
}

func (r *groupRepository) GetByName(ctx context.Context, name string) (*domain.Group, error) {
	op := "GroupRepository.GetByName"
	query := `
		SELECT g.id, g.name, g.description, 
		       COALESCE(COUNT(gr.receiver_id), 0) AS receivers_count, 
		       g.created_at, g.updated_at
		FROM groups g
		LEFT JOIN groups_receivers gr ON g.id = gr.group_id
		WHERE g.name = $1
		GROUP BY g.id
	`
	executor := ExtractTx(ctx, r.DB)
	var group domain.Group
	err := sqlx.GetContext(ctx, executor, &group, query, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrNoGroup)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &group, nil
}

func (r *groupRepository) List(ctx context.Context, limit, offset int) ([]domain.Group, error) {
	op := "GroupRepository.List"
	query := `
		SELECT g.id, g.name, g.description, 
		       COALESCE(COUNT(gr.receiver_id), 0) AS receivers_count, 
		       g.created_at, g.updated_at
		FROM groups g
		LEFT JOIN groups_receivers gr ON g.id = gr.group_id
		GROUP BY g.id
		ORDER BY g.created_at DESC
		LIMIT $1 OFFSET $2
	`
	executor := ExtractTx(ctx, r.DB)
	var groups []domain.Group
	err := sqlx.SelectContext(ctx, executor, &groups, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return groups, nil
}

func (r *groupRepository) Create(ctx context.Context, group *domain.Group) error {
	op := "GroupRepository.Create"
	query := `
		INSERT INTO groups (id, name, description, created_at, updated_at)
		VALUES (:id, :name, :description, :created_at, :updated_at)
	`
	executor := ExtractTx(ctx, r.DB)
	_, err := sqlx.NamedExecContext(ctx, executor, query, group)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *groupRepository) Update(ctx context.Context, group *domain.Group) error {
	op := "GroupRepository.Update"
	query := `
		UPDATE groups
		SET name = :name, description = :description, updated_at = :updated_at
		WHERE id = :id
	`
	executor := ExtractTx(ctx, r.DB)
	_, err := sqlx.NamedExecContext(ctx, executor, query, group)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *groupRepository) Delete(ctx context.Context, id string) error {
	op := "GroupRepository.Delete"
	query := `
		DELETE FROM groups
		WHERE id = $1
	`
	executor := ExtractTx(ctx, r.DB)
	_, err := executor.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *groupRepository) IsReceiverInGroup(ctx context.Context, groupID, receiverID string) (bool, error) {
	op := "GroupRepository.IsReceiverInGroup"
	query := `
		SELECT EXISTS(
			SELECT 1
			FROM groups_receivers
			WHERE group_id = $1 AND receiver_id = $2
		)
	`
	executor := ExtractTx(ctx, r.DB)
	var exists bool
	err := sqlx.GetContext(ctx, executor, &exists, query, groupID, receiverID)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return exists, nil
}

func (r *groupRepository) AddReceiver(ctx context.Context, groupID, receiverID string) error {
	op := "GroupRepository.AddReceiver"
	query := `
		INSERT INTO groups_receivers (group_id, receiver_id)
		VALUES ($1, $2)
	`
	executor := ExtractTx(ctx, r.DB)
	_, err := executor.ExecContext(ctx, query, groupID, receiverID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *groupRepository) RemoveReceiver(ctx context.Context, groupID, receiverID string) error {
	op := "GroupRepository.RemoveReceiver"
	query := `
		DELETE FROM groups_receivers
		WHERE group_id = $1 AND receiver_id = $2
	`
	executor := ExtractTx(ctx, r.DB)
	_, err := executor.ExecContext(ctx, query, groupID, receiverID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *groupRepository) ListReceivers(ctx context.Context, groupID string, limit, offset int) ([]domain.Receiver, error) {
	op := "GroupRepository.ListReceivers"
	query := `
		SELECT r.id, r.name, r.email, r.created_at, r.updated_at
		FROM receivers r
		JOIN groups_receivers gr ON r.id = gr.receiver_id
		WHERE gr.group_id = $1
		LIMIT $2 OFFSET $3
	`
	executor := ExtractTx(ctx, r.DB)
	var receivers []domain.Receiver
	err := sqlx.SelectContext(ctx, executor, &receivers, query, groupID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return receivers, nil
}
