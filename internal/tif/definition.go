package tif

import (
	"encoding/hex"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

type TifDefinition struct {
	Schema string `json:"$schema"`
	Header struct {
		TifVersions  []string `json:"tifVersions"`
		ProtocolName string   `json:"protocolName"`
	} `json:"header"`
	AttributesV2 []AttributeV2Definition `json:"attributes-v2"`
	Methods      []MethodDefinition      `json:"methods"`
}

type AttributeV2Definition struct {
	Family      string `json:"family"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Params      []struct {
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"params"`
	Operations []string `json:"operations"`
	Read       struct {
		Command struct {
			Family string `json:"family"`
			Name   string `json:"name"`
		} `json:"command"`
	} `json:"read,omitempty"`
	Write struct {
		Command struct {
			Family string `json:"family"`
			Name   string `json:"name"`
		} `json:"command"`
	} `json:"write,omitempty"`
	Protocol struct {
		Read struct {
			LoginLevels []string `json:"loginLevels"`
		} `json:"read"`
	} `json:"protocol"`
	Tags []string `json:"tags,omitempty"`
	List struct {
		Family string `json:"family"`
		Name   string `json:"name"`
	} `json:"list,omitempty"`
}

func (attr AttributeV2Definition) ReadCommand() (string, bool) {
	if !slices.Contains(attr.Operations, "read") {
		return "", false
	}

	return fmt.Sprintf("%s.%s", attr.Read.Command.Family, attr.Read.Command.Name), true
}

func (attr AttributeV2Definition) WriteCommand() (string, bool) {
	if !slices.Contains(attr.Operations, "write") {
		return "", false
	}

	return fmt.Sprintf("%s.%s", attr.Write.Command.Family, attr.Write.Command.Name), true
}

type InputParameter struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Length string `json:"length"`
}

type MethodDefinition struct {
	Family      string           `json:"family"`
	Command     string           `json:"command"`
	Description string           `json:"description,omitempty"`
	ElementType string           `json:"elementType"`
	InParams    []InputParameter `json:"inParams"`
	OutParams   []struct {
		Name string   `json:"name"`
		Type string   `json:"type"`
		Tags []string `json:"tags,omitempty"`
	} `json:"outParams"`
	Protocol []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"protocol"`
	LoginLevels     []string `json:"loginLevels,omitempty"`
	Tags            []string `json:"tags,omitempty"`
	MaxResponseTime string   `json:"maxResponseTime,omitempty"`
}

func (m MethodDefinition) Name() string {
	return fmt.Sprintf("%s.%s", m.Family, m.Command)
}

// Todo: use time.Time instead
type TypeUnixTime = int32
type TypeUCS2 = string
type TypeAscii = string
type TypeByteArray = []byte
type TypeBool = bool
type TypeUint8 = uint8
type TypeUint16 = uint16
type TypeUint32 = uint32
type TypeUint64 = uint64
type TypeInt8 = int8
type TypeInt16 = int16
type TypeInt32 = int32
type TypeInt64 = int64
type TypeFloat = float32
type TypeBit = byte

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
	default:
		return nil, &ErrUnkownType{Type: tifType}
		// case "bit":
		// case "tSimpleVersion":
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
