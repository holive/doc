package docApi

import "context"

type Repository interface {
	Create(ctx context.Context, docApi *DocApi) (*DocApi, error)
	Find(ctx context.Context, squad string, projeto string, versao string) (*DocApi, error)
	Delete(ctx context.Context, squad string, projeto string, versao string) error
	FindAll(ctx context.Context, limit string, offset string) (*SearchResult, error)
}
