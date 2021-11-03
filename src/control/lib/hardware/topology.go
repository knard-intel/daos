//
// (C) Copyright 2021 Intel Corporation.
//
// SPDX-License-Identifier: BSD-2-Clause-Patent
//

package hardware

import "context"

type (
	// Topology is a hierarchy of hardware devices grouped under NUMA nodes.
	Topology struct {
		// NUMANodes is the set of NUMA nodes mapped by their ID.
		NUMANodes map[uint]*NUMANode
	}

	// NUMANode represents an individual NUMA node in the system and the devices associated with it.
	NUMANode struct {
		ID       uint
		NumCores uint
		Devices  PCIDevices
	}

	// PCIDevices groups hardware devices by PCI address.
	PCIDevices map[string][]*Device
)

// Add adds a device to the PCIDevices.
func (m PCIDevices) Add(dev *Device) {
	if m == nil || dev == nil {
		return
	}
	addr := dev.PCIAddr
	m[addr] = append(m[addr], dev)
}

// DeviceType indicates the type of a hardware device.
type DeviceType uint

const (
	// DeviceTypeUnknown indicates a device type that is not recognized.
	DeviceTypeUnknown DeviceType = iota
	// DeviceTypeNetwork indicates a standard network device.
	DeviceTypeNetwork
	// DeviceTypeOpenFabrics indicates an OpenFabrics device.
	DeviceTypeOpenFabrics
)

func (t DeviceType) String() string {
	switch t {
	case DeviceTypeNetwork:
		return "Network"
	case DeviceTypeOpenFabrics:
		return "OpenFabrics"
	}

	return "Unknown"
}

// Device represents an individual hardware device.
type Device struct {
	Name    string
	Type    DeviceType
	PCIAddr string
}

// TopologyProvider is an interface for acquiring a system topology.
type TopologyProvider interface {
	GetTopology(context.Context) (*Topology, error)
}
