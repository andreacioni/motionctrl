package utils

import "reflect"

func InSlice(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}

	return
}

func BlockSlideSlice(array interface{}, blockSize int, f func(interface{})) {
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)
		n := s.Len()

		for i := 0; i < n; i += blockSize {
			if i+blockSize > n { //TODO improve here (?)
				f(s.Slice(i, n).Interface())
			} else {
				f(s.Slice(i, i+blockSize).Interface())
			}
		}
	}
}
