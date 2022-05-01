package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

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

func File2Commands(filename string) ([]Command, error) {
	f, _ := os.Open(filename)
	defer f.Close()
	bcommands, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	commands := []Command{}
	//if err, file is empty, in which case i just return initialized new(Commands)
	json.Unmarshal(bcommands, &commands)
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

func WriteCommands(commands []Command) error {
	f, err := os.Create(cmdFile_path)
	if err != nil {
		return err
	}
	defer f.Close()

	bcommands, err := json.Marshal(&commands)
	_, err = f.Write(bcommands)
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

func Copy(str string) {
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
		//err during copy does not matter so much
		if err != nil {
			fmt.Println("[-] Failed to copy")
			return
		}
		fmt.Println("[+] Copied")
		return
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
			fmt.Println("[-] Failed to copy")
			return
		}
		fmt.Println("[+] Copied")
		return
	}

	//if pbcopy and xclip was not found
	fmt.Println("[-] Install pbcopy for OSX, xclip for Linux")
}

func ExecuteCmd(command string) {
	fmt.Println("[+] Executing:", command)
	execcmd := exec.Command("bash", "-c", command)
	stdin, _ := execcmd.StdinPipe()
	stdout, _ := execcmd.StdoutPipe()
	stderr, _ := execcmd.StderrPipe()
	err := execcmd.Start()
	if err != nil {
		fmt.Println("[-] Failed to execute command")
		return
	}
	go func() {
		defer stdin.Close()
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			io.WriteString(stdin, scanner.Text()+"\r\n")
		}
	}()
	go func() {
		defer stdout.Close()
		scanner := bufio.NewReader(stdout)
		for {
			obyte, _ := scanner.ReadByte()
			fmt.Print(string(obyte))
		}
	}()
	go func() {
		defer stderr.Close()
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			fmt.Fprintln(os.Stderr, scanner.Text())
		}
	}()
	execcmd.Wait()
	os.Exit(execcmd.ProcessState.ExitCode())

}

func FillVar(varmap map[string]string, str, variable string) string {
	if !strings.Contains(str, "${"+variable+"}") {
		return str
	}
	str = strings.Replace(str, "${"+variable+"}", varmap[variable], -1)
	return str
}

//return -1 for error
func Key2Id(commands []Command, key string) int {
	for i, v := range commands {
		if v.Key == key {
			return i
		}
	}
	return -1
}
