package websitepoller

import (
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	randomdata "github.com/Pallinder/go-randomdata"
)

func parseHTTPMethod(method *string) (string, error) {
	m := "GET"
	if method != nil {
		if len(*method) > 0 {
			m = *method
		}
	}

	switch m := strings.ToUpper(m); m {
	case http.MethodGet:
		return http.MethodGet, nil
	case http.MethodHead:
		return http.MethodHead, nil
	case http.MethodPost:
		return http.MethodPost, nil
	case http.MethodPut:
		return http.MethodPut, nil
	case http.MethodPatch:
		return http.MethodPatch, nil
	case http.MethodDelete:
		return http.MethodDelete, nil
	case http.MethodConnect:
		return http.MethodConnect, nil
	case http.MethodOptions:
		return http.MethodOptions, nil
	case http.MethodTrace:
		return http.MethodTrace, nil
	default:
		return "", ErrUnrecognizedHTTPMethod
	}
}

func parseURL(rawurl string) (*url.URL, error) {
	parsed, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}

	if !parsed.IsAbs() {
		return nil, ErrURLNoScheme
	}

	return parsed, nil
}

func parsePollOptions(id string, opts *PollOptions) (randFreq bool, freq int, offset int) {
	l := log.With().Str("id", id).Logger()
	randFreq, freq, offset = false, defaultFrequency, 0

	if opts == nil {
		l.Debug().Msg("no poll options, returning default values...")
		return
	}

	if opts.Frequency >= minFrequency {
		freq = opts.Frequency
	} else {
		l.Error().Int("frequency", freq).Int("default", defaultFrequency).Msg("invalid frequency provided, using default value...")
	}

	if !opts.RandomFrequency {
		return
	}

	randFreq = true
	offset = defaultOffsetRange
	if opts.OffsetRange != nil {
		if *opts.OffsetRange >= minOffset {
			offset = *opts.OffsetRange
		} else {
			l.Warn().Int("offset", *opts.OffsetRange).Int("default", defaultOffsetRange).Msg("invalid offset range provided, reverting to default...")
		}
	}

	if freq-offset >= minFrequency {
		return
	}

	l.Warn().Int("range", offset).Int("frequency", freq).Msg("offset is too low, reverting to default...")
	offset = defaultOffsetRange
	freq = defaultFrequency
	return
}

func parseUserAgentOptions(id string, opts *UserAgentOptions) (randUA bool, uas []string) {
	l := log.With().Str("id", id).Logger()
	randUA, uas = false, []string{}

	if opts == nil {
		l.Warn().Msg("no user agents provided, you should provide at least one or enable random user agents")
		return
	}

	randUA = opts.RandomUA
	if len(opts.UserAgents) > 0 {
		uas = opts.UserAgents
		return
	}

	if randUA {
		l.Debug().Msg("random user agents will be used")
		return
	}

	l.Warn().Msg("no user agents provided, you should provide at least one or enable random user agents")
	return
}

func nextRandomTick(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Intn(max-min)
}

func getNextUA(id string, userAgents []string, random bool, last int) (ua string, index int) {
	l := log.With().Str("id", id).Logger()
	ua, index = "", -1

	if len(userAgents) == 0 {
		if random {
			ua = randomdata.UserAgentString()
			l.Debug().Str("user-agent", ua).Msg("generated random user agent")
		}

		return
	}

	length := len(userAgents)
	index = last
	if !random {
		index = (index + 1) % length
		ua = userAgents[index]
		l.Debug().Str("user-agent", ua).Int("index", index).Msg("rotated user agent")
	} else {
		rand.Seed(time.Now().UnixNano())
		for index == last {
			index = rand.Intn(length - 1)
		}
		ua = userAgents[index]
		l.Debug().Str("user-agent", ua).Int("index", index).Msg("picked random user agent")
	}

	return
}
