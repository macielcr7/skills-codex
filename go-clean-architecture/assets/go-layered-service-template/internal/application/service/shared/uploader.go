package shared

import (
	"context"
	"io"
)

// Uploader uploads objects to a remote storage.
type Uploader interface {
	PutObject(ctx context.Context, key string, body io.Reader) error
}
