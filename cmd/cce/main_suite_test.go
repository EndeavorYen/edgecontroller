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

package main_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"

	cce "github.com/smartedgemec/controller-ce"
)

var (
	ctrl *gexec.Session
	node *gexec.Session
	ctx  = context.Background()
)

var _ = BeforeSuite(func() {
	startup()
})

var _ = AfterSuite(func() {
	shutdown()
})

func TestApplicationClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Controller CE API Suite")
}

func startup() {
	By("Building the controller")
	exe, err := gexec.Build("github.com/smartedgemec/controller-ce/cmd/cce")
	Expect(err).ToNot(HaveOccurred(), "Problem building service")

	cmd := exec.Command(exe,
		"-dsn", "root:beer@tcp(:8083)/controller_ce",
		"-httpPort", "8080",
		"-grpcPort", "8081")

	By("Starting the controller")
	ctrl, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).ToNot(HaveOccurred(), "Problem starting service")

	By("Verifying that the controller started successfully")
	Eventually(ctrl.Err, 3).Should(gbytes.Say(
		"Controller CE ready"),
		"Service did not start in time")

	By("Building the node")
	exe, err = gexec.Build(
		"github.com/smartedgemec/controller-ce/test/node/grpc")
	Expect(err).ToNot(HaveOccurred(), "Problem building node")

	cmd = exec.Command(exe,
		"-port", "8082")

	By("Starting the node")
	node, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).ToNot(HaveOccurred(), "Problem starting node")

	By("Verifying that the node started successfully")
	Eventually(node.Err, 3).Should(gbytes.Say(
		"test-node: listening on port: 8082"),
		"Node did not start in time")
}

func shutdown() {
	if ctrl != nil {
		By("Stopping the controller service")
		ctrl.Kill()
	}
	if node != nil {
		By("Stopping the test node")
		node.Kill()
	}
}

func postApps(appType string) (id string) {
	By("Sending a POST /apps request")
	resp, err := http.Post(
		"http://127.0.0.1:8080/apps",
		"application/json",
		strings.NewReader(fmt.Sprintf(`
			{
				"type": "%s",
				"name": "%s app",
				"vendor": "smart edge",
				"description": "my %s app",
				"image": "http://www.test.com/my_%s_app.tar.gz",
				"cores": 4,
				"memory": 1024
			}`, appType, appType, appType, appType)))

	By("Verifying a 201 Created response")
	Expect(err).ToNot(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(http.StatusCreated))

	By("Reading the response body")
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	var respBody struct {
		ID string
	}

	By("Unmarshalling the response")
	Expect(json.Unmarshal(body, &respBody)).To(Succeed())

	return respBody.ID
}

func getApp(id string) *cce.App {
	By("Sending a GET /apps/{id} request")
	resp, err := http.Get(
		fmt.Sprintf("http://127.0.0.1:8080/apps/%s", id))

	By("Verifying a 200 OK response")
	Expect(err).ToNot(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(http.StatusOK))

	By("Reading the response body")
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	var app cce.App

	By("Unmarshalling the response")
	Expect(json.Unmarshal(body, &app)).To(Succeed())

	return &app
}

func postVNFs(vnfType string) (id string) {
	By("Sending a POST /vnfs request")
	resp, err := http.Post(
		"http://127.0.0.1:8080/vnfs",
		"application/json",
		strings.NewReader(fmt.Sprintf(`
			{
				"type": "%s",
				"name": "%s vnf",
				"vendor": "smart edge",
				"description": "my %s vnf",
				"image": "http://www.test.com/my_%s_vnf.tar.gz",
				"cores": 4,
				"memory": 1024
			}`, vnfType, vnfType, vnfType, vnfType)))

	By("Verifying a 201 Created response")
	Expect(err).ToNot(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(http.StatusCreated))

	By("Reading the response body")
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	var respBody struct {
		ID string
	}

	By("Unmarshalling the response")
	Expect(json.Unmarshal(body, &respBody)).To(Succeed())

	return respBody.ID
}

func getVNF(id string) *cce.VNF {
	By("Sending a GET /vnfs/{id} request")
	resp, err := http.Get(
		fmt.Sprintf("http://127.0.0.1:8080/vnfs/%s", id))

	By("Verifying a 200 OK response")
	Expect(err).ToNot(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(http.StatusOK))

	By("Reading the response body")
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	var vnf cce.VNF

	By("Unmarshalling the response")
	Expect(json.Unmarshal(body, &vnf)).To(Succeed())

	return &vnf
}

