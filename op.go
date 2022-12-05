package libexpression

import (
	"encoding/json"
	"reflect"
)

type OpType uint8

const (
	OpTypeUnspecified = iota
	OpTypeValue
	OpTypeAnd
	OpTypeOr
)

type MarkFlag uint8

const (
	MarkFlagUnspecified = iota
	MarkFlagTrue
	MarkFlagFalse
)

type Op struct {
	OpType OpType   `json:"ot,omitempty"`
	Values []*Op    `json:"vs,omitempty"`
	Flag   MarkFlag `json:"f,omitempty"`

	ValueType string      `json:"vt,omitempty"`
	Value     interface{} `json:"v,omitempty"`
}

type Op4JSONParse Op

func (op *Op) UnmarshalJSON(data []byte) error {
	var tmpJSON Op4JSONParse

	err := json.Unmarshal(data, &tmpJSON)
	if err != nil {
		return err
	}

	err = decodeCustom(&tmpJSON, data)
	if err != nil {
		return err
	}

	*op = Op(tmpJSON)

	return nil
}

func decodeCustom(tmpJSON *Op4JSONParse, data []byte) (err error) {
	if tmpJSON.ValueType == "" {
		return
	}

	m := getOpValueTypeM()
	if len(m) == 0 {
		return
	}

	ty, found := m[tmpJSON.ValueType]
	if !found {
		return
	}

	tmpJSON.Value = reflect.New(ty).Interface()

	err = json.Unmarshal(data, tmpJSON)

	return
}
