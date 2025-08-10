package storage

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

type TestCase struct {
	name  string
	key   string
	value any
}

func TestSetGet(t *testing.T) {
	cases := []TestCase{
		{"1", "first", "world"},
		{"2", "second", 2},
		{"3", "third", "3d"},
	}

	s, err := NewStorage()
	if err != nil {
		t.Errorf("no storage: %v", err)
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s.Set(c.key, c.value)

			sValue := s.Get(c.key)

			if *sValue != c.value {
				t.Errorf("values not equal")
			}
		})
	}
}

func TestKind(t *testing.T) {
	cases := []TestCase{
		{"1", "first", "world"},
		{"2", "second", 2},
		{"3", "third", "3d"},
	}

	s, err := NewStorage()
	if err != nil {
		t.Errorf("no storage: %v", err)
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s.Set(c.key, c.value)
			var kind Kind
			kind = getType(c.value)

			if s.inner[c.key].ValueType != kind {
				t.Errorf("kinds not equal")
			}
		})
	}
}

var casebench = []TestCase{
	{"Hello world", "Hello", "world"},
	{"number1221", "number", 1221},
	{"number666", "number", 666},
}

func BenchmarkGet(b *testing.B) {
	s, _ := NewStorage()

	for _, c := range casebench {
		s.Set(c.key, c.value)
	}

	var results strings.Builder

	for _, c := range casebench {
		b.Run(c.name, func(bb *testing.B) {
			bb.ResetTimer()

			for i := 0; i < bb.N; i++ {
				_ = s.Get(c.key)
			}

			results.WriteString(fmt.Sprintf(
				"%s: %d ops, %v/op\n",
				c.name,
				bb.N,
				bb.Elapsed()/time.Duration(bb.N),
			))
		})
	}
	SaveBenchMarkResults(results.String())

}

func BenchmarkSet(b *testing.B) {

	var results strings.Builder

	for _, c := range casebench {
		b.Run(c.name, func(bb *testing.B) {
			s, _ := NewStorage()
			bb.ResetTimer()

			for i := 0; i < bb.N; i++ {
				s.Set(c.key, c.value)
			}

			results.WriteString(fmt.Sprintf(
				"%s: %d ops, %v/op\n",
				c.name,
				bb.N,
				bb.Elapsed()/time.Duration(bb.N),
			))
		})
	}
	SaveBenchMarkResults(results.String())

}

func BenchmarkSetGet(b *testing.B) {

	var results strings.Builder

	for _, c := range casebench {
		b.Run(c.name, func(bb *testing.B) {
			s, _ := NewStorage()
			bb.ResetTimer()

			for i := 0; i < bb.N; i++ {
				s.Set(c.key, c.value)
				_ = s.Get(c.key)
			}

			results.WriteString(fmt.Sprintf(
				"%s: %d ops, %v/op\n",
				c.name,
				bb.N,
				bb.Elapsed()/time.Duration(bb.N),
			))

		})
	}
	SaveBenchMarkResults(results.String())
}

func SaveBenchMarkResults(output string) {
	if err := os.MkdirAll("benchmark_results", 0755); err != nil {
		panic(err)
	}

	file, err := os.OpenFile(
		"benchmark_results/results.txt",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		panic(err)
	}

	defer file.Close()
	if _, err := file.WriteString(output + "\n"); err != nil {
		panic(err)
	}
}
