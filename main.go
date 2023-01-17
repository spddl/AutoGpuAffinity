package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"unsafe"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	var err error

	UnzipFiles()
	defer Cleanup()

	if totaltrials == -1 {
		resultChan := make(chan int, 1)
		if _, err := tea.NewProgram(initialModel("Enter how many trials you would like to test for each CPU:", "3", resultChan), tea.WithAltScreen()).Run(); err != nil {
			log.Println(err)
			return
		}
		totaltrials = <-resultChan
	}

	if trialtime == -1 {
		resultChan := make(chan int, 1)
		if _, err := tea.NewProgram(initialModel("Enter how many seconds you want each trial to last:", "30", resultChan), tea.WithAltScreen()).Run(); err != nil {
			log.Println(err)
			return
		}
		trialtime = <-resultChan
	}

	estimated := ((((trialtime + 5) * totaltrials) * Cores) / 60)

	if !cliMode {
		resultChan := make(chan int, 1)
		if _, err := tea.NewProgram(rebootModel{choice: resultChan, choices: []string{fmt.Sprintf("Start Now (%d mins)", estimated), "Set AutoGpuAffinity on startup and reboot"}}, tea.WithAltScreen()).Run(); err != nil {
			log.Println(err)
			return
		}
		restart := <-resultChan
		if restart != 0 {
			ex, err := os.Executable()
			if err != nil {
				panic(err)
			}
			SetRunOnce(fmt.Sprintf("%s -totaltrials=%d -trialtime=%d", ex, totaltrials, trialtime))
			return
		}
	}

	GPUdevices, handle = FindAllDevices()
	defer SetupDiDestroyDeviceInfoList(handle)
	defaultSettings.GPUdevices = append([]Device(nil), GPUdevices...)

	fmt.Printf("DISCLAIMER: nobody is responsible if you damage your computer. run at your own risk\n\ndo not touch your PC at all while this tool runs to avoid collecting invalid data\n\nclose any background apps you have open\n\nestimated time for completetion: %d mins\n\n", estimated)
	time.Sleep(5 * time.Second)

	lava = Program{
		exe:      LavaPath,
		highPrio: true,
		async:    true,
	}

	if err := os.MkdirAll(filepath.Join(os.Getenv("AppData"), "liblava", "lava triangle"), os.ModePerm); err != nil {
		log.Println(err)
	}

	tempFolder, err = os.MkdirTemp("", "PresentMon")
	if err != nil {
		log.Println(err)
	}

	lavaconfig := filepath.Join(os.Getenv("AppData"), "liblava", "lava triangle", "window.json")
	if FileExists(lavaconfig) {
		err := os.Remove(lavaconfig)
		if err != nil {
			log.Println(err)
		}
	}
	createFile(lavaconfig, `{
	"default": {
		"decorated": true,
		"floating": false,
		"fullscreen": true,
		"height": 540,
		"maximized": true,
		"monitor": 0,
		"resizable": true,
		"width": 960,
		"x": 480,
		"y": 270
	}
}`)

	resultDone := make(chan struct{}, 1)
	go func(resultDone chan struct{}) {
		tea.NewProgram(progressbarModel{
			progress: progress.New(progress.WithGradient("#008000", "#00FF00")),
			callback: resultDone,
		}).Run()
		resultDone <- struct{}{}
	}(resultDone)

	percent100 := float64(totaltrials * Cores)
	var index float64
	for trial := 0; trial < totaltrials; trial++ {
		for cpu := 0; cpu < Cores; cpu++ {
			for {
				if benchmark(cpu, trial) {
					index++
					progressbarValue = (100 / percent100 * index) / 100
					break
				}
			}
		}
	}
	<-resultDone

	tableResult := []CPUResultList{}
	for _, cpulist := range result.CPUList {
		all := CPUResultList{}
		for _, v := range cpulist {
			all.Max += v.Max
			all.Avg += v.Avg
			all.Min += v.Min

			all.LowsPoint005 += v.LowsPoint005
			all.LowsPoint01 += v.LowsPoint01
			all.LowsPoint1 += v.LowsPoint1
			all.LowsOne += v.LowsOne

			all.Percent1 += v.Percent1
			all.Percent01 += v.Percent01
			all.Percent001 += v.Percent001
			all.Percent0005 += v.Percent0005
		}
		length := float64(len(cpulist))
		all.Max /= length
		all.Avg /= length
		all.Min /= length

		all.LowsPoint005 /= length
		all.LowsPoint01 /= length
		all.LowsPoint1 /= length
		all.LowsOne /= length

		all.Percent1 /= length
		all.Percent01 /= length
		all.Percent001 /= length
		all.Percent0005 /= length
		tableResult = append(tableResult, all)
	}
	tableOutput(tableResult)

	for i := range defaultSettings.GPUdevices {
		setAffinityPolicy(&defaultSettings.GPUdevices[i]) // set old settings, if you exit the program now and restart the PC everything is as before
	}

	Promt_cpu()

	if err := os.RemoveAll(tempFolder); err != nil {
		log.Println(err)
	}
}

