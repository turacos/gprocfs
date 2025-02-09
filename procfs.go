package main


import (

  "fmt"
  "time"
  "sync"
  "github.com/shirou/gopsutil/v4/process"
  "github.com/shirou/gopsutil/v4/cpu"
  "github.com/shirou/gopsutil/v4/load"
  "github.com/shirou/gopsutil/v4/host"
  "github.com/shirou/gopsutil/v4/disk"
  "github.com/shirou/gopsutil/v4/net" 
  
)


func debug(a [][]string) {

  for _,v := range a {
  
    fmt.Printf("%v\n",v)
  
  } 

}

func getProcesses() [][]string {

  var a [][]string
  
  p,_ := process.Processes()
  
  a = append(a,[]string{"PID","PPID","Name","cpupct","mempct","RSS","","HWM","Data","Stack","Locked","Swap","Start time","Username","Cmd"})
  

  for _,pp := range p {
  
  mempct,_ := pp.MemoryPercent()
  cpupct,_ := pp.CPUPercent()
  name,_ := pp.Name()
  ppid,_ := pp.Ppid()
  cmd,_ := pp.Cmdline()
  memStat,_ := pp.MemoryInfo()
  ttime,_ :=  pp.CreateTime()
  tm := time.UnixMilli(ttime)
  usrname,_ := pp.Username()
  
  a = append(a,[]string{
  					fmt.Sprintf("%d",pp.Pid),
  					fmt.Sprintf("%d",ppid),
  					name,
  					fmt.Sprintf("%f",cpupct),
  					fmt.Sprintf("%f",mempct),
  					fmt.Sprintf("%d",memStat.RSS),
  					fmt.Sprintf("%d",memStat.VMS),
  					fmt.Sprintf("%d",memStat.HWM),
  					fmt.Sprintf("%d",memStat.Data),
  					fmt.Sprintf("%d",memStat.Stack),
  					fmt.Sprintf("%d",memStat.Locked),
  					fmt.Sprintf("%d",memStat.Swap),
  					tm.Format(time.RFC3339),
  					usrname,
  					cmd})
  }
  
  
  
  return a
}

func getHardware() [][]string {

  var a [][]string
  
  h,_ := cpu.Info()
  
  a = append(a,[]string{"Name","Value"})
  

  nbc := fmt.Sprintf("%d",len(h))

  a = append(a,[]string{"numher of cpu threads",nbc })

  hi,_ := host.Info()
  bttm := time.Unix(int64(hi.BootTime),0)  
    
    
  a = append(a,[]string{"Hostname", hi.Hostname})
//  a = append(a,[]string{"Uptime",uttm.Format(time.RFC3339) })
  a = append(a,[]string{"BootTime",bttm.Format(time.RFC3339) })
  a = append(a,[]string{"number of processes",fmt.Sprintf("%d",hi.Procs)})
  a = append(a,[]string{"OS",hi.OS })
  a = append(a,[]string{"Platform",hi.Platform }) 
  a = append(a,[]string{"PlatformFamily",hi.PlatformFamily })  
  a = append(a,[]string{"PlatformVersion",hi.PlatformVersion })
  a = append(a,[]string{"KernelVersion",hi.KernelVersion }) 
  a = append(a,[]string{"KernelArch",hi.KernelArch })  
  a = append(a,[]string{"VirtualizationSystem",hi.VirtualizationSystem })
  a = append(a,[]string{"VirtualizationRole",hi.VirtualizationRole }) 
  a = append(a,[]string{"HostID",hi.HostID })  

//  l,_ :=  load.Misc()
  
 // fmt.Printf("%+v\n",l)
      
  return a
}


func getCPULoad() [][]string {

  var a [][]string
  
  l,_ := load.Avg()
  
  a = append(a,[]string{"load1","load5","load15"})
  
   // fmt.Printf("%v\n",pp)
  a = append(a,[]string{
    		fmt.Sprintf("%.2f",l.Load1),
				fmt.Sprintf("%.2f",l.Load5),
    		fmt.Sprintf("%.2f",l.Load15)})

 // fmt.Printf("%v\n",l)
  
  return a
}

func getCPUPercent() []string {

  var a []string
  
  c,_ := cpu.Percent(100 * time.Millisecond,true)
  
  for _,p := range c {
    a = append(a,fmt.Sprintf("%.2f",p))
  }
  
  return a

}
	
