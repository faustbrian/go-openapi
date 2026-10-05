// Package security evaluates OpenAPI Security Requirement Objects without
// coupling authorization decisions to an HTTP framework.
package security

import (
	"errors"
	"fmt"

	"github.com/faustbrian/go-openapi/v2/jsonvalue"
)

var (
	// ErrInvalidRequirements reports malformed Security Requirement values.
	ErrInvalidRequirements = errors.New("invalid OpenAPI security requirements")
	// ErrLimitExceeded reports bounded security evaluation exhaustion.
	ErrLimitExceeded = errors.New("OpenAPI security evaluation limit exceeded")
)

// MaxCredentialBytes bounds the total bytes in credential scheme names and
// granted scope occurrences accepted by one evaluation. The inclusive fixed
// budget is independent of MaxRequirementBytes.
const MaxCredentialBytes = 1 << 20

// MaxRequirementBytes bounds the total bytes in requirement scheme names and
// required scope occurrences accepted by one evaluation. The inclusive fixed
// budget is independent of MaxCredentialBytes.
const MaxRequirementBytes = 1 << 20

// Credentials maps an available security scheme to its granted OAuth or
// OpenID Connect scopes. Presence with an empty scope list satisfies schemes
// that do not use scopes. Satisfied borrows this map and its slices only during
// the call; callers must not mutate them concurrently with evaluation.
type Credentials map[string][]string

type requirementScheme struct {
	name   string
	scopes []jsonvalue.Value
}

// Limits bounds one Security Requirement array evaluation. MaxSchemes and
// MaxScopes apply independently to requirements and supplied credentials.
// Zero fields use DefaultLimits; negative fields are invalid. Positive overrides
// remain caller-selected finite limits, without an additional upper ceiling.
type Limits struct {
	MaxAlternatives int
	MaxSchemes      int
	MaxScopes       int
}

// DefaultLimits returns conservative limits for untrusted requirements.
func DefaultLimits() Limits {
	return Limits{MaxAlternatives: 1_000, MaxSchemes: 1_000, MaxScopes: 10_000}
}

