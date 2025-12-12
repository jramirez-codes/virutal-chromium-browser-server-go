package util

import (
	"time"

	"virtual-browser/internal/types"

	"github.com/shirou/gopsutil/v3/cpu"
)

// Server Stats
func StartServerStats(serverStats *types.ServerStatsResponse, sleepTime time.Duration) {
	for {
		percentCPU, _ := cpu.Percent(time.Second, false)
		serverStats.Mu.Lock()
		serverStats.CPUUsage = percentCPU[0]
		serverStats.Mu.Unlock()
		time.Sleep(sleepTime)
	}
}
