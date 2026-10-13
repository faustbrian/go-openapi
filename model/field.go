// Package model preserves absent, null, typed, and invalid field states and
// insertion-ordered collections for generated OpenAPI views. Field.Valid checks
// representation, not specification legality or requiredness.
//
// Generic values are copied shallowly: container structure is immutable, but
// mutable pointers, maps, or slices inside T remain caller-owned. Concurrent use
// requires immutable elements or caller synchronization. Built-in generated
// models use immutable JSON values and typed views as elements.
package model

import "github.com/faustbrian/go-openapi/v2/jsonvalue"

type fieldState uint8

const (
	fieldAbsent fieldState = iota
	fieldNull
	fieldPresent
	fieldInvalid
)

// Field preserves the distinction between an absent field, an explicit null,
// and a concrete value including its zero or empty value.
type Field[T any] struct {
	state fieldState
	value T
	raw   jsonvalue.Value
}

// Absent constructs an absent field.
func Absent[T any]() Field[T] {
	return Field[T]{state: fieldAbsent}
}

// Null constructs an explicitly null field.
func Null[T any]() Field[T] {
	return Field[T]{state: fieldNull}
}

// Present stores a shallow copy of value without applying defaults or cloning
// reference-valued contents. The caller retains responsibility for mutable T.
func Present[T any](value T) Field[T] {
	return Field[T]{state: fieldPresent, value: value}
}

// Invalid constructs a present, non-null field whose representation cannot be
// interpreted as T. The raw value remains available for lossless diagnostics.
func Invalid[T any](raw jsonvalue.Value) Field[T] {
	return Field[T]{state: fieldInvalid, raw: raw}
}

// Present reports whether the field occurred, including as explicit null.
func (field Field[T]) Present() bool {
	return field.state != fieldAbsent
}

// Null reports whether the field occurred as explicit null.
func (field Field[T]) Null() bool {
	return field.state == fieldNull
}

// Valid reports whether the representation can be interpreted as T. Absent
// and explicit-null fields are representation-valid; version validation
// decides whether null is legal at the field's location.
func (field Field[T]) Valid() bool {
	return field.state != fieldInvalid
}

// Value returns a shallow copy and true only for a present, non-null typed field.
// Absent, null, and invalid fields return the zero T and false.
func (field Field[T]) Value() (T, bool) {
	return field.value, field.state == fieldPresent
}

// Raw returns the lossless representation of an invalid typed field.
func (field Field[T]) Raw() (jsonvalue.Value, bool) {
	return field.raw, field.state == fieldInvalid
}
