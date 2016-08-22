# memdisk-cloudwatch 

## Monitoring Memory and Disk Metrics for Amazon EC2 Linux Instances

This is the replacement for example CloudWatch scripts by Amazon (CloudWatchMonitoringScripts, see [Monitoring Memory and Disk Metrics for Amazon EC2 Linux Instances](http://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/mon-scripts.html) ) 

This program is written in Go, the binary is statically linked and does not require any dependencies. 

These monitoring program is intended for use with Amazon EC2 instances running Linux operating systems. 
It has been tested on the 64-bit versions of the following Amazon Machine Images (AMIs):

- Amazon Linux 2014.09.2
- Ubuntu Server 16.04
- CentOS 6.x 

## Metrics 

The program only reports **MemoryUtilization** (percentage) and **DiskSpaceUtilization** (percentage).

**DiskSpaceUtilization** is reported only for root volume. 

There is no other configuration. Namespace is `System/Linux`, default metric names and units are the same as
 CloudWatchClient.pm / mon-put-instance-data.pl.
 
## Installation 

See [install_systemd.sh](install_systemd.sh) and [install_sysv.sh](install_sysv.sh)

It is considered to be insecure to run them as: 

```sh 
curl -L https://raw.githubusercontent.com/AndrianBdn/memdisk-cloudwatch/master/install_systemd.sh | sh 
```

or 

```sh 
curl -L https://raw.githubusercontent.com/AndrianBdn/memdisk-cloudwatch/master/install_sysv.sh | sh 
```
