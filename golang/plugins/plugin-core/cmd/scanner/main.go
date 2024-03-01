package main

import (
	"black_hat_go/plugins/plugin-core/scanner"
	"fmt"
	"log"
	"os"
	"plugin"
)

const PluginDir = "../../plugins"

func main() {
	var (
		files []os.FileInfo
		err   error
		p     *plugin.Plugin
		n     plugin.Symbol
		check scanner.Checker
		res   *scanner.Result
	)
	filesInDir, err := os.ReadDir(PluginDir)
	if err != nil {
		log.Fatalln(err)
	}
	for _, file := range filesInDir {
		fileInfo, err := file.Info()
		if err != nil {
			log.Fatalln(err)
		}
		files = append(files, fileInfo)
	}
	for idx := range files {
		fmt.Printf("Found plugin: %s\n", files[idx].Name())
		if p, err = plugin.Open(PluginDir + "/" + files[idx].Name()); err != nil {
			log.Fatalln(err)
		}
		if n, err = p.Lookup("New"); err != nil {
			log.Fatalln(err)
		}
		newFunc, ok := n.(func() scanner.Checker)
		if !ok {
			log.Fatalln("Plugin entry point is no good. Expecting: func New() scanner. Checker{...}")
		}
		check = newFunc()
		res = check.Check("10.0.0.20", 8080)
		if res.Vulnerable {
			log.Printf("Host is vulnerable: %s\n", res.Details)
		} else {
			log.Println("Host is not vulnerable")
		}

	}
}
