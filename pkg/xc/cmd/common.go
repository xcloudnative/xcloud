package cmd

import (
	"fmt"
	vaultoperatorclient "github.com/banzaicloud/bank-vaults/operator/pkg/client/clientset/versioned"
	"github.com/jenkins-x/golang-jenkins"
	"github.com/xcloudnative/xcloud/pkg/kube"
	"time"

	buildclient "github.com/knative/build/pkg/client/clientset/versioned"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/xcloudnative/xcloud/pkg/client/clientset/versioned"
	"github.com/xcloudnative/xcloud/pkg/helm"
	"github.com/xcloudnative/xcloud/pkg/auth"
	"github.com/xcloudnative/xcloud/pkg/log"
	"github.com/xcloudnative/xcloud/pkg/util"
	"github.com/xcloudnative/xcloud/pkg/gits"
	"gopkg.in/AlecAivazis/survey.v1/terminal"
	"io"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/kubernetes"
	"os"
)

const (
	optionServerName        = "name"
	optionServerURL         = "url"
	optionBatchMode         = "batch-mode"
	optionVerbose           = "verbose"
	optionLogLevel          = "log-level"
	optionHeadless          = "headless"
	optionNoBrew            = "no-brew"
	optionInstallDeps       = "install-dependencies"
	optionSkipAuthSecMerge  = "skip-auth-secrets-merge"
	optionPullSecrets       = "pull-secrets"
	exposecontrollerVersion = "2.3.82"
	exposecontroller        = "exposecontroller"
	exposecontrollerChart   = "jenkins-x/exposecontroller"
)

// CommonOptions contains common options and helper methods
type CommonOptions struct {
	Factory                Factory
	In                     terminal.FileReader
	Out                    terminal.FileWriter
	Err                    io.Writer
	Cmd                    *cobra.Command
	Args                   []string
	BatchMode              bool
	Verbose                bool
	LogLevel               string
	Headless               bool
	NoBrew                 bool
	InstallDependencies    bool
	SkipAuthSecretsMerge   bool
	ServiceAccount         string
	Username               string
	ExternalJenkinsBaseURL string
	PullSecrets            string

	// common cached clients
	KubeClientCached    kubernetes.Interface
	apiExtensionsClient apiextensionsclientset.Interface
	currentNamespace    string
	devNamespace        string
	jxClient            versioned.Interface
	knbClient           buildclient.Interface
	jenkinsClient       gojenkins.JenkinsClient
	GitClient           gits.Gitter
	helm                helm.Helmer
	Kuber               kube.Kuber
	vaultOperatorClient vaultoperatorclient.Interface

	Prow
}

type ServerFlags struct {
	ServerName string
	ServerURL  string
}

func (f *ServerFlags) IsEmpty() bool {
	return f.ServerName == "" && f.ServerURL == ""
}

//func (c *CommonOptions) CreateTable() table.Table {
//	return c.Factory.CreateTable(c.Out)
//}

// NewCommonOptions a helper method to create a new CommonOptions instance
// pre configured in a specific devNamespace
func NewCommonOptions(devNamespace string, factory Factory) CommonOptions {
	return CommonOptions{
		Factory:          factory,
		Out:              os.Stdout,
		Err:              os.Stderr,
		currentNamespace: devNamespace,
		devNamespace:     devNamespace,
	}
}

// SetDevNamespace configures the current dev namespace
func (c *CommonOptions) SetDevNamespace(ns string) {
	c.devNamespace = ns
	c.currentNamespace = ns
	c.KubeClientCached = nil
}

// Debugf outputs the given text to the console if verbose mode is enabled
func (c *CommonOptions) Debugf(format string, a ...interface{}) {
	if c.Verbose {
		log.Infof(format, a...)
	}
}

