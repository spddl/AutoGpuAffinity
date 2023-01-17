package main

import (
	"log"
	"os/exec"
	"syscall"
	"time"
)

type Program struct {
	exe      string
	arg      []string
	highPrio bool
	async    bool
	cmd      *exec.Cmd
}

func (p *Program) start() uint32 {
	p.cmd = exec.Command(p.exe, p.arg...)

	// p.cmd.Stdout = os.Stdout
	// p.cmd.Stderr = os.Stderr
	if p.highPrio {
		p.cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x00000080} // HIGH_PRIORITY_CLASS
	}

	if p.async {
		err := p.cmd.Start()
		if err != nil {
			log.Printf("Failed to start " + err.Error())
			log.Printf("%#v\n ", err)
		}

		return uint32(p.cmd.Process.Pid)
	} else {
		err := p.cmd.Run()
		if err != nil {
			log.Printf("Failed to start " + err.Error())
			log.Printf("%#v\n ", err)
		}
	}
	return 0
}

func (p *Program) kill() bool {
	if p.cmd.Process == nil {
		return false
	}
	return p.cmd.Process.Kill() == nil
}

func StartUntilItExists(prog *Program) (uint32, bool) {
	var index uint32
	var pid uint32
	for {
		time.Sleep(time.Second)
		pid = prog.start()
		if pid == 0 {
			if index > 10 {
				return 0, false
			}
			continue
		}
		running, _ := pid_is_running(pid)
		if !running {
			if index > 10 {
				return 0, false
			}
			continue
		}

		index++
		break
	}
	return pid, true
}

func FindWindowAndSetOnTop(title string) bool {
	var index uint32
	for {
		index++
		time.Sleep(time.Second)
		hwnd, err := FindWindow(title)
		if err != nil || hwnd == 0 {
			if index > 10 {
				return false // Timeout
			}
			continue
		}

		if SetTopWindow(hwnd) {
			return true
		}
	}
}
