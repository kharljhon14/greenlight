package data

import (
	"fmt"
	"strconv"
)

type Runtime int32

func (r Runtime) MarshalJSON() ([]byte, error) {
	// Generate a string containing the movie runtime in the required format
	jsonValue := fmt.Sprintf("%d mins", r)

	// Wrap string in double quotes
	quotedJSONValue := strconv.Quote(jsonValue)

	// Convert the quoted string value to a byte slice and reurn it
	return []byte(quotedJSONValue), nil
}