func (options *CommonOptions) addCommonFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&options.BatchMode, optionBatchMode, "b", false, "In batch mode the command never prompts for user input")
	cmd.Flags().BoolVarP(&options.Verbose, optionVerbose, "", false, "Enable verbose logging")
	cmd.Flags().StringVarP(&options.LogLevel, optionLogLevel, "", logrus.InfoLevel.String(), "Logging level. Possible values - panic, fatal, error, warning, info, debug.")
	cmd.Flags().BoolVarP(&options.Headless, optionHeadless, "", false, "Enable headless operation if using browser automation")
	cmd.Flags().BoolVarP(&options.NoBrew, optionNoBrew, "", false, "Disables the use of brew on macOS to install or upgrade command line dependencies")
	cmd.Flags().BoolVarP(&options.InstallDependencies, optionInstallDeps, "", false, "Should any required dependencies be installed automatically")
	cmd.Flags().BoolVarP(&options.SkipAuthSecretsMerge, optionSkipAuthSecMerge, "", false, "Skips merging a local git auth yaml file with any pipeline secrets that are found")
	cmd.Flags().StringVarP(&options.PullSecrets, optionPullSecrets, "", "", "The pull secrets the service account created should have (useful when deploying to your own private registry): provide multiple pull secrets by providing them in a singular block of quotes e.g. --pull-secrets \"foo, bar, baz\"")

	options.Cmd = cmd
}

func (o *CommonOptions) CreateApiExtensionsClient() (apiextensionsclientset.Interface, error) {
	var err error
	if o.apiExtensionsClient == nil {
		o.apiExtensionsClient, err = o.Factory.CreateApiExtensionsClient()
		if err != nil {
			return nil, err
		}
	}
	return o.apiExtensionsClient, nil
}

func (o *CommonOptions) KubeClient() (kubernetes.Interface, string, error) {
	if o.KubeClientCached == nil {
		kubeClient, currentNs, err := o.Factory.CreateClient()
		if err != nil {
			return nil, "", err
		}
		o.KubeClientCached = kubeClient
		o.currentNamespace = currentNs

	}
	return o.KubeClientCached, o.currentNamespace, nil
}

//// KubeClientAndDevNamespace returns a kube client and the development namespace
//func (o *CommonOptions) KubeClientAndDevNamespace() (kubernetes.Interface, string, error) {
//	kubeClient, curNs, err := o.KubeClient()
//	if err != nil {
//		return nil, "", err
//	}
//	if o.devNamespace == "" {
//		o.devNamespace, _, err = kube.GetDevNamespace(kubeClient, curNs)
//	}
//	return kubeClient, o.devNamespace, err
//}

func (o *CommonOptions) JXClient() (versioned.Interface, string, error) {
	if o.Factory == nil {
		return nil, "", errors.New("command factory is not initialized")
	}
	if o.jxClient == nil {
		jxClient, ns, err := o.Factory.CreateJXClient()
		if err != nil {
			return nil, ns, err
		}
		o.jxClient = jxClient
		if o.currentNamespace == "" {
			o.currentNamespace = ns
		}
	}
	return o.jxClient, o.currentNamespace, nil
}

//func (o *CommonOptions) KnativeBuildClient() (buildclient.Interface, string, error) {
//	if o.Factory == nil {
//		return nil, "", errors.New("command factory is not initialized")
//	}
//	if o.knbClient == nil {
//		knbClient, ns, err := o.Factory.CreateKnativeBuildClient()
//		if err != nil {
//			return nil, ns, err
//		}
//		o.knbClient = knbClient
//		if o.currentNamespace == "" {
//			o.currentNamespace = ns
//		}
//	}
//	return o.knbClient, o.currentNamespace, nil
//}
//
//func (o *CommonOptions) JXClientAndAdminNamespace() (versioned.Interface, string, error) {
//	kubeClient, _, err := o.KubeClient()
//	if err != nil {
//		return nil, "", err
//	}
//	jxClient, devNs, err := o.JXClientAndDevNamespace()
//	if err != nil {
//		return nil, "", err
//	}
//
//	ns, err := kube.GetAdminNamespace(kubeClient, devNs)
//	return jxClient, ns, err
//}

func (o *CommonOptions) JXClientAndDevNamespace() (versioned.Interface, string, error) {
	if o.jxClient == nil {
		jxClient, ns, err := o.JXClient()
		if err != nil {
			return nil, ns, err
		}
		o.jxClient = jxClient
		if o.currentNamespace == "" {
			o.currentNamespace = ns
		}
	}
	if o.devNamespace == "" {
		client, ns, err := o.KubeClient()
		if err != nil {
			return nil, "", err
		}
		devNs, _, err := kube.GetDevNamespace(client, ns)
		if err != nil {
			return nil, "", err
		}
		o.devNamespace = devNs
	}
	return o.jxClient, o.devNamespace, nil
}

