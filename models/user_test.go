package models

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// パスワードの正常系

func TestCreateUserPassword(t *testing.T) {
	os.Chdir("..")                 // プロジェクトのルートに移動する
	db, mock, err := sqlmock.New() // モックデータベース接続
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a gorm database", err)
	}

	// ユーザーデータ
	email := "test@gmail.com"
	password := "Abc123!?"
	hashedPassword := Encrypt(password) // パスワードをハッシュ化

	// Mock設定
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users" ("email","password","created_at") VALUES ($1,$2,$3) RETURNING "id"`)).
		WithArgs(email, hashedPassword, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	// 実行
	repo = &Repository{DB: gdb}
	repo.CreateUser(&User{Email: email, Password: password})

	// モックが期待通りの動作をしたか確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

// メールアドレスの正常系

func TestCreateUserFailsEmailLengthMaXEq(t *testing.T) {
	os.Chdir("..")                 // プロジェクトのルートに移動する
	db, mock, err := sqlmock.New() // モックデータベース接続
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a gorm database", err)
	}

	// ユーザーデータ
	var email string
	emailMaxLength := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa@gmail.com"

	if len(emailMaxLength) == 74 {
		fmt.Println(len(emailMaxLength))
		email = emailMaxLength // メールアドレスは 5~74 文字 (74文字)
	} else {
		fmt.Println(len(emailMaxLength))
		email = ""
	}

	// ユーザーデータ
	password := "Abc123!?"
	hashedPassword := Encrypt(password) // パスワードをハッシュ化

	// Mock設定
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users" ("email","password","created_at") VALUES ($1,$2,$3) RETURNING "id"`)).
		WithArgs(email, hashedPassword, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	// 実行
	repo = &Repository{DB: gdb}
	repo.CreateUser(&User{Email: email, Password: password})

	// モックが期待通りの動作をしたか確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

// メールアドレスの異常系

func TestCreateUserFailsEmailReq(t *testing.T) {
	os.Chdir("..") // プロジェクトのルートに移動する

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a gorm database", err)
	}

	// ユーザーデータ
	email := "" // メールアドレス必須入力
	password := "Abc123!?"

	errMsg, err := LoadConfig("settings/error_messages.json")
	if err != nil {
		fmt.Println("Error loading config:", err)
	}

	repo := Repository{DB: gdb}

	// 実行
	_, err = repo.CreateUser(&User{Email: email, Password: password})
	if err == nil {
		t.Errorf("expected an error but got none")
	} else {
		expectedErrMsg := "email: " + errMsg.EmailRequired + "."
		if err.Error() != expectedErrMsg {
			t.Errorf(`expected error message "%s", but got "%s"`, expectedErrMsg, err.Error())
		}
	}
}

func TestCreateUserFailsEmailFormatMissingLocalPart(t *testing.T) {
	os.Chdir("..") // プロジェクトのルートに移動する

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a gorm database", err)
	}

	// ユーザーデータ
	email := "@gmail.com" // @の前方が無くてもエラー
	password := "Abc123!?"

	errMsg, err := LoadConfig("settings/error_messages.json")
	if err != nil {
		fmt.Println("Error loading config:", err)
	}

	repo := Repository{DB: gdb}

	// 実行
	_, err = repo.CreateUser(&User{Email: email, Password: password})
	if err == nil {
		t.Errorf("expected an error but got none")
	} else {
		expectedErrMsg := "email: " + errMsg.EmailFormat + "."
		if err.Error() != expectedErrMsg {
			t.Errorf(`expected error message "%s", but got "%s"`, expectedErrMsg, err.Error())
		}
	}
}