func postDNSConfigs() (id string) {
	By("Sending a POST /dns_configs request")
	resp, err := http.Post(
		"http://127.0.0.1:8080/dns_configs",
		"application/json",
		strings.NewReader(`
			{
				"name": "dns config 123",
				"a_records": [{
					"name": "a record 1",
					"description": "description 1",
					"ips": [
						"172.16.55.43",
						"172.16.55.44"
					]
				}],
				"forwarders": [{
					"name": "forwarder 1",
					"description": "description 1",
					"ip": "8.8.8.8"
				}, {
					"name": "forwarder 2",
					"description": "description 2",
					"ip": "1.1.1.1"
				}]
			}`))

	By("Verifying a 201 Created response")
	Expect(err).ToNot(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(http.StatusCreated))

	By("Reading the response body")
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	var respBody struct {
		ID string
	}

	By("Unmarshalling the response")
	Expect(json.Unmarshal(body, &respBody)).To(Succeed())

	return respBody.ID
}

func getDNSConfig(id string) *cce.DNSConfig {
	By("Sending a GET /dns_configs/{id} request")
	resp, err := http.Get(
		fmt.Sprintf("http://127.0.0.1:8080/dns_configs/%s", id))

	By("Verifying a 200 OK response")
	Expect(err).ToNot(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(http.StatusOK))

	By("Reading the response body")
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	var dnsConfig cce.DNSConfig

	By("Unmarshalling the response")
	Expect(json.Unmarshal(body, &dnsConfig)).To(Succeed())

	return &dnsConfig
}

func postDNSConfigsAppAliases(
	dnsConfigID string,
	appID string,
) (id string) {
	By("Sending a POST /dns_configs_app_aliases request")
	resp, err := http.Post(
		"http://127.0.0.1:8080/dns_configs_app_aliases",
		"application/json",
		strings.NewReader(fmt.Sprintf(`
			{
				"dns_config_id": "%s",
				"name": "dns config app alias",
				"description": "my dns config app alias",
				"app_id": "%s"
			}`, dnsConfigID, appID)))

	By("Verifying a 201 Created response")
	Expect(err).ToNot(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(http.StatusCreated))

	By("Reading the response body")
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	var respBody struct {
		ID string
	}

	By("Unmarshalling the response")
	Expect(json.Unmarshal(body, &respBody)).To(Succeed())

	return respBody.ID
}

func getDNSConfigsAppAlias(
	id string,
) *cce.DNSConfigAppAlias {
	By("Sending a GET /dns_configs_app_aliases/{id} request")
	resp, err := http.Get(
		fmt.Sprintf("http://127.0.0.1:8080/dns_configs_app_aliases/%s", id))

	By("Verifying a 200 OK response")
	Expect(err).ToNot(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(http.StatusOK))

	By("Reading the response body")
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	var dnsConfigAppAlias cce.DNSConfigAppAlias

	By("Unmarshalling the response")
	Expect(json.Unmarshal(body, &dnsConfigAppAlias)).To(Succeed())

	return &dnsConfigAppAlias
}

func postDNSConfigsVNFAliases(
	dnsConfigID string,
	vnfID string,
) (id string) {
	By("Sending a POST /dns_configs_vnf_aliases request")
	resp, err := http.Post(
		"http://127.0.0.1:8080/dns_configs_vnf_aliases",
		"application/json",
		strings.NewReader(fmt.Sprintf(`
			{
				"dns_config_id": "%s",
				"name": "dns config vnf alias",
				"description": "my dns config vnf alias",
				"vnf_id": "%s"
			}`, dnsConfigID, vnfID)))

	By("Verifying a 201 Created response")
	Expect(err).ToNot(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(http.StatusCreated))

	By("Reading the response body")
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	var respBody struct {
		ID string
	}

	By("Unmarshalling the response")
	Expect(json.Unmarshal(body, &respBody)).To(Succeed())

	return respBody.ID
}

func getDNSConfigsVNFAlias(
	id string,
) *cce.DNSConfigVNFAlias {
	By("Sending a GET /dns_configs_vnf_aliases/{id} request")
	resp, err := http.Get(
		fmt.Sprintf("http://127.0.0.1:8080/dns_configs_vnf_aliases/%s", id))

	By("Verifying a 200 OK response")
	Expect(err).ToNot(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(http.StatusOK))

	By("Reading the response body")
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	var dnsConfigVNFAlias cce.DNSConfigVNFAlias

	By("Unmarshalling the response")
	Expect(json.Unmarshal(body, &dnsConfigVNFAlias)).To(Succeed())

	return &dnsConfigVNFAlias
}

func postNodes() (id string) {
	By("Sending a POST /nodes request")
	resp, err := http.Post(
		"http://127.0.0.1:8080/nodes",
		"application/json",
		strings.NewReader(`
			{
				"name": "node123",
				"location": "smart edge lab",
				"serial": "abc123"
			}`))

	By("Verifying a 201 Created response")
	Expect(err).ToNot(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(http.StatusCreated))

	By("Reading the response body")
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	var respBody struct {
		ID string
	}

	By("Unmarshalling the response")
	Expect(json.Unmarshal(body, &respBody)).To(Succeed())

	return respBody.ID
}

