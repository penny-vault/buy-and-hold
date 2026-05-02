// Copyright 2026
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	_ "embed"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/penny-vault/pvbt/asset"
	"github.com/penny-vault/pvbt/engine"
	"github.com/penny-vault/pvbt/portfolio"
)

//go:embed README.md
var description string

// BuyAndHold allocates to a user-specified set of tickers and weights and
// rebalances back to those target weights once per year.
type BuyAndHold struct {
	Holdings string `pvbt:"holdings" desc:"Comma-separated TICKER:WEIGHT pairs (e.g. SPY:0.6,QQQ:0.3,BND:0.1). Weights summing below 1.0 leave the remainder in cash." default:"SPY:1.0" suggest:"SP500=SPY:1.0|SixtyForty=SPY:0.6,IEF:0.4|ThreeFund=VTI:0.6,VXUS:0.3,BND:0.1"`

	parsed   []targetWeight
	parseErr error
}

type targetWeight struct {
	Ticker string
	Weight float64
}

func (s *BuyAndHold) Name() string { return "Buy and Hold" }

func (s *BuyAndHold) Setup(_ *engine.Engine) {
	s.parsed, s.parseErr = parseHoldings(s.Holdings)
}

func (s *BuyAndHold) Describe() engine.StrategyDescription {
	return engine.StrategyDescription{
		ShortCode:   "bnh",
		Description: description,
		Source:      "",
		Version:     "1.0.0",
		VersionDate: time.Date(2026, 5, 2, 0, 0, 0, 0, time.UTC),
		Schedule:    "@yearend",
		Benchmark:   "SPY",
	}
}

func (s *BuyAndHold) Compute(ctx context.Context, eng *engine.Engine, _ portfolio.Portfolio, batch *portfolio.Batch) error {
	if s.parseErr != nil {
		return s.parseErr
	}

	members := make(map[asset.Asset]float64, len(s.parsed))

	parts := make([]string, 0, len(s.parsed))
	for _, tw := range s.parsed {
		members[eng.Asset(tw.Ticker)] = tw.Weight
		parts = append(parts, fmt.Sprintf("%.1f%% %s", tw.Weight*100, tw.Ticker))
	}

	justification := "buy and hold: " + strings.Join(parts, ", ")
	batch.Annotate("justification", justification)

	allocation := portfolio.Allocation{
		Date:          eng.CurrentDate(),
		Members:       members,
		Justification: justification,
	}

	if err := batch.RebalanceTo(ctx, allocation); err != nil {
		return fmt.Errorf("rebalance failed: %w", err)
	}

	return nil
}

func parseHoldings(spec string) ([]targetWeight, error) {
	seen := make(map[string]bool)

	var (
		out   []targetWeight
		total float64
	)

	for _, pair := range strings.Split(spec, ",") {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}

		kv := strings.SplitN(pair, ":", 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid holdings entry %q (expected TICKER:WEIGHT)", pair)
		}

		ticker := strings.ToUpper(strings.TrimSpace(kv[0]))

		weight, err := strconv.ParseFloat(strings.TrimSpace(kv[1]), 64)
		if err != nil {
			return nil, fmt.Errorf("invalid weight in %q: %w", pair, err)
		}

		if ticker == "" {
			return nil, fmt.Errorf("empty ticker in %q", pair)
		}

		if weight <= 0 {
			return nil, fmt.Errorf("weight must be positive in %q", pair)
		}

		if seen[ticker] {
			return nil, fmt.Errorf("duplicate ticker %q", ticker)
		}

		seen[ticker] = true
		out = append(out, targetWeight{Ticker: ticker, Weight: weight})
		total += weight
	}

	if len(out) == 0 {
		return nil, fmt.Errorf("no holdings specified")
	}

	if total > 1.0001 {
		return nil, fmt.Errorf("weights sum to %.4f; long-only buy-and-hold requires sum <= 1.0", total)
	}

	sort.Slice(out, func(i, j int) bool { return out[i].Ticker < out[j].Ticker })

	return out, nil
}
