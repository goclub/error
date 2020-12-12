package xerr

import (
	"os"
	"time"
)

func MustUint(i uint, err error) uint {
	if err != nil { panic(err) }
	return i
}
func MustUints(i []uint, err error) []uint {
	if err != nil { panic(err) }
	return i
}
func MustInt(i int, err error) int {
	if err != nil { panic(err) }
	return i
}
func MustInts(i []int, err error) []int {
	if err != nil { panic(err) }
	return i
}
func MustInt32s(i []int32, err error) []int32 {
	if err != nil { panic(err) }
	return i
}

func MustFloat64(f float64, err error) float64 {
	if err != nil { panic(err) }
	return f
}
func MustFloat64s(f []float64, err error) []float64 {
	if err != nil { panic(err) }
	return f
}
func MustFloat32(f float32, err error) float32 {
	if err != nil { panic(err) }
	return f
}
func MustFloat32s(f []float32, err error) []float32 {
	if err != nil { panic(err) }
	return f
}

func MustBytes(data []byte, err error) []byte {
	if err != nil { panic(err) }
	return data
}
func MustString(s string, err error) string {
	if err != nil { panic(err) }
	return s
}
func MustStrings(s []string, err error) []string {
	if err != nil { panic(err) }
	return s
}

func MustBool(b bool, err error) bool {
	if err != nil { panic(err) }
	return b
}
func MustBools(b []bool, err error) []bool {
	if err != nil { panic(err) }
	return b
}

func MustAny(v interface{}, err error) interface{} {
	if err != nil { panic(err) }
	return v
}
func MustAnys(v []interface{}, err error) []interface{} {
	if err != nil { panic(err) }
	return v
}
func MustTime(v time.Time, err error) time.Time {
	if err != nil { panic(err) }
	return v
}
func MustFile(file *os.File, err error) *os.File {
	if err != nil { panic(err) }
	return file
}