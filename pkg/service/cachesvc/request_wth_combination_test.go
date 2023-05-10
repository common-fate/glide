package cachesvc

// import (
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// )

// func TestRequestWithCombination(t *testing.T) {
// 	type testcase struct {
// 		name    string
// 		crw     CreateRequestWith
// 		want    RequestArgumentCombinations
// 		wantErr error
// 	}

// 	testcases := []testcase{
// 		{
// 			name: "ok one argument",
// 			crw: CreateRequestWith{
// 				AdditionalProperties: map[string][]string{
// 					"accountId": {"a", "b"},
// 				},
// 			},
// 			want: []map[string]string{{"accountId": "a"}, {"accountId": "b"}},
// 		},
// 		{
// 			name: "ok two arguments",
// 			crw: CreateRequestWith{
// 				AdditionalProperties: map[string][]string{
// 					"accountId":        {"a", "b"},
// 					"permissionSetArn": {"c", "d"},
// 				},
// 			},
// 			want: []map[string]string{{"permissionSetArn": "d", "accountId": "a"}, {"accountId": "b", "permissionSetArn": "c"}, {"accountId": "b", "permissionSetArn": "d"}, {"accountId": "a", "permissionSetArn": "c"}},
// 		},
// 		{
// 			name: "invalid if one of the arguments has no values",
// 			crw: CreateRequestWith{
// 				AdditionalProperties: map[string][]string{
// 					"accountId":        {},
// 					"permissionSetArn": {"c", "d"},
// 				},
// 			},
// 			wantErr: ArgumentHasNoValuesError{Argument: "accountId"},
// 		},
// 	}

// 	for _, tc := range testcases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			got, err := tc.crw.ArgumentCombinations()
// 			assert.ElementsMatch(t, tc.want, got)
// 			assert.Equal(t, tc.wantErr, err)
// 		})
// 	}
// }

// func TestHasDuplicates(t *testing.T) {
// 	type testcase struct {
// 		name string
// 		rac  RequestArgumentCombinations
// 		want bool
// 	}

// 	testcases := []testcase{
// 		{
// 			name: "no duplicates",
// 			rac: RequestArgumentCombinations{
// 				{
// 					"a": "b",
// 				},
// 				{
// 					"a": "c",
// 				},
// 			},
// 			want: false,
// 		},
// 		{
// 			name: "has duplicates",
// 			rac: RequestArgumentCombinations{
// 				{
// 					"a": "b",
// 				},
// 				{
// 					"a": "b",
// 				},
// 			},
// 			want: true,
// 		},
// 		{
// 			name: "has duplicates with lots of values",
// 			rac: RequestArgumentCombinations{
// 				{
// 					"a": "b",
// 					"c": "d",
// 				},
// 				{
// 					"a": "b",
// 					"c": "e",
// 				},
// 				{
// 					"a": "b",
// 					"c": "d",
// 				},
// 			},
// 			want: true,
// 		},
// 	}

// 	for _, tc := range testcases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			got := tc.rac.HasDuplicates()
// 			assert.Equal(t, tc.want, got)
// 		})
// 	}
// }
