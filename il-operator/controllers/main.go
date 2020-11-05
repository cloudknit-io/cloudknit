package main

import github "github.com/compuzest/environment-operator/controllers/github"

func main() {
        github.Setup("CompuZest", "terraform-environment", "1/dev/dev.yaml,1/dev/networking.yaml,1/dev/platform-eks.yaml", "master", "Adarsh Shah", "shahadarsh@gmail.com")
}
