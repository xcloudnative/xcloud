package cmd

import (

	"github.com/xcloudnative/xcloud/pkg/pkg/apis/jenkins.io/v1"
	"github.com/xcloudnative/xcloud/pkg/pkg/config"
	"github.com/xcloudnative/xcloud/pkg/pkg/gits"
	"github.com/xcloudnative/xcloud/pkg/pkg/xc/cmd/templates"
)

var (
	env_description = `		
	An Environment maps to a Kubernetes cluster and namespace and is a place that your team's applications can be promoted to via Continous Delivery.

	You can optionally use GitOps to manage the configuration of an Environment by storing all configuration in a Git repository and then only changing it via Pull Requests and CI/CD.

	For more documentation on Environments see: [https://jenkins-x.io/about/features/#environments](https://jenkins-x.io/about/features/#environments)
	`
	create_env_long = templates.LongDesc(`
		Creates a new Environment
        ` + env_description + `
`)

	create_env_example = templates.Examples(`
		# Create a new Environment, prompting for the required data
		jx create env

		# Creates a new Environment passing in the required data on the command line
		jx create env -n prod -l Production --no-gitops --namespace my-prod
	`)
)

// CreateEnvOptions the options for the create env command
type CreateEnvOptions struct {
	CreateOptions

	Options                v1.Environment
	HelmValuesConfig       config.HelmValuesConfig
	PromotionStrategy      string
	NoGitOps               bool
	Prow                   bool
	ForkEnvironmentGitRepo string
	EnvJobCredentials      string
	GitRepositoryOptions   gits.GitRepositoryOptions
	Prefix                 string
	BranchPattern          string
}

// NewCmdCreateEnv creates a command object for the "create" command
// func NewCmdCreateEnv(f Factory, in terminal.FileReader, out terminal.FileWriter, errOut io.Writer) *cobra.Command {
// 	options := &CreateEnvOptions{
// 		HelmValuesConfig: config.HelmValuesConfig{
// 			ExposeController: &config.ExposeController{},
// 		},
// 		CreateOptions: CreateOptions{
// 			CommonOptions: CommonOptions{
// 				Factory: f,
// 				In:      in,
// 				Out:     out,
// 				Err:     errOut,
// 			},
// 		},
// 	}

// 	cmd := &cobra.Command{
// 		Use:     "environment",
// 		Short:   "Create a new Environment which is used to promote your Team's Applications via Continuous Delivery",
// 		Aliases: []string{"env"},
// 		Long:    create_env_long,
// 		Example: create_env_example,
// 		Run: func(cmd *cobra.Command, args []string) {
// 			options.Cmd = cmd
// 			options.Args = args
// 			err := options.Run()
// 			CheckErr(err)
// 		},
// 	}
// 	//addCreateAppFlags(cmd, &options.CreateOptions)

// 	cmd.Flags().StringVarP(&options.Options.Name, kube.OptionName, "n", "", "The Environment resource name. Must follow the Kubernetes name conventions like Services, Namespaces")
// 	cmd.Flags().StringVarP(&options.Options.Spec.Label, "label", "l", "", "The Environment label which is a descriptive string like 'Production' or 'Staging'")

// 	cmd.Flags().StringVarP(&options.Options.Spec.Namespace, kube.OptionNamespace, "s", "", "The Kubernetes namespace for the Environment")
// 	cmd.Flags().StringVarP(&options.Options.Spec.Cluster, "cluster", "c", "", "The Kubernetes cluster for the Environment. If blank and a namespace is specified assumes the current cluster")
// 	cmd.Flags().StringVarP(&options.Options.Spec.Source.URL, "git-url", "g", "", "The Git clone URL for the source code for GitOps based Environments")
// 	cmd.Flags().StringVarP(&options.Options.Spec.Source.Ref, "git-ref", "r", "", "The Git repo reference for the source code for GitOps based Environments")
// 	cmd.Flags().Int32VarP(&options.Options.Spec.Order, "order", "o", 100, "The order weighting of the Environment so that they can be sorted by this order before name")
// 	cmd.Flags().StringVarP(&options.Prefix, "prefix", "", "jx", "Environment repo prefix, your Git repo will be of the form 'environment-$prefix-$envName'")

