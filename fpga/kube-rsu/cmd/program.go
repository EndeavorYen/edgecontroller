// Copyright 2019 Intel Corporation. All rights reserved
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

package rsu

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// programCmd represents the program command
var programCmd = &cobra.Command{
	Use:   "program",
	Short: "Program an FPGA device on a target node with an RTL image",
	Args:  cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		RTLFile, _ := cmd.Flags().GetString("filename")
		if RTLFile == "" {
			fmt.Println(errors.New("RTL image file missing"))
			return
		}

		node, _ := cmd.Flags().GetString("node")
		if node == "" {
			fmt.Println(errors.New("target node missing"))
			return
		}

		device, _ := cmd.Flags().GetString("device")
		if device == "" {
			fmt.Println(errors.New("target PCI device missing"))
			return
		}

		// retrieve .kube/config file
		kubeconfig := filepath.Join(
			os.Getenv("HOME"), ".kube", "config",
		)

		// use the current context in kubeconfig
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// create the clientset
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// edit K8 job with `program` command specifics
		podSpec := &(RSUJob.Spec.Template.Spec)
		containerSpec := &(RSUJob.Spec.Template.Spec.Containers[0])
		// assign unique job name for the target device
		d := strings.ReplaceAll(device, ":", "")
		RSUJob.ObjectMeta.Name = "fpga-opae-" + node + "-" + d

		containerSpec.Args = []string{
			"./check_if_modules_loaded.sh && fpgasupdate " +
				"/root/images/" + RTLFile + " " +
				device + " && rsu bmcimg " + device,
		}

		containerSpec.VolumeMounts = []corev1.VolumeMount{
			{
				Name:      "class",
				MountPath: "/sys/devices",
				ReadOnly:  false,
			},
			{
				Name:      "image-dir",
				MountPath: "/root/images",
				ReadOnly:  false,
			},
		}
		podSpec.NodeSelector["kubernetes.io/hostname"] = node
		podSpec.Volumes = []corev1.Volume{
			{
				Name: "class",
				VolumeSource: corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{
						Path: "/sys/devices",
					},
				},
			},
			{
				Name: "image-dir",
				VolumeSource: corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{
						Path: "/temp/vran_images",
					},
				},
			},
		}

		// create job in K8 environment
		jobsClient := clientset.BatchV1().Jobs(namespace)
		k8Job, err := jobsClient.Create(RSUJob)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		// print logs from pod
		logProcess, err := PrintJobLogs(clientset, k8Job)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer logProcess.Process.Kill()
		defer logProcess.Wait()

		// no timeout (this is a long process)
		for {
			// wait
			time.Sleep(1 * time.Second)
			// get job
			k8Job, err := jobsClient.Get(RSUJob.Name, metav1.GetOptions{})
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			if k8Job.Status.Failed > 0 {
				fmt.Println("Job `" + k8Job.Name + "` failed!")
				break
			}
			if (k8Job.Status.Succeeded > 0) && (k8Job.Status.Active == 0) {
				fmt.Println("Job `" + k8Job.Name + "` completed successfully!")
				break
			}
		}

		// delete job after completion
		err = jobsClient.Delete(k8Job.Name, &metav1.DeleteOptions{})
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	},
}

func init() {

	const help = `Program an FPGA device on a target node with an RTL image

Usage:
  rsu program -f <signed-RTL-img-file> -n <target-node> -d <target-device>

Flags:
  -h, --help       help
  -f, --filename   signed RTL image file
  -n, --node       where the target FPGA card is plugged in
  -d, --device     PCI ID of the target FPGA card
`
	// add `program` command
	rsuCmd.AddCommand(programCmd)
	programCmd.Flags().StringP("filename", "f", "", "RTL image file")
	programCmd.MarkFlagRequired("filename")
	programCmd.Flags().StringP("node", "n", "", "where the target FPGA card is plugged in")
	programCmd.MarkFlagRequired("node")
	programCmd.Flags().StringP("device", "d", "", "PCI ID of the target FPGA card")
	programCmd.MarkFlagRequired("device")
	programCmd.SetHelpTemplate(help)
}