//func (o *CommonOptions) JenkinsClient() (gojenkins.JenkinsClient, error) {
//	if o.jenkinsClient == nil {
//		kubeClient, ns, err := o.KubeClientAndDevNamespace()
//		if err != nil {
//			return nil, err
//		}
//
//		jenkins, err := o.Factory.CreateJenkinsClient(kubeClient, ns, o.In, o.Out, o.Err)
//
//		if err != nil {
//			return nil, err
//		}
//		o.jenkinsClient = jenkins
//	}
//	return o.jenkinsClient, nil
//}
//func (o *CommonOptions) GetJenkinsURL() (string, error) {
//	kubeClient, ns, err := o.KubeClient()
//	if err != nil {
//		return "", err
//	}
//
//	return o.Factory.GetJenkinsURL(kubeClient, ns)
//}

func (o *CommonOptions) Git() gits.Gitter {
	if o.GitClient == nil {
		o.GitClient = gits.NewGitCLI()
	}
	return o.GitClient
}

func (o *CommonOptions) Helm() helm.Helmer {
	if o.helm == nil {
		helmBinary, noTiller, helmTemplate, err := o.TeamHelmBin()
		if err != nil {
			helmBinary = defaultHelmBin
		}
		featureFlag := "none"
		if helmTemplate {
			featureFlag = "template-mode"
		} else if noTiller {
			featureFlag = "no-tiller-server"
		}
		log.Infof("Using helmBinary %s with feature flag: %s\n", util.ColorInfo(helmBinary), util.ColorInfo(featureFlag))
		helmCLI := helm.NewHelmCLI(helmBinary, helm.V2, "", o.Verbose)
		o.helm = helmCLI
		if helmTemplate {
			kubeClient, ns, _ := o.KubeClient()
			o.helm = helm.NewHelmTemplate(helmCLI, "", kubeClient, ns)
		} else {
			o.helm = helmCLI
		}
		if noTiller {
			o.helm.SetHost(o.tillerAddress())
			o.startLocalTillerIfNotRunning()
		}
	}
	return o.helm
}

func (o *CommonOptions) Kube() kube.Kuber {
	if o.Kuber == nil {
		o.Kuber = kube.NewKubeConfig()
	}
	return o.Kuber
}
//
//func (o *CommonOptions) TeamAndEnvironmentNames() (string, string, error) {
//	kubeClient, currentNs, err := o.KubeClient()
//	if err != nil {
//		return "", "", err
//	}
//	return kube.GetDevNamespace(kubeClient, currentNs)
//}
//
//func (o *CommonOptions) GetImagePullSecrets() []string {
//	pullSecrets := strings.Fields(o.PullSecrets)
//	return pullSecrets
//}
//
//func (o *ServerFlags) addGitServerFlags(cmd *cobra.Command) {
//	cmd.Flags().StringVarP(&o.ServerName, optionServerName, "n", "", "The name of the Git server to add a user")
//	cmd.Flags().StringVarP(&o.ServerURL, optionServerURL, "u", "", "The URL of the Git server to add a user")
//}
//
//// findGitServer finds the Git server from the given flags or returns an error
//func (o *CommonOptions) findGitServer(config *auth.AuthConfig, serverFlags *ServerFlags) (*auth.AuthServer, error) {
//	return o.findServer(config, serverFlags, "git", "Try creating one via: jx create git server", false)
//}
//
//// findIssueTrackerServer finds the issue tracker server from the given flags or returns an error
//func (o *CommonOptions) findIssueTrackerServer(config *auth.AuthConfig, serverFlags *ServerFlags) (*auth.AuthServer, error) {
//	return o.findServer(config, serverFlags, "issues", "Try creating one via: jx create tracker server", false)
//}
//
//// findChatServer finds the chat server from the given flags or returns an error
//func (o *CommonOptions) findChatServer(config *auth.AuthConfig, serverFlags *ServerFlags) (*auth.AuthServer, error) {
//	return o.findServer(config, serverFlags, "chat", "Try creating one via: jx create chat server", false)
//}
//
//// findAddonServer finds the addon server from the given flags or returns an error
//func (o *CommonOptions) findAddonServer(config *auth.AuthConfig, serverFlags *ServerFlags, kind string) (*auth.AuthServer, error) {
//	return o.findServer(config, serverFlags, kind, "Try creating one via: jx create addon", true)
//}

