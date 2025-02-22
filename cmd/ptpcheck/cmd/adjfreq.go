/*
Copyright (c) Facebook, Inc. and its affiliates.

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

package cmd

import (
	"fmt"
	"math"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/facebook/time/phc"
)

// flag
var dev string
var freq float64

func init() {
	RootCmd.AddCommand(adjFreqCmd)
	adjFreqCmd.Flags().StringVarP(&dev, "device", "d", "/dev/ptp0", "PTP device to get frequnency from")
	adjFreqCmd.Flags().Float64VarP(&freq, "set", "s", math.NaN(), "New PHC frequency")
}

func doAdjFreq(device string, freq float64) error {
	curFreq, err := phc.FrequencyPPBFromDevice(device)
	if err != nil {
		return err
	}
	log.Infof("Current device frequency: %f", curFreq)

	maxFreq, err := phc.MaxFreqAdjPPBFromDevice(device)
	if err != nil {
		return err
	}
	log.Infof("Device supports frequency range [%f,%f]", -maxFreq, maxFreq)

	if math.IsNaN(freq) {
		return nil
	}

	if freq < -maxFreq || freq > maxFreq {
		return fmt.Errorf("Frequncy %f is out supported range", freq)
	}

	log.Infof("Setting new frequency value %f", freq)
	err = phc.ClockAdjFreq(device, freq)

	return err
}

var adjFreqCmd = &cobra.Command{
	Use:   "adjfreq",
	Short: "Print PHC frequency information. Use `-set <freq>` to change the frequency",
	Run: func(cmd *cobra.Command, args []string) {
		ConfigureVerbosity()

		if err := doAdjFreq(dev, freq); err != nil {
			log.Fatal(err)
		}
	},
}
