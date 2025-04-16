package api

import (
	"log"
	"os"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/process"
	"github.com/teamleaderleo/potato-quality-image-compressor/internal/metrics"
)

// StartResourceMonitor starts monitoring system resources like memory and CPU
func StartResourceMonitor(interval time.Duration) {
	go monitorResourceUsage(interval)
}

// monitorResourceUsage periodically updates memory and CPU usage metrics
func monitorResourceUsage(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Get process info for this application
	proc, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		log.Printf("Error setting up process monitoring: %v", err)
		// Fall back to basic monitoring if process monitoring fails
		go monitorBasicResourceUsage(interval)
		return
	}

	for range ticker.C {
		// Memory metrics (runtime)
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		metrics.UpdateMemoryUsage(m.Alloc)

		// Memory metrics (system)
		if vmem, err := mem.VirtualMemory(); err == nil {
			metrics.UpdateSystemMemoryUsage(vmem.Used)
			metrics.UpdateSystemMemoryPercent(vmem.UsedPercent)
		}

		// CPU metrics (process)
		if cpuPercent, err := proc.CPUPercent(); err == nil {
			metrics.UpdateCPUUsage(cpuPercent)
		}

		// CPU metrics (system)
		if cpuPercents, err := cpu.Percent(0, false); err == nil && len(cpuPercents) > 0 {
			metrics.UpdateSystemCPUUsage(cpuPercents[0])
		}
	}
}

// monitorBasicResourceUsage is a fallback that uses only runtime metrics
func monitorBasicResourceUsage(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	var m runtime.MemStats
	for range ticker.C {
		runtime.ReadMemStats(&m)
		metrics.UpdateMemoryUsage(m.Alloc)
		
		log.Printf("Using basic monitoring - CPU metrics unavailable")
	}
}