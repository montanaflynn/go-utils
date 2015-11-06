package utils

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

func TestNewPort(t *testing.T) {
	port, err := NewPort()
	if err != nil {
		t.Fatalf("Got err from NewPort(): %s", err)
	}

	if port == 0 {
		t.Fatal("Got port 0")
	}
}

func TestNewPortUnavailable(t *testing.T) {
	port, err := NewPort()
	if err != nil {
		t.Fatalf("Got err from NewPort(): %s", err)
	}

	listener, err := net.Listen("tcp", "127.0.0.1:"+strconv.Itoa(port))
	if err != nil {
		t.Fatalf("Got err from net.Listen: %s", err)
	}
	listener.Close()
}

func TestCheckPort(t *testing.T) {
	status, err := CheckPort(44444)
	if err != nil {
		t.Fatalf("Got err from CheckPort: %s", err)
	}

	if status == false {
		t.Fatal("Got false status for port 44444")
	}
}

func TestCheckPortUnavailable(t *testing.T) {
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

	status, err := CheckPort(portInt)
	if err == nil {
		t.Fatalf("Didn't get err from CheckPort: %s", err)
	}

	if status {
		t.Fatal("Got true status from CheckPort when port is in use")
	}
}

func TestPortToString(t *testing.T) {
	portString := PortToString(4444)
	if portString != ":4444" {
		t.Fatalf("Didn't get correct port string: %s != %s", portString, ":4444")
	}
}
