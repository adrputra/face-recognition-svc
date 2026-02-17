package domain

import "context"

type Repository interface {
	Create(ctx context.Context, item *Entity) (*Entity, error)
	Get(ctx context.Context, id string) (*Entity, error)
	List(ctx context.Context, filter Filter) ([]*Entity, error)
	Update(ctx context.Context, id string, patch Patch) (*Entity, error)
	Delete(ctx context.Context, id string) error
}
