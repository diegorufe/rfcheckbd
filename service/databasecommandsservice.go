package service

import (
	"io/fs"
	"io/ioutil"
	"log"
	"os"
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
			stat, err := os.Stat(configuration.ConfigurationDatabase.PathMigrations)

			// Comrpobamos que el directorio de módulos exista
			if err != nil {
				log.Panicf("Se ha producido un error al comprobar el path de la migración %s", err)
			} else {
				// En el caso de ser una carpeta continuamos leyendo el directorio
				if stat.IsDir() {
					files, err := ioutil.ReadDir(configuration.ConfigurationDatabase.PathMigrations)

					if err != nil {
						log.Panicf("Se ha producido un error al listar el directorio de migraciones %s", err)
					} else {
						processModulesMigratios(configuration, files)
					}
					// Lanzamos error en el caso de que no sea un directorio
				} else {
					log.Panicf("La ruta donde se encuentran los módulos de migraciones no es un directorio. Ruta de módulos de migraciones %s", stat.Name())
				}
			}
		}
	}
}

// processModulesMigratios. Método para procesar los módulos de las migraciones
// @parameter configuration configuración que tiene todos los parámetros de configuración y comandos a procesar
//
// @parameter files array de información de ficheros que serán los módulos a migrar
//
// @returns ---
func processModulesMigratios(configuration beans.Configuration, files []fs.FileInfo) {
	for _, module := range files {
		if module.IsDir() {
			log.Printf("Procesando módulo a migrar %s", module.Name())
		} else {
			log.Printf("Módulo a migrar %s no es un directorio", module.Name())
		}
	}
}
