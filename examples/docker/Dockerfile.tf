# Dockerfile example running ntpd
FROM {{.Image}}:{{.ImageVersion}}

# Configure timezone
RUN ln -sf /usr/share/zoneinfo/{{.File.TimeZone}} /etc/localtime
COPY etc/sysconfig/clock /etc/sysconfig/clock

# Configure ntpdate
RUN yum install -y ntpdate 
COPY etc/sysconfig/ntpdate /etc/sysconfig/ntpdate

# Configure ntpd
RUN yum install -y ntp 
COPY etc/ntp.conf /etc/ntp.conf
