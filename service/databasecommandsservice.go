package service

import (
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"rfcheckbd/beans"
	"rfcheckbd/constants"
	"sort"
	"strconv"
)

// ProccessDatabseCommands método para procesar los comandos de base de datos
//
// @parameter configuration configuración que tiene todos los parámetros de configuración y comandos a procesar
//
// @returns ---
func ProccessDatabseCommands(configuration beans.Configuration) {
	log.Println("Procesando comandos base de datos")

	var cacheProcess beans.CacheProcess

	// Procesamos los comandos
	for _, command := range configuration.ConfigurationDatabase.Commands {
		log.Printf("Procesando tipo de comando %d", command)

		switch command {
		case constants.Backup:
		case constants.Migrate:

			// Connectamos contra la base de datos
			connectDatabase(cacheProcess, configuration)

			// Miramos la existencia del directorio de migraciones
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
						processModulesMigratios(cacheProcess, configuration, files)
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
//
// @parameters cacheProcess cache donde guardar ciertos datos del procesado
//
// @parameter configuration configuración que tiene todos los parámetros de configuración y comandos a procesar
//
// @parameter files array de información de ficheros que serán los módulos a migrar
//
// @returns ---
func processModulesMigratios(cacheProcess beans.CacheProcess, configuration beans.Configuration, files []fs.FileInfo) {
	var pathMigrations string = configuration.ConfigurationDatabase.PathMigrations
	for _, module := range files {
		// En el caso de que el módulo sea un directorio listamos sus versiones
		if module.IsDir() {
			log.Printf("Procesando módulo a migrar %s", module.Name())

			filesModule, err := ioutil.ReadDir(filepath.Join(pathMigrations, module.Name()))

			if err != nil {
				log.Panicf("Se ha producido un error al listar el directorio del modulo %s. Error: %s", module.Name(), err)
			} else {
				// Ordenamos por version
				sort.Slice(filesModule, func(first, second int) bool {
					firstFile, errFirstFile := strconv.Atoi(filesModule[first].Name())
					secondFIle, errSecondFile := strconv.Atoi(filesModule[second].Name())

					if errFirstFile != nil {
						log.Panicf("Se ha producido un error al ordenar los directorios del módulo. Direcotrio: %s. Error %s", files[first].Name(), errFirstFile)
					}

					if errSecondFile != nil {
						log.Panicf("Se ha producido un error al ordenar los directorios del módulo. Direcotrio: %s. Error %s", files[second].Name(), errFirstFile)
					}

					return errFirstFile != nil && errSecondFile != nil && firstFile < secondFIle
				})
			}

			// Recorremos los ficheros del módulo que serán las versiones
			for _, fileVersionModule := range filesModule {
				log.Printf("Procesando versión a migrar %s del módulo %s", fileVersionModule.Name(), module.Name())

				// Guardamos el path de la versión para listar los ficheros
				pathVersion := filepath.Join(pathMigrations, module.Name(), fileVersionModule.Name())

				// Procesamos la versión
				processVersion(cacheProcess, configuration, fileVersionModule, module.Name(), pathVersion)
			}

		} else {
			log.Printf("Módulo a migrar %s no es un directorio", module.Name())
		}
	}
}

// processVersion: Método para procesar la versión
//
// @parameters cacheProcess cache donde guardar ciertos datos del procesado
//
// @parameter configuration configuración que tiene todos los parámetros de configuración y comandos a procesar
//
// @parameter fileVersion fichero de la versión que será un directorio
//
// @parameter moduleName es el nombre del módulo de la versión
//
// @paremeter pathVersion es la ruta donde se encuentra la versión para listar los ficheros
//
// @returns ---
func processVersion(cacheProcess beans.CacheProcess, configuration beans.Configuration, fileVersion fs.FileInfo, moduleName string, pathVersion string) {

}

// connectDatabase. Métoodo para connectar con la base de datos
//
// @parameters cacheProcess cache donde guardar ciertos datos del procesado
//
// @parameter configuration configuración que tiene todos los parámetros de configuración y comandos a procesar
//
// @returns --
func connectDatabase(cacheProcess beans.CacheProcess, configuration beans.Configuration) {
	switch configuration.ConfigurationDatabase.Type {
	case constants.Mysql:
		ConnectMysqlDatabase(cacheProcess, configuration)
	}
}
