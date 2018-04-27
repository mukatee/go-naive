package tests

//https://blog.alexellis.io/golang-writing-unit-tests/

import (
	"testing"
	"time"
	"fmt"
	"strconv"
)

func TestSum(t *testing.T) {
	total := 5+5
	if total != 10 {
		t.Errorf("Sum was incorrect, got: %d, want: %d.", total, 10)
	} else {
		println("pass")
	}
}

func TestTimeDiff(t *testing.T) {
	//NOTE: this has to be exactly 24.03 since the layout has the zero in it..
	start, _ := time.Parse("02.01.2006 15:04:05", "24.03.2018 19:00:00")
	end, _ := time.Parse("02.01.2006 15:04:05", "24.03.2018 19:05:00")
	diffTime := end.Sub(start)
	fmt.Printf("start = %v\n", start)
	fmt.Printf("end = %v\n", end)
	fmt.Printf("difference = %v\n", diffTime)
	//https://stackoverflow.com/questions/24122821/go-golang-time-now-unixnano-convert-to-milliseconds
	diff := diffTime.Minutes()
	diffStr := strconv.FormatFloat(diff, 'f', 2, 64)
	fmt.Println(diffStr)
	if diff != 5 {
		t.Errorf("Diff of dates with 5 min difference incorrect, expected %d minutes, got %f", 5, diff)
	}

	start, _ = time.Parse("02.01.2006 15:04:05", "24.03.2018 19:00:00")
	end, _ = time.Parse("02.01.2006 15:04:05", "24.03.2018 21:05:00")
	diffTime = end.Sub(start)
	diff = diffTime.Minutes()
	if diff != 125 {
		t.Errorf("Diff of dates with 2h5m difference incorrect, expected %d minutes, got %f", 125, diff)
	}
}
