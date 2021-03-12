/* Copyright (C) 2020 CompuZest, Inc. - All Rights Reserved
 *
 * Unauthorized copying of this file, via any medium, is strictly prohibited
 * Proprietary and confidential
 *
 * NOTICE: All information contained herein is, and remains the property of
 * CompuZest, Inc. The intellectual and technical concepts contained herein are
 * proprietary to CompuZest, Inc. and are protected by trade secret or copyright
 * law. Dissemination of this information or reproduction of this material is
 * strictly forbidden unless prior written permission is obtained from CompuZest, Inc.
 */

package github

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/compuzest/zlifecycle-il-operator/controllers/util/common"
	"github.com/go-logr/logr"

	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

var (
	sourceOwner   string
	sourceRepo    string
	commitMessage string
	commitBranch  string
	baseBranch    string
	prRepoOwner   string
	prRepo        string
	prBranch      string
	prSubject     string
	prDescription string
	sourceFiles   string
	authorName    string
	authorEmail   string
)

var client *github.Client
var ctx = context.Background()

// CreateRepoWebhook tries to create a repository webhook
// It returns true if webhook already exists, false if new webhook was created, or any kind of error.
func CreateRepoWebhook(log logr.Logger, api RepositoryApi, repoUrl string, payloadUrl string) (bool, error) {
	owner, repo, err := parseRepoUrl(repoUrl)
	if err != nil {
		log.Error(err, "Error while parsing owner and repo name from repo url", "url", repoUrl)
		return false, err
	}
	log.Info("Parsed repo url", "url", repoUrl, "owner", owner, "repo", repo)

	hooks, resp1, err := api.ListHooks(owner, repo, nil)
	if err != nil {
		log.Error(err, "Error while fetching list of webhooks",
			"url", repoUrl,
			"owner", owner,
			"repo", repo,
		)
		return false, err
	}
	defer resp1.Body.Close()

	exists, err := checkIsHookRegistered(
		log,
		hooks,
		payloadUrl,
	)
	if exists {
		log.Info(
			"Hook already exists",
			"url", repoUrl,
			"owner", owner,
			"repo", repo,
			"payloadUrl", payloadUrl,
		)
		return true, nil
	}
	if err != nil {
		log.Error(err, "Error while checking is hook already registered", "hooks", hooks)
		return false, err
	}

	h := newHook(payloadUrl)
	hook, resp2, err := api.CreateHook(owner, repo, &h)
	if err != nil {
		log.Error(
			err, "Error while calling create repository hook",
			"url", repoUrl,
			"owner", owner,
			"repo", repo,
			"cfg", h.Config,
		)
		return false, err
	}
	defer resp2.Body.Close()

	log.Info(
		"Successfully created repository webhook",
		"url", repoUrl,
		"owner", owner,
		"repo",	repo,
		"hookId", *hook.ID,
		"payloadUrl", *hook.URL,
		"cfg", h.Config,
	)
	return false, nil
}

func parseRepoUrl(url string) (Owner, Repo, error) {
	if url == "" {
		return "", "", errors.New("URL cannot be empty")
	}

	owner := url[strings.LastIndex(url, ":")+1 : strings.LastIndex(url, "/")]
	repoUri := url[strings.LastIndex(url, "/")+1:]
	repo := strings.TrimSuffix(repoUri, ".git")

	return owner, repo, nil
}

func checkIsHookRegistered(log logr.Logger, hooks []*github.Hook, payloadUrl string) (bool, error) {
	for _, h := range hooks {
		cfg := new(HookCfg)
		err := common.FromJsonMap(log, h.Config, cfg)
		if err != nil {
			return false, err
		}
		if cfg.Url == payloadUrl {
			return true, nil
		}
	}

	return false, nil
}

func newHook(payloadUrl string) github.Hook {
	isActive := true
	events := []string{"push"}
	cfg := map[string]interface{}{
		"url":          payloadUrl,
		"content_type": "json",
	}
	return github.Hook{Events: events, Active: &isActive, Config: cfg}
}

// getRef returns the commit branch reference object if it exists or return error
func getRef() (ref *github.Reference, err error) {
	if commitBranch == "" {
		return nil, errors.New("The `-commit-branch` should not be set to an empty string")
	}

	if ref, _, err = client.Git.GetRef(ctx, sourceOwner, sourceRepo, "refs/heads/"+commitBranch); err == nil {
		return ref, nil
	}

	return ref, err
}

