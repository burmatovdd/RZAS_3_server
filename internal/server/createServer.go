package server

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"rzas_3/internal/server/resource"
)

type Service struct {
	service *Create
}

type Create interface {
	CreateServer()
}

func (service *Service) CreateServer() {
	method := resource.PgService{}
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders: []string{"Content-Type, access-control-allow-origin, access-control-allow-headers, Authorization"},
	}))

	router.GET("/get-info", method.GetInfo)
	router.GET("/get-counter", method.Counter)
	router.GET("/get-db-version", method.GetDBVersion)
	router.GET("/get-db-size", method.GetDBSize)
	router.GET("/get-tb-size", method.GetTbSize)

	router.POST("/add-resource", method.AddResource)
	router.POST("/delete-resource", method.DeleteResource)
	router.POST("/edit-resource", method.EditResource)

	err := router.Run(":8083")
	if err != nil {
		fmt.Println("err: ", err)
	}
}
