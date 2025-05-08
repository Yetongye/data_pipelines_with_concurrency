// benchmark tests for the main package.

package main

import (
	"testing"
)

var testPaths = []string{
	"images/image1.jpeg",
	"images/image2.jpeg",
	"images/image5.jpeg",
	"images/image6.jpeg",
}

func BenchmarkMainConcurrent(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runConcurrent(testPaths)
	}
}

func BenchmarkMainSequential(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runSequential(testPaths)
	}
}