func getNode(id string) *cce.Node {
	By("Sending a GET /nodes/{id} request")
	resp, err := http.Get(
		fmt.Sprintf("http://127.0.0.1:8080/nodes/%s", id))

	By("Verifying a 200 OK response")
	Expect(err).ToNot(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(http.StatusOK))

	By("Reading the response body")
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	var node cce.Node

	By("Unmarshalling the response")
	Expect(json.Unmarshal(body, &node)).To(Succeed())

	return &node
}

func postNodesDNSConfigs(nodeID, dnsConfigID string) (id string) {
	By("Sending a POST /nodes_dns_configs request")
	resp, err := http.Post(
		"http://127.0.0.1:8080/nodes_dns_configs",
		"application/json",
		strings.NewReader(fmt.Sprintf(`
			{
				"node_id": "%s",
				"dns_config_id": "%s"
			}`, nodeID, dnsConfigID)))

	By("Verifying a 201 Created response")
	Expect(err).ToNot(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(http.StatusCreated))

	By("Reading the response body")
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	var respBody struct {
		ID string
	}

	By("Unmarshalling the response")
	Expect(json.Unmarshal(body, &respBody)).To(Succeed())

	fmt.Println(respBody.ID)
	return respBody.ID
}

func getNodesDNSConfig(id string) *cce.NodeDNSConfig {
	By("Sending a GET /nodes_dns_configs/{id} request")
	resp, err := http.Get(
		fmt.Sprintf("http://127.0.0.1:8080/nodes_dns_configs/%s", id))

	By("Verifying a 200 OK response")
	Expect(err).ToNot(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(http.StatusOK))

	By("Reading the response body")
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	var nodeDNSConfig cce.NodeDNSConfig

	By("Unmarshalling the response")
	Expect(json.Unmarshal(body, &nodeDNSConfig)).To(Succeed())

	return &nodeDNSConfig
}

func postTrafficPolicies() (id string) {
	By("Sending a POST /traffic_policies request")
	resp, err := http.Post(
		"http://127.0.0.1:8080/traffic_policies",
		"application/json",
		strings.NewReader(`
		{
			"rules": [{
				"description": "test-rule-1",
				"priority": 1,
				"source": {
					"description": "test-source-1",
					"macs": {
						"mac_addresses": [
							"F0-59-8E-7B-36-8A",
							"23-20-8E-15-89-D1",
							"35-A4-38-73-35-45"
						]
					},
					"ip": {
						"address": "223.1.1.0",
						"mask": 16,
						"begin_port": 2000,
						"end_port": 2012,
						"protocol": "tcp"
					},
					"gtp": {
						"address": "10.6.7.2",
						"mask": 12,
						"imsis": [
							"310150123456789",
							"310150123456790",
							"310150123456791"
						]
					}
				},
				"destination": {
					"description": "test-destination-1",
					"macs": {
						"mac_addresses": [
							"7D-C2-3A-1C-63-D9",
							"E9-6B-D1-D2-1A-6B",
							"C8-32-A9-43-85-55"
						]
					},
					"ip": {
						"address": "64.1.1.0",
						"mask": 16,
						"begin_port": 1000,
						"end_port": 1012,
						"protocol": "tcp"
					},
					"gtp": {
						"address": "108.6.7.2",
						"mask": 4,
						"imsis": [
							"310150123456792",
							"310150123456793",
							"310150123456794"
						]
					}
				},
				"target": {
					"description": "test-target-1",
					"action": "accept",
					"mac": {
						"mac_address": "C7-5A-E7-98-1B-A3"
					},
					"ip": {
						"address": "123.2.3.4",
						"port": 1600
					}
				}
			}]
		}`))

	By("Verifying a 201 Created response")
	Expect(err).ToNot(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(http.StatusCreated))

	By("Reading the response body")
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	var respBody struct {
		ID string
	}

	By("Unmarshalling the response")
	Expect(json.Unmarshal(body, &respBody)).To(Succeed())

	return respBody.ID
}

func getTrafficPolicy(id string) *cce.TrafficPolicy {
	By("Sending a GET /traffic_policies/{id} request")
	resp, err := http.Get(
		fmt.Sprintf("http://127.0.0.1:8080/traffic_policies/%s", id))

	By("Verifying a 200 OK response")
	Expect(err).ToNot(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(http.StatusOK))

	By("Reading the response body")
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	var trafficPolicy cce.TrafficPolicy

	By("Unmarshalling the response")
	Expect(json.Unmarshal(body, &trafficPolicy)).To(Succeed())

	return &trafficPolicy
}
