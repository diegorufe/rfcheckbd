package service

import (
	"log"
	"rfcheckbd/beans"
	"rfcheckbd/utils"
)

// ProccessCommands método para procesar los comandos ya cargados en la configuración
//
// @parameter configuration configuración que tiene todos los parámetros de configuración y comandos a procesar
//
// @returns ---
func ProccessCommands(configuration beans.Configuration) {

	log.Println("*** Procesando comandos inicio ***")

	start := utils.MakeTimestamp()

	// Procesamos comandos de base de datos
	ProccessDatabseCommands(configuration)

	end := utils.MakeTimestamp()

	log.Printf("Tiempo tardado en procesar comandos: %d", (end - start))

	log.Println("*** Procesando comandos fin ***")
}
