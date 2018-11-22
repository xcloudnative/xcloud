package cmd

import (
	// this is so that we load the auth plugins so we can connect to, say, GCP

	"github.com/xcloudnative/xcloud/pkg/client/clientset/versioned"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"


	"github.com/heptio/sonobuoy/pkg/dynamic"
	"io"

	"github.com/heptio/sonobuoy/pkg/client"
	"github.com/xcloudnative/xcloud/pkg/gits"
	"github.com/xcloudnative/xcloud/pkg/table"
	"gopkg.in/AlecAivazis/survey.v1/terminal"

	"github.com/jenkins-x/golang-jenkins"
	"github.com/xcloudnative/xcloud/pkg/auth"
	corev1 "k8s.io/api/core/v1"

	vaultoperatorclient "github.com/banzaicloud/bank-vaults/operator/pkg/client/clientset/versioned"
	buildclient "github.com/knative/build/pkg/client/clientset/versioned"
	metricsclient "k8s.io/metrics/pkg/client/clientset_generated/clientset"

	_ "k8s.io/client-go/plugin/pkg/client/auth"

)

// Factory is the interface defined for jx interactions via the cli
//go:generate pegomock generate github.com/xcloudnative/xcloud/pkg/jx/cmd Factory -o mocks/factory.go --generate-matchers
type Factory interface {
	WithBearerToken(token string) Factory
	
	ImpersonateUser(user string) Factory
	
	CreateJenkinsClient(kubeClient kubernetes.Interface, ns string, in terminal.FileReader, out terminal.FileWriter, errOut io.Writer) (gojenkins.JenkinsClient, error)
	
	GetJenkinsURL(kubeClient kubernetes.Interface, ns string) (string, error)
	
	CreateAuthConfigService(fileName string) (auth.AuthConfigService, error)
	
	CreateJenkinsAuthConfigService(kubernetes.Interface, string) (auth.AuthConfigService, error)
	
	CreateChartmuseumAuthConfigService() (auth.AuthConfigService, error)
	
	CreateIssueTrackerAuthConfigService(secrets *corev1.SecretList) (auth.AuthConfigService, error)
	
	CreateChatAuthConfigService(secrets *corev1.SecretList) (auth.AuthConfigService, error)
	
	CreateAddonAuthConfigService(secrets *corev1.SecretList) (auth.AuthConfigService, error)
	
	CreateClient() (kubernetes.Interface, string, error)
	
	CreateGitProvider(string, string, auth.AuthConfigService, string, bool, gits.Gitter, terminal.FileReader, terminal.FileWriter, io.Writer) (gits.GitProvider, error)
	
	CreateKubeConfig() (*rest.Config, error)

	CreateJXClient() (versioned.Interface, string, error)
	
	CreateApiExtensionsClient() (apiextensionsclientset.Interface, error)
	
	CreateDynamicClient() (*dynamic.APIHelper, string, error)
	
	CreateMetricsClient() (*metricsclient.Clientset, error)
	
	CreateComplianceClient() (*client.SonobuoyClient, error)
	
	CreateKnativeBuildClient() (buildclient.Interface, string, error)
	
	CreateTable(out io.Writer) table.Table
	
	SetBatch(batch bool)
	
	IsInCluster() bool
	
	IsInCDPipeline() bool
	
	AuthMergePipelineSecrets(config *auth.AuthConfig, secrets *corev1.SecretList, kind string, isCDPipeline bool) error
	
	CreateVaultOperatorClient() (vaultoperatorclient.Interface, error)
}
