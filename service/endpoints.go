package service

import (
	"context"
	"fmt"
	"log"

	"github.com/nori-io/nori-common/endpoint"
)

func MakeSignUpEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, r interface{}) (interface{}, error) {
		req := r.(SignUpRequest)
		resp := s.SignUp(ctx, req)
		return *resp, resp.Error()
	}

}

func MakeSignInEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, r interface{}) (interface{}, error) {

		req := r.(SignInRequest)
		resp := s.SignIn(ctx, req)
		return *resp, resp.Error()
	}
}

func MakeSignOutEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, r interface{}) (interface{}, error) {
		req := r.(SignOutRequest)
		resp := s.SignOut(ctx, req)
		return *resp, resp.Err
	}
}

func MakeRecoveryCodesEndpoint(s Service) endpoint.Endpoint {
	log.Println("Encode")
	return func(ctx context.Context, r interface{}) (interface{}, error) {
		req := r.(RecoveryCodesRequest)
		resp := s.RecoveryCodes(ctx, req)
		fmt.Print(resp.Error())
		return *resp, resp.Error()
	}
}
