package cranmer

// AddRequest represents add request body.
// UserID represents unique user that
// we want to increment count based on EntityID.
type AddRequest struct {
	UserID   int32 `json:"user_id"`
	EntityID int64 `json:"entity_id"`
}
