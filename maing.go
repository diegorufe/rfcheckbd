package main

import (
	"log"
	"rfcheckbd/config"
	"rfcheckbd/service"
	"rfcheckbd/utils"
	"time"
)

// Método para lanzar la aplicación
func main() {

	start := utils.MakeTimestamp()

	execDate := time.Now().Format("2006-01-02 03:04:05")

	// Arrancamos la configuración
	configuration := config.Config(execDate)

	//Arrancamos la lectura de comadnso
	service.ProccessCommands(configuration, execDate)

	end := utils.MakeTimestamp()

	log.Printf("*** Tiempo tardado en realizar operaciones del ejecutable: %d ms", (end - start))
}
