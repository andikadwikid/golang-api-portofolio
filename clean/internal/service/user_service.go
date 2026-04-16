package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"portofolio-api/clean/internal/domain"
	"portofolio-api/clean/internal/pkg/utils"
	"portofolio-api/clean/internal/repository"
)

type UserService interface {
	Register(ctx context.Context, input domain.CreateUserInput) (*domain.User, error)
	Login(ctx context.Context, input domain.UserLoginInput) (string, error)
	GetAllUsers(ctx context.Context) ([]domain.UserResponse, error)
	Update(ctx context.Context, id string, input domain.UpdateUserInput) error
	Delete(ctx context.Context, id string) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) Register(ctx context.Context, input domain.CreateUserInput) (*domain.User, error) {
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))

	existingUser, err := s.userRepo.FindByEmail(ctx, input.Email)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, errors.New("gagal memeriksa email")
	}
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := &domain.User{
		Name:     input.Name,
		Username: input.Username,
		Email:    input.Email,
		Password: hashedPassword,
		Avatar:   "",
		Bio:      "",
	}

	insertedID, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, errors.New("failed to create user")
	}

	user.ID = insertedID
	user.CreatedAt = time.Now()
	return user, nil
}

func (s *userService) Login(ctx context.Context, input domain.UserLoginInput) (string, error) {
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))

	user, err := s.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", errors.New("user not found")
		}
		return "", errors.New("failed to find user")
	}

	if !utils.CheckPasswordHash(input.Password, user.Password) {
		return "", errors.New("invalid email or password")
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}

func (s *userService) GetAllUsers(ctx context.Context) ([]domain.UserResponse, error) {
	users, err := s.userRepo.GetAllUsers(ctx)
	if err != nil {
		return nil, errors.New("failed to fetch users")
	}

	usersResponse := make([]domain.UserResponse, 0, len(users))
	for _, u := range users {
		usersResponse = append(usersResponse, domain.UserResponse{
			ID:        u.ID.Hex(),
			Name:      u.Name,
			Username:  u.Username,
			Email:     u.Email,
			Avatar:    u.Avatar,
			Bio:       u.Bio,
			CreatedAt: u.CreatedAt,
		})
	}
	return usersResponse, nil
}

func (s *userService) Update(ctx context.Context, id string, input domain.UpdateUserInput) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id")
	}

	updateData := bson.M{
		"updated_at": time.Now(),
	}

	if input.Name != nil {
		updateData["name"] = *input.Name
	}
	if input.Username != nil {
		updateData["username"] = *input.Username
	}
	if input.Email != nil {
		updateData["email"] = strings.ToLower(strings.TrimSpace(*input.Email))
	}
	if input.Password != nil {
		hashedPassword, err := utils.HashPassword(*input.Password)
		if err != nil {
			return errors.New("failed to hash password")
		}
		updateData["password"] = hashedPassword
	}
	if input.Avatar != nil {
		updateData["avatar"] = *input.Avatar
	}
	if input.Bio != nil {
		updateData["bio"] = *input.Bio
	}

	if len(updateData) == 1 { // Only updated_at
		return errors.New("no fields to update")
	}

	result, err := s.userRepo.Update(ctx, objID, updateData)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (s *userService) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid id")
	}

	result, err := s.userRepo.Delete(ctx, objID)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("user not found")
	}

	return nil
}
