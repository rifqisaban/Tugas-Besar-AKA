package main

import (
	"fmt"
	"math"
	"time"
)

const (
	numRuns    = 5      // Jumlah pengujian untuk perbandingan
	warmUpRuns = 1000   // Jumlah iterasi pemanasan (warm-up)
	iterations = 100000 // Jumlah iterasi untuk pengukuran waktu
	epsilon    = 1e-10  // Konstanta untuk perbandingan floating point
)

// GeometricCalculator holds the parameters for a geometric sequence
type GeometricCalculator struct {
	a float64 // Suku pertama
	r float64 // Rasio
	n int     // Jumlah suku
}

// measureExecutionTime measures the execution time of a function in nanoseconds
func measureExecutionTime(f func()) float64 {
	// Warm-up phase to stabilize any jitter
	for i := 0; i < warmUpRuns; i++ {
		f()
	}

	// Measure execution time
	var totalDuration time.Duration
	for run := 0; run < iterations; run++ {
		start := time.Now()
		f()
		totalDuration += time.Since(start)
	}

	// Return average duration in nanoseconds
	return float64(totalDuration.Nanoseconds()) / float64(iterations)
}

// GeometricSumIterative calculates the sum of a geometric sequence using iteration
func (g *GeometricCalculator) GeometricSumIterative() float64 {
	sum := 0.0
	term := g.a
	for i := 0; i < g.n; i++ {
		sum += term
		term *= g.r
	}
	return sum
}

// GeometricSumRecursive calculates the sum of a geometric sequence using recursion
func (g *GeometricCalculator) GeometricSumRecursive() float64 {
	memo := make(map[int]float64)

	var recursive func(a, r float64, n int) float64
	recursive = func(a, r float64, n int) float64 {
		if n == 0 {
			return 0
		}
		if val, found := memo[n]; found {
			return val
		}
		memo[n] = a + recursive(a*r, r, n-1)
		return memo[n]
	}

	return recursive(g.a, g.r, g.n)
}

// GeometricSumFormula calculates the sum of a geometric sequence using the closed-form formula
func (g *GeometricCalculator) GeometricSumFormula() float64 {
	if math.Abs(g.r-1.0) < epsilon {
		return g.a * float64(g.n)
	}
	return g.a * (1 - math.Pow(g.r, float64(g.n))) / (1 - g.r)
}

// validateInput prompts the user to input valid parameters for the geometric sequence
func validateInput() (float64, float64, int, error) {
	var a, r float64
	var n int

	fmt.Print("Suku pertama (a): ")
	if _, err := fmt.Scan(&a); err != nil || a <= 0 {
		return 0, 0, 0, fmt.Errorf("harap masukkan nilai a > 0")
	}

	fmt.Print("Rasio (r): ")
	if _, err := fmt.Scan(&r); err != nil || r <= 0 {
		return 0, 0, 0, fmt.Errorf("harap masukkan nilai r > 0")
	}

	fmt.Print("Jumlah suku (n): ")
	if _, err := fmt.Scan(&n); err != nil || n <= 0 {
		return 0, 0, 0, fmt.Errorf("harap masukkan nilai n > 0")
	}

	return a, r, n, nil
}

// ComparisonProgram runs the comparison between iterative and recursive methods
func ComparisonProgram() {
	fmt.Println("\n=== Perbandingan Metode ===")
	a, r, n, err := validateInput()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	calc := &GeometricCalculator{a: a, r: r, n: n}

	// Measure iterative time
	iterativeTimes := make([]float64, numRuns)
	var resultIterative float64
	for i := 0; i < numRuns; i++ {
		iterativeTimes[i] = measureExecutionTime(func() {
			resultIterative = calc.GeometricSumIterative()
		})
	}

	// Measure recursive time
	recursiveTimes := make([]float64, numRuns)
	var resultRecursive float64
	for i := 0; i < numRuns; i++ {
		recursiveTimes[i] = measureExecutionTime(func() {
			resultRecursive = calc.GeometricSumRecursive()
		})
	}

	// Calculate average times
	avgIterativeTime := 0.0
	avgRecursiveTime := 0.0
	for i := 0; i < numRuns; i++ {
		avgIterativeTime += iterativeTimes[i]
		avgRecursiveTime += recursiveTimes[i]
	}
	avgIterativeTime /= float64(numRuns)
	avgRecursiveTime /= float64(numRuns)

	// Formula result
	resultFormula := calc.GeometricSumFormula()

	// Output results
	fmt.Println("\n=== Hasil Perbandingan ===")
	fmt.Printf("Iteratif: %.3f (waktu: %.3f ns)\n", resultIterative, avgIterativeTime)
	fmt.Printf("Rekursif: %.3f (waktu: %.3f ns)\n", resultRecursive, avgRecursiveTime)
	fmt.Printf("Hasil: %.2f\n", resultFormula)

	// Performance ratio
	if avgIterativeTime > 0 {
		ratio := avgRecursiveTime / avgIterativeTime
		fmt.Printf("\nPerbandingan waktu (Rekursif/Iteratif): %.2fx\n", ratio)
		if ratio > 1 {
			fmt.Printf("Metode iteratif lebih cepat sebesar %.2f%%\n", (ratio-1)*100)
		} else {
			fmt.Printf("Metode rekursif lebih cepat sebesar %.2f%%\n", (1-ratio)*100)
		}
	}
}

// main is the entry point of the program
func main() {
	for {
		fmt.Println("========================================================")
		fmt.Println("   PERBANDINGAN ALGORITMA ITERATIF DAN REKURSIF")
		fmt.Println("           DALAM MENGHITUNG DERET GEOMETRI")
		fmt.Println("========================================================")
		fmt.Println("\nPilih mode program:")
		fmt.Println("1. Perbandingan Metode iteratif dan rekursif")
		fmt.Println("2. Keluar")
		fmt.Print("\nMasukkan pilihan Anda (1/2): ")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			ComparisonProgram()
		case 2:
			fmt.Println("Terima kasih telah menggunakan program ini. Sampai jumpa!")
			return
		default:
			fmt.Println("Pilihan tidak valid! Harap pilih 1 atau 2.")
		}

		fmt.Println("\nTekan Enter untuk kembali ke menu utama...")
		fmt.Scanln() // Tunggu pengguna untuk melanjutkan
	}
}
