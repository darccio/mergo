package mergo_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"dario.cat/mergo"
)

func TestJsonNumber(t *testing.T) {
	jsonSampleData := `
{
	"amount": 1234
}
`

	type Data struct {
		Amount int64 `json:"amount"`
	}

	foo := make(map[string]interface{})

	decoder := json.NewDecoder(bytes.NewReader([]byte(jsonSampleData)))
	decoder.UseNumber()
	decoder.Decode(&foo)

	data := Data{}
	err := mergo.Map(&data, foo)
	if err != nil {
		t.Errorf("failed to merge with json.Number: %+v", err)
	}
	if data.Amount != 1234 {
		t.Errorf("merged amount does not match the json value! expected: 1234 got: %v", data.Amount)
	}
}
