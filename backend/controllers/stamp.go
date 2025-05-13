package controllers

import (
	"database/sql"
	"encoding/json"
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

	jsonBytes, err := json.Marshal(stamp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode JSON"})
		return
	}

	_, err = scdb.db.Exec(`
        INSERT INTO stamps (json_data, created_at)
        VALUES (?, NOW())`, string(jsonBytes))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "saved"})
}

func (sc *StampController) GETStampHandler(c *gin.Context) {
	rows, err := sc.db.Query(`SELECT json_data FROM stamps ORDER BY created_at DESC LIMIT 20`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		return
	}
	defer rows.Close()

	var results []map[string]interface{}

	for rows.Next() {
		var jsonStr string
		if err := rows.Scan(&jsonStr); err != nil {
			continue
		}

		var data map[string]interface{}
		if err := json.Unmarshal([]byte(jsonStr), &data); err == nil {
			results = append(results, data)
		}
	}

	c.JSON(http.StatusOK, results)
}