// Satisfied reports whether credentials satisfy a Security Requirement array.
// Array entries are alternatives (OR), while names in one entry are combined
// requirements (AND). An empty array or an empty requirement permits anonymous
// access without inspecting credentials that cannot affect that result. All
// requirements are admitted and validated before any successful result. Collection
// lengths are checked before copies; credential counts and bytes are checked
// before building one reusable call-local index. Limit errors include no labels.
// Construction of the caller's semantic values and credentials is outside this
// evaluator's resource policy. Evaluation is synchronous and retains no state.
func Satisfied(
	requirements jsonvalue.Value,
	credentials Credentials,
	limits Limits,
) (bool, error) {
	limits, valid := effectiveLimits(limits)
	if !valid {
		return false, fmt.Errorf("%w: invalid limits", ErrInvalidRequirements)
	}
	if requirements.Kind() != jsonvalue.ArrayKind {
		return false, fmt.Errorf("%w: expected array", ErrInvalidRequirements)
	}
	alternativeCount, _ := requirements.Length()
	if alternativeCount > limits.MaxAlternatives {
		return false, ErrLimitExceeded
	}
	if alternativeCount == 0 {
		return true, nil
	}
	alternatives, _ := requirements.Elements()

	remainingSchemes := limits.MaxSchemes
	remainingScopes := limits.MaxScopes
	remainingRequirementBytes := MaxRequirementBytes
	anonymous := false
	requirementIndex := make([]requirementScheme, 0)
	alternativeEnds := make([]int, 0, len(alternatives))
	for _, alternative := range alternatives {
		if alternative.Kind() != jsonvalue.ObjectKind {
			return false, fmt.Errorf("%w: expected requirement object", ErrInvalidRequirements)
		}
		schemeCount, _ := alternative.Length()
		if !consumeBudget(&remainingSchemes, schemeCount) {
			return false, ErrLimitExceeded
		}
		schemes, _ := alternative.Members()
		anonymous = anonymous || schemeCount == 0
		for _, scheme := range schemes {
			if !consumeBudget(&remainingRequirementBytes, len(scheme.Name)) {
				return false, ErrLimitExceeded
			}
			if scheme.Value.Kind() != jsonvalue.ArrayKind {
				return false, fmt.Errorf("%w: expected scope array", ErrInvalidRequirements)
			}
			scopeCount, _ := scheme.Value.Length()
			if !consumeBudget(&remainingScopes, scopeCount) {
				return false, ErrLimitExceeded
			}
			requiredScopes, _ := scheme.Value.Elements()
			for _, rawScope := range requiredScopes {
				scope, valid := rawScope.Text()
				if !valid {
					return false, fmt.Errorf("%w: scope must be a string", ErrInvalidRequirements)
				}
				if !consumeBudget(&remainingRequirementBytes, len(scope)) {
					return false, ErrLimitExceeded
				}
			}
			requirementIndex = append(requirementIndex, requirementScheme{
				name: scheme.Name, scopes: requiredScopes,
			})
		}
		alternativeEnds = append(alternativeEnds, len(requirementIndex))
	}
	if anonymous {
		return true, nil
	}

	credentialIndex, err := indexCredentials(credentials, limits)
	if err != nil {
		return false, err
	}
	alternativeStart := 0
	for _, alternativeEnd := range alternativeEnds {
		matched := true
		for _, scheme := range requirementIndex[alternativeStart:alternativeEnd] {
			grantedSet, available := credentialIndex[scheme.name]
			if !available {
				matched = false
			}
			for _, rawScope := range scheme.scopes {
				scope, _ := rawScope.Text()
				if _, exists := grantedSet[scope]; !exists {
					matched = false
				}
			}
		}
		if matched {
			return true, nil
		}
		alternativeStart = alternativeEnd
	}
	return false, nil
}

func indexCredentials(credentials Credentials, limits Limits) (map[string]map[string]struct{}, error) {
	if len(credentials) > limits.MaxSchemes {
		return nil, ErrLimitExceeded
	}
	remainingScopes := limits.MaxScopes
	remainingBytes := MaxCredentialBytes
	for scheme, scopes := range credentials {
		if !consumeBudget(&remainingScopes, len(scopes)) {
			return nil, ErrLimitExceeded
		}
		if !consumeBudget(&remainingBytes, len(scheme)) {
			return nil, ErrLimitExceeded
		}
		for _, scope := range scopes {
			if !consumeBudget(&remainingBytes, len(scope)) {
				return nil, ErrLimitExceeded
			}
		}
	}

	index := make(map[string]map[string]struct{}, len(credentials))
	for scheme, scopes := range credentials {
		granted := make(map[string]struct{}, len(scopes))
		for _, scope := range scopes {
			granted[scope] = struct{}{}
		}
		index[scheme] = granted
	}
	return index, nil
}

// consumeBudget compares nonnegative lengths with a nonnegative remaining
// budget before subtracting, avoiding overflow-prone cumulative addition.
func consumeBudget(remaining *int, amount int) bool {
	if amount > *remaining {
		return false
	}
	*remaining -= amount
	return true
}

func effectiveLimits(limits Limits) (Limits, bool) {
	if limits.MaxAlternatives < 0 || limits.MaxSchemes < 0 || limits.MaxScopes < 0 {
		return Limits{}, false
	}
	defaults := DefaultLimits()
	if limits.MaxAlternatives == 0 {
		limits.MaxAlternatives = defaults.MaxAlternatives
	}
	if limits.MaxSchemes == 0 {
		limits.MaxSchemes = defaults.MaxSchemes
	}
	if limits.MaxScopes == 0 {
		limits.MaxScopes = defaults.MaxScopes
	}
	return limits, true
}
