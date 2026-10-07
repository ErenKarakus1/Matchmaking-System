package validation

import "testing"

func TestValidateUsernameAcceptsAllowedUsernames(t *testing.T) {
	usernames := []string{
		"abc",
		"player_1",
		"player123",
		"abc_123_xyz",
	}

	for _, username := range usernames {
		t.Run(username, func(t *testing.T) {
			if err := ValidateUsername(username); err != nil {
				t.Fatalf("ValidateUsername(%q) returned error: %v", username, err)
			}
		})
	}
}

func TestValidateUsernameRejectsInvalidUsernames(t *testing.T) {
	usernames := []string{
		"",
		"ab",
		"has-dash",
		"has space",
		"UPPERCASE",
		"with.dot",
		"with@symbol",
		"this_username_is_way_too_long_for_the_limit",
	}

	for _, username := range usernames {
		t.Run(username, func(t *testing.T) {
			if err := ValidateUsername(username); err == nil {
				t.Fatalf("ValidateUsername(%q) returned nil error, want error", username)
			}
		})
	}
}
