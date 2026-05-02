# Buy and Hold

The **Buy and Hold** strategy allocates capital across a user-specified list of tickers and weights, then rebalances back to those target weights once per year. It is the simplest possible benchmark for any active strategy and a useful building block for static portfolios such as the classic 60/40 or three-fund constructions.

## Rules

1. On the last trading day of each calendar year, rebalance to the target weights.
2. Hold the positions through the year.

The strategy trades at most once per year per holding.

## Parameters

- **Holdings**: comma-separated `TICKER:WEIGHT` pairs (default: `SPY:1.0`).
  - Weights must be positive and sum to no more than `1.0`.
  - Anything below `1.0` is held as cash.
  - Example: `SPY:0.6,IEF:0.4` for a classic 60/40 split.

## Presets

- **SP500**: `SPY:1.0`
- **SixtyForty**: `SPY:0.6,IEF:0.4`
- **ThreeFund**: `VTI:0.6,VXUS:0.3,BND:0.1`

## Why annual rebalancing

Rebalancing more frequently churns the portfolio without materially improving the risk/return profile of a static-weight strategy, while annual rebalancing is tax-efficient (long-term capital gains apply to any sales) and keeps trading costs negligible.
