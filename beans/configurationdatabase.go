package beans

import (
	"rfcheckbd/constants"
)

// ConfigurationDatabase para indicar la configuración de la base datos
type ConfigurationDatabase struct {
	Type           constants.DatabaseType            // Tipo de base de datos mysql, posgrees etc ...
	PathSaveBackup string                            // Ruta donde guardar el fichero de backup. Tiene que ser un directorio existente
	Commands       []constants.DatabaseCommandsTypes // Comandos a ejecutar
	User           string                            // Usuario para realizar la conexión
	Password       string                            // Contraseña para realizar la conexión
	DatabaseName   string                            // Nobre de la base de datos para conexión o realizar backup
}
