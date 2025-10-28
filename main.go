package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/h2non/bimg"
)

var apiKey string

func init() {
	apiKey = os.Getenv("IMAGE_API_KEY")
	if apiKey == "" {
		panic("IMAGE_API_KEY is not set")
	}
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("X-API-KEY")
		if key != apiKey {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid API key"})
			return
		}
		c.Next()
	}
}

func convertHandler(c *gin.Context) {
	file, _, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no image provided"})
		return
	}
	defer file.Close()

	buf, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read image"})
		return
	}

	format := strings.ToLower(c.DefaultQuery("format", "avif"))
	width := c.DefaultQuery("width", "")
	height := c.DefaultQuery("height", "")

	options := bimg.Options{}
	switch format {
	case "avif":
		options.Type = bimg.AVIF
	case "webp":
		options.Type = bimg.WEBP
	case "jpeg":
		options.Type = bimg.JPEG
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported format"})
		return
	}
	if width != "" {
		fmt.Sscanf(width, "%d", &options.Width)
	}
	if height != "" {
		fmt.Sscanf(height, "%d", &options.Height)
	}

	img, err := bimg.NewImage(buf).Process(options)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "processing failed"})
		return
	}

	c.Data(http.StatusOK, "image/"+format, img)
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(authMiddleware())
	router.POST("/convert", convertHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
