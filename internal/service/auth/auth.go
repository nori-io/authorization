package auth

import (
	"context"
	"crypto/rand"
	"time"

	"github.com/nori-plugins/authentication/pkg/enum/hash_algorithm"
	"github.com/nori-plugins/authentication/pkg/enum/users_action"
	"github.com/nori-plugins/authentication/pkg/enum/users_status"
	"github.com/nori-plugins/authentication/pkg/enum/users_type"

	"github.com/nori-plugins/authentication/pkg/enum/session_status"

	s "github.com/nori-io/interfaces/nori/session"

	service "github.com/nori-plugins/authentication/internal/domain/service"

	"github.com/nori-plugins/authentication/internal/domain/entity"
)

func (srv AuthenticationService) SignUp(ctx context.Context, data service.SignUpData) (*entity.User, error) {
	if err := data.Validate(); err != nil {
		return nil, err
	}

	var user *entity.User

	//@todo определить на уровне конфига верификацию email и на основании этого заполнить информацию о статусе
	//@todo заполнить оставшиеся поля
	user = &entity.User{
		Status:          users_status.Active,
		UserType:        users_type.User,
		MfaType:         0,
		Email:           data.Email,
		Password:        data.Password,
		HashAlgorithm:   hash_algorithm.Bcrypt,
		IsEmailVerified: false,
		IsPhoneVerified: false,
		CreatedAt:       time.Now(),
	}

	if err := srv.UserRepository.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (srv *AuthenticationService) SignIn(ctx context.Context, data service.SignInData) (*entity.Session, error) {
	var err error
	if err = data.Validate(); err != nil {
		return nil, err
	}

	var user *entity.User
	user, err = srv.UserRepository.FindByEmail(ctx, data.Email)
	if err != nil {
		return nil, err
	}

	sid, err := srv.getToken()
	if err != nil {
		return nil, err
	}

	if err := srv.SessionRepository.Create(ctx, &entity.Session{
		ID:         0,
		UserID:     user.ID,
		SessionKey: sid,
		Status:     session_status.Active,
		OpenedAt:   time.Now(),
	}); err != nil {
		return nil, err
	}

	session, err := srv.SessionRepository.FindBySessionKey(ctx, string(sid))
	if err != nil {
		return nil, err
	}

	if err = srv.AuthenticationLogRepository.Create(ctx, &entity.AuthenticationLog{
		ID:     0,
		UserID: user.ID,
		Action: users_action.SignUp,
		//@todo заполнить метаданные айпи адресом и городом или чем-то ещё?
		Meta:      "",
		SessionID: session.ID,
		CreatedAt: time.Now(),
	}); err != nil {
		return nil, err
	}

	return &entity.Session{
		SessionKey: sid,
	}, nil
}

func (srv *AuthenticationService) SignOut(ctx context.Context, data *entity.Session) error {
	if err := srv.Session.Delete(data.SessionKey); err != nil {
		return err
	}

	session, err := srv.SessionRepository.FindBySessionKey(ctx, string(data.SessionKey))
	if err != nil {
		return err
	}

	if err := srv.SessionRepository.Update(ctx, &entity.Session{
		ID:        session.ID,
		UserID:    session.UserID,
		Status:    session_status.Inactive,
		ClosedAt:  time.Now(),
		UpdatedAt: time.Now(),
	}); err != nil {
		return err
	}

	// if err:=srv.AuthenticationLogRepository.Create()
	// if err:=srv.AuthenticationLogRepository.Update(ctx, )

	return err
}

func (srv *AuthenticationService) getToken() ([]byte, error) {
	sid := make([]byte, 32)

	if _, err := rand.Read(sid); err != nil {
		return nil, err
	}
	if err := srv.Session.Get(sid, s.SessionActive); err != nil {
		srv.Session.Save(sid, s.SessionActive, 0)
		return sid, nil
	}
	return srv.getToken()
}
