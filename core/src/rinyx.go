package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

//import docker "modules/docker"

func main() {
	//docker.Affiche()
	/*argsWithProg := os.Args
	argsWithoutProg := os.Args[1:]

	arg := os.Args[1]

	fmt.Println(argsWithProg)
	fmt.Println(argsWithoutProg)
	fmt.Println(arg)
	fmt.Println("hh")*/

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "init":
			initiallisation()
		case "generate":
			generate(os.Args[2:])
		case "build":
			fmt.Println("build all services")
		case "start":
			fmt.Println("start serveur dev")
		case "stop":
			fmt.Println("stop serveur dev")
		case "deploy":
			fmt.Println("deploy project")
		case "help":
			help()

		default:
			help()
		}

	} else {
		help()
	}

}

func help(){
	fmt.Println("this is help for use Rinix framework")
	fmt.Println("rinyx init (init the new Project")
	fmt.Println("rinyx generate application -language lang -backend | -frontend")
	fmt.Println("rinyx generate module application -backend | -frontend")
	fmt.Println("rinyx build (build all service)")
	fmt.Println("rinyx start ( run mode dev)")
	fmt.Println("rinyx stop ( stop mode dev)")
	fmt.Println("rinyx deploy -targz|-docker ( generate  dockerfile and docker-compose.yml or tar.gz )")
}

func generate(args []string) {

	var c cproject
	if args != nil && len(args) > 0 {

		for _, arg := range args {
			if _, err := os.Stat("settings/project/" + arg + "/" + arg + ".yml"); !os.IsNotExist(err) {
				gendprojet(c, arg)
			} else {
				fmt.Println("settings of project " + arg + " not exist")
			}

		}
		/*var name string
		if args[0] == "module" {
			fmt.Println("create module on project")

		}*/
	} else {

		files, _ := ioutil.ReadDir("settings/project/")
		for _, f := range files {
			name := f.Name()
			gendprojet(c, name)
		}
	}

}

func gendprojet(c cproject, name string) {  /*Generate directory for  project*/

	c.getConf("settings/project/" + name + "/" + name + ".yml")

	if c.Type == "service" {
		if c.Git != "" {
			gitclone(c.Git, "project/backoffice/")
		}
		os.Mkdir("storage/"+name, 0777)
		os.Mkdir("project/backoffice/"+name, 0777)
		os.Chdir("project/backoffice/" + name)
		os.Symlink("../../../settings/project/"+name, "settings")
		os.Symlink("../../../storage/"+name, "storage")
		fmt.Println("create service")
	} else if c.Type == "application" {
		if c.Git != "" {
			gitclone(c.Git, "project/front/")
		}
		fmt.Println("create application")
		os.Mkdir("storage/"+name, 0777)
		os.Mkdir("project/front/"+name, 0777)
		os.Chdir("project/front/" + name)
		os.Symlink("../../../settings/project/"+name, "settings")
		os.Symlink("../../../storage/"+name, "storage")
	} else {
		fmt.Println("generate type config is not valide")
	}
	genvendor(c.Lang)
	os.Chdir("../../..")
}
func gitclone(git string, dest string) {
	os.Chdir(dest)
	cmd := exec.Command("git", "clone", git)
	cmd.Stdin = strings.NewReader("some input")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(out.String())
	os.Chdir("../../")
}

func genvendor(lang string) {
	switch lang {
	case "nodejs":
		if _, err := os.Stat("../../../vendor/nodejs"); os.IsNotExist(err) {
			os.Mkdir("../../../vendor/nodejs", 0777)
		}
		os.Symlink("../../../vendor/nodejs", "node_modules")
	case "php":
		if _, err := os.Stat("../../../vendor/php"); os.IsNotExist(err) {
			os.Mkdir("../../../vendor/php", 0777)
		}
		os.Symlink("../../../vendor/php", "vendor")

	default:
		fmt.Println("not language supported")
	}
}

func initiallisation() {

	fmt.Println("init")
	folders := []string { "models","project","project/backoffice","project/front","settings","settings/project","settings/core","storage","vendor" }

	for _, folder := range folders {
        if _, err := os.Stat(folder); os.IsNotExist(err) {
			os.Mkdir(folder, 0777)
		}
    }
	generate(nil)

}
