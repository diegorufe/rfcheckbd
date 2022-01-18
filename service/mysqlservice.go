package service

import (
	"container/list"
	"context"
	"database/sql"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"rfcheckbd/beans"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// MysqlDatabaseService servicio para definir la funcionalidad de backup, ejercución de scripts etc ... para mysql
type MysqlDatabaseService struct {
}

// errorIsNoDataFound: Método para saber si un error es un error de no datos encontrados en el scan
//
// @parameter err error a verificar
//
// @returns true si es not null y no es error de no encontrados registros
func errorIsNoDataFound(err error) bool {
	return err != nil && err.Error() == "sql: no rows in result set"
}

func (service MysqlDatabaseService) ConnectDatabase(cacheProcess beans.CacheProcess, configuration beans.Configuration) beans.CacheProcess {
	log.Println("Connectando con la base de datos de tipo mysql")

	// TODO pedir los datos por línea de comandos

	strConnection := configuration.ConfigurationDatabase.User + ":" + configuration.ConfigurationDatabase.Password + "@tcp(" + configuration.ConfigurationDatabase.Host + ":" + strconv.FormatInt(configuration.ConfigurationDatabase.Port, 10) + ")/" + configuration.ConfigurationDatabase.DatabaseName + "?multiStatements=true&interpolateParams=true"

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

	return cacheProcess
}

func (service MysqlDatabaseService) FindVersionModule(cacheProcess beans.CacheProcess, configuration beans.Configuration, moduleName string) int64 {
	var version int64

	query := "SELECT version from rfcheckbd_migrations where module = ?"

	cacheProcess.DBSql.QueryRow(query, moduleName)

	err := cacheProcess.DBSql.QueryRow(query, moduleName).Scan(&version)

	if err != nil && !errorIsNoDataFound(err) {
		log.Panicf("Se ha produdio un error a la hora de buscar la versión del módulo. %s", err)
	}

	return version
}

func (service MysqlDatabaseService) ProcessFileInVersion(cacheProcess beans.CacheProcess, configuration beans.Configuration, moduleName string, version string, versionInt int, pathVersion string, fileInVersion fs.FileInfo) {
	// Tenemos que buscar si está en el los registros de ficheros ejecutados
	var id int64

	queryFindFile := "SELECT h.id from rfcheckbd_migrations_history h INNER JOIN rfcheckbd_migrations m ON h.rfcheckbd_migrations_id = m.id  where h.fileName = ? and m.module = ?"

	err := cacheProcess.DBSql.QueryRow(queryFindFile, fileInVersion.Name(), moduleName).Scan(&id)

	if err != nil && !errorIsNoDataFound(err) {
		log.Panicf("Se ha produdio un error a la hora de buscar la el fichero de la versión %s para el módulo %s. Error %s", version, moduleName, err)
	}

	// En el caso de que no exista lo proceso
	if id <= 0 {

		dataFile, errReadFile := os.ReadFile(filepath.Join(pathVersion, fileInVersion.Name()))

		if errReadFile != nil {
			log.Panicf("Se ha producido un error leer el fichero %s para la versión %s del módulo %s. Error %s", fileInVersion.Name(), version, moduleName, errReadFile)
		}

		// Datos del fichero a string
		queryProcessFile := string(dataFile)

		var listParams list.List

		execDate := time.Now().Format("2006-01-02 03:04:05")

		// En el caso de que la versión sea menor o igual a cero la insertamos inicialmente
		if versionInt <= 0 {
			queryInsertMigrataions := "INSERT IGNORE INTO rfcheckbd_migrations (`version`, `module`, `execDate`) " +
				" VALUES ( " +
				" ?, " +
				" ?, " +
				" ? );"

			listParams.PushBack(versionInt)
			listParams.PushBack(moduleName)
			listParams.PushBack(execDate)

			queryProcessFile = queryProcessFile + " " + queryInsertMigrataions
		}

		// guardamos el fichero en el historico
		queryInsertInHistory := "INSERT INTO rfcheckbd_migrations_history (`rfcheckbd_migrations_id`, `fileName` , `execDate` ) " +
			" VALUES ( " +
			" (SELECT m.id FROM  rfcheckbd_migrations m where module = ? ), " +
			" ?, " +
			" ? );"

		listParams.PushBack(moduleName)
		listParams.PushBack(fileInVersion.Name())
		listParams.PushBack(execDate)

		queryProcessFile = queryProcessFile + " " + queryInsertInHistory

		// Query para actualizar la versión siempre por que puede que el insert de la migración ya existiera antes y por lo tanto no realizó

		queryUpdateVersion := "UPDATE rfcheckbd_migrations " +
			"SET " +
			" `version`  = ?, " +
			" `execDate` = ? " +
			"WHERE `module` = ? ; "

		queryProcessFile = queryProcessFile + " " + queryUpdateVersion

		listParams.PushBack(versionInt)
		listParams.PushBack(execDate)
		listParams.PushBack(moduleName)

		// Paso la lista de parámetros a un array
		arrayParams := make([]interface{}, listParams.Len())
		var counter uint64 = 0

		for element := listParams.Front(); element != nil; element = element.Next() {
			// siempre utilziar el value del elemento "element.Value"
			arrayParams[counter] = element.Value
			counter = counter + 1
		}

		// Abrimos transacción
		ctx := context.Background()
		tx, errTx := cacheProcess.DBSql.BeginTx(ctx, nil)

		if errTx != nil {
			log.Panicf("Se ha producido un error al abrir una transación para ejecutar el contenido del fichero %s para la versión %s del módulo %s. Error %s", fileInVersion.Name(), version, moduleName, errTx)
		}

		// Ejecutamos el proceso del fichero
		_, errProcessFile := tx.ExecContext(ctx, queryProcessFile, arrayParams...)

		if errProcessFile != nil {

			errorRollback := tx.Rollback()

			if errorRollback != nil {
				log.Panicf("Se ha producido un error al realizar el rollback de la transación al ejecutar el contenido del fichero %s para la versión %s del módulo %s. Error %s. Error procesado fichero %s", fileInVersion.Name(), version, moduleName, errorRollback, errProcessFile)
			}

			log.Panicf("Se ha producido un error al ejecutar el contenido del fichero %s para la versión %s del módulo %s. Error %s", fileInVersion.Name(), version, moduleName, errProcessFile)

		}

		errCommit := tx.Commit()

		if errCommit != nil {
			log.Panicf("Se ha producido un error al realizar el commit de la transación al ejecutar el contenido del fichero %s para la versión %s del módulo %s. Error %s", fileInVersion.Name(), version, moduleName, errCommit)
		}

		log.Printf("Fichero %s procesado con éxito para la versión %s y módulo %s", fileInVersion.Name(), version, pathVersion)

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

	query := "CREATE TABLE IF NOT EXISTS `rfcheckbd_migrations` (" +
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

	query := "CREATE TABLE IF NOT EXISTS `rfcheckbd_migrations_history` (" +
		"	`id` int(11) NOT NULL auto_increment,   " +
		"	`rfcheckbd_migrations_id` int(11) NOT NULL COMMENT 'Módulo de la migración',     " +
		"	`fileName` varchar(500)  NOT NULL COMMENT 'Nombre del fichero que se ejecuto' ,  " +
		"	`execDate` DATETIME  NOT NULL COMMENT 'Fecha de ejecución del fichero' , " +
		"	UNIQUE KEY `module_UNIQUE` (`rfcheckbd_migrations_id`, `fileName`),  " +
		"	 PRIMARY KEY  (`id`),  " +
		"    FOREIGN KEY (rfcheckbd_migrations_id) REFERENCES rfcheckbd_migrations(id) " +
		"  ) COMMENT 'Tabla que contiene el histórico de ficheros pasados en cada migración';"

	_, err := cacheProcess.DBSql.Exec(query)

	if err != nil {
		log.Panicf("Se ha produdio un error en la creación de la tabla de historico de migraciones. %s", err)
	}
}
