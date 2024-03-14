package security

import (
	"fmt"
	"unsafe"

	xsys "golang.org/x/sys/windows"
)

type cryptoFlag uint32

const CRYPTPROTECT_UI_FORBIDDEN = 0x1

var (
	crypt32            = xsys.NewLazyDLL("Crypt32.dll")
	cryptProtectProc   = crypt32.NewProc("CryptProtectData")
	cryptUnProtectProc = crypt32.NewProc("CryptUnprotectData")
)

/*
The CRYPT_INTEGER_BLOB structure contains an arbitrary array of bytes.
The structure definition includes aliases appropriate to the various functions that use it.
"https://learn.microsoft.com/en-us/previous-versions/windows/desktop/legacy/aa381414(v=vs.85)""
*/
type cryptoAPIBlob struct {
	cbData uint32 // Number of bytes in the pbData. DWORD == uint32.
	pbData *byte  // Pointer to first byte of data
}

func (b *cryptoAPIBlob) free() error {
	// LocalFree https://learn.microsoft.com/en-us/windows/win32/api/winbase/nf-winbase-localfree
	_, err := xsys.LocalFree(xsys.Handle(unsafe.Pointer(b.pbData)))
	return err
}

func (b *cryptoAPIBlob) toBytes() []byte {
	// Convert sliceHeader to a []byte slice using unsafe.Slice
	data := unsafe.Slice((*byte)(unsafe.Pointer(b.pbData)), b.cbData)

	// data's memory is unmanaged by Go, so we need to make a copy
	// this is because of unsafe
	result := make([]byte, len(data))
	copy(result, data)
	return result
}

func newCryptoAPIBlob(data []byte) *cryptoAPIBlob {
	return &cryptoAPIBlob{
		cbData: uint32(len(data)),
		pbData: &data[0],
	}
}

// https://learn.microsoft.com/en-us/windows/win32/api/dpapi/nf-dpapi-cryptprotectdata
func Encrypt(data []byte, entropy []byte) ([]byte, error) {
	var inputData = newCryptoAPIBlob(data)
	var entropyData = newCryptoAPIBlob(entropy)
	var outData cryptoAPIBlob

	inputDataPtr := uintptr(unsafe.Pointer(inputData))
	entropyDataPtr := uintptr(unsafe.Pointer(entropyData))
	outDataPtr := uintptr(unsafe.Pointer(&outData))
	dwFlags := uintptr(CRYPTPROTECT_UI_FORBIDDEN)
	didSucceed, _, err := cryptProtectProc.Call(inputDataPtr, 0, entropyDataPtr, 0, 0, dwFlags, outDataPtr)
	if didSucceed == 0 {
		return nil, fmt.Errorf("CryptProtectData failed: %w", err)
	}

	return outData.toBytes(), outData.free()
}

// https://learn.microsoft.com/en-us/windows/win32/api/dpapi/nf-dpapi-cryptunprotectdata
func Decrypt(data []byte, entropy []byte) ([]byte, error) {
	var inputData = newCryptoAPIBlob(data)
	var entropyData = newCryptoAPIBlob(entropy)
	var outData cryptoAPIBlob

	inputDataPtr := uintptr(unsafe.Pointer(inputData))
	entropyDataPtr := uintptr(unsafe.Pointer(entropyData))
	outDataPtr := uintptr(unsafe.Pointer(&outData))
	dwFlags := uintptr(CRYPTPROTECT_UI_FORBIDDEN)
	didSucceed, _, err := cryptUnProtectProc.Call(inputDataPtr, 0, entropyDataPtr, 0, 0, dwFlags, outDataPtr)
	if didSucceed == 0 {
		return nil, fmt.Errorf("CryptUnprotectData failed: %w", err)
	}

	return outData.toBytes(), outData.free()
}
