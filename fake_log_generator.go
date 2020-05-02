package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

// ApacheCommonLog : {remotehost} {rfc931} {authuser} [{datetime}] "{method} {request} {protocole}" {status} {bytes}
const (
	ApacheCommonLog      = "%s %s %s [%s] \"%s %s %s\" %d %d"
	ApacheCommonDatetime = "02/Jan/2006:15:04:05 -0700"
)

func IPv4Address() string {
	num := func() int { return 2 + rand.Intn(254) }
	return fmt.Sprintf("%d.%d.%d.%d", num(), num(), num(), num())
}

func RandAuthUserID() string {
	candidates := []string{"claire", "david", "frank", "james", "mary", "johnny", "louise", "tom"}
	return candidates[rand.Intn(8)]
}

func RandHttpMethod() string {
	candidates := []string{"GET", "POST", "PUT", "DELETE"}
	return candidates[rand.Intn(4)]
}

func RandNumber(min int, max int) int {
	return rand.Intn((max+1)-min) + min
}

// Return Most Common Status Code
func RandStatusCode() int {
	candidates := []int{200, 301, 302, 401, 403, 404, 500, 503, 504}
	return candidates[rand.Intn(9)]
}

func RandHTTPVersion() string {
	versions := []string{"HTTP/1.0", "HTTP/1.1", "HTTP/2.0"}
	return versions[rand.Intn(3)]
}

func RandSection() string {
	candidates := []string{"/api", "/api/user", "/v2/api/groups", "/sections",
		"/api/search", "/sections/users", "/api/user/place",
		"/sections/search", "/test", "/test/users",
		"/home", "/admin/users", "/admin", "/bag", "/bag/buy",
		"/error", "/error/section"}
	return candidates[rand.Intn(8)]
}

func NewApacheLog(t time.Time) string {
	return fmt.Sprintf(
		ApacheCommonLog,
		IPv4Address(),
		"-",
		RandAuthUserID(),
		t.Format(ApacheCommonDatetime),
		RandHttpMethod(),
		RandSection(),
		RandHTTPVersion(),
		RandStatusCode(),
		RandNumber(0, 5000),
	)
}

type Option struct {
	Delay  float64
	Output string
}

func OpenOrCreate(path string) *os.File {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		fmt.Printf("Error: Can't create config file, %s\n", err)
		os.Exit(1)
	}
	return f
}

func defaultOptions() *Option {
	return &Option{
		Output: "/tmp/access.log",
		Delay:  1000,
	}
}

func ParseArgs() *Option {
	opt := defaultOptions()
	output := flag.String("output", opt.Output, "Output filename. Path-like filename is allowed")
	delay := flag.Float64("delay", opt.Delay, "Delay log generation speed (in seconds)")
	flag.Parse()
	opt.Output = *output
	opt.Delay = *delay
	return opt
}

func Run(f *os.File, opt *Option) {
	delay := time.Duration(opt.Delay)
	for {
		time.Sleep(delay * time.Millisecond)
		log := NewApacheLog(time.Now())
		f.Write([]byte(log + "\n"))
	}
}

func main() {
	opt := ParseArgs()
	f := OpenOrCreate(opt.Output)
	defer f.Close()
	rand.Seed(time.Now().UnixNano())
	Run(f, opt)
}
