package params

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestParams(t *testing.T) {
	data := map[string]string{
		"uuid": uuid.New().String(),
	}

	payload, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	var params Params
	if err = Decode(payload, &params); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		expected string
		actual   string
		err      string
	}{
		{"uuid", data["uuid"], params.Uuid.String(), "uuid not match"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.expected != test.actual {
				t.Error(test.err)
			}
		})
	}

	// Extend params
	data = map[string]string{
		"uuid":  uuid.New().String(),
		"other": "other params",
	}

	payload, err = json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	extededParams := struct {
		Params

		OtherParam string `json:"other"`
	}{}

	if err = Decode(payload, &extededParams); err != nil {
		t.Fatal(err)
	}

	tests = []struct {
		name     string
		expected string
		actual   string
		err      string
	}{
		{"uuid", data["uuid"], extededParams.Uuid.String(), "uuid not match"},
		{"other", data["other"], extededParams.OtherParam, "other params not match"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.expected != test.actual {
				t.Error(test.err)
			}
		})
	}
}

func TestParamsInvalidUUID(t *testing.T) {
	data := map[string]string{
		"uuid": uuid.New().String(),
	}

	payload, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	var params Params
	err = Decode(payload, &params)
	require.ErrorIs(t, nil, err, "should not error")

	data = map[string]string{
		"uuid": "invalid uuid",
	}

	payload, err = json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	err = Decode(payload, &params)
	require.EqualError(t, err, "invalid params", "expecting error invalid params")

	data = map[string]string{
		"uuid": "",
	}

	payload, err = json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	err = Decode(payload, &params)
	require.EqualError(t, err, "invalid params", "expecting error invalid params")
}
