// Code generated by 'go generate'; DO NOT EDIT.

package main

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var _ unsafe.Pointer

// Do the interface allocations only once for common
// Errno values.
const (
	errnoERROR_IO_PENDING = 997
)

var (
	errERROR_IO_PENDING error = syscall.Errno(errnoERROR_IO_PENDING)
)

// errnoErr returns common boxed Errno values, to prevent
// allocations at runtime.
func errnoErr(e syscall.Errno) error {
	switch e {
	case 0:
		return nil
	case errnoERROR_IO_PENDING:
		return errERROR_IO_PENDING
	}
	// TODO: add more here, after collecting data on the common
	// error values see on Windows. (perhaps when running
	// all.bat?)
	return e
}

var (
	modsetupapi = windows.NewLazySystemDLL("setupapi.dll")

	procSetupDiCreateDeviceInfoListExW    = modsetupapi.NewProc("SetupDiCreateDeviceInfoListExW")
	procSetupDiGetDeviceInfoListDetailW   = modsetupapi.NewProc("SetupDiGetDeviceInfoListDetailW")
	procSetupDiCreateDeviceInfoW          = modsetupapi.NewProc("SetupDiCreateDeviceInfoW")
	procSetupDiEnumDeviceInfo             = modsetupapi.NewProc("SetupDiEnumDeviceInfo")
	procSetupDiDestroyDeviceInfoList      = modsetupapi.NewProc("SetupDiDestroyDeviceInfoList")
	procSetupDiBuildDriverInfoList        = modsetupapi.NewProc("SetupDiBuildDriverInfoList")
	procSetupDiCancelDriverInfoSearch     = modsetupapi.NewProc("SetupDiCancelDriverInfoSearch")
	procSetupDiEnumDriverInfoW            = modsetupapi.NewProc("SetupDiEnumDriverInfoW")
	procSetupDiGetSelectedDriverW         = modsetupapi.NewProc("SetupDiGetSelectedDriverW")
	procSetupDiSetSelectedDriverW         = modsetupapi.NewProc("SetupDiSetSelectedDriverW")
	procSetupDiGetDriverInfoDetailW       = modsetupapi.NewProc("SetupDiGetDriverInfoDetailW")
	procSetupDiDestroyDriverInfoList      = modsetupapi.NewProc("SetupDiDestroyDriverInfoList")
	procSetupDiGetClassDevsExW            = modsetupapi.NewProc("SetupDiGetClassDevsExW")
	procSetupDiCallClassInstaller         = modsetupapi.NewProc("SetupDiCallClassInstaller")
	procSetupDiOpenDevRegKey              = modsetupapi.NewProc("SetupDiOpenDevRegKey")
	procSetupDiGetDeviceRegistryPropertyW = modsetupapi.NewProc("SetupDiGetDeviceRegistryPropertyW")
	procSetupDiSetDeviceRegistryPropertyW = modsetupapi.NewProc("SetupDiSetDeviceRegistryPropertyW")
	procSetupDiGetDeviceInstallParamsW    = modsetupapi.NewProc("SetupDiGetDeviceInstallParamsW")
	procSetupDiGetClassInstallParamsW     = modsetupapi.NewProc("SetupDiGetClassInstallParamsW")
	procSetupDiSetDeviceInstallParamsW    = modsetupapi.NewProc("SetupDiSetDeviceInstallParamsW")
	procSetupDiSetClassInstallParamsW     = modsetupapi.NewProc("SetupDiSetClassInstallParamsW")
	procSetupDiClassNameFromGuidExW       = modsetupapi.NewProc("SetupDiClassNameFromGuidExW")
	procSetupDiClassGuidsFromNameExW      = modsetupapi.NewProc("SetupDiClassGuidsFromNameExW")
	procSetupDiGetSelectedDevice          = modsetupapi.NewProc("SetupDiGetSelectedDevice")
	procSetupDiSetSelectedDevice          = modsetupapi.NewProc("SetupDiSetSelectedDevice")
	procSetupDiGetDevicePropertyW         = modsetupapi.NewProc("SetupDiGetDevicePropertyW")
	procSetupDiRestartDevices             = modsetupapi.NewProc("SetupDiRestartDevices")

	modCfgMgr32 = windows.NewLazySystemDLL("CfgMgr32.dll")
	procCM_Get_DevNode_Status = modCfgMgr32.NewProc("CM_Get_DevNode_Status")
)

