package main

import (
	"database/sql"
	"fmt"
	"image/color"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"backend/controllers"

	"github.com/fogleman/gg"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type RequestData struct {
	Text  string `json:"text"`
	Color string `json:"textColor"`
}

func hexToRGBA(hex string) color.RGBA {
	hex = strings.TrimPrefix(hex, "#")

	var r, g, b, a uint8 = 0, 0, 0, 255
	switch len(hex) {
	case 6:
		fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	case 8:
		fmt.Sscanf(hex, "%02x%02x%02x%02x", &r, &g, &b, &a)
	}

	return color.RGBA{r, g, b, a}
}

func generateImage(text, hexColor string) string {
	const width = 400
	const height = 200

	dc := gg.NewContext(width, height)
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	dc.SetColor(hexToRGBA(hexColor))
	if err := dc.LoadFontFace("/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf", 40); err != nil {
		log.Fatal(err)
	}

	dc.DrawStringAnchored(text, width/2, height/2, 0.5, 0.5)
	outputPath := "output.png"
	dc.SavePNG(outputPath)

	return outputPath
}

func generateHandler(c *gin.Context) {
	var requestData RequestData
	if err := c.ShouldBindJSON(&requestData); err != nil { // BindJsonでも良さそう
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	log.Println("受信データ:", requestData)

	imgPath := generateImage(requestData.Text, requestData.Color)
	if _, err := os.Stat(imgPath); os.IsNotExist(err) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate image"})
		return
	}
	c.File(imgPath)
}

var db *sql.DB

func initDB() {
	var err error
	dsn := "root:root_pass@tcp(db:3306)/database"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 簡単なテーブル作成（必要に応じて実行）
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		password VARCHAR(255) NOT NULL
	);`)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}

	fmt.Println("Database connected!")
}

// ユーザー登録
func registerUser(c *gin.Context) {
	var req struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// パスワードをハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// DB に保存
	result, err := db.Exec("INSERT INTO users (name, password) VALUES (?, ?)", req.Name, string(hashedPassword))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	id, _ := result.LastInsertId()
	c.JSON(http.StatusOK, gin.H{"id": id, "name": req.Name})
}

// JWTの秘密鍵
const SECRET_KEY = "SECRET"

// ログイン処理
func loginUser(c *gin.Context) {
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
	err := db.QueryRow("SELECT password FROM users WHERE name = ?", req.Name).Scan(&storedHashedPassword)
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
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
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

func authMiddleware(c *gin.Context) {
	// Authorizationヘッダーからトークンを取得
	tokenString := c.GetHeader("Authorization")

	// トークンの検証
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		c.Abort()
		return
	}

	c.Next()
}

func main() {
	r := gin.Default()

	//	CORSの設定
	r.Use(cors.Default())

	public := r.Group("/api")

	initDB()
	defer db.Close()

	authGroup := r.Group("/auth")

	authGroup.POST("/register", registerUser)
	authGroup.POST("/login", loginUser)

	public.POST("/generate", generateHandler)

	public.POST("/register", controllers.Register)

	log.Println("Server running on :8080")
	r.Run(":8080")
}
