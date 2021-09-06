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

	r.Step(demo.S(
		"Search for a policy in hub.kubewarden.io",
	), nil)

	r.Step(demo.S(
		"List policies",
	), demo.S("kwctl policies"))

	r.Step(demo.S(
		"Pull a policy",
	), demo.S("kwctl pull registry://ghcr.io/kubewarden/policies/safe-annotations:v0.1.0"))

	r.Step(demo.S(
		"List policies",
	), demo.S("kwctl policies"))

	r.Step(demo.S(
		"Inspect policy",
	), demo.S("kwctl inspect registry://ghcr.io/kubewarden/policies/safe-annotations:v0.1.0"))

	r.Step(demo.S(
		"Request with a letsencrypt-production issuer",
	), demo.S("bat test_data/valid-ingress.json"))

	r.Step(demo.S(
		"Evaluate request with a letsencrypt-production issuer",
	), demo.S(`kwctl run \
      --settings-json '{"constrained_annotations": {"cert-manager.io/cluster-issuer": "letsencrypt-production"}}' \
      --request-path test_data/valid-ingress.json \
      registry://ghcr.io/kubewarden/policies/safe-annotations:v0.1.0 | jq`))

	r.Step(demo.S(
		"Request with a letsencrypt-staging issuer",
	), demo.S("bat test_data/invalid-ingress.json"))

	r.Step(demo.S(
		"Evaluate request with a letsencrypt-staging issuer",
	), demo.S(`kwctl run \
      --settings-json '{"constrained_annotations": {"cert-manager.io/cluster-issuer": "letsencrypt-production"}}' \
      --request-path test_data/invalid-ingress.json \
      registry://ghcr.io/kubewarden/policies/safe-annotations:v0.1.0 | jq`))

	return r
}

func policyServerRun() *demo.Run {
	r := demo.NewRun(
		"Running policies on the policy-server",
	)

	r.Setup(setupKubernetes)
	r.Cleanup(cleanupKubernetes)

	r.Step(demo.S(
		"List policies",
	), demo.S("kwctl policies"))

	r.Step(demo.S(
		"Pull a policy",
	), demo.S("kwctl pull registry://ghcr.io/kubewarden/policies/safe-annotations:v0.1.0"))

	r.Step(demo.S(
		"List policies",
	), demo.S("kwctl policies"))

	r.Step(demo.S(
		"Generate Kubernetes manifest",
	), demo.S(`kwctl manifest \
      --type ClusterAdmissionPolicy \
      registry://ghcr.io/kubewarden/policies/safe-annotations:v0.1.0 | \
        yq '.metadata.name = "kubecon-na-21"' | \
        yq '.spec.settings.constrained_annotations."cert-manager.io/cluster-issuer" = "letsencrypt-production"'`))

	r.Step(demo.S(
		"Apply Kubernetes manifest",
	), demo.S(`kwctl manifest \
      --type ClusterAdmissionPolicy \
      registry://ghcr.io/kubewarden/policies/safe-annotations:v0.1.0 | \
        yq '.metadata.name = "kubecon-na-21"' | \
        yq '.spec.settings.constrained_annotations."cert-manager.io/cluster-issuer" = "letsencrypt-production"' | \
          kubectl apply -f -`))

	return r
}

func cleanupKwctl() error {
	os.RemoveAll(filepath.Join(os.Getenv("HOME"), ".cache", "kubewarden"))
	return nil
}

func setupKubernetes() error {
	cleanupKwctl()
	return nil
}

func cleanupKubernetes() error {
	cleanupKwctl()
	return nil
}
