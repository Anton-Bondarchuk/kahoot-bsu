package ports

import (
	"context"
	"kahoot_bsu/internal/domain/models"
)

type UserRepository interface {	
	Update(
		ctx context.Context, 
		userID int64, 
		updateFn func(innerCtx context.Context, user *models.User) error,
	) error
	
	UpdateOrCreate(ctx context.Context, user *models.User) error
}
