package handlers

import (
	"slices"
	"testing"
)

func TestStrip(t *testing.T) {
	{
		data := []int{5, 6, 8, 0, 0, 0, -2, 0, 3, 4, 5, 6}
		expect := []int{5, 6, 8, -2, 3, 4, 5, 6}

		strip(&data)
		if !slices.Equal(data, expect) {
			t.Errorf("expected slice to be '%v', got '%v'", expect, data)
		}
	}

	{
		data := []int{0, 0, 0, 0, 0}
		expect := []int{}

		strip(&data)
		if !slices.Equal(data, expect) {
			t.Errorf("expected slice to be '%v', got '%v'", expect, data)
		}
	}

	{
		data := []int{0, 0, 0, 0, 1}
		expect := []int{1}

		strip(&data)
		if !slices.Equal(data, expect) {
			t.Errorf("expected slice to be '%v', got '%v'", expect, data)
		}
	}

	{
		data := []int{1, 0, 0, 0, 0}
		expect := []int{1}

		strip(&data)
		if !slices.Equal(data, expect) {
			t.Errorf("expected slice to be '%v', got '%v'", expect, data)
		}
	}

	{
		data := []int{1, 2, 3, 4, 5}
		expect := []int{1, 2, 3, 4, 5}

		strip(&data)
		if !slices.Equal(data, expect) {
			t.Errorf("expected slice to be '%v', got '%v'", expect, data)
		}
	}

	{
		type s struct {
			msg string
		}

		data := []s{{}, {}, {msg: "hello"}, {}, {msg: "world"}, {}, {}, {}, {}}
		expect := []s{{msg: "hello"}, {msg: "world"}}

		strip(&data)
		if !slices.Equal(data, expect) {
			t.Errorf("expected slice to be '%v', got '%v'", expect, data)
		}
	}
}

func BenchmarkStrip(b *testing.B) {
	for i := 0; i < b.N; i++ {
		data := []int{1, 2, 0, 3, 0, 0, 0, 4, 5, 6, 0, 7, 0, 8, 0, 9, 0}
		strip(&data)
	}
}