func getCPUStats() [][]string {

  var a [][]string
  
  h,_ := cpu.Times(true)
  
  a = append(a,[]string{"CPU","CPUpct","User","System","Idle","Nice","Iowait","Irq","Softirq","Steal","Guest","GuestNice"})
  
  c := getCPUPercent()
  
  i := 0
  for _,pp := range h {
  
   // fmt.Printf("%v\n",pp)
  a = append(a,[]string{
     		pp.CPU, 
        c[i],      
    		fmt.Sprintf("%f",pp.User),
				fmt.Sprintf("%f",pp.System),
    		fmt.Sprintf("%f",pp.Idle),
    		fmt.Sprintf("%f",pp.Nice),
				fmt.Sprintf("%f",pp.Iowait),
    		fmt.Sprintf("%f",pp.Irq),
    		fmt.Sprintf("%f",pp.Softirq),
				fmt.Sprintf("%f",pp.Steal),
    		fmt.Sprintf("%f",pp.Guest),
      	fmt.Sprintf("%f",pp.GuestNice)})
    i=i+1
  }
  
  //fmt.Printf("%v\n",h)
  
  return a
}


func getDiskUsage() [][]string {

  var a [][]string
  
  p,_ := disk.Partitions(true)

  a = append(a,[]string{"Path","Fstype","Total","Free","Used","UsedPercent","InodesTotal","InodesUsed","InodesFree","InodesUsedPercent"})
  
  
  for _,pp := range p {
  
    d,err := disk.Usage(pp.Mountpoint)
    if err == nil {
    //fmt.Printf("%s - %v\n",pp.Mountpoint,err)
//fmt.Printf("%v\n",d)
        if d.Total > 0 {
				a = append(a,[]string{d.Path,
				d.Fstype,
				fmt.Sprintf("%d",d.Total),
				fmt.Sprintf("%d",d.Free),
				fmt.Sprintf("%d",d.Used),
				fmt.Sprintf("%.2f",d.UsedPercent),
				fmt.Sprintf("%d",d.InodesTotal),
				fmt.Sprintf("%d",d.InodesUsed),
				fmt.Sprintf("%d",d.InodesFree),
				fmt.Sprintf("%.2f",d.InodesUsedPercent)})		
        }
    }
  }
  
  return a

}


func getNetDev() [][]string {

  var a [][]string
  var f string
  var addr string

  ios := make(map[string][]string)
  i,_ := net.Interfaces()

  io,_ := net.IOCounters(true)
  
 // fmt.Printf("%+v\n",io)

  a = append(a,[]string{"Name","HardwareAddr","Status","Addrs","BytesSent","BytesRecv","PacketsSent","PacketsRecv","Errin","Errout","Dropin","Dropout"})
 
  for _,i := range io { 
  

     ios[i.Name] = []string{fmt.Sprintf("%d",i.BytesSent),
     fmt.Sprintf("%d",i.BytesRecv),
     fmt.Sprintf("%d",i.PacketsSent),
     fmt.Sprintf("%d",i.PacketsRecv),
     fmt.Sprintf("%d",i.Errin),
     fmt.Sprintf("%d",i.Errout),
     fmt.Sprintf("%d",i.Dropin),
     fmt.Sprintf("%d",i.Dropout)}

  }
  
  for _,n := range i {
 
        if len(n.Flags) == 3 {
          f = n.Flags[0]
        }else{
          f=""
        }
        if len(n.Addrs) > 0 {
          addr = n.Addrs[0].Addr
        }else{
         addr=""
        } 
        
 				a = append(a,[]string{n.Name,
				n.HardwareAddr,
				f,
				addr,
				ios[n.Name][0],
				ios[n.Name][1],
				ios[n.Name][2],
				ios[n.Name][3],
				ios[n.Name][4],
				ios[n.Name][5],
				ios[n.Name][6],
				ios[n.Name][7]})	 
  
  }
  

  return a

}

func execCommand(f func() [][]string, i uint) {
  
  for {
    debug(f())
  //  f()
    time.Sleep( time.Duration(i) * time.Second )
  }
}


func main() {


//debug(getProcesses())
//debug(getHardware())
//debug(getCPUStats())
//debug(getCPULoad())
//debug(getCPUPercent())
//debug(getNetDev())

var wg sync.WaitGroup

	
//	wg.Add(2)
wg.Add(1)	
	go execCommand(getNetDev,20)
/*
wg.Add(1)
	go execCommand(getProcesses,20)
wg.Add(1)	
	go execCommand(getHardware,120)
wg.Add(1)	
	go execCommand(getCPUStats,20)
wg.Add(1)	
	go execCommand(getDiskUsage,20)
wg.Add(1)	
	go execCommand(getNetDev,20)
*/
wg.Wait()
/*
  for {
  
debug(getProcesses())
debug(getHardware())
debug(getCPUStats())
//debug(getCPULoad())/debug(getCPUPercent())
debug(getDiskUsage())
debug(getNetDev())
      time.Sleep( 20000 * time.Millisecond)
	//	wg.Wait()
  }
*/
}
