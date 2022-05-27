package services

import (
	"fmt"
	"math"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

type SysinfoService struct {
}

type Sysinfo struct {
	CpuPer    float64 `json:"cpuper"`
	CpuCount  int     `json:"cpu_count"`
	MemPer    float64 `json:"memper"`
	MemTotal  uint64  `json:"mem_total"`
	DiskPer   float64 `json:"diskper"`
	DiskTotal uint64  `json:"disk_total"`
	HostName  string  `json:"hostname"`
	Uptime    string  `json:"uptime"`
	Os        string  `json:"os"`
	Arch      string  `json:"arch"`
}

// 截取float64的小数 prec 代表小数位数
func TrunFloat(f float64, prec int) float64 {
	x := math.Pow10(prec)
	return math.Trunc(f*x) / x
}

// 读取服务器运行状态信息
func (*SysinfoService) GetInfoByPsutil() (sysinfo Sysinfo, err error) {

	sysinfo = Sysinfo{}

	percent, err := cpu.Percent(time.Second, false)

	sysinfo.CpuPer = TrunFloat(percent[0], 2)
	sysinfo.CpuCount, err = cpu.Counts(true) //cpu逻辑数量

	memInfo, err := mem.VirtualMemory()
	sysinfo.MemPer = TrunFloat(memInfo.UsedPercent, 2)
	sysinfo.MemTotal = memInfo.Total / (1024 * 1024 * 1024) // 单位G

	parts, err := disk.Partitions(true)
	diskInfo, err := disk.Usage(parts[0].Mountpoint)

	sysinfo.DiskPer = TrunFloat(diskInfo.UsedPercent, 2)
	sysinfo.DiskTotal = diskInfo.Total / (1024 * 1024 * 1024) //单位G

	info, err := host.Info()

	sysinfo.HostName = info.Hostname
	sysinfo.Uptime = fmt.Sprintf("%d天%d小时%d分钟%d秒", info.Uptime/(3600*24), info.Uptime%(3600*24)/3600, info.Uptime%3600/60, info.Uptime%60)
	sysinfo.Os = info.OS
	sysinfo.Arch = info.KernelArch
	logs.Info(info)
	logs.Info(sysinfo)
	return
}
