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
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/compuzest/zlifecycle-il-operator/controllers/util/common"
	"github.com/go-logr/logr"

	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

var (
	mutex         sync.Mutex
	sourceOwner   string
	sourceRepo    string
	commitMessage string
	commitBranch  string
	sourceFiles   string
	authorName    string
	authorEmail   string
)

var (
	client *github.Client
	ctx    = context.Background()
)

func DownloadFile(
	api RepositoryAPI,
	repoURL string,
	ref string,
	path string,
) (file io.ReadCloser, exists bool, err error) {
	owner, repo, err := parseRepoURL(repoURL)
	if ref == "" {
		ref = "HEAD"
	}
	if err != nil {
		return nil, false, err
	}
	rc, exists, err := api.DownloadContents(owner, repo, ref, path)
	if err != nil {
		return nil, false, err
	}
	return rc, exists, nil
}

func DeletePatternsFromRootTree(
	log logr.Logger,
	api GitAPI,
	owner string,
	repo string,
	branch string,
	team string,
	patterns []string,
	commitAuthor *github.CommitAuthor,
	commitMessage string,
) error {
	tree, err := removeEnvironmentObjectsFromTree(log, api, owner, repo, branch, team, patterns)
	if err != nil {
		return err
	}
	err = commitAndPushTree(log, api, owner, repo, branch, tree, commitAuthor, commitMessage)
	if err != nil {
		return err
	}

	return nil
}

func removeEnvironmentObjectsFromTree(
	log logr.Logger,
	api GitAPI,
	owner string,
	repo string,
	branch string,
	team string,
	patterns []string,
) (*github.Tree, error) {
	// get base tree
	log.Info("Fetching base tree...", "owner", owner, "repo", repo, "branch", branch)
	baseTree, baseTreeResp, baseTreeErr := api.GetTree(owner, repo, branch, false)
	if baseTreeErr != nil {
		return nil, baseTreeErr
	}
	defer common.CloseBody(baseTreeResp.Body)

	// team root tree
	teamRootPath := "team"
	log.Info("Fetching team root tree...", "owner", owner, "repo", repo, "path", teamRootPath)
	var teamRootTree *github.Tree
	for _, entry := range baseTree.Entries {
		if *entry.Path == teamRootPath {
			tree, treeResp, treeErr := api.GetTree(owner, repo, *entry.SHA, false)
			if treeErr != nil {
				return nil, treeErr
			}
			common.CloseBody(treeResp.Body)
			teamRootTree = tree
		}
	}
	if teamRootTree == nil {
		return nil, errors.New("missing team root tree")
	}

	// team tree
	teamPath := team
	log.Info("Fetching team tree...", "owner", owner, "repo", repo, "path", teamPath)
	var teamTree *github.Tree
	for _, entry := range teamRootTree.Entries {
		if *entry.Path == teamPath {
			tree, treeResp, treeErr := api.GetTree(owner, repo, *entry.SHA, false)
			if treeErr != nil {
				return nil, treeErr
			}
			common.CloseBody(treeResp.Body)
			teamTree = tree
		}
	}
	if teamTree == nil {
		return nil, errors.New("missing team tree")
	}

	// exclude entries
	entries := removePathsFromTree(log, teamTree, patterns)

	// update team tree
	log.Info("Creating new team subtree...")
	newTeamTree, newTeamTreeResp, newTeamTreeErr := api.CreateTree(owner, repo, "", entries)
	if newTeamTreeErr != nil {
		return nil, newTeamTreeErr
	}
	defer common.CloseBody(newTeamTreeResp.Body)

	// update root team tree
	for _, entry := range teamRootTree.Entries {
		if teamPath == *entry.Path {
			log.Info("Updating SHA value for team subtree", "path", *entry.Path, "oldSHA", *entry.SHA, "newSHA", newTeamTree.SHA)
			entry.SHA = newTeamTree.SHA
			entry.URL = nil
		}
	}

	log.Info("Creating new team root tree...")
	newTeamRootTree, newTeamRootTreeResp, newTeamRootTreeErr := api.CreateTree(owner, repo, "", teamRootTree.Entries)
	if newTeamRootTreeErr != nil {
		return nil, newTeamRootTreeErr
	}
	defer common.CloseBody(newTeamRootTreeResp.Body)

	// update base tree
	for _, entry := range baseTree.Entries {
		if teamRootPath == *entry.Path {
			log.Info("Updating SHA value for root team subtree", "path", *entry.Path, "oldSHA", *entry.SHA, "newSHA", newTeamRootTree.SHA)
			entry.SHA = newTeamRootTree.SHA
			entry.URL = nil
		}
	}

	log.Info("Creating new base tree...")
	newBaseTree, newBaseTreeResp, newBaseTreeErr := api.CreateTree(owner, repo, "", baseTree.Entries)
	if newBaseTreeErr != nil {
		return nil, newBaseTreeErr
	}
	defer common.CloseBody(newBaseTreeResp.Body)

	return newBaseTree, nil
}

