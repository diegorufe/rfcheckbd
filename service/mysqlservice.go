package service

import (
	"database/sql"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"rfcheckbd/beans"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

// MysqlDatabaseService servicio para definir la funcionalidad de backup, ejercución de scripts etc ... para mysql
type MysqlDatabaseService struct {
}

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

func (service MysqlDatabaseService) FindVersionModule(cacheProcess beans.CacheProcess, configuration beans.Configuration, moduleName string) int64 {
	var version int64

	query := "SELECT version from rfchecbd_migrations where module = %s"

	err := cacheProcess.DBSql.QueryRow(query, moduleName).Scan(&version)

	if err != nil {
		log.Panicf("Se ha produdio un error a la hora de buscar la versión del módulo. %s", err)
	}

	return version
}

func (service MysqlDatabaseService) ProcessFileInVersion(cacheProcess beans.CacheProcess, configuration beans.Configuration, moduleName string, version string, versionInt int, pathVersion string, fileInVersion fs.FileInfo) {
	// Tenemos que buscar si está en el los registros de ficheros ejecutados
	var id int64

	queryFindFile := "SELECT h.id from rfchecbd_migrations_history h INNER JOIN rfchecbd_migrations m ON h.rfchecbd_migrations_id = m.id  where h.fileName = %s and m.module = %s"

	err := cacheProcess.DBSql.QueryRow(queryFindFile, moduleName).Scan(&id)

	if err != nil {
		log.Panicf("Se ha produdio un error a la hora de buscar la el fichero de la versión %s para el módulo %s. Error %s", version, moduleName, err)
	}

	// En el caso de que no exista lo proceso
	if id <= 0 {

		dataFile, errReadFile := os.ReadFile(filepath.Join(pathVersion, fileInVersion.Name()))

		if errReadFile != nil {
			log.Panicf("Se ha producido un error leer el fichero %s para la versión %s del módulo %s. Error %s", fileInVersion.Name(), version, moduleName, errReadFile)
		}

		// Datos del fichero a string
		queryProcessFIle := string(dataFile)

		// guardamos el fichero en el historico
		queryInsertInHistory := "INSERT INTO rfchecbd_migrations_history (`rfchecbd_migrations_id`, `fileName`, `execDate` ) " +
			" VALUES ( " +
			" (SELECT m.id FROM  rfchecbd_migrations m where module = %s), " +
			" %s, " +
			" %s);"

		queryProcessFIle = queryProcessFIle + " " + queryInsertInHistory

		// En el caso de que la versión sea menor o igual a cero la insertamos inicialmente
		if versionInt <= 0 {
			// TODO
			_ = "INSERT IGNORE INTO rfchecbd_migrations (`version`, `module`, `execDate`) " +
				" VALUES ( " +
				" %d, " +
				" %s, " +
				" %s); "
		}

		// Query para actualizar la versión
		// TODO
		_ = "UPDATE rfchecbd_migrations " +
			"SET " +
			" `version` = %d, " +
			" `execDate` = %s " +
			"WHERE `module` = %s; "

		// Ejecutamos el proceso del fichero
		// TODO revisar paráemtros de tipo date
		_, errProcessFile := cacheProcess.DBSql.Exec(queryProcessFIle, moduleName, fileInVersion.Name(), "")

		if errProcessFile != nil {
			log.Panicf("Se ha producido un error leer el procesar el fichero %s para la versión %s del módulo %s. Error %s", fileInVersion.Name(), version, moduleName, errProcessFile)
		}

	} else {
		log.Printf("Fichero %s ya procesado para la versión %s y módulo %s", fileInVersion.Name(), version, pathVersion)
	}
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
