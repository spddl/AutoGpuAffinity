package main

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

var (
	// Library
	modSetupapi = windows.NewLazyDLL("setupapi.dll")

	// Functions
	procSetupDiGetClassDevsW = modSetupapi.NewProc("SetupDiGetClassDevsW")
)

type Device struct {
	Idata DevInfoData
	reg   registry.Key

	// AffinityPolicy
	DevicePolicy          uint32
	AssignmentSetOverride Bits
}

func FindAllDevices() ([]Device, DevInfo) {
	var allDevices []Device
	handle, err := SetupDiGetClassDevs(nil, nil, 0, uint32(DIGCF_ALLCLASSES|DIGCF_PRESENT|DIGCF_DEVICEINTERFACE))
	if err != nil {
		panic(err)
	}

	var index = 0
	for {
		idata, err := SetupDiEnumDeviceInfo(handle, index)
		if err != nil { // ERROR_NO_MORE_ITEMS
			break
		}
		index++

		dev := Device{
			Idata: *idata,
		}

		val, err := SetupDiGetDeviceRegistryProperty(handle, idata, SPDRP_CLASSGUID)
		if err == nil {
			if val.(string) != "{4d36e968-e325-11ce-bfc1-08002be10318}" {
				// Display Adapters
				// Class = Display
				// ClassGuid = {4d36e968-e325-11ce-bfc1-08002be10318}
				// This class includes video adapters. Drivers for this class include display drivers and video miniport drivers.
				// https://docs.microsoft.com/en-us/windows-hardware/drivers/install/system-defined-device-setup-classes-available-to-vendors
				continue
			}
		} else {
			continue
		}

		dev.reg, _ = SetupDiOpenDevRegKey(handle, idata, DICS_FLAG_GLOBAL, 0, DIREG_DEV, windows.KEY_READ)

		affinityPolicyKey, _ := registry.OpenKey(dev.reg, `Interrupt Management\Affinity Policy`, registry.QUERY_VALUE)
		dev.DevicePolicy = GetDWORDuint32Value(affinityPolicyKey, "DevicePolicy")               // REG_DWORD
		AssignmentSetOverrideByte := GetBinaryValue(affinityPolicyKey, "AssignmentSetOverride") // REG_BINARY
		affinityPolicyKey.Close()

		if len(AssignmentSetOverrideByte) != 0 {
			AssignmentSetOverrideBytes := make([]byte, 8)
			copy(AssignmentSetOverrideBytes, AssignmentSetOverrideByte)
			dev.AssignmentSetOverride = Bits(btoi64(AssignmentSetOverrideBytes))
		}

		allDevices = append(allDevices, dev)
	}
	return allDevices, handle
}

func SetupDiGetClassDevs(classGuid *windows.GUID, enumerator *uint16, hwndParent uintptr, flags uint32) (handle DevInfo, err error) {
	r0, _, e1 := syscall.SyscallN(procSetupDiGetClassDevsW.Addr(), uintptr(unsafe.Pointer(classGuid)), uintptr(unsafe.Pointer(enumerator)), uintptr(hwndParent), uintptr(flags))
	handle = DevInfo(r0)
	if handle == DevInfo(windows.InvalidHandle) {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func GetDeviceProperty(dis DevInfo, devInfoData *DevInfoData, devPropKey DEVPROPKEY) ([]byte, error) {
	var propt, size uint32
	buf := make([]byte, 16)
	run := true
	for run {
		err := SetupDiGetDeviceProperty(dis, devInfoData, &devPropKey, &propt, &buf[0], uint32(len(buf)), &size, 0)
		switch {
		case size > uint32(len(buf)):
			buf = make([]byte, size+16)
		case err != nil:
			return buf, err
		default:
			run = false
		}
	}

	return buf, nil
}
