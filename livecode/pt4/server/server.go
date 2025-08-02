package server

import (
	"encoding/json"
	"myproj/internal/pkg/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	//engine *gin.Engine
	host    string
	storage *storage.Storage
}

type Entry struct {
	Value string
}

func New(host string, st *storage.Storage) *Server {
	s := &Server{
		host:    host,
		storage: st,
	}

	return s
}

// gin framework
func (r *Server) newAPI() *gin.Engine {
	engine := gin.New()

	// engine.GET("/", func(ctx *gin.Context) { //  func(ctx *gin.Context) - это хендлер
	// 	ctx.String(http.StatusOK, "Hello World!")
	// })

	engine.GET("hello-world", func(ctx *gin.Context) { //  func(ctx *gin.Context) - это хендлер
		ctx.JSON(http.StatusOK, "Hello World!")
	})
	// /key/a  /key/abcd
	engine.PUT("/key/:key", r.handlerSet)
	engine.GET("/key/:key", r.handlerGet)
	return engine
}
func (r *Server) handlerSet(ctx *gin.Context) {
	key := ctx.Param("key")

	var v Entry

	if err := json.NewDecoder(ctx.Request.Body).Decode(&v); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return 
	}

	r.storage.Set(key, v.Value)

	ctx.Status(http.StatusOK)
}
func (r *Server) handlerGet(ctx *gin.Context) {
	key := ctx.Param("key")

	v := r.storage.Get(key)
	if v == nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctx.JSON(http.StatusOK, Entry{Value : *v})
}
func (r *Server) Start() {
	// var req http.Request // это структура
	// var resp http.Response // это структура

	// это штука, которая строит дерево и в зависимости от того какой реквест пришел, он подбирает хендлер определенный для этого дерева
	// http.NewServeMux()

	r.newAPI().Run(r.host) // запустит сервер на каком-то адресе

}
