package config

import (
	"encoding/json"
	"log"
	"os"
	"rfcheckbd/beans"
)

// Config Método para realizar la configuración de la aplicación
func Config() {

	// Flags del log para dejarlo "BONITO"
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile)

	configuration, err := loadConfiguration("")

	if err == nil {

		// Ponemos ciertos datos por defecto de la aplicación
		loadDefaultConfig(&configuration)

		// Configuramos los logs
		ConfigLog(configuration)

	} else {
		log.Panic("No se ha podido cargar la configuración de la aplicación")
	}
}

// loadDefaultConfig: Método para cargar ciertos datos de la configuración por defecto
//
// @param configuration configuración a la cual realizar ciertas operaciones por defecto
//
// @retuns --
func loadDefaultConfig(configuration *beans.Configuration) {

	configuration.ConfigurationLog.Compress = true

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
		configuration.ConfigurationLog.PathFileLogging = "rfcheckbd.log"
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

	file, err := os.Open(pathJsonFile)

	if err == nil {
		defer file.Close()
		decoder := json.NewDecoder(file)
		configuration := beans.Configuration{}

		err = decoder.Decode(&configuration)
	} else {
		err = nil
		configuration = beans.Configuration{}
	}

	return configuration, err
}
