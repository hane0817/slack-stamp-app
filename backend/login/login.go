package login

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

func login() {
    router := gin.Default()

    public := router.Group("/api")

    public.POST("/register", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "data": "this is the register endpoint.",
        })
    })

    router.Run(":8080")
}