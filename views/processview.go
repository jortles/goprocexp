package views

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shirou/gopsutil/process"
	"math"
	"sort"
	"strconv"
)

type ProcCategories struct {
	ProcessName string
	ProcessID   int32
	CPUPercent  float64
	Username    string
	ProcessExec string
}

func (p *ProcCategories) ProcessList() {
	procs, _ := process.Processes()
	for _, proc := range procs {
		p.ProcessName, _ = proc.Name()
		p.ProcessID = proc.Pid
		p.CPUPercent, _ = proc.CPUPercent()
		p.Username, _ = proc.Username()
		p.ProcessExec, _ = proc.Cmdline()
	}
}

func (p *ProcCategories) ProcessNames() string {
	procs, _ := process.Processes()
	for _, proc := range procs {
		p.ProcessName, _ = proc.Name()
	}

	return p.ProcessName
}

func (p *ProcCategories) ProcessId() int32 {
	procs, _ := process.Processes()
	for _, proc := range procs {
		p.ProcessID = proc.Pid
	}

	return p.ProcessID
}

func (p *ProcCategories) ProcessCpu() float64 {
	procs, _ := process.Processes()
	for _, proc := range procs {
		p.CPUPercent, _ = proc.CPUPercent()
	}

	return p.CPUPercent
}

func (p *ProcCategories) ProcessUser() string {
	procs, _ := process.Processes()
	for _, proc := range procs {
		p.Username, _ = proc.Username()
	}

	return p.Username
}

func (p *ProcCategories) ProcessCmd() string {
	procs, _ := process.Processes()
	for _, proc := range procs {
		p.ProcessExec, _ = proc.Exe()
	}

	return p.ProcessExec
}

func (p *ProcCategories) ProcessExe() []string {
	procs, _ := process.Processes()
	var allExe []string
	for _, proc := range procs {
		p.ProcessExec, _ = proc.Cmdline()
		allExe = append(allExe, p.ProcessExec)
	}

	return allExe
}

func (p *ProcCategories) ProcessIds() []string {
	procs, _ := process.Processes()
	var allPids []string
	for _, proc := range procs {
		p.ProcessID = proc.Pid
		allPids = append(allPids, strconv.Itoa(int(p.ProcessID)))
	}
	return allPids
}

func (p *ProcCategories) ProcessCpus() []float64 {
	procs, _ := process.Processes()
	var cpuAll []float64
	for _, proc := range procs {
		p.CPUPercent, _ = proc.CPUPercent()
		cpuAll = append(cpuAll, math.Round(p.CPUPercent*100)/100)
		sort.Sort(sort.Reverse(sort.Float64Slice(cpuAll)))
	}
	return cpuAll
}

func (p *ProcCategories) ProcessUsers() []string {
	procs, _ := process.Processes()
	var procUsers []string
	for _, proc := range procs {
		p.Username, _ = proc.Username()
		procUsers = append(procUsers, p.Username)
	}
	return procUsers
}

func ProcessExeView() (processExes *tview.TextView) {
	var p ProcCategories
	procExe := p.ProcessExe()
	processExes = tview.NewTextView()

	processExes.SetBorder(true).SetTitle(" PROCESS DETAILS ").SetBorderColor(tcell.ColorIndigo)
	for _, e := range procExe {
		fmt.Fprintf(processExes, "%s\n", e)
	}

	return processExes
}

func ProcessNamesView() (processNames *tview.TextView) {
	var p ProcCategories
	procNames := p.ProcessNames()
	processNames = tview.NewTextView()

	processNames.SetBorder(true).SetTitle("> Process Name <").SetBorderColor(tcell.ColorDarkCyan)
	for _, n := range procNames {
		fmt.Fprintf(processNames, "%v\n", n)
	}

	return processNames
}
