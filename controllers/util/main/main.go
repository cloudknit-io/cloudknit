package main

import (
	"context"
	"fmt"
	"github.com/compuzest/zlifecycle-il-operator/controllers/util/git"
	"log"
)

func main() {
	ctx := context.Background()
	g, err := git.NewGoGit(ctx)
	if err != nil {
		log.Fatal(err)
	}
	//dir := "/tmp/il"
	repo := "https://github.com/zlifecycle-il/dev-zl-il.git"
	tempDir, cleanup, err := git.CloneTemp(g, repo)
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()
	fmt.Println(tempDir)

	//if err := os.Remove(dir + "/team/checkout-team-environment/prod-environment-component/networking/terraform/terraforms.tf"); err != nil {
	//	log.Fatal(err)
	//}

	//nfo := git.CommitInfo{Author: "Dejan", Email: "dejan@compuzest.com", Msg: "test commit"}
	//_, err = g.Commit(&nfo)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//if err := g.Push(); err != nil {
	//	log.Fatal(err)
	//}

}
