package main

import (
	"fmt"
	"image/color"
	"log"
	"net/http"
	"os"
	"strings"

	"backend/controllers"

	"github.com/fogleman/gg"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type RequestData struct {
	Text     string `json:"text"`
	Color    string `json:"textColor"`
	Language string `json:"language"`
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

func generateImage(text, hexColor, language string) string {
	const width = 400
	const height = 200

	dc := gg.NewContext(width, height)
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	dc.SetColor(hexToRGBA(hexColor))

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

	imgPath := generateImage(requestData.Text, requestData.Color, requestData.Language)
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
