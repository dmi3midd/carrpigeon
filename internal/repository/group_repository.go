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
	List(ctx context.Context, limit, offset int64) ([]domain.Group, error)
	// Create creates a new group
	Create(ctx context.Context, group *domain.Group) error
	// Update updates an existing group
	Update(ctx context.Context, group *domain.Group) error
	// Delete deletes a group by its ID
	Delete(ctx context.Context, id string) error

	// AddReceiver adds a receiver to a group
	AddReceiver(ctx context.Context, groupID, receiverID string) error
	// RemoveReceiver removes a receiver from a group
	RemoveReceiver(ctx context.Context, groupID, receiverID string) error
	// ListReceivers returns a list of receivers in a group
	ListReceivers(ctx context.Context, groupID string, limit, offset int64) ([]domain.EmailReceiver, error)
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
		SELECT id, name, description, created_at, updated_at
		FROM groups
		WHERE id = $1
	`
	var group domain.Group
	err := r.DB.GetContext(ctx, &group, query, id)
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
		SELECT id, name, description, created_at, updated_at
		FROM groups
		WHERE name = $1
	`
	var group domain.Group
	err := r.DB.GetContext(ctx, &group, query, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrNoGroup)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &group, nil
}

func (r *groupRepository) List(ctx context.Context, limit, offset int64) ([]domain.Group, error) {
	op := "GroupRepository.List"
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM groups
		LIMIT $1 OFFSET $2
	`
	var groups []domain.Group
	err := r.DB.SelectContext(ctx, &groups, query, limit, offset)
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
	_, err := r.DB.NamedExecContext(ctx, query, group)
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
	_, err := r.DB.NamedExecContext(ctx, query, group)
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
	_, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *groupRepository) AddReceiver(ctx context.Context, groupID, receiverID string) error {
	op := "GroupRepository.AddReceiver"
	query := `
		INSERT INTO groups_receivers (group_id, receiver_id)
		VALUES ($1, $2)
	`
	_, err := r.DB.ExecContext(ctx, query, groupID, receiverID)
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
	_, err := r.DB.ExecContext(ctx, query, groupID, receiverID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *groupRepository) ListReceivers(ctx context.Context, groupID string, limit, offset int64) ([]domain.EmailReceiver, error) {
	op := "GroupRepository.ListReceivers"
	query := `
		SELECT r.id, r.name, r.email, r.created_at, r.updated_at
		FROM email_receivers r
		JOIN groups_receivers gr ON r.id = gr.receiver_id
		WHERE gr.group_id = $1
		LIMIT $2 OFFSET $3
	`
	var receivers []domain.EmailReceiver
	err := r.DB.SelectContext(ctx, &receivers, query, groupID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return receivers, nil
}
