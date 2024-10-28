package main

import (
	"fmt"
	"log"
	"time"

	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
)

const (
	THRESHOLD = 70
	INTERVAL  = 5
)

// Format any form of "number"
func fmtPercentage(percentage interface{}) string {
	return fmt.Sprintf("%.2f", percentage)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	for range time.Tick(time.Second * INTERVAL) {
		log.Println("[+] Checking...")

		memoryStats, err := mem.VirtualMemory()
		checkErr(err)

		percentage := (float64(memoryStats.Used) / float64(memoryStats.Total)) * 100

		log.Println("[+] Current percentage: ", fmtPercentage(percentage))

		if percentage > THRESHOLD {

			log.Println("[+] Spike Detected: ", fmtPercentage(percentage))

			processList, err := process.Processes()
			checkErr(err)

			type Process struct {
				Name        string
				MemoryUsage float32
			}

			// Condense all child processes into one
			db := make(map[string]Process)

			for _, currentProcess := range processList {
				name, err := currentProcess.Name()
				checkErr(err)

				memoryUsage, err := currentProcess.MemoryPercent()
				checkErr(err)

				db[name] = Process{
					Name:        name,
					MemoryUsage: db[name].MemoryUsage + memoryUsage,
				}

			}

			var process Process

			for _, value := range db {
				if value.MemoryUsage > process.MemoryUsage {
					process = value
				}
			}

			log.Println("[+] Largest process detected: ", process.Name, fmtPercentage(process.MemoryUsage))

			for _, currentProcess := range processList {
				name, _ := currentProcess.Name()

				if name == process.Name {
					err := currentProcess.Kill()
					checkErr(err)
				}
			}

		}
	}
}
