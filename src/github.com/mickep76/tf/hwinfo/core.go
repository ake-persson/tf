package hwinfo

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func cpuInfo() (map[string]string, error) {
	d := make(map[string]string)
	logical := int64(0)

	b, err := ioutil.ReadFile("/proc/cpuinfo")
	if err != nil {
		return map[string]string{}, errors.New(fmt.Sprintf("can't read file: %s", err))
	}

	cpu_id := int64(-1)
	cpu_ids := make(map[int64]int64)
	cores := int64(0)
	for _, line := range strings.Split(string(b), "\n") {
		values := strings.Split(line, ":")
		if len(values) < 1 {
			continue
		} else if _, ok := d["cpu.model"]; !ok && strings.HasPrefix(line, "model name") {
			d["cpu.model"] = strings.Trim(strings.Join(values[1:], " "), " ")
		} else if _, ok := d["cpu.flags"]; !ok && strings.HasPrefix(line, "flags") {
			d["cpu.flags"] = strings.Trim(strings.Join(values[1:], " "), " ")
		} else if _, ok := d["cpu cores"]; !ok && strings.HasPrefix(line, "cpu cores") {
			cores, _ = strconv.ParseInt(strings.Trim(strings.Join(values[1:], " "), " "), 10, 0)
		} else if strings.HasPrefix(line, "processor") {
			logical++
		} else if strings.HasPrefix(line, "physical id") {
			cpu_id, _ = strconv.ParseInt(strings.Trim(strings.Join(values[1:], " "), " "), 10, 0)
			cpu_ids[cpu_id] = cpu_ids[cpu_id] + 1
		}
	}
	d["cpu.logical"] = strconv.FormatInt(logical, 10)
	sockets := int64(len(cpu_ids))
	d["cpu.sockets"] = strconv.FormatInt(sockets, 10)
	d["cpu.cores_per_socket"] = strconv.FormatInt(cores, 10)
	physical := int64(len(cpu_ids)) * cores
	d["cpu.physical"] = strconv.FormatInt(physical, 10)
	t := logical / sockets / cores
	d["cpu.threads_per_core"] = strconv.FormatInt(t, 10)

	return d, nil
}

func loadFiles(files map[string]string) (map[string]string, error) {
	d := make(map[string]string)

	for k, v := range files {
		if _, err := os.Stat(v); os.IsNotExist(err) {
			return map[string]string{}, errors.New("file doesn't exist")
		}

		b, err := ioutil.ReadFile(v)
		if err != nil {
			return map[string]string{}, errors.New(fmt.Sprintf("can't read file: %s", err))
		}

		d[k] = strings.Trim(string(b), "\n")
	}

	return d, nil
}

func loadFile(file string, del string, fields map[string]string) (map[string]string, error) {
	d := make(map[string]string)

	out, err := ioutil.ReadFile(file)
	if err != nil {
		return map[string]string{}, errors.New(fmt.Sprintf("can't read file: %s", err))
	}

	for _, line := range strings.Split(string(out), "\n") {
		values := strings.Split(line, del)
		if len(values) < 1 {
			continue
		}

		for k, v := range fields {
			if strings.HasPrefix(line, v) {
				d[k] = strings.Trim(strings.Join(values[1:], " "), " \t")
			}
		}
	}

	return d, nil
}

func execCmd(cmd string, args []string, del string, fields map[string]string) (map[string]string, error) {
	d := make(map[string]string)

	out, err := exec.Command(cmd, args...).Output()
	if err != nil {
		return map[string]string{}, err
	}

	for _, line := range strings.Split(string(out), "\n") {
		values := strings.Split(line, del)
		if len(values) < 1 {
			continue
		}

		for k, v := range fields {
			if strings.HasPrefix(line, v) {
				d[k] = strings.Trim(strings.Join(values[1:], " "), " \t")
			}
		}
	}

	return d, nil
}

func merge(a map[string]string, b map[string]string) {
	for k, v := range b {
		a[k] = v
	}
}

