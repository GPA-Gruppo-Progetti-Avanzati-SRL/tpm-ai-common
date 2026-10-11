package client

import (
	"fmt"
	"strings"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/rs/zerolog"
)

// UsageAccumulator sums anthropic.Usage values across any number of API calls.
//
// It is the single place that knows how to (a) fold usage together field by
// field — every counter, not just input/output — and (b) render it, both as a
// human string and as structured zerolog fields. Because it holds whatever has
// been added to it, the same formatting serves one call or a whole run:
//
//	// aggregate across a run
//	var acc client.UsageAccumulator
//	for _, r := range results {
//	    acc.Add(r.Usage)
//	}
//	log.Info().Object("usage", &acc).Msg("run complete") // structured
//	fmt.Println(acc.String())                             // human
//
//	// log a single response — fold one value, then render
//	var one client.UsageAccumulator
//	one.Add(resp.Usage)
//	log.Info().Object("usage", &one).Msg("request done")
//
// The zero value is ready to use. It is not safe for concurrent Add calls;
// guard it with a mutex if multiple goroutines report into one accumulator.
type UsageAccumulator struct {
	total anthropic.Usage
	calls int
}

// Add folds one Usage value into the running total, summing every token and
// request counter the SDK reports — input/output, both cache columns and their
// 1h/5m TTL buckets, thinking tokens, and server-tool (web fetch/search)
// requests. The metadata fields (InferenceGeo, ServiceTier) are not counters
// and are left out of the total.
func (a *UsageAccumulator) Add(u anthropic.Usage) {
	a.calls++
	a.total.InputTokens += u.InputTokens
	a.total.OutputTokens += u.OutputTokens
	a.total.CacheReadInputTokens += u.CacheReadInputTokens
	a.total.CacheCreationInputTokens += u.CacheCreationInputTokens
	a.total.CacheCreation.Ephemeral1hInputTokens += u.CacheCreation.Ephemeral1hInputTokens
	a.total.CacheCreation.Ephemeral5mInputTokens += u.CacheCreation.Ephemeral5mInputTokens
	a.total.OutputTokensDetails.ThinkingTokens += u.OutputTokensDetails.ThinkingTokens
	a.total.ServerToolUse.WebFetchRequests += u.ServerToolUse.WebFetchRequests
	a.total.ServerToolUse.WebSearchRequests += u.ServerToolUse.WebSearchRequests
}

// Total returns the accumulated usage. The counter fields are populated; the
// metadata fields (InferenceGeo, ServiceTier) are zero since Add does not
// carry them.
func (a *UsageAccumulator) Total() anthropic.Usage { return a.total }

// Calls reports how many times Add has been called — i.e. how many API
// responses were folded into the total.
func (a *UsageAccumulator) Calls() int { return a.calls }

// String renders the accumulated usage as a compact single line. Input and
// output are always shown; the cache, thinking, and server-tool groups appear
// only when non-zero, so a plain call without caching or thinking stays short.
//
//	in=1234 out=567 cache_r=890 cache_w=12 think=42 web_fetch=1
func (a *UsageAccumulator) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "in=%d out=%d", a.total.InputTokens, a.total.OutputTokens)

	if a.total.CacheReadInputTokens > 0 || a.total.CacheCreationInputTokens > 0 {
		fmt.Fprintf(&b, " cache_r=%d cache_w=%d", a.total.CacheReadInputTokens, a.total.CacheCreationInputTokens)
	}
	if t := a.total.OutputTokensDetails.ThinkingTokens; t > 0 {
		fmt.Fprintf(&b, " think=%d", t)
	}
	if r := a.total.ServerToolUse.WebFetchRequests; r > 0 {
		fmt.Fprintf(&b, " web_fetch=%d", r)
	}
	if r := a.total.ServerToolUse.WebSearchRequests; r > 0 {
		fmt.Fprintf(&b, " web_search=%d", r)
	}
	return b.String()
}

// MarshalZerologObject implements zerolog.LogObjectMarshaler, so an accumulator
// can be logged as a structured, queryable sub-object:
//
//	log.Info().Object("usage", &acc).Msg("run complete")
//	// {"usage":{"input":1234,"output":567,"cache_read":890,...}}
//
// Every counter is emitted (including zeros) so log queries can rely on the
// fields being present; calls carries how many responses were folded in.
func (a *UsageAccumulator) MarshalZerologObject(e *zerolog.Event) {
	e.Int("calls", a.calls).
		Int64("input", a.total.InputTokens).
		Int64("output", a.total.OutputTokens).
		Int64("cache_read", a.total.CacheReadInputTokens).
		Int64("cache_creation", a.total.CacheCreationInputTokens).
		Int64("cache_1h", a.total.CacheCreation.Ephemeral1hInputTokens).
		Int64("cache_5m", a.total.CacheCreation.Ephemeral5mInputTokens).
		Int64("thinking", a.total.OutputTokensDetails.ThinkingTokens).
		Int64("web_fetch", a.total.ServerToolUse.WebFetchRequests).
		Int64("web_search", a.total.ServerToolUse.WebSearchRequests)
}
