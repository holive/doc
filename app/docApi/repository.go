package docApi

import "context"

type Repository interface {
	Create(ctx context.Context, docApi *DocApi) (*DocApi, error)
	Delete(ctx context.Context, squad string, versao string, projeto string) error
	FindAll(ctx context.Context, limit string, offset string) (*SearchResult, error)
}