func HWInfo() (map[string]string, error) {
	sys_files := map[string]string{
		"serial_number":   "/sys/devices/virtual/dmi/id/product_serial",
		"manufacturer":    "/sys/devices/virtual/dmi/id/chassis_vendor",
		"product_version": "/sys/devices/virtual/dmi/id/product_version",
		"product":         "/sys/devices/virtual/dmi/id/product_name",
		"bios.date":       "/sys/devices/virtual/dmi/id/bios_date",
		"bios.vendor":     "/sys/devices/virtual/dmi/id/bios_vendor",
		"bios.version":    "/sys/devices/virtual/dmi/id/bios_version",
	}

	sysctl_fields := map[string]string{
		"mem.total.b":          "hw.memsize",
		"cpu.cores_per_socket": "machdep.cpu.core_count",
		"cpu.physical":         "hw.physicalcpu_max",
		"cpu.logical":          "hw.logicalcpu_max",
		"cpu.model":            "machdep.cpu.brand_string",
		"cpu.flags":            "machdep.cpu.features",
	}

	sw_vers_fields := map[string]string{
		"os.name":    "ProductName",
		"os.version": "ProductVersion",
	}

	lsb_release_fields := map[string]string{
		"os.name":    "Distributor ID",
		"os.version": "Release",
	}

	meminfo_fields := map[string]string{
		"mem.total.kb": "MemTotal",
	}

	sys := make(map[string]string)

	sys["os.kernel"] = runtime.GOOS

	h, err := os.Hostname()
	if err != nil {
		return map[string]string{}, err
	}
	sys["fqdn"] = h

	addrs, _ := net.LookupIP(h)
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			sys["fqdn_ip"] = ipv4.String()
		}
	}

	switch runtime.GOOS {
	case "darwin":
		o, err := execCmd("/usr/sbin/sysctl", []string{"-a"}, ":", sysctl_fields)
		if err != nil {
			return map[string]string{}, err
		}

		merge(sys, o)

		b, err := strconv.ParseUint(sys["mem.total.b"], 10, 64)
		if err != nil {
			return map[string]string{}, err
		} else {
			kb := b / 1024
			mb := kb / 1024
			gb := mb / 1024
			sys["mem.total.kb"] = strconv.FormatUint(kb, 10)
			sys["mem.total.mb"] = strconv.FormatUint(mb, 10)
			sys["mem.total.gb"] = strconv.FormatUint(gb, 10)
		}

		c, _ := strconv.ParseUint(sys["cpu.cores_per_socket"], 10, 64)
		p, _ := strconv.ParseUint(sys["cpu.physical"], 10, 64)
		l, _ := strconv.ParseUint(sys["cpu.logical"], 10, 64)
		s := p / c
		sys["cpu.sockets"] = strconv.FormatUint(s, 10)
		t := l / s / c
		sys["cpu.threads_per_core"] = strconv.FormatUint(t, 10)

		sys["cpu.flags"] = strings.ToLower(sys["cpu.flags"])

		o2, err2 := execCmd("/usr/bin/sw_vers", []string{}, ":", sw_vers_fields)
		if err2 != nil {
			return map[string]string{}, err2
		}

		merge(sys, o2)

	case "linux":
		o, err := loadFiles(sys_files)
		if err != nil {
			return map[string]string{}, err
		}
		merge(sys, o)

		if strings.Contains(sys["product_version"], "amazon") {
			sys["virtual"] = "Amazon EC2"
		}

		o2, err2 := execCmd("/usr/bin/lsb_release", []string{"-a"}, ":", lsb_release_fields)
		if err2 != nil {
			return map[string]string{}, err2
		}
		merge(sys, o2)

		o3, err3 := cpuInfo()
		if err3 != nil {
			return map[string]string{}, err3
		}
		merge(sys, o3)

		o4, err4 := loadFile("/proc/meminfo", ":", meminfo_fields)
		if err4 != nil {
			return map[string]string{}, err4
		}
		merge(sys, o4)

		sys["mem.total.kb"] = strings.Trim(sys["mem.total.kb"], " kB")

		kb, err := strconv.ParseUint(sys["mem.total.kb"], 10, 64)
		if err != nil {
			return map[string]string{}, err
		} else {
			b := kb * 1024
			mb := kb / 1024
			gb := mb / 1024
			sys["mem.total.b"] = strconv.FormatUint(b, 10)
			sys["mem.total.mb"] = strconv.FormatUint(mb, 10)
			sys["mem.total.gb"] = strconv.FormatUint(gb, 10)
		}

	default:
		return map[string]string{}, errors.New(fmt.Sprintf("unsupported plattform (%s), needs to be either linux or darwin", runtime.GOOS))
	}

	return sys, nil
}
