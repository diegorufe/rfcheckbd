package service

import (
	"database/sql"
	"log"
	"rfcheckbd/beans"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

// MysqlDatabaseService servicio para definir la funcionalidad de backup, ejercución de scripts etc ... para mysql
type MysqlDatabaseService struct {
}

// ConnectMysqlDatabase. Métoodo para connectar con la base de datos
//
// @parameters cacheProcess cache donde guardar ciertos datos del procesado
//
// @parameter configuration configuración que tiene todos los parámetros de configuración y comandos a procesar
//
// @returns --
func (service MysqlDatabaseService) ConnectDatabase(cacheProcess beans.CacheProcess, configuration beans.Configuration) {
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

	// Creamos tabla de historico de migraciones
	createHistoryTable(cacheProcess, configuration)
}

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
func (service MysqlDatabaseService) FindVersionModule(cacheProcess beans.CacheProcess, configuration beans.Configuration, moduleName string) int64 {
	var version int64

	query := "SELECT version from rfchecbd_migrations where module = %s"

	err := cacheProcess.DBSql.QueryRow(query, moduleName).Scan(&version)

	if err != nil {
		log.Panicf("Se ha produdio un error a la hora de buscar la versión del módulo. %s", err)
	}

	return version
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
		"  ) COMMENT 'Tabla que contiene las versiones de las migraciones realizdadas';"

	_, err := cacheProcess.DBSql.Exec(query)

	if err != nil {
		log.Panicf("Se ha produdio un error en la creación de la tabla de versionado de la migración. %s", err)
	}
}

// createVersionTable Método para crear el versioando de la tabla
//
// @parameters cacheProcess cache donde guardar ciertos datos del procesado
//
// @parameter configuration configuración que tiene todos los parámetros de configuración y comandos a procesar
//
// @returns --
func createHistoryTable(cacheProcess beans.CacheProcess, configuration beans.Configuration) {

	query := "CREATE TABLE IF NOT EXISTS `rfchecbd_migrations_history` (" +
		"	`id` int(11) NOT NULL auto_increment,   " +
		"	`rfchecbd_migrations_id` int(11) NOT NULL COMMENT 'Módulo de la migración',     " +
		"	`fileName` varchar(500)  NOT NULL COMMENT 'Nombre del fichero que se ejecuto' ,  " +
		"	`execDate` DATETIME  NOT NULL COMMENT 'Fecha de ejecución del fichero' , " +
		"	UNIQUE KEY `module_UNIQUE` (`rfchecbd_migrations_id`, `fileName`),  " +
		"	 PRIMARY KEY  (`id`),  " +
		"    FOREIGN KEY (rfchecbd_migrations_id) REFERENCES rfchecbd_migrations(id) " +
		"  ) COMMENT 'Tabla que contiene el histórico de ficheros pasados en cada migración';"

	_, err := cacheProcess.DBSql.Exec(query)

	if err != nil {
		log.Panicf("Se ha produdio un error en la creación de la tabla de historico de migraciones. %s", err)
	}
}
