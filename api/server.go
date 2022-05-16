package api

import (
	db "github.com/CM-IV/mef-api/db/sqlc"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

//Serves HTTP Requests for Posts
type Server struct {
	store  db.Store
	router *gin.Engine
}

//Create new HTTP Server and setup routes
func NewServer(store db.Store) *Server {

	server := &Server{store: store}
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "POST", "HEAD", "DELETE", "OPTIONS", "GET"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.SetTrustedProxies(nil)

	//API GROUP
	api := router.Group("/api")
	{
		//POSTS ENDPOINTS
		api.POST("/posts", server.createPost)
		api.GET("/posts/:id", server.getPost)
		api.GET("/posts", server.listPost)
		api.PUT("/posts/:id", server.updatePost)
		api.DELETE("/posts/:id", server.deletePost)

		//USERS ENDPOINTS
		api.POST("/users", server.createUser)
	}

	server.router = router
	return server

}

//Start runs HTTP Server on a specific address
func (server *Server) Start(address string) error {

	return server.router.Run(address)

}

func errorResponse(err error) gin.H {

	return gin.H{"error": err.Error()}

}

// func CORSMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {

// 		c.Header("Access-Control-Allow-Origin", "*")
// 		c.Header("Access-Control-Allow-Credentials", "true")
// 		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
// 		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

// 		if c.Request.Method == "OPTIONS" {
// 			c.AbortWithStatus(204)
// 			return
// 		}

// 		c.Next()
// 	}
// }
