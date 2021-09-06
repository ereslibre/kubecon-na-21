package main

import (
	"os"
	"path/filepath"

	demo "github.com/saschagrunert/demo"
)

func main() {
	d := demo.New()
	d.Add(kwctlRun(), "kwctl demo", "kwctl demo")
	d.Add(policyServerRun(), "policy-server demo", "policy-server demo")
	d.Run()
}

func kwctlRun() *demo.Run {
	r := demo.NewRun(
		"Running policies with kwctl",
	)

	r.Setup(cleanupKwctl)
	r.Cleanup(cleanupKwctl)

	return r
}

func policyServerRun() *demo.Run {
	r := demo.NewRun(
		"Running policies on the policy-server",
	)

	r.Setup(setupKubernetes)
	r.Cleanup(cleanupKubernetes)

	return r
}

func cleanupKwctl() error {
	os.RemoveAll(filepath.Join(os.Getenv("HOME"), ".cache", "kubewarden"))
	return nil
}

func setupKubernetes() error {
	return nil
}

func cleanupKubernetes() error {
	return nil
}
