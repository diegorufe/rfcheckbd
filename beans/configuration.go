package beans

// Configuration estructura donde guardar la configuración pasada en un fichero json al arrancar la aplicación
type Configuration struct {
	AppName               string                `json:"AppName"`               // Indicamos el nombre de la aplicación en el caso de no decirle nada se llamará rfcheckbd
	ConfigurationLog      ConfigurationLog      `json:"ConfigurationLog"`      // Datos para la configuración del log
	ConfigurationDatabase ConfigurationDatabase `json:"ConfigurationDatabase"` // Datos para la configuración de la base de datos
}
