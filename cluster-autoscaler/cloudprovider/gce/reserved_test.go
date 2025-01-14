/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package gce

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateKernelReservedLinux(t *testing.T) {
	type testCase struct {
		physicalMemory int64
		reservedMemory int64
		osDistribution OperatingSystemDistribution
	}
	testCases := []testCase{
		{
			physicalMemory: 256 * MiB,
			reservedMemory: 4*MiB + kernelReservedMemory,
			osDistribution: OperatingSystemDistributionCOS,
		},
		{
			physicalMemory: 2 * GiB,
			reservedMemory: 32*MiB + kernelReservedMemory,
			osDistribution: OperatingSystemDistributionCOS,
		},
		{
			physicalMemory: 3 * GiB,
			reservedMemory: 48*MiB + kernelReservedMemory,
			osDistribution: OperatingSystemDistributionCOS,
		},
		{
			physicalMemory: 3.25 * GiB,
			reservedMemory: 52*MiB + kernelReservedMemory + swiotlbReservedMemory,
			osDistribution: OperatingSystemDistributionCOS,
		},
		{
			physicalMemory: 4 * GiB,
			reservedMemory: 64*MiB + kernelReservedMemory + swiotlbReservedMemory,
			osDistribution: OperatingSystemDistributionCOS,
		},
		{
			physicalMemory: 128 * GiB,
			reservedMemory: 2*GiB + kernelReservedMemory + swiotlbReservedMemory,
			osDistribution: OperatingSystemDistributionUbuntu,
		},
		{
			physicalMemory: 256 * MiB,
			reservedMemory: 4*MiB + kernelReservedMemory,
			osDistribution: OperatingSystemDistributionUbuntu,
		},
		{
			physicalMemory: 2 * GiB,
			reservedMemory: 32*MiB + kernelReservedMemory,
			osDistribution: OperatingSystemDistributionUbuntu,
		},
		{
			physicalMemory: 3 * GiB,
			reservedMemory: 48*MiB + kernelReservedMemory,
			osDistribution: OperatingSystemDistributionUbuntu,
		},
		{
			physicalMemory: 3.25 * GiB,
			reservedMemory: 52*MiB + kernelReservedMemory + swiotlbReservedMemory,
			osDistribution: OperatingSystemDistributionUbuntu,
		},
		{
			physicalMemory: 4 * GiB,
			reservedMemory: 64*MiB + kernelReservedMemory + swiotlbReservedMemory,
			osDistribution: OperatingSystemDistributionUbuntu,
		},
		{
			physicalMemory: 128 * GiB,
			reservedMemory: 2*GiB + kernelReservedMemory + swiotlbReservedMemory,
			osDistribution: OperatingSystemDistributionUbuntu,
		},
	}
	for idx, tc := range testCases {
		t.Run(fmt.Sprintf("%v", idx), func(t *testing.T) {
			reserved := CalculateKernelReserved(tc.physicalMemory, OperatingSystemLinux, tc.osDistribution)
			if tc.osDistribution == OperatingSystemDistributionUbuntu {
				assert.Equal(t, tc.reservedMemory+int64(math.Min(correctionConstant*float64(tc.physicalMemory), maximumCorrectionValue)+ubuntuSpecificOffset), reserved)
			} else if tc.osDistribution == OperatingSystemDistributionCOS {
				assert.Equal(t, tc.reservedMemory+int64(math.Min(correctionConstant*float64(tc.physicalMemory), maximumCorrectionValue)), reserved)
			}
		})
	}
}
