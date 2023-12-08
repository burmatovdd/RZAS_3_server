package main

import (
	"log"
	"rzas_3/configs/serverConf"
	"rzas_3/internal/server"
)

func main() {
	//загружаем конфиг на старте приложения
	_, err := serverConf.LoadConfig("configs/serverConf")
	if err != nil {
		log.Fatalln("error in config: ", err)
	}
	service := server.Service{}
	service.CreateServer()
}
