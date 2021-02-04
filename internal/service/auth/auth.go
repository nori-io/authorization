package auth

import (
	"context"
	"crypto/rand"

	"github.com/nori-plugins/authentication/internal/domain/entity"

	"github.com/nori-plugins/authentication/internal/domain/repository"

	s "github.com/nori-io/interfaces/nori/session"
	serv "github.com/nori-plugins/authentication/internal/domain/service"
)

type service struct {
	session s.Session
	db      repository.UserRepository
}

func New(sessionInstance s.Session, dbInstance repository.UserRepository) serv.AuthenticationService {
	return &service{
		session: sessionInstance,
		db:      dbInstance,
	}
}

func (srv *service) SignUp(ctx context.Context, data serv.SignUpData) (*entity.User, error) {
	if err := data.Validate(); err != nil {
		return nil, err
	}

	var user *entity.User

	user = &entity.User{
		Email:    data.Email,
		Password: data.Password,
	}

	if err := srv.db.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (srv *service) SignIn(ctx context.Context, data serv.SignInData) (*entity.Session, error) {
	if err := data.Validate(); err != nil {
		return nil, err
	}

	user := &entity.User{
		Email:    data.Email,
		Password: data.Password,
	}

	var err error
	user, err = srv.db.GetByEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	}

	sid, err := srv.getToken()
	if err != nil {
		return nil, err
	}
	return &entity.Session{SessionKey: sid}, nil
}

func (srv *service) SignOut(ctx context.Context, data *entity.Session) error {
	err := srv.session.Delete([]byte(data.SessionKey))
	return err
}

func (srv *service) MfaRecoveryCodes(ctx context.Context, data *entity.Session) error {
	err := srv.db.
	return err
}

func (srv *service) getToken() ([]byte, error) {
	sid := make([]byte, 32)

	if _, err := rand.Read(sid); err != nil {
		return nil, err
	}
	if err := srv.session.Get(sid, s.SessionActive); err != nil {
		srv.session.Save(sid, s.SessionActive, 0)
		return sid, nil
	}
	return srv.getToken()
}
