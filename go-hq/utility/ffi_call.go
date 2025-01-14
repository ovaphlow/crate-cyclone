package utility

/*
#cgo LDFLAGS: -L../lib/ -lcrate_shared
#include <stdlib.h>

extern char* generate_ksuid();
*/
import "C"
import "unsafe"

// CallGenerateKsuid 调用 libcrate_shared.so 中的 generate_ksuid 函数
func CallGenerateKsuid() string {
	// 调用 C 函数
	cStr := C.generate_ksuid()
	// 将 C 字符串转换为 Go 字符串
	goStr := C.GoString(cStr)
	// 释放 C 字符串
	C.free(unsafe.Pointer(cStr))
	return goStr
}
