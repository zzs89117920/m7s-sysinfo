package sysinfo

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	. "m7s.live/engine/v4"
	"m7s.live/engine/v4/config"
)
type SysInfoConfig struct {
    config.HTTP
		
}

var SysInfoPluginConfig = &SysInfoConfig{}

var _ = InstallPlugin(SysInfoPluginConfig)

func (conf *SysInfoConfig) OnEvent(event any){
 	switch event.(type) {
    case FirstConfig: //插件初始化逻辑
    case *Stream://按需拉流逻辑
    case SEwaitPublish://由于发布者掉线等待发布者
    case SEpublish://首次进入发布状态
    case SErepublish://再次进入发布状态
    case SEwaitClose://由于最后一个订阅者离开等待关闭流
    case SEclose://关闭流
    case UnsubscribeEvent://订阅者离开
  }
}
type Result[T []*DiskInfo | *NetworkInfo | *CpuInfo | *MemInfo] struct {
	Code int
	Msg string
	Data T
}
var (
	timeFormat = "2006-01-02 15:04:05"
)
type MemInfo struct {
	Time string
	Percent float64
}

func (p *SysInfoConfig) API_MemInfo(w http.ResponseWriter, r *http.Request) {
    memInfo, error := mem.VirtualMemory()
		var result Result[*MemInfo] = Result[*MemInfo]{ }
		if(error == nil){
			
			memInfo := &MemInfo{
				Time: time.Now().Format(timeFormat),
				Percent: memInfo.UsedPercent,
			}

			result.Code = 200
			result.Data = memInfo
			result.Msg = ""
			jsonBytes, _ := json.Marshal(result)
			w.Write([]byte(jsonBytes))
			return
		}else{
			result.Code = 400
			result.Msg = error.Error()
			jsonBytes, _ := json.Marshal(result)
			w.Write([]byte(jsonBytes))
			return
		}
}

type CpuInfo struct {
	Time string
	Percent float64
}
func (p *SysInfoConfig) API_CpuInfo(w http.ResponseWriter, r *http.Request) {
    percent, error := cpu.Percent(time.Second, false)
		var result Result[*CpuInfo] = Result[*CpuInfo]{ }
		if(error == nil){

			cpuInfo := &CpuInfo{
				Time: time.Now().Format(timeFormat),
				Percent: percent[0],
			}

			result.Code = 200
			result.Data = cpuInfo
			result.Msg = ""
			jsonBytes, _ := json.Marshal(result)
			w.Write([]byte(jsonBytes))
			return
		}else{
			result.Code = 400
			result.Msg = error.Error()
			jsonBytes, _ := json.Marshal(result)
			w.Write([]byte(jsonBytes))
			return
		}
}


type DiskInfo struct {
	Time string
	Device   string
	Used   uint64
	Free   uint64
}
func (p *SysInfoConfig) API_DiskInfo(w http.ResponseWriter, r *http.Request) {
		parts, error := disk.Partitions(true)
		var result Result[[]*DiskInfo] = Result[[]*DiskInfo]{ }

		if(error == nil){
			
			var diskInfos []*DiskInfo
			for _, part := range parts {
        diskInfo, _ := disk.Usage(part.Mountpoint)
				
				diskInfos = append(diskInfos, &DiskInfo{ 
						Time: time.Now().Format(timeFormat),
						Device:part.Device, 
						Used:diskInfo.Used, 
						Free: diskInfo.Free,
					})
    	}

			result.Code = 200
			result.Data = diskInfos
			result.Msg = ""
			jsonBytes, _ := json.Marshal(result)
			w.Write([]byte(jsonBytes))
			return
		}else{
			result.Code = 400
			result.Msg = error.Error()
			jsonBytes, _ := json.Marshal(result)
			w.Write([]byte(jsonBytes))
			return
		}
}

type NetworkInfo struct {
	Time string
	Download uint64
	Upload uint64
}

func (p *SysInfoConfig) API_NetworkInfo(w http.ResponseWriter, r *http.Request) {
    info, error := net.IOCounters(false)
		var result Result[*NetworkInfo] = Result[*NetworkInfo]{ }
		if(error == nil){
			bytesRecv := info[0].BytesRecv
			bytesSend := info[0].BytesSent
			timer := time.NewTimer(1 * time.Second)

		  <-timer.C

			info1, _ := net.IOCounters(false)
		 	bytesRecv1 := info1[0].BytesRecv
			bytesSend1 := info1[0].BytesSent

			newBytesRecv := bytesRecv1 - bytesRecv
			newBytesSend := bytesSend1 - bytesSend

			networkInfo := &NetworkInfo{
				Time: time.Now().Format(timeFormat),
				Download: newBytesRecv,
				Upload: newBytesSend,
			}

			result.Code = 200
			result.Data = networkInfo
			result.Msg = ""
			jsonBytes, _ := json.Marshal(result)
			w.Write([]byte(jsonBytes))
			return
		}else{
			result.Code = 400
			result.Msg = error.Error()
			jsonBytes, _ := json.Marshal(result)
			w.Write([]byte(jsonBytes))
			return
		}
}