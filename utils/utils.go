// Copyright 2015 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"errors"
	"fmt"
	"github.com/prometheus/common/log"
	"io"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const CommandExecutionTimeout = 5 * time.Second

func splitToInts(str string, sep string) (ints []int, err error) {
	for _, part := range strings.Split(str, sep) {
		i, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("could not split '%s' because %s is no int: %s", str, part, err)
		}
		ints = append(ints, i)
	}
	return ints, nil
}

func readUintFromFile(path string) (uint64, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, err
	}
	value, err := strconv.ParseUint(strings.TrimSpace(string(data)), 10, 64)
	if err != nil {
		return 0, err
	}
	return value, nil
}
func cmd(stdin string, name string) (string, error) {

	cmd := exec.Command("sh", "-c", name)
	stdinPipe, _ := cmd.StdinPipe()
	stdoutPipe, _ := cmd.StdoutPipe()
	stderrPipe, _ := cmd.StderrPipe()

	defer func() {
		stdinPipe.Close()
		stdoutPipe.Close()
		stderrPipe.Close()
	}()

	if err := cmd.Start(); err != nil {
		return "", err
	}

	var (
		errc                     = make(chan error, 1)
		stdoutBytes, stderrBytes []byte
	)

	readAll := func(out *[]byte, reader *io.ReadCloser) {
		*out, _ = ioutil.ReadAll(*reader)
	}

	go func() {
		readAll(&stdoutBytes, &stdoutPipe)
		readAll(&stderrBytes, &stderrPipe)
		errc <- cmd.Wait()
	}()

	if stdin != "" {
		_, err := io.WriteString(stdinPipe, stdin)
		if err != nil {
			return "", err
		}
	}

	times := 1
	for {
		select {
		case <-time.After(CommandExecutionTimeout):
			log.Errorf(`command "%s" run more than %d seconds.`, name, int(CommandExecutionTimeout.Seconds())*times)
		case err := <-errc:
			stdout := string(stdoutBytes)
			stderr := string(stderrBytes)
			log.Infoln("=============", "Command Execution", "=============")
			log.Infoln("command:", "", name)
			log.Infoln("output:", "", stdout)
			log.Infoln("stderr:", "", stderr)
			if err != nil {
				return stdout, errors.New(fmt.Sprintf(`execute error with stderr: %s`, stderr))
			} else {
				return string(stdoutBytes), nil
			}
		}
		times++
	}
}

func Cmd(name string) (string, error) {
	return cmd("", name)
}

func execCommand(stdin string, name string, arg ...string) (string, error) {
	return cmd(stdin, fmt.Sprintf("%s %s", name, strings.Join(arg, " ")))
}
