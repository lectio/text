package text

import (
	"context"
	"fmt"
	"regexp"
)

var sourceNameAfterPipeRegEx = regexp.MustCompile(` \| .*$`)   // Matches " | Healthcare IT News" from a title like "xyz title | Healthcare IT News"
var sourceNameAfterHyphenRegEx = regexp.MustCompile(` \- .*$`) // Matches " - Healthcare IT News" from a title like "xyz title - Healthcare IT News"
var firstSentenceRegExp = regexp.MustCompile(`^(.*?)[.?!]`)

type textTransformer interface {
	TransformText(ctx context.Context, from string, audit func(ctx context.Context, message string)) string
}

// RemovePipedSuffixFromText removes all text after a space and pipe (Matches " | Healthcare IT News" from a title like "xyz title | Healthcare IT News")
func RemovePipedSuffixFromText(ctx context.Context, from string, audit func(ctx context.Context, message string)) string {
	result := sourceNameAfterPipeRegEx.ReplaceAllString(from, "")
	if result != from && audit != nil {
		audit(ctx, fmt.Sprintf("Removed piped suffix from %q, now %q", from, result))
	}
	return result
}

// RemoveHyphenatedSuffix removes all text after a space and hyphen (// Matches " - Healthcare IT News" from a title like "xyz title - Healthcare IT News")
func RemoveHyphenatedSuffix(ctx context.Context, from string, audit func(ctx context.Context, message string)) string {
	result := sourceNameAfterHyphenRegEx.ReplaceAllString(from, "")
	if result != from && audit != nil {
		audit(ctx, fmt.Sprintf("Removed hyphen suffix from %q, now %q", from, result))
	}
	return result
}

// WarnHyphenatedSuffix adds a message to the audit log when hyphenated text is found
func WarnHyphenatedSuffix(ctx context.Context, from string, audit func(ctx context.Context, message string)) string {
	if sourceNameAfterHyphenRegEx.MatchString(from) && audit != nil {
		audit(ctx, fmt.Sprintf("Hyphenated suffix found in %q", from))
	}
	return from // we don't transform it, we just audit
}

// TransformText uses the options and transforms the from string through all the transformations
func TransformText(ctx context.Context, from string, audit func(ctx context.Context, message string), options ...interface{}) string {
	result := from
	for _, option := range options {
		if v, ok := option.(textTransformer); ok {
			result = v.TransformText(ctx, result, audit)
		}
		if transformText, ok := option.(func(ctx context.Context, from string, audit func(ctx context.Context, message string)) string); ok {
			result = transformText(ctx, result, audit)
		}
	}
	return result
}