// 	cmd.Flags().StringVarP(&options.PromotionStrategy, "promotion", "p", "", "The promotion strategy")
// 	cmd.Flags().StringVarP(&options.ForkEnvironmentGitRepo, "fork-git-repo", "f", kube.DefaultEnvironmentGitRepoURL, "The Git repository used as the fork when creating new Environment Git repos")
// 	cmd.Flags().StringVarP(&options.EnvJobCredentials, "env-job-credentials", "", "", "The Jenkins credentials used by the GitOps Job for this environment")
// 	cmd.Flags().StringVarP(&options.BranchPattern, "branches", "", "", "The branch pattern for branches to trigger CI/CD pipelines on the environment Git repository")

// 	cmd.Flags().BoolVarP(&options.NoGitOps, "no-gitops", "x", false, "Disables the use of GitOps on the environment so that promotion is implemented by directly modifying the resources via helm instead of using a Git repository")
// 	cmd.Flags().BoolVarP(&options.Prow, "prow", "", false, "Install and use Prow for environment promotion")

// 	addGitRepoOptionsArguments(cmd, &options.GitRepositoryOptions)
// 	options.HelmValuesConfig.AddExposeControllerValues(cmd, false)

// 	options.addCommonFlags(cmd)

// 	return cmd
// }

// // Run implements the command
// func (o *CreateEnvOptions) Run() error {
// 	args := o.Args
// 	if len(args) > 0 && o.Options.Name == "" {
// 		o.Options.Name = args[0]
// 	}
// 	//_, currentNs, err := o.JXClientAndDevNamespace()
// 	jxClient, currentNs, err := o.JXClientAndDevNamespace()
// 	if err != nil {
// 		return err
// 	}
// 	kubeClient, _, err := o.KubeClient()
// 	if err != nil {
// 		return err
// 	}
// 	apisClient, err := o.CreateApiExtensionsClient()
// 	if err != nil {
// 		return err
// 	}
// 	kube.RegisterEnvironmentCRD(apisClient)

// 	//_, _, err = kube.GetDevNamespace(kubeClient, currentNs)
// 	ns, _, err := kube.GetDevNamespace(kubeClient, currentNs)
// 	if err != nil {
// 		return err
// 	}
// 	_, err = util.EnvironmentsDir()
// 	envDir, err := util.EnvironmentsDir()
// 	if err != nil {
// 		return err
// 	}
// 	//_, err = o.CreateGitAuthConfigService()
// 	authConfigSvc, err := o.CreateGitAuthConfigService()
// 	if err != nil {
// 		return err
// 	}
// 	devEnv, err := kube.EnsureDevEnvironmentSetup(jxClient, ns)
// 	//_, err = kube.EnsureDevEnvironmentSetup(jxClient, ns)
// 	if err != nil {
// 		return err
// 	}

// 	prowFlag, err := o.isProw()
// 	if err != nil {
// 		return err
// 	}
// 	if prowFlag && !o.Prow {
// 		o.Prow = true
// 	}
// 	if o.Prow {
// 		devEnv.Spec.TeamSettings.PromotionEngine = v1.PromotionEngineProw
// 		devEnv, err = jxClient.JenkinsV1().Environments(ns).Update(devEnv)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	env := v1.Environment{}
// 	o.Options.Spec.PromotionStrategy = v1.PromotionStrategyType(o.PromotionStrategy)
// 	gitProvider, err := kube.CreateEnvironmentSurvey(o.BatchMode, authConfigSvc, devEnv, &env, &o.Options, o.ForkEnvironmentGitRepo, ns,
// 		jxClient, kubeClient, envDir, &o.GitRepositoryOptions, o.HelmValuesConfig, o.Prefix, o.Git(), o.In, o.Out, o.Err)
// 	if err != nil {
// 		return err
// 	}
// 	_, err = jxClient.JenkinsV1().Environments(ns).Create(&env)
// 	if err != nil {
// 		return err
// 	}

// 	log.Infof("Created environment %s\n", util.ColorInfo(env.Name))

