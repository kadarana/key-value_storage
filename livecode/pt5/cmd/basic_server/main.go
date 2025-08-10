package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	types "myproj/livecode/pt5/internal/pkg/basic_http"
)

/*
http:
"GET" "PUT" "POST" "DELETE"
url
тело запроса
(хедеры)

тело ответа
код ответа
*/
func main() {
	engine := gin.New()

	// REST/SOAP
	// REST - это паттерн, как мы будем пользоваться http, те REST говорит, сто мы будем использовать вот эти глаголы "GET" "PUT" "POST" "DELETE", для описания того, что делает метод, мы будем использовать статусы
	// SOAP говорит о том, что мы не будем использовать разные глаголы, статусы. Вся суть того, счто происходит будет в теле ответа

	engine.GET("/health", func(ctx *gin.Context) {
		time.Sleep(30 * time.Second) // имитируем долгий запрос
		ctx.Status(http.StatusOK)
	})

	engine.POST("/x2", func(ctx *gin.Context) {
		var req types.X2Request
		if err := ctx.BindJSON(&req); err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}

		res := types.X2Response{
			Val: req.Val * 2,
		}

		ctx.JSON(http.StatusOK, res)
	})

	// частый паттер в go
	go func() {
		if err := engine.Run(":7500"); err != nil {
			log.Fatal(err)
		}
	}()
	// engine.Run(":7500")

	// по сети можно общаться на разных уровнях osi
	// самый базовый протокол, на котором можно рабоать udp/tcp

}
