package server

import (
	"BuntServer/internal/coderun"
	"BuntServer/internal/models"
	"net/http"
	"strings"

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
	r.GET("/projects/tags", s.getProjectTags)
	r.POST("/projects", s.postProject)
	r.DELETE("/projects/:id", s.deleteProject)

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

	tagQuery := c.DefaultQuery("tags", "notags")
	if tagQuery != "notags" {
		tags = strings.Split(tagQuery, ",")
	}

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

func (s *Server) getProjectTags(c *gin.Context) {
	tags := s.db.GetAllTags()
	c.JSON(http.StatusOK, map[string]interface{}{"data": tags})
}

func (s *Server) postProject(c *gin.Context) {
	var project models.Project
	if err := c.Bind(&project); err != nil {
		// Return error
		c.JSON(http.StatusBadRequest, "Failed to post project")
	}

	// Upset into DB
	id := s.db.PostProject(project)

	project.Id = id

	c.IndentedJSON(http.StatusCreated, project)
}

func (s *Server) deleteProject(c *gin.Context) {
	id := c.Param("id")

	res := s.db.DeleteProject(id)

	if res {
		c.JSON(http.StatusOK, "Document deleted")
	} else {
		c.JSON(http.StatusBadRequest, "Bad request")
	}
}
