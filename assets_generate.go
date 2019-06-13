// +build ignore

package main

import (
	"log"

	"github.com/CyberGRX/api-connector-bulk/assets"
	"github.com/shurcooL/vfsgen"
)

func main() {
	err := vfsgen.Generate(assets.Assets, vfsgen.Options{
		PackageName:  "assets",
		BuildTags:    "!dev",
		VariableName: "Assets",
	})
	if err != nil {
		log.Fatalln(err)
	}
}