func (o *CommonOptions) findServer(config *auth.AuthConfig, serverFlags *ServerFlags, defaultKind string, missingServerDescription string, lazyCreate bool) (*auth.AuthServer, error) {
	kind := defaultKind
	var server *auth.AuthServer
	if serverFlags.ServerURL != "" {
		server = config.GetServer(serverFlags.ServerURL)
		if server == nil {
			if lazyCreate {
				return config.GetOrCreateServerName(serverFlags.ServerURL, serverFlags.ServerName, kind), nil
			}
			return nil, util.InvalidOption(optionServerURL, serverFlags.ServerURL, config.GetServerURLs())
		}
	}
	if server == nil && serverFlags.ServerName != "" {
		name := serverFlags.ServerName
		if lazyCreate {
			server = config.GetOrCreateServerName(serverFlags.ServerURL, name, kind)
		} else {
			server = config.GetServerByName(name)
		}
		if server == nil {
			return nil, util.InvalidOption(optionServerName, name, config.GetServerNames())
		}
	}
	if server == nil {
		name := config.CurrentServer
		if name != "" && o.BatchMode {
			server = config.GetServerByName(name)
			if server == nil {
				log.Warnf("Current server %s no longer exists\n", name)
			}
		}
	}
	if server == nil && len(config.Servers) == 1 {
		server = config.Servers[0]
	}
	if server == nil && len(config.Servers) > 1 {
		if o.BatchMode {
			return nil, fmt.Errorf("Multiple servers found. Please specify one via the %s option", optionServerName)
		}
		defaultServerName := ""
		if config.CurrentServer != "" {
			s := config.GetServer(config.CurrentServer)
			if s != nil {
				defaultServerName = s.Name
			}
		}
		name, err := util.PickNameWithDefault(config.GetServerNames(), "Pick server to use: ", defaultServerName, "", o.In, o.Out, o.Err)
		if err != nil {
			return nil, err
		}
		server = config.GetServerByName(name)
		if server == nil {
			return nil, fmt.Errorf("Could not find the server for name %s", name)
		}
	}
	if server == nil {
		return nil, fmt.Errorf("Could not find a %s. %s", kind, missingServerDescription)
	}
	return server, nil
}

//func (o *CommonOptions) findService(name string) (string, error) {
//	client, ns, err := o.KubeClient()
//	if err != nil {
//		return "", err
//	}
//	devNs, _, err := kube.GetDevNamespace(client, ns)
//	if err != nil {
//		return "", err
//	}
//	url, err := services.FindServiceURL(client, ns, name)
//	if url == "" {
//		url, err = services.FindServiceURL(client, devNs, name)
//	}
//	if url == "" {
//		names, err := services.GetServiceNames(client, ns, name)
//		if err != nil {
//			return "", err
//		}
//		if len(names) > 1 {
//			name, err = util.PickName(names, "Pick service to open: ", "", o.In, o.Out, o.Err)
//			if err != nil {
//				return "", err
//			}
//			if name != "" {
//				url, err = services.FindServiceURL(client, ns, name)
//			}
//		} else if len(names) == 1 {
//			// must have been a filter
//			url, err = services.FindServiceURL(client, ns, names[0])
//		}
//		if url == "" {
//			return "", fmt.Errorf("Could not find URL for service %s in namespace %s", name, ns)
//		}
//	}
//	return url, nil
//}
//
//func (o *CommonOptions) findEnvironmentNamespace(envName string) (string, error) {
//	client, ns, err := o.KubeClient()
//	if err != nil {
//		return "", err
//	}
//	jxClient, _, err := o.JXClient()
//	if err != nil {
//		return "", err
//	}
//
//	devNs, _, err := kube.GetDevNamespace(client, ns)
//	if err != nil {
//		return "", err
//	}
//
//	envMap, envNames, err := kube.GetEnvironments(jxClient, devNs)
//	if err != nil {
//		return "", err
//	}
//	env := envMap[envName]
//	if env == nil {
//		return "", util.InvalidOption(optionEnvironment, envName, envNames)
//	}
//	answer := env.Spec.Namespace
//	if answer == "" {
//		return "", fmt.Errorf("Environment %s does not have a Namespace!", envName)
//	}
//	return answer, nil
//}
//
//func (o *CommonOptions) findServiceInNamespace(name string, ns string) (string, error) {
//	client, curNs, err := o.KubeClient()
//	if err != nil {
//		return "", err
//	}
//	if ns == "" {
//		ns = curNs
//	}
//	url, err := services.FindServiceURL(client, ns, name)
//	if url == "" {
//		names, err := services.GetServiceNames(client, ns, name)
//		if err != nil {
//			return "", err
//		}
//		if len(names) > 1 {
//			name, err = util.PickName(names, "Pick service to open: ", "", o.In, o.Out, o.Err)
//			if err != nil {
//				return "", err
//			}
//			if name != "" {
//				url, err = services.FindServiceURL(client, ns, name)
//			}
//		} else if len(names) == 1 {
//			// must have been a filter
//			url, err = services.FindServiceURL(client, ns, names[0])
//		}
//		if url == "" {
//			return "", fmt.Errorf("Could not find URL for service %s in namespace %s", name, ns)
//		}
//	}
//	return url, nil
//}

