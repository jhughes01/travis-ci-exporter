package main

import (
	"bufio"
	"context"
	"flag"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/shuheiktgw/go-travis"

	log "github.com/sirupsen/logrus"
)

var (
	travisToken      string
	travisApiUrl     string
	travisApiVersion string
	logLevel         string
	repoFile         string
)

var travisBuildStatus = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "travis_build_status", Help: "Reports the last status of a Travis CI build job",
})

func main() {
	// rs, err := readRepoSlugsFromFile(repoFile)
	// if err != nil {
	// 	log.Error(err)
	// }
	// for _, r := range rs {
	// 	getBuildStatus(r)
	// }
	prometheus.MustRegister(travisBuildStatus)
	http.Handle("/metrics", promhttp.Handler())
	log.Info("Beginning to serve on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func init() {
	// Expose flags for configuration
	flag.StringVar(&logLevel, "log-level", "Info", "Exporter log level")
	flag.StringVar(&travisToken, "token", "", "Travis API Token")
	flag.StringVar(&travisApiUrl, "api-url", "https://api.travis-ci.org/", "Travis CI API url")
	flag.StringVar(&travisApiVersion, "api-version", "3", "Travis CI API version")
	flag.StringVar(&repoFile, "repo-file", "", "File containing a list of repositories")
	// Parse flags
	flag.Parse()

	// Set log level
	ll, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Error(err)
	}
	log.SetLevel(ll)
	log.SetOutput(os.Stdout)
}

func checkFlags() error {
	return nil
}

func readRepoSlugsFromFile(filename string) (repos []string, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	repos = []string{}
	s := bufio.NewScanner(f)
	for s.Scan() {
		repos = append(repos, s.Text())
	}
	return repos, nil
}

func getBuildStatus(repoSlug string) {
	client := travis.NewClient(travisApiUrl, travisToken)
	builds, _, err := client.Builds.ListByRepoSlug(context.Background(), repoSlug, nil)
	if err != nil {
		log.Error(err)
	}
	latestBuild := builds[0]
	log.Info(*latestBuild.Repository.Slug + " " + *latestBuild.State + " " + *latestBuild.FinishedAt)
}
