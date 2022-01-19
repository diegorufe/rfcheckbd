package beans

// ConfigurationLog estructura para la configuración del log
type ConfigurationLog struct {
	PathFileLogging string `json:"PathFileLogging"` // Ruta donde se guardará el fichero de loggin
	MaxSize         int    `json:"MaxSize"`         // Tamaño máximo del fichero en megabites
	MaxBackups      int    `json:"MaxBackups"`      // Número máximo de backups
	MaxAge          int    `json:"MaxAge"`          // Número de días que se tendran los ficheros de log
	Compress        bool   `json:"Compress"`        // Indicamos si queremos que se comprima o no
	LogPerProcess   bool   `json:"LogPerProcess"`   // Si es true indicamos que se tiene que generar un fichero de log por cada vez que se ejecute la aplicación
}
