package tif

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

type MethodDefinition struct {
	Command     string `json:"command"`
	ElementType string `json:"elementType"`
	Family      string `json:"family"`
	InParams    []struct {
		Name   string `json:"name"`
		Type   string `json:"type"`
		Length string `json:"length"`
	} `json:"inParams"`
	OutParams []interface{} `json:"outParams"`
	Protocol  []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"protocol"`
	LoginLevels     []string `json:"loginLevels,omitempty"`
	Tags            []string `json:"tags,omitempty"`
	Description     string   `json:"description,omitempty"`
	MaxResponseTime string   `json:"maxResponseTime,omitempty"`
}
