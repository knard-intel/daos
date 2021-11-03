//
// (C) Copyright 2021 Intel Corporation.
//
// SPDX-License-Identifier: BSD-2-Clause-Patent
//

package hwloc

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/daos-stack/daos/src/control/common"
	"github.com/daos-stack/daos/src/control/lib/hardware"
	"github.com/daos-stack/daos/src/control/logging"
)

func TestHwlocProvider_GetTopology_Samples(t *testing.T) {
	// sample hwloc topologies taken from real systems
	for name, tc := range map[string]struct {
		hwlocXMLFile string
		expResult    *hardware.Topology
	}{
		"boro-84": {
			hwlocXMLFile: "testdata/boro-84.xml",
			expResult: &hardware.Topology{
				NUMANodes: map[uint]*hardware.NUMANode{
					0: {
						ID:       0,
						NumCores: 24,
						Devices: hardware.PCIDevices{
							"0000:18:00.0": {
								{
									Name:    "ib0",
									Type:    hardware.DeviceTypeNetwork,
									PCIAddr: "0000:18:00.0",
								},
								{
									Name:    "hfi1_0",
									Type:    hardware.DeviceTypeOpenFabrics,
									PCIAddr: "0000:18:00.0",
								},
							},
							"0000:3d:00.0": {
								{
									Name:    "eth0",
									Type:    hardware.DeviceTypeNetwork,
									PCIAddr: "0000:3d:00.0",
								},
								{
									Name:    "i40iw1",
									Type:    hardware.DeviceTypeOpenFabrics,
									PCIAddr: "0000:3d:00.0",
								},
							},
							"0000:3d:00.1": {
								{
									Name:    "eth1",
									Type:    hardware.DeviceTypeNetwork,
									PCIAddr: "0000:3d:00.1",
								},
								{
									Name:    "i40iw0",
									Type:    hardware.DeviceTypeOpenFabrics,
									PCIAddr: "0000:3d:00.1",
								},
							},
						},
					},
					1: {
						ID:       1,
						NumCores: 24,
						Devices:  hardware.PCIDevices{},
					},
				},
			},
		},
		"wolf-133": {
			hwlocXMLFile: "testdata/wolf-133.xml",
			expResult: &hardware.Topology{
				NUMANodes: map[uint]*hardware.NUMANode{
					0: {
						ID:       0,
						NumCores: 24,
						Devices: hardware.PCIDevices{
							"0000:18:00.0": {
								{
									Name:    "ib0",
									Type:    hardware.DeviceTypeNetwork,
									PCIAddr: "0000:18:00.0",
								},
								{
									Name:    "hfi1_0",
									Type:    hardware.DeviceTypeOpenFabrics,
									PCIAddr: "0000:18:00.0",
								},
							},
							"0000:3d:00.0": {
								{
									Name:    "eth0",
									Type:    hardware.DeviceTypeNetwork,
									PCIAddr: "0000:3d:00.0",
								},
								{
									Name:    "i40iw1",
									Type:    hardware.DeviceTypeOpenFabrics,
									PCIAddr: "0000:3d:00.0",
								},
							},
							"0000:3d:00.1": {
								{
									Name:    "eth1",
									Type:    hardware.DeviceTypeNetwork,
									PCIAddr: "0000:3d:00.1",
								},
								{
									Name:    "i40iw0",
									Type:    hardware.DeviceTypeOpenFabrics,
									PCIAddr: "0000:3d:00.1",
								},
							},
						},
					},
					1: {
						ID:       1,
						NumCores: 24,
						Devices: hardware.PCIDevices{
							"0000:af:00.0": {
								{
									Name:    "ib1",
									Type:    hardware.DeviceTypeNetwork,
									PCIAddr: "0000:af:00.0",
								},
								{
									Name:    "hfi1_1",
									Type:    hardware.DeviceTypeOpenFabrics,
									PCIAddr: "0000:af:00.0",
								},
							},
						},
					},
				},
			},
		},
		"no devices": {
			hwlocXMLFile: "testdata/gcp_topology.xml",
			expResult: &hardware.Topology{
				NUMANodes: map[uint]*hardware.NUMANode{
					0: {
						ID:       0,
						NumCores: 8,
						Devices:  hardware.PCIDevices{},
					},
					1: {
						ID:       1,
						NumCores: 8,
						Devices:  hardware.PCIDevices{},
					},
				},
			},
		},
		"multiport": {
			hwlocXMLFile: "testdata/multiport_hfi_topology.xml",
			expResult: &hardware.Topology{
				NUMANodes: map[uint]*hardware.NUMANode{
					0: {
						ID:       0,
						NumCores: 8,
						Devices: hardware.PCIDevices{
							"0000:02:00.0": {
								{
									Name:    "ib0",
									Type:    hardware.DeviceTypeNetwork,
									PCIAddr: "0000:02:00.0",
								},
								{
									Name:    "enp2s0",
									Type:    hardware.DeviceTypeNetwork,
									PCIAddr: "0000:02:00.0",
								},
								{
									Name:    "mlx4_0",
									Type:    hardware.DeviceTypeOpenFabrics,
									PCIAddr: "0000:02:00.0",
								},
							},
							"0000:06:00.0": {
								{
									Name:    "enp6s0",
									Type:    hardware.DeviceTypeNetwork,
									PCIAddr: "0000:06:00.0",
								},
							},
						},
					},
					1: {
						ID:       1,
						NumCores: 8,
						Devices: hardware.PCIDevices{
							"0000:83:00.0": {
								{
									Name:    "ib1",
									Type:    hardware.DeviceTypeNetwork,
									PCIAddr: "0000:83:00.0",
								},
								{
									Name:    "ib2",
									Type:    hardware.DeviceTypeNetwork,
									PCIAddr: "0000:83:00.0",
								},
								{
									Name:    "mlx4_1",
									Type:    hardware.DeviceTypeOpenFabrics,
									PCIAddr: "0000:83:00.0",
								},
							},
						},
					},
				},
			},
		},
		"no NUMA nodes in topology": {
			hwlocXMLFile: "testdata/no-numa-nodes.xml",
			expResult: &hardware.Topology{
				NUMANodes: map[uint]*hardware.NUMANode{
					0: {
						ID:       0,
						NumCores: 4,
						Devices: hardware.PCIDevices{
							"0000:18:00.0": {
								{
									Name:    "ib0",
									Type:    hardware.DeviceTypeNetwork,
									PCIAddr: "0000:18:00.0",
								},
								{
									Name:    "hfi1_0",
									Type:    hardware.DeviceTypeOpenFabrics,
									PCIAddr: "0000:18:00.0",
								},
							},
						},
					},
				},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			log, buf := logging.NewTestLogger(name)
			defer common.ShowBufferOnFailure(t, buf)

			_, err := os.Stat(tc.hwlocXMLFile)
			common.AssertEqual(t, err, nil, "unable to read hwloc XML file")
			os.Setenv("HWLOC_XMLFILE", tc.hwlocXMLFile)
			defer os.Unsetenv("HWLOC_XMLFILE")

			provider := NewProvider(log)
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			result, err := provider.GetTopology(ctx)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(tc.expResult, result); diff != "" {
				t.Errorf("(-want, +got)\n%s\n", diff)
			}
		})

	}
}
