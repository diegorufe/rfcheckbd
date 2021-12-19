package service

import (
	"log"
	"rfcheckbd/beans"
	"rfcheckbd/constants"
)

// ProccessDatabseCommands método para procesar los comandos de base de datos
//
// @parameter configuration configuración que tiene todos los parámetros de configuración y comandos a procesar
//
// @returns ---
func ProccessDatabseCommands(configuration beans.Configuration) {
	log.Println("Procesando comandos base de datos")

	// Procesamos los comandos
	for _, command := range configuration.ConfigurationDatabase.Commands {
		log.Printf("Procesando tipo de comando %d", command)

		switch command {
		case constants.Backup:
		case constants.Migrate:
		}
	}
}
