package shared

import "github.com/google/uuid"

// UUIDGenerator generates UUID-based IDs.
type UUIDGenerator struct{}

// NewID returns a new UUID string.
func (UUIDGenerator) NewID() string {
	return uuid.NewString()
}
