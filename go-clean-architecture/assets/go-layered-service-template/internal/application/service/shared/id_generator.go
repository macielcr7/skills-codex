package shared

// IDGenerator generates unique IDs for domain entities.
type IDGenerator interface {
	NewID() string
}
