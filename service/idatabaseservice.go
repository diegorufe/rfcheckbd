package service

import (
	"io/fs"
	"rfcheckbd/beans"
)

// IDabaseService defincición base para los servicios
type IDabaseService interface {
	// ConnectDatabase. Métoodo para connectar con la base de datos
	//
	// @parameters cacheProcess cache donde guardar ciertos datos del procesado
	//
	// @parameter configuration configuración que tiene todos los parámetros de configuración y comandos a procesar
	//
	// @returns --
	ConnectDatabase(cacheProcess beans.CacheProcess, configuration beans.Configuration) beans.CacheProcess

	// FindVersionModuleMysql : Método apra buscar la versión del módulo para mysql
	//
	// @parameters cacheProcess cache donde guardar ciertos datos del procesado
	//
	// @parameter configuration configuración que tiene todos los parámetros de configuración y comandos a procesar
	//
	// @parameter moduleName es el nombre del módulo
	//
	// @returns --
	FindVersionModule(cacheProcess beans.CacheProcess, configuration beans.Configuration, moduleName string) int64

	// ProcessFileInVersion: Método para procesar el fichero de la versión. En caso de no estar procesado ya el fichero lo procesa y actualiza la versión para el módulo
	//
	// @parameters cacheProcess cache donde guardar ciertos datos del procesado
	//
	// @parameter configuration configuración que tiene todos los parámetros de configuración y comandos a procesar
	//
	// @parameter moduleName es el nombre del módulo
	//
	// @parameter version es la versión a ejecutar
	//
	// @parameter versionInt es la versión ya convertida a entero
	//
	// @parameter pathVersion ruta donde están los ficheros de la versión
	//
	// @parameter fileInVersion es el fichero a procesar
	//
	// @returns --
	ProcessFileInVersion(cacheProcess beans.CacheProcess, configuration beans.Configuration, moduleName string, version string, versionInt int, pathVersion string, fileInVersion fs.FileInfo)
}
