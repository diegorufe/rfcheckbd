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
// @parameter execDate fecha de ejecucion en formato string
//
// @returns ---
func ProccessDatabseCommands(configuration beans.Configuration, execDate string) {
	log.Println("Procesando comandos base de datos")

	var cacheProcess beans.CacheProcess
	var databaseService IDabaseService

	// Nos quedamos con la fecha de ejercicón en la cache de los procesos
	cacheProcess.ExecDate = execDate

	switch configuration.ConfigurationDatabase.Type {
	case constants.Mysql:
		databaseService = MysqlDatabaseService{}
	}

	// Procesamos los comandos
	for _, command := range configuration.ConfigurationDatabase.Commands {
		log.Printf("Procesando tipo de comando %d", command)

		switch command {
		case constants.Backup:
		case constants.Migrate:

			// Connectamos contra la base de datos
			cacheProcess = databaseService.ConnectDatabase(cacheProcess, configuration)
			// Nos quedamos con la fecha de ejercicón en la cache de los procesos
			cacheProcess.ExecDate = execDate

			// Miramos la existencia del directorio de migraciones
			stat, err := os.Stat(configuration.ConfigurationDatabase.PathMigrations)

			// Comrpobamos que el directorio de módulos exista
			if err != nil {
				log.Panicf("Se ha producido un error al comprobar el path de la migración %s", err)
			}

			// En el caso de ser una carpeta continuamos leyendo el directorio
			if !stat.IsDir() {
				log.Panicf("La ruta donde se encuentran los módulos de migraciones no es un directorio. Ruta de módulos de migraciones %s", stat.Name())
			}

			files, err := ioutil.ReadDir(configuration.ConfigurationDatabase.PathMigrations)

			if err != nil {
				log.Panicf("Se ha producido un error al listar el directorio de migraciones %s", err)
			}

			processModulesMigrations(cacheProcess, configuration, databaseService, files)

		}
	}
}

// processModulesMigrations. Método para procesar los módulos de las migraciones
//
// @parameters cacheProcess cache donde guardar ciertos datos del procesado
//
// @parameter configuration configuración que tiene todos los parámetros de configuración y comandos a procesar
//
// @parameter databaseService servicio para realizar la lógica funcional de la base de datos
//
// @parameter files array de información de ficheros que serán los módulos a migrar
//
// @returns ---
func processModulesMigrations(cacheProcess beans.CacheProcess, configuration beans.Configuration, databaseService IDabaseService, files []fs.FileInfo) {
	var pathMigrations string = configuration.ConfigurationDatabase.PathMigrations
	for _, module := range files {
		// En el caso de que el módulo sea un directorio listamos sus versiones
		if module.IsDir() {
			log.Printf("Procesando módulo a migrar %s", module.Name())

			filesModule, err := ioutil.ReadDir(filepath.Join(pathMigrations, module.Name()))

			if err != nil {
				log.Panicf("Se ha producido un error al listar el directorio del modulo %s. Error: %s", module.Name(), err)
			}

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

			// Recorremos los ficheros del módulo que serán las versiones
			for _, fileVersionModule := range filesModule {
				log.Printf("Procesando versión a migrar %s del módulo %s", fileVersionModule.Name(), module.Name())

				// Guardamos el path de la versión para listar los ficheros
				pathVersion := filepath.Join(pathMigrations, module.Name(), fileVersionModule.Name())

				// Procesamos la versión
				processVersion(cacheProcess, configuration, databaseService, fileVersionModule, module.Name(), pathVersion)
			}

		} else {
			log.Printf("Módulo a migrar %s no es un directorio", module.Name())
		}
	}
}

// processVersion: Método para procesar la versión. Ejemplo versión 10 del módulo de contabiliad
//
// @parameter cacheProcess cache donde guardar ciertos datos del procesado
//
// @parameter configuration configuración que tiene todos los parámetros de configuración y comandos a procesar
//
// @parameter databaseService servicio para realizar la lógica funcional de la base de datos
//
// @parameter fileVersion fichero de la versión que será un directorio
//
// @parameter moduleName es el nombre del módulo de la versión
//
// @paremeter pathVersion es la ruta donde se encuentra la versión para listar los ficheros
//
// @returns ---
func processVersion(cacheProcess beans.CacheProcess, configuration beans.Configuration, databaseService IDabaseService, fileVersion fs.FileInfo, moduleName string, pathVersion string) {
	log.Printf("Procesando versión %s del módulo %s", fileVersion.Name(), moduleName)

	versionToProcess, err := strconv.Atoi(fileVersion.Name())

	if err != nil {
		log.Panicf("Error al convertir a entero la versión del módulo. %s", err)
	}

	// Buscamos la versión del módulo
	cacheProcess.VersionModule = int(databaseService.FindVersionModule(cacheProcess, configuration, moduleName))

	if versionToProcess >= cacheProcess.VersionModule {
		// Leemos los ficheros de la versión
		filesVersion, err := ioutil.ReadDir(pathVersion)

		if err != nil {
			log.Panicf("Se ha producido un error al buscar los ficheros de la versión: %s para el módulo %s. Error %s", fileVersion.Name(), moduleName, err)
		}

		// Ordenamos los ficheros por nombre
		sort.Slice(filesVersion, func(first, second int) bool {
			firstFile := filesVersion[first].Name()
			secondFIle := filesVersion[second].Name()

			return firstFile < secondFIle
		})

		// Recorremos los ficheros
		for _, fileInVersion := range filesVersion {
			log.Printf("Procesando fichero %s de la versión %s del módulo %s", fileInVersion.Name(), fileVersion.Name(), moduleName)

			// Procesamos el ficheros
			databaseService.ProcessFileInVersion(cacheProcess, configuration, moduleName, fileVersion.Name(), versionToProcess, pathVersion, fileInVersion)
		}

	} else {
		log.Printf("No se procesará la versión %s del módulo %s por que es inferior a la versión actual %s", fileVersion.Name(), moduleName, strconv.Itoa(cacheProcess.VersionModule))
	}
}
