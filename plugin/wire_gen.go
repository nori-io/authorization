// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package main

import (
	"github.com/google/wire"
	"github.com/nori-io/authentication/internal/handler/http"
	"github.com/nori-io/authentication/internal/repository/user"
	"github.com/nori-io/authentication/internal/service/auth"
	"github.com/nori-io/common/v4/pkg/domain/registry"
	"github.com/nori-io/interfaces/database/orm/gorm"
	"github.com/nori-io/interfaces/nori/http"
	"github.com/nori-io/interfaces/nori/session"
)

// Injectors from wire.go:

func Initialize(registry2 registry.Registry, urlPrefix string) (*http_handler.Handler, error) {
	httpHttp, err := http.GetHttp(registry2)
	if err != nil {
		return nil, err
	}
	sessionSession, err := session.GetSession(registry2)
	if err != nil {
		return nil, err
	}
	db, err := pg.GetGorm(registry2)
	if err != nil {
		return nil, err
	}
	userRepository := user.New(db)
	authenticationService := auth.New(sessionSession, userRepository)
	handler := &http_handler.Handler{
		R:         httpHttp,
		Auth:      authenticationService,
		UrlPrefix: urlPrefix,
	}
	return handler, nil
}

// wire.go:

var set1 = wire.NewSet(auth.New, pg.GetGorm, session.GetSession, user.New, wire.Struct(new(http_handler.Handler), "R", "Auth", "UrlPrefix"), http.GetHttp)
