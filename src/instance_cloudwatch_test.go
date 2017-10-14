package main

import (
	"testing"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

func TestPrefill(t *testing.T) {
	t.Run("prefill", func(t *testing.T) {
		icw := instanceCloudwatch{}
		icw.prefill("blockDevice", "mountPath")

		if len(icw.memDimensions) != 1 {
			t.Fatalf("len(icw.memDimensions) want:1 got:%v", len(icw.memDimensions))
		}
		if len(icw.diskDimentions) != 3 {
			t.Fatalf("len(icw.diskDimentions) want:3 got:%v", len(icw.diskDimentions))
		}
	})
}

func TestResetMetrics(t *testing.T) {
	t.Run("ResetMetrics", func(t *testing.T) {
		icw := instanceCloudwatch{}
		icw.metrics = make([]*cloudwatch.MetricDatum, 0)
		icw.metrics = append(icw.metrics, &cloudwatch.MetricDatum{})

		if len(icw.metrics) != 1 {
			t.Fatalf("len(icw.metrics) want:1 got:%v", len(icw.metrics))
		}

		icw.ResetMetrics()

		if len(icw.metrics) != 0 {
			t.Fatalf("len(icw.metrics) want:0 got:%v", len(icw.metrics))
		}
	})
}

func TestAddDiskMetric(t *testing.T) {
	t.Run("addDiskMetric", func(t *testing.T) {
		icw := instanceCloudwatch{}
		icw.metrics = make([]*cloudwatch.MetricDatum, 0)
		icw.addDiskMetric("Name", "Unit", 10)

		if len(icw.metrics) != 1 {
			t.Fatalf("len(icw.metrics) want:1 got:%v", len(icw.metrics))
		}
	})
}

func TestAddDiskGigabytes(t *testing.T) {
	t.Run("AddDiskGigabytes", func(t *testing.T) {
		icw := instanceCloudwatch{}
		icw.metrics = make([]*cloudwatch.MetricDatum, 0)
		icw.AddDiskGigabytes("Name", 10)

		metrics := make([]cloudwatch.MetricDatum, 0)
		for _, m := range icw.metrics {
			if *m.MetricName == "Name" {
				metrics = append(metrics, *m)
			}
		}

		if len(icw.metrics) != 1 {
			t.Fatalf("len(icw.metrics) want:1 got:%v", len(icw.metrics))
		}
		if *icw.metrics[0].Unit != "Gigabytes" {
			t.Fatalf("icw.metrics[0].Unit want:Gigabytes got:%v", icw.metrics[0].Unit)
		}
	})
}

func TestAddDiskPercent(t *testing.T) {
	t.Run("AddDiskPercent", func(t *testing.T) {
		icw := instanceCloudwatch{}
		icw.metrics = make([]*cloudwatch.MetricDatum, 0)
		icw.AddDiskPercent("Name", 10)

		metrics := make([]cloudwatch.MetricDatum, 0)
		for _, m := range icw.metrics {
			if *m.MetricName == "Name" {
				metrics = append(metrics, *m)
			}
		}

		if len(icw.metrics) != 1 {
			t.Fatalf("len(icw.metrics) want:1 got:%v", len(icw.metrics))
		}
		if *icw.metrics[0].Unit != "Percent" {
			t.Fatalf("icw.metrics[0].Unit want:Percent got:%v", icw.metrics[0].Unit)
		}
	})
}

func TestAddMemoryMetric(t *testing.T) {
	t.Run("addMemoryMetric", func(t *testing.T) {
		icw := instanceCloudwatch{}
		icw.metrics = make([]*cloudwatch.MetricDatum, 0)
		icw.addMemoryMetric("Name", "Unit", 10)

		if len(icw.metrics) != 1 {
			t.Fatalf("len(icw.metrics) want:1 got:%v", len(icw.metrics))
		}
	})
}

func TestAddMemoryMegabytes(t *testing.T) {
	t.Run("AddMemoryMegabytes", func(t *testing.T) {
		icw := instanceCloudwatch{}
		icw.metrics = make([]*cloudwatch.MetricDatum, 0)
		icw.AddMemoryMegabytes("Name", 10)

		metrics := make([]cloudwatch.MetricDatum, 0)
		for _, m := range icw.metrics {
			if *m.MetricName == "Name" {
				metrics = append(metrics, *m)
			}
		}

		if len(icw.metrics) != 1 {
			t.Fatalf("len(icw.metrics) want:1 got:%v", len(icw.metrics))
		}
		if *icw.metrics[0].Unit != "Megabytes" {
			t.Fatalf("icw.metrics[0].Unit want:Megabytes got:%v", icw.metrics[0].Unit)
		}
	})
}

func TestAddMemoryPercent(t *testing.T) {
	t.Run("AddMemoryPercent", func(t *testing.T) {
		icw := instanceCloudwatch{}
		icw.metrics = make([]*cloudwatch.MetricDatum, 0)
		icw.AddMemoryPercent("Name", 10)

		metrics := make([]cloudwatch.MetricDatum, 0)
		for _, m := range icw.metrics {
			if *m.MetricName == "Name" {
				metrics = append(metrics, *m)
			}
		}

		if len(icw.metrics) != 1 {
			t.Fatalf("len(icw.metrics) want:1 got:%v", len(icw.metrics))
		}
		if *icw.metrics[0].Unit != "Percent" {
			t.Fatalf("icw.metrics[0].Unit want:Percent got:%v", icw.metrics[0].Unit)
		}
	})
}
