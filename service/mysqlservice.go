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
}
