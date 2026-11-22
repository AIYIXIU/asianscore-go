# asianscore-go

[![Go Reference](https://pkg.go.dev/badge/github.com/AIYIXIU/asianscore-go.svg)](https://pkg.go.dev/github.com/AIYIXIU/asianscore-go)

Zero-dependency Go client and sports quantitative modeling library for **[AsianScore](https://asianscore.com/)** — real-time Football & NBA live scores, Dixon-Coles xG probabilities, and pre-kickoff SHA-256 verifiable prediction ledger.

### Links

- **Homepage:** https://asianscore.com/
- **Live Scoreboard:** https://asianscore.com/live.html
- **SHA-256 Verifiable AI Prediction Ledger:** https://asianscore.com/ai.html
- **Interactive Web Quant Tools:** https://asianscore.com/tools.html
- **Open Dataset (CC BY 4.0):** https://github.com/AIYIXIU/asianscore-xg-ledger-dataset
- **Shadow DOM Live Widget:** https://aiyixiu.github.io/asianscore-live-widget/

### Install

```bash
go get github.com/AIYIXIU/asianscore-go@v1.0.0
```

### Example

```go
summary, err := asianscore.FetchSummary(context.Background(), nil)
fair, _ := asianscore.DevigPower([]float64{1.95, 3.50, 4.10})
probs := asianscore.DixonColes(1.65, 1.05, -0.05)
stake := asianscore.KellyFraction(probs.Home, 1.95, 0.25)
```

## License

MIT
