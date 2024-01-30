// SizeStruct project sizestruct.go
package sizestruct

import (
	"fmt"
	"reflect"
)

type sStruct struct {
	npm   map[uintptr]bool
	exNum int
}

func SizeOf(data interface{}) int {
	var npm = &sStruct{make(map[uintptr]bool), 0}
	num := npm.sizeof(reflect.ValueOf(data))
	return num //+ npm.exNum
}

func SizeTOf(data interface{}) int {
	var npm = &sStruct{make(map[uintptr]bool), 0}
	num := npm.sizeof(reflect.ValueOf(data))
	return num + npm.exNum
}

func (s *sStruct) sizeofSlice(v reflect.Value) (sum int) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("panic sizeofSlice %v\n", err)
		}
	}()

	for i, n := 0, v.Len(); i < n; i++ {
		num := s.sizeof(v.Index(i))
		if num < 0 {
			return -1
		}
		sum += num
	}
	s.exNum += int(v.Type().Size())
	return sum
}

func (s *sStruct) sizeofArray(v reflect.Value) (sum int) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("panic sizeofArray %v\n", err)
		}
	}()
	for i, n := 0, v.Len(); i < n; i++ {
		num := s.sizeof(v.Index(i))
		if num < 0 {
			return -1
		}
		sum += num
	}
	return sum
}

func (s *sStruct) sizeofMap(v reflect.Value) (sum int) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("panic sizeofMap %v\n", err)
		}
	}()

	keys := v.MapKeys()
	for i := 0; i < len(keys); i++ {
		mapkey := keys[i]
		num := s.sizeof(mapkey)
		if num < 0 {
			return -1
		}
		sum += num
		num = s.sizeof(v.MapIndex(mapkey))
		if num < 0 {
			return -1
		}
		sum += num
	}
	s.exNum += int(v.Type().Size())
	return sum
}
func (s *sStruct) sizeof(v reflect.Value) (sum int) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("panic sizeof %v\n", err)
		}
	}()

	switch v.Kind() {
	case reflect.Map:
		return s.sizeofMap(v)

	case reflect.Slice:
		return s.sizeofSlice(v)
	case reflect.Array:
		return s.sizeofArray(v)

	case reflect.String:
		vs := v.Interface().(string)
		s.exNum += int(v.Type().Size())
		return len(vs)

	case reflect.Ptr:
		s.exNum += int(v.Type().Size())
		if v.IsNil() {
			return 0
		}
		//fmt.Println(v.Pointer())
		if _, ok := s.npm[v.Pointer()]; ok {
			return 0
		} else {
			s.npm[v.Pointer()] = true
		}
		return s.sizeof(v.Elem())

	case reflect.Interface:
		s.exNum += int(v.Type().Size())
		if v.IsNil() {
			return 0
		}
		return s.sizeof(v.Elem())

	case reflect.Uintptr: //Don't think it's Pointer 不认为是指针
		return int(v.Type().Size())

	case reflect.UnsafePointer: //Don't think it's Pointer 不认为是指针
		return int(v.Type().Size())

	case reflect.Struct:
		for i, n := 0, v.NumField(); i < n; i++ {
			if v.Type().Field(i).Tag.Get("ss") == "-" {
				continue
			}
			num := s.sizeof(v.Field(i))
			if num < 0 {
				return -1
			}
			sum += num
		}
		return sum

	case reflect.Func, reflect.Chan:
		s.exNum += int(v.Type().Size())
		if v.IsNil() {
			return 0
		}
		return 0 //Temporary non handling func,chan.
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128,
		reflect.Int:
		return int(v.Type().Size())
	case reflect.Bool:
		return int(v.Type().Size())
	default:
		fmt.Println("t.Kind() no found:", v.Kind())
	}

	return -1
}
