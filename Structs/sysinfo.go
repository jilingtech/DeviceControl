package Structs

type DiskInfo struct {
	FreeSpace int `json:"free_space,omitempty"`
	Fstype string `json:"fstype,omitempty"`
	TotalSpace int `json:"total_space,omitempty"`
}

type Environment struct {
	GoPath string `json:"GOPATH,omitempty"`
	ZfsPath string `json:"ZFS_PATH,omitempty"`
}

type Memory struct {
	Swap int `json:"swap,omitempty"`
	Virt int `json:"virt,omitempty"`
}

type Net struct {
	InterfaceAddresses []string `json:"interface_addresses,omitempty"`
	Online bool `json:"online,omitempty"`
}

type RunTime struct {
	Arch string `json:"arch,omitempty"`
	Compiler string `json:"compiler,omitempty"`
	Gomaxprocs int `json:"gomaxprocs,omitempty"`
	NumCpu int `json:"numcpu,omitempty"`
	NumGoroutines int `json:"numgoroutines,omitempty"`
	Os string `json:"os,omitempty"`
	Version string `json:"version,omitempty"`
}


type SysInfo struct {
	DiskInfo *DiskInfo `json:"diskinfo,omitempty"`
	Environment *Environment `json:"environment,omitempty"`
	Memory *Memory `json:"memory,omitempty"`
	Net *Net `json:"net,omitempty"`
	RunTime *RunTime `json:"runtime,omitempty"`
	ZfsCommit string `json:"zfs_commit,omitempty"`
	ZfsVersion string `json:"zfs_version,omitempty"`
}
