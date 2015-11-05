package port

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

func TestGet(t *testing.T) {
	port, err := Get()
	if err != nil {
		t.Fatalf("Got err from Get: %s", err)
	}

	if port == 0 {
		t.Fatal("Got port 0")
	}
}

func TestGetUnavailablePortCheck(t *testing.T) {
	port, err := Get()
	if err != nil {
		t.Fatalf("Got err from Get: %s", err)
	}

	listener, err := net.Listen("tcp", "127.0.0.1:"+strconv.Itoa(port))
	if err != nil {
		t.Fatalf("Got err from net.Listen: %s", err)
	}
	listener.Close()
}

func TestCheck(t *testing.T) {
	status, err := Check(44444)
	if err != nil {
		t.Fatalf("Got err from Check: %s", err)
	}

	if status == false {
		t.Fatal("Got false status for port 44444")
	}
}

func TestCheckUnavailablePort(t *testing.T) {
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello Tester")
	})

	ts := httptest.NewServer(testHandler)
	defer ts.Close()

	parsedURL, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatalf("Got err from url.Parse: %s", err)
	}

	_, portString, err := net.SplitHostPort(parsedURL.Host)
	if err != nil {
		t.Fatalf("Got err from net.SplitHostPort: %s", err)
	}

	portInt, err := strconv.Atoi(portString)
	if err != nil {
		t.Fatalf("Got err from strconv.Atoi: %s", err)
	}

	status, err := Check(portInt)
	if err == nil {
		t.Fatalf("Didn't get err from Check: %s", err)
	}

	if status {
		t.Fatal("Got true status from Check when port is in use")
	}
}

func TestToString(t *testing.T) {
	portString := ToString(4444)
	if portString != ":4444" {
		t.Fatalf("Didn't get correct port string: %s != %s", portString, ":4444")
	}
}
