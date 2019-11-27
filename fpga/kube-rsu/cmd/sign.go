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
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// copy image for signing
func copyRTLFile(node string, file string) error {
	var err error
	var cmd *exec.Cmd

	// #nosec
	cmd = exec.Command("scp", file, node+":/temp/vran_images/")

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	go func() {
		if _, err = io.Copy(os.Stdout, stdout); err != nil {
			fmt.Println(err.Error())
		}
	}()
	go func() {
		if _, err = io.Copy(os.Stderr, stderr); err != nil {
			fmt.Println(err.Error())
		}
	}()

	err = cmd.Start()
	if err != nil {
		return err
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}

// signCmd represents the sign command
var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "Sign FPGA RTL image for RSU",
	Args:  cobra.MaximumNArgs(1),
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

		// copy RTL image to target node
		err := copyRTLFile(node, RTLFile)
		if err != nil {
			fmt.Println(err.Error())
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
		RSUJob.ObjectMeta.Name = "fpga-opae-"+node

		containerSpec.Args = []string{
			"./check_if_modules_loaded.sh && yes Y | " +
				"python3 /usr/local/bin/PACSign SR -t UPDATE -H openssl_manager -i " +
				"/root/images/" + RTLFile + " -o /root/images/SIGNED_" + RTLFile,
		}

		containerSpec.VolumeMounts = []corev1.VolumeMount{
			{
				Name:      "image-dir",
				MountPath: "/root/images",
				ReadOnly:  false,
			},
		}
		podSpec.NodeSelector["kubernetes.io/hostname"] = node
		podSpec.Volumes = []corev1.Volume{
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

		for i := 0; i < jobTimeout; i++ {
			// wait
			time.Sleep(1 * time.Second)
			// get job
			k8Job, err := jobsClient.Get(RSUJob.Name, metav1.GetOptions{})
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			if (k8Job.Status.Failed > 0) {
				fmt.Println("Job `"+k8Job.Name+"` failed!")
				break
			}

			if (k8Job.Status.Succeeded > 0) && (k8Job.Status.Active == 0) {
				fmt.Println("Job `"+k8Job.Name+"` completed successfully!")
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

	const help = `Sign FPGA RTL image for RSU

Usage:
  rsu sign -f <unsigned-RTL-img-file> -n <target-node>

Flags:
  -h, --help       help
  -f, --filename   unsigned RTL image file
  -n, --node       where the target FPGA card is plugged in
`
	// add `sign` command
	rsuCmd.AddCommand(signCmd)
	signCmd.Flags().StringP("filename", "f", "", "RTL image file")
	signCmd.MarkFlagRequired("filename")
	signCmd.Flags().StringP("node", "n", "", "where the target FPGA card is plugged in")
	signCmd.MarkFlagRequired("node")
	signCmd.SetHelpTemplate(help)
}
