package utils

import "time"

// MakeTimestamp método para obtener el tiempo actual en milisegundos
//
// @returns el tiempo actual en milisegundos
func MakeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
