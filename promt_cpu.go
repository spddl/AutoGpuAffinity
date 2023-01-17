package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"unsafe"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var quitTextStyle = lipgloss.NewStyle().Margin(1, 0, 2, 4)
var SpecifiedProc uint32 = 4

type cpumodel struct {
	quitting bool
	cursor   int
	choices  []string
	selected map[int]struct{}
}

func initialCPUModel() cpumodel {
	var choices []string
	for _, v := range CPUArray {
		if HT {
			i, _ := strconv.Atoi(v)
			if i%2 == 0 {
				choices = append(choices, "CPU "+v)
			} else {
				choices = append(choices, fmt.Sprintf("Thread %d", i-1))
			}
		} else {
			choices = append(choices, "CPU"+v)
		}
	}

	selected := make(map[int]struct{})
	if len(defaultSettings.GPUdevices) == 1 { // nur die erste GPU wird für die Vorauswahl benutzt
		if defaultSettings.GPUdevices[0].DevicePolicy == SpecifiedProc {
			for bit, cpu := range CPUMap {
				if Has(bit, defaultSettings.GPUdevices[0].AssignmentSetOverride) {
					i, _ := strconv.Atoi(cpu)
					selected[i] = struct{}{}
				}
			}
		}
	}

	return cpumodel{
		choices:  choices,
		selected: selected,
	}
}

func (m cpumodel) Init() tea.Cmd {
	return nil
}

func (m cpumodel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "ctrl+c", "esc":
			os.Exit(1)
			return m, tea.Quit
		case "q":
			m.quitting = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	return m, nil
}

func (m cpumodel) View() string {
	if m.quitting {
		var setBits Bits
		if len(m.selected) != 0 {
			for i := range m.selected {
				setBits |= CPUBits[i]
			}
			SetGPUandRestart(setBits)
			return quitTextStyle.Render("Done.")
		}

		for i := range GPUdevices {
			propChangeParams := PropChangeParams{
				ClassInstallHeader: *MakeClassInstallHeader(DIF_PROPERTYCHANGE),
				StateChange:        DICS_PROPCHANGE,
				Scope:              DICS_FLAG_GLOBAL,
			}
			for c := 0; c < 2; c++ {
				if err := SetupDiSetClassInstallParams(handle, &GPUdevices[i].Idata, &propChangeParams.ClassInstallHeader, uint32(unsafe.Sizeof(propChangeParams))); err != nil {
					log.Println(err)
				}
				if err := SetupDiCallClassInstaller(DIF_PROPERTYCHANGE, handle, &GPUdevices[i].Idata); err != nil {
					log.Println(err)
				}
			}
		}

		return quitTextStyle.Render("set old settings and restarts driver")
	}

	s := fmt.Sprintf("\nThe separate logs are located here: %s\nYou can upload and view them here https://boringboredom.github.io/Frame-Time-Analysis/ for better overview,\nthey will be deleted on exit.\n\nOn which CPU should the GPU be assigned?\n\n", tempFolder)

	changeTest := false
	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
			changeTest = true
		}

		switch CPUPoints[i] {
		case topCPU1:
			s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, defaultStyle.Foreground(highest[2]).Render(choice))
		case topCPU2:
			s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, defaultStyle.Foreground(highest[1]).Render(choice))
		case topCPU3:
			s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, defaultStyle.Foreground(highest[0]).Render(choice))
		default:
			s += defaultStyle.Render(fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice))
		}

	}

	if changeTest {
		s += "\nPress q to apply. (esc to quit)\n"
	} else {
		s += "\nPress q to restart driver. (esc to quit)\n"
	}

	return s
}

func Promt_cpu() {
	p := tea.NewProgram(initialCPUModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
