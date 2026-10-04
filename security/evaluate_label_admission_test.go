package security_test

import (
	"strings"
	"testing"

	"github.com/faustbrian/go-openapi/v2/jsonvalue"
	"github.com/faustbrian/go-openapi/v2/security"
)

func TestSatisfiedLabelByteAdmission(t *testing.T) {
	const budget = 1 << 20
	limits := security.Limits{MaxAlternatives: 1, MaxSchemes: 2, MaxScopes: 2}
	assertAdmitted := func(t *testing.T, requirements jsonvalue.Value, credentials security.Credentials, want bool) {
		t.Helper()
		matched, err := security.Satisfied(requirements, credentials, limits)
		if err != nil || matched != want {
			t.Fatal("inclusive label budget changed admission or matching")
		}
	}

	t.Run("requirement scheme name", func(t *testing.T) {
		name := strings.Repeat("r", budget+1)
		assertAdmitted(t, arrayValue(t, objectValue(t,
			jsonvalue.Member{Name: name[:budget], Value: arrayValue(t)},
		)), nil, false)
		matched, err := security.Satisfied(arrayValue(t, objectValue(t,
			jsonvalue.Member{Name: name, Value: arrayValue(t)},
		)), nil, limits)
		assertAdmissionLimit(t, matched, err)
	})

	t.Run("requirement scope aggregate", func(t *testing.T) {
		// The scheme name and scope share one requirement-byte budget.
		scope := strings.Repeat("r", budget-len("OAuth")+1)
		assertAdmitted(t, arrayValue(t, objectValue(t,
			jsonvalue.Member{Name: "OAuth", Value: arrayValue(t, stringValue(t, scope[:len(scope)-1]))},
		)), nil, false)
		matched, err := security.Satisfied(arrayValue(t, objectValue(t,
			jsonvalue.Member{Name: "OAuth", Value: arrayValue(t, stringValue(t, scope))},
		)), nil, limits)
		assertAdmissionLimit(t, matched, err)
	})

	t.Run("credential scheme name aggregate", func(t *testing.T) {
		requirements := arrayValue(t, objectValue(t,
			jsonvalue.Member{Name: "OAuth", Value: arrayValue(t)},
		))
		// An unused scheme still consumes the shared credential-byte budget.
		name := strings.Repeat("c", budget-len("OAuth")+1)
		assertAdmitted(t, requirements, security.Credentials{
			"OAuth": {}, name[:len(name)-1]: {},
		}, true)
		matched, err := security.Satisfied(requirements, security.Credentials{
			"OAuth": {}, name: {},
		}, limits)
		assertAdmissionLimit(t, matched, err)
	})

	t.Run("credential scope aggregate", func(t *testing.T) {
		requirements := arrayValue(t, objectValue(t,
			jsonvalue.Member{Name: "OAuth", Value: arrayValue(t, stringValue(t, "read"))},
		))
		// The scheme and every granted scope occurrence consume the budget.
		scope := strings.Repeat("c", budget-len("OAuth")-len("read")+1)
		assertAdmitted(t, requirements, security.Credentials{
			"OAuth": {"read", scope[:len(scope)-1]},
		}, true)
		matched, err := security.Satisfied(requirements, security.Credentials{
			"OAuth": {"read", scope},
		}, limits)
		assertAdmissionLimit(t, matched, err)
	})

	t.Run("independent requirement and credential domains", func(t *testing.T) {
		name := strings.Repeat("s", budget+1)
		requirements := arrayValue(t, objectValue(t,
			jsonvalue.Member{Name: name[:budget], Value: arrayValue(t)},
		))
		// Each domain may independently consume the complete fixed budget.
		assertAdmitted(t, requirements, security.Credentials{name[:budget]: {}}, true)
		matched, err := security.Satisfied(requirements, security.Credentials{name: {}}, limits)
		assertAdmissionLimit(t, matched, err)
	})
}
