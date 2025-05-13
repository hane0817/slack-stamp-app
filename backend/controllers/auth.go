package controllers

import (
	"backend/usecase"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": "validated!",
	})
}

// ユーザー登録　オニオン　外側　infrastructure
func RegisterUser(c *gin.Context) {
	var req struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// コントローラーがアプリケーションに依存するように入れなおす
	input := usecase.RegisterUserInput{
		Name:     req.Name,
		Password: req.Password,
	}

	result, err := usecase.RegisterUser(Db, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": result.ID, "name": result.Name})
}

// JWTの秘密鍵
const SECRET_KEY = "SECRET"

// ログイン処理
func LoginUser(c *gin.Context) {
	var req struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// ユーザーを取得
	var storedHashedPassword string
	err := Db.QueryRow("SELECT password FROM users WHERE name = ?", req.Name).Scan(&storedHashedPassword)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// パスワードの検証
	if err := bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// トークンを発行

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": req.Name,
		"exp":  time.Now().Add(time.Minute * 3).Unix(),
	})
	tokenString, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error signing token"})
		return
	}

	// ヘッダーにトークンをセット
	c.Header("Authorization", tokenString)

	c.JSON(http.StatusOK, gin.H{"message": "Login successful!"})
}

var Db *sql.DB

// main関数で呼び出す　定義
func InitDB() {
	var err error
	dsn := "root:root_pass@tcp(db:3306)/database"
	Db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 簡単なテーブル作成（必要に応じて実行）
	_, err = Db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		password VARCHAR(255) NOT NULL
	);`)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}

	fmt.Println("Database connected!")
}
