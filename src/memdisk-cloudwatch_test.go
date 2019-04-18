package main

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

func TestToMB(t *testing.T) {
	data := []struct {
		desc string
		got  uint64
		want float64
	}{
		{"Convert and round to 0MB when result not gt 1MB", 123456, 0},
		{"Convert and round to 1MB when result not gt 2MB", 1234567, 1},
		{"Convert 1024000000 bytes to MB", 1024000000, 976},
	}

	for _, d := range data {
		t.Run(d.desc, func(t *testing.T) {
			got := toMB(d.got)
			if got != d.want {
				t.Fatalf("val want:%f got:%f", d.want, got)
			}
		})
	}
}

func TestToGB(t *testing.T) {
	data := []struct {
		desc string
		got  uint64
		want float64
	}{
		{"Convert and round to 1GB", 1073741820, 1},
		{"Convert and round to 1GB when result not gt 2GB", 1073741825, 1},
		{"Convert 2147483648 bytes", 2147483648, 2},
	}

	for _, d := range data {
		t.Run(d.desc, func(t *testing.T) {
			got := toGB(d.got)
			if got != d.want {
				t.Fatalf("val want:%f got:%f", d.want, got)
			}
		})
	}
}

func TestFsDeviceForMountPath(t *testing.T) {
	t.Run("fsDeviceForMountPath", func(t *testing.T) {
		path := fsDeviceForMountPath("/")
		if path == "" {
			t.Fatalf("val path is empty")
		}
	})
}

func TestReportDisk(t *testing.T) {
	t.Run("reportDisk", func(t *testing.T) {
		icw := instanceCloudwatch{}
		reportDisk(&icw)

		if len(icw.metrics) != 1 {
			t.Fatalf("len(icw.metrics) want:1 got:%v", len(icw.metrics))
		}
		metrics := make([]cloudwatch.MetricDatum, 0)
		for _, m := range icw.metrics {
			if *m.MetricName == "DiskSpaceUtilization" {
				metrics = append(metrics, *m)
			}
		}
		if len(metrics) != 1 {
			t.Fatalf("len(icw.metrics) want:1 got:%v", len(metrics))
		}
	})
}

func TestReportMemory(t *testing.T) {
	t.Run("reportMemory", func(t *testing.T) {
		icw := instanceCloudwatch{}
		reportMemory(&icw)

		if len(icw.metrics) != 1 {
			t.Fatalf("len(icw.metrics) want:1 got:%v", len(icw.metrics))
		}
		metrics := make([]cloudwatch.MetricDatum, 0)
		for _, m := range icw.metrics {
			if *m.MetricName == "MemoryUtilization" {
				metrics = append(metrics, *m)
			}
		}
		if len(metrics) != 1 {
			t.Fatalf("len(icw.metrics) want:1 got:%v", len(metrics))
		}
	})
}
