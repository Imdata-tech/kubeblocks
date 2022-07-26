/*
Copyright © 2022 The OpenCli Authors

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

package utils

import (
	"fmt"
	"github.com/infracreate/opencli/pkg/types"
	"github.com/pkg/errors"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

var (
	// Info print message
	Info func(a ...interface{})
	// Infof print message with format
	Infof func(format string, a ...interface{})
	// InfoP print message with padding
	InfoP func(padding int, a ...interface{})
	// Errf print error with format
	Errf func(format string, a ...interface{})
)

func init() {
	Info = func(a ...interface{}) {
		fmt.Println(a...)
	}
	Errf = func(format string, a ...interface{}) {
		fmt.Printf(format, a...)
	}
	Infof = func(format string, a ...interface{}) {
		fmt.Printf(format, a...)
	}
	InfoP = func(padding int, a ...interface{}) {
		fmt.Printf("%*s", padding, "")
		fmt.Println(a...)
	}
	if _, err := GetCliHomeDir(); err != nil {
		fmt.Println("Failed to build opencli home dir:", err)
	}
}

// CleanUpPlayground removes the playground directory
func CleanUpPlayground() error {
	tempDir, err := GetTempDir()
	if err != nil {
		return err
	}
	return os.RemoveAll(tempDir)
}

// CloseQuietly closes `io.Closer` quietly. Very handy and helpful for code
// quality too.
func CloseQuietly(d io.Closer) {
	_ = d.Close()
}

// GetCliHomeDir return opencli home dir
func GetCliHomeDir() (string, error) {
	var cliHome string
	if custom := os.Getenv(types.OpenCliHomeEnv); custom != "" {
		cliHome = custom
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		cliHome = filepath.Join(home, types.OpenCliDefaultHome)
	}
	if _, err := os.Stat(cliHome); err != nil && os.IsNotExist(err) {
		err := os.MkdirAll(cliHome, 0750)
		if err != nil {
			return "", errors.Wrap(err, "error when create opencli home directory")
		}
	}
	return cliHome, nil
}

func GetTempDir() (string, error) {
	dir, err := GetCliHomeDir()
	if err != nil {
		return "", err
	}
	tmpDir := filepath.Join(dir, "tmp")
	if err := os.MkdirAll(tmpDir, 0700); err != nil {
		return "", err
	}
	return tmpDir, nil
}

func SaveToTemp(file fs.File, format string) (string, error) {
	tempDir, err := GetTempDir()
	if err != nil {
		return "", err
	}
	tempFile, err := ioutil.TempFile(tempDir, format)
	if err != nil {
		return "", err
	}
	defer CloseQuietly(tempFile)

	if _, err := io.Copy(tempFile, file); err != nil {
		return "", err
	}

	return tempFile.Name(), nil
}

// InfoBytes is a helper function to print a byte array
func InfoBytes(b []byte) {
	if len(b) != 0 {
		Info(string(b))
	}
}

// GetKubeconfigDir returns the kubeconfig directory.
func GetKubeconfigDir() string {
	var kubeconfigDir string
	switch runtime.GOOS {
	case types.GoosDarwin, types.GoosLinux:
		kubeconfigDir = filepath.Join(os.Getenv("HOME"), ".kube")
	case types.GoosWindows:
		kubeconfigDir = filepath.Join(os.Getenv("USERPROFILE"), ".kube")
	}
	return kubeconfigDir
}

func ConfigPath(name string) string {
	return filepath.Join(GetKubeconfigDir(), name)
}

func RemoveConfig(name string) error {
	if err := os.Remove(ConfigPath(name)); err != nil {
		return err
	}
	return nil
}
