package main

import (
	"database/sql"
	"fmt"
	"image/color"
	"log"
	"net/http"
	"os"

	"backend/controllers"

	"github.com/fogleman/gg"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type RGBA struct {
	R uint8   `json:"r"`
	G uint8   `json:"g"`
	B uint8   `json:"b"`
	A float64 `json:"a"` // 0.0 ~ 1.0
}

type Effect struct {
	None   bool `json:"none"`
	Glitch bool `json:"glitch"`
	Jitter bool `json:"jitter"`
	Rotate bool `json:"rotate"`
	Shadow bool `json:"shadow"`
	Blur   bool `json:"blur"`
}

type RequestData struct {
	Text            string `json:"text"`
	TextColor       RGBA   `json:"textColor"`
	BackgroundColor RGBA   `json:"backgroundColor"`
	Language        string `json:"language"`
}

func toColorRGBA(input RGBA) color.RGBA {
	return color.RGBA{
		R: input.R,
		G: input.G,
		B: input.B,
		A: uint8(input.A * 255), // float64 → uint8 に変換
	}
}

func generateImage(text string, TextColor RGBA, BackgroundColor RGBA, language string) string {
	const width = 400
	const height = 200

	dc := gg.NewContext(width, height)

	if BackgroundColor.A > 0 {
		bg := color.NRGBA{
			R: BackgroundColor.R,
			G: BackgroundColor.G,
			B: BackgroundColor.B,
			A: uint8(BackgroundColor.A * 255),
		}
		dc.SetColor(bg)
		dc.Clear()
	}
	// 透明の場合は Clear()

	dc.SetColor(toColorRGBA(TextColor))

	if language == "japanese" {
		if err := dc.LoadFontFace("/go/src/font/NotoSansJP-VariableFont_wght.ttf", 40); err != nil {
			log.Fatal(err)
		}
	} else if language == "chinese" {
		if err := dc.LoadFontFace("/go/src/font/NotoSerifSC-Regular.ttf", 40); err != nil {
			log.Fatal(err)
		}
	}

	dc.DrawStringAnchored(text, width/2, height/2, 0.5, 0.5)
	outputPath := "output.png"
	dc.SavePNG(outputPath)

	return outputPath
}

func generateHandler(c *gin.Context) {
	var requestData RequestData
	if err := c.ShouldBindJSON(&requestData); err != nil {
		log.Println("JSON decode error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	log.Println("受信データ:", requestData)

	imgPath := generateImage(requestData.Text, requestData.TextColor, requestData.BackgroundColor, requestData.Language)
	if _, err := os.Stat(imgPath); os.IsNotExist(err) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate image"})
		return
	}
	c.File(imgPath)
}

func InitDB() *sql.DB {
	var err error
	dsn := "root:root_pass@tcp(db:3306)/database"
	db, err := sql.Open("mysql", dsn)
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

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS stamps (
		id INT AUTO_INCREMENT PRIMARY KEY,
		text TEXT NOT NULL,
		language VARCHAR(50) NOT NULL,
		text_color VARCHAR(50) NOT NULL,
		selected_effects TEXT NOT NULL,
		background_color VARCHAR(50) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}

	// レコードがすでに存在するか確認
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM stamps").Scan(&count)
	if err != nil {
		log.Fatal("failed to count stamps: %w", err)
	}

	if count == 0 {
		text := "Hello"
		language := "japanese"
		textColor := "rgba(255, 0, 0, 1)"
		selectedeffects := "none"
		backgroundColor := "rgba(255, 255, 255, 1)"

		_, err = db.Exec("INSERT INTO stamps (`text`, language, text_color, selected_effects, background_color) VALUES (?, ?, ?, ?, ?)", text, language, textColor, selectedeffects, backgroundColor)
		if err != nil {
			log.Fatal("Failed to insert stamp or exit stampdata:", err)
		}
		log.Println("初期データを挿入")
	} else {
		log.Printf("初期データの挿入はスキップ\n")
	}

	fmt.Println("Database connected!")
	return db
}

func main() {
	r := gin.Default()

	db := InitDB()

	//	CORSの設定
	r.Use(cors.Default())

	public := r.Group("/api")

	controllers.InitDB()
	defer controllers.Db.Close()

	authGroup := r.Group("/auth")
	userGroup := r.Group("/user")

	userController := controllers.NewUserController(db)
	userGroup.POST("/register", userController.RegisterUser)

	authGroup.POST("/register", controllers.RegisterUser)
	authGroup.POST("/login", controllers.LoginUser)

	public.POST("/generate", generateHandler)

	stampController := controllers.NewStampController(db)

	public.POST("/stamp/post", stampController.PostStampHandler)
	public.GET("/stamp/get", stampController.GETStampHandler)

	log.Println("Server running on :8080")
	r.Run(":8080")
}
