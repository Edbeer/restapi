package utils

import (
	"context"

	"github.com/Edbeer/restapi/pkg/httpe"
	"github.com/Edbeer/restapi/pkg/logger"
)

// Validate is user from owner of content
func ValidateIsOwner(ctx context.Context, creatorID string, logger logger.Logger) error {
	user, err := GetUserFromCtx(ctx)
	if err != nil {
		return err
	}

	if user.ID.String() != creatorID {
		logger.Errorf(
			"ValidateIsOwner, userID: %v, creatorID: %v",
			user.ID.String(),
			creatorID,
		)
		return httpe.Forbidden
	}
	return nil
}