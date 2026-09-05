package usecase

import (
	"fmt"

	"github.com/IlucielI/insurance-policy-core-api/internal/domain"
	"github.com/IlucielI/insurance-policy-core-api/internal/repository"
)

type NotificationPreferencesUsecase struct {
	repo *repository.NotificationPreferencesRepository
}

func NewNotificationPreferencesUsecase(repo *repository.NotificationPreferencesRepository) *NotificationPreferencesUsecase {
	return &NotificationPreferencesUsecase{repo: repo}
}

func (u *NotificationPreferencesUsecase) GetPreferences(userID string) (*domain.NotificationPreferences, error) {
	prefs, err := u.repo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get preferences: %w", err)
	}
	return prefs, nil
}

func (u *NotificationPreferencesUsecase) UpdatePreferences(userID string, prefs *domain.NotificationPreferences) error {
	// Ensure user_id matches
	prefs.UserID = userID

	err := u.repo.Update(userID, prefs)
	if err != nil {
		return fmt.Errorf("failed to update preferences: %w", err)
	}

	return nil
}

func (u *NotificationPreferencesUsecase) ShouldSendEmail(userID, emailType string) (bool, error) {
	enabled, err := u.repo.CheckPreference(userID, emailType)
	if err != nil {
		// Default to true on error to avoid blocking emails
		return true, nil
	}
	return enabled, nil
}
