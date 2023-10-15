package problemdetail_test

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/josestg/problemdetail"
)

// BalanceProblemDetail is a sample problem detail with extension by embedding ProblemDetail.
type BalanceProblemDetail struct {
	*problemdetail.ProblemDetail
	Balance  int32    `json:"balance" xml:"balance"`
	Accounts []string `json:"accounts" xml:"accounts"`
}

func TestWriteJSON_WithExtension(t *testing.T) {
	data := BalanceProblemDetail{
		ProblemDetail: problemdetail.New(
			"https://example.com/probs/out-of-credit",
			problemdetail.WithDetail("Your current balance is 30, but that costs 50."),
			problemdetail.WithInstance("/account/12345/abc"),
			problemdetail.WithTitle("You do not have enough credit."),
		),
		Balance:  30,
		Accounts: []string{"/account/12345", "/account/67890"},
	}

	rec := httptest.NewRecorder()
	err := problemdetail.WriteJSON(rec, &data, 403)
	expectTrue(t, err == nil)

	expRaw := `{"type":"https://example.com/probs/out-of-credit","title":"You do not have enough credit.","status":403,"detail":"Your current balance is 30, but that costs 50.","instance":"/account/12345/abc","balance":30,"accounts":["/account/12345","/account/67890"]}`
	gotRaw := strings.TrimSpace(rec.Body.String())

	expectTrue(t, gotRaw == expRaw)
	expectTrue(t, rec.Code == 403)
	expectTrue(t, rec.Header().Get("Content-Type") == "application/problem+json; charset=utf-8")
}

func TestWriteJSON_WithUntyped(t *testing.T) {
	data := problemdetail.New(problemdetail.Untyped, problemdetail.WithValidateLevel(problemdetail.LStandard))

	rec := httptest.NewRecorder()
	err := problemdetail.WriteJSON(rec, data, 403)
	expectTrue(t, err == nil)

	expRaw := `{"type":"about:blank","title":"Forbidden","status":403}`
	gotRaw := strings.TrimSpace(rec.Body.String())

	expectTrue(t, gotRaw == expRaw)
	expectTrue(t, rec.Code == 403)
	expectTrue(t, rec.Header().Get("Content-Type") == "application/problem+json; charset=utf-8")
}

func TestWriteJSON_WithTyped(t *testing.T) {
	data := problemdetail.New(
		"https://example.com/probs/out-of-credit",
		problemdetail.WithDetail("Your current balance is 30, but that costs 50."),
		problemdetail.WithInstance("/account/12345/abc"),
		problemdetail.WithTitle("You do not have enough credit."),
	)

	rec := httptest.NewRecorder()
	err := problemdetail.WriteJSON(rec, data, 403)
	expectTrue(t, err == nil)

	expRaw := `{"type":"https://example.com/probs/out-of-credit","title":"You do not have enough credit.","status":403,"detail":"Your current balance is 30, but that costs 50.","instance":"/account/12345/abc"}`
	gotRaw := strings.TrimSpace(rec.Body.String())

	expectTrue(t, gotRaw == expRaw)
	expectTrue(t, rec.Code == 403)
	expectTrue(t, rec.Header().Get("Content-Type") == "application/problem+json; charset=utf-8")
}

func TestWriteJSON_WithStrictButAllEmpty(t *testing.T) {
	data := problemdetail.New("")

	rec := httptest.NewRecorder()
	err := problemdetail.WriteJSON(rec, data, 0)
	expectTrue(t, err != nil)

	expectTrue(t, errors.Is(err, problemdetail.ErrTypeRequired))
	expectTrue(t, errors.Is(err, problemdetail.ErrTitleRequired))
	expectTrue(t, errors.Is(err, problemdetail.ErrStatusRequired))
	expectTrue(t, errors.Is(err, problemdetail.ErrDetailRequired))
	expectTrue(t, errors.Is(err, problemdetail.ErrInstanceRequired))
	expectTrue(t, !errors.Is(err, problemdetail.ErrTypeFormat))
	expectTrue(t, !errors.Is(err, problemdetail.ErrInstanceFormat))
}

func TestWriteJSON_WithTypedStrictButTypeAndInstanceInvalidFormat(t *testing.T) {
	data := problemdetail.New("--not-\n/a/valid/uri--",
		problemdetail.WithInstance("\n-not/a/valid/path\n"),
		problemdetail.WithInstance("\n-not/a/valid/path\n"),
	)
	rec := httptest.NewRecorder()
	err := problemdetail.WriteJSON(rec, data, 0)
	expectTrue(t, err != nil)
	expectTrue(t, errors.Is(err, problemdetail.ErrTypeFormat))
	expectTrue(t, errors.Is(err, problemdetail.ErrInstanceFormat))
}

