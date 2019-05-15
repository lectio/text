package text

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

type TextSuite struct {
	suite.Suite
}

func (suite *TextSuite) SetupSuite() {
}

func (suite *TextSuite) TearDownSuite() {
}

func (suite *TextSuite) TestPipeTransform() {
	ctx := context.Background()
	original := "this is an xyz title | Healthcare IT News"
	final := TransformText(ctx, original, func(ctx context.Context, message string) {
		fmt.Println(message)
	}, RemovePipedSuffixFromText)
	suite.Equal("this is an xyz title", final, "Everything after the pipe should be removed")
}

func (suite *TextSuite) TestHyphenTransform() {
	ctx := context.Background()
	original := "this is an xyz title, with hyphen-nated word too - The Wall Street Journal"
	final := TransformText(ctx, original, func(ctx context.Context, message string) {
		fmt.Println(message)
	}, RemoveHyphenatedSuffix)
	suite.Equal("this is an xyz title, with hyphen-nated word too", final, "Everything after the last hyphen should be removed")
}
func (suite *TextSuite) TestHyphenWarning() {
	ctx := context.Background()
	original := "this is an xyz title, with hyphen-nated word too - The Wall Street Journal"
	var warning string
	final := TransformText(ctx, original, func(ctx context.Context, message string) {
		warning = message
	}, WarnHyphenatedSuffix)
	suite.Equal("Hyphenated suffix found in \""+original+"\"", warning, "Should have gotten a warning")
	suite.Equal(final, original, "Should not have been changed")
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(TextSuite))
}
