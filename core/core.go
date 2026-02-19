package core

import (
	"errors"
	"regexp"

	"github.com/rogonion/go-json/path"
)

const (
	// ArrayPathPlaceholder represents the wildcard index `[*]` used in JSON paths.
	ArrayPathPlaceholder = path.JsonpathLeftBracket + path.JsonpathKeyIndexAll + path.JsonpathRightBracket
	// GroupJsonPathPrefix represents the standard prefix for group fields in the metadata model path structure.
	GroupJsonPathPrefix = path.JsonpathDotNotation + GroupFields + ArrayPathPlaceholder + path.JsonpathDotNotation
)

// ArrayPathRegexSearch returns a compiled regex to find `[*]` in strings.
func ArrayPathRegexSearch() *regexp.Regexp {
	return regexp.MustCompile(`\[\*]`)
}

// GroupFieldsPathRegexSearch returns a compiled regex to find `GroupFields[*]` followed optionally by a dot.
func GroupFieldsPathRegexSearch() *regexp.Regexp {
	return regexp.MustCompile(`GroupFields\[\*](?:\.|)`)
}

// GroupFieldsRegexSearch returns a compiled regex to find `GroupFields` preceded optionally by a dot.
func GroupFieldsRegexSearch() *regexp.Regexp {
	return regexp.MustCompile(`(?:\.|)GroupFields`)
}

// SpecialCharsRegexSearch returns a compiled regex to find any character that is not alphanumeric.
func SpecialCharsRegexSearch() *regexp.Regexp {
	return regexp.MustCompile(`[^a-zA-Z0-9]+`)
}

// MetadataModelGroupReadOrderOfFields represents type for GroupReadOrderOfFields property.
type MetadataModelGroupReadOrderOfFields []string

var (
	// ErrArgumentInvalid For when an argument is not valid for the current action.
	ErrArgumentInvalid = errors.New("invalid argument")

	// ErrPathContainsIndexPlaceholders For when preparing JSONPath from metadata-model to value in object still contains ARRAY_PATH_PLACEHOLDER.
	ErrPathContainsIndexPlaceholders = errors.New("path contains index placeholders")
)
