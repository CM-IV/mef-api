package api

import (
	"fmt"

	db "github.com/CM-IV/mef-api/db/sqlc"
	"github.com/CM-IV/mef-api/token"
	"github.com/CM-IV/mef-api/util"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

//Serves HTTP Requests for Posts
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

//Create new HTTP Server and setup routes
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	server.setupRouter()
	return server, nil

}

func (server *Server) setupRouter() {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "POST", "HEAD", "DELETE", "OPTIONS", "GET"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.SetTrustedProxies(nil)

	//API GROUP
	api := router.Group("/api")
	{

		authRoutes := router.Group("/api")
		authRoutes.Use(authMiddleware(server.tokenMaker))
		{
			//PROTECTED ENDPOINTS
			//POSTS ENDPOINTS
			authRoutes.PUT("/posts/:id", server.updatePost)
			authRoutes.DELETE("/posts/:id", server.deletePost)
			authRoutes.POST("/posts", server.createPost)
		}
		api.GET("/posts/:id", server.getPost)
		api.GET("/posts", server.listPost)

		//USERS ENDPOINTS
		api.POST("/users", server.createUser)
		api.POST("/users/login", server.loginUser)
	}

	server.router = router
}

//Start runs HTTP Server on a specific address
func (server *Server) Start(address string) error {

	return server.router.Run(address)

}

func errorResponse(err error) gin.H {

	return gin.H{"error": err.Error()}

}
