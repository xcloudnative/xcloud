package kube

import (
	"fmt"

	"github.com/xcloudnative/xcloud/pkg/apis/jenkins.io/v1"
	"github.com/xcloudnative/xcloud/pkg/client/clientset/versioned"
	"github.com/xcloudnative/xcloud/pkg/gits"
)

var (
	DefaultQuickstartLocations = []v1.QuickStartLocation{
		{
			GitURL:   gits.GitHubURL,
			GitKind:  gits.KindGitHub,
			Owner:    "jenkins-x-quickstarts",
			Includes: []string{"*"},
			Excludes: []string{"WIP-*"},
		},
	}
)

// GetQuickstartLocations returns the current quickstart locations. If no locations are defined
// yet lets return the defaults
func GetQuickstartLocations(jxClient versioned.Interface, ns string) ([]v1.QuickStartLocation, error) {
	var answer []v1.QuickStartLocation
	env, err := EnsureDevEnvironmentSetup(jxClient, ns)
	if err != nil {
		return answer, err
	}
	if env == nil {
		return answer, fmt.Errorf("No Development environment found for namespace %s", ns)
	}

	answer = env.Spec.TeamSettings.QuickstartLocations
	if len(answer) == 0 {
		answer = DefaultQuickstartLocations
	}
	return answer, nil
}
