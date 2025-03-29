package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/fogleman/gg"
)

type RequestData struct {
	Text  string `json:"text"`
	Color string `json:"textColor"`
}

// HEX 文字列を color.RGBA に変換する関数
func hexToRGBA(hex string) color.RGBA {
	hex = strings.TrimPrefix(hex, "#") // `#` を削除

	var r, g, b, a uint8 = 0, 0, 0, 255 // デフォルトの色（黒）+ 不透明
	switch len(hex) {
	case 6: // RGB (例: #FF0000)
		fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	case 8: // RGBA (例: #FF0000FF)
		fmt.Sscanf(hex, "%02x%02x%02x%02x", &r, &g, &b, &a)
	}

	return color.RGBA{r, g, b, a}
}

func generateImage(text, hexColor string) string {
	const width = 400
	const height = 200

	dc := gg.NewContext(width, height)
	dc.SetRGB(1, 1, 1) // 背景色（白）
	dc.Clear()

	dc.SetColor(hexToRGBA(hexColor)) // ユーザー指定の色を適用
	if err := dc.LoadFontFace("/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf", 40); err != nil {
		log.Fatal(err)
	}

	dc.DrawStringAnchored(text, width/2, height/2, 0.5, 0.5)
	outputPath := "output.png"
	dc.SavePNG(outputPath)

	return outputPath
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // CORS 対応
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var requestData RequestData
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	log.Println("受信データ:", requestData) // ← デバッグ用のログを追加

	imgPath := generateImage(requestData.Text, requestData.Color)
	imgFile, err := os.Open(imgPath)
	if err != nil {
		http.Error(w, "Failed to generate image", http.StatusInternalServerError)
		return
	}
	defer imgFile.Close()

	w.Header().Set("Content-Type", "image/png")
	http.ServeFile(w, r, imgPath)
}

func main() {
	http.HandleFunc("/generate", handler)
	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}
