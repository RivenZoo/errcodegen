package patch

import (
	"errors"
	"reflect"
)

// CoverStructsField assign exported fields of src to dst if field is not empty.
// Param src should be the same type of dst.
// Param dst should be pointer.
func CoverStructsField(src, dst interface{}) error {
	if reflect.TypeOf(dst).Kind() != reflect.Ptr {
		return errors.New("param dst should be ptr")
	}

	vsrc := reflect.ValueOf(src)
	vdst := reflect.ValueOf(dst)
	if !vsrc.IsValid() {
		return nil
	}
	if !vdst.IsValid() || vdst.IsNil() {
		return nil
	}

	if reflect.Indirect(vsrc).Type() != reflect.Indirect(vdst).Type() {
		return errors.New("param src and dst should be same type")
	}

	vsrc = reflect.Indirect(vsrc)
	vdst = reflect.Indirect(vdst)

	for i := 0; i < vdst.NumField(); i++ {
		fsrc := vsrc.Field(i)
		fdst := vdst.Field(i)

		if !reflect.DeepEqual(fsrc.Interface(), reflect.Zero(fsrc.Type()).Interface()) {
			fdst.Set(fsrc)
		}
	}
	return nil
}
