package mergo_test

import "testing"

type testCase struct {
	desc string
	test func(*testing.T)
}

func runTests(t *testing.T, testCases []testCase) {
	for _, testCase := range testCases {
		t.Run(testCase.desc, testCase.test)
	}
}

func unexpected(t *testing.T, expected, got interface{}) {
	t.Errorf("expected %+v, got %+v", expected, got)
}
