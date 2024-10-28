package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
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
	threshold := flag.Int("threshold", 70, "Set RAM threshold in percentage")
	interval := flag.Int("interval", 5, "Interval to keep checking in seconds")

	flag.Parse()

	log.Println("[+] Starting...")

	for range time.Tick(time.Second * time.Duration(*interval)) {
		log.Println("[+] Checking...")

		memoryStats, err := mem.VirtualMemory()
		checkErr(err)

		percentage := (float64(memoryStats.Used) / float64(memoryStats.Total)) * 100

		log.Println("[+] Current percentage: ", fmtPercentage(percentage))

		if percentage > float64(*threshold) {

			log.Println("[+] Spike Detected: ", fmtPercentage(percentage))

			beeep.Alert("Ram Guard", fmt.Sprintf("Memory at %s", fmtPercentage(percentage)), "")

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
			beeep.Alert("Ram Guard", fmt.Sprintf("Largest process detected", fmtPercentage(percentage)), "")

			for _, currentProcess := range processList {
				name, _ := currentProcess.Name()

				if name == process.Name {
					err := currentProcess.Kill()
					if err == nil {
						beeep.Alert("Ram Guard", fmt.Sprintf("Killed %s", name), "")
					}
				}
			}

		}
	}
}
