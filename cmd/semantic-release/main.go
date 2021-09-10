package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Masterminds/semver/v3"
	"github.com/duanqy/semantic-release/pkg/config"
	"github.com/duanqy/semantic-release/pkg/plugin"
	"github.com/duanqy/semantic-release/pkg/semrel"
	"github.com/spf13/cobra"

	_ "github.com/duanqy/semantic-release/plugin/changelog_generator"
	_ "github.com/duanqy/semantic-release/plugin/commit_analyzer"
	_ "github.com/duanqy/semantic-release/plugin/condition/default_condition"
	_ "github.com/duanqy/semantic-release/plugin/condition/github_condition"
	_ "github.com/duanqy/semantic-release/plugin/condition/gitlab_condition"
	_ "github.com/duanqy/semantic-release/plugin/files_updater_npm"
	_ "github.com/duanqy/semantic-release/plugin/provider/githubprovider"
	_ "github.com/duanqy/semantic-release/plugin/provider/gitlabprovider"
	_ "github.com/duanqy/semantic-release/plugin/provider/gitprovider"
)

// SRVERSION is the semantic-release version (added at compile time)
var SRVERSION string = "dqy"

var exitHandler func()

func errorHandler(logger *log.Logger) func(error, ...int) {
	return func(err error, exitCode ...int) {
		if err != nil {
			logger.Println(err)
			if exitHandler != nil {
				exitHandler()
			}
			if len(exitCode) == 1 {
				os.Exit(exitCode[0])
				return
			}
			os.Exit(1)
		}
	}
}

func main() {
	cmd := &cobra.Command{
		Use:     "semantic-release",
		Short:   "semantic-release - fully automated package/module/image publishing",
		Run:     cliHandler,
		Version: SRVERSION,
	}

	err := config.InitConfig(cmd)
	if err != nil {
		fmt.Printf("\nConfig error: %s\n", err.Error())
		os.Exit(1)
		return
	}
	err = cmd.Execute()
	if err != nil {
		fmt.Printf("\n%s\n", err.Error())
		os.Exit(1)
	}
}

