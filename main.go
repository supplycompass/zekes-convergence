package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.File("index.html")
	})

	r.POST("/upload", func(c *gin.Context) {
		img, err := Convert(c.Request, "file")
		if err != nil {
			panic(err)
		}

		width, err := strconv.Atoi(c.PostForm("width"))
		height, err := strconv.Atoi(c.PostForm("height"))

		thumb, err := ThumbnailJPEG(img, width, height, 100)
		if err != nil {
			fmt.Print(err)
		}

		thumb.Write(c.Writer)
	})

	r.Run(":5000")

}
