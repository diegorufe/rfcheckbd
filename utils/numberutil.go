package utils

import "strconv"

// IncrementeCounter método para incremetar el valor de un puntero contador
//
// @parameter counter puntero del contador a incrementa
//
// return el puntero con el valor incrementado
func IncrementeCounter(counter *int64) int64 {
	*counter = *counter + 1
	return *counter
}

// IncrementeCounterAndReturnToString método para incremetar el valor de un puntero contador y convertirlo a un string
//
// @parameter counter puntero del contador a incrementa
//
// return el puntero con el valor incrementado en formato string
func IncrementeCounterAndReturnToString(counter *int64) string {
	return strconv.FormatInt(IncrementeCounter(counter), 10)
}
