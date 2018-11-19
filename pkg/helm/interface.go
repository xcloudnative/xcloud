package helm

// Helmer defines common helm actions used within Jenkins X
//go:generate pegomock generate github.com/xcloudnative/xcloud/pkg/helm Helmer -o mocks/helmer.go
type Helmer interface {
	//SetCWD(dir string)
	HelmBinary() string
	SetHelmBinary(binary string)
	Init(clientOnly bool, serviceAccount string, tillerNamespace string, upgrade bool) error
	//AddRepo(repo string, URL string) error
	//RemoveRepo(repo string) error
	//ListRepos() (map[string]string, error)
	//UpdateRepo() error
	//IsRepoMissing(URL string) (bool, error)
	//RemoveRequirementsLock() error
	//BuildDependency() error
	InstallChart(chart string, releaseName string, ns string, version *string, timeout *int,
		values []string, valueFiles []string) error
	//FetchChart(chart string, version *string, untar bool, untardir string) error
	//UpgradeChart(chart string, releaseName string, ns string, version *string, install bool,
	//	timeout *int, force bool, wait bool, values []string, valueFiles []string) error
	//DeleteRelease(ns string, releaseName string, purge bool) error
	//ListCharts() (string, error)
	//SearchChartVersions(chart string) ([]string, error)
	//FindChart() (string, error)
	//PackageChart() error
	//StatusRelease(ns string, releaseName string) error
	//StatusReleases(ns string) (map[string]string, error)
	//Lint() (string, error)
	//Version(tls bool) (string, error)
	//SearchCharts(filter string) ([]ChartSummary, error)
	SetHost(host string)
	//Env() map[string]string
	//DecryptSecrets(location string) error
}
