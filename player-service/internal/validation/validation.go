package validation

import (
	"errors"
	"regexp"

	"github.com/ErenKarakus1/Matchmaking-System/player-service/internal/models"
)

var usernameRegex = regexp.MustCompile(`^[a-z0-9_]+$`)

func ValidateCreatePlayerRequest(req models.CreatePlayerRequest) error {
	if req.Username == "" {
		return errors.New("username is required")
	}
	if len(req.Username) < 3 {
		return errors.New("username must be at least 3 characters")
	}
	if len(req.Username) > 32 {
		return errors.New("username must be at most 32 characters")
	}
	if !usernameRegex.MatchString(req.Username) {
		return errors.New("username can only contain lowercase letters, numbers, and underscores")
	}
	return nil
}
