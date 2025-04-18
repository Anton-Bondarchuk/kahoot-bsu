package models

import "context"

type User struct {
	ID    int64  `json:"id"`
	Login string `json:"login"`
	Role  int64  `json:"role"`
}

type UserRepository interface {	
	Update(
		ctx context.Context, 
		userID int64, 
		updateFn func(innerCtx context.Context, user *User) error,
	) error
	UpdateOrCreate(
		ctx context.Context, 
		user *User,
	) error
}