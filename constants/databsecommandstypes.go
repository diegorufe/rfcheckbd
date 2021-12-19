package constants

// Indicamos comandos a utilizar
type DatabaseCommandsTypes int

const (
	Backup  DatabaseCommandsTypes = 0 // Backup de la base datos
	Migrate DatabaseCommandsTypes = 1 // Migración de la base de datos
)
