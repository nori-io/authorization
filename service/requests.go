package service

import (
	"github.com/asaskevich/govalidator"
	"github.com/cheebo/gorest"
)

// SignUp Request
type SignUpRequest struct {
	Email            string `json:"email" validate:"email"`
	PhoneCountryCode string `json:"phone_country_code" validate:"phone_country_code"`
	PhoneNumber      string `json:"phone_number" validate:"phone_number"`
	Password         string `json:"password" validate:"password"`
	Type             string `json:"user_type" validate:"user_type"`
	Meta             string `json:"meta" validate:"meta"`
}

func (r SignUpRequest) Validate() error {

	/*if (r.Email == "") && (r.Phone == "") {
		errorText = errorText + "Fields 'email' and 'phone' are unavailable on frontend \n"
		errCommon = errors.New(errorText)
	}

	if body.Password == "" {
		errorText = errorText + "Field 'password' is unavailable on frontend \n"
		errCommon = errors.New(errorText)

	}*/

	_, err := govalidator.ValidateStruct(r)
	return rest.ValidateResponse(err)
}

func (r SignUpRequest) ValidateOnlyByMail() error {
	_, err := govalidator.ValidateStruct(r)
	govalidator.IsEmail(r.Email)
	return rest.ValidateResponse(err)
}

func (r SignUpRequest) ValidateOnlyByPhone() error {
	_, err := govalidator.ValidateStruct(r)

	//re := regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)

	govalidator.IsEmail(r.PhoneCountryCode + r.PhoneNumber)
	return rest.ValidateResponse(err)
}

// SignIn Request
type SignInRequest struct {
	Name     string `json:"name" validate:"name"`
	Password string `json:"password" validate:"password"`
}

func (r SignInRequest) Validate() error {
	_, err := govalidator.ValidateStruct(r)
	return rest.ValidateResponse(err)
}

// SignOut Request
type SignOutRequest struct{}
