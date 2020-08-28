package api


import (
	"github.com/gin-gonic/gin"
	"github.com/zdarovich/win-lose-api/handlers"
	"github.com/zdarovich/win-lose-api/middleware/contenttype"
	"net/http"
)

type (
	router struct {
	}
	// IRouter irouter
	IRouter interface {
		GetEngine() IGINEngine
	}
	// IGINEngine iginengine
	IGINEngine interface {
		Run(addr ...string) (err error)
		ServeHTTP(w http.ResponseWriter, req *http.Request)
	}
)

func NewRouter() IRouter {
	return &router{}
}

func (apiRouter *router) GetEngine() IGINEngine {

	router := gin.Default()

	router.Use(gin.Recovery())
	router.Use(contenttype.New())

	resq := handlers.TrxHandler{}
	router.POST("your_url", resq.PostTransaction)

	return router
}
