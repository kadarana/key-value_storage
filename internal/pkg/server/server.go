package server

import (
	"myproj/internal/pkg/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	host    string
	storage *storage.Storage
}

type Entry struct {
	Value any `json:"value"`
}

type EntryArray struct {
	Value []any `json:"value"`
}

type EntryList struct {
	Slice []int `json:"slice"`
}

type EntryLSET struct {
	Index   int `json:"index"`
	Element any `json:"element"`
}

type EntryLGET struct {
	Index int `json:"index"`
}

func New(st *storage.Storage) *Server {
	s := &Server{
		host:    ":8090",
		storage: st,
	}

	return s
}

func (r *Server) newAPI() *gin.Engine {
	engine := gin.New()

	engine.GET("health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "OK")
	})

	engine.POST("/scalar/set/:key", r.handlerSet)
	engine.GET("/scalar/get/:key", r.handlerGet)

	engine.POST("/hash/set/:key/:field", r.handlerHSET)
	engine.POST("hash/get/:key/:field", r.handlerHGET)

	engine.POST("/array/lpush/:key", r.handlerLPUSH)
	engine.GET("/array/lpop/:key", r.handlerLPOP)

	engine.POST("/array/rpush/:key", r.handlerRPUSH)
	engine.POST("/array/raddtoset/:key", r.handlerRADDTOSET)
	engine.GET("/array/rpop/:key", r.handlerRPOP)

	engine.POST("/array/lset/:key", r.handlerLSET)
	engine.GET("/array/lget/:key", r.handlerLGET)

	return engine
}