func cliHandler(cmd *cobra.Command, args []string) {
	logger := log.New(os.Stderr, "[go-semantic-release]: ", 0)
	exitIfError := errorHandler(logger)

	logger.Printf("version: %s\n", SRVERSION)

	conf, err := config.NewConfig(cmd)
	exitIfError(err)

	pluginManager, err := plugin.NewManager(conf)
	exitIfError(err)
	exitHandler = func() {
		pluginManager.Stop()
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		exitIfError(errors.New("terminating..."))
	}()

	ci, err := pluginManager.GetCICondition()
	exitIfError(err)
	logger.Printf("ci-condition plugin: %s@%s\n", ci.Name(), ci.Version())

	prov, err := pluginManager.GetProvider()
	exitIfError(err)
	logger.Printf("provider plugin: %s@%s\n", prov.Name(), prov.Version())

	if conf.ProviderOpts["token"] == "" {
		conf.ProviderOpts["token"] = conf.Token
	}
	err = prov.Init(conf.ProviderOpts)
	exitIfError(err)

	logger.Println("getting default branch...")
	repoInfo, err := prov.GetInfo()
	exitIfError(err)
	logger.Println("found default branch: " + repoInfo.DefaultBranch)
	if repoInfo.Private {
		logger.Println("repo is private")
	}

	currentBranch := ci.GetCurrentBranch()
	if currentBranch == "" {
		exitIfError(fmt.Errorf("current branch not found"))
	}
	logger.Println("found current branch: " + currentBranch)

	if conf.MaintainedVersion != "" && currentBranch == repoInfo.DefaultBranch {
		exitIfError(fmt.Errorf("maintained version not allowed on default branch"))
	}

	if conf.MaintainedVersion != "" {
		logger.Println("found maintained version: " + conf.MaintainedVersion)
		repoInfo.DefaultBranch = "*"
	}

	currentSha := ci.GetCurrentSHA()
	logger.Println("found current sha: " + currentSha)

	hooksExecutor, err := pluginManager.GetChainedHooksExecutor()
	exitIfError(err)

	hooksNames := hooksExecutor.GetNameVersionPairs()
	if len(hooksNames) > 0 {
		logger.Printf("hooks plugins: %s\n", strings.Join(hooksNames, ", "))
	}

	exitIfError(hooksExecutor.Init(conf.HooksOpts))

	if !conf.NoCI {
		logger.Println("running CI condition...")
		conditionConfig := map[string]string{
			"token":         conf.Token,
			"defaultBranch": repoInfo.DefaultBranch,
			"private":       fmt.Sprintf("%t", repoInfo.Private),
		}
		for k, v := range conf.CIConditionOpts {
			conditionConfig[k] = v
		}
		err = ci.RunCondition(conditionConfig)
		if err != nil {
			herr := hooksExecutor.NoRelease(plugin.NoReleaseReasonCondition, err.Error())
			if herr != nil {
				logger.Printf("there was an error executing the hooks plugins: %s", herr.Error())
			}
			exitIfError(err, 66)
		}

	}

	logger.Println("getting latest release...")
	matchRegex := ""
	match := strings.TrimSpace(conf.Match)
	if match != "" {
		logger.Printf("getting latest release matching %s...", match)
		matchRegex = "^" + match
	}
	releases, err := prov.GetReleases(matchRegex)
	exitIfError(err)
	release, err := semrel.GetLatestReleaseFromReleases(releases, conf.MaintainedVersion)
	exitIfError(err)
	logger.Println("found version: " + release.Version)

	if strings.Contains(conf.MaintainedVersion, "-") && semver.MustParse(release.Version).Prerelease() == "" {
		exitIfError(fmt.Errorf("no pre-release for this version possible"))
	}

	if release.SHA == currentSha {
		logger.Println("no new commits,write last release version")
		exitIfError(ioutil.WriteFile(".version", []byte(release.Version), 0644))
		return
	}

	logger.Println("getting commits...")
	rawCommits, err := prov.GetCommits(release.SHA, currentSha)
	exitIfError(err)

	logger.Println("analyzing commits...")
	commitAnalyzer, err := pluginManager.GetCommitAnalyzer()
	exitIfError(err)
	logger.Printf("commit-analyzer plugin: %s@%s\n", commitAnalyzer.Name(), commitAnalyzer.Version())
	exitIfError(commitAnalyzer.Init(conf.ChangelogGeneratorOpts))

	commits := commitAnalyzer.Analyze(rawCommits)

	logger.Println("calculating new version...")
	newVer := semrel.GetNewVersion(conf, commits, release)
	if newVer == "" {
		herr := hooksExecutor.NoRelease(plugin.NoReleaseReasonNoChange, "")
		if herr != nil {
			logger.Printf("there was an error executing the hooks plugins: %s", herr.Error())
		}
		errNoChange := errors.New("no change")
		if conf.AllowNoChanges {
			exitIfError(errNoChange, 0)
		} else {
			exitIfError(errNoChange, 65)
		}
	}
	logger.Println("new version: " + newVer)

	logger.Println("generating changelog...")
	changelogGenerator, err := pluginManager.GetChangelogGenerator()
	exitIfError(err)
	logger.Printf("changelog-generator plugin: %s@%s\n", changelogGenerator.Name(), changelogGenerator.Version())
	exitIfError(changelogGenerator.Init(conf.ChangelogGeneratorOpts))

	changelogRes := changelogGenerator.Generate(commits, release, newVer)
	if conf.Changelog != "" {
		oldFile := make([]byte, 0)
		if conf.PrependChangelog {
			oldFileData, err := ioutil.ReadFile(conf.Changelog)
			if err == nil {
				oldFile = append([]byte("\n"), oldFileData...)
			}
		}
		changelogData := append([]byte(changelogRes), oldFile...)
		exitIfError(ioutil.WriteFile(conf.Changelog, changelogData, 0644))
	}

	if conf.Dry {
		if conf.VersionFile {
			exitIfError(ioutil.WriteFile(".version-unreleased", []byte(newVer), 0644))
		}
		exitIfError(errors.New("DRY RUN: no release was created"), 0)
	}

	logger.Println("creating release...")
	newRelease := &plugin.CreateReleaseConfig{
		Changelog:  changelogRes,
		NewVersion: newVer,
		Prerelease: conf.Prerelease,
		Branch:     currentBranch,
		SHA:        currentSha,
	}
	exitIfError(prov.CreateRelease(newRelease))

	if conf.Ghr {
		exitIfError(ioutil.WriteFile(".ghr", []byte(fmt.Sprintf("-u %s -r %s v%s", repoInfo.Owner, repoInfo.Repo, newVer)), 0644))
	}

	if conf.VersionFile {
		exitIfError(ioutil.WriteFile(".version", []byte(newVer), 0644))
	}

	if len(conf.UpdateFiles) > 0 {
		logger.Println("updating files...")
		updater, err := pluginManager.GetChainedUpdater()
		exitIfError(err)
		logger.Printf("files-updater plugins: %s\n", strings.Join(updater.GetNameVersionPairs(), ", "))
		exitIfError(updater.Init(conf.FilesUpdaterOpts))

		for _, f := range conf.UpdateFiles {
			exitIfError(updater.Apply(f, newVer))
		}
	}

	herr := hooksExecutor.Success(commits, release, &semrel.Release{
		SHA:     currentSha,
		Version: newVer,
	}, changelogRes, repoInfo)

	if herr != nil {
		logger.Printf("there was an error executing the hooks plugins: %s", herr.Error())
	}

	logger.Println("done.")
}