func TestCreateUserFailsEmailFormat(t *testing.T) {
	os.Chdir("..") // プロジェクトのルートに移動する

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a gorm database", err)
	}

	// ユーザーデータ
	email := "test.gmail.com" // メールアドレスとしての体を成していないパターン
	password := "Abc123!?"

	errMsg, err := LoadConfig("settings/error_messages.json")
	if err != nil {
		fmt.Println("Error loading config:", err)
	}

	repo := Repository{DB: gdb}

	// 実行
	_, err = repo.CreateUser(&User{Email: email, Password: password})
	if err == nil {
		t.Errorf("expected an error but got none")
	} else {
		expectedErrMsg := "email: " + errMsg.EmailFormat + "."
		if err.Error() != expectedErrMsg {
			t.Errorf(`expected error message "%s", but got "%s"`, expectedErrMsg, err.Error())
		}
	}
}

func TestCreateUserFailsEmailLengthTooLong(t *testing.T) {
	os.Chdir("..") // プロジェクトのルートに移動する

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a gorm database", err)
	}

	// ユーザーデータ
	var email string
	tooLongEmail := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa@gmail.com"

	if len(tooLongEmail) == 75 {
		fmt.Println(len(tooLongEmail))
		email = tooLongEmail // メールアドレスの最大文字数 (75文字)
	} else {
		fmt.Println(len(tooLongEmail))
		email = ""
	}

	password := "Abc123!?"

	errMsg, err := LoadConfig("settings/error_messages.json")
	if err != nil {
		fmt.Println("Error loading config:", err)
	}

	repo := Repository{DB: gdb}

	// 実行
	_, err = repo.CreateUser(&User{Email: email, Password: password})
	if err == nil {
		t.Errorf("expected an error but got none")
	} else {
		expectedErrMsg := "email: " + errMsg.EmailFormat + "."
		if err.Error() != expectedErrMsg {
			t.Errorf(`expected error message "%s", but got "%s"`, expectedErrMsg, err.Error())
		}
	}
}

func TestCreateUserFailsEmailFormatNotExistDns(t *testing.T) {
	os.Chdir("..") // プロジェクトのルートに移動する

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a gorm database", err)
	}

	// ユーザーデータ
	email := "t@mail.com.jp" // 存在しないDNS
	password := "Abc123!?"

	errMsg, err := LoadConfig("settings/error_messages.json")
	if err != nil {
		fmt.Println("Error loading config:", err)
	}

	repo := Repository{DB: gdb}

	// 実行
	_, err = repo.CreateUser(&User{Email: email, Password: password})
	if err == nil {
		t.Errorf("expected an error but got none")
	} else {
		expectedErrMsg := "email: " + errMsg.EmailFormat + "."
		if err.Error() != expectedErrMsg {
			t.Errorf(`expected error message "%s", but got "%s"`, expectedErrMsg, err.Error())
		}
	}
}

func TestCreateUserFailsEmailUnique(t *testing.T) {
	os.Chdir("..")                 // プロジェクトのルートに移動する
	db, mock, err := sqlmock.New() // モックデータベース接続
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a gorm database", err)
	}

	// ユーザーデータ1
	email_1 := "testuser1@gmail.com"
	password_1 := "Abc123!?"
	hashedPassword_1 := Encrypt(password_1) // パスワードをハッシュ化

	// ユーザーデータ2
	email_2 := "testuser1@gmail.com" // 既に存在するメールアドレス
	password_2 := "Abc1234!?"
	hashedPassword_2 := Encrypt(password_2) // パスワードをハッシュ化

	// ユーザーデータ1の作成（成功するはず）
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users" ("email","password","created_at") VALUES ($1,$2,$3) RETURNING "id"`)).
		WithArgs(email_1, hashedPassword_1, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	expectedErrMsg := "duplicate key value violates unique constraint" // TranslateErrorsを使用する前のものと比較する

	// ユーザーデータ2の作成（ユニーク制約違反エラーを期待）
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users" ("email","password","created_at") VALUES ($1,$2,$3) RETURNING "id"`)).
		WithArgs(email_2, hashedPassword_2, sqlmock.AnyArg()).
		WillReturnError(fmt.Errorf(expectedErrMsg))
	mock.ExpectRollback()

	// ユーザーデータ1の作成（成功するはず）
	repo = &Repository{DB: gdb}
	repo.CreateUser(&User{Email: email_1, Password: password_1})

	// ユーザーデータ1の作成（ユニーク制約違反エラーを期待）
	_, err = repo.CreateUser(&User{Email: email_2, Password: password_2})
	if err == nil {
		t.Errorf("expected an error but got none")
	} else if err.Error() != expectedErrMsg {
		t.Errorf(`expected error message "%s", but got "%s"`, expectedErrMsg, err.Error())
	}

	// モックが期待通りの動作をしたか確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

// パスワードの異常系

func TestCreateUserFailsPasswordReq(t *testing.T) {
	os.Chdir("..") // プロジェクトのルートに移動する

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a gorm database", err)
	}

	// ユーザーデータ
	email := "test@gmail.com"
	password := "" // パスワードの必須入力

	errMsg, err := LoadConfig("settings/error_messages.json")
	if err != nil {
		fmt.Println("Error loading config:", err)
	}

	repo := Repository{DB: gdb}

	// 実行
	_, err = repo.CreateUser(&User{Email: email, Password: password})
	if err == nil {
		t.Errorf("expected an error but got none")
	} else {
		expectedErrMsg := "password: " + errMsg.PasswordRequired + "."
		if err.Error() != expectedErrMsg {
			t.Errorf(`expected error message "%s", but got "%s"`, expectedErrMsg, err.Error())
		}
	}
}

func TestCreateUserFailsPasswordSpecialChar(t *testing.T) {
	os.Chdir("..") // プロジェクトのルートに移動する

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a gorm database", err)
	}

	// ユーザーデータ
	email := "test@gmail.com"
	password := "password123" // !?-_ の記号が含まれない

	errMsg, err := LoadConfig("settings/error_messages.json")
	if err != nil {
		fmt.Println("Error loading config:", err)
	}

	repo = &Repository{DB: gdb}

	// 実行
	_, err = repo.CreateUser(&User{Email: email, Password: password})
	if err == nil {
		t.Errorf("expected an error but got none")
	} else {
		expectedErrMsg := "password: " + errMsg.PasswordSpecialChar + "."
		if err.Error() != expectedErrMsg {
			t.Errorf(`expected error message "%s", but got "%s"`, expectedErrMsg, err.Error())
		}
	}
}

func TestCreateUserFailsPasswordAlphabet(t *testing.T) {
	os.Chdir("..") // プロジェクトのルートに移動する

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a gorm database", err)
	}

	// ユーザーデータ
	email := "test@gmail.com"
	password := "123456789" // 数字だけもx

	errMsg, err := LoadConfig("settings/error_messages.json")
	if err != nil {
		fmt.Println("Error loading config:", err)
	}

	repo := Repository{DB: gdb}

	// 実行
	_, err = repo.CreateUser(&User{Email: email, Password: password})
	if err == nil {
		t.Errorf("expected an error but got none")
	} else {
		expectedErrMsg := "password: " + errMsg.PasswordAlphabet + "."
		if err.Error() != expectedErrMsg {
			t.Errorf(`expected error message "%s", but got "%s"`, expectedErrMsg, err.Error())
		}
	}
}

func TestCreateUserFailsPasswordNumber(t *testing.T) {
	os.Chdir("..") // プロジェクトのルートに移動する

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a gorm database", err)
	}

	// ユーザーデータ
	email := "test@gmail.com"
	password := "onlyAlphabet" // 英字だけはx

	errMsg, err := LoadConfig("settings/error_messages.json")
	if err != nil {
		fmt.Println("Error loading config:", err)
	}

	repo := Repository{DB: gdb}

	// 実行
	_, err = repo.CreateUser(&User{Email: email, Password: password})
	if err == nil {
		t.Errorf("expected an error but got none")
	} else {
		expectedErrMsg := "password: " + errMsg.PasswordNumber + "."
		if err.Error() != expectedErrMsg {
			t.Errorf(`expected error message "%s", but got "%s"`, expectedErrMsg, err.Error())
		}
	}
}

func TestCreateUserFailsPasswordTooShort(t *testing.T) {
	os.Chdir("..") // プロジェクトのルートに移動する

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a gorm database", err)
	}

	// ユーザーデータ
	email := "test@gmail.com"
	password := "abc1234" // パスワードは8文字以上(20)以下を満たすこと 7文字

	errMsg, err := LoadConfig("settings/error_messages.json")
	if err != nil {
		fmt.Println("Error loading config:", err)
	}

	repo := Repository{DB: gdb}

	// 実行
	_, err = repo.CreateUser(&User{Email: email, Password: password})
	if err == nil {
		t.Errorf("expected an error but got none")
	} else {
		expectedErrMsg := "password: " + errMsg.PasswordLength + "."
		if err.Error() != expectedErrMsg {
			t.Errorf(`expected error message "%s", but got "%s"`, expectedErrMsg, err.Error())
		}
	}
}

func TestCreateUserPasswordTooLong(t *testing.T) {
	os.Chdir("..") // プロジェクトのルートに移動する

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a gorm database", err)
	}

	// ユーザーデータ
	email := "test@gmail.com"
	password := "Abc123!?sdsdssssddddsss" // パスワードは8文字以上(20)以下を満たすこと 21文字

	errMsg, err := LoadConfig("settings/error_messages.json")
	if err != nil {
		fmt.Println("Error loading config:", err)
	}

	repo := Repository{DB: gdb}

	// 実行
	_, err = repo.CreateUser(&User{Email: email, Password: password})
	if err == nil {
		t.Errorf("expected an error but got none")
	} else {
		expectedErrMsg := "password: " + errMsg.PasswordLength + "."
		if err.Error() != expectedErrMsg {
			t.Errorf(`expected error message "%s", but got "%s"`, expectedErrMsg, err.Error())
		}
	}
}

func TestCreateUserPasswordMixedCaseLower(t *testing.T) {
	os.Chdir("..") // プロジェクトのルートに移動する

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a gorm database", err)
	}

	// ユーザーデータ
	email := "test@gmail.com"
	password := "abc123!?" // 英字は小文字大文字混合　小文字のみ

	errMsg, err := LoadConfig("settings/error_messages.json")
	if err != nil {
		fmt.Println("Error loading config:", err)
	}

	repo := Repository{DB: gdb}

	// 実行
	_, err = repo.CreateUser(&User{Email: email, Password: password})
	if err == nil {
		t.Errorf("expected an error but got none")
	} else {
		expectedErrMsg := "password: " + errMsg.PasswordMixedCase + "."
		if err.Error() != expectedErrMsg {
			t.Errorf(`expected error message "%s", but got "%s"`, expectedErrMsg, err.Error())
		}
	}
}

func TestCreateUserPasswordMixedCase(t *testing.T) {
	os.Chdir("..") // プロジェクトのルートに移動する

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a gorm database", err)
	}

	// ユーザーデータ
	email := "test@gmail.com"
	password := "ABC123!?" // 英字は小文字大文字混合 大文字のみ

	errMsg, err := LoadConfig("settings/error_messages.json")
	if err != nil {
		fmt.Println("Error loading config:", err)
	}

	repo := Repository{DB: gdb}

	// 実行
	_, err = repo.CreateUser(&User{Email: email, Password: password})
	if err == nil {
		t.Errorf("expected an error but got none")
	} else {
		expectedErrMsg := "password: " + errMsg.PasswordMixedCase + "."
		if err.Error() != expectedErrMsg {
			t.Errorf(`expected error message "%s", but got "%s"`, expectedErrMsg, err.Error())
		}
	}
}
