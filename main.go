package main

import (
	"log"
	"time"

	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
)

const (
	THRESHOLD = 50
	INTERVAL  = 5
)

func main() {
	for range time.Tick(time.Second * INTERVAL) {

		memoryStats, err := mem.VirtualMemory()
		if err != nil {
			log.Fatal(err)
		}

		percentage := (float64(memoryStats.Used) / float64(memoryStats.Total)) * 100

		log.Println("[+] Current percentage: ", percentage)

		if percentage > THRESHOLD {

			processList, err := process.Processes()
			if err != nil {
				log.Fatal(err)
			}

			type Process struct {
				Name        string
				MemoryUsage float32
			}

			db := make(map[string]Process)

			for _, currentProcess := range processList {
				name, _ := currentProcess.Name()
				memory_usage, _ := currentProcess.MemoryPercent()

				db[name] = Process{
					Name:        name,
					MemoryUsage: db[name].MemoryUsage + memory_usage,
				}

			}

			var process Process

			for _, value := range db {
				if value.MemoryUsage > process.MemoryUsage {
					process = value
				}
			}

			log.Println("[+] Spike detected: ", process.Name, process.MemoryUsage)

			for _, currentProcess := range processList {
				name, _ := currentProcess.Name()

				if name == process.Name {
					err := currentProcess.Kill()
					if err != nil {
						log.Println(err)
					}
				}
			}

		}
	}
}
