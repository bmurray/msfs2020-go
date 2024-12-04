package simconnect

import (
	"fmt"
	"log/slog"
	"reflect"
	"unsafe"
)

// IsReport Convenience function to check if the data is the correct type
func IsReport[T any](s *SimConnect, ppData *RecvSimobjectDataByType) (*T, bool) {
	var typed *T
	defineId := s.GetDefineID(typed)
	if ppData.DefineID == defineId {
		return (*T)(unsafe.Pointer(ppData)), true
	}
	return nil, false
}

// RequestData Convenience function to request data
func RequestData[T any](s *SimConnect) error {
	var report *T
	defineId := s.GetDefineID(report)
	reqId := defineId
	return s.RequestDataOnSimObjectType(reqId, defineId, 0, SIMOBJECT_TYPE_USER)
}

// SetData currently only supports float64 fields
func (s *SimConnect) SetData(fr any) error {
	defineId := s.GetDefineID(fr)

	cnt := 0

	val := reflect.ValueOf(fr)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	typ := val.Type()
	if typ.Kind() != reflect.Struct {
		return fmt.Errorf("not a struct: %s", typ.Kind().String())
	}
	buf := []float64{}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		name := field.Tag.Get("name")
		if name == "" {
			continue
		}
		if field.Type.Kind() != reflect.Float64 {
			// if field.Name == "RecvSimobjectDataByType" {
			// 	continue
			// }
			return fmt.Errorf("not a float64: %s -- %s", field.Name, field.Type.Kind().String())
		}
		buf = append(buf, val.Field(i).Float())
		cnt++
	}

	size := DWORD(cnt * 8)
	slog.Debug("Setting data", "defineid", defineId, "count", cnt, "size", size)
	return s.SetDataOnSimObject(defineId, OBJECT_ID_USER, 0, 0, size, unsafe.Pointer(&buf[0]))

}
