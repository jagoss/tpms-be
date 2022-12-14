package router

import (
	"be-tpms/src/api/environment"
	"be-tpms/src/api/io"
	"be-tpms/src/api/usecases/cvmodel"
	"be-tpms/src/api/usecases/dogs"
	"be-tpms/src/api/usecases/posts"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// RegisterNewPost godoc
// @Summary Register new post
// @Schemes
// @Description Register new post
// @Tags        posts
// @Accept      json
// @Produce     json
// @Param		dog body model.PostRequest false  "post"
// @Success     200 {object} model.PostResponse
// @Failure		400 {object} object{error=string,message=string}
// @Failure		401 {object} object{error=string,message=string}
// @Failure		422 {object} object{error=string,message=string}
// @Failure		500 {object} object{error=string,message=string}
// @Router      /post [post]
func RegisterNewPost(c *gin.Context, env environment.Env) {
	jsonBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("error reading request body: %v", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   err.Error(),
			"message": "error reading request body!",
		})
		return
	}
	reqPost, err := io.DeserializePost(jsonBody)
	if err != nil {
		log.Printf("error unmarshalling post body: %v", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   err.Error(),
			"message": "error reading post body!",
		})
		return
	}
	post := io.MapFromPostRequest(reqPost)
	if post == nil {
		log.Printf("error mapping dog request when parsing uint")
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   "error parsing uint",
			"message": "error parsing dogRequest to dog!",
		})
		return
	}

	dogManager := dogs.NewDogManager(env.DogPersister, env.Storage)
	dog, err := dogManager.RegisterPostDog(post, reqPost.Image)
	if err != nil {
		log.Printf("%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": fmt.Sprintf("error inserting new dog: %v", err),
		})
		return
	}

	predictionService := cvmodel.NewPrediction(env.DogPersister, env.CVModelRestClient, env.Storage)
	if err = predictionService.CalculateEmbedding(uint(dog.ID)); err != nil {
		log.Printf("%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": fmt.Sprintf("error calculating new dog %d vector: %v", dog.ID, err),
		})
		return
	}

	post.DogId = dog.ID
	postManager := posts.NewPostManager(env.PostPersister)
	post, err = postManager.RegisterPost(post)
	if err != nil {
		log.Printf("%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": fmt.Sprintf("error persisting post: %s", err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, io.MapToPostResponse(post))
}

// GetPost godoc
// @Summary Get post given its ID
// @Schemes
// @Description Get post given its ID
// @Tags        posts
// @Accept      json
// @Produce     json
// @Param		post path string false  "post ID"
// @Success     200 {object} model.PostResponse
// @Failure		400 {object} object{error=string,message=string}
// @Failure		401 {object} object{error=string,message=string}
// @Failure		500 {object} object{error=string,message=string}
// @Router      /post/:id [get]
func GetPost(c *gin.Context, env environment.Env) {
	postID := c.Param("id")
	if postID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Variable missing",
			"message": "Missing dog ID",
		})
		return
	}
	postManager := posts.NewPostManager(env.PostPersister)
	id, _ := strconv.ParseInt(postID, 10, 64)
	post, err := postManager.GetPost(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "error getting post",
		})
		return
	}
	c.JSON(http.StatusOK, io.MapToPostResponse(post))
}

// GetAllPost godoc
// @Summary Get all post given its ID
// @Schemes
// @Description Get all post given its ID
// @Tags        posts
// @Accept      json
// @Produce     json
// @Param		post path string false  "post ID"
// @Success     200 {object} model.PostResponse
// @Failure		400 {object} object{error=string,message=string}
// @Failure		401 {object} object{error=string,message=string}
// @Failure		500 {object} object{error=string,message=string}
// @Router      /post [get]
func GetAllPost(c *gin.Context, env environment.Env) {
	postManager := posts.NewPostManager(env.PostPersister)
	postsList, err := postManager.GetAllPost()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "error getting posts",
		})
		return
	}
	c.JSON(http.StatusOK, io.MapToPostResponseList(postsList))
}
