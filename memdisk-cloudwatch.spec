# don't install debuginfo -- this package is just a Go binary
%global debug_package %{nil}
%global __debug_install_post /bin/true

Name:           memdisk-cloudwatch
Version:        0.9.2
Release:        1%{?dist}
Summary:        Monitoring Memory and Disk Metrics for Amazon EC2 Linux Instances
Group:          Development/Tools
License:        MIT
URL:            https://github.com/AndrianBdn/memdisk-cloudwatch/
Source0:        %{name}
Source1:        %{name}.service
BuildArch:      x86_64
BuildRequires:  systemd

%description
This is the replacement for example CloudWatch scripts by Amazon
(CloudWatchMonitoringScripts, see Monitoring Memory and Disk Metrics for Amazon
EC2 Linux Instances.)

This program is written in Go, the binary is statically linked and does not
require any dependencies.

This monitoring program is intended for use with Amazon EC2 instances running
Linux operating systems. It has been tested on the 64-bit versions of the
following Amazon Machine Images (AMIs):

Amazon Linux 2014.09.2
Ubuntu Server 16.04
CentOS 6.x

# not using %setup since the upstream package is not a tarball
%prep

%build

%install
%{__install} -p -m 0755 -D %{SOURCE0} %{buildroot}%{_bindir}/%{name}
%{__install} -p -m 0644 -D %{SOURCE1} %{buildroot}%{_unitdir}/%{name}.service

%post
%systemd_post %{name}.service

%preun
%systemd_preun %{name}.service

%postun
%systemd_postun_with_restart %{name}.service

%clean
rm -rf %{buildroot}

%files
%defattr(-,root,root,-)
%{_bindir}/%{name}
%{_unitdir}/%{name}.service

%changelog
* Wed Oct 17 2018 Evan Zacks <zackse@gmail.com> 0.9.2-1
- Initial RPM package.
