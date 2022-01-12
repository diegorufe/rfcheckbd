# rfcheckbd

El objetivo de este proyecto es poder realizar migraciones y comprobaciones de bases de datos sin tener que depender de proyectos externos como flyaway, liquidbase o herramientas interas de migración de frameworks de desarrollo como por ejemplo en .net, laravel ruby on rails etc ... 

## Funcionalidades

- (Pendiente) Realizar backup base de datos mysql
- (Pendiente) Implemetación de migración de bases de datos mysql
- (Pendiente) Separación de versionado de base de datos por módulos 
- (Pendiente) Realizar backup base de datos postgress
- (Pendiente) Implemetación de migración de bases de datos postgress
- (Pendiente) Modo servidor web

## Ventajas 

- Es idependiente del lenguage de programción y no requiere instalación de terceros como pueda ser la jvm 
- (Pendiente) Permite la separación de versionado de base de datos por módulo. Ejemplo: Es común en proyectos monoliticos modulares o en servicios rest tener una base de datos común, para ello cada módulo tiene unas tablas en especifico por ejemplo el módulo de contabilidad tendrá tablas como asientos y apuntes contables pero módulos como gestión de calidad solo tendrán parámetros y metricadas de calidad. En estos casos lo ideal sería poder separar el versionado de base de datos puesto que si mi proyecto monolitico modular no tiene el módulo de calidad no tiene ningún sentido que yo actualice sus tablas o si las tengo que actualizar por lo menos saber cual ha sido el último cambio por el equipo de desarrollo en esas tablas para ver si rompen con la funcionadlidad de mi aplicación
- (Pendiente) Al poder separar versionado de base de datos por módulo nos da la flexibidad de poder separar equipos en proyectos monoliticos modulares o servicios rest que comparten la misma base datos sin el temor de que toca "que" y mayor control operacional sobre el equipo
- (Pendiente) Al poder separar versionado de base de datos por módulo nos permite poder ejecutar una versión en especifico de ese módulo sin tener que cambiar toda la base de datos entera para movernos a una versión. Esta parte quizas es la más interesante para probar rendimiento entre 2 versiones de un módulo en cuestión

## Desventajas 

- Hay que manterlo e implementar cada tipo distinto de motor de base datos 
- Tener en cuenta ecosistemas ORM y OGM (SQL y NOSQL)
- No genera ningún beneficio economico directo su desarrolo

### Instalación

Simplemente es desacargase el fuente del código de github tener go instalado en algún equipo y construir la aplicación para el sistema operativo que se necesite. (Pendiente) Cuando se tenga modo web se podría desplegar en un docker/kubernetes
