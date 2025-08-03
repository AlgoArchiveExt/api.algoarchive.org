package tests

import (
	githubutils "main/infra/utils/github"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapper(t *testing.T) {
	assert := assert.New(t)
	inputs := []string{"python3", "scala", "java", "cpp", "csharp", "", "yellow", "hihi"}

	expected := []string{"py", "scala", "java", "cpp", "cs", "", "", ""}

	for i, language := range inputs {
		result, _ := githubutils.MapLanguageStringToFileExtension(language)
		expectedOutput := expected[i]

		assert.Equal(result, expectedOutput)
	}
}
