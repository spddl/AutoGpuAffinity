package main

import (
	"fmt"
	"log"
	"syscall"
	"time"
	"unsafe"
)

var (
	// Library
	user32 = syscall.MustLoadDLL("user32.dll")

	// Functions
	procEnumWindows         = user32.MustFindProc("EnumWindows")         // https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-enumwindows
	procGetWindowTextW      = user32.MustFindProc("GetWindowTextW")      // https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getwindowtextw
	procSetForegroundWindow = user32.MustFindProc("SetForegroundWindow") // https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-setforegroundwindow
	procGetForegroundWindow = user32.MustFindProc("GetForegroundWindow") // https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getforegroundwindow
	procDestroyWindow       = user32.MustFindProc("DestroyWindow")       // https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-destroywindow
	procOpenInputDesktop    = user32.MustFindProc("OpenInputDesktop")    // https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-openinputdesktop
	procCloseDesktop        = user32.MustFindProc("CloseDesktop")        // https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-closedesktop
	procIsIconic            = user32.MustFindProc("IsIconic")            // https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-isiconic
	// procShowWindow          = user32.MustFindProc("ShowWindow")          // https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-showwindow
	procIsWindowVisible = user32.MustFindProc("IsWindowVisible") // https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-iswindowvisible
)

const PROCESS_QUERY_LIMITED_INFORMATION = 0x0400

func EnumWindows(enumFunc uintptr, lparam uintptr) (err error) {
	r1, _, e1 := syscall.SyscallN(procEnumWindows.Addr(), uintptr(enumFunc), uintptr(lparam))
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func GetWindowText(hwnd uintptr, str *uint16, maxCount int32) (len int32, err error) {
	r0, _, e1 := syscall.SyscallN(procGetWindowTextW.Addr(), hwnd, uintptr(unsafe.Pointer(str)), uintptr(maxCount))
	len = int32(r0)
	if len == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func FindWindow(title string) (uintptr, error) {
	var hwnd uintptr
	cb := syscall.NewCallback(func(h uintptr, p uintptr) uintptr {
		b := make([]uint16, 200)
		_, err := GetWindowText(h, &b[0], int32(len(b)))
		if err != nil {
			// ignore the error
			return 1 // continue enumeration
		}
		// log.Println("syscall.UTF16ToString(b)", syscall.UTF16ToString(b))
		if syscall.UTF16ToString(b) == title {
			// note the window
			hwnd = uintptr(h)
			return 0 // stop enumeration
		}
		return 1 // continue enumeration
	})
	EnumWindows(cb, 0)
	if hwnd == 0 {
		return 0, fmt.Errorf("no window with title '%s' found", title)
	}
	return hwnd, nil
}

func pid_is_running(pid uint32) (bool, syscall.Handle) {
	handle, err := syscall.OpenProcess(PROCESS_QUERY_LIMITED_INFORMATION, false, pid)
	defer syscall.CloseHandle(handle)
	if err != nil {
		log.Println(err)
		return false, handle
	}
	if handle == 0 {
		return false, handle
	}
	var ec uint32
	err = syscall.GetExitCodeProcess(handle, &ec) // https://docs.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-getexitcodeprocess
	if ec == 259 {                                // STILL_ACTIVE
		return true, handle
	}
	if err != nil {
		log.Println(err)
		return false, handle
	}
	return true, handle
}

func SetTopWindow(hwnd uintptr) bool {
	for {
		time.Sleep(400 * time.Millisecond)

		ret, _, _ := procIsWindowVisible.Call(hwnd)
		if ret != 1 { // window is not visible
			return false
		}

		ret, _, _ = procIsIconic.Call(hwnd)
		if ret == 1 { // IsIconic
			return false
			// UIPI "User Interface Privilege Isolation"
			// ret, _, err := syscall.Syscall(procShowWindow.Addr(), 2, hwnd, uintptr(SW_RESTORE), 0)
			// log.Println("procShowWindow\t\t", hwnd, (ret == 1), ret, err)
		}

		ret, _, _ = procGetForegroundWindow.Call()
		if hwnd == ret { // hwnd == foregroundwindow
			return true
		}

		procSetForegroundWindow.Call(hwnd)
	}
}

// DestroyWindow destroys the specified window.
func DestroyWindow(hwnd uintptr) error {
	ret, _, err := procDestroyWindow.Call(hwnd)
	if ret == 0 {
		return err
	}
	return nil
}

func openInputDesktop(flags uint32, isInherit bool, accessMask uint32) uintptr {
	var isInheritValue uint32
	if isInherit {
		isInheritValue = 1
	}
	r, _, _ := procOpenInputDesktop.Call(
		uintptr(flags),
		uintptr(isInheritValue),
		uintptr(accessMask))
	return uintptr(r)
}

func closeDesktop(hDesk uintptr) bool {
	r, _, _ := procCloseDesktop.Call(hDesk)
	if r == 0 {
		return false
	} else {
		return true
	}
}
