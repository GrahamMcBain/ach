// Copyright 2017 The ACH Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package ach

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Converters handles golang to ACH type Converters
type Converters struct{}

func (c *Converters) parseNumField(r string) (s int) {
	s, err := strconv.Atoi(strings.TrimSpace(r))
	if err != nil {
		// TODO: This is horrible
		fmt.Printf("%v", err)
		return
	}
	return s
}

// formatSimpleDate takes a time.Time and returns a string of YYMMDD
func (c *Converters) formatSimpleDate(t time.Time) string {
	return t.Format("060102")
}

// parseSimpleDate returns a time.Time when passed time as YYMMDD
func (c *Converters) parseSimpleDate(s string) time.Time {
	t, _ := time.Parse("060102", s)
	return t
}

// formatSimpleTime returns a string of HHMM when  passed a time.Time
func (c *Converters) formatSimpleTime(t time.Time) string {
	return t.Format("1504")
}

// parseSimpleTime returns a time.Time when passed a string of HHMM
func (c *Converters) parseSimpleTime(s string) time.Time {
	t, _ := time.Parse("1504", s)
	return t
}

//func (v *Converters) numericField()

// alphaField Alphanumeric and Alphabetic fields are left-justified and space filled.
func (c *Converters) alphaField(s string, max uint) string {
	ln := uint(len(s))
	if ln > max {
		return s[:max]
	}
	s += strings.Repeat(" ", int(max-ln))
	return s
}

// numericField right-justified, unisigned, and zero filled
func (c *Converters) numericField(n int, max uint) string {
	// @TODO remove decimel space from amount int

	s := strconv.Itoa(n)
	ln := uint(len(s))
	if ln > max {
		return s[ln-max:]
	}
	s = strings.Repeat("0", int(max-ln)) + s
	return s
}
