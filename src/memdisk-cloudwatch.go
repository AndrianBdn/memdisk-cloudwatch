package main

import (
	"flag"
	"log"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

func main() {
	crontabPtr := flag.Bool("crontab", false, "run from cron, with random 0..20 delay (by default is off, run as foreground service with 5 min timer)")
	runoncePtr := flag.Bool("runonce", false, "like 'crontab' but without delay")
	addcontainer := flag.Bool("addcontainer", false, "add status of containers")
	flag.Parse()
	runonce := *runoncePtr || *crontabPtr
	if *crontabPtr {
		time.Sleep(time.Duration(rand.Intn(20)) * time.Second)
	}

	icw := NewInstanceCloudwatch(fsDeviceForMountPath("/"), "/")
	reportMetricsOnce(icw, *addcontainer)
	if !runonce {
		startFiveMinuteTicker(icw, *addcontainer)
		select {}
	}
}

func reportContainers(icw *instanceCloudwatch) {
	containers, err := DockeInspect()
	if err != nil {
		log.Fatal("Failed to get info of containers:", err)
	}
	for _, container := range containers {
		icw.AddContainerMetric(container.name, container.status)
	}
}

func reportMemory(icw *instanceCloudwatch) {

	v, err := mem.VirtualMemory()
	if err != nil {
		log.Fatal("Failed to get memory usage:", err)
	}

	icw.AddMemoryPercent("MemoryUtilization", v.UsedPercent)

	/// Normally I prefer to avoid commending out code blocks
	/// However in this case this is a step before adding possibly non-necessary configuration

	/// memory metrics (megabytes)
	//icw.AddMemoryMegabytes("MemoryUsed", toMB(v.Used))
	//icw.AddMemoryMegabytes("MemoryAvailable", toMB(v.Available))

	/// swap status (swap is usually turned off in EC2)
	//s, err := mem.SwapMemory()
	//if err != nil {
	//	log.Fatal("Failed to get swap usage:", err)
	//}
	//icw.AddMemoryPercent("SwapUtilization", s.UsedPercent)
}

func reportDisk(icw *instanceCloudwatch) {
	du, err := disk.Usage("/")
	if err != nil {
		log.Printf("Failed to get disk usage: %s", err)
		os.Exit(2)
	}

	icw.AddDiskPercent("DiskSpaceUtilization", du.UsedPercent)

	/// Normally I prefer to avoid commending out code blocks
	/// However in this case this is a step before adding possibly non-necessary configuration

	//icw.AddDiskGigabytes("DiskSpaceUsed", toGB(du.Used))
	//icw.AddDiskGigabytes("DiskSpaceAvailable", toGB(du.Free))
	//icw.AddDiskPercent("DiskInodesUtilization", du.InodesUsedPercent)
}

func reportMetricsOnce(icw *instanceCloudwatch, addcontainer bool) {
	icw.ResetMetrics()
	reportMemory(icw)
	reportDisk(icw)
	if addcontainer {
		reportContainers(icw)
	}
	icw.Send()
}

func startFiveMinuteTicker(icw *instanceCloudwatch, addcontainer bool) {
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				reportMetricsOnce(icw, addcontainer)
			}
		}
	}()
}

func fsDeviceForMountPath(mountPath string) string {

	partitions, err := disk.Partitions(false)
	if err != nil {
		log.Fatal("Can't get partitions list: ", err)
	}

	var rootdevice string
	for _, val := range partitions {
		if val.Mountpoint == mountPath {
			rootdevice = val.Device
		}
	}

	return rootdevice
}

func toMB(bytes uint64) float64 {
	return float64(bytes / 1048576)
	// integer division, only whole megabytes
	// we don't need half of megabyte in 2016
}

func toGB(bytes uint64) float64 {
	gb := toMB(bytes) / 1024.0
	return math.Ceil(gb*100.0) / 100.0
	// only two points after dot
}
