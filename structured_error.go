/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package klog

import "slices"

// ErrorDetailer provides additional information about an error.
// When an error value implements this additional interface,
// the result of ErrorDetails will be logged in a separate key/value
// pair. The result of Error is logged as usual.
//
// In Kubernetes, text and JSON output backends (aka klog and zapr)
// will support this with "<error key>Details" (typically "errDetails")
// as key for the additional value.
//
// Other backends might not support this, so all relevant information
// should be in the error string.
type ErrorDetailer interface {
	ErrorDetails() any
}

// ErrorWithDetails adds additional details to an error for logging.
// If the base error already has such additional details, they
// will be included in a list of details.
//
// A [PseudoStruct] can be used to log some key/value pairs as
// if they were in a struct, without having to define such a struct.
// The formatting may be nicer, too.
func ErrorWithDetails(err error, details any) error {
	// This could be implemented as ErrorWithDetailsFunc(err, func() { return details }),
	// but having the details visible in the error instance may be more useful for
	// interactive debugging.
	return &errWithDetails{err, details}
}

type errWithDetails struct {
	error
	details any
}

var _ error = &errWithDetails{}
var _ ErrorDetailer = &errWithDetails{}

func (err *errWithDetails) ErrorDetails() any {
	if base, ok := err.error.(ErrorDetailer); ok {
		baseDetails := base.ErrorDetails()
		if baseDetailsList, ok := baseDetails.([]any); ok {
			// Flatten the list.
			return append(slices.Clone(baseDetailsList), err.details)
		}
		// Use a pair of values in a slice which gets detected above when nesting multiple times.
		return []any{baseDetails, err.details}
	}
	return err.details
}

// ErrorWithDetailsFunc adds additional details to an error for logging.
// In contrast to [ErrorWithDetails], the additional details are provided
// by the given function, which will be called only when needed. This
// can be used to avoid building some potentially expensive data structure
// that will not be needed when the error does not get logged.
//
// If the base error already has such additional details, they
// will be included in a list of details.
//
// A [PseudoStruct] can be used to log some key/value pairs as
// if they were in a struct, without having to define such a struct.
// The formatting may be nicer, too.
func ErrorWithDetailsFunc(err error, details func() any) error {
	return &errWithDetailsFunc{err, details}
}

type errWithDetailsFunc struct {
	error
	details func() any
}

var _ error = &errWithDetailsFunc{}
var _ ErrorDetailer = &errWithDetailsFunc{}

func (err *errWithDetailsFunc) ErrorDetails() any {
	if base, ok := err.error.(ErrorDetailer); ok {
		baseDetails := base.ErrorDetails()
		if baseDetailsList, ok := baseDetails.([]any); ok {
			// Flatten the list.
			return append(slices.Clone(baseDetailsList), err.details())
		}
		// Use a pair of values in a slice which gets detected above when nesting multiple times.
		return []any{baseDetails, err.details()}
	}
	return err.details()
}
