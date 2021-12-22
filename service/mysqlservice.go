package service

import (
	"database/sql"
	"log"
	"rfcheckbd/beans"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

// ConnectMysqlDatabase. Métoodo para connectar con la base de datos
//
// @parameters cacheProcess cache donde guardar ciertos datos del procesado
//
// @parameter configuration configuración que tiene todos los parámetros de configuración y comandos a procesar
//
// @returns --
func ConnectMysqlDatabase(cacheProcess beans.CacheProcess, configuration beans.Configuration) {
	log.Println("Connectando con la base de datos de tipo mysql")

	// TODO pedir los datos por línea de comandos

	strConnection := configuration.ConfigurationDatabase.User + ":" + configuration.ConfigurationDatabase.Password + "@tcp(" + configuration.ConfigurationDatabase.Host + ":" + strconv.FormatInt(configuration.ConfigurationDatabase.Port, 10) + ")/" + configuration.ConfigurationDatabase.DatabaseName

	db, err := sql.Open("mysql", strConnection)

	// if there is an error opening the connection, handle it
	if err != nil {
		log.Panicf("Se ha produdio un error en la conexión de la base de datos. %s", err)
	}

	// Nos quedamos con la conexión de la base de datos
	cacheProcess.DBSql = db

	// Creamos versionado de la tabla
	createVersionTable(cacheProcess, configuration)
}

// createVersionTable Método para crear el versioando de la tabla
//
// @parameters cacheProcess cache donde guardar ciertos datos del procesado
//
// @parameter configuration configuración que tiene todos los parámetros de configuración y comandos a procesar
//
// @returns --
func createVersionTable(cacheProcess beans.CacheProcess, configuration beans.Configuration) {

	query := "CREATE TABLE IF NOT EXISTS `rfchecbd_migrations` (" +
		"	`id` int(11) NOT NULL auto_increment,   " +
		"	`version` int(11) NOT NULL COMMENT 'Versión del módulo',     " +
		"	`module` varchar(250)  NOT NULL COMMENT 'Nombre del módulo' ,  " +
		"	`execDate` DATETIME  NOT NULL COMMENT 'Fecha de última ejecución del módulo' , " +
		"	UNIQUE KEY `module_UNIQUE` (`module`),  " +
		"	 PRIMARY KEY  (`id`)  " +
		"  );"

	_, err := cacheProcess.DBSql.Exec(query)

	if err != nil {
		log.Panicf("Se ha produdio un error en la creación de la tabla de versionado de la migración. %s", err)
	}
}
