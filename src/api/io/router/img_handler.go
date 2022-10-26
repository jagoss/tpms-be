package router

import (
	"be-tpms/src/api/environment"
	"be-tpms/src/api/io/storage"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// AddImg godoc
// @Summary Add Image
// @Schemes
// @Description Add image to storage
// @Tags        imgs
// @Param 		img formData []byte false "dog img to save"
// @Accept      json
// @Produce     json
// @Success     200 {string} img "name" saved correctly!
// @Failure		400 {object} object{error=string,message=string}
// @Failure		401 {object} object{error=string,message=string}
// @Failure		500 {object} object{error=string,message=string}
// @Router      /img [post]
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
	var encodedImgs []string
	for _, img := range imgBuffArray {
		encoded := base64.StdEncoding.EncodeToString(img)
		encodedImgs = append(encodedImgs, encoded)
	}
	imgs, err := env.Storage.SaveImgs(encodedImgs)
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
