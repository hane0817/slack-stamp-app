package main

import (
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
	if err := c.ShouldBindJSON(&requestData); err != nil { // BindJsonでも良さそう
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

func main() {
	r := gin.Default()

	//	CORSの設定
	r.Use(cors.Default())

	public := r.Group("/api")

	controllers.InitDB()
	defer controllers.Db.Close()

	authGroup := r.Group("/auth")

	authGroup.POST("/register", controllers.RegisterUser)
	authGroup.POST("/login", controllers.LoginUser)

	public.POST("/generate", generateHandler)

	public.POST("/register", controllers.Register)

	log.Println("Server running on :8080")
	r.Run(":8080")
}
