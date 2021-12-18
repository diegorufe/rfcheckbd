package beans

// ConfigurationLog estructura para la configuración del log
type ConfigurationLog struct {
	PathFileLogging string // Ruta donde se guardará el fichero de loggin
	MaxSize         int    // Tamaño máximo del fichero en megabites
	MaxBackups      int    // Número máximo de backups
	MaxAge          int    // Número de días que se tendran los ficheros de log
	Compress        bool   // Indicamos si queremos que se comprima o no
}
