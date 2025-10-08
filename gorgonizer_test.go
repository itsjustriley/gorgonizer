package main

import (
	"testing"
)

func TestIsExactFolder(t *testing.T) {
    var tests = []struct {
        name     string
        expected bool
    }{
        {"NoExt", true},
        {"PNG", true},
        {"MOV", true},
        {"Images", false},
        {"MY MOV", false},
        {"", false},
    }

    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            if output := isExactFolder(test.name); output != test.expected {
                t.Errorf("input %q, expected %v, got %v", test.name, test.expected, output)
            }
        })
    }
}


func TestHumanizeBytes(t *testing.T) {
    var tests = []struct {
        size     int64
        expected string
    }{
        {1023, "1023 B"},
        {1024, "1.00 KB"},
        {1536, "1.50 KB"},
        {1024 * 1024, "1.00 MB"},
        {1024 * 1024 * 1024, "1.00 GB"},
    }
    for _, test := range tests {
        t.Run(test.expected, func(t *testing.T) {
            if output := humanizeBytes(test.size); output != test.expected {
                t.Errorf("input %d, expected %q, got %q", test.size, test.expected, output)
            }
        })
    }
}

func BenchmarkHumanizeBytes(b *testing.B) {
    for i := 0; i < b.N; i++ {
        humanizeBytes(1 << 40)
    }
}