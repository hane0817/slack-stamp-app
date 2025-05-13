package usecase

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

type RegisterUserInput struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type RegisterUserResult struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// オニオン　内側　applicaion

func RegisterUser(db *sql.DB, registerUserInput RegisterUserInput) (*RegisterUserResult, error) {
	// パスワードをハッシュ化　分割
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerUserInput.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// DB に保存 明示的にトランザクションをはろう　分割
	// db.BeginTx()
	result, err := db.Exec("INSERT INTO users (name, password) VALUES (?, ?)", registerUserInput.Name, string(hashedPassword))
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	return &RegisterUserResult{
		ID:   id,
		Name: registerUserInput.Name,
	}, nil
}