func setupDiCreateDeviceInfoListEx(classGUID *windows.GUID, hwndParent uintptr, machineName *uint16, reserved uintptr) (handle DevInfo, err error) {
	r0, _, e1 := syscall.SyscallN(procSetupDiCreateDeviceInfoListExW.Addr(), uintptr(unsafe.Pointer(classGUID)), uintptr(hwndParent), uintptr(unsafe.Pointer(machineName)), uintptr(reserved))
	handle = DevInfo(r0)
	if handle == DevInfo(windows.InvalidHandle) {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func setupDiGetDeviceInfoListDetail(deviceInfoSet DevInfo, deviceInfoSetDetailData *DevInfoListDetailData) (err error) {
	r1, _, e1 := syscall.SyscallN(procSetupDiGetDeviceInfoListDetailW.Addr(), uintptr(deviceInfoSet), uintptr(unsafe.Pointer(deviceInfoSetDetailData)))
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func setupDiCreateDeviceInfo(deviceInfoSet DevInfo, DeviceName *uint16, classGUID *windows.GUID, DeviceDescription *uint16, hwndParent uintptr, CreationFlags DICD, deviceInfoData *DevInfoData) (err error) {
	r1, _, e1 := syscall.SyscallN(procSetupDiCreateDeviceInfoW.Addr(), uintptr(deviceInfoSet), uintptr(unsafe.Pointer(DeviceName)), uintptr(unsafe.Pointer(classGUID)), uintptr(unsafe.Pointer(DeviceDescription)), uintptr(hwndParent), uintptr(CreationFlags), uintptr(unsafe.Pointer(deviceInfoData)))
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func setupDiEnumDeviceInfo(deviceInfoSet DevInfo, memberIndex uint32, deviceInfoData *DevInfoData) (err error) {
	r1, _, e1 := syscall.SyscallN(procSetupDiEnumDeviceInfo.Addr(), uintptr(deviceInfoSet), uintptr(memberIndex), uintptr(unsafe.Pointer(deviceInfoData)))
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func SetupDiDestroyDeviceInfoList(deviceInfoSet DevInfo) (err error) {
	r1, _, e1 := syscall.SyscallN(procSetupDiDestroyDeviceInfoList.Addr(), uintptr(deviceInfoSet))
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func SetupDiBuildDriverInfoList(deviceInfoSet DevInfo, deviceInfoData *DevInfoData, driverType SPDIT) (err error) {
	r1, _, e1 := syscall.SyscallN(procSetupDiBuildDriverInfoList.Addr(), uintptr(deviceInfoSet), uintptr(unsafe.Pointer(deviceInfoData)), uintptr(driverType))
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func SetupDiCancelDriverInfoSearch(deviceInfoSet DevInfo) (err error) {
	r1, _, e1 := syscall.SyscallN(procSetupDiCancelDriverInfoSearch.Addr(), uintptr(deviceInfoSet))
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func setupDiEnumDriverInfo(deviceInfoSet DevInfo, deviceInfoData *DevInfoData, driverType SPDIT, memberIndex uint32, driverInfoData *DrvInfoData) (err error) {
	r1, _, e1 := syscall.SyscallN(procSetupDiEnumDriverInfoW.Addr(), uintptr(deviceInfoSet), uintptr(unsafe.Pointer(deviceInfoData)), uintptr(driverType), uintptr(memberIndex), uintptr(unsafe.Pointer(driverInfoData)))
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func setupDiGetSelectedDriver(deviceInfoSet DevInfo, deviceInfoData *DevInfoData, driverInfoData *DrvInfoData) (err error) {
	r1, _, e1 := syscall.SyscallN(procSetupDiGetSelectedDriverW.Addr(), uintptr(deviceInfoSet), uintptr(unsafe.Pointer(deviceInfoData)), uintptr(unsafe.Pointer(driverInfoData)))
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func SetupDiSetSelectedDriver(deviceInfoSet DevInfo, deviceInfoData *DevInfoData, driverInfoData *DrvInfoData) (err error) {
	r1, _, e1 := syscall.SyscallN(procSetupDiSetSelectedDriverW.Addr(), uintptr(deviceInfoSet), uintptr(unsafe.Pointer(deviceInfoData)), uintptr(unsafe.Pointer(driverInfoData)))
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func setupDiGetDriverInfoDetail(deviceInfoSet DevInfo, deviceInfoData *DevInfoData, driverInfoData *DrvInfoData, driverInfoDetailData *DrvInfoDetailData, driverInfoDetailDataSize uint32, requiredSize *uint32) (err error) {
	r1, _, e1 := syscall.SyscallN(procSetupDiGetDriverInfoDetailW.Addr(), uintptr(deviceInfoSet), uintptr(unsafe.Pointer(deviceInfoData)), uintptr(unsafe.Pointer(driverInfoData)), uintptr(unsafe.Pointer(driverInfoDetailData)), uintptr(driverInfoDetailDataSize), uintptr(unsafe.Pointer(requiredSize)))
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func SetupDiDestroyDriverInfoList(deviceInfoSet DevInfo, deviceInfoData *DevInfoData, driverType SPDIT) (err error) {
	r1, _, e1 := syscall.SyscallN(procSetupDiDestroyDriverInfoList.Addr(), uintptr(deviceInfoSet), uintptr(unsafe.Pointer(deviceInfoData)), uintptr(driverType))
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func setupDiGetClassDevsEx(classGUID *windows.GUID, Enumerator *uint16, hwndParent uintptr, Flags DIGCF, deviceInfoSet DevInfo, machineName *uint16, reserved uintptr) (handle DevInfo, err error) {
	r0, _, e1 := syscall.SyscallN(procSetupDiGetClassDevsExW.Addr(), uintptr(unsafe.Pointer(classGUID)), uintptr(unsafe.Pointer(Enumerator)), uintptr(hwndParent), uintptr(Flags), uintptr(deviceInfoSet), uintptr(unsafe.Pointer(machineName)), uintptr(reserved))
	handle = DevInfo(r0)
	if handle == DevInfo(windows.InvalidHandle) {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func SetupDiCallClassInstaller(installFunction DI_FUNCTION, deviceInfoSet DevInfo, deviceInfoData *DevInfoData) (err error) {
	r1, _, e1 := syscall.SyscallN(
		procSetupDiCallClassInstaller.Addr(),
		uintptr(installFunction),
		uintptr(deviceInfoSet),
		uintptr(unsafe.Pointer(deviceInfoData)),
	)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func setupDiOpenDevRegKey(deviceInfoSet DevInfo, deviceInfoData *DevInfoData, Scope DICS_FLAG, HwProfile uint32, KeyType DIREG, samDesired uint32) (key windows.Handle, err error) {
	r0, _, e1 := syscall.SyscallN(procSetupDiOpenDevRegKey.Addr(), uintptr(deviceInfoSet), uintptr(unsafe.Pointer(deviceInfoData)), uintptr(Scope), uintptr(HwProfile), uintptr(KeyType), uintptr(samDesired))
	key = windows.Handle(r0)
	if key == windows.InvalidHandle {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func setupDiGetDeviceRegistryProperty(deviceInfoSet DevInfo, deviceInfoData *DevInfoData, property SPDRP, propertyRegDataType *uint32, propertyBuffer *byte, propertyBufferSize uint32, requiredSize *uint32) (err error) {
	r1, _, e1 := syscall.SyscallN(procSetupDiGetDeviceRegistryPropertyW.Addr(), uintptr(deviceInfoSet), uintptr(unsafe.Pointer(deviceInfoData)), uintptr(property), uintptr(unsafe.Pointer(propertyRegDataType)), uintptr(unsafe.Pointer(propertyBuffer)), uintptr(propertyBufferSize), uintptr(unsafe.Pointer(requiredSize)))
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func setupDiSetDeviceRegistryProperty(deviceInfoSet DevInfo, deviceInfoData *DevInfoData, property SPDRP, propertyBuffer *byte, propertyBufferSize uint32) (err error) {
	r1, _, e1 := syscall.SyscallN(procSetupDiSetDeviceRegistryPropertyW.Addr(), uintptr(deviceInfoSet), uintptr(unsafe.Pointer(deviceInfoData)), uintptr(property), uintptr(unsafe.Pointer(propertyBuffer)), uintptr(propertyBufferSize))
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func setupDiGetDeviceInstallParams(deviceInfoSet DevInfo, deviceInfoData *DevInfoData, deviceInstallParams *DevInstallParams) (err error) {
	r1, _, e1 := syscall.SyscallN(procSetupDiGetDeviceInstallParamsW.Addr(), uintptr(deviceInfoSet), uintptr(unsafe.Pointer(deviceInfoData)), uintptr(unsafe.Pointer(deviceInstallParams)))
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func SetupDiGetClassInstallParams(deviceInfoSet DevInfo, deviceInfoData *DevInfoData, classInstallParams *ClassInstallHeader, classInstallParamsSize uint32, requiredSize *uint32) (err error) {
	r1, _, e1 := syscall.SyscallN(procSetupDiGetClassInstallParamsW.Addr(), uintptr(deviceInfoSet), uintptr(unsafe.Pointer(deviceInfoData)), uintptr(unsafe.Pointer(classInstallParams)), uintptr(classInstallParamsSize), uintptr(unsafe.Pointer(requiredSize)))
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func SetupDiSetDeviceInstallParams(deviceInfoSet DevInfo, deviceInfoData *DevInfoData, deviceInstallParams *DevInstallParams) (err error) {
	r1, _, e1 := syscall.SyscallN(procSetupDiSetDeviceInstallParamsW.Addr(), uintptr(deviceInfoSet), uintptr(unsafe.Pointer(deviceInfoData)), uintptr(unsafe.Pointer(deviceInstallParams)))
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

// WINSETUPAPI BOOL SetupDiRestartDevices(
//   [in]      HDEVINFO         DeviceInfoSet,
//   [in, out] PSP_DEVINFO_DATA DeviceInfoData
// );
func SetupDiRestartDevices(deviceInfoSet DevInfo, deviceInfoData *DevInfoData) (err error) {
	r1, _, e1 := syscall.SyscallN(
		procSetupDiRestartDevices.Addr(),
		uintptr(deviceInfoSet),
		uintptr(unsafe.Pointer(deviceInfoData)),
	)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func SetupDiSetClassInstallParams(deviceInfoSet DevInfo, deviceInfoData *DevInfoData, classInstallParams *ClassInstallHeader, classInstallParamsSize uint32) (err error) {
	r1, _, e1 := syscall.SyscallN(procSetupDiSetClassInstallParamsW.Addr(),
		uintptr(deviceInfoSet),
		uintptr(unsafe.Pointer(deviceInfoData)),
		uintptr(unsafe.Pointer(classInstallParams)),
		uintptr(classInstallParamsSize),
	)
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func setupDiClassNameFromGuidEx(classGUID *windows.GUID, className *uint16, classNameSize uint32, requiredSize *uint32, machineName *uint16, reserved uintptr) (err error) {
	r1, _, e1 := syscall.SyscallN(procSetupDiClassNameFromGuidExW.Addr(), uintptr(unsafe.Pointer(classGUID)), uintptr(unsafe.Pointer(className)), uintptr(classNameSize), uintptr(unsafe.Pointer(requiredSize)), uintptr(unsafe.Pointer(machineName)), uintptr(reserved))
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func setupDiClassGuidsFromNameEx(className *uint16, classGuidList *windows.GUID, classGuidListSize uint32, requiredSize *uint32, machineName *uint16, reserved uintptr) (err error) {
	r1, _, e1 := syscall.SyscallN(procSetupDiClassGuidsFromNameExW.Addr(), uintptr(unsafe.Pointer(className)), uintptr(unsafe.Pointer(classGuidList)), uintptr(classGuidListSize), uintptr(unsafe.Pointer(requiredSize)), uintptr(unsafe.Pointer(machineName)), uintptr(reserved))
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func setupDiGetSelectedDevice(deviceInfoSet DevInfo, deviceInfoData *DevInfoData) (err error) {
	r1, _, e1 := syscall.SyscallN(procSetupDiGetSelectedDevice.Addr(), uintptr(deviceInfoSet), uintptr(unsafe.Pointer(deviceInfoData)))
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func SetupDiSetSelectedDevice(deviceInfoSet DevInfo, deviceInfoData *DevInfoData) (err error) {
	r1, _, e1 := syscall.SyscallN(procSetupDiSetSelectedDevice.Addr(), uintptr(deviceInfoSet), uintptr(unsafe.Pointer(deviceInfoData)))
	if r1 == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

// https://github.com/tajtiattila/hid/blob/2bd63ffd4c8c0b81e5999c7b183cc4325f773527/platform/zsys_windows.go
func SetupDiGetDeviceProperty(devInfoSet DevInfo, devInfoData *DevInfoData, propKey *DEVPROPKEY, propType *uint32, propBuf *byte, propBufSize uint32, reqsize *uint32, flags uint32) (err error) {
	r1, _, e1 := syscall.SyscallN(procSetupDiGetDevicePropertyW.Addr(), uintptr(devInfoSet), uintptr(unsafe.Pointer(devInfoData)), uintptr(unsafe.Pointer(propKey)), uintptr(unsafe.Pointer(propType)), uintptr(unsafe.Pointer(propBuf)), uintptr(propBufSize), uintptr(unsafe.Pointer(reqsize)), uintptr(flags))
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func cm_Get_DevNode_Status(status *uint32, problemNumber *uint32, devInst uint32, flags uint32) uint32 {
	r0, r1, er := syscall.SyscallN(procCM_Get_DevNode_Status.Addr(), uintptr(unsafe.Pointer(status)), uintptr(unsafe.Pointer(problemNumber)), uintptr(devInst), uintptr(flags))
	println(r1)
	println(er)
	return uint32(r0)
}