func (o *CommonOptions) retry(attempts int, sleep time.Duration, call func() error) (err error) {
	for i := 0; ; i++ {
		err = call()
		if err == nil {
			return
		}

		if i >= (attempts - 1) {
			break
		}

		time.Sleep(sleep)

		log.Infof("retrying after error:%s\n", err)
	}
	return fmt.Errorf("after %d attempts, last error: %s", attempts, err)
}

func (o *CommonOptions) retryQuiet(attempts int, sleep time.Duration, call func() error) (err error) {
	lastMessage := ""
	dot := false

	for i := 0; ; i++ {
		err = call()
		if err == nil {
			if dot {
				log.Blank()
			}
			return
		}

		if i >= (attempts - 1) {
			break
		}

		time.Sleep(sleep)

		message := fmt.Sprintf("retrying after error: %s", err)
		if lastMessage == message {
			log.Info(".")
			dot = true
		} else {
			lastMessage = message
			if dot {
				dot = false
				log.Blank()
			}
			log.Infof("%s\n", lastMessage)
		}
	}
	return fmt.Errorf("after %d attempts, last error: %s", attempts, err)
}

//func (o *CommonOptions) retryQuietlyUntilTimeout(timeout time.Duration, sleep time.Duration, call func() error) (err error) {
//	timeoutTime := time.Now().Add(timeout)
//
//	lastMessage := ""
//	dot := false
//
//	for i := 0; ; i++ {
//		err = call()
//		if err == nil {
//			if dot {
//				log.Blank()
//			}
//			return
//		}
//
//		if time.Now().After(timeoutTime) {
//			return fmt.Errorf("Timed out after %s, last error: %s", timeout.String(), err)
//		}
//
//		time.Sleep(sleep)
//
//		message := fmt.Sprintf("retrying after error: %s", err)
//		if lastMessage == message {
//			log.Info(".")
//			dot = true
//		} else {
//			lastMessage = message
//			if dot {
//				dot = false
//				log.Blank()
//			}
//			log.Infof("%s\n", lastMessage)
//		}
//	}
//}
//
//// retryUntilTrueOrTimeout waits until complete is true, an error occurs or the timeout
//func (o *CommonOptions) retryUntilTrueOrTimeout(timeout time.Duration, sleep time.Duration, call func() (bool, error)) (err error) {
//	timeoutTime := time.Now().Add(timeout)
//
//	for i := 0; ; i++ {
//		complete, err := call()
//		if complete || err != nil {
//			return err
//		}
//		if time.Now().After(timeoutTime) {
//			return fmt.Errorf("Timed out after %s, last error: %s", timeout.String(), err)
//		}
//
//		time.Sleep(sleep)
//	}
//}
//
//func (o *CommonOptions) getJobMap(filter string) (map[string]gojenkins.Job, error) {
//	jobMap := map[string]gojenkins.Job{}
//	jenkins, err := o.JenkinsClient()
//	if err != nil {
//		return jobMap, err
//	}
//	jobs, err := jenkins.GetJobs()
//	if err != nil {
//		return jobMap, err
//	}
//	o.addJobs(&jobMap, filter, "", jobs)
//	return jobMap, nil
//}
//
//func (o *CommonOptions) addJobs(jobMap *map[string]gojenkins.Job, filter string, prefix string, jobs []gojenkins.Job) {
//	jenkins, err := o.JenkinsClient()
//	if err != nil {
//		return
//	}
//
//	for _, j := range jobs {
//		name := jobName(prefix, &j)
//		if IsPipeline(&j) {
//			if filter == "" || strings.Contains(name, filter) {
//				(*jobMap)[name] = j
//				continue
//			}
//		}
//		if j.Jobs != nil {
//			o.addJobs(jobMap, filter, name, j.Jobs)
//		} else {
//			job, err := jenkins.GetJob(name)
//			if err == nil && job.Jobs != nil {
//				o.addJobs(jobMap, filter, name, job.Jobs)
//			}
//		}
//	}
//}
//func (o *CommonOptions) tailBuild(jobName string, build *gojenkins.Build) error {
//	jenkins, err := o.JenkinsClient()
//	if err != nil {
//		return nil
//	}
//
//	u, err := url.Parse(build.Url)
//	if err != nil {
//		return err
//	}
//	buildPath := u.Path
//	log.Infof("%s %s\n", "tailing the log of", fmt.Sprintf("%s #%d", jobName, build.Number))
//	// TODO Logger
//	return jenkins.TailLog(buildPath, o.Out, time.Second, time.Hour*100)
//}
//
//func (o *CommonOptions) pickRemoteURL(config *gitcfg.Config) (string, error) {
//	surveyOpts := survey.WithStdio(o.In, o.Out, o.Err)
//	urls := []string{}
//	if config.Remotes != nil {
//		for _, r := range config.Remotes {
//			if r.URLs != nil {
//				for _, u := range r.URLs {
//					urls = append(urls, u)
//				}
//			}
//		}
//	}
//	if len(urls) == 1 {
//		return urls[0], nil
//	}
//	url := ""
//	if len(urls) > 1 {
//		prompt := &survey.Select{
//			Message: "Choose a remote git URL:",
//			Options: urls,
//		}
//		err := survey.AskOne(prompt, &url, nil, surveyOpts)
//		if err != nil {
//			return "", err
//		}
//	}
//	return url, nil
//}
//
//// todo switch to using exposecontroller as a jx plugin
//// get existing config from the devNamespace and run exposecontroller in the target environment
//func (o *CommonOptions) expose(devNamespace, targetNamespace, password string) error {
//
//	_, err := o.KubeClientCached.CoreV1().Secrets(targetNamespace).Get(kube.SecretBasicAuth, v1.GetOptions{})
//	if err != nil {
//		data := make(map[string][]byte)
//
//		if password != "" {
//			hash := config.HashSha(password)
//			data[kube.AUTH] = []byte(fmt.Sprintf("admin:{SHA}%s", hash))
//		} else {
//			basicAuth, err := o.KubeClientCached.CoreV1().Secrets(devNamespace).Get(kube.SecretBasicAuth, v1.GetOptions{})
//			if err != nil {
//				return fmt.Errorf("cannot find secret %s in namespace %s: %v", kube.SecretBasicAuth, devNamespace, err)
//			}
//			data = basicAuth.Data
//		}
//
//		sec := &core_v1.Secret{
//			Data: data,
//			ObjectMeta: v1.ObjectMeta{
//				Name: kube.SecretBasicAuth,
//			},
//		}
//		_, err := o.KubeClientCached.CoreV1().Secrets(targetNamespace).Create(sec)
//		if err != nil {
//			return fmt.Errorf("cannot create secret %s in target namespace %s: %v", kube.SecretBasicAuth, targetNamespace, err)
//		}
//	}
//
//	ic, err := kube.GetIngressConfig(o.KubeClientCached, devNamespace)
//	if err != nil {
//		return fmt.Errorf("cannot get existing team exposecontroller config from namespace %s: %v", devNamespace, err)
//	}
//
//	err = services.AnnotateNamespaceServicesWithCertManager(o.KubeClientCached, targetNamespace, ic.Issuer)
//	if err != nil {
//		return err
//	}
//
//	// if targetnamespace is different than dev check if there's any certmanager CRDs, if not check dev and copy any found across
//	err = o.copyCertmanagerResources(targetNamespace, ic)
//	if err != nil {
//		return fmt.Errorf("failed to copy certmanager resources from %s to %s namespace: %v", devNamespace, targetNamespace, err)
//	}
//
//	return o.runExposecontroller(devNamespace, targetNamespace, ic)
//}
//
//func (o *CommonOptions) exposeService(service, devNamespace, targetNamespace string) error {
//	ic, err := kube.GetIngressConfig(o.KubeClientCached, devNamespace)
//	if err != nil {
//		return fmt.Errorf("cannot get existing team exposecontroller config from namespace %s: %v", devNamespace, err)
//	}
//	err = services.AnnotateNamespaceServicesWithCertManager(o.KubeClientCached, targetNamespace, ic.Issuer, service)
//	if err != nil {
//		return err
//	}
//
//	err = o.copyCertmanagerResources(targetNamespace, ic)
//	if err != nil {
//		return fmt.Errorf("failed to copy certmanager resources from %s to %s namespace: %v", devNamespace, targetNamespace, err)
//	}
//
//	return o.runExposecontroller(devNamespace, targetNamespace, ic, service)
//}
//
//func (o *CommonOptions) runExposecontroller(devNamespace, targetNamespace string, ic kube.IngressConfig, services ...string) error {
//
//	o.CleanExposecontrollerReources(targetNamespace)
//
//	exValues := []string{
//		"config.exposer=" + ic.Exposer,
//		"config.domain=" + ic.Domain,
//		"config.tlsacme=" + strconv.FormatBool(ic.TLS),
//	}
//
//	if !ic.TLS && ic.Issuer != "" {
//		exValues = append(exValues, "config.http=true")
//	}
//
//	if len(services) > 0 {
//		serviceCfg := "config.extravalues.services={"
//		for i, service := range services {
//			if i > 0 {
//				serviceCfg += ","
//			}
//			serviceCfg += service
//		}
//		serviceCfg += "}"
//		exValues = append(exValues, serviceCfg)
//	}
//
//	helmRelease := "expose-" + strings.ToLower(randomdata.SillyName())
//	err := o.installChartOptions(InstallChartOptions{
//		ReleaseName: helmRelease,
//		Chart:       exposecontrollerChart,
//		Version:     exposecontrollerVersion,
//		Ns:          targetNamespace,
//		HelmUpdate:  true,
//		SetValues:   exValues,
//	})
//	if err != nil {
//		return fmt.Errorf("exposecontroller deployment failed: %v", err)
//	}
//	err = kube.WaitForJobToSucceeded(o.KubeClientCached, targetNamespace, exposecontroller, 5*time.Minute)
//	if err != nil {
//		return fmt.Errorf("failed waiting for exposecontroller job to succeed: %v", err)
//	}
//	return o.helm.DeleteRelease(targetNamespace, helmRelease, true)
//
//}
//
//// CleanExposecontrollerReources cleans expose controller resources
//func (o *CommonOptions) CleanExposecontrollerReources(ns string) {
//
//	// let's not error if nothing to cleanup
//	o.KubeClientCached.RbacV1().Roles(ns).Delete(exposecontroller, &metav1.DeleteOptions{})
//	o.KubeClientCached.RbacV1().RoleBindings(ns).Delete(exposecontroller, &metav1.DeleteOptions{})
//	o.KubeClientCached.RbacV1().ClusterRoleBindings().Delete(exposecontroller, &metav1.DeleteOptions{})
//	o.KubeClientCached.CoreV1().ConfigMaps(ns).Delete(exposecontroller, &metav1.DeleteOptions{})
//	o.KubeClientCached.CoreV1().ServiceAccounts(ns).Delete(exposecontroller, &metav1.DeleteOptions{})
//	o.KubeClientCached.BatchV1().Jobs(ns).Delete(exposecontroller, &metav1.DeleteOptions{})
//
//}
//
//func (o *CommonOptions) getDefaultAdminPassword(devNamespace string) (string, error) {
//	basicAuth, err := o.KubeClientCached.CoreV1().Secrets(devNamespace).Get(JXInstallConfig, v1.GetOptions{})
//	if err != nil {
//		return "", fmt.Errorf("cannot find secret %s in namespace %s: %v", kube.SecretBasicAuth, devNamespace, err)
//	}
//	adminSecrets := basicAuth.Data[AdminSecretsFile]
//	adminConfig := config.AdminSecretsConfig{}
//
//	err = yaml.Unmarshal(adminSecrets, &adminConfig)
//	if err != nil {
//		return "", err
//	}
//	return adminConfig.Jenkins.JenkinsSecret.Password, nil
//}
//
//func (o *CommonOptions) ensureAddonServiceAvailable(serviceName string) (string, error) {
//	present, err := services.IsServicePresent(o.KubeClientCached, serviceName, o.currentNamespace)
//	if err != nil {
//		return "", fmt.Errorf("no %s provider service found, are you in your teams dev environment?  Type `jx ns` to switch.", serviceName)
//	}
//	if present {
//		url, err := services.GetServiceURLFromName(o.KubeClientCached, serviceName, o.currentNamespace)
//		if err != nil {
//			return "", fmt.Errorf("no %s provider service found, are you in your teams dev environment?  Type `jx ns` to switch.", serviceName)
//		}
//		return url, nil
//	}
//
//	// todo ask if user wants to install addon?
//	return "", nil
//}
//
//func (o *CommonOptions) copyCertmanagerResources(targetNamespace string, ic kube.IngressConfig) error {
//	if ic.TLS {
//		err := kube.CleanCertmanagerResources(o.KubeClientCached, targetNamespace, ic)
//		if err != nil {
//			return fmt.Errorf("failed to create certmanager resources in target namespace %s: %v", targetNamespace, err)
//		}
//	}
//
//	return nil
//}
//
//func (o *CommonOptions) getJobName() string {
//	owner := os.Getenv("REPO_OWNER")
//	repo := os.Getenv("REPO_NAME")
//	branch := os.Getenv("BRANCH_NAME")
//
//	if owner != "" && repo != "" && branch != "" {
//		return fmt.Sprintf("%s/%s/%s", owner, repo, branch)
//	}
//
//	job := os.Getenv("JOB_NAME")
//	if job != "" {
//		return job
//	}
//	return ""
//}
//
//func (o *CommonOptions) getBuildNumber() string {
//	buildNumber := os.Getenv("JX_BUILD_NUMBER")
//	if buildNumber != "" {
//		return buildNumber
//	}
//	buildNumber = os.Getenv("BUILD_NUMBER")
//	if buildNumber != "" {
//		return buildNumber
//	}
//	buildID := os.Getenv("BUILD_ID")
//	if buildID != "" {
//		return buildID
//	}
//	return ""
//}
//
//func (o *CommonOptions) VaultOperatorClient() (vaultoperatorclient.Interface, error) {
//	if o.Factory == nil {
//		return nil, errors.New("command factory is not initialized")
//	}
//	if o.vaultOperatorClient == nil {
//		vaultOperatorClient, err := o.Factory.CreateVaultOperatorClient()
//		if err != nil {
//			return nil, err
//		}
//		o.vaultOperatorClient = vaultOperatorClient
//	}
//	return o.vaultOperatorClient, nil
//}
//
//func (o *CommonOptions) GetWebHookEndpoint() (string, error) {
//	_, _, err := o.JXClient()
//	if err != nil {
//		return "", errors.Wrap(err, "failed to get jxclient")
//	}
//
//	_, _, err = o.KubeClient()
//	if err != nil {
//		return "", errors.Wrap(err, "failed to get kube client")
//	}
//
//	isProwEnabled, err := o.isProw()
//	if err != nil {
//		return "", err
//	}
//
//	ns, _, err := kube.GetDevNamespace(o.KubeClientCached, o.currentNamespace)
//	if err != nil {
//		return "", err
//	}
//
//	var webHookUrl string
//
//	if isProwEnabled {
//		baseURL, err := services.GetServiceURLFromName(o.KubeClientCached, "hook", ns)
//		if err != nil {
//			return "", err
//		}
//
//		webHookUrl = util.UrlJoin(baseURL, "hook")
//	} else {
//		baseURL, err := services.GetServiceURLFromName(o.KubeClientCached, "jenkins", ns)
//		if err != nil {
//			return "", err
//		}
//
//		webHookUrl = util.UrlJoin(baseURL, "github-webhook/")
//	}
//
//	return webHookUrl, nil
//}
//
//func (o *CommonOptions) GetIn() terminal.FileReader {
//	return o.In
//}
//
//func (o *CommonOptions) GetOut() terminal.FileWriter {
//	return o.Out
//}
//
//func (o *CommonOptions) GetErr() io.Writer {
//	return o.Err
//}
