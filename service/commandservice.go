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
// @parameter execDate fecha de ejecucion en formato string
//
// @returns ---
func ProccessCommands(configuration beans.Configuration, execDate string) {

	log.Println("*** Procesando comandos inicio ***")

	start := utils.MakeTimestamp()

	// Procesamos comandos de base de datos
	ProccessDatabseCommands(configuration, execDate)

	end := utils.MakeTimestamp()

	log.Printf("Tiempo tardado en procesar comandos: %d", (end - start))

	log.Println("*** Procesando comandos fin ***")
}
