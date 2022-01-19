package config

import (
	"io"
	"log"
	"os"
	"rfcheckbd/beans"
	"strings"

	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

// ConfigLog: Método para configurar el log de la aplicación
//
// @param configuration configuración pasada para poder instanciar el log
//
// @parameter execDate fecha de ejecucion en formato string
//
// @returns ---
func ConfigLog(configuration beans.Configuration, excecDate string) {

	log.Println("*** Configurando logs inicio ***")

	log.Printf("Configuración del log. FileName: %s | MaxSize: %d | MaxBackups: %d | MaxAge: %d | Compress: %t", configuration.ConfigurationLog.PathFileLogging, configuration.ConfigurationLog.MaxSize, configuration.ConfigurationLog.MaxBackups, configuration.ConfigurationLog.MaxAge, configuration.ConfigurationLog.Compress)

	var pahtFileLogging string = configuration.ConfigurationLog.PathFileLogging

	// En el caso de que queramos log por prceso quito la parte de log y le añado la fecha de ejercución sin puntos
	if configuration.ConfigurationLog.LogPerProcess {
		pahtFileLogging = strings.ReplaceAll(pahtFileLogging, ".log", "")

		pahtFileLogging = pahtFileLogging + "_" + strings.ReplaceAll(excecDate, ":", "") + ".log"
	}

	multiWritter := io.MultiWriter(&lumberjack.Logger{
		Filename:   pahtFileLogging,                           // ruta del fichero
		MaxSize:    configuration.ConfigurationLog.MaxSize,    // megabytes
		MaxBackups: configuration.ConfigurationLog.MaxBackups, // Máximo de backups
		MaxAge:     configuration.ConfigurationLog.MaxAge,     // días que se tendrán los ficheros
		Compress:   configuration.ConfigurationLog.Compress,   // Compresión deshabilitado por defecto
	}, os.Stdout)

	log.SetOutput(multiWritter)

	log.Println("*** Configurando logs fin ***")

}
