package pkg

import (
	"context"
	"errors"

	"github.com/UniqueStudio/open-platform/database"
)

type AccessUser struct {
	UserID  string
	IsAdmin bool
}

func UserToAccessUser(user *database.User) *AccessUser {
	return &AccessUser{
		UserID:  user.UID,
		IsAdmin: user.Role >= ROLE_BRANCH_SECRETARY,
	}
}

func AddAccessUserIntoContext(ctx context.Context, au *AccessUser) context.Context {
	return context.WithValue(ctx, AccessUser{}, au)
}

func GetAccessUserFromContext(ctx context.Context) (*AccessUser, error) {
	value := ctx.Value(AccessUser{})
	au, ok := value.(*AccessUser)
	if !ok || au == nil {
		return nil, errors.New("can't get AccessUser from context")
	}
	return au, nil
}
