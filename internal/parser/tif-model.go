package parser

type TifDefinition struct {
	Schema string `json:"$schema"`
	Header struct {
		TifVersions  []string `json:"tifVersions"`
		ProtocolName string   `json:"protocolName"`
	} `json:"header"`
	Attributes     []Attribute       `json:"attributes"`
	AttributesV2   []AttributeV2     `json:"attributes-v2"`
	AttributeLists []AttributeList   `json:"attribute-lists"`
	ListsV2        []AttributeListV2 `json:"lists-v2"`
	Types          []Type            `json:"types"`
	TypesV2        []TypeV2          `json:"types-v2"`
	Methods        []Method          `json:"methods"`
	Actions        []Action          `json:"actions"`
	ActionsV2      []ActionV2        `json:"actions-v2"`
	Events         []Event           `json:"events"`
	Documentations []Documentation   `json:"documentations"`
	StateMachines  []StateMachine    `json:"statemachines"`
	Logstrings     []struct {
		ID   string `json:"id"`
		Text string `json:"text"`
	} `json:"logstrings"`
}

type Attribute struct {
	Family      string   `json:"family"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Get         string   `json:"get,omitempty"`
	Params      []string `json:"params,omitempty"`
	Type        struct {
		Kind string `json:"kind"`
	} `json:"type"`
	Protocol struct {
		FamilyID    string `json:"familyId"`
		AttributeID string `json:"attributeId"`
		RawOnly     bool   `json:"rawOnly"`
	} `json:"protocol,omitempty"`
	Set        string `json:"set,omitempty"`
	ParamIndex string `json:"paramIndex,omitempty"`
}

type AttributeV2 struct {
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
	Protocol struct {
		Read struct {
			LoginLevels []string `json:"loginLevels"`
		} `json:"read"`
	} `json:"protocol"`
	Tags  []string `json:"tags,omitempty"`
	Write struct {
		Command struct {
			Family string `json:"family"`
			Name   string `json:"name"`
		} `json:"command"`
	} `json:"write,omitempty"`
	List struct {
		Family string `json:"family"`
		Name   string `json:"name"`
	} `json:"list,omitempty"`
}

type AttributeList struct {
	Family      string `json:"family"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Tags        string `json:"tags,omitempty"`
	IndexParam  struct {
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"index-param"`
	Protocol struct {
		FamilyID string `json:"familyId"`
		ListID   string `json:"listId"`
	} `json:"protocol,omitempty"`
}

type AttributeListV2 struct {
	Family string   `json:"family"`
	Name   string   `json:"name"`
	Tags   []string `json:"tags,omitempty"`
	Key    struct {
		Type       string   `json:"type"`
		Operations []string `json:"operations"`
	} `json:"key"`
	Protocol struct {
		FamilyID string `json:"familyId"`
		ListID   string `json:"listId"`
	} `json:"protocol,omitempty"`
	Operations  []string `json:"operations"`
	Description string   `json:"description,omitempty"`
}

type Type struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type"`
	Postfix     string `json:"postfix,omitempty"`
	Range       struct {
		Type  string `json:"type"`
		Enums []struct {
			Key         string `json:"key"`
			Description string `json:"description"`
			Value       int    `json:"value"`
		} `json:"enums"`
	} `json:"range,omitempty"`
}

type TypeV2 struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type"`
	Postfix     string `json:"postfix,omitempty"`
	Range       struct {
		Type  string `json:"type"`
		Enums []struct {
			Key         string `json:"key"`
			Description string `json:"description"`
			Value       int    `json:"value"`
		} `json:"enums"`
	} `json:"range,omitempty"`
}

type Method struct {
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

type Action struct {
	Family      string `json:"family"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Set         string `json:"set"`
}

type ActionV2 struct {
	Family      string   `json:"family"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	InParams    []struct {
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"inParams"`
	OutParams       []interface{} `json:"outParams"`
	MaxResponseTime int           `json:"maxResponseTime,omitempty"`
	Protocol        struct {
		FamilyID    string   `json:"familyId"`
		ActionID    string   `json:"actionId"`
		LoginLevels []string `json:"loginLevels"`
	} `json:"protocol"`
	List struct {
		Family string `json:"family"`
		Name   string `json:"name"`
	} `json:"list,omitempty"`
}

type Event struct {
	Family string   `json:"family"`
	Name   string   `json:"name"`
	Tags   []string `json:"tags,omitempty"`
	Params []struct {
		Name string   `json:"name"`
		Type string   `json:"type"`
		Tags []string `json:"tags"`
	} `json:"params"`
	Protocol []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"protocol"`
}

type Documentation struct {
	Family      string   `json:"family"`
	Name        string   `json:"name"`
	Tags        []string `json:"tags,omitempty"`
	Description string   `json:"description"`
}

type StateMachine struct {
	Family      string `json:"family"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Attribute   struct {
		Family string `json:"family"`
		Name   string `json:"name"`
		Param  string `json:"param"`
	} `json:"attribute"`
	Transitions []struct {
		Trigger   string `json:"trigger"`
		From      string `json:"from"`
		To        string `json:"to"`
		Condition string `json:"condition"`
		Action    string `json:"action"`
	} `json:"transitions"`
}
