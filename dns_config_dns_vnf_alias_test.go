// Copyright 2019 Smart-Edge.com, Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cce_test

import (
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	cce "github.com/smartedgemec/controller-ce"
)

var _ = Describe("Entities: DNSConfigVNFAlias", func() {
	var (
		cfgAlias *cce.DNSConfigVNFAlias
	)

	BeforeEach(func() {
		cfgAlias = &cce.DNSConfigVNFAlias{
			ID:          "7557c16e-1777-490e-9f4f-c59936b2468f",
			DNSConfigID: "84c1f7b9-53e7-408e-9223-deab73befc54",
			Name:        "test-dns-config-vnf-alias",
			Description: "test-description",
			VNFID:       "28bbfdb2-dace-421d-a680-9ae893a95d37",
		}
	})

	Describe("GetTableName", func() {
		It(`Should return "dns_configs_vnf_aliases"`, func() {
			Expect(cfgAlias.GetTableName()).To(Equal(
				"dns_configs_vnf_aliases"))
		})
	})

	Describe("GetID", func() {
		It("Should return the ID", func() {
			Expect(cfgAlias.GetID()).To(Equal(
				"7557c16e-1777-490e-9f4f-c59936b2468f"))
		})
	})

	Describe("SetID", func() {
		It("Should set and return the updated ID", func() {
			By("Setting the ID")
			cfgAlias.SetID("456")

			By("Getting the updated ID")
			Expect(cfgAlias.ID).To(Equal("456"))
		})
	})

	Describe("Validate", func() {
		It("Should return an error if ID is not a UUID", func() {
			cfgAlias.ID = "123"
			Expect(cfgAlias.Validate()).To(MatchError("id not a valid uuid"))
		})

		It("Should return an error if DNSConfigID is not a UUID", func() {
			cfgAlias.DNSConfigID = "123"
			Expect(cfgAlias.Validate()).To(MatchError(
				"dns_config_id not a valid uuid"))
		})

		It("Should return an error if Name is empty", func() {
			cfgAlias.Name = ""
			Expect(cfgAlias.Validate()).To(MatchError("name cannot be empty"))
		})

		It("Should return an error if Description is empty", func() {
			cfgAlias.Description = ""
			Expect(cfgAlias.Validate()).To(MatchError(
				"description cannot be empty"))
		})

		It("Should return an error if VNFID is not a UUID", func() {
			cfgAlias.VNFID = "123"
			Expect(cfgAlias.Validate()).To(MatchError(
				"vnf_id not a valid uuid"))
		})
	})

	Describe("String", func() {
		It("Should return the string value", func() {
			Expect(cfgAlias.String()).To(Equal(strings.TrimSpace(`
DNSConfigVNFAlias[
    ID: 7557c16e-1777-490e-9f4f-c59936b2468f
    DNSConfigID: 84c1f7b9-53e7-408e-9223-deab73befc54
    VNFID: 28bbfdb2-dace-421d-a680-9ae893a95d37
]`,
			)))
		})
	})
})
