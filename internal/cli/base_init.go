package cli

import (
	"errors"
	"fmt"
	"path/filepath"

	clientpkg "github.com/hashicorp/vagrant/internal/client"
	configpkg "github.com/hashicorp/vagrant/internal/config"
)

// This file contains the various methods that are used to perform
// the Init call on baseCommand. They are broken down into individual
// smaller methods for readability but more importantly to power the
// "init" subcommand. This allows us to share as much logic as possible
// between Init and "init" to help ensure that "init" succeeding means that
// other commands will succeed as well.

// initConfig initializes the configuration.
func (c *baseCommand) initConfig(optional bool) (*configpkg.Config, error) {
	path, err := c.initConfigPath()
	if err != nil {
		return nil, err
	}

	if path == "" {
		if optional {
			return nil, nil
		}

		return nil, errors.New("A Vagrant configuration file is required but wasn't found.")
	}

	return c.initConfigLoad(path)
}

// initConfigPath returns the configuration path to load.
func (c *baseCommand) initConfigPath() (string, error) {
	// This configuarion is for the Vagrant process, not the same as a Vagrantfile
	path, err := configpkg.FindPath("", "vagrant-config.hcl")
	if err != nil {
		return "", fmt.Errorf("Error looking for a Vagrant configuration: %s", err)
	}

	return path, nil
}

// initConfigLoad loads the configuration at the given path.
func (c *baseCommand) initConfigLoad(path string) (*configpkg.Config, error) {
	cfg, err := configpkg.Load(path, filepath.Dir(path))
	if err != nil {
		return nil, err
	}

	// Validate
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// initClient initializes the client.
func (c *baseCommand) initClient() (*clientpkg.Basis, error) {
	// Start building our client options
	opts := []clientpkg.Option{
		clientpkg.WithLabels(c.flagLabels),
		clientpkg.WithSourceOverrides(c.flagRemoteSource),
		clientpkg.WithConfig(c.cfg),
	}

	// Create our client
	return clientpkg.New(c.Ctx, opts...)
}