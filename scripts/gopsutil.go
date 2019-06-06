package main

import (
	"encoding/json"
	"fmt"
	//"fmt"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	//"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	//"github.com/shirou/gopsutil/net"
)

type ResourceInfo struct {
	ID      int64
	HostID  int64
	GroupID int64
	//单位MB
	MemTotal       uint64
	MemUsedPercent float64
	//单位核
	CpuTotal       uint64
	CpuUsedPercent float64
	//空间大小单位GB
	DiskTotal       uint64
	DiskUsedPercent float64
	CreateTime      int64
}

func main() {
	v, _ := mem.VirtualMemory()
	c, _ := cpu.Info()
	cc, _ := cpu.Percent(time.Second, false)
	d, _ := disk.Usage("/")
	var cores int32
	for _, sub_cpu := range c {
		//modelname := sub_cpu.ModelName
		cores = cores + sub_cpu.Cores
		//fmt.Printf("        CPU       : %v   %v cores \n", modelname, cores)
	}

	//fmt.Printf("        Network: %v bytes / %v bytes\n", nv[0].BytesRecv, nv[0].BytesSent)
	//fmt.Printf("        SystemBoot:%v\n", btime)
	//for _, p := range cc {
	//	fmt.Printf("        CPU Used    : used %f%% \n", p)
	//}

	//fmt.Printf("        Disk        : %v GB  Free: %v GB Usage:%f%%\n", d.Total/1024/1024/1024, d.Free/1024/1024/1024, d.UsedPercent)
	//fmt.Printf("        OS        : %v(%v)   %v  \n", n.Platform, n.PlatformFamily, n.PlatformVersion)
	//fmt.Printf("        Hostname  : %v  \n", n.Hostname)
	resourceInfo := &ResourceInfo{
		MemTotal:        v.Total / 1024 / 1024,
		MemUsedPercent:  v.UsedPercent,
		CpuTotal:        uint64(cores),
		CpuUsedPercent:  cc[0],
		DiskTotal:       d.Total / 1024 / 1024 / 1024,
		DiskUsedPercent: d.UsedPercent,
		CreateTime:      time.Now().Unix(),
	}
	data, _ := json.Marshal(resourceInfo)
	fmt.Printf("data=============:%v  \n", string(data))
}
