package controllers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StampController struct {
	//dbを内側に閉じる　Dbとの違い
	db *sql.DB
}

type StampRequestData struct {
	Text            string `json:"text"`
	TextColor       string `json:"textColor"`
	BackgroundColor string `json:"backgroundColor"`
	Language        string `json:"language"`
	SelectedEffect  string `json:"selectedEffect"`
}

func NewStampController(db *sql.DB) *StampController {
	return &StampController{db: db}
}

func (scdb *StampController) PostStampHandler(c *gin.Context) {
	var stamp StampRequestData
	if err := c.ShouldBindJSON(&stamp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := scdb.db.Exec(`
    INSERT INTO stamps (text, language, text_color, selected_effects, background_color, created_at)
    VALUES (?, ?, ?, ?, ?, NOW())`,
		stamp.Text, stamp.Language, stamp.TextColor, stamp.SelectedEffect, stamp.BackgroundColor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "saved"})
}

func (sc *StampController) GETStampHandler(c *gin.Context) {
	rows, err := sc.db.Query(`
		SELECT text, language, text_color, selected_effects, background_color
		FROM stamps
		ORDER BY created_at DESC
		LIMIT 10;
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		return
	}
	defer rows.Close()

	var results []StampRequestData

	for rows.Next() {
		var s StampRequestData
		err := rows.Scan(&s.Text, &s.Language, &s.TextColor, &s.SelectedEffect, &s.BackgroundColor)
		if err != nil {
			continue
		}
		results = append(results, s)
	}

	c.JSON(http.StatusOK, results)
}
