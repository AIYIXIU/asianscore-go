# asianscore-go

[![Go Reference](https://pkg.go.dev/badge/github.com/AIYIXIU/asianscore-go.svg)](https://pkg.go.dev/github.com/AIYIXIU/asianscore-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://github.com/AIYIXIU/asianscore-go/blob/main/LICENSE)

Zero-dependency Go client and sports quantitative modeling library for **[AsianScore](https://asianscore.com/)** — real-time Football & NBA live scores, Dixon-Coles xG probabilities, Power odds devigging, fractional Kelly sizing, and pre-kickoff SHA-256 verifiable prediction ledger verification.

## Official Links

- **Homepage:** https://asianscore.com/
- **Live Football & NBA Scoreboard:** https://asianscore.com/scores.html
- **SHA-256 Verifiable AI Prediction Ledger:** https://asianscore.com/models.html
- **Interactive Web Quant Tools:** https://asianscore.com/tools.html
- **Public REST API Reference (Read the Docs):** https://asianscore-quant.readthedocs.io/en/latest/rest-api/
- **Open Dataset (CC BY 4.0 / Zenodo DOI):** https://doi.org/10.5281/zenodo.24441036
- **Companion SDKs:** [Python (PyPI)](https://pypi.org/project/asianscore-quant/) · [Ruby (RubyGems)](https://rubygems.org/gems/asianscore_quant) · [Web Component (npm)](https://www.npmjs.com/package/asianscore-widget)

## Install

```bash
go get github.com/AIYIXIU/asianscore-go@v1.0.1
```

## Quickstart

```go
package main

import (
    "context"
    "fmt"
    asianscore "github.com/AIYIXIU/asianscore-go"
)

func main() {
    summary, err := asianscore.FetchSummary(context.Background(), nil)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Locked=%d Settled=%d AvgBrier=%.4f\n",
        summary.Ledger.TotalLocked, summary.Ledger.Settled, summary.Ledger.AvgBrier)

    fair, _ := asianscore.DevigPower([]float64{1.95, 3.50, 4.10})
    probs := asianscore.DixonColes(1.65, 1.05, -0.05)
    stake := asianscore.KellyFraction(probs.Home, 1.95, 0.25)
    fmt.Printf("Fair=%v HomeWin=%.3f Kelly=%.3f\n", fair, probs.Home, stake)
}
```

## License

MIT License — see [LICENSE](LICENSE).
