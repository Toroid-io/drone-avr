package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

type (
	Netrc struct {
		Machine  string
		Login    string
		Password string
	}

	// Commit handles commit information
	Commit struct {
		Tag string // tag if tag event
		Sha string // commit sha
	}

	Project struct {
		Source       string   `json:"source"`       // Project source folder
		Dependencies []string `json:"dependencies"` // Projects dependencies
		Command      string   `json:"command"`      // Build command
		Store        string   `json:"store"`        // Store artifacts
		Arguments    string   `json:"arguments"`    // Make arguments
	}

	// Plugin defines the KiCad plugin parameters
	Plugin struct {
		Projects []Project // Projects configuration
		Netrc    Netrc     // Authentication
		Commit   Commit    // Commit information
	}
)

func (p Plugin) Exec() error {

	err := writeNetrc(p.Netrc.Machine, p.Netrc.Login, p.Netrc.Password)
	if err != nil {
		return err
	}

	var cmds []*exec.Cmd

	for _, project := range p.Projects {

		if len(project.Source) == 0 {
			return errors.New("Source is not defined")
		}

		for _, dep := range project.Dependencies {
			cmds = append(cmds, commandClone(dep, project.Source))
		}

		// Make the project
		cmds = append(cmds, commandMake(project))

		// Store build artifacts
		if len(project.Store) > 0 {
			cmds = append(cmds, commandStore(project))
		}
	}

	// execute all commands in batch mode.
	for _, cmd := range cmds {
		if cmd != nil {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			trace(cmd)

			err := cmd.Run()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func commandClone(url string, dir string) *exec.Cmd {

	var cmd []string
	cmd = append(cmd, "cd", dir, "&&", "git", "clone", url)

	return exec.Command(
		"/bin/sh",
		"-c",
		strings.Join(cmd, " "),
	)
}

func commandMake(project Project) *exec.Cmd {

	var cmd []string
	if len(project.Command) > 0 {
		cmd = append(cmd, "cd", project.Source, "&&", project.Command)
	} else {
		cmd = append(cmd, "make", "-C", project.Source)
		if len(project.Arguments) > 0 {
			cmd = append(cmd, project.Arguments)
		}
	}

	return exec.Command(
		"/bin/sh",
		"-c",
		strings.Join(cmd, " "),
	)
}

func commandStore(project Project) *exec.Cmd {
	var opts []string
	opts = append(opts, "cd", project.Source, "&&", project.Store)

	return exec.Command(
		"/bin/sh",
		"-c",
		strings.Join(opts, " "),
	)
}

// trace writes each command to stdout with the command wrapped in an xml
// tag so that it can be extracted and displayed in the logs.
func trace(cmd *exec.Cmd) {
	fmt.Fprintf(os.Stdout, "+ %s\n", strings.Join(cmd.Args, " "))
}

// helper function to write a netrc file. [From drone-git]
func writeNetrc(machine, login, password string) error {
	if machine == "" {
		return nil
	}
	out := fmt.Sprintf(
		netrcFile,
		machine,
		login,
		password,
	)

	home := "/root"
	u, err := user.Current()
	if err == nil {
		home = u.HomeDir
	}
	path := filepath.Join(home, ".netrc")
	return ioutil.WriteFile(path, []byte(out), 0600)
}

const netrcFile = `
machine %s
login %s
password %s
`
