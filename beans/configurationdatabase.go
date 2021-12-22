package beans

import (
	"rfcheckbd/constants"
)

// ConfigurationDatabase para indicar la configuración de la base datos
type ConfigurationDatabase struct {
	Type           constants.DatabaseType            `json:"Type"`           // Tipo de base de datos mysql, posgrees etc ...
	PathSaveBackup string                            `json:"PathSaveBackup"` // Ruta donde guardar el fichero de backup. Tiene que ser un directorio existente
	Commands       []constants.DatabaseCommandsTypes `json:"Commands"`       // Comandos a ejecutar
	User           string                            `json:"User"`           // Usuario para realizar la conexión
	Password       string                            `json:"Password"`       // Contraseña para realizar la conexión
	DatabaseName   string                            `json:"DatabaseName"`   // Nobre de la base de datos para conexión o realizar backup
	PathMigrations string                            `json:"PathMigrations"` // Ruta donde se encuentran las migraciones para ejecutar
	Host           string                            `json:"Host"`           // Host donde conectarse de la base de datos
	Port           int64                             `json:"Port"`           // Port puerto del host de connexión
}
