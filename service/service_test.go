package service_test

import (
	"context"
	"database/sql/driver"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	rest "github.com/cheebo/gorest"
	"github.com/nori-io/nori-common/meta"
	"github.com/nori-io/nori-common/mocks"
	"github.com/nori-io/nori-interfaces/interfaces"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/nori-io/authentication/service"
	"github.com/nori-io/authentication/service/database"
)

type AnyTime struct {
}

type AnyByteArray struct {
}

type AnyString struct {
}

func TestService_SignUp_Email_UserExists(t *testing.T) {
	mockRegistry := &mocks.Registry{}
	mockRegistry.On("Interface", meta.Interface("Auth")).Return(interfaces.AuthInterface)
	auth := *new(interfaces.Auth)
	mockRegistry.On("Interface", meta.Interface("Cache")).Return(interfaces.CacheInterface)
	cache := *new(interfaces.Cache)
	cfg := &service.Config{}
	mockDatabase, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	db := database.DB(mockDatabase, logrus.New())
	mockRegistry.On("Interface", meta.Interface("Mail")).Return(interfaces.MailInterface)
	mail := *new(interfaces.Mail)

	mockRegistry.On("Interface", meta.Interface("Session")).Return(interfaces.SessionInterface)
	session := *new(interfaces.Session)

	serviceTest := service.NewService(auth, cache, cfg, db, new(logrus.Logger), mail, session)
	signUpRequest := service.SignUpRequest{Email: "test@example.com", Password: "pass"}
	errField := rest.ErrFieldResp{
		Meta: rest.ErrMeta{
			ErrCode:    0,
			ErrMessage: "",
		},
	}
	errField.AddError("phone, email", 400, "User already exists.")

	respExpected := service.SignUpResponse{Email: "test@example.com", PhoneNumber: "", PhoneCountryCode: "", Err: errField}

	nonEmptyRows := sqlmock.NewRows([]string{"id", "email", "password", "salt"}).
		AddRow(1, "test@example.com", "pass", "salt")

	mock.ExpectQuery("SELECT id, email,password,salt FROM auth WHERE email = ? LIMIT 1").
		WithArgs("test@example.com").WillReturnRows(nonEmptyRows)
	pluginParamaters := service.PluginParameters{ActivationCode: true}

	resp := serviceTest.SignUp(context.Background(), signUpRequest, pluginParamaters)

	assert.Equal(t, &respExpected, resp)
}

