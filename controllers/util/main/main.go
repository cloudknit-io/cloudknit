package main

import (
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/git"
	"log"
	"os"
)

func main() {
	g, err := git.NewGoGit()
	if err != nil {
		log.Fatal(err)
	}
	dir := "/tmp/il"
	if err := g.Clone("https://github.com/zlifecycle-il/dev-zl-il.git", dir); err != nil {
		log.Fatal(err)
	}

	if err := os.Remove(dir+"/team/checkout-team-environment/prod-environment-component/networking/terraform/test.txt"); err != nil {
		log.Fatal(err)
	}

	nfo := git.CommitInfo{Author: "Dejan", Email: "dejan@compuzest.com", Msg: "test commit"}
	_, err = g.Commit(&nfo)
	if err != nil {
		log.Fatal(err)
	}

	if err := g.Push(); err != nil {
		log.Fatal(err)
	}

	if err := os.RemoveAll(dir); err != nil {
		log.Fatal(err)
	}
}
