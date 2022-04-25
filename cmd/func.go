package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func CheckErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

//File2Map reads filename which is expected to be exists, and return map of file
func File2Map(filename string) (map[string]string, error) {
	//TODO
	f, _ := os.Open(filename)
	defer f.Close()
	bvarjson, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	varmap := map[string]string{}
	//json to map
	//if err, file is empty, in which case we just return non-nil varmap (addvar.go will create json file with newly added var)
	json.Unmarshal(bvarjson, &varmap)
	return varmap, nil
}

func File2Commands(filename string) (*Commands, error) {
	f, _ := os.Open(filename)
	defer f.Close()
	bcommands, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	commands := new(Commands)
	//if err, file is empty, in which case i just return initialized new(Commands)
	json.Unmarshal(bcommands, commands)
	return commands, nil
}

func WriteMap(contents map[string]string, filename string) error {
	bcontents, err := json.Marshal(contents)
	if err != nil {
		return err
	}
	//by creating file, overwrite
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(bcontents)
	if err != nil {
		return err
	}

	return nil
}

func WriteCommands(contents *Commands) error {
	f, err := os.Create(cmdFile_path)
	if err != nil {
		return err
	}
	defer f.Close()

	bcontents, err := json.Marshal(contents)
	_, err = f.Write(bcontents)
	if err != nil {
		return err
	}

	return nil
}

//Contains checks if slice contains element. argument slice must be []string or []int
func Contains(slice interface{}, element interface{}) bool {
	switch slice.(type) {
	case []string:
		for _, v := range slice.([]string) {
			if v == element {
				return true
			}
		}
	case []int:
		for _, v := range slice.([]int) {
			if v == element {
				return true
			}
		}
	}
	return false
}

func Copy(str string) string {
	//copy to clipboar
	if _, err := exec.LookPath("pbcopy"); err == nil && runtime.GOOS == "darwin" {
		//for MacOS
		//Automatically Escape character with &q
		str = fmt.Sprintf("%q", str)
		//prevent "$NotAddedToMMM" from filled in
		str = strings.Replace(str, "$", `\$`, -1) //Change $ to \$ so echo -n %q won't change $VARIABLE to value (using ``)
		echopy := fmt.Sprintf(`echo -n %s | pbcopy`, str)
		execcmd := exec.Command("bash", "-c", echopy)
		err = execcmd.Run()
		//err during copy does not matter so much
		if err != nil {
			return "Failed to copy."
		}
		return "Copied"
	}
	if _, err := exec.LookPath("xclip"); err == nil && runtime.GOOS == "linux" {
		//for linux
		str = fmt.Sprintf("%q", str)
		//prevent "$VARIABLE" from filled in
		str = strings.Replace(str, "$", `\$`, -1) //Change $ to \$ so echo -n %q won't change $VARIABLE to value (using ``)
		echolip := fmt.Sprintf(`echo -n %s | xclip -selection c`, str)
		execcmd := exec.Command("bash", "-c", echolip)
		err = execcmd.Run()
		//err during copy does not matter so much
		if err != nil {
			return "Failed to copy."
		}
		return "Copied"
	}

	//if pbcopy and xclip was not found
	return "[*] You should install pbcopy if you are Mac or xclip if you are Linux, to enable kop to have access to system clipboard"
}

func FillVar(str, variable string) string {
	if strings.Contains(str, "${"+variable+"}") {
		varmap, err := File2Map(varFile_path)
		CheckErr(err)

		if varmap[variable] != "" {
			str = strings.Replace(str, "${"+variable+"}", varmap[variable], -1)
		}
	}

	return str
}