func commitAndPushTree(
	log logr.Logger,
	api GitAPI,
	owner string,
	repo string,
	branch string,
	tree *github.Tree,
	commitAuthor *github.CommitAuthor,
	commitMessage string,
) error {
	// get base tree
	refFormat := fmt.Sprintf("refs/heads/%s", branch)
	log.Info("Fetching ref...", "owner", owner, "repo", repo, "ref", refFormat)
	ref, refResp, refErr := api.GetRef(owner, repo, refFormat)
	if refErr != nil {
		return refErr
	}
	defer common.CloseBody(refResp.Body)

	log.Info("Getting parent commit...", "parentSHA", ref.Object.SHA)
	parentCommit, commitResp, commitErr := api.GetCommit(owner, repo, *ref.Object.SHA)
	if commitErr != nil {
		return commitErr
	}
	defer common.CloseBody(commitResp.Body)

	log.Info("Parent commit info",
		"parentTreeSha", parentCommit.Tree.SHA,
		"parentSha", parentCommit.SHA,
	)

	c := gitCommit(tree, commitAuthor, commitMessage, parentCommit)

	log.Info(
		"Creating commit for tree...",
		"parentCommitSha", parentCommit.SHA,
		"parentTreeSha", parentCommit.Tree.SHA,
		"newTreeSha", tree.SHA,
	)
	commit, commitResp, commitErr := api.CreateCommit(owner, repo, c)
	if commitErr != nil {
		return commitErr
	}
	defer common.CloseBody(commitResp.Body)

	if !hasChangesToCommit(parentCommit, tree) {
		log.Info(
			"No changes in il repo for deletion. Ignoring empty commit.",
			"parentCommitSha", parentCommit.SHA,
			"parentTreeSha", parentCommit.Tree.SHA,
			"newTreeSha", tree.SHA,
		)
		return nil
	}

	log.Info("Updating ref with new commit SHA...", "oldSha", ref.Object.SHA, "newSha", commit.SHA)
	ref.Object.SHA = commit.SHA
	_, newRefResp, newRefErr := api.UpdateRef(owner, repo, ref, false)
	if newRefErr != nil {
		return newRefErr
	}
	defer common.CloseBody(newRefResp.Body)

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
	return false
}

// TryCreateRepository tries to create a private repository in a organization (enter blank string for owner if it is a user repo)
// It returns true if repository is created, false if repository already exists, or any kind of error.
func TryCreateRepository(log logr.Logger, api RepositoryAPI, owner string, repo string) (bool, error) {
	log.Info("Checking does repo exist on GitHub", "owner", owner, "repo", repo)
	r, resp1, err := api.GetRepository(owner, repo)
	if err != nil {
		log.Error(err, "Error while fetching repository from GitHub API",
			"owner", owner,
			"repo", repo,
		)
		return false, err
	}
	defer common.CloseBody(resp1.Body)

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

	defer common.CloseBody(resp2.Body)

	return nr != nil, nil
}

// CreateRepoWebhook tries to create a repository webhook
// It returns true if webhook already exists, false if new webhook was created, or any kind of error.
func CreateRepoWebhook(log logr.Logger, api RepositoryAPI, repoURL string, payloadURL string, webhookSecret string) (bool, error) {
	owner, repo, err := parseRepoURL(repoURL)
	if err != nil {
		log.Error(err, "Error while parsing owner and repo name from repo url", "url", repoURL)
		return false, err
	}
	log.Info("Parsed repo url", "url", repoURL, "owner", owner, "repo", repo)

	hooks, resp1, err := api.ListHooks(owner, repo, nil)
	if err != nil {
		log.Error(err, "Error while fetching list of webhooks",
			"url", repoURL,
			"owner", owner,
			"repo", repo,
		)
		return false, err
	}
	defer common.CloseBody(resp1.Body)

	exists, err := checkIsHookRegistered(
		hooks,
		payloadURL,
	)
	if exists {
		log.Info(
			"Hook already exists",
			"url", repoURL,
			"owner", owner,
			"repo", repo,
			"payloadUrl", payloadURL,
		)
		return true, nil
	}
	if err != nil {
		log.Error(err, "Error while checking is hook already registered", "hooks", hooks)
		return false, err
	}

	h := newHook(payloadURL, webhookSecret)
	hook, resp2, err := api.CreateHook(owner, repo, &h)
	if err != nil {
		log.Error(
			err, "Error while calling create repository hook",
			"url", repoURL,
			"owner", owner,
			"repo", repo,
		)
		return false, err
	}
	defer common.CloseBody(resp2.Body)

	log.Info(
		"Successfully created repository webhook",
		"url", repoURL,
		"owner", owner,
		"repo", repo,
		"hookId", *hook.ID,
		"payloadUrl", *hook.URL,
	)
	return false, nil
}

func parseRepoURL(url string) (Owner, Repo, error) {
	if url == "" {
		return "", "", errors.New("URL cannot be empty")
	}

	owner := url[strings.LastIndex(url, ":")+1 : strings.LastIndex(url, "/")]
	repoURI := url[strings.LastIndex(url, "/")+1:]
	repo := strings.TrimSuffix(repoURI, ".git")

	return owner, repo, nil
}

