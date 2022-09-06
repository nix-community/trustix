// Copyright (C) 2022 adisbladis
// Copyright © 2020-2022 The Trustix Authors
//
// SPDX-License-Identifier: GPL-3.0-only

package lib

type Result[T any] struct {
	value T
	err   error
}

func NewResult[T any](value T, err error) *Result[T] {
	return &Result[T]{
		value: value,
		err:   err,
	}
}

func (r *Result[T]) Unwrap() (T, error) {
	return r.value, r.err
}
