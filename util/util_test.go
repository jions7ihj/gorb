/*
   Copyright (c) 2015 Andrey Sibiryov <me@kobology.ru>
   Copyright (c) 2015 Other contributors as noted in the AUTHORS file.

   This file is part of GORB - Go Routing and Balancing.

   GORB is free software; you can redistribute it and/or modify
   it under the terms of the GNU Lesser General Public License as published by
   the Free Software Foundation; either version 3 of the License, or
   (at your option) any later version.

   GORB is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
   GNU Lesser General Public License for more details.

   You should have received a copy of the GNU Lesser General Public License
   along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package util

import (
	"encoding/json"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMustMarshall(t *testing.T) {
	input := struct {
		Answer int `json:"answer"`
	}{
		Answer: 42,
	}

	rv := map[JSONOptions][]byte{}

	normal := JSONOptions{}
	indent := JSONOptions{Indent: true}

	rv[normal], _ = json.Marshal(&input)
	rv[indent], _ = json.MarshalIndent(&input, "", "\t")

	tests := []struct {
		in   interface{}
		opts JSONOptions
		rv   []byte
	}{
		{in: input, opts: normal, rv: rv[normal]},
		{in: input, opts: indent, rv: rv[indent]},
	}

	for _, test := range tests {
		assert.Equal(t, test.rv, MustMarshal(test.in, test.opts))
	}
}

func TestMustMarshalPanic(t *testing.T) {
	assert.Panics(t, func() {
		// Map key type is not string.
		MustMarshal(map[int]int{}, JSONOptions{})
	})
}

func TestAddrFamily(t *testing.T) {
	tests := []struct {
		in net.IP
		rv int
	}{
		{in: net.ParseIP("10.0.0.1"), rv: IPv4},
		{in: net.ParseIP("fd11:bcb5:61df::1"), rv: IPv6},
		{in: net.ParseIP("10.0.0.1").To4(), rv: IPv4},
	}

	for _, test := range tests {
		assert.Equal(t, test.rv, AddrFamily(test.in))
	}
}

func TestParseInterval(t *testing.T) {
	tests := []struct {
		in string
		rv time.Duration
	}{
		{in: "600s", rv: 600 * time.Second},
		{in: "2m", rv: 2 * time.Minute},
		{in: "24h", rv: 24 * time.Hour},
	}

	for _, test := range tests {
		rv, err := ParseInterval(test.in)

		require.NoError(t, err)
		assert.Equal(t, test.rv, rv)
	}
}

func TestParseIntervalErrors(t *testing.T) {
	tests := []struct {
		in  string
		err error
	}{
		// Missing unit.
		{in: "600", err: errInvalidIntervalFormat},
		// Missing number.
		{in: "foos", err: errInvalidIntervalFormat},
		// Unknown unit.
		{in: "24z", err: errInvalidIntervalFormat},
	}

	for _, test := range tests {
		_, err := ParseInterval(test.in)

		require.Error(t, err)
		assert.Equal(t, test.err, err)
	}
}
