package squads

import "context"

type Repository interface {
	Create(ctx context.Context, squad Squad) (Squad, error)
	GetByKey(ctx context.Context, key string) (Squad, error)
}
