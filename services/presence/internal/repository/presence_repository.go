package repository

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	domain "github.com/adrputra/face-recognition-svc/presence-svc/internal/domain"
)

type PresenceRepository struct {
	mu    sync.RWMutex
	items map[string]*domain.Entity
}

func NewPresenceRepository() *PresenceRepository {
	return &PresenceRepository{
		items: make(map[string]*domain.Entity),
	}
}

func (r *PresenceRepository) Create(_ context.Context, item *domain.Entity) (*domain.Entity, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	cp := clone(item)
	r.items[cp.ID] = cp
	return clone(cp), nil
}

func (r *PresenceRepository) Get(_ context.Context, id string) (*domain.Entity, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, ok := r.items[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return clone(item), nil
}

func (r *PresenceRepository) List(_ context.Context, filter domain.Filter) ([]*domain.Entity, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*domain.Entity, 0, len(r.items))
	for _, item := range r.items {
		if filter.UserID != "" && item.UserID != filter.UserID {
			continue
		}
		if filter.InstitutionID != "" && item.InstitutionID != filter.InstitutionID {
			continue
		}
		if filter.Status != "" && !strings.EqualFold(item.Status, filter.Status) {
			continue
		}
		result = append(result, clone(item))
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.After(result[j].CreatedAt)
	})

	return result, nil
}

func (r *PresenceRepository) Update(_ context.Context, id string, patch domain.Patch) (*domain.Entity, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	item, ok := r.items[id]
	if !ok {
		return nil, domain.ErrNotFound
	}

	if patch.UserID != nil {
		item.UserID = *patch.UserID
	}
	if patch.InstitutionID != nil {
		item.InstitutionID = *patch.InstitutionID
	}
	if patch.Status != nil {
		item.Status = *patch.Status
	}
	if patch.Note != nil {
		item.Note = *patch.Note
	}
	item.UpdatedAt = time.Now().UTC()

	return clone(item), nil
}

func (r *PresenceRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.items[id]; !ok {
		return domain.ErrNotFound
	}
	delete(r.items, id)
	return nil
}

func clone(item *domain.Entity) *domain.Entity {
	if item == nil {
		return nil
	}
	cp := *item
	return &cp
}
