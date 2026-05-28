package urlTools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateAlias(t *testing.T) {
	testData := []struct {
		testName    string
		alias       string
		expectError bool
	}{
		{"aliasTosmall", "a", true},
		{"aliasToBig", "aliasToBigggggggggggggggggg", true},
		{"aliasWithSpace", "alias with space", true},
		{"aliasWithControl", "aliasWithControl\n", true},
		{"validAlias", "validAlias", false},
	}
	for _, testDataItem := range testData {
		t.Run(testDataItem.testName, func(t *testing.T) {
			err := ValidateAlias(testDataItem.alias)
			assert.Equal(t, testDataItem.expectError, err != nil)
		})
	}
}
