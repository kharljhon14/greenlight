package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrInvalidRuntimeFormat = errors.New("invalid runtime format")

type Runtime int32

func (r Runtime) MarshalJSON() ([]byte, error) {
	// Generate a string containing the movie runtime in the required format
	jsonValue := fmt.Sprintf("%d mins", r)

	// Wrap string in double quotes
	quotedJSONValue := strconv.Quote(jsonValue)

	// Convert the quoted string value to a byte slice and reurn it
	return []byte(quotedJSONValue), nil
}

func (r *Runtime) UnmarshalJSON(jsonValue []byte) error {
	// Remove quotes if incoming JSON value will be a string
	unqoutedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	// Isolate the part that contains the number
	parts := strings.Split(unqoutedJSONValue, " ")

	// Sanity check the parts
	if len(parts) != 2 || parts[1] != "mins" {
		return ErrInvalidRuntimeFormat
	}

	// Parse the string containing the number into a int32
	i, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	// Convert int32 to Runtime type
	*r = Runtime(i)

	return nil
}
