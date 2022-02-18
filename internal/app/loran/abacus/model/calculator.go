package model

import "context"

// Calculator calculates new count.
type Calculator struct {
	UserID   int32
	EntityID int64
}

// CalculatorRepo represents the interface for working with calculator in db.
type CalculatorRepo interface {
	Add(ctx context.Context, userID int32, entityID int64) error
	Count(ctx context.Context, entityID int64) (int64, error)
}