// getTree generates the tree to commit based on the given files and the commit
// of the ref you got in getRef.
func getTree(ref *github.Reference) (tree *github.Tree, err error) {
	// Create a tree with what to commit.
	entries := []*github.TreeEntry{}

	// Load each file into the tree.
	for _, fileArg := range strings.Split(sourceFiles, ",") {
		file, content, err := getFileContent(fileArg)
		if err != nil {
			return nil, err
		}
		entries = append(entries, &github.TreeEntry{Path: github.String(file), Type: github.String("blob"), Content: github.String(string(content)), Mode: github.String("100644")})
	}

	tree, _, err = client.Git.CreateTree(ctx, sourceOwner, sourceRepo, *ref.Object.SHA, entries)
	return tree, err
}

// getFileContent loads the local content of a file and return the target name
// of the file in the target repository and its contents.
func getFileContent(fileArg string) (targetName string, b []byte, err error) {
	var localFile string
	files := strings.Split(fileArg, ":")
	switch {
	case len(files) < 1:
		return "", nil, errors.New("empty `-files` parameter")
	case len(files) == 1:
		localFile = files[0]
		targetName = files[0]
	default:
		localFile = files[0]
		targetName = files[1]
	}

	b, err = ioutil.ReadFile(localFile)
	return targetName, b, err
}

// pushCommit creates the commit in the given reference using the given tree.
func pushCommit(ref *github.Reference, tree *github.Tree) (err error) {
	// Get the parent commit to attach the commit to.
	parent, _, err := client.Repositories.GetCommit(ctx, sourceOwner, sourceRepo, *ref.Object.SHA)
	if err != nil {
		return err
	}
	// This is not always populated, but is needed.
	parent.Commit.SHA = parent.SHA

	parentSha := *parent.Commit.Tree.SHA
	newSha := *tree.SHA
	// log.Printf(string(json.Marshal(tree))) for debugging

	if parentSha == newSha {
		log.Printf("No git changes to commit, no-op reconciliation.")
	} else {
		// Create the commit using the tree.
		date := time.Now()
		author := &github.CommitAuthor{Date: &date, Name: &authorName, Email: &authorEmail}
		commit := &github.Commit{Author: author, Message: &commitMessage, Tree: tree, Parents: []*github.Commit{parent.Commit}}
		newCommit, _, err := client.Git.CreateCommit(ctx, sourceOwner, sourceRepo, commit)
		if err != nil {
			return err
		}

		// Attach the commit to the master branch.
		ref.Object.SHA = newCommit.SHA
		_, _, err = client.Git.UpdateRef(ctx, sourceOwner, sourceRepo, ref, false)
	}
	return err
}

func CommitAndPushFiles(_sourceOwner string, _sourceRepo string, _sourceFolders []string,
	_commitBranch string, _commitMessage string, _authorName string, _authorEmail string) (err error) {
	sourceOwner = _sourceOwner
	sourceRepo = _sourceRepo
	commitBranch = _commitBranch
	commitMessage = _commitMessage
	baseBranch = _commitBranch
	authorName = _authorName
	authorEmail = _authorEmail

	sourceFiles = getFileNames(_sourceFolders)

	token := os.Getenv("GITHUB_AUTH_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}
	if sourceOwner == "" || sourceRepo == "" || commitBranch == "" || sourceFiles == "" || authorName == "" || authorEmail == "" {
		log.Fatal("You need to specify a non-empty value for the flags `-source-owner`, `-source-repo`, `-commit-branch`, `-files`, `-author-name` and `-author-email`")
	}
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client = github.NewClient(tc)

	ref, err := getRef()
	if err != nil {
		log.Fatalf("Unable to get/create the commit reference: %s\n", err)
	}
	if ref == nil {
		log.Fatalf("No error where returned but the reference is nil")
	}

	tree, err := getTree(ref)
	if err != nil {
		log.Fatalf("Unable to create the tree based on the provided files: %s\n", err)
	}

	if err := pushCommit(ref, tree); err != nil {
		log.Fatalf("Unable to create the commit: %s\n", err)
	}

	return err
}

func getFileNames(folderPaths []string) (fileNames string) {
	var s, sep string
	for _, f := range folderPaths {
		err := filepath.Walk(f,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					log.Fatalf("Error fetching fileNames: %s\n", err)
					return err
				}

				if !info.IsDir() {
					s += sep + path
					sep = ","
				}
				return nil
			})

		if err != nil {
			log.Fatalf("Error fetching fileNames: %s\n", err)
		}
	}
	return s
}
