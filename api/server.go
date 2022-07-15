package api

import (
	"fmt"

	db "github.com/CM-IV/mef-api/db/sqlc"
	"github.com/CM-IV/mef-api/token"
	"github.com/CM-IV/mef-api/util"
	"github.com/Depado/ginprom"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	p := ginprom.New(
		ginprom.Engine(router),
		ginprom.Subsystem("gin"),
		ginprom.Ignore("/api/metrics"),
	)
	router.Use(p.Instrument())

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"PUT", "POST", "HEAD", "DELETE", "OPTIONS", "GET"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.SetTrustedProxies(nil)

	//API GROUP
	api := router.Group("/api")
	{
		api.GET("/posts/:id", server.getPost)
		api.GET("/posts", server.listPost)

		//USERS ENDPOINTS
		api.POST("/users", server.createUser)
		api.POST("/users/login", server.loginUser)

		authRoutes := router.Group("/api")
		authRoutes.Use(authMiddleware(server.tokenMaker))
		{
			//PROTECTED ENDPOINTS
			//POSTS ENDPOINTS
			authRoutes.PUT("/posts/:id", server.updatePost)
			authRoutes.DELETE("/posts/:id", server.deletePost)
			authRoutes.POST("/posts", server.createPost)
			authRoutes.GET("/metrics", gin.WrapH(promhttp.Handler()))
		}

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
