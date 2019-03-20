package database_test

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"golang.org/x/net/context"

	"github.com/nori-io/auth/service/database/sql_scripts"
)


type (
  AnyTime struct{}
)

func TestUsers_Create(t *testing.T) {

    mockDatabase, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	/*defer mockDatabase.Close()
	createTables(mockDatabase)

	testDatabase:=database.DB(mockDatabase,nil)

	modelUsers:=&database.UsersModel{
		Type:    "vendor",
		Created:time.Now(),
		Updated:time.Now(),

	}
	modelAuth:= &database.AuthModel{
		Email:    "test@mail.ru",
		Password: "pass",
		Created:time.Now(),
		Updated:time.Now(),

	}


	if err1:= testDatabase.Users().Create(modelAuth,modelUsers); err1 != nil {
		t.Errorf("error was not expected while updating stats: %s", err.Error())
	}*/

  rows := sqlmock.NewRows([]string{"id", "status_account", "type", "created","updated", "mfa_type"}).
		AddRow(0, "active", "vendor", time.Now(),time.Now(),"")

	mock.ExpectQuery("SELECT").WillReturnRows(rows)


	rs,err2:= mockDatabase.Query("SELECT")
	fmt.Println(rs)
   if err2!=nil {
		fmt.Println("failed to match expected query",err2)
		return
	}
	defer rs.Close()



	for rs.Next() {
		var id int
		var status_account string
		rs.Scan(&id, &status_account)
		fmt.Println("scanned id:", id, "and status_account:", status_account)
	}
	if rs.Err() != nil {
		fmt.Println("got rows error:", rs.Err())
	}
}
func TestUsers_Create2(t *testing.T) {
	var err error
	t.Parallel()


	mockDatabase, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}


	defer mockDatabase.Close()

	mock.ExpectExec("INSERT INTO users (status_account, type, created, updated) VALUES(?,?,?,?)").
		WithArgs("active","vendor",AnyTime{}, AnyTime{}).WillReturnResult(sqlmock.NewResult(0,1))

  _, err = mockDatabase.Exec("INSERT INTO users (status_account, type, created, updated) VALUES(?,?,?,?)",
  	"active", "vendor",time.Now(),time.Now())
	if err != nil {
		t.Errorf("error '%s' was not expected, while inserting a row", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}




}

func TestUsers_Create3(t *testing.T) {
	var err error
	t.Parallel()


	mockDatabase, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}


	defer mockDatabase.Close()
	mock.ExpectExec("INSERT INTO users (status_account, type, created, updated) VALUES(?,?,?,?)").
		WithArgs("active","vendor",AnyTime{}, AnyTime{}).WillReturnResult(sqlmock.NewResult(0,1))
	mock.ExpectExec("INSERT INTO auth (user_id,  email, password, salt, created, updated, is_email_verified, is_phone_verified) VALUES(?,?,?,?,?,?,?,?)").
		WithArgs(0,"1@mail.ru","pass","",AnyTime{},AnyTime{},0,0)


	result, _:= mockDatabase.Exec("INSERT INTO users (status_account, type, created, updated) VALUES(?,?,?,?)",
		"active", "vendor",time.Now(),time.Now())

	lastId:=int64(-1)
	lastId,_= result.LastInsertId()

	mockDatabase.Exec("INSERT INTO auth (user_id,  email, password, salt, created, updated, is_email_verified, is_phone_verified) VALUES(?,?,?,?,?,?,?,?)",
		lastId,"1@mail.ru","pass","",time.Now(),time.Now(),0,0)

	t.Log(lastId)
	if err != nil {
		t.Errorf("error '%s' was not expected, while inserting a row", err)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	fmt.Println(err)




}


func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}
/*func TestUsers_Create2(t *testing.T) {
	var err error


	mockDatabase, mock, err := sqlmock.New()
	if !assert.Nilf(t,
		err,
		"Error on trying to start up a stub database connection") {
		t.Fatal()
	}
	defer mockDatabase.Close()

	testDatabase:=database.DB(mockDatabase,nil)



	modelUsers:=&database.UsersModel{
		Type:    "vendor",
		Created:time.Now(),
		Updated:time.Now(),

	}
	modelAuth:= &database.AuthModel{
		Email:    "test@mail.ru",
		Password: "pass",
		Created:time.Now(),
		Updated:time.Now(),

	}

	t.Run("WantCommitError", func(t *testing.T) {
		// define sql behavior
		mock.ExpectBegin()

		testDatabase.Users().Create(modelAuth,modelUsers)

		defer func() {
			mock.ExpectCommit().
				WillReturnError(nil)


		}()

		mock.ExpectationsWereMet()

	})


	if err = testDatabase.Users().Create(modelAuth,modelUsers); err != nil {
		t.Errorf("error was not expected while updating stats: %s", err.Error())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}




}
*/

func createTables(db *sql.DB) error {
	ctx := context.Background()


	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})

	if err != nil {
		log.Fatal(err)
	}

	_, execErr := tx.Exec(
		sql_scripts.SetDatabaseSettings)
	if execErr != nil {
		_ = tx.Rollback()
		log.Fatal(execErr)

	}

	_, execErr= tx.Exec(
		sql_scripts.SetDatabaseStricts)
	if execErr != nil {
		_ = tx.Rollback()
		log.Fatal(execErr)

	}

	_, execErr= tx.Exec(
		sql_scripts.CreateTableUsers)
	if execErr != nil {
		_ = tx.Rollback()
		log.Fatal(execErr)

	}
	_, execErr = tx.Exec(
		sql_scripts.CreateTableAuth)
	if execErr != nil {
		_ = tx.Rollback()
		log.Fatal(execErr)
	}
	_, execErr = tx.Exec(
		sql_scripts.CreateTableAuthProviders)
	if execErr != nil {
		_ = tx.Rollback()
		log.Fatal(execErr)
	}

	_, execErr = tx.Exec(
		sql_scripts.CreateTableAuthentificationHistory)
	if execErr != nil {
		_ = tx.Rollback()
		log.Fatal(execErr)
	}

	_, execErr = tx.Exec(
		sql_scripts.CreateTableUserMfaCode)
	if execErr != nil {
		_ = tx.Rollback()
		log.Fatal(execErr)
	}

	_, execErr = tx.Exec(
		sql_scripts.CreateTableUsersMfaPhone)
	if execErr != nil {
		_ = tx.Rollback()
		log.Fatal(execErr)
	}
	_, execErr = tx.Exec(
		sql_scripts.CreateTableUsersMfaSecret)
	if execErr != nil {
		_ = tx.Rollback()
		log.Fatal(execErr)
	}

	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}
	return nil
}

func TestAnyTimeArgument(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("INSERT INTO users").
		WithArgs("john", AnyTime{}).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = db.Exec("INSERT INTO users(name, created_at) VALUES (?, ?)", "john", time.Now())
	if err != nil {
		t.Errorf("error '%s' was not expected, while inserting a row", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}


func ExampleRows_rowError() {
	db, mock, err := sqlmock.New()
	if err != nil {
		fmt.Println("failed to open sqlmock database:", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "title"}).
		AddRow(0, "one").
		AddRow(1, "two").
		RowError(1, fmt.Errorf("row error"))
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	rs, _ := db.Query("SELECT")
	defer rs.Close()
   fmt.Println(rs)
	for rs.Next() {
		var id int
		var title string
		rs.Scan(&id, &title)
		fmt.Println("scanned id:", id, "and title:", title)
	}

	if rs.Err() != nil {
		fmt.Println("got rows error:", rs.Err())
	}
	// Output: scanned id: 0 and title: one
	// got rows error: row error
}