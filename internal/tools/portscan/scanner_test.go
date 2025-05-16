package portscan

import (
	"testing"
	"time"
)

func TestNewScanner(t *testing.T) {
	target := "localhost"
	startPort := 1
	endPort := 1024

	scanner := NewScanner(target, startPort, endPort)

	if scanner.target != target {
		t.Errorf("Expected target %s, got %s", target, scanner.target)
	}

	if scanner.startPort != startPort {
		t.Errorf("Expected start port %d, got %d", startPort, scanner.startPort)
	}

	if scanner.endPort != endPort {
		t.Errorf("Expected end port %d, got %d", endPort, scanner.endPort)
	}
}

func TestSetTimeout(t *testing.T) {
	scanner := NewScanner("localhost", 1, 100)
	timeout := time.Second * 5

	scanner.SetTimeout(timeout)

	if scanner.timeout != timeout {
		t.Errorf("Expected timeout %v, got %v", timeout, scanner.timeout)
	}
}

func TestSetConcurrent(t *testing.T) {
	scanner := NewScanner("localhost", 1, 100)
	concurrent := 50

	scanner.SetConcurrent(concurrent)

	if scanner.concurrent != concurrent {
		t.Errorf("Expected concurrent %d, got %d", concurrent, scanner.concurrent)
	}
}

func TestScanPort(t *testing.T) {
	scanner := NewScanner("localhost", 1, 100)
	scanner.SetTimeout(time.Millisecond * 100) // Short timeout for tests

	// Test an invalid port (should return nil)
	result := scanner.scanPort(0)
	if result != nil {
		t.Errorf("Expected nil result for invalid port, got %+v", result)
	}

	// Test a valid port
	result = scanner.scanPort(80)
	// Note: We can't reliably test for open ports in unit tests
	// as they depend on the system state
	if result != nil {
		if result.Port != 80 {
			t.Errorf("Expected port 80, got %d", result.Port)
		}
		if result.State != "open" {
			t.Errorf("Expected state 'open', got %s", result.State)
		}
		if result.Service == "" {
			t.Error("Expected non-empty service name")
		}
	}
}
