package config

import (
	"io"
	"log"
	"os"
	"rfcheckbd/beans"

	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

// ConfigLog: Método para configurar el log de la aplicación
//
// @param configuration configuración pasada para poder instanciar el log
func ConfigLog(configuration beans.Configuration) {

	log.Println("*** Configurando logs inicio ***")

	log.Printf("Configuración del log. FileName: %s | MaxSize: %d | MaxBackups: %d | MaxAge: %d | Compress: %t", configuration.ConfigurationLog.PathFileLogging, configuration.ConfigurationLog.MaxSize, configuration.ConfigurationLog.MaxBackups, configuration.ConfigurationLog.MaxAge, configuration.ConfigurationLog.Compress)

	multiWritter := io.MultiWriter(&lumberjack.Logger{
		Filename:   configuration.ConfigurationLog.PathFileLogging, // ruta del fichero
		MaxSize:    configuration.ConfigurationLog.MaxSize,         // megabytes
		MaxBackups: configuration.ConfigurationLog.MaxBackups,      // Máximo de backups
		MaxAge:     configuration.ConfigurationLog.MaxAge,          // días que se tendrán los ficheros
		Compress:   configuration.ConfigurationLog.Compress,        // Compresión deshabilitado por defecto
	}, os.Stdout)

	log.SetOutput(multiWritter)

	log.Println("*** Configurando logs fin ***")

}