func SetGPUandRestart(cpubit Bits) {
	for i := range GPUdevices {
		GPUdevices[i].DevicePolicy = 4
		GPUdevices[i].AssignmentSetOverride = cpubit
		setAffinityPolicy(&GPUdevices[i])

		propChangeParams := PropChangeParams{
			ClassInstallHeader: *MakeClassInstallHeader(DIF_PROPERTYCHANGE),
			StateChange:        DICS_PROPCHANGE,
			Scope:              DICS_FLAG_GLOBAL,
		}
		for c := 0; c < 2; c++ {
			if err := SetupDiSetClassInstallParams(handle, &GPUdevices[i].Idata, &propChangeParams.ClassInstallHeader, uint32(unsafe.Sizeof(propChangeParams))); err != nil {
				log.Println(err)
				return
			}
			if err := SetupDiCallClassInstaller(DIF_PROPERTYCHANGE, handle, &GPUdevices[i].Idata); err != nil {
				log.Println(err)
				return
			}
		}

		DeviceInstallParams, err := SetupDiGetDeviceInstallParams(handle, &GPUdevices[i].Idata)
		if err != nil {
			log.Println(err)
			return
		}

		if DeviceInstallParams.Flags&DI_NEEDREBOOT != 0 {
			fmt.Println("Device could not be restarted. Changes will take effect the next time you reboot.")
		}
	}
}

func benchmark(cpu, trial int) bool {
	if HT {
		SetGPUandRestart(CPUBits[cpu*2])
	} else {
		SetGPUandRestart(CPUBits[cpu])
	}
	WaitForDesktop()

	// Start Lava
	var pid uint32
	for {
		// if the program was started successfully there is also this pid
		var progExists bool
		pid, progExists = StartUntilItExists(&lava)
		if !progExists {
			continue
		}

		// to put the program in the foreground we need the window handle
		if !FindWindowAndSetOnTop("lava triangle") {
			if lava.kill() { // kill directly the process
				hwnd, _ := FindWindow("lava triangle")
				DestroyWindow(hwnd)
			}
			continue
		}
		break
	}

	tempfile := fmt.Sprintf("PresentMon_CPU%d_%d", cpu, trial)
	presentMon := Program{
		exe: PresentMonPath,
		arg: []string{
			"-session_name", tempfile,
			"-process_id", strconv.FormatInt(int64(pid), 10),
			"-output_file", filepath.Join(tempFolder, tempfile+".csv"),
			"-no_top",
			"-timed", fmt.Sprintf("%d", trialtime),
			"-terminate_after_timed",
			"-terminate_on_proc_exit",
		},
		highPrio: true,
	}
	presentMon.start()
	if presentMon.kill() {
		log.Println("error? -timed doesn't trigger")
	}

	if lava.kill() { // kill directly the process
		hwnd, _ := FindWindow("lava triangle")
		DestroyWindow(hwnd)
	}

	if FileExists(filepath.Join(tempFolder, tempfile+".csv")) {
		list := calc(filepath.Join(tempFolder, tempfile+".csv"))
		if trial == 0 {
			result.CPUList = append(result.CPUList, []CPUResultList{})
		}
		result.CPUList[cpu] = append(result.CPUList[cpu], list)
		return true
	}
	return false // run invalid
}

func WaitForDesktop() bool {
	for i := 0; i < 20; i++ {
		desktop := openInputDesktop(0, false, 0x02000000) // MAXIMUM_ALLOWED
		if desktop != 0 {
			closeDesktop(desktop)
			return true
		}
		time.Sleep(200 * time.Millisecond)
	}
	return false
}
