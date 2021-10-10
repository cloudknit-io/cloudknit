package main

import "github.com/compuzest/zlifecycle-state-manager/web"

func main() {
	//ctx := context.Background()
	//zs := &il.ZState{
	//	RepoURL: "https://github.com/CompuZest/zmart-sandbox-il.git",
	//	Meta: &il.ComponentMeta{
	//		IL: "zmart-sandbox-il",
	//		Team: "design",
	//		Environment: "demo",
	//		Component: "overlays",
	//	},
	//}
	web.NewServer()
}