/*
func TestService_SignUp_Phone_UserExists(t *testing.T) {
	auth := &mocks.Auth{}

	cache := &mocks.Cache{}
	cfg := &service.Config{}
	mockDatabase, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	db := database.DB(mockDatabase, logrus.New())
	mail := &mocks.Mail{}
	session := &mocks.Session{}

	serviceTest := service.NewService(auth, cache, cfg, db, new(logrus.Logger), mail, session)
	errField := rest.ErrFieldResp{
		Meta: rest.ErrFieldRespMeta{
			ErrCode:    0,
			ErrMessage: "",
		},
	}
	errField.AddError("phone, email", 400, "User already exists.")

	signUpRequest := service.SignUpRequest{PhoneCountryCode: "1", PhoneNumber: "234567890", Password: "pass"}

	//respExcepted.Err = errField

	respExpected := service.SignUpResponse{Email: "", PhoneCountryCode: "1", PhoneNumber: "234567890", Err: errField}

	nonEmptyRows := sqlmock.NewRows([]string{"id", "phone_country_code", "phone_number", "password", "salt"}).
		AddRow(1, "1", "234567890", "pass", "salt")

	mock.ExpectQuery("SELECT id, phone_country_code, phone_number, password,salt FROM auth WHERE concat(phone_country_code,phone_number)=?  LIMIT 1").WithArgs("1234567890").WillReturnRows(nonEmptyRows)

	mock.ExpectBegin()

	//mock.ExpectQuery("SELECT id, email,password,salt FROM auth WHERE email = ? LIMIT 1").WillReturnRows(nil)

	resp := serviceTest.SignUp(context.Background(), signUpRequest)

	assert.Equal(t, &respExpected, resp)
}

func TestService_SignUp_Phone_UserNotExist(t *testing.T) {
	auth := &mocks.Auth{}

	cache := &mocks.Cache{}
	cfg := &service.Config{}
	mockDatabase, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	db := database.DB(mockDatabase, logrus.New())
	mail := &mocks.Mail{}
	session := &mocks.Session{}

	serviceTest := service.NewService(auth, cache, cfg, db, new(logrus.Logger), mail, session)
	signUpRequest := service.SignUpRequest{PhoneCountryCode: "1", PhoneNumber: "234567890", Password: "pass"}

	respExpected := service.SignUpResponse{PhoneCountryCode: "1", PhoneNumber: "234567890", Err: nil}

	emptyRows := sqlmock.NewRows([]string{"id", "phone_country_code", "phone_number", "password", "salt"}).
		AddRow(nil, nil, nil, nil, nil)
	mock.ExpectQuery("SELECT id, phone_country_code, phone_number, password,salt FROM auth WHERE concat(phone_country_code,phone_number)=? LIMIT 1").
		WillReturnRows(emptyRows)

	mock.ExpectBegin()
	mock.ExpectPrepare("INSERT INTO users (status_account, type, created, updated) VALUES(?,?,?,?)").
		ExpectExec().WithArgs("locked", "", AnyTime{}, AnyTime{}).
		WillReturnResult(sqlmock.NewResult(1, 1))

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(1).
		RowError(1, fmt.Errorf("row error"))
	mock.ExpectQuery("SELECT LAST_INSERT_ID()").WillReturnRows(rows)

	mock.ExpectPrepare("INSERT INTO auth (user_id, phone_country_code, phone_number, password, salt, created, updated, is_email_verified, is_phone_verified) VALUES(?,?,?,?,?,?,?,?,?)").
		ExpectExec().
		WithArgs(1, "1", "234567890", AnyByteArray{}, AnyByteArray{}, AnyTime{}, AnyTime{}, false, false).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	resp := serviceTest.SignUp(context.Background(), signUpRequest)

	assert.Equal(t, &respExpected, resp)
}

func TestService_SignIn_Email_UserExist_CorrectPassword(t *testing.T) {

	auth := &mocks.Auth{}

	cache := &mocks.Cache{}
	cfg := &service.Config{}
	mockDatabase, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	db := database.DB(mockDatabase, logrus.New())
	mail := &mocks.Mail{}
	session := &mocks.Session{}

	serviceTest := service.NewService(auth, cache, cfg, db, new(logrus.Logger), mail, session)
	signInRequest := service.SignInRequest{Name: "test@example.com", Password: "pass"}

	respExpected := service.SignInResponse{Id: 1, User: service.UserResponse{UserName: "test@example.com"}, HttpStatusCode: 0, Token: mock2.Anything}

	salt, err := database.CreateSalt()
	if err != nil {
		t.Log(err)
	}

	password, err := database.Hash([]byte("pass"), salt)

	nonEmptyRows := sqlmock.NewRows([]string{"id", "email", "password", "salt"}).
		AddRow(1, "test@example.com", password, salt)

	mock.ExpectQuery("SELECT id, email,password,salt FROM auth WHERE email = ? LIMIT 1").
		WithArgs("test@example.com").WillReturnRows(nonEmptyRows)

	mock.ExpectExec("INSERT INTO authentication_history (user_id, signin, meta) VALUES(?,?,?)").
		WithArgs(1, AnyTime{}, "").WillReturnResult(sqlmock.NewResult(1, 1))
	auth.On("AccessToken", mock2.Anything).Return(mock2.Anything, nil)
	session.On("Save", mock2.Anything, mock2.Anything, mock2.Anything).Return(nil)

	pluginParamaters := service.PluginParameters{UserRegistrationByEmailAddress: true}
	resp := serviceTest.SignIn(context.Background(), signInRequest, pluginParamaters)

	assert.Equal(t, &respExpected, resp)

}

func TestService_SignIn_Email_UserExist_IncorrectPassword(t *testing.T) {
	auth := &mocks.Auth{}

	cache := &mocks.Cache{}
	cfg := &service.Config{}
	mockDatabase, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	db := database.DB(mockDatabase, logrus.New())
	mail := &mocks.Mail{}
	session := &mocks.Session{}

	serviceTest := service.NewService(auth, cache, cfg, db, new(logrus.Logger), mail, session)
	signInRequest := service.SignInRequest{Name: "test@example.com", Password: "pass"}
	Err := rest.ErrResp{Meta: rest.ErrMeta{ErrMessage: "Incorrect Password", ErrCode: 0}}

	respExpected := service.SignInResponse{Id: 1, User: service.UserResponse{UserName: "test@example.com"}, HttpStatusCode: 0, Err: Err}

	salt, err := database.CreateSalt()
	if err != nil {
		t.Log(err)
	}

	password, err := database.Hash([]byte("pass1"), salt)

	nonEmptyRows := sqlmock.NewRows([]string{"id", "email", "password", "salt"}).
		AddRow(1, "test@example.com", password, salt)

	mock.ExpectQuery("SELECT id, email,password,salt FROM auth WHERE email = ? LIMIT 1").
		WithArgs("test@example.com").WillReturnRows(nonEmptyRows)

	pluginParamaters := service.PluginParameters{UserRegistrationByEmailAddress: true}
	resp := serviceTest.SignIn(context.Background(), signInRequest, pluginParamaters)

	assert.Equal(t, &respExpected, resp)

}

func TestService_SignIn_Email_UserExist_IncorrectUserName(t *testing.T) {
	auth := &mocks.Auth{}

	cache := &mocks.Cache{}
	cfg := &service.Config{}
	mockDatabase, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	db := database.DB(mockDatabase, logrus.New())
	mail := &mocks.Mail{}
	session := &mocks.Session{}

	serviceTest := service.NewService(auth, cache, cfg, db, new(logrus.Logger), mail, session)
	signInRequest := service.SignInRequest{Name: "testNot@example.com", Password: "pass"}

	Err := rest.ErrResp{Meta: rest.ErrMeta{ErrMessage: "User not found", ErrCode: 0}}

	respExpected := service.SignInResponse{Id: 0, User: service.UserResponse{UserName: "testNot@example.com"}, HttpStatusCode: 0, Err: Err}

	emptyRows := sqlmock.NewRows([]string{"id", "email", "password", "salt"}).
		AddRow(nil, nil, nil, nil)

	mock.ExpectQuery("SELECT id, email,password,salt FROM auth WHERE email = ? LIMIT 1").WillReturnRows(emptyRows)

	pluginParamaters := service.PluginParameters{UserRegistrationByEmailAddress: true}
	resp := serviceTest.SignIn(context.Background(), signInRequest, pluginParamaters)

	assert.Equal(t, &respExpected, resp)
}

func TestService_SignUp_Email_UserNotExist(t *testing.T) {
	auth := &mocks.Auth{}

	cache := &mocks.Cache{}
	cfg := &service.Config{}
	mockDatabase, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	db := database.DB(mockDatabase, logrus.New())
	mail := &mocks.Mail{}
	session := &mocks.Session{}

	serviceTest := service.NewService(auth, cache, cfg, db, new(logrus.Logger), mail, session)
	signUpRequest := service.SignUpRequest{Email: "test@example.com", Password: "pass"}

	respExpected := service.SignUpResponse{Email: "test@example.com", PhoneNumber: "", PhoneCountryCode: "", Err: nil}

	emptyRows := sqlmock.NewRows([]string{"id", "email", "password", "salt"}).
		AddRow(nil, nil, nil, nil)

	mock.ExpectQuery("SELECT id, email,password,salt FROM auth WHERE email = ? LIMIT 1").WillReturnRows(emptyRows)

	mock.ExpectBegin()
	mock.ExpectPrepare("INSERT INTO users (status_account, type, created, updated) VALUES(?,?,?,?)").
		ExpectExec().WithArgs("locked", "", AnyTime{}, AnyTime{}).
		WillReturnResult(sqlmock.NewResult(1, 1))

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(1).
		RowError(1, fmt.Errorf("row error"))
	mock.ExpectQuery("SELECT LAST_INSERT_ID()").WillReturnRows(rows)

	mock.ExpectPrepare("INSERT INTO auth (user_id, email, password, salt, created, updated, is_email_verified, is_phone_verified) VALUES(?,?,?,?,?,?,?,?)").
		ExpectExec().
		WithArgs(1, "test@example.com", AnyByteArray{}, AnyByteArray{}, AnyTime{}, AnyTime{}, false, false).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()
	pluginParamaters := service.PluginParameters{ActivationCode: true}

	resp := serviceTest.SignUp(context.Background(), signUpRequest, pluginParamaters)

	assert.Equal(t, &respExpected, resp)
}

func TestService_SignIn_Phone_UserExist_CorrectPassword(t *testing.T) {
	auth := &mocks.Auth{}

	cache := &mocks.Cache{}
	cfg := &service.Config{}
	mockDatabase, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	db := database.DB(mockDatabase, logrus.New())
	mail := &mocks.Mail{}
	session := &mocks.Session{}

	serviceTest := service.NewService(auth, cache, cfg, db, new(logrus.Logger), mail, session)
	signInRequest := service.SignInRequest{Name: "1234567890", Password: "pass"}

	respExpected := service.SignInResponse{Id: 1, User: service.UserResponse{UserName: "1234567890"}, HttpStatusCode: 0, Token: mock2.Anything}

	salt, err := database.CreateSalt()
	if err != nil {
		t.Log(err)
	}

	password, err := database.Hash([]byte("pass"), salt)

	nonEmptyRowsPhone := sqlmock.NewRows([]string{"id", "phone_country_code", "phone_number", "password", "salt"}).
		AddRow(1, "1", "234567890", password, salt)

	mock.ExpectQuery("SELECT id, phone_country_code, phone_number, password,salt FROM auth WHERE concat(phone_country_code,phone_number)=?  LIMIT 1").
		WithArgs("1234567890").WillReturnRows(nonEmptyRowsPhone)

	mock.ExpectExec("INSERT INTO authentication_history (user_id, signin, meta) VALUES(?,?,?)").
		WithArgs(1, AnyTime{}, "").WillReturnResult(sqlmock.NewResult(1, 1))
	auth.On("AccessToken", mock2.Anything).Return(mock2.Anything, nil)
	session.On("Save", mock2.Anything, mock2.Anything, mock2.Anything).Return(nil)

	pluginParamaters := service.PluginParameters{UserRegistrationByPhoneNumber: true}
	resp := serviceTest.SignIn(context.Background(), signInRequest, pluginParamaters)

	assert.Equal(t, &respExpected, resp)
}

func TestService_SignIn_Phone_UserExist_IncorrectPassword(t *testing.T) {
	auth := &mocks.Auth{}

	cache := &mocks.Cache{}
	cfg := &service.Config{}
	mockDatabase, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	db := database.DB(mockDatabase, logrus.New())
	mail := &mocks.Mail{}
	session := &mocks.Session{}

	serviceTest := service.NewService(auth, cache, cfg, db, new(logrus.Logger), mail, session)
	signInRequest := service.SignInRequest{Name: "1234567890", Password: "pass"}
	Err := rest.ErrResp{Meta: rest.ErrMeta{ErrMessage: "Incorrect Password", ErrCode: 0}}

	respExpected := service.SignInResponse{Id: 1, User: service.UserResponse{UserName: "1234567890"}, HttpStatusCode: 0, Err: Err}

	salt, err := database.CreateSalt()
	if err != nil {
		t.Log(err)
	}

	password, err := database.Hash([]byte("pass1"), salt)

	nonEmptyRows := sqlmock.NewRows([]string{"id", "phone_country_code", "phone_number", "password", "salt"}).
		AddRow(1, "1", "234567890", password, salt)

	mock.ExpectQuery("SELECT id, phone_country_code, phone_number, password,salt FROM auth WHERE concat(phone_country_code,phone_number)=?  LIMIT 1").
		WithArgs("1234567890").WillReturnRows(nonEmptyRows)

	pluginParamaters := service.PluginParameters{UserRegistrationByPhoneNumber: true}
	resp := serviceTest.SignIn(context.Background(), signInRequest, pluginParamaters)

	assert.Equal(t, &respExpected, resp)

}

func TestService_SignIn_Phone_UserExist_IncorrectUserName(t *testing.T) {
	auth := &mocks.Auth{}

	cache := &mocks.Cache{}
	cfg := &service.Config{}
	mockDatabase, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	db := database.DB(mockDatabase, logrus.New())
	mail := &mocks.Mail{}
	session := &mocks.Session{}

	serviceTest := service.NewService(auth, cache, cfg, db, new(logrus.Logger), mail, session)
	signInRequest := service.SignInRequest{Name: "1234567890", Password: "pass"}

	Err := rest.ErrResp{Meta: rest.ErrMeta{ErrMessage: "User not found", ErrCode: 0}}

	respExpected := service.SignInResponse{Id: 0, User: service.UserResponse{UserName: "1234567890"}, HttpStatusCode: 0, Err: Err}

	emptyRows := sqlmock.NewRows([]string{"id", "phone_country_code", "phone_number", "password", "salt"}).
		AddRow(nil, nil, nil, nil, nil)

	mock.ExpectQuery("SELECT id, phone_country_code, phone_number, password,salt FROM auth WHERE concat(phone_country_code,phone_number)=?  LIMIT 1").WillReturnRows(emptyRows)
	auth.On("AccessToken", mock2.Anything).Return(mock2.Anything, nil)
	session.On("Save", mock2.Anything, mock2.Anything, mock2.Anything).Return(nil)
	pluginParamaters := service.PluginParameters{UserRegistrationByPhoneNumber: true}
	resp := serviceTest.SignIn(context.Background(), signInRequest, pluginParamaters)

	assert.Equal(t, &respExpected, resp)

}

func TestService_SignOut(t *testing.T) {

	auth := &mocks.Auth{}

	cache := &mocks.Cache{}
	cfg := &service.Config{}
	mockDatabase, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	db := database.DB(mockDatabase, logrus.New())
	mail := &mocks.Mail{}
	session := &mocks.Session{}

	serviceTest := service.NewService(auth, cache, cfg, db, new(logrus.Logger), mail, session)
	signOutRequest := service.SignOutRequest{Name: "1234567890"}

	respExpected := &service.SignOutResponse{HttpStatusCode: 0, Err: nil}

	type mapClaims jwt.MapClaims

	type any interface{}

	contextTest := make(jwt.MapClaims)
	contextTest["exp"] = time.Now().Hour() + time.Now().Minute() + time.Now().Second()
	contextTest["iat"] = time.Now().Hour() + time.Now().Minute() + time.Now().Second()
	contextTest["iss"] = "nori/api"
	contextTest["nbf"] = time.Now().Hour() + time.Now().Minute() + time.Now().Second()
	contextTest["raw"] = map[string]interface{}{
		"id":   "",
		"name": "test@example.com",
	}
	contextTest["sub"] = "nori"

	ctxAuthData := context.WithValue(context.Background(), "nori.auth.data", contextTest)
	ctx := context.WithValue(ctxAuthData, "nori.session.id", "irf7VYww6w57KzlVELHp6DvzCNiLjgqU")

	nonEmptyRows := sqlmock.NewRows([]string{"id", "email", "password", "salt"}).
		AddRow(1, "test@example.com", "pass", "salt")

	mock.ExpectQuery("SELECT id, email,password,salt FROM auth WHERE email = ? LIMIT 1").
		WithArgs("test@example.com").WillReturnRows(nonEmptyRows)

	mock.ExpectExec("UPDATE authentication_history SET  signout = ?   WHERE user_id = ? ORDER BY id DESC LIMIT 1").
		WithArgs(AnyTime{}, 1).WillReturnResult(sqlmock.NewResult(1, 0))

	session.On("Delete", []byte("irf7VYww6w57KzlVELHp6DvzCNiLjgqU")).Return(nil)

	session.On("SessionId", ctx).Return([]byte("irf7VYww6w57KzlVELHp6DvzCNiLjgqU"))

	resp := serviceTest.SignOut(ctx, signOutRequest)

	assert.Equal(t, respExpected, resp)

}
*/
/*func TestService_ActivationCode(t *testing.T) {

	auth := &mocks.Auth{}

	cache := &mocks.Cache{}
	cfg := &service.Config{}
	mockDatabase, _, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	db := database.DB(mockDatabase, logrus.New())
	mail := &mocks.Mail{}
	session := &mocks.Session{}

	//mail.On("Send").

	//serviceTest := service.NewService(auth, cache, cfg, db, new(logrus.Logger), mail, session)

	//SignUpRequest:=service.SignUpRequest{Email:"test@example.com", Password:"pass", Type:"vendor"}

	//serviceTest.SignUp(context.Background(), SignUpRequest)

	code := rand.RandomAlphaNum(65)
	salt, codeForSend, err := getActivationCode([]byte(code))
	if err == nil {
		//activationRequest:=service.ActivationCodeRequest{ActivationCode:codeForSend}
		activate([]byte(code), salt, codeForSend)
		mail.Send(activate([]byte(code), salt, codeForSend))
	}

}

func getActivationCode(codeMessage []byte) ([]byte, []byte, error) {
	salt := rand.RandomAlphaNum(65)

	codeForSend, err := database.Hash(codeMessage, []byte(salt))

	return []byte(salt), codeForSend, err
}

func activate(code []byte, salt []byte, codeForSend []byte) string {
	result, _ := database.VerifyPassword([]byte(code), salt, codeForSend)
	b := strconv.FormatBool(result)

	return b
}*/

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func (s AnyByteArray) Match(v driver.Value) bool {
	_, ok := v.([]byte)
	return ok
}

func (s AnyString) Match(v driver.Value) bool {
	_, ok := v.(string)
	return ok
}
