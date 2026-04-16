package service

import (
	"context"
	"errors"
	"portofolio-api/clean/internal/domain"
	"portofolio-api/clean/internal/repository"
	"strings"
)

type SocialMediaService interface {
	Create(ctx context.Context, input domain.SocialMediaCreateInput) (*domain.SocialMedia, error)
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
