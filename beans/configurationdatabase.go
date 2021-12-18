package beans

import (
	"rfcheckbd/constants"
)

// ConfigurationDatabase para indicar la configuración de la base datos
type ConfigurationDatabase struct {
	Type constants.DatabaseType // Tipo de base de datos mysql, posgrees etc ...
}
