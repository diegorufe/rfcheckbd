package constants

// Indicamos que tipo de base de datos es
type DatabaseType int

const (
	Mysql     DatabaseType = 0 // Tipo de base datos mysql
	Postgress DatabaseType = 1 // Tipo de base de datos postgrees
)
