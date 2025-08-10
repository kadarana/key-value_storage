package server

import (
	"encoding/json"
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

	engine.POST("/array/lpush/:key", r.handlerLPUSH)
	engine.GET("/array/lpop/:key", r.handlerLPOP)

	engine.POST("/array/rpush/:key", r.handlerRPUSH)
	engine.POST("/array/raddtoset/:key", r.handlerRADDTOSET)
	engine.GET("/array/rpop/:key", r.handlerRPOP)

	engine.POST("/array/lset/:key", r.handlerLSET)
	engine.GET("/array/lget/:key", r.handlerLGET)

	return engine
}

func (r *Server) handlerSet(ctx *gin.Context) {
	key := ctx.Param("key")

	var v Entry
	if err := json.NewDecoder(ctx.Request.Body).Decode(&v); err != nil {
		ctx.AbortWithStatus(http.StatusBadGateway)
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

	ctx.JSON(http.StatusOK, Entry{Value: *v})
}

func (r *Server) handlerLPUSH(ctx *gin.Context) {
	key := ctx.Param("key")

	var v EntryArray
	if err := json.NewDecoder(ctx.Request.Body).Decode(&v); err != nil {
		ctx.AbortWithStatus(http.StatusBadGateway)
		return
	}

	err := r.storage.LPUSH(key, v.Value)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
			"status":  "false",
			"message": err.Error(),
		})

		return
	}

	ctx.Status(http.StatusOK)

}

func (r *Server) handlerLPOP(ctx *gin.Context) {
	key := ctx.Param("key")

	var v EntryList
	if err := json.NewDecoder(ctx.Request.Body).Decode(&v); err != nil {
		ctx.AbortWithStatus(http.StatusBadGateway)
		return
	}

	val, err := r.storage.LPOP(key, v.Slice...)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
			"status":  "false",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Entry{Value: val})
}

func (r *Server) handlerRPUSH(ctx *gin.Context) {
	key := ctx.Param("key")

	var v EntryArray
	if err := json.NewDecoder(ctx.Request.Body).Decode(&v); err != nil {
		ctx.AbortWithStatus(http.StatusBadGateway)
		return
	}

	err := r.storage.RPUSH(key, v.Value)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
			"status":  "false",
			"message": err.Error(),
		})

		return
	}

	ctx.Status(http.StatusOK)
}

func (r *Server) handlerRADDTOSET(ctx *gin.Context) {
	key := ctx.Param("key")

	var v EntryArray
	if err := json.NewDecoder(ctx.Request.Body).Decode(&v); err != nil {
		ctx.AbortWithStatus(http.StatusBadGateway)
		return
	}

	err := r.storage.RADDTOSET(key, v.Value)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
			"status":  "false",
			"message": err.Error(),
		})
		return
	}

	ctx.Status(http.StatusOK)
}

func (r *Server) handlerRPOP(ctx *gin.Context) {
	key := ctx.Param("key")

	var v EntryList
	if err := json.NewDecoder(ctx.Request.Body).Decode(&v); err != nil {
		ctx.AbortWithStatus(http.StatusBadGateway)
		return
	}

	val, err := r.storage.RPOP(key, v.Slice...)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
			"status":  "false",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Entry{Value: val})

}
func (r *Server) handlerLSET(ctx *gin.Context) {
	key := ctx.Param("key")

	var v EntryLSET
	if err := json.NewDecoder(ctx.Request.Body).Decode(&v); err != nil {
		ctx.AbortWithStatus(http.StatusBadGateway)
		return
	}

	val, err := r.storage.LSET(key, v.Index, v.Element)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  "false",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Entry{Value: val})
}

func (r *Server) handlerLGET(ctx *gin.Context) {
	key := ctx.Param("key")

	var v EntryLGET
	if err := json.NewDecoder(ctx.Request.Body).Decode(&v); err != nil {
		ctx.AbortWithStatus(http.StatusBadGateway)
		return
	}

	val, err := r.storage.LGET(key, v.Index)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
			"status":  "fasle",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Entry{Value: val})
}
