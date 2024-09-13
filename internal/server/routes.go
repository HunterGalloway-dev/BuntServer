package server

import (
	"BuntServer/internal/coderun"
	"BuntServer/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.Use(gin.Recovery())
	r.Use(CORSMiddleware())

	r.GET("/health", s.serverHealthHandler)
	r.GET("/db_heatlh", s.dbHealthHandler)
	r.GET("/projects", s.getAllProjectsHandler)
	r.GET("/anagrams", s.getAnagramSolution)

	r.POST("/coderun", s.postCodeRunResult)

	return r
}

func (s *Server) serverHealthHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "I'm online"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) dbHealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}

func (s *Server) getAllProjectsHandler(c *gin.Context) {
	var tags []string
	tags = append(tags, "Tag 1")

	projects := s.db.GetAllProjects(tags)

	c.JSON(http.StatusOK, map[string]interface{}{"data": projects})
}

func (s *Server) getAnagramSolution(c *gin.Context) {
	input := c.Query("input")
	words := s.anagramSolver.GetWords(input)

	c.JSON(http.StatusOK, map[string]interface{}{"data": words})
}

func (s *Server) postCodeRunResult(c *gin.Context) {
	var crRequest models.CodeRunRequest

	if err := c.BindJSON(&crRequest); err != nil {
		// Handle error
	}

	result, err := coderun.ProcessCodeRunRequest(&crRequest)

	if err != nil {

	}

	c.IndentedJSON(http.StatusOK, result)
}

func (s *Server) GetProjectTags(c *gin.Context) {
	tags := s.db.GetAllTags()
	c.JSON(http.StatusOK, map[string]interface{}{"data": tags})
}
