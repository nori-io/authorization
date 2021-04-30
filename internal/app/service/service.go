package service

import (
	"github.com/google/wire"
	"github.com/nori-plugins/authentication/internal/service"
	"github.com/nori-plugins/authentication/internal/service/auth"
	"github.com/nori-plugins/authentication/internal/service/mfa_recovery_code"
	"github.com/nori-plugins/authentication/internal/service/mfa_totp"
	"github.com/nori-plugins/authentication/internal/service/reset_password"
	"github.com/nori-plugins/authentication/internal/service/session"
	"github.com/nori-plugins/authentication/internal/service/settings"
	"github.com/nori-plugins/authentication/internal/service/social_provider"
	"github.com/nori-plugins/authentication/internal/service/user"
	"github.com/nori-plugins/authentication/internal/service/user_log"
)

var ServiceSet = wire.NewSet(
	wire.Struct(new(auth.Params),
		"UserLogService",
		"MfaRecoveryCodeService",
		"MfaTotpService",
		"SessionService",
		"SocialProviderService",
		"UserService",
		"SecurityHelper",
		"Config",
		"Transactor"),
	auth.New,
	wire.Struct(new(user_log.Params),
		"UserLogRepository",
		"Transactor"),
	user_log.New,
	wire.Struct(new(mfa_recovery_code.Params),
		"UserLogService",
		"SessionService",
		"UserService",
		"MfaRecoveryCodeRepository",
		"MfaRecoveryCodeHelper",
		"Config",
		"Transactor"),
	mfa_recovery_code.New,
	wire.Struct(new(mfa_totp.Params),
		"UserLogService",
		"SessionService",
		"UserService",
		"MfaTotpRepository",
		"MfaTotpHelper",
		"Config",
		"Transactor"),
	mfa_totp.New,
	wire.Struct(new(reset_password.Params),
		"UserLogService",
		"UserService",
		"ResetPasswordRepository",
		"SecurityHelper",
		"Config",
		"Transactor"),
	reset_password.New,
	wire.Struct(new(session.Params), "SessionRepository", "Transactor"),
	session.New,
	wire.Struct(new(settings.Params),
		"UserLogService",
		"UserService",
		"SessionRepository",
		"SecurityHelper",
		"Config",
		"Transactor"),
	settings.New,
	wire.Struct(new(social_provider.Params), "SocialProviderRepository"),
	social_provider.New,
	wire.Struct(new(user.Params),
		"UserLogService",
		"UserRepository",
		"SecurityHelper",
		"Transactor",
		"Config"),
	user.New,
	wire.Struct(new(service.Params),
		"AuthenticationService",
		"UserLogService",
		"MfaRecoveryCodeService",
		"MfaTotpService",
		"ResetPasswordService",
		"SessionService",
		"SettingsService",
		"SocialProviderService",
		"UserService"),
	service.New,
)
