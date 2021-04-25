// Copyright (C) 2021 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package config

type GRPCTransport struct {
	Remote string `toml:"remote"`
}

type Transport struct {
	Type string         `toml:"type"`
	GRPC *GRPCTransport `toml:"grpc"`
}

type Subscriber struct {
	Transport *Transport        `toml:"transport"`
	PublicKey *PublicKey        `toml:"key"`
	Meta      map[string]string `toml:"meta"`
}

func (s *Subscriber) Validate() error {
	return nil
}