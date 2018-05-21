package main

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"path"
	"path/filepath"
	"os"
	"github.com/gin-contrib/cors"
)

func main() {
	basePath := os.Args[1]
	r := gin.Default()
	r.Use(cors.Default())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/boards", func(c *gin.Context) {
		boards, err := listFiles(basePath,"boards", true)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}

		c.JSON(200, boards)
	})
	r.GET("/boards/:board", func(c *gin.Context) {
		board := c.Param("board")
		boardPath := filepath.Join(basePath, filepath.FromSlash(path.Clean("/"+board)))
		prefix, _ := filepath.Rel(basePath, boardPath)
		prefix = "boards/" + prefix

		albums, err := listFiles(boardPath, prefix, false)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}

		c.JSON(200, albums)
	})
	r.GET("/boards/:board/:album", func(c *gin.Context) {
		board := c.Param("board")
		album := c.Param("album")
		albumPath := filepath.Join(basePath, filepath.FromSlash(path.Clean("/"+board)), filepath.FromSlash(path.Clean("/"+album)))
		prefix, _ := filepath.Rel(basePath, albumPath)
		prefix = "boards/" + prefix

		photos, err := listFiles(albumPath, prefix, false)
		if err != nil {
			c.AbortWithError(500, err)
			return
		}

		c.JSON(200, photos)
	})
	r.GET("/boards/:board/:album/:photo", func(c *gin.Context) {
		board := c.Param("board")
		album := c.Param("album")
		photo := c.Param("photo")
		photoPath := filepath.Join(basePath, filepath.FromSlash(path.Clean("/"+board)), filepath.FromSlash(path.Clean("/"+album)), filepath.FromSlash(path.Clean("/"+photo)))
		if _, err := os.Stat(photoPath); !os.IsNotExist(err) {
			c.File(photoPath)
		} else {
			c.AbortWithError(404, err)
		}
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}

func listFiles(path string, prefix string, reverse bool) ([]map[string]interface{}, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	fileNames := make([]map[string]interface{}, len(files))
	for i := range files {
		j := i
		if reverse {
			j = len(files) - 1 - i
		}
		fileNames[j] = gin.H{
			"title": files[i].Name(),
			"path":  prefix + "/" + files[i].Name(),
		}
	}
	return fileNames, nil
}
