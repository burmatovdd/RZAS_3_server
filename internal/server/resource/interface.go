package resource

import "github.com/gin-gonic/gin"

type PgService struct {
	service *Service
}

type Service interface {
	GetInfo(c *gin.Context)
	AddResource(c *gin.Context)
	DeleteResource(c *gin.Context)
	EditResource(c *gin.Context)
	Counter(c *gin.Context)
	GetDBVersion(c *gin.Context)
	GetDBSize(c *gin.Context)
	GetTbSize(c *gin.Context)
}
