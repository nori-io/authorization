package errors

import "github.com/nori-plugins/authentication/pkg/errors"

var (
	ActiveSessionAlreadyExists = errors.New("user_session_exists", "users already sign in, sign up isn't possible", errors.ErrAlreadyExists)
	UserNotFound               = errors.New("authentication.user_not_found", "user not found", errors.ErrNotFound)
	TokenNotFound              = errors.New("authentication.token_not_found", "token not found", errors.ErrNotFound)
	SessionNotFound            = errors.New("session.session_not_found", "session not found", errors.ErrNotFound)
	MfaRecoveryCodeNotFound    = errors.New("mfa.recovery_code_not_found", "mfa recovery code not found", errors.ErrNotFound)
	MfaTotpSecretNotFound      = errors.New("mfa.totp_secret_not_found", "mfa totp secret not found", errors.ErrNotFound)
	EmailAlreadyTaken          = errors.New("authentication.email_already_taken", "email already taken", errors.ErrAlreadyExists)
	SocialProviderNotFound     = errors.New("social_provider_not_found", "social_provider_not_found", errors.ErrNotFound)
)
