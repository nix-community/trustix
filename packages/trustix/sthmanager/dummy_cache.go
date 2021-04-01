// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package sthmanager

import (
	"fmt"

	"github.com/tweag/trustix/packages/trustix-proto/schema"
)

type dummySthCache struct {
	handler func() (*schema.STH, error)
}

func NewDummySTHCache(handler func() (*schema.STH, error)) STHCache {
	return &dummySthCache{
		handler: handler,
	}
}

func (c *dummySthCache) Set(sth *schema.STH) error {
	return fmt.Errorf("Operation not supported")
}

func (c *dummySthCache) Get() (*schema.STH, error) {
	return c.handler()
}

func (c *dummySthCache) Close() {}