// 	err = kube.EnsureEnvironmentNamespaceSetup(kubeClient, jxClient, &env, ns)
// 	if err != nil {
// 		return err
// 	}
// 	gitURL := env.Spec.Source.URL
// 	gitInfo, err := gits.ParseGitURL(gitURL)
// 	if err != nil {
// 		return err
// 	}
// 	if o.Prow {
// 		repo := fmt.Sprintf("%s/environment-%s-%s", gitInfo.Organisation, o.Prefix, o.Options.Name)
// 		err = prow.AddEnvironment(o.KubeClientCached, []string{repo}, devEnv.Spec.Namespace, env.Spec.Namespace)
// 		if err != nil {
// 			return fmt.Errorf("failed to add repo %s to Prow config in namespace %s: %v", repo, env.Spec.Namespace, err)
// 		}
// 	}
// 	/* It is important this pull secret handling goes after any namespace creation code; the service account exists in the created namespace */

// 	if o.PullSecrets != "" {
// 		// We need the namespace to be created first - do the check
// 		err = kube.EnsureEnvironmentNamespaceSetup(kubeClient, jxClient, &env, env.Spec.Namespace)
// 		if err != nil {
// 			// This can happen if, for whatever reason, the namespace takes a while to create. That shouldn't stop the entire process though
// 			log.Warnf("Namespace %s does not exist for jx to patch the service account for, you should patch the service account manually with your pull secret(s) \n", env.Spec.Namespace)
// 		}
// 		// It's a common option, see addCommonFlags in common.go
// 		imagePullSecrets := o.GetImagePullSecrets()
// 		saName := "default"
// 		//log.Infof("Patching the secrets %s for the service account %s\n", imagePullSecrets, saName)
// 		err = serviceaccount.PatchImagePullSecrets(kubeClient, env.Spec.Namespace, saName, imagePullSecrets)
// 		if err != nil {
// 			return fmt.Errorf("Failed to add pull secrets %s to service account %s in namespace %s: %v", imagePullSecrets, saName, env.Spec.Namespace, err)
// 		} else {
// 			log.Infof("Service account \"%s\" in namespace \"%s\" configured to use pull secret(s) %s \n", saName, env.Spec.Namespace, imagePullSecrets)
// 			log.Infof("Pull secret(s) must exist in namespace %s before deploying your applications in this environment \n", env.Spec.Namespace)
// 		}
// 	}

// 	if gitURL != "" {
// 		if gitProvider == nil {
// 			authConfigSvc, err := o.CreateGitAuthConfigService()
// 			if err != nil {
// 				return err
// 			}
// 			gitKind, err := o.GitServerKind(gitInfo)
// 			if err != nil {
// 				return err
// 			}
// 			message := "user name to create the Git repository"
// 			p, err := o.CreateOptions.CommonOptions.Factory.CreateGitProvider(gitURL, message, authConfigSvc, gitKind, o.BatchMode, o.Git(), o.In, o.Out, o.Err)
// 			if err != nil {
// 				return err
// 			}
// 			gitProvider = p
// 		}
// 		if o.Prow {
// 			config := authConfigSvc.Config()
// 			u := gitInfo.HostURL()
// 			server := config.GetOrCreateServer(u)
// 			if len(server.Users) == 0 {
// 				// lets check if the host was used in `~/.jx/gitAuth.yaml` instead of URL
// 				s2 := config.GetOrCreateServer(gitInfo.Host)
// 				if s2 != nil && len(s2.Users) > 0 {
// 					server = s2
// 					u = gitInfo.Host
// 				}
// 			}
// 			user, err := config.PickServerUserAuth(server, "user name for the Pipeline", o.BatchMode, "", o.In, o.Out, o.Err)
// 			if err != nil {
// 				return err
// 			}
// 			if user.Username == "" {
// 				return fmt.Errorf("Could not find a username for git server %s", u)
// 			}
// 			_, err = o.updatePipelineGitCredentialsSecret(server, user)
// 			if err != nil {
// 				return err
// 			}
// 			// register the webhook
// 			return o.createWebhookProw(gitURL, gitProvider)
// 		}
// 		return o.ImportProject(gitURL, envDir, jenkins.DefaultJenkinsfile, o.BranchPattern, o.EnvJobCredentials, false, gitProvider, authConfigSvc, true, o.BatchMode)
// 	}

// 	return nil
// }
