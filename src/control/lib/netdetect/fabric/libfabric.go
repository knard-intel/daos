//
// (C) Copyright 2021 Intel Corporation.
//
// SPDX-License-Identifier: BSD-2-Clause-Patent
//

package fabric

/*
#cgo LDFLAGS: -lfabric
#include <rdma/fabric.h>
#include <rdma/fi_domain.h>
#include <rdma/fi_endpoint.h>
#include <rdma/fi_cm.h>
#include <rdma/fi_tagged.h>
#include <rdma/fi_rma.h>
#include <rdma/fi_errno.h>

typedef struct {
	uint64_t reserved_1;
	uint8_t  reserved_2;
	int8_t   unit;
	uint8_t  port;
	uint8_t  reserved_3;
	uint32_t service;
} psmx2_ep_name;

int getHFIUnit(void *src_addr) {
	psmx2_ep_name *psmx2;
	psmx2 = (psmx2_ep_name *)src_addr;
	return psmx2->unit;
}
*/
import "C"

import 	"github.com/pkg/errors"

const (
	libFabricMajorVersion = 1
	libFabricMinorVersion = 7
	allHFIUsed            = -1
)

type FIInfo struct {
	info *C.struct_fi_info
}

func (f *FIInfo) GetHFIUnit() (int, error) {
	if f.info.src_addr == nil {
		return 0, errors.New("nil source address")
	}

	hfiUnit := C.getHFIUnit(f.info.src_addr)
	if hfiUnit == allHFIUsed {
		return 0, errors.New("all HFI used")
	}
	return int(hfiUnit), nil
}

func (f *FIInfo) Domain() string {
	if !f.HasDomain() {
		return ""
	}
	return C.GoString(f.info.domain_attr.name)
}

func (f *FIInfo) HasDomain() bool {
	return f != nil && f.info != nil && f.info.domain_attr != nil && f.info.domain_attr.name != nil
}

func (f *FIInfo) Provider() string {
	if !f.HasProvider() {
		return ""
	}
	return C.GoString(f.info.fabric_attr.prov_name)
}

func (f *FIInfo) HasProvider() bool {
	return f != nil && f.info != nil && f.info.fabric_attr != nil && f.info.fabric_attr.prov_name != nil
}

func GetFIList(provider string) ([]*FIInfo, func(), error) {
	info, cleanup, err := getFIInfo(provider)
	if err != nil {
		return nil, nil, err
	}

	result := []*FIInfo{}
	for fi := info; fi != nil; fi = fi.next {
		result = append(result, &FIInfo{info: fi})
	}
	return result, cleanup, nil
}

func getFIInfo(provider string) (*C.struct_fi_info, func(), error) {
	hints := C.fi_allocinfo()
	if hints == nil {
		return nil, nil, errors.New("fi_allocinfo failed for hint")
	}
	defer C.fi_freeinfo(hints)

	hints.fabric_attr.prov_name = C.CString(provider)

	var info *C.struct_fi_info
	C.fi_getinfo(cVersion(), nil, nil, 0, hints, &info)
	if info == nil {
		return nil, nil, errors.Errorf("fi_getinfo failed for provider %q", provider)
	}

	return info, func(){
		C.fi_freeinfo(info)
	}, nil
}

func cVersion() C.uint {
	return C.uint(libFabricMajorVersion<<16 | libFabricMinorVersion)
}
