package portscan

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// Result represents a port scan result
type Result struct {
	Port    int
	State   string
	Service string
}

// Scanner represents a port scanner
type Scanner struct {
	target     string
	startPort  int
	endPort    int
	timeout    time.Duration
	concurrent int
}

// NewScanner creates a new port scanner
func NewScanner(target string, startPort, endPort int) *Scanner {
	return &Scanner{
		target:     target,
		startPort:  startPort,
		endPort:    endPort,
		timeout:    time.Second * 2,
		concurrent: 100,
	}
}

// SetTimeout sets the connection timeout
func (s *Scanner) SetTimeout(timeout time.Duration) {
	s.timeout = timeout
}

// SetConcurrent sets the number of concurrent scans
func (s *Scanner) SetConcurrent(n int) {
	s.concurrent = n
}

// Scan performs the port scan
func (s *Scanner) Scan(progress chan<- float64) ([]Result, error) {
	var results []Result
	var mutex sync.Mutex
	var wg sync.WaitGroup

	// Create work channel
	ports := make(chan int, s.concurrent)

	// Start workers
	for i := 0; i < s.concurrent; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for port := range ports {
				result := s.scanPort(port)
				if result != nil {
					mutex.Lock()
					results = append(results, *result)
					mutex.Unlock()
				}
				// Report progress
				if progress != nil {
					progress <- float64(port-s.startPort) / float64(s.endPort-s.startPort) * 100
				}
			}
		}()
	}

	// Send work
	for port := s.startPort; port <= s.endPort; port++ {
		ports <- port
	}
	close(ports)

	// Wait for completion
	wg.Wait()

	return results, nil
}

// scanPort scans a single port
func (s *Scanner) scanPort(port int) *Result {
	target := fmt.Sprintf("%s:%d", s.target, port)
	conn, err := net.DialTimeout("tcp", target, s.timeout)

	if err != nil {
		return nil
	}

	defer conn.Close()

	service := getServiceName(port)
	return &Result{
		Port:    port,
		State:   "open",
		Service: service,
	}
}

// getServiceName returns the common service name for a port
func getServiceName(port int) string {
	services := map[int]string{
		20:   "FTP-DATA",
		21:   "FTP",
		22:   "SSH",
		23:   "TELNET",
		25:   "SMTP",
		53:   "DNS",
		80:   "HTTP",
		110:  "POP3",
		143:  "IMAP",
		443:  "HTTPS",
		3306: "MySQL",
		5432: "PostgreSQL",
		8080: "HTTP-ALT",
	}

	if service, ok := services[port]; ok {
		return service
	}
	return "unknown"
}
