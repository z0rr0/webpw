package main

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type TestCase struct {
	URL    string
	Params string
	Code   int
}

func Test(t *testing.T) {
	tpl, err := template.ParseFiles(indexTemplate)
	if err != nil {
		t.Fatal(err)
	}
	f := prepareForm()
	testCases := []TestCase{
		{"/", "", http.StatusOK},
		{"/bad", "", http.StatusNotFound},
		{"/", "length=12&number=2&type=0", http.StatusOK},
		{"/", "length=12&number=2&type=bad", http.StatusOK},
		{"/", "length=12&number=2&type=1", http.StatusOK},
		{"/", "length=12&number=2&type=2", http.StatusOK},
		{"/", "length=12&number=2&type=0&no_capitalize=on", http.StatusOK},
		{"/", "length=12&number=2&type=0&vowels=on", http.StatusOK},
		{"/", "length=12&number=2&type=0&ambiguous=on", http.StatusOK},
		{"/", "length=12&number=2&type=0&ambiguous=on&o_capitalize=on&ambiguous=on", http.StatusOK},
		{"/", "length=0&number=2&type=0", http.StatusOK},
		{"/", "length=0&number=0&type=0", http.StatusOK},
	}
	for i, tc := range testCases {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", tc.URL, strings.NewReader(tc.Params))
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		handler(w, r, tpl, f)
		if w.Code != tc.Code {
			t.Errorf("failed code [%v]: %v != %v", i, w.Code, tc.Code)
		}
		if tc.Code == http.StatusOK {
			result := w.Body.String()
			strings.Contains(result, "Generated passwords")
		}
	}
}
