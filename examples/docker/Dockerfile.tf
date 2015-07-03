# Dockerfile example running ntpd
FROM {{.Image}}:{{.ImageVersion}}

# Configure timezone
RUN ln -sf /usr/share/zoneinfo/{{.TimeZone}} /etc/localtime
COPY etc/sysconfig/clock /etc/sysconfig/clock

# Configure ntpdate
RUN yum install -y ntpdate 
COPY etc/sysconfig/ntpdate /etc/sysconfig/ntpdate

# Remove chrony since we're using ntpd
RUN yum remove -y chrony

# Configure ntpd
RUN yum install -y ntp 
COPY etc/ntp.conf /etc/ntp.conf
