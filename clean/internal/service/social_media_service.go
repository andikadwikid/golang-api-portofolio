package service

import (
	"context"
	"errors"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"portofolio-api/clean/internal/domain"
	"portofolio-api/clean/internal/repository"
)

type SocialMediaService interface {
	Create(ctx context.Context, input domain.SocialMediaCreateInput) (*domain.SocialMedia, error)
	GetAll(ctx context.Context) ([]domain.SocialMedia, error)
	GetById(ctx context.Context, id primitive.ObjectID) (domain.SocialMedia, error)
	Update(ctx context.Context, id primitive.ObjectID, input domain.SocialMediaUpdateInput) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type socialMediaService struct {
	socialMediaRepo repository.SocialMediaRepository
}

func NewSocialMediaService(socialMediaRepo repository.SocialMediaRepository) SocialMediaService {
	return &socialMediaService{socialMediaRepo: socialMediaRepo}
}

func (s *socialMediaService) Create(ctx context.Context, input domain.SocialMediaCreateInput) (*domain.SocialMedia, error) {
	input.Name = strings.ToLower(strings.TrimSpace(input.Name))
	existingSocialMedia, err := s.socialMediaRepo.CheckExistingSocialMedia(ctx, input.Name)
	if err != nil {
		return nil, err
	}
	if existingSocialMedia {
		return nil, errors.New("social media already exists")
	}
	socialMedia := &domain.SocialMedia{
		Icon:      input.Icon,
		Name:      input.Name,
		IsDeleted: false,
	}

	insertedID, err := s.socialMediaRepo.Create(ctx, socialMedia)
	if err != nil {
		return nil, errors.New("failed to create social media")
	}

	socialMedia.ID = insertedID
	return socialMedia, nil
}

func (s *socialMediaService) GetAll(ctx context.Context) ([]domain.SocialMedia, error) {
	return s.socialMediaRepo.GetAllSocialMedia(ctx)
}

func (s *socialMediaService) GetById(ctx context.Context, id primitive.ObjectID) (domain.SocialMedia, error) {
	return s.socialMediaRepo.GetSocialMediaById(ctx, id)
}

func (s *socialMediaService) Update(ctx context.Context, id primitive.ObjectID, input domain.SocialMediaUpdateInput) error {
	socialMedia := &domain.SocialMedia{
		Icon:      input.Icon,
		Name:      input.Name,
		IsDeleted: false,
	}
	return s.socialMediaRepo.UpdateSocialMedia(ctx, id, socialMedia)
}

func (s *socialMediaService) Delete(ctx context.Context, id primitive.ObjectID) error {
	return s.socialMediaRepo.DeleteSocialMedia(ctx, id)
}
