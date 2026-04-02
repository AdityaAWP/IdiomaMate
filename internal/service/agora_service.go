package service

import (
	"github.com/AdityaAWP/IdiomaMate/internal/domain"
	rtctokenbuilder "github.com/AgoraIO-Community/go-tokenbuilder/rtctokenbuilder2"
	"github.com/google/uuid"
)

type agoraService struct {
	appID          string
	appCertificate string
}

func NewAgoraService(appID, appCertificate string) domain.AgoraService {
	return &agoraService{
		appID:          appID,
		appCertificate: appCertificate,
	}
}

func (s *agoraService) GenerateRTCToken(channelName string, userID uuid.UUID) (string, error) {
	if s.appID == "" || s.appCertificate == "" {
		// Log or return empty if not properly configured yet.
		// For development, we might not have Agora set up but shouldn't block the app.
		return "", nil
	}

	tokenExpireTimeInSeconds := uint32(3600) // 1 hour
	privilegeExpireTimeInSeconds := uint32(3600)

	token, err := rtctokenbuilder.BuildTokenWithUserAccount(
		s.appID,
		s.appCertificate,
		channelName,
		userID.String(),
		rtctokenbuilder.RolePublisher,
		tokenExpireTimeInSeconds,
		privilegeExpireTimeInSeconds,
	)

	return token, err
}
