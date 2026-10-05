package security_test

import (
	"testing"

	"github.com/faustbrian/go-openapi/v2/jsonvalue"
	"github.com/faustbrian/go-openapi/v2/security"
)

func TestSatisfiedAdmissionCredentialCounts(t *testing.T) {
	requirements := arrayValue(t, objectValue(t,
		jsonvalue.Member{Name: "OAuth", Value: arrayValue(t, stringValue(t, "read"))},
	))
	for _, test := range []struct {
		name      string
		limits    security.Limits
		inclusive security.Credentials
		over      security.Credentials
	}{
		{
			name:      "schemes including unused schemes",
			limits:    security.Limits{MaxAlternatives: 1, MaxSchemes: 1, MaxScopes: 2},
			inclusive: security.Credentials{"OAuth": {"read"}},
			over:      security.Credentials{"OAuth": {"read"}, "Other": {}},
		},
		{
			name:      "granted scopes",
			limits:    security.Limits{MaxAlternatives: 1, MaxSchemes: 1, MaxScopes: 1},
			inclusive: security.Credentials{"OAuth": {"read"}},
			over:      security.Credentials{"OAuth": {"read", "write"}},
		},
		{
			name:      "aggregate scopes across schemes",
			limits:    security.Limits{MaxAlternatives: 1, MaxSchemes: 2, MaxScopes: 2},
			inclusive: security.Credentials{"OAuth": {"read"}, "Other": {"write"}},
			over:      security.Credentials{"OAuth": {"read"}, "Other": {"write", "edit"}},
		},
		{
			name:      "duplicate granted scope occurrences",
			limits:    security.Limits{MaxAlternatives: 1, MaxSchemes: 1, MaxScopes: 1},
			inclusive: security.Credentials{"OAuth": {"read"}},
			over:      security.Credentials{"OAuth": {"read", "read"}},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			matched, err := security.Satisfied(requirements, test.inclusive, test.limits)
			if err != nil || !matched {
				t.Fatal("inclusive ordinary credentials were not admitted")
			}
			matched, err = security.Satisfied(requirements, test.over, test.limits)
			assertAdmissionLimit(t, matched, err)
		})
	}
}

func TestSatisfiedAdmissionCountsAreIndependent(t *testing.T) {
	requirements := arrayValue(t, objectValue(t,
		jsonvalue.Member{Name: "OAuth", Value: arrayValue(t, stringValue(t, "read"))},
		jsonvalue.Member{Name: "Other", Value: arrayValue(t, stringValue(t, "write"))},
	))
	matched, err := security.Satisfied(requirements,
		security.Credentials{"OAuth": {"read"}, "Other": {"write"}},
		security.Limits{MaxAlternatives: 1, MaxSchemes: 2, MaxScopes: 2})
	if err != nil || !matched {
		t.Fatal("requirement and credential counts were combined instead of independently admitted")
	}
}

func TestSatisfiedAdmissionAllRequirementsBeforeSuccess(t *testing.T) {
	for _, test := range []struct {
		name         string
		requirements jsonvalue.Value
		credentials  security.Credentials
		limits       security.Limits
	}{
		{
			name: "later scheme after matching alternative",
			requirements: arrayValue(t,
				objectValue(t, jsonvalue.Member{Name: "OAuth", Value: arrayValue(t)}),
				objectValue(t, jsonvalue.Member{Name: "Other", Value: arrayValue(t)})),
			credentials: security.Credentials{"OAuth": {}},
			limits:      security.Limits{MaxAlternatives: 2, MaxSchemes: 1, MaxScopes: 1},
		},
		{
			name: "later scope after matching alternative",
			requirements: arrayValue(t,
				objectValue(t, jsonvalue.Member{Name: "OAuth", Value: arrayValue(t, stringValue(t, "read"))}),
				objectValue(t, jsonvalue.Member{Name: "OAuth", Value: arrayValue(t, stringValue(t, "write"))})),
			credentials: security.Credentials{"OAuth": {"read"}},
			limits:      security.Limits{MaxAlternatives: 2, MaxSchemes: 2, MaxScopes: 1},
		},
		{
			name: "later scopes after anonymous alternative",
			requirements: arrayValue(t, objectValue(t),
				objectValue(t, jsonvalue.Member{Name: "OAuth", Value: arrayValue(t,
					stringValue(t, "read"), stringValue(t, "write"))})),
			limits: security.Limits{MaxAlternatives: 2, MaxSchemes: 1, MaxScopes: 1},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			matched, err := security.Satisfied(test.requirements, test.credentials, test.limits)
			assertAdmissionLimit(t, matched, err)
		})
	}
}

func TestSatisfiedAdmissionAnonymousBypassesCredentials(t *testing.T) {
	credentials := security.Credentials{"OAuth": {"read", "write"}, "Other": {}}
	for _, requirements := range []jsonvalue.Value{
		arrayValue(t),
		arrayValue(t,
			objectValue(t, jsonvalue.Member{Name: "OAuth", Value: arrayValue(t, stringValue(t, "read"))}),
			objectValue(t)),
	} {
		matched, err := security.Satisfied(requirements, credentials,
			security.Limits{MaxAlternatives: 2, MaxSchemes: 1, MaxScopes: 1})
		if err != nil || !matched {
			t.Fatal("valid anonymous requirements unnecessarily inspected irrelevant credentials")
		}
	}
}

func TestSatisfiedAdmissionRepeatedSchemeAlternatives(t *testing.T) {
	requirements := arrayValue(t,
		objectValue(t, jsonvalue.Member{Name: "OAuth", Value: arrayValue(t, stringValue(t, "write"))}),
		objectValue(t, jsonvalue.Member{Name: "OAuth", Value: arrayValue(t, stringValue(t, "read"))}))
	credentials := security.Credentials{"OAuth": {"read"}}
	limits := security.Limits{MaxAlternatives: 2, MaxSchemes: 2, MaxScopes: 2}
	matched, err := security.Satisfied(requirements, credentials, limits)
	if err != nil || !matched {
		t.Fatal("repeated scheme alternatives lost an admitted matching scope")
	}
	credentials["OAuth"] = []string{"edit"}
	matched, err = security.Satisfied(requirements, credentials, limits)
	if err != nil || matched {
		t.Fatal("call-local credential state was reused across evaluations")
	}
}

func assertAdmissionLimit(t *testing.T, matched bool, err error) {
	t.Helper()
	if matched || err != security.ErrLimitExceeded {
		t.Fatal("one-over ordinary evaluation did not fail closed with the fixed limit sentinel")
	}
	if err.Error() != "OpenAPI security evaluation limit exceeded" {
		t.Fatal("limit error contained unexpected detail")
	}
}
