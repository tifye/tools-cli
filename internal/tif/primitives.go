package tif

import (
	"encoding/hex"
	"strconv"
	"strings"
)

// Todo: use time.Time instead
type TypeUnixTime int32
type TUnixTime struct{}

func (t *TUnixTime) ParseString(v string) (any, error) {
	base := getValueBase(v)
	integer, err := strconv.ParseInt(v, base, 32)
	if err != nil {
		return nil, err
	}
	return TypeUnixTime(integer), nil
}

type TypeUCS2 string
type TUCS2 struct{}

func (t *TUCS2) ParseString(v string) (any, error) {
	return TypeUCS2(v), nil
}

type TypeAscii string
type TAscii struct{}

func (t *TAscii) ParseString(v string) (any, error) {
	return TypeAscii(v), nil
}

type TypeByteArray []byte
type TByteArray struct{}

func (t *TByteArray) ParseString(v string) (any, error) {
	val, err := hex.DecodeString(v)
	return TypeByteArray(val), err
}

type TypeBool bool
type TBool struct{}

func (t *TBool) ParseString(v string) (any, error) {
	val, err := strconv.ParseBool(v)
	return TypeBool(val), err
}

type TypeUint8 uint8
type TUint8 struct{}

func (t *TUint8) ParseString(v string) (any, error) {
	return parseUint8(v)
}

type TypeUint16 uint16
type TUint16 struct{}

func (t *TUint16) ParseString(v string) (any, error) {
	return parseUint16(v)
}

type TypeUint32 uint32
type TUint32 struct{}

func (t *TUint32) ParseString(v string) (any, error) {
	return parseUint32(v)
}

type TypeUint64 uint64
type TUint64 struct{}

func (t *TUint64) ParseString(v string) (any, error) {
	return parseUint64(v)
}

type TypeInt8 int8
type TInt8 struct{}

func (t *TInt8) ParseString(v string) (any, error) {
	return parseInt8(v)
}

type TypeInt16 int16
type TInt16 struct{}

func (t *TInt16) ParseString(v string) (any, error) {
	return parseInt16(v)
}

type TypeInt32 int32
type TInt32 struct{}

func (t *TInt32) ParseString(v string) (any, error) {
	return parseInt32(v)
}

type TypeInt64 int64
type TInt64 struct{}

func (t *TInt64) ParseString(v string) (any, error) {
	return parseInt64(v)
}

type TypeFloat float32
type TFloat struct{}

func (t *TFloat) ParseString(v string) (any, error) {
	val, err := strconv.ParseFloat(v, 32)
	return TypeFloat(val), err
}

type TypeBit byte
type TBit struct{}

func (t *TBit) ParseString(v string) (any, error) {
	val, err := strconv.ParseUint(v, 2, 8)
	return TypeBit(val), err
}

type TypeRoboticsVersion uint16
type TRoboticsVersion struct{}

func (t *TRoboticsVersion) ParseString(v string) (any, error) {
	val, err := parseInt16(v)
	return TypeRoboticsVersion(val), err
}

type ErrUnkownType struct {
	Type string
}

func (e *ErrUnkownType) Error() string {
	return "unkown type passed to [ParseType]; passed type: " + e.Type
}

func ParseType(tifType, data string) (any, error) {
	switch tifType {
	case "tUnixTime":
		base := getValueBase(data)
		integer, err := strconv.ParseInt(data, base, 32)
		if err != nil {
			return nil, err
		}
		return TypeUnixTime(integer), nil
	case "ascii":
		return TypeAscii(data), nil
	case "byteArray":
		val, err := hex.DecodeString(data)
		return TypeByteArray(val), err
	case "bool":
		val, err := strconv.ParseBool(data)
		return TypeBool(val), err
	case "uint8":
		return parseUint8(data)
	case "uint16":
		return parseUint16(data)
	case "uint32":
		return parseUint32(data)
	case "uint64":
		return parseUint64(data)
	case "sint8":
		return parseInt8(data)
	case "sint16":
		return parseInt16(data)
	case "sint32":
		return parseInt32(data)
	case "sint64":
		return parseInt64(data)
	case "tUCS2":
		return TypeUCS2(data), nil
	case "tSimpleVersion":
		val, err := parseInt16(data)
		return TypeRoboticsVersion(val), err
	default:
		return nil, &ErrUnkownType{Type: tifType}
		// case "bit":
		// case "dateTime":
		// case "float":
	}
}

func getValueBase(v string) int {
	if strings.HasPrefix(v, "0x") {
		return 16
	}
	return 10
}

func parseUint8(v string) (TypeUint8, error) {
	base := getValueBase(v)
	integer, err := strconv.ParseUint(v, base, 8)
	if err != nil {
		return 0, err
	}
	return TypeUint8(integer), nil
}

func parseUint16(v string) (TypeUint16, error) {
	base := getValueBase(v)
	integer, err := strconv.ParseUint(v, base, 16)
	if err != nil {
		return 0, err
	}
	return TypeUint16(integer), nil
}

func parseUint32(v string) (TypeUint32, error) {
	base := getValueBase(v)
	integer, err := strconv.ParseUint(v, base, 32)
	if err != nil {
		return 0, err
	}
	return TypeUint32(integer), nil
}

func parseUint64(v string) (TypeUint64, error) {
	base := getValueBase(v)
	integer, err := strconv.ParseUint(v, base, 64)
	if err != nil {
		return 0, err
	}
	return TypeUint64(integer), nil
}

func parseInt8(v string) (TypeInt8, error) {
	base := getValueBase(v)
	integer, err := strconv.ParseInt(v, base, 8)
	if err != nil {
		return 0, err
	}
	return TypeInt8(integer), nil
}

func parseInt16(v string) (TypeInt16, error) {
	base := getValueBase(v)
	integer, err := strconv.ParseInt(v, base, 16)
	if err != nil {
		return 0, err
	}
	return TypeInt16(integer), nil
}

func parseInt32(v string) (TypeInt32, error) {
	base := getValueBase(v)
	integer, err := strconv.ParseInt(v, base, 32)
	if err != nil {
		return 0, err
	}
	return TypeInt32(integer), nil
}

func parseInt64(v string) (TypeInt64, error) {
	base := getValueBase(v)
	integer, err := strconv.ParseInt(v, base, 64)
	if err != nil {
		return 0, err
	}
	return TypeInt64(integer), nil
}
