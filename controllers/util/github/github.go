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
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
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

func RemoveObjectsFromBranch(
	log logr.Logger,
	api GitApi,
	owner string,
	repo string,
	branch string,
	baseDir string,
	paths []string,
	commitAuthor *github.CommitAuthor,
	commitMessage string,
) error {
	refFormat := fmt.Sprintf("refs/heads/%s", branch)
	log.Info("Fetching ref...", "owner", owner, "repo", repo, "ref", refFormat)
	ref, refResp, refErr := api.GetRef(owner, repo, refFormat)
	if refErr != nil {
		return refErr
	}
	defer refResp.Body.Close()

	log.Info("Fetching git tree...", "owner", owner, "repo", repo, "ref", *ref.Ref)
	tree, treeResp, treeErr := api.GetTree(owner, repo, *ref.Ref, true)
	if treeErr != nil {
		return treeErr
	}
	defer treeResp.Body.Close()

	var subtree *github.Tree

	if baseDir != "" {
		var subtreeSHA string
		for _, entry := range tree.Entries {
			if baseDir == *entry.Path {
				log.Info("Found subtree", "tree", tree.SHA, "subtree", entry.Path)
				subtreeSHA = *entry.SHA
				break
			}
		}
		tree, treeResp, treeErr := api.GetTree(owner, repo, subtreeSHA, true)
		if treeErr != nil {
			return treeErr
		}
		defer treeResp.Body.Close()

		log.Info("Deleting paths from tree...", "sha", tree.SHA, "paths", paths)
		entries := removePathsFromTree(log, tree, paths)

		log.Info("Creating new subtree without the excluded values...")
		newTree, newTreeResp, newTreeErr := api.CreateTree(owner, repo, "", entries)
		if newTreeErr != nil {
			return newTreeErr
		}
		defer newTreeResp.Body.Close()

		subtree = newTree
	}

	//for _, entry := range tree.Entries {
	//	if baseDir == *entry.Path {
	//		log.Info("Updating subtree", "path", *entry.Path, "oldSHA", *entry.SHA, "newSHA", subtree.SHA)
	//		entry.SHA = subtree.SHA
	//		entry.Content = subtree.Con
	//		break
	//	}
	//}

	log.Info("Creating new subtree without the excluded values...")
	newTree, newTreeResp, newTreeErr := api.CreateTree(owner, repo, *tree.SHA, subtree.Entries)
	if newTreeErr != nil {
		return newTreeErr
	}
	defer newTreeResp.Body.Close()

	log.Info("Getting parent commit...", "parentSha", ref.Object.SHA)
	parentCommit, commitResp, commitErr := api.GetCommit(owner, repo, *ref.Object.SHA)
	if commitErr != nil {
		return commitErr
	}
	defer commitResp.Body.Close()

	log.Info("Parent commit info",
		"parentTreeSha", parentCommit.Tree.SHA,
		"parentSha", parentCommit.SHA,
	)

	commit := gitCommit(newTree, commitAuthor, commitMessage, parentCommit)

	log.Info(
		"Creating commit for updated tree...",
		"parentCommitSha", parentCommit.SHA,
		"parentTreeSha", parentCommit.Tree.SHA,
		"baseTreeSha", tree.SHA,
		"newTreeSha", newTree.SHA,
	)
	newCommit, _, commitErr := api.CreateCommit(owner, repo, commit)
	if commitErr != nil {
		return commitErr
	}

	if !hasChangesToCommit(parentCommit, newTree) {
		log.Info(
			"No changes in il repo for deletion. Ignoring empty commit.",
			"parentCommitSha", parentCommit.SHA,
			"parentTreeSha", parentCommit.Tree.SHA,
			"baseTreeSha", tree.SHA,
			"newTreeSha", newTree.SHA,
		)
		return nil
	}

	log.Info("Updating ref with new commit SHA...", "oldSha", ref.Object.SHA, "newSha", newCommit.SHA)
	ref.Object.SHA = newCommit.SHA
	_, newRefResp, newRefErr := api.UpdateRef(owner, repo, ref, false)
	if newRefErr != nil {
		return newRefErr
	}
	defer newRefResp.Body.Close()

	return nil
}

func hasChangesToCommit(parent *github.Commit, tree *github.Tree) bool {
	parentSha := *parent.Tree.SHA
	newSha := *tree.SHA

	return parentSha != newSha
}

func gitCommit(tree *github.Tree, author *github.CommitAuthor, message string, parentCommit *github.Commit) *github.Commit {
	return &github.Commit{Author: author, Message: &message, Tree: tree, Parents: []*github.Commit{parentCommit}}
}

func removePathsFromTree(log logr.Logger, tree *github.Tree, paths []string) []*github.TreeEntry {
	var entries []*github.TreeEntry
	for _, entry := range tree.Entries {
		if shouldExclude(entry, paths) {
			log.Info("Excluding path from base tree", "path", entry.Path)
		} else {
			e := &github.TreeEntry{
				SHA:     entry.SHA,
				Content: entry.Content,
				Path:    entry.Path,
				Size:    entry.Size,
				URL:     entry.URL,
				Type:    entry.Type,
				Mode:    entry.Mode,
			}
			entries = append(entries, e)
		}
	}
	return entries
}

func shouldExclude(entry *github.TreeEntry, paths []string) bool {
	for _, pathPattern := range paths {
		if exclude, _ := regexp.MatchString(pathPattern, *entry.Path); exclude {
			return true
		}
	}
	return  false
}

// TryCreateRepository tries to create a private repository in a organization (enter blank string for owner if it is a user repo)
// It returns true if repository is created, false if repository already exists, or any kind of error.
func TryCreateRepository(log logr.Logger, api RepositoryApi, owner string, repo string) (bool, error) {
	log.Info("Checking does repo exist on GitHub", "owner", owner, "repo", repo)
	r, resp1, err := api.GetRepository(owner, repo)
	if err != nil {
		log.Error(err, "Error while fetching repository from GitHub API",
			"owner", owner,
			"repo", repo,
		)
		return false, err
	}
	defer resp1.Body.Close()

	log.Info("Call to GitHub API get repo succeeded", "code", resp1.StatusCode)

	if r != nil {
		log.Info("GitHub repository already exists", "owner", owner, "repo", repo)
		return false, nil
	}

	log.Info("Creating new private repository in GitHub", "owner", owner, "repo", repo)

	nr, resp2, err := api.CreateRepository(owner, repo)
	if err != nil {
		log.Error(err, "Error while creating private repository on GitHub",
			"owner", owner,
			"repo", repo,
		)
		return false, err
	}

	log.Info("Call to GitHub API create repo succeeded", "code", resp2.StatusCode)

	defer resp2.Body.Close()

	return nr != nil, nil
}

// CreateRepoWebhook tries to create a repository webhook
// It returns true if webhook already exists, false if new webhook was created, or any kind of error.
func CreateRepoWebhook(log logr.Logger, api RepositoryApi, repoUrl string, payloadUrl string, webhookSecret string) (bool, error) {
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

	h := newHook(payloadUrl, webhookSecret)
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
		"repo", repo,
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

func newHook(payloadUrl string, secret string) github.Hook {
	events := []string{"push"}
	cfg := map[string]interface{}{
		"url":          payloadUrl,
		"content_type": "json",
		"secret":       secret,
	}
	return github.Hook{Events: events, Active: github.Bool(true), Config: cfg}
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
