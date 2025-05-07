package main

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math/rand"
	"net/http"
	"time"

	"backend/controllers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

type RGBA struct {
	R uint8 `json:"r"`
	G uint8 `json:"g"`
	B uint8 `json:"b"`
	A uint8 `json:"a"`
}

type EffectRequest struct {
	Text            string          `json:"text"`
	TextColor       RGBA            `json:"textColor"`
	BackgroundColor RGBA            `json:"backgroundColor"`
	Language        string          `json:"language"`
	Effect          map[string]bool `json:"effect"`
}

// func toColorRGBA(input RGBA) color.RGBA {
// 	return color.RGBA{
// 		R: input.R,
// 		G: input.G,
// 		B: input.B,
// 		A: uint8(input.A * 255), // float64 → uint8 に変換
// 	}
// }

// func generateImage(text string, TextColor RGBA, BackgroundColor RGBA, language string) string {
// 	const width = 400
// 	const height = 200

// 	dc := gg.NewContext(width, height)

// 	if BackgroundColor.A > 0 {
// 		bg := color.NRGBA{
// 			R: BackgroundColor.R,
// 			G: BackgroundColor.G,
// 			B: BackgroundColor.B,
// 			A: uint8(BackgroundColor.A * 255),
// 		}
// 		dc.SetColor(bg)
// 		dc.Clear()
// 	}

// 	dc.SetColor(toColorRGBA(TextColor))

// 	if language == "japanese" {
// 		if err := dc.LoadFontFace("/go/src/font/NotoSansJP-VariableFont_wght.ttf", 40); err != nil {
// 			log.Fatal(err)
// 		}
// 	} else if language == "chinese" {
// 		if err := dc.LoadFontFace("/go/src/font/NotoSerifSC-Regular.ttf", 40); err != nil {
// 			log.Fatal(err)
// 		}
// 	}

// 	dc.DrawStringAnchored(text, width/2, height/2, 0.5, 0.5)
// 	outputPath := "output.png"
// 	dc.SavePNG(outputPath)

// 	return outputPath
// }

// func generateHandler(c *gin.Context) {
// 	var requestData RequestData
// 	if err := c.ShouldBindJSON(&requestData); err != nil { // BindJsonでも良さそう
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
// 		return
// 	}

// 	log.Println("受信データ:", requestData)

// 	imgPath := generateImage(requestData.Text, requestData.TextColor, requestData.BackgroundColor, requestData.Language)
// 	if _, err := os.Stat(imgPath); os.IsNotExist(err) {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate image"})
// 		return
// 	}
// 	c.File(imgPath)
// }

func encodeToPNG(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	return buf.Bytes(), err
}

func applyFog(img *ebiten.Image) {
	w, h := img.Size()
	fog := ebiten.NewImage(w, h)
	alpha := uint8(40)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c := color.RGBA{200, 200, 200, alpha}
			fog.Set(x, y, c)
		}
	}
	img.DrawImage(fog, nil)
}

func applyGlitch(img *ebiten.Image) {
	w, h := img.Size()
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 5; i++ {
		x := rand.Intn(w / 2)
		y := rand.Intn(h)
		sw := rand.Intn(w - x)
		sh := rand.Intn(10)
		src := img.SubImage(image.Rect(x, y, x+sw, y+sh)).(*ebiten.Image)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x+rand.Intn(10)-5), float64(y+rand.Intn(5)-2))
		img.DrawImage(src, op)
	}
}

func applySparkle(img *ebiten.Image) {
	w, h := img.Size()
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 50; i++ {
		x := rand.Intn(w)
		y := rand.Intn(h)
		img.Set(x, y, color.White)
	}
}

func toRGBA(c RGBA) color.RGBA {
	return color.RGBA{R: c.R, G: c.G, B: c.B, A: c.A}
}

func GenerateImage(req EffectRequest) (image.Image, error) {
	width, height := 256, 256

	bg := toRGBA(req.BackgroundColor)
	canvas := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(canvas, canvas.Bounds(), &image.Uniform{bg}, image.Point{}, draw.Src)

	ebitenImg := ebiten.NewImageFromImage(canvas)

	textColor := toRGBA(req.TextColor)
	text.Draw(ebitenImg, req.Text, basicfont.Face7x13, 20, height/2, textColor)

	if req.Effect["fog"] {
		applyFog(ebitenImg)
	}
	if req.Effect["glitch"] {
		applyGlitch(ebitenImg)
	}
	if req.Effect["sparkle"] {
		applySparkle(ebitenImg)
	}

	final := image.NewRGBA(ebitenImg.Bounds())
	ebitenImg.ReadPixels(final.Pix)
	return final, nil
}

func generateHandlerEbiten(c *gin.Context) {
	var req EffectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	img, err := GenerateImage(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "image generation failed"})
		return
	}

	c.Header("Content-Type", "image/png")
	pngBytes, _ := encodeToPNG(img)
	c.Writer.Write(pngBytes)
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

	public.POST("/generate", generateHandlerEbiten)

	public.POST("/register", controllers.Register)

	log.Println("Server running on :8080")
	r.Run(":8080")
}
