# Dockerfile example running ntpd
FROM centos:centos7.1.1503

# Configure timezone
RUN ln -sf /usr/share/zoneinfo/Etc/UTC /etc/localtime
COPY etc/sysconfig/clock /etc/sysconfig/clock

# Configure ntpdate
RUN yum install -y ntpdate 
COPY etc/sysconfig/ntpdate /etc/sysconfig/ntpdate

# Configure ntpd
RUN yum install -y ntp 
COPY etc/ntp.conf /etc/ntp.conf
