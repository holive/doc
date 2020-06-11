package docApi

import "context"

type Repository interface {
	Create(ctx context.Context, docApi *DocApi) error
	Find(ctx context.Context, squad string, projeto string, versao string) (*DocApi, error)
	FindBySquad(ctx context.Context, squad string, limit string, offset string) (*SearchResult, error)
	Delete(ctx context.Context, squad string, projeto string, versao string) error
	FindAll(ctx context.Context, limit string, offset string) (*SearchResult, error)
}
