package main

import (
	"log"
	"rfcheckbd/config"
	"rfcheckbd/service"
	"rfcheckbd/utils"
)

// Método para lanzar la aplicación
func main() {

	start := utils.MakeTimestamp()

	// Arrancamos la configuración
	configuration := config.Config()

	//Arrancamos la lectura de comadnso
	service.ProccessCommands(configuration)

	end := utils.MakeTimestamp()

	log.Printf("*** Tiempo tardado en realizar operaciones del ejecutable: %d", (end - start))
}
