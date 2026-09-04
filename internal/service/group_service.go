package service

import (
	"carrpigeo/internal/domain"
	"carrpigeo/internal/repository"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rs/xid"
)

var (
	ErrGroupAlreadyExists     = errors.New("group already exists")
	ErrGroupNotFound          = errors.New("group not found")
	ErrReceiverAlreadyInGroup = errors.New("receiver already in group")
)

type GroupService interface {
	// GetByID returns a group by its ID
	// Returns [ErrGroupNotFound] if no group is found
	GetByID(ctx context.Context, id string) (*domain.Group, error)
	// List returns a list of groups
	List(ctx context.Context, limit, offset int) ([]domain.Group, error)
	// Create creates a new group
	// Returns [ErrGroupAlreadyExists] if a group with the same name already exists
	Create(ctx context.Context, name, description string) (string, error)
	// Update updates an existing group
	// Returns [ErrGroupNotFound] if no group is found.
	Update(ctx context.Context, id, name, description string) (string, error)
	// Delete deletes a group by its ID
	Delete(ctx context.Context, id string) error

	// AddReceiver adds a receiver to a group.
	// Returns [ErrGroupNotFound] if no group is found.
	// Returns [ErrReceiverNotFound] if no receiver is found.
	// Returns [ErrReceiverAlreadyInGroup] if the receiver is already in the group.
	AddReceiver(ctx context.Context, groupID, receiverID string) error
	// RemoveReceiver removes a receiver from a group.
	// Returns [ErrGroupNotFound] if no group is found.
	RemoveReceiver(ctx context.Context, groupID, receiverID string) error
	// ListReceivers returns a list of receivers in a group.
	// Returns [ErrGroupNotFound] if no group is found
	ListReceivers(ctx context.Context, groupID string, limit, offset int) ([]domain.EmailReceiver, error)
}

type groupService struct {
	groupRepo    repository.GroupRepository
	receiverRepo repository.EmailReceiverRepository
}

func NewGroupService(groupRepo repository.GroupRepository, receiverRepo repository.EmailReceiverRepository) GroupService {
	return &groupService{
		groupRepo:    groupRepo,
		receiverRepo: receiverRepo,
	}
}

func (s *groupService) GetByID(ctx context.Context, id string) (*domain.Group, error) {
	op := "GroupService.GetByID"
	group, err := s.groupRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNoGroup) {
			return nil, fmt.Errorf("%s: %w", op, ErrGroupNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return group, nil
}

func (s *groupService) List(ctx context.Context, limit, offset int) ([]domain.Group, error) {
	op := "GroupService.List"
	groups, err := s.groupRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return groups, nil
}

func (s *groupService) Create(ctx context.Context, name, description string) (string, error) {
	op := "GroupService.Create"
	candidate, err := s.groupRepo.GetByName(ctx, name)
	if err != nil && !errors.Is(err, repository.ErrNoGroup) {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	if candidate != nil {
		return "", fmt.Errorf("%s: %w", op, ErrGroupAlreadyExists)
	}
	group := &domain.Group{
		ID:          xid.New().String(),
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := s.groupRepo.Create(ctx, group); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return group.ID, nil
}

func (s *groupService) Update(ctx context.Context, id, name, description string) (string, error) {
	op := "GroupService.Update"
	candidate, err := s.groupRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNoGroup) {
			return "", fmt.Errorf("%s: %w", op, ErrGroupNotFound)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}
	candidate.Name = name
	candidate.Description = description
	candidate.UpdatedAt = time.Now()
	if err := s.groupRepo.Update(ctx, candidate); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return candidate.ID, nil
}

func (s *groupService) Delete(ctx context.Context, id string) error {
	op := "GroupService.Delete"
	if err := s.groupRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *groupService) AddReceiver(ctx context.Context, groupID, receiverID string) error {
	op := "GroupService.AddReceiver"

	if _, err := s.groupRepo.GetByID(ctx, groupID); err != nil {
		if errors.Is(err, repository.ErrNoGroup) {
			return fmt.Errorf("%s: %w", op, ErrGroupNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	if _, err := s.receiverRepo.GetById(ctx, receiverID); err != nil {
		if errors.Is(err, repository.ErrNoReceiver) {
			return fmt.Errorf("%s: %w", op, ErrReceiverNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	exists, err := s.groupRepo.IsReceiverInGroup(ctx, groupID, receiverID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if exists {
		return fmt.Errorf("%s: %w", op, ErrReceiverAlreadyInGroup)
	}

	if err := s.groupRepo.AddReceiver(ctx, groupID, receiverID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *groupService) RemoveReceiver(ctx context.Context, groupID, receiverID string) error {
	op := "GroupService.RemoveReceiver"
	if err := s.groupRepo.RemoveReceiver(ctx, groupID, receiverID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *groupService) ListReceivers(ctx context.Context, groupID string, limit, offset int) ([]domain.EmailReceiver, error) {
	op := "GroupService.ListReceivers"

	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	if _, err := s.groupRepo.GetByID(ctx, groupID); err != nil {
		if errors.Is(err, repository.ErrNoGroup) {
			return nil, fmt.Errorf("%s: %w", op, ErrGroupNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	receivers, err := s.groupRepo.ListReceivers(ctx, groupID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return receivers, nil
}
