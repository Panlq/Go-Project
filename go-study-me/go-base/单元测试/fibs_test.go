package split

import (
	"testing"
)

func TestFib(t *testing.T) {
	var (
		in = 7
		expected = 13
	)
	actual := Fib(in)
	if actual != expected {
		t.Errorf("Fib(%d) = %d; expectd %d", in, actual, expected)
	}
}

// func TestFibMulti(t *testing.T) {
// 	var fibTests = []struct {
// 		in int   //input
// 		expected int // expected result
// 	}{
// 		{1, 1},
//         {2, 1},
//         {3, 2},
//         {4, 3},
//         {5, 5},
//         {6, 8},
//         {7, 13},
// 	}

// 	for _, tt := range fibTests {
// 		actual := Fib(tt.in)
// 		if actual != tt.expected {
// 			t.Errorf("Fib(%d) = %d; expectd %d", tt.in, actual, tt.expected)
// 		}
// 	}
// }


// 基准测试函数
func benchmarkFib(i int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		Fib(i)
	}
}

func BenchmarkFib1(b *testing.B)  { benchmarkFib(1, b) }
func BenchmarkFib2(b *testing.B)  { benchmarkFib(2, b) }
func BenchmarkFib3(b *testing.B)  { benchmarkFib(3, b) }
func BenchmarkFib10(b *testing.B) { benchmarkFib(10, b) }
func BenchmarkFib20(b *testing.B) { benchmarkFib(20, b) }
func BenchmarkFib40(b *testing.B) { benchmarkFib(40, b) }
