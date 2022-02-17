package cranmer

// AddRequest represents add request body.
// UserID represents unique user that
// we want to increment count based on EntityID.
type AddRequest struct {
	UserID   int32
	EntityID int64
}
