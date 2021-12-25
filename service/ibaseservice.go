package service

import "rfcheckbd/beans"

// IDabaseService defincición base para los servicios
type IDabaseService interface {
	// ConnectDatabase. Métoodo para connectar con la base de datos
	//
	// @parameters cacheProcess cache donde guardar ciertos datos del procesado
	//
	// @parameter configuration configuración que tiene todos los parámetros de configuración y comandos a procesar
	//
	// @returns --
	ConnectDatabase(cacheProcess beans.CacheProcess, configuration beans.Configuration)

	// FindVersionModuleMysql : Método apra buscar la versión del módulo para mysql
	//
	// @parameters cacheProcess cache donde guardar ciertos datos del procesado
	//
	// @parameter configuration configuración que tiene todos los parámetros de configuración y comandos a procesar
	//
	//
	// @parameter moduleName es el nombre del módulo
	//
	// @returns --
	FindVersionModule(cacheProcess beans.CacheProcess, configuration beans.Configuration, moduleName string) int64
}
