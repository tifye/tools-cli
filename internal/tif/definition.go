package tif

import (
	"fmt"
	"slices"
)

type TifDefinition struct {
	Schema string `json:"$schema"`
	Header struct {
		TifVersions  []string `json:"tifVersions"`
		ProtocolName string   `json:"protocolName"`
	} `json:"header"`
	AttributesV2 []AttributeV2Definition `json:"attributes-v2"`
	Methods      []MethodDefinition      `json:"methods"`
	TypesV2      []TypeV2Definition      `json:"types-v2"`
}

type TypeV2Definition struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type"`
	Postfix     string `json:"postfix"`
	Range       struct {
		Type  string `json:"type"`
		Start int    `json:"start"`
		Stop  int    `json:"stop"`
		Enums []struct {
			Key         string `json:"key"`
			Value       int    `json:"value"`
			Description string `json:"description"`
		} `json:"enum"`
	} `json:"range"`
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