func checkIsHookRegistered(hooks []*github.Hook, payloadURL string) (bool, error) {
	for _, h := range hooks {
		cfg := new(HookCfg)
		err := common.FromJSONMap(h.Config, cfg)
		if err != nil {
			return false, err
		}
		if cfg.URL == payloadURL {
			return true, nil
		}
	}

	return false, nil
}

func newHook(payloadURL string, secret string) github.Hook {
	events := []string{"push"}
	cfg := map[string]interface{}{
		"url":          payloadURL,
		"content_type": "json",
		"secret":       secret,
	}
	return github.Hook{Events: events, Active: github.Bool(true), Config: cfg}
}

// getRef returns the commit branch reference object if it exists or return error
func getRef() (ref *github.Reference, err error) {
	if commitBranch == "" {
		return nil, errors.New("the `-commit-branch` should not be set to an empty string")
	}
	if ref, _, err = client.Git.GetRef(ctx, sourceOwner, sourceRepo, fmt.Sprintf("refs/heads/%s", commitBranch)); err == nil {
		return ref, nil
	}

	return ref, err
}

// getTree generates the tree to commit based on the given files and the commit
// of the ref you got in getRef.
func getTree(ref *github.Reference) (tree *github.Tree, err error) {
	// Create a tree with what to commit.
	splittedFiles := strings.Split(sourceFiles, ",")
	entries := make([]*github.TreeEntry, 0, len(splittedFiles))

	// Load each file into the tree.
	for _, fileArg := range splittedFiles {
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
func pushCommit(ref *github.Reference, tree *github.Tree) (dirty bool, err error) {
	dirty = false
	// Get the parent commit to attach the commit to.
	parent, _, err := client.Repositories.GetCommit(ctx, sourceOwner, sourceRepo, *ref.Object.SHA)
	if err != nil {
		return dirty, err
	}
	// This is not always populated, but is needed.
	parent.Commit.SHA = parent.SHA

	parentSha := *parent.Commit.Tree.SHA

	// log.Printf(string(json.Marshal(tree))) for debugging
	if newSha := *tree.SHA; parentSha != newSha {
		dirty = true
		// Create the commit using the tree.
		date := time.Now()
		author := &github.CommitAuthor{Date: &date, Name: &authorName, Email: &authorEmail}
		commit := &github.Commit{Author: author, Message: &commitMessage, Tree: tree, Parents: []*github.Commit{parent.Commit}}
		newCommit, _, err := client.Git.CreateCommit(ctx, sourceOwner, sourceRepo, commit)
		if err != nil {
			return dirty, err
		}

		// Attach the commit to the master branch.
		ref.Object.SHA = newCommit.SHA
		_, _, err = client.Git.UpdateRef(ctx, sourceOwner, sourceRepo, ref, false)
		if err != nil {
			return dirty, err
		}
	}
	return dirty, nil
}

func CommitAndPushFiles(_sourceOwner string, _sourceRepo string, _sourceFolders []string,
	_commitBranch string, _commitMessage string, _authorName string, _authorEmail string) (dirty bool, err error) {
	mutex.Lock()
	defer mutex.Unlock()
	sourceOwner = _sourceOwner
	sourceRepo = _sourceRepo
	commitBranch = _commitBranch
	commitMessage = _commitMessage
	authorName = _authorName
	authorEmail = _authorEmail

	sourceFiles, err = getFileNames(_sourceFolders)
	if err != nil {
		return false, err
	}

	token := os.Getenv("GITHUB_AUTH_TOKEN")
	if token == "" {
		return false, fmt.Errorf("unauthorized: No token present")
	}
	if sourceOwner == "" || sourceRepo == "" || commitBranch == "" || sourceFiles == "" || authorName == "" || authorEmail == "" {
		return false, fmt.Errorf("you need to specify a non-empty value for the flags `-source-owner`, `-source-repo`, `-commit-branch`, `-files`, `-author-name` and `-author-email`")
	}
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client = github.NewClient(tc)

	ref, err := getRef()
	if err != nil {
		return false, fmt.Errorf("unable to get/create the commit reference: %w", err)
	}
	if ref == nil {
		return false, fmt.Errorf("no error where returned but the reference is nil")
	}

	tree, err := getTree(ref)
	if err != nil {
		return false, fmt.Errorf("unable to create the tree based on the provided files: %w", err)
	}

	dirty, err = pushCommit(ref, tree)
	if err != nil {
		return false, fmt.Errorf("unable to create the commit: %w", err)
	}

	return dirty, nil
}

func getFileNames(folderPaths []string) (fileNames string, err error) {
	var s, sep string
	for _, f := range folderPaths {
		err := filepath.Walk(f,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if !info.IsDir() {
					s += sep + path
					sep = ","
				}
				return nil
			})
		if err != nil {
			return "", err
		}
	}
	return s, nil
}
