package systemauthrestrictions

import (
	"context"
	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"
	"github.com/kyma-incubator/compass/components/director/pkg/resource"
	"github.com/pkg/errors"
)

//go:generate mockery -name=UIDService -output=automock -outpkg=automock -case=underscore
type UIDService interface {
	Generate() string
}

//go:generate mockery -name=SystemAuthsRepository -output=automock -outpkg=automock -case=underscore
type SystemAuthsRepository interface {
	ListForObjectGlobal(ctx context.Context, objectType model.SystemAuthReferenceObjectType, objectID string) ([]model.SystemAuth, error)
}

//go:generate mockery -name=SystemAuthRestrictionsRepository -output=automock -outpkg=automock -case=underscore
type SystemAuthRestrictionsRepository interface {
	Create(ctx context.Context, item model.SystemAuthRestrictions) error
}

type service struct {
	uidService                 UIDService
	systemAuthsRepo            SystemAuthsRepository
	systemAuthRestrictionsRepo SystemAuthRestrictionsRepository
}

func NewService(uidService UIDService, systemAuthsRepo SystemAuthsRepository, systemAuthRestrictionsRepo SystemAuthRestrictionsRepository) *service {
	return &service{
		uidService:                 uidService,
		systemAuthsRepo:            systemAuthsRepo,
		systemAuthRestrictionsRepo: systemAuthRestrictionsRepo,
	}
}

func (s *service) CreateMany(ctx context.Context, in model.SystemAuthRestrictionsInput) error {
	var systemAuths []model.SystemAuth

	//TODO can be done with single query for listing and single query for inserting
	if in.To.ApplicationID != nil {
		systemAuthsForApp, err := s.systemAuthsRepo.ListForObjectGlobal(ctx, model.ApplicationReference, *in.To.ApplicationID)
		if err != nil {
			return err
		}
		systemAuths = append(systemAuths, systemAuthsForApp...)
	}

	if in.To.RuntimeID != nil {
		systemAuthsForRuntime, err := s.systemAuthsRepo.ListForObjectGlobal(ctx, model.RuntimeReference, *in.To.RuntimeID)
		if err != nil {
			return err
		}

		systemAuths = append(systemAuths, systemAuthsForRuntime...)
	}

	if in.To.IntegrationSystemID != nil {
		systemAuthsForIntegrationSystem, err := s.systemAuthsRepo.ListForObjectGlobal(ctx, model.IntegrationSystemReference, *in.To.IntegrationSystemID)
		if err != nil {
			return err
		}

		systemAuths = append(systemAuths, systemAuthsForIntegrationSystem...)
	}

	if len(systemAuths) == 0 {
		return apperrors.NewNotFoundErrorWithType(resource.SystemAuth)
	}

	for _, auth := range systemAuths {
		systemAuthRestrictions := in.ToSystemAuthRestrictions(s.uidService.Generate(), auth.ID)
		if err := s.systemAuthRestrictionsRepo.Create(ctx, systemAuthRestrictions); err != nil {
			return errors.Wrapf(err, "while creating System Auth restrictions for system auth with ID %s", systemAuthRestrictions.SystemAuthID)
		}
	}

	return nil
}

func (s *service) DeleteMany(ctx context.Context, in model.SystemAuthRestrictionsInput) error {
	//TODO implement
	return nil
}
