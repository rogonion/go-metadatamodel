package core

import (
	"errors"
	"regexp"

	"github.com/rogonion/go-json/path"
)

const (
	ArrayPathPlaceholder = path.JsonpathLeftBracket + path.JsonpathKeyIndexAll + path.JsonpathRightBracket
	GroupJsonPathPrefix  = path.JsonpathDotNotation + GroupFields + ArrayPathPlaceholder + path.JsonpathDotNotation
)

func ArrayPathRegexSearch() *regexp.Regexp {
	return regexp.MustCompile(`\[\*]`)
}

func GroupFieldsPathRegexSearch() *regexp.Regexp {
	return regexp.MustCompile(`GroupFields\[\*](?:\.|)`)
}

func GroupFieldsRegexSearch() *regexp.Regexp {
	return regexp.MustCompile(`(?:\.|)GroupFields`)
}

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
