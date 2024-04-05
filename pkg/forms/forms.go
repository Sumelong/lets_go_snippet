package forms

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

// EmailRX Use the regexp.MustCompile() function to parse a pattern and compile a
// regular expression for sanity checking the format of an email address.
// This returns a *regexp.Regexp object, or panics in the event of an error.
// Doing this once at runtime, and storing the compiled regular expression
// object in a variable, is more performant than re-compiling the pattern with
// every request.
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Form Create a custom Form struct, which anonymously embeds a url.Values object
// (to hold the form data) and an Errors field to hold any validation errors
// for the form data.
type Form struct {
	Values url.Values
	Errors errors
}

// NewForm Define a New function to initialize a custom Form struct. Notice that
// this takes the form data as the parameter?
func NewForm(val url.Values) *Form {
	return &Form{
		Values: val,
		Errors: make(errors),
	}
}

// Required Implement a Required method to check that specific fields in the form
// data are present and not blank. If any fields fail this check, add the
// appropriate message to the form errors.
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		val := f.Values.Get(field)
		if strings.TrimSpace(val) == "" {
			msg := fmt.Sprintf("%s cannot be empty", field)
			f.Errors.Add(field, msg)
		}
	}
}

// MaxLength Implement a MaxLength method to check that a specific field in the form
// contains a maximum number of characters. If the check fails then add the
// appropriate message to the form errors.
func (f *Form) MaxLength(field string, d int) {
	val := f.Values.Get(field)
	if len(val) == 0 {
		return
	}
	if utf8.RuneCountInString(val) > d {
		msg := fmt.Sprintf("%s is too long.(maximum allowed character is %d)", field, d)
		f.Errors.Add(field, msg)
	}
}

// PermittedValues Implement a PermittedValues method to check that a specific field in the form
// matches one of a set of specific permitted values. If the check fails
// then add the appropriate message to the form errors.
func (f *Form) PermittedValues(field string, opts ...string) {
	val := f.Values.Get(field)
	if len(val) == 0 {
		return
	}

	for _, opt := range opts {
		if val == opt {
			return
		}
	}
	msg := fmt.Sprintf("%s is not valid %s", val, field)
	f.Errors.Add(field, msg)
}

// IsString  Implement a IsString method to check that a specific field in the form
// is sending valid string values. If the check fails
// then add the appropriate message to the form errors.
func (f *Form) IsString(fields ...string) {
	for _, field := range fields {
		val := f.Values.Get(field)

		// Define a regular expression to match any special characters
		re := regexp.MustCompile("^[a-zA-Z0-9s!?.,:;-]+$")
		// Check if the string contains any matches for the regular expression
		match := re.MatchString(val)

		if match {
			msg := fmt.Sprintf("%s cannot have special characters", field)
			f.Errors.Add(field, msg)
		}
	}

}

// MinLength Implement a MinLength method to check that a specific field in the form
// contains a minimum number of characters. If the check fails then add the
// appropriate message to the form errors.
func (f *Form) MinLength(field string, d int) {
	val := f.Values.Get(field)
	if val == "" {
		return
	}
	if utf8.RuneCountInString(val) < d {
		f.Errors.Add(field, fmt.Sprintf("This field is too short (minimum is %d characters)", d))
	}
}

// MatchesPattern Implement a MatchesPattern method to check that a specific field in the form
// matches a regular expression. If the check fails then add the
// appropriate message to the form errors.
func (f *Form) MatchesPattern(field string, pattern *regexp.Regexp) {
	val := f.Values.Get(field)
	if val == "" {
		return
	}
	if !pattern.MatchString(val) {
		f.Errors.Add(field, "This field is invalid")
	}
}

// Valid Implement a Valid method which returns true if there are no errors.
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
