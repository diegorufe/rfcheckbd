package beans

import "database/sql"

// CacheProcess cache donde guardar los datos que se necesitan para ir realizando las acciones del procesado
type CacheProcess struct {
	DBSql *sql.DB // Connexión sql con la base de datos
}
