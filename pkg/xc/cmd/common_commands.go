package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/xcloudnative/xcloud/pkg/log"
	"github.com/xcloudnative/xcloud/pkg/util"
)

// TODO Refactor to use util.Run or util.RunWithoutRetry?
//
//func (o *CommonOptions) runCommandFromDir(dir, name string, args ...string) error {
//	e := exec.Command(name, args...)
//	if dir != "" {
//		e.Dir = dir
//	}
//	e.Stdout = o.Out
//	e.Stderr = o.Err
//	os.Setenv("PATH", util.PathWithBinary())
//	err := e.Run()
//	if err != nil {
//		log.Errorf("Error: Command failed  %s %s\n", name, strings.Join(args, " "))
//	}
//	return err
//}

// RunCommand runs a command
func (o *CommonOptions) RunCommand(name string, args ...string) error {
	e := exec.Command(name, args...)
	if o.Verbose {
		e.Stdout = o.Out
		e.Stderr = o.Err
	}
	os.Setenv("PATH", util.PathWithBinary())
	err := e.Run()
	if err != nil {
		log.Errorf("Error: Command failed  %s %s\n", name, strings.Join(args, " "))
	}
	return err
}

//func (o *CommonOptions) runCommandVerbose(name string, args ...string) error {
//	e := exec.Command(name, args...)
//	e.Stdout = o.Out
//	e.Stderr = o.Err
//	os.Setenv("PATH", util.PathWithBinary())
//	err := e.Run()
//	if err != nil {
//		log.Errorf("Error: Command failed  %s %s\n", name, strings.Join(args, " "))
//	}
//	return err
//}

func (o *CommonOptions) runCommandBackground(name string, output io.Writer, verbose bool, args ...string) error {
	e := exec.Command(name, args...)
	e.Stdout = output
	e.Stderr = output
	os.Setenv("PATH", util.PathWithBinary())
	err := e.Start()
	if err != nil && verbose {
		log.Errorf("Error: Command failed to start  %s %s\n", name, strings.Join(args, " "))
	}
	return err
}

//func (o *CommonOptions) runCommandVerboseAt(dir string, name string, args ...string) error {
//	e := exec.Command(name, args...)
//	if dir != "" {
//		e.Dir = dir
//	}
//	e.Stdout = o.Out
//	e.Stderr = o.Err
//	os.Setenv("PATH", util.PathWithBinary())
//	err := e.Run()
//	if err != nil {
//		log.Errorf("Error: Command failed  %s %s\n", name, strings.Join(args, " "))
//	}
//	return err
//}
//
//func (o *CommonOptions) runCommandQuietly(name string, args ...string) error {
//	e := exec.Command(name, args...)
//	e.Stdout = ioutil.Discard
//	e.Stderr = ioutil.Discard
//	os.Setenv("PATH", util.PathWithBinary())
//	return e.Run()
//}
//
//func (o *CommonOptions) runCommandInteractive(interactive bool, name string, args ...string) error {
//	e := exec.Command(name, args...)
//	e.Stdout = o.Out
//	e.Stderr = o.Err
//	if interactive {
//		e.Stdin = os.Stdin
//	}
//	os.Setenv("PATH", util.PathWithBinary())
//	err := e.Run()
//	if err != nil {
//		log.Errorf("Error: Command failed  %s %s\n", name, strings.Join(args, " "))
//	}
//	return err
//}
//
//func (o *CommonOptions) runCommandInteractiveInDir(interactive bool, dir string, name string, args ...string) error {
//	e := exec.Command(name, args...)
//	e.Stdout = o.Out
//	e.Stderr = o.Err
//	if interactive {
//		e.Stdin = os.Stdin
//	}
//	if dir != "" {
//		e.Dir = dir
//	}
//	os.Setenv("PATH", util.PathWithBinary())
//	err := e.Run()
//	if err != nil {
//		log.Errorf("Error: Command failed  %s %s\n", name, strings.Join(args, " "))
//	}
//	return err
//}
//
// getCommandOutput evaluates the given command and returns the trimmed output
func (o *CommonOptions) getCommandOutput(dir string, name string, args ...string) (string, error) {
	os.Setenv("PATH", util.PathWithBinary())
	e := exec.Command(name, args...)
	if dir != "" {
		e.Dir = dir
	}
	data, err := e.CombinedOutput()
	text := string(data)
	text = strings.TrimSpace(text)
	if err != nil {
		return "", fmt.Errorf("Command failed '%s %s': %s %s\n", name, strings.Join(args, " "), text, err)
	}
	return text, err
}
