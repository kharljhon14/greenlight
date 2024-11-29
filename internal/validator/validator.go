package validator

import "regexp"

// Regular expression for sanity checking the format of email addresses
var (
	EmailRX = regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`)
)

type Validator struct {
	Errors map[string]string
}

// Creates a new validator instance with an empty error map
func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

// Check if the errors map contains any entries
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// Add and error message to the map only if there is no entry already exists for the given key
func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// Check adds an error message to the map only if a validation check is not ok
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// Check if a specific value is in a list of strings
func In(value string, list ...string) bool {
	for i := range list {
		if value == list[i] {
			return true
		}
	}

	return false
}

// Check if a value matches a regex
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// Check if all string values are unique
func Unique(values []string) bool {
	uniqueValues := make(map[string]bool)

	for _, value := range values {
		uniqueValues[value] = true
	}

	return len(values) == len(uniqueValues)
}
