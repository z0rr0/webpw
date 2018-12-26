package main

import (
	"context"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/z0rr0/gopwgen/pwgen"
)

const (
	indexTemplate = "index.html"

	// Length is default password length.
	Length = 8
	// Number is default number of generated passwords.
	Number = 5
	// MaxLength is limit for password length.
	MaxLength = 100
	// MaxNumbers is limit of generated numbers.
	MaxNumbers = 100

	typeANS = "0"
	typeAN  = "1"
	typeA   = "2"
)

var (
	// internal loggers
	loggerError = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
	loggerInfo  = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime)
	// Types is available types of passwords
	Types = map[string]string{
		typeANS: "alphabet symbols + numbers + special symbols",
		typeAN:  "alphabet symbols + numbers",
		typeA:   "alphabet symbols",
	}
)

// Form is struct of HTML form parameters.
type Form struct {
	Length []int
	Number []int
	Types  [][2]string
}

// Result is struct of user's data.
type Result struct {
	Length       int
	Number       int
	Type         string
	NoCapitalize bool
	NoVowels     bool
	NoAmbiguous  bool
	Passwords    chan string
}

// Data is data for html template execution.
type Data struct {
	F *Form
	R *Result
}

func main() {
	host := flag.String("host", "0.0.0.0", "host")
	port := flag.Uint("port", 30080, "port")
	timeout := flag.Uint64("timeout", 30, "handling timeout, seconds")
	index := flag.String("index", indexTemplate, "HTML template file")
	flag.Parse()

	tpl, err := template.ParseFiles(*index)
	if err != nil {
		loggerError.Panic(err)
	}
	portStr := strconv.FormatUint(uint64(*port), 10)
	address := net.JoinHostPort(*host, portStr)
	loggerInfo.Printf("Listen %v\n", address)

	tmt := time.Duration(*timeout) * time.Second
	srv := &http.Server{
		Addr:           address,
		Handler:        http.DefaultServeMux,
		ReadTimeout:    tmt,
		WriteTimeout:   tmt,
		MaxHeaderBytes: 1 << 20, // 1MB
	}
	f := creteForm()
	// there is only one handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, tpl, f)
	})
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, os.Signal(syscall.SIGTERM))
		<-sigint

		if err := srv.Shutdown(context.Background()); err != nil {
			loggerError.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		loggerError.Printf("HTTP server ListenAndServe: %v", err)
	}
	<-idleConnsClosed
	loggerInfo.Printf("stopped")
}

// validate checks request data and initializes passwords generation.
func validate(r *http.Request) (*Result, error) {
	// default values
	length := Length
	number := Number
	t := typeANS
	symbols := true
	noNumerals := false
	noCapitalize := false
	noVowels := true
	noAmbiguous := true

	// password length
	value := r.FormValue("length")
	if value != "" {
		l, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("failed length %v: %v", value, err)
		}
		if (l < 1) || (l > MaxLength) {
			return nil, fmt.Errorf("invalid length value %v", l)
		}
		length = l
	}
	// number of passwords
	value = r.FormValue("number")
	if value != "" {
		n, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("failed number %v: %v", value, err)
		}
		if (n < 1) || (n > MaxNumbers) {
			return nil, fmt.Errorf("invalid number value %v", n)
		}
		number = n
	}
	// type specific
	value = strings.Trim(r.FormValue("type"), " ")
	if value != "" {
		_, ok := Types[value]
		if ok {
			t = value
		}
	}
	switch t {
	case typeAN: // alphabet symbols + numbers
		symbols = false
	case typeA: // alphabet symbols
		noNumerals = true
		symbols = false
	}
	// capitalize
	value = r.FormValue("no_capitalize")
	if value == "on" {
		noCapitalize = true
	}
	// vowels
	value = r.FormValue("vowels")
	if value == "on" {
		noVowels = false
	}
	// ambiguous
	value = r.FormValue("ambiguous")
	if value == "on" {
		noAmbiguous = false
	}
	// passwords generation
	pw, err := pwgen.New(
		length, number, "", "",
		noNumerals, true, false,
		noCapitalize, noAmbiguous, symbols, noVowels, true,
	)
	if err != nil {
		return nil, err
	}
	result := &Result{
		Length:       length,
		Number:       number,
		Type:         t,
		NoCapitalize: noCapitalize,
		NoVowels:     noVowels,
		NoAmbiguous:  noAmbiguous,
		Passwords:    pw.Passwords(),
	}
	return result, nil
}

// creteForm creates html form data.
func creteForm() *Form {
	formLength := make([]int, 95)
	for i := 0; i < 95; i++ {
		formLength[i] = i + 6
	}
	types := make([][2]string, 3)
	for i := range types {
		types[i] = [2]string{strconv.Itoa(i), Types[strconv.Itoa(i)]}
	}
	return &Form{
		Length: formLength,
		Number: []int{5, 10, 15},
		Types:  types,
	}
}

// handler handles incoming requests.
func handler(w http.ResponseWriter, r *http.Request, tpl *template.Template, f *Form) {
	start, code := time.Now(), http.StatusOK
	defer func() {
		loggerInfo.Printf("%-5v %v\t%-12v\t%v",
			r.Method,
			code,
			time.Since(start),
			r.URL.String(),
		)
	}()
	if r.URL.Path != "/" {
		code = http.StatusNotFound
		http.NotFound(w, r)
		return
	}
	result, err := validate(r)
	if err != nil {
		loggerError.Println(err)
		code = http.StatusBadRequest
		http.Error(w, "Bad request", code)
		return
	}
	data := &Data{
		F: f,
		R: result,
	}
	err = tpl.Execute(w, data)
	if err != nil {
		loggerError.Println(err)
		code = http.StatusInternalServerError
		http.Error(w, "Internal Server Error", code)
	}
}
