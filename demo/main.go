package main

import (
	"os"
	"os/exec"
	"path/filepath"

	demo "github.com/saschagrunert/demo"
)

func main() {
	d := demo.New()
	d.Add(policyServerRun(), "policy-server demo", "policy-server demo")
	d.Add(gatekeeperPolicyBuildAndRun(), "gatekeeper policy build and run demo", "gatekeeper policy build and run demo")
	d.Run()
}

func policyServerRun() *demo.Run {
	r := demo.NewRun(
		"Running policies on the policy-server",
	)

	r.Setup(setupKubernetes)
	r.Cleanup(cleanupKubernetes)

	policyServer(r)

	return r
}

func policyServer(r *demo.Run) {
	r.Step(demo.S(
		"Show cluster admission policy",
	), demo.S("bat test_data/letsencrypt-production-manifest.yaml"))

	r.Step(demo.S(
		"Deploy cluster admission policy",
	), demo.S(
		"kubectl apply -f test_data/letsencrypt-production-manifest.yaml",
	))

	r.Step(demo.S(
		"Wait for our policy to be active",
	), demo.S(
		"kubectl wait --for=condition=PolicyServerWebhookConfigurationReconciled clusteradmissionpolicy letsencrypt-production-ingress",
	))

	r.Step(demo.S(
		"Ingress with a letsencrypt-production issuer",
	), demo.S("bat test_data/production-ingress-resource.yaml"))

	r.Step(demo.S(
		"Deploy an Ingress resource with a letsencrypt-production issuer",
	), demo.S("kubectl apply -f test_data/production-ingress-resource.yaml"))

	r.Step(demo.S(
		"Ingress with a letsencrypt-staging issuer",
	), demo.S("bat test_data/staging-ingress-resource.yaml"))

	r.StepCanFail(demo.S(
		"Deploy an Ingress resource with a letsencrypt-staging issuer",
	), demo.S("kubectl apply -f test_data/staging-ingress-resource.yaml"))
}

func gatekeeperPolicyBuildAndRun() *demo.Run {
	r := demo.NewRun(
		"Running a gatekeeper policy",
	)

	r.Step(demo.S(
		"Build policy",
	), demo.S(
		"opa build -t wasm -e echo/violation -o gatekeeper/bundle.tar.gz gatekeeper/echo.rego &&",
		"tar -C gatekeeper -xf gatekeeper/bundle.tar.gz /policy.wasm",
	))

	r.Step(demo.S(
		"Run policy: accept the request",
	), demo.S(
		"kwctl run -e gatekeeper",
		`--settings-json '{"reject":false}'`,
		"--request-path test_data/having-label-ingress.json",
		"gatekeeper/policy.wasm | jq",
	))

	r.Step(demo.S(
		"Run policy: reject the request",
	), demo.S(
		"kwctl run -e gatekeeper",
		`--settings-json '{"reject":true, "rejection_message": "this is going to be rejected, no matter what"}'`,
		"--request-path test_data/having-label-ingress.json",
		"gatekeeper/policy.wasm | jq",
	))

	return r
}

func cleanupKwctl() error {
	os.RemoveAll(filepath.Join(os.Getenv("HOME"), ".cache", "kubewarden"))
	return nil
}

func setupKubernetes() error {
	cleanupKwctl()
	cleanupKubernetes()
	exec.Command("kubectl", "create", "namespace", "kubecon-na-21").Run()
	exec.Command("kubectl", "delete", "clusteradmissionpolicy", "--all").Run()
	return nil
}

func cleanupKubernetes() error {
	cleanupKwctl()
	exec.Command("kubectl", "delete", "namespace", "kubecon-na-21").Run()
	exec.Command("kubectl", "delete", "clusteradmissionpolicy", "--all").Run()
	return nil
}