func TestWriteXML_WithExtension(t *testing.T) {
	data := BalanceProblemDetail{
		ProblemDetail: problemdetail.New(
			"https://example.com/probs/out-of-credit",
			problemdetail.WithDetail("Your current balance is 30, but that costs 50."),
			problemdetail.WithInstance("/account/12345/abc"),
			problemdetail.WithTitle("You do not have enough credit."),
		),
		Balance:  30,
		Accounts: []string{"/account/12345", "/account/67890"},
	}

	rec := httptest.NewRecorder()
	err := problemdetail.WriteXML(rec, &data, 403)
	expectTrue(t, err == nil)

	rawExp := `<problem xmlns="urn:ietf:rfc:7807"><type>https://example.com/probs/out-of-credit</type><title>You do not have enough credit.</title><status>403</status><detail>Your current balance is 30, but that costs 50.</detail><instance>/account/12345/abc</instance><balance>30</balance><accounts>/account/12345</accounts><accounts>/account/67890</accounts></problem>`
	rawGot := strings.TrimSpace(rec.Body.String())

	expectTrue(t, rawGot == rawExp)
	expectTrue(t, rec.Code == 403)
	expectTrue(t, rec.Header().Get("Content-Type") == "application/problem+xml; charset=utf-8")
}

func TestWriteXML_WithUntyped(t *testing.T) {
	data := problemdetail.New(problemdetail.Untyped, problemdetail.WithValidateLevel(problemdetail.LStandard))

	rec := httptest.NewRecorder()
	err := problemdetail.WriteXML(rec, data, 403)
	expectTrue(t, err == nil)

	rawExp := `<problem xmlns="urn:ietf:rfc:7807"><type>about:blank</type><title>Forbidden</title><status>403</status></problem>`
	rawGot := strings.TrimSpace(rec.Body.String())

	expectTrue(t, rawGot == rawExp)
	expectTrue(t, rec.Code == 403)
	expectTrue(t, rec.Header().Get("Content-Type") == "application/problem+xml; charset=utf-8")
}

func TestWriteXML_WithTyped(t *testing.T) {
	data := problemdetail.New(
		"https://example.com/probs/out-of-credit",
		problemdetail.WithDetail("Your current balance is 30, but that costs 50."),
		problemdetail.WithInstance("/account/12345/abc"),
		problemdetail.WithTitle("You do not have enough credit."),
	)

	rec := httptest.NewRecorder()
	err := problemdetail.WriteXML(rec, data, 403)
	expectTrue(t, err == nil)

	rawExp := `<problem xmlns="urn:ietf:rfc:7807"><type>https://example.com/probs/out-of-credit</type><title>You do not have enough credit.</title><status>403</status><detail>Your current balance is 30, but that costs 50.</detail><instance>/account/12345/abc</instance></problem>`
	rawGot := strings.TrimSpace(rec.Body.String())

	expectTrue(t, rawGot == rawExp)
	expectTrue(t, rec.Code == 403)
	expectTrue(t, rec.Header().Get("Content-Type") == "application/problem+xml; charset=utf-8")
}

func TestWriteXML_WithStrictButAllEmpty(t *testing.T) {
	data := problemdetail.New("")

	rec := httptest.NewRecorder()
	err := problemdetail.WriteXML(rec, data, 0)
	expectTrue(t, err != nil)

	expectTrue(t, errors.Is(err, problemdetail.ErrTypeRequired))
	expectTrue(t, errors.Is(err, problemdetail.ErrTitleRequired))
	expectTrue(t, errors.Is(err, problemdetail.ErrStatusRequired))
	expectTrue(t, errors.Is(err, problemdetail.ErrDetailRequired))
	expectTrue(t, errors.Is(err, problemdetail.ErrInstanceRequired))
	expectTrue(t, !errors.Is(err, problemdetail.ErrTypeFormat))
	expectTrue(t, !errors.Is(err, problemdetail.ErrInstanceFormat))
}

func TestWriteXML_WithTypedStrictButTypeAndInstanceInvalidFormat(t *testing.T) {
	data := problemdetail.New("--not-\n/a/valid/uri--",
		problemdetail.WithInstance("\n-not/a/valid/path\n"),
		problemdetail.WithInstance("\n-not/a/valid/path\n"),
	)
	rec := httptest.NewRecorder()
	err := problemdetail.WriteXML(rec, data, 0)
	expectTrue(t, err != nil)
	expectTrue(t, errors.Is(err, problemdetail.ErrTypeFormat))
	expectTrue(t, errors.Is(err, problemdetail.ErrInstanceFormat))
}

func TestProblemDetail_Error(t *testing.T) {
	pd := problemdetail.New("https://example.com/probs/out-of-credit",
		problemdetail.WithDetail("Your current balance is 30, but that costs 50."),
		problemdetail.WithInstance("/account/12345/abc"),
		problemdetail.WithTitle("You do not have enough credit."),
	)

	err := fmt.Errorf("error: %w", pd)

	var pdErr *problemdetail.ProblemDetail
	expectTrue(t, errors.As(err, &pdErr))
	expectTrue(t, pdErr == pd)
	expectTrue(t, pdErr.Error() == "problem detail: https://example.com/probs/out-of-credit")
}

func expectTrue(t *testing.T, b bool) {
	t.Helper()
	if !b {
		t.Fatal("expected true, got false")
	}
}
