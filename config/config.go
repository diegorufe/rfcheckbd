package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"rfcheckbd/beans"
	"strings"
)

// Config Método para realizar la configuración de la aplicación
//
// @parameter execDate fecha de ejecucion en formato string
//
// @returns --
func Config(execDate string) beans.Configuration {

	// Flags del log para dejarlo "BONITO"
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile)

	configuration, err := loadConfiguration("")

	if err == nil {

		// Ponemos ciertos datos por defecto de la aplicación
		loadDefaultConfig(&configuration)

		// Configuramos los logs
		ConfigLog(configuration, execDate)

	} else {
		log.Panic("No se ha podido cargar la configuración de la aplicación")
	}

	return configuration
}

// loadDefaultConfig: Método para cargar ciertos datos de la configuración por defecto
//
// @param configuration configuración a la cual realizar ciertas operaciones por defecto
//
// @retuns --
func loadDefaultConfig(configuration *beans.Configuration) {

	// Configuración por defecto de log
	loadDefaultConfigLog(configuration)

	// Configuración por defecto para la base datos
	loadDefaultConfigDatabase(configuration)
}

// loadDefaultConfigLog configuración por defecto para log
//
// @parameters configuration  configuración a la cual realizar ciertas operaciones por defecto
//
// @retuns --
func loadDefaultConfigLog(configuration *beans.Configuration) {
	configuration.ConfigurationLog.Compress = true

	// Si el nombre de la aplicación está vacío ponemos uno por defecto
	if configuration.AppName == "" {
		configuration.AppName = "rfcheckbd"
	}

	// Si no se pasa tiempo máximo lo ponemos a 30
	if configuration.ConfigurationLog.MaxAge <= 0 {
		configuration.ConfigurationLog.MaxAge = 30
	}

	// Si no tenemos un númeor máximo de backups lo dejamos a 5
	if configuration.ConfigurationLog.MaxBackups <= 0 {
		configuration.ConfigurationLog.MaxBackups = 5
	}

	// Si no se pasa tamaño máximo lo ponemos a 500
	if configuration.ConfigurationLog.MaxSize <= 0 {
		configuration.ConfigurationLog.MaxSize = 500
	}

	// En el caso de no poner una ruta de fichero se pone la ruta del ejecutable y se pone el nombre del proyecto como fichero de log
	if configuration.ConfigurationLog.PathFileLogging == "" {
		configuration.ConfigurationLog.PathFileLogging = configuration.AppName + ".log"
	}
}

// loadDefaultConfigDatabase configuración por defecto para la base datos
//
// @parameters configuration  configuración a la cual realizar ciertas operaciones por defecto
//
// @retuns --
func loadDefaultConfigDatabase(configuration *beans.Configuration) {

	// En el caso de peder usuario y contraseña de la base de datos tenemos que pedirlo por consola
	if configuration.ConfigurationDatabase.AskForUserPassword {
		reader := bufio.NewReader(os.Stdin)

		fmt.Println("Introduce usuario de la base de datos: ")
		user, err := reader.ReadString('\n')

		if err != nil {
			log.Panicf("No se ha podido procesar el usuario de la base de datos. Error %s", err)
		}

		user = strings.TrimRight(user, "\r\n")

		configuration.ConfigurationDatabase.User = user

		fmt.Println("Introduce contraseña de la base de datos: ")
		password, err := reader.ReadString('\n')

		if err != nil {
			log.Panicf("No se ha podido procesar la contraseña de la base de datos. Error %s", err)
		}

		password = strings.TrimRight(password, "\r\n")

		configuration.ConfigurationDatabase.Password = password

	}
}

// loadConfiguration : Método para cargar la configuración
//
// @params pathJsonFile: Ruta donde está el fichero de configuración
//
// @retuns
//
// - La configuración en caso de poder cargarla
//
// - Error en el caso de no poder cargar el json de la configuración
func loadConfiguration(pathJsonFile string) (beans.Configuration, error) {
	// Cargamos la configuración de las url a interceptar y el host de destino
	var configuration beans.Configuration
	var err error

	if pathJsonFile == "" {
		pathJsonFile = "rfcheckbd.json"
	}

	file, err := os.Open(pathJsonFile)

	if err == nil {
		defer file.Close()
		decoder := json.NewDecoder(file)

		err = decoder.Decode(&configuration)
	} else {
		err = nil
		configuration = beans.Configuration{}
	}

	return configuration, err
}
