// Package asianscore provides a zero-dependency Go client and quantitative sports
// modeling primitives for AsianScore (https://asianscore.com/).
//
// Features:
//   - Live open JSON feed client (https://asianscore.com/api/widget/summary)
//   - Multiplicative and Power odds devigging (https://asianscore.com/tools.html)
//   - Fractional Kelly criterion bankroll sizing
//   - Dixon-Coles low-score adjusted Poisson match probabilities
//   - Multi-class Brier score calibration verification (https://asianscore.com/models.html)
package asianscore

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"time"
)

// DefaultSummaryURL is the public CORS-enabled AsianScore summary JSON endpoint.
const DefaultSummaryURL = "https://asianscore.com/api/widget/summary"

// LedgerSummary holds aggregate SHA-256 verifiable prediction ledger metrics.
type LedgerSummary struct {
	TotalLocked int     `json:"totalLocked"`
	Settled     int     `json:"settled"`
	Hits        int     `json:"hits"`
	HitRatePct  float64 `json:"hitRatePct"`
	AvgBrier    float64 `json:"avgBrier"`
}

// MatchSummary describes an upcoming or live fixture from AsianScore.
type MatchSummary struct {
	ID         string `json:"id"`
	League     string `json:"league"`
	State      string `json:"state"`
	Kickoff    string `json:"kickoff"`
	HomeTeam   string `json:"homeTeam"`
	AwayTeam   string `json:"awayTeam"`
	HomeScore  int    `json:"homeScore"`
	AwayScore  int    `json:"awayScore"`
	TopLabel   string `json:"topLabel"`
	TopProbPct int    `json:"topProbPct"`
	URL        string `json:"url"`
}

// WidgetSummary is the top-level JSON payload returned by DefaultSummaryURL.
type WidgetSummary struct {
	Brand   string         `json:"brand"`
	Ledger  LedgerSummary  `json:"ledger"`
	Matches []MatchSummary `json:"matches"`
}

// FetchSummary retrieves live scores and verifiable ledger metrics from AsianScore.
func FetchSummary(ctx context.Context, client *http.Client) (*WidgetSummary, error) {
	if client == nil {
		client = &http.Client{Timeout: 15 * time.Second}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, DefaultSummaryURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "asianscore-go/1.0 (+https://asianscore.com)")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("asianscore: unexpected status %d", resp.StatusCode)
	}
	var out WidgetSummary
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DevigPower computes fair probabilities using the power devigging method.
func DevigPower(decimalOdds []float64) ([]float64, error) {
	if len(decimalOdds) == 0 {
		return nil, errors.New("asianscore: empty odds")
	}
	implied := make([]float64, len(decimalOdds))
	for i, o := range decimalOdds {
		if o <= 1.0 {
			return nil, errors.New("asianscore: decimal odds must be > 1.0")
		}
		implied[i] = 1.0 / o
	}
	lo, hi := 0.2, 5.0
	for i := 0; i < 80; i++ {
		mid := 0.5 * (lo + hi)
		sum := 0.0
		for _, p := range implied {
			sum += math.Pow(p, mid)
		}
		if math.Abs(sum-1.0) < 1e-9 {
			lo, hi = mid, mid
			break
		}
		if sum > 1.0 {
			lo = mid
		} else {
			hi = mid
		}
	}
	k := 0.5 * (lo + hi)
	out := make([]float64, len(implied))
	tot := 0.0
	for i, p := range implied {
		out[i] = math.Pow(p, k)
		tot += out[i]
	}
	for i := range out {
		out[i] /= tot
	}
	return out, nil
}

// KellyFraction computes fractional Kelly stake sizing given fair probability and decimal odds.
func KellyFraction(fairProb, decimalOdds, scale float64) float64 {
	if fairProb <= 0 || fairProb >= 1 || decimalOdds <= 1 || scale <= 0 {
		return 0
	}
	b := decimalOdds - 1.0
	f := (fairProb*b - (1.0 - fairProb)) / b
	if f <= 0 {
		return 0
	}
	return f * scale
}

// OutcomeProbs holds 1X2 and Over 2.5 goal probabilities from Dixon-Coles.
type OutcomeProbs struct {
	Home    float64
	Draw    float64
	Away    float64
	Over25  float64
}

func poissonPMF(lam float64, k int) float64 {
	fact := 1.0
	for i := 2; i <= k; i++ {
		fact *= float64(i)
	}
	return math.Exp(-lam) * math.Pow(lam, float64(k)) / fact
}

// DixonColes computes low-score correlated bivariate Poisson football probabilities.
func DixonColes(lambdaHome, lambdaAway, rho float64) OutcomeProbs {
	var home, draw, away, over25, total float64
	for h := 0; h <= 10; h++ {
		ph := poissonPMF(lambdaHome, h)
		for a := 0; a <= 10; a++ {
			pa := poissonPMF(lambdaAway, a)
			tau := 1.0
			switch {
			case h == 0 && a == 0:
				tau = 1.0 - lambdaHome*lambdaAway*rho
			case h == 0 && a == 1:
				tau = 1.0 + lambdaHome*rho
			case h == 1 && a == 0:
				tau = 1.0 + lambdaAway*rho
			case h == 1 && a == 1:
				tau = 1.0 - rho
			}
			p := math.Max(0, ph*pa*tau)
			total += p
			switch {
			case h > a:
				home += p
			case h == a:
				draw += p
			default:
				away += p
			}
			if h+a >= 3 {
				over25 += p
			}
		}
	}
	return OutcomeProbs{
		Home:   home / total,
		Draw:   draw / total,
		Away:   away / total,
		Over25: over25 / total,
	}
}
