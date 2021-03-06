package database

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strconv"
)

var(
	ctx= context.Background()

)

type auth struct {
	db *sql.DB
}

func (a *auth) Create(model *AuthModel) error {

	tx, err := a.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.Fatal(err)
	}

	_, execErr:= tx.Exec ("INSERT INTO users (kind, status_id, type, created, updated, mfa_type) VALUES(?,?,?,?,?,?)",
		model.Kind_Users, model.StatusId_Users, model.Type_Users, model.Created_Users, model.Updated_Users, model.Mfa_type_Users)
	if execErr != nil {
		_ = tx.Rollback()
		log.Fatalf("Insert table 'users' error",execErr)
	}

    lastid,err:=a.db.Exec("select * from users where id = (select max(id) from users")
	if err != nil {
		log.Fatalf("Select table 'users' error ",err)
	}
	number,err:=lastid.LastInsertId()
	if err != nil {
		log.Fatalf("LastInsertId() taking error",err)
	}
	_, execErr = tx.Exec("INSERT INTO auth (user_id, phone, email, password, salt, created, updated, is_email_verified, is_phone_verified) VALUES("+(strconv.FormatInt(number,10))+",?,?,?,?,?,?,?,?)",
		model.UserId_Auth, model.Phone_Auth, model.Email_Auth, model.Password_Auth, model.Salt_Auth, model.Created_Auth, model.Updated_Auth, model.IsEmailVerified_Auth, model.IsPhoneVerified_Auth)
	if execErr != nil {
		_ = tx.Rollback()
		log.Fatalf("Insert table 'auth' error",execErr)
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("Commit transaction error",err)
	}

	return nil

}

func (a *auth) Update(model *AuthModel) error {
	if model.Id_Auth == 0 {
		return errors.New("Empty model")
	}
	_, err := a.db.Exec("UPDATE auth SET profile_user_id = ?, phone = ?, email = ?, password = ? salt = ? created =? WHERE id = ? ",
		model.UserId_Auth, model.Id_Auth)
	return err
}

func (a *auth) FindByEmail(email string) (model *AuthModel, err error) {
	rows, err := a.db.Query("SELECT user_id, phone, email, password, salt, created, updated, is_email_verified, is_phone_verified FROM auth WHERE email = ? LIMIT 1", email)
	if err != nil {
		return nil, err
	}
	model = &AuthModel{}

	defer rows.Close()
	for rows.Next() {
		var m AuthModel
		rows.Scan(&m.Id_Auth, &m.Email_Auth)
		model.Id_Auth = m.Id_Auth
		model.Email_Auth = m.Email_Auth
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return model, nil
}
