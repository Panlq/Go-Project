package unit_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

var queryArgs []string

func init() {
	for i := 0; i < 5000; i++ {
		queryArgs = append(queryArgs, str.NewSeqID())
	}
}

func ConcatenationOperator(strSLice []string) string {
	var q string
	for _, v := range strSLice {
		q = q + "," + v
	}
	return q
}

func BenchmarkConcatenationOperator(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ConcatenationOperator(queryArgs)
	}
	b.ReportAllocs()
}

func FmtSprint(strSLice []string) string {
	var q string
	for _, v := range strSLice {
		q = fmt.Sprint(q, ",", v)
	}
	return q
}

func BenchmarkFmtSprint(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FmtSprint(queryArgs)
	}
	b.ReportAllocs()
}

func StringsJoin(strSLice []string) string {
	return strings.Join(strSLice, ",")
}

func BenchmarkStringsJoin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		StringsJoin(queryArgs)
	}
	b.ReportAllocs()
}

func BytesBuffer(strSlice []string) string {
	var q bytes.Buffer

	q.Grow(36*len(strSlice) + len(strSlice) - 1)

	for i, v := range strSlice {
		q.WriteString(v)
		if i == len(strSlice)-1 {
			continue
		}
		q.WriteString(",")
	}

	return q.String()
}

func BenchmarkBytesBuffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		BytesBuffer(queryArgs)
	}
	b.ReportAllocs()
}

func StringBuilder(strSlice []string) string {
	var q strings.Builder
	// 当最终结果的最大大小已知时，用于预分配内存。确保这样做才能在这种情况下获得最大效率
	q.Grow(36*len(strSlice) + len(strSlice) - 1)

	for i, v := range strSlice {
		q.WriteString(v)
		if i == len(strSlice)-1 {
			continue
		}
		q.WriteString(",")
	}

	return q.String()
}

func BenchmarkStringBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		StringBuilder(queryArgs)
	}
	b.ReportAllocs()
}
