package main

import (
	"fmt"
	"runtime"
)

func main() {
	cpuCount := runtime.NumCPU()
	fmt.Println("Количество логических ядер:", cpuCount)
}
