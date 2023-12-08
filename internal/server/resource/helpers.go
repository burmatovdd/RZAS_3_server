package resource

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"rzas_3/configs/serverConf"
	"rzas_3/internal/server/postgresql/helpers"
)

func toJson(variable any) []byte {
	jsonStruct, err := json.Marshal(variable)
	if err != nil {
		log.Println("Error: ", err.Error())
	}
	return jsonStruct
}

func (service *PgService) EditResource(c *gin.Context) {
	data := Info{}
	err := c.BindJSON(&data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
		})
		return
	}
	fmt.Println("data: ", data)

	res := helpers.Exec("update resource set waf = $2, ssl = $3, active = $4, owner =$5 where domain = $1",
		[]any{
			data.Domain,
			data.Waf,
			data.SSL,
			data.Active,
			data.Owner,
		},
		serverConf.DefaultConfig)
	if !res {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
	})
}

func (service *PgService) DeleteResource(c *gin.Context) {
	var data string
	err := c.BindJSON(&data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
		})
		return
	}
	res := helpers.Exec("delete from resource where domain = $1", []any{data}, serverConf.DefaultConfig)
	if !res {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"body": res,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"body": res,
	})
}

func (service *PgService) AddResource(c *gin.Context) {
	data := Info{}
	err := c.BindJSON(&data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
		})
		return
	}

	res := helpers.Exec("insert into resource (domain, waf, ssl, active, owner) values ($1,$2,$3,$4,$5)",
		[]any{
			data.Domain,
			data.Waf,
			data.SSL,
			data.Active,
			data.Owner,
		},
		serverConf.DefaultConfig)
	if !res {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
	})
}

func (service *PgService) GetInfo(c *gin.Context) {
	rows, err := helpers.Select("select * from resource ", nil, serverConf.DefaultConfig)
	defer rows.Close()

	info := []Info{}

	for rows.Next() {
		p := Info{}
		err = rows.Scan(
			&p.Domain,
			&p.Waf,
			&p.SSL,
			&p.Active,
			&p.Owner,
		)

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": http.StatusInternalServerError,
			})
			return
		}

		info = append(info, Info{
			Domain: p.Domain,
			Waf:    p.Waf,
			SSL:    p.SSL,
			Active: p.Active,
			Owner:  p.Owner,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"body": info,
	})
}

func (service *PgService) Counter(c *gin.Context) {
	var resources int
	var active int
	var waf int
	var ssl int

	rows, err := helpers.Select("select count(*) from resource", nil, serverConf.DefaultConfig)
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&resources); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
			})
			return
		}
	}
	rows, err = helpers.Select("select count(*) from resource where active = $1", []any{true}, serverConf.DefaultConfig)
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&active); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
			})
			return
		}
	}
	rows, err = helpers.Select("select count(*) from resource where waf = $1", []any{true}, serverConf.DefaultConfig)
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&waf); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
			})
			return
		}
	}
	rows, err = helpers.Select("select count(*) from resource where ssl = $1", []any{true}, serverConf.DefaultConfig)
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&ssl); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
			})
			return
		}
	}

	info := struct {
		Resources int `json:"resources"`
		Waf       int `json:"waf"`
		SSL       int `json:"SSL"`
		Active    int `json:"active"`
	}{
		resources, waf, ssl, active,
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"body": info,
	})
}

func (service *PgService) GetDBVersion(c *gin.Context) {
	var version string

	rows, err := helpers.Select("select version()", nil, serverConf.DefaultConfig)
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&version); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
			})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"body": version,
	})
}

func (service *PgService) GetDBSize(c *gin.Context) {
	var size string

	rows, err := helpers.Select("SELECT pg_size_pretty( pg_database_size( 'RZAS_3' ) );", nil, serverConf.DefaultConfig)
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&size); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
			})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"body": size,
	})
}

func (service *PgService) GetTbSize(c *gin.Context) {
	var size string
	rows, err := helpers.Select("SELECT pg_size_pretty(pg_total_relation_size('public.resource'))", nil, serverConf.DefaultConfig)
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&size); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
			})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"body": size,
	})
}
