package application

import (
	"context"
	"time"

	"github.com/ramdhanrizkij/next-store-api/internal/modules/identity/domain"
	appErrors "github.com/ramdhanrizkij/next-store-api/internal/shared/errors"
	"github.com/ramdhanrizkij/next-store-api/internal/shared/pagination"
)

type UserService interface {
	GetByID(ctx context.Context, id string) (*UserResponse, error)
	Update(ctx context.Context, id string, req *UpdateUserRequest) (*UserResponse, error)
	List(ctx context.Context, p pagination.Pagination) ([]*UserResponse, *pagination.Meta, error)
	Delete(ctx context.Context, id string) error
}

type userService struct {
	userRepo domain.UserRepository
}

func NewUserService(userRepo domain.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetByID(ctx context.Context, id string) (*UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, appErrors.ErrNotFound
	}
	return ToUserResponse(user), nil
}

func (s *userService) Update(ctx context.Context, id string, req *UpdateUserRequest) (*UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, appErrors.ErrNotFound
	}

	user.Name = req.Name
	if req.Phone != nil {
		user.Phone = req.Phone
	}
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return ToUserResponse(user), nil
}

func (s *userService) List(ctx context.Context, p pagination.Pagination) ([]*UserResponse, *pagination.Meta, error) {
	users, total, err := s.userRepo.FindAll(ctx, p.Limit, p.Offset())
	if err != nil {
		return nil, nil, err
	}

	meta := pagination.BuildMeta(total, p.Page, p.Limit)
	return ToUserResponses(users), meta, nil
}

func (s *userService) Delete(ctx context.Context, id string) error {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return appErrors.ErrNotFound
	}

	return s.userRepo.Delete(ctx, id)
}
