# Problem Details

This package implements [RFC 7807](https://datatracker.ietf.org/doc/html/rfc7807) (Problem Details for HTTP APIs) for Go. 
It provides an idiomatic way to use RFC 7807 in Go and offers both JSON and XML writers.

## Installation

```bash
go get github.com/josestg/problemdetail
```

## Examples

### Problem Details as an Error

```go
const (
	TypOutOfCredit     = "https://example.com/probs/out-of-credit"
	TypProductNotFound = "https://example.com/probs/product-not-found"
)

func service() error {
	// do something...

	// simulate 30% error rate for each type of error.
	n := rand.Intn(9)
	if n < 3 {
		return problemdetail.New(TypOutOfCredit,
			problemdetail.WithValidateLevel(problemdetail.LStandard),
			problemdetail.WithTitle("You do not have enough credit."),
			problemdetail.WithDetail("Your current balance is 30, but that costs 50."),
		)
	}

	if n < 6 {
		return problemdetail.New(TypProductNotFound,
			problemdetail.WithValidateLevel(problemdetail.LStandard),
			problemdetail.WithTitle("The product was not found."),
			problemdetail.WithDetail("The product you requested was not found in the system."),
		)
	}

	return errors.New("unknown error")
}

// handler is a sample handler for HTTP server.
// you can make this as a centralized error handler middleware.
func handler(w http.ResponseWriter, _ *http.Request) {
	err := service()
	if err != nil {
		// read the error as problemdetail.ProblemDetailer.
		var pd problemdetail.ProblemDetailer
		if !errors.As(err, &pd) {
			// untyped error for generic error handling.
			untyped := problemdetail.New(
				problemdetail.Untyped,
				problemdetail.WithValidateLevel(problemdetail.LStandard),
			)
			_ = problemdetail.WriteJSON(w, untyped, http.StatusInternalServerError)
			return
		}

		// typed error for specific error handling.
		switch pd.Kind() {
		case TypOutOfCredit:
			_ = problemdetail.WriteJSON(w, pd, http.StatusForbidden)
			// or problemdetail.WriteXML(w, pd, http.StatusForbidden) for XML
		case TypProductNotFound:
			_ = problemdetail.WriteJSON(w, pd, http.StatusNotFound)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### Problem Details with Extensions

```go
const (
	TypOutOfCredit     = "https://example.com/probs/out-of-credit"
	TypProductNotFound = "https://example.com/probs/product-not-found"
)

// BalanceProblemDetail is a sample problem detail with extension by embedding ProblemDetail.
type BalanceProblemDetail struct {
	*problemdetail.ProblemDetail
	Balance  int64    `json:"balance" xml:"balance"`
	Accounts []string `json:"accounts" xml:"accounts"`
}

func service() error {
	// do something...

	// simulate 30% error rate for each type of error.
	n := rand.Intn(9)
	if n < 3 {
		pd := problemdetail.New(TypOutOfCredit,
			problemdetail.WithValidateLevel(problemdetail.LStandard),
			problemdetail.WithTitle("You do not have enough credit."),
			problemdetail.WithDetail("Your current balance is 30, but that costs 50."),
		)
		return &BalanceProblemDetail{
			ProblemDetail: pd,
			Balance:       30,
			Accounts:      []string{"/account/12345", "/account/67890"},
		}
	}

	if n < 6 {
		return problemdetail.New(TypProductNotFound,
			problemdetail.WithValidateLevel(problemdetail.LStandard),
			problemdetail.WithTitle("The product was not found."),
			problemdetail.WithDetail("The product you requested was not found in the system."),
		)
	}

	return errors.New("unknown error")
}

// handler is a sample handler for HTTP server.
// you can make this as a centralized error handler middleware.
func handler(w http.ResponseWriter, _ *http.Request) {
	err := service()
	if err != nil {
		// read the error as problemdetail.ProblemDetailer.
		var pd problemdetail.ProblemDetailer
		if !errors.As(err, &pd) {
			// untyped error for generic error handling.
			untyped := problemdetail.New(
				problemdetail.Untyped,
				problemdetail.WithValidateLevel(problemdetail.LStandard),
			)
			_ = problemdetail.WriteJSON(w, untyped, http.StatusInternalServerError)
			return
		}

		// typed error for specific error handling.
		switch pd.Kind() {
		case TypOutOfCredit:
			_ = problemdetail.WriteJSON(w, pd, http.StatusForbidden)
			// or problemdetail.WriteXML(w, pd, http.StatusForbidden) for XML
		case TypProductNotFound:
			_ = problemdetail.WriteJSON(w, pd, http.StatusNotFound)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```