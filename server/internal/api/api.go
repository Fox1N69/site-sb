package api

import (
	"site-sb/common/http/middleware"
	"site-sb/common/http/request"
	"site-sb/infra"
	v1 "site-sb/internal/api/v1"
	"site-sb/internal/manager"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type Server interface {
	Run()
}

type server struct {
	infra       infra.Infra
	gin         *gin.Engine
	service     manager.ServiceManager
	middleware  middleware.Middleware
	store       sessions.Store
	redisClient *redis.Client
}

func NewServer(infra infra.Infra, redisClient *redis.Client) Server {
	store := sessions.NewCookieStore([]byte(infra.Config().GetString("secret.key")))

	return &server{
		infra:       infra,
		gin:         gin.Default(),
		service:     manager.NewServiceManager(infra),
		middleware:  middleware.NewMiddleware(infra.Config().GetString("secret.key")),
		store:       store,
		redisClient: redisClient,
	}
}

func (c *server) Run() {
	c.gin.Use(c.middleware.CORS())
	c.handlers()
	c.v1()

	c.gin.Run(c.infra.Port())
}

func (c *server) handlers() {
	h := request.DefaultHandler()

	c.gin.NoRoute(h.NoRoute)
	c.gin.GET("/", h.Index)
}

func (c *server) v1() {
	authHandler := v1.NewAuthHandler(c.service.AuthService(), c.infra)
	postHandler := v1.NewPostHandler(c.service.PostService())

	c.gin.Use(sessions.Sessions("user", c.store))

	admin := c.gin.Group("/admin")
	{
		admin.Use(c.middleware.Role("admin"))

		post := admin.Group("/post")
		{
			post.POST("/create", postHandler.Create)
			post.PATCH("/update/:id", postHandler.Update)
			post.DELETE("/delete/:id", postHandler.Delete)
		}
	}

	post := c.gin.Group("/post")
	{
		post.GET("/", postHandler.GetAll)
		post.GET("/:id", postHandler.GetPostByID)
	}

	v1 := c.gin.Group("v1/account")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}
	}
}
