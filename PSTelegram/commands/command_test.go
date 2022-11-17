package commands

import (
	"testing"
)

func TestValidateArgumentsCounts(t *testing.T) {

}

func TestValidateArgumentsValues(t *testing.T) {
	testCases := []struct {
		arguments         []string
		expectedArguments []string
		isErrorExpected   bool
	}{
		{
			arguments:         []string{"asc"},
			expectedArguments: []string{"s(asc|desc)"},
			isErrorExpected:   false,
		},
		{
			arguments:         []string{"1"},
			expectedArguments: []string{"i"},
			isErrorExpected:   false,
		},
		{
			arguments:         []string{"asd"},
			expectedArguments: []string{"i"},
			isErrorExpected:   true,
		},
		{
			arguments:         []string{"a"},
			expectedArguments: []string{"s"},
			isErrorExpected:   false,
		},
		{
			arguments:         []string{"2022-01-01 00:00:23"},
			expectedArguments: []string{"d"},
			isErrorExpected:   false,
		},
		{
			arguments:         []string{"2022-01-123f asd01 00:00:23"},
			expectedArguments: []string{"d"},
			isErrorExpected:   true,
		},
		{
			arguments:         []string{"2022-01-01 25:00:23"},
			expectedArguments: []string{"d"},
			isErrorExpected:   true,
		},
	}

	for _, testCase := range testCases {
		command := StructureOfCommand{ExpectedArguments: testCase.expectedArguments}
		err := command.validateArgumentsValues(testCase.arguments)
		if (err != nil) != testCase.isErrorExpected {
			t.Errorf("Is error expected: %v, Is error returned: %v", testCase.isErrorExpected, err != nil)
		}
	}
}
