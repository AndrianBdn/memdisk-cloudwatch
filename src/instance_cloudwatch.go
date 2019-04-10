package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

type instanceCloudwatch struct {
	cloudwatch       *cloudwatch.CloudWatch
	instansDimention *cloudwatch.Dimension
	instanceIdentity ec2metadata.EC2InstanceIdentityDocument
	memDimensions    []*cloudwatch.Dimension
	diskDimentions   []*cloudwatch.Dimension
	metrics          []*cloudwatch.MetricDatum
}

const CLOUDWATCH_NAMESPACE string = "System/Linux"

func NewInstanceCloudwatch(blockDevice, mountPath string) *instanceCloudwatch {
	ec2_sess := session.New()
	meta := ec2metadata.New(ec2_sess)

	instanceIdentity, err := meta.GetInstanceIdentityDocument()

	if err != nil {
		log.Fatal("Unable to get AWS instance identity, error: ", err)
	}

	cwSess := session.New(&aws.Config{Region: aws.String(instanceIdentity.Region)})

	retval := &instanceCloudwatch{
		cloudwatch:       cloudwatch.New(cwSess),
		instanceIdentity: instanceIdentity,
	}

	retval.prefill(blockDevice, mountPath)

	return retval
}

func (a *instanceCloudwatch) prefill(blockDevice, mountPath string) {

	a.instansDimention = &cloudwatch.Dimension{
		Name:  aws.String("InstanceId"),
		Value: aws.String(a.instanceIdentity.InstanceID),
	}

	a.memDimensions = append(a.memDimensions, a.instansDimention)
	a.diskDimentions = append(a.diskDimentions, a.instansDimention)
	a.diskDimentions = append(a.diskDimentions, &cloudwatch.Dimension{
		Name:  aws.String("Filesystem"),
		Value: aws.String(blockDevice),
	})
	a.diskDimentions = append(a.diskDimentions, &cloudwatch.Dimension{
		Name:  aws.String("MountPath"),
		Value: aws.String(mountPath),
	})
}

func (a *instanceCloudwatch) ResetMetrics() {
	a.metrics = make([]*cloudwatch.MetricDatum, 0)
}

func (a *instanceCloudwatch) addDiskMetric(name, unit string, value float64) {

	a.metrics = append(a.metrics, &cloudwatch.MetricDatum{
		MetricName: aws.String(name),
		Unit:       aws.String(unit),
		Value:      aws.Float64(value),
		Dimensions: a.diskDimentions,
	})

}

func (a *instanceCloudwatch) AddDiskGigabytes(name string, gigabytes float64) {
	a.addDiskMetric(name, "Gigabytes", gigabytes)
}

func (a *instanceCloudwatch) AddDiskPercent(name string, percent float64) {
	a.addDiskMetric(name, "Percent", percent)
}

func (a *instanceCloudwatch) addMemoryMetric(name, unit string, value float64) {
	a.metrics = append(a.metrics, &cloudwatch.MetricDatum{
		MetricName: aws.String(name),
		Unit:       aws.String(unit),
		Value:      aws.Float64(value),
		Dimensions: a.memDimensions,
	})
}

func (a *instanceCloudwatch) AddContainerMetric(name string, value int) {
	dimensions := []*cloudwatch.Dimension{
		a.instansDimention,
		&cloudwatch.Dimension{
			Name:  aws.String("ContainerName"),
			Value: aws.String(name),
		},
	}

	a.metrics = append(a.metrics, &cloudwatch.MetricDatum{
		MetricName: aws.String("ContainerStatus"),
		Unit:       aws.String("Count"),
		Value:      aws.Float64(float64(value)),
		Dimensions: dimensions,
	})
}

func (a *instanceCloudwatch) AddMemoryMegabytes(name string, megabytes float64) {
	a.addMemoryMetric(name, "Megabytes", megabytes)
}

func (a *instanceCloudwatch) AddMemoryPercent(name string, percent float64) {
	a.addMemoryMetric(name, "Percent", percent)
}

func (a *instanceCloudwatch) Send() {
	start := 0
	len := len(a.metrics)
	for start < len {
		end := start + 20
		if end > len {
			end = len
		}
		a.send(a.metrics[start:end])
		start = end
	}
}

func (a *instanceCloudwatch) send(metrics []*cloudwatch.MetricDatum) {
	_, err := a.cloudwatch.PutMetricData(&cloudwatch.PutMetricDataInput{
		MetricData: metrics,
		Namespace:  aws.String(CLOUDWATCH_NAMESPACE),
	})

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			log.Fatal("PutMetricData Error ", awsErr.Code(), ": ", awsErr.Message())
		}
	}
}
