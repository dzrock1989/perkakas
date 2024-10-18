package params

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type Params struct {
	Uuid uuid.UUID `json:"uuid"`
}

// Decode help to decoded params
// @params byte of params from rpc
// @p is pointer to params struct for decoded result
func Decode[T any](params []byte, p *T) (err error) {
	if err = json.Unmarshal(params, p); err != nil {
		err = fmt.Errorf("invalid params")
		return
	}

	return
}
