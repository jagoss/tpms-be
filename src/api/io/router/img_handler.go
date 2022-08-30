package router

import (
	"be-tpms/src/api/environment"
	"be-tpms/src/api/io/storage"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func AddImg(c *gin.Context, env environment.Env) {
	_, fileHeader, err := c.Request.FormFile("img")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Failed to upload img",
		})
		return
	}
	imgBuffArray, path, err := storage.SaveTempImg(c, fileHeader)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to save temporal img",
		})
		return
	}

	imgs, err := env.Storage.SaveImgs(imgBuffArray)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to upload img",
		})
		return
	}
	err = storage.DeleteImg(path)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to upload img",
		})
		return
	}

	message := fmt.Sprintf("img %s saved correctly!", imgs)
	c.String(http.StatusOK, message)
}
