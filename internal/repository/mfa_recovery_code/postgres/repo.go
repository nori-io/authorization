package postgres

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/nori-plugins/authentication/internal/domain/entity"
)

type MfaRecoveryCodeRepository struct {
	Db *gorm.DB
}

func (m MfaRecoveryCodeRepository) Use(ctx context.Context, e *entity.MfaRecoveryCode) error {
	panic("implement me")
}

func (m MfaRecoveryCodeRepository) Create(ctx context.Context, userID uint64, mfaRecoveryCode string) (entity.MfaRecoveryCode, error) {
	panic("implement me")
}