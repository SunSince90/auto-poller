package websitepoller

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseHTTPMethod(t *testing.T) {
	a := assert.New(t)

	methodEmpty := ""
	unrecognized := "test"
	methodGet := "head"

	cases := []struct {
		method    *string
		expMethod string
		expErr    error
	}{
		{
			expMethod: http.MethodGet,
		},
		{
			method:    &methodEmpty,
			expMethod: http.MethodGet,
		},
		{
			method:    &unrecognized,
			expMethod: http.MethodGet,
			expErr:    ErrUnrecognizedHTTPMethod,
		},
		{
			method:    &methodGet,
			expMethod: http.MethodHead,
		},
	}

	for _, currCase := range cases {
		method, err := parseHTTPMethod(currCase.method)

		if currCase.expErr != nil {
			a.Error(err)
			a.Equal(currCase.expErr, err)
		} else {
			a.NoError(err)
			a.Equal(currCase.expMethod, method)
		}
	}

}

func TestParsePollOptions(t *testing.T) {
	a := assert.New(t)

	invalidOff := -10
	tooLowFreqRange := 24
	cases := []struct {
		arg     *PollOptions
		expRand bool
		expFreq int
		expOff  int
	}{
		{
			expRand: false,
			expFreq: defaultFrequency,
			expOff:  0,
		},
		{
			arg: &PollOptions{
				Frequency: -1,
			},
			expRand: false,
			expFreq: defaultFrequency,
			expOff:  0,
		},
		{
			arg: &PollOptions{
				Frequency: 2,
			},
			expRand: false,
			expFreq: defaultFrequency,
			expOff:  0,
		},
		{
			arg: &PollOptions{
				Frequency: 25,
			},
			expRand: false,
			expFreq: 25,
			expOff:  0,
		},
		{
			arg: &PollOptions{
				Frequency:       25,
				RandomFrequency: true,
			},
			expRand: true,
			expFreq: 25,
			expOff:  defaultOffsetRange,
		},
		{
			arg: &PollOptions{
				Frequency:       25,
				RandomFrequency: true,
				OffsetRange:     &invalidOff,
			},
			expRand: true,
			expFreq: 25,
			expOff:  defaultOffsetRange,
		},
		{
			arg: &PollOptions{
				Frequency:       25,
				RandomFrequency: true,
				OffsetRange:     &tooLowFreqRange,
			},
			expRand: true,
			expFreq: defaultFrequency,
			expOff:  defaultOffsetRange,
		},
	}

	for i, currCase := range cases {
		rand, freq, off := parsePollOptions("", currCase.arg)

		errRand := a.Equal(currCase.expRand, rand)
		errFreq := a.Equal(currCase.expFreq, freq)
		errOff := a.Equal(currCase.expOff, off)
		if !errRand || !errFreq || !errOff {
			a.FailNow(fmt.Sprintf("case %d failed", i))
		}
	}
}

func TestParseUserAgentOptions(t *testing.T) {
	a := assert.New(t)

	cases := []struct {
		arg     *UserAgentOptions
		expRand bool
		expUas  []string
	}{
		{
			expRand: false,
			expUas:  []string{},
		},
		{
			arg: &UserAgentOptions{
				UserAgents: []string{"one", "two"},
				RandomUA:   true,
			},
			expRand: true,
			expUas:  []string{"one", "two"},
		},
	}

	for i, currCase := range cases {
		rand, uas := parseUserAgentOptions("", currCase.arg)

		errRand := a.Equal(currCase.expRand, rand)
		errUAs := a.Equal(currCase.expUas, uas)
		if !errRand || !errUAs {
			a.FailNow(fmt.Sprintf("case %d failed", i))
		}
	}
}

func TestGetNextUA(t *testing.T) {
	a := assert.New(t)
	userAgents := []string{"zero", "one", "two", "three"}

	// Test rotate
	cases := []struct {
		uas      []string
		rand     bool
		last     int
		expUA    string
		expIndex int
	}{
		{
			last:     -1,
			expIndex: -1,
		},
		{
			last:     5,
			expIndex: -1,
		},
		// {
		// 	rand:     true,
		// 	last:     -1,
		// 	expIndex: -1,
		// },
		{
			uas:      userAgents,
			rand:     false,
			last:     -1,
			expUA:    "zero",
			expIndex: 0,
		},
		{
			uas:      userAgents,
			rand:     false,
			last:     0,
			expUA:    "one",
			expIndex: 1,
		},
		{
			uas:      userAgents,
			rand:     false,
			last:     3,
			expUA:    "zero",
			expIndex: 0,
		},
	}

	for i, currCase := range cases {
		ua, ind := getNextUA("", currCase.uas, currCase.rand, currCase.last)

		errUA := a.Equal(currCase.expUA, ua)
		errInd := a.Equal(currCase.expIndex, ind)
		if !errUA || !errInd {
			a.FailNow(fmt.Sprintf("case %d failed", i))
		}
	}

	// Test random
	ua, ind := getNextUA("", []string{}, true, -1)
	a.Equal(-1, ind)
	a.NotContains(userAgents, ua)

	last := 3
	ua, ind = getNextUA("", userAgents, true, last)
	a.NotEqual(-1, ind)
	a.NotEqual(last, ind)
	a.Contains(userAgents, ua)
}
