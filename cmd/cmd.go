package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

type Commands struct {
	Cmd []string `json:cmd,omitempty`
}

//addcmd adds strigns provided in arguments to cmdFile_path as a command
var addcmd = &cobra.Command{
	Use:   "ac",
	Short: "add command",
	Long:  "`kop ac` to add command. Type command you want to add to kop.",
	Run: func(cmd *cobra.Command, args []string) {
		commands, err := File2Commands(cmdFile_path)
		CheckErr(err)

		var newcommand string = ""
		fmt.Print("cmd: ")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			if scanner.Text() != "" {
				newcommand = scanner.Text()
			} else {
				log.Fatal("cannot add empty command")
			}
		}
		//fmt.Printf("Adding Command: %+v\n", newcommand)
		commands.Cmd = append(commands.Cmd, newcommand)

		err = WriteCommands(commands)
		CheckErr(err)

		fmt.Printf("[+] Added\n")
	},
}

//getcmd gets commands in cmdFile_path, and provide them with index number of Commands.Cmd array.
var searchcmd = &cobra.Command{
	Use:   "sc [<search_query>]",
	Short: "search commands. If no query, get all commands",
	Run: func(cmd *cobra.Command, args []string) {
		///unmarshaled json stored in
		commands, err := File2Commands(cmdFile_path)
		CheckErr(err)
		//if no query specified
		if len(args) < 1 {
			for i, v := range commands.Cmd {
				fmt.Printf("%v: %v\n", i, v)
			}
		} else {
			query := strings.Join(args, " ")
			for i, v := range commands.Cmd {
				if strings.Contains(v, query) {
					fmt.Printf("%v: %v\n", i, v)
				}
			}
		}
	},
}

//removecmd removes a command specified in argument.
var removecmd = &cobra.Command{
	Use:   "rc <index>...",
	Short: "remove commands specified by index",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Too few arnuments to call rc.")
		} else {
			return nil
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		//Read and to commands
		commands, err := File2Commands(cmdFile_path)
		CheckErr(err)

		//Do not initialize, because elements will be added later
		new_commands := &Commands{}
		removed_indexes := []int{}
		for i := 0; i < len(args); i++ {
			iii, _ := strconv.Atoi(args[i])
			removed_indexes = append(removed_indexes, iii)
		}
		var isfound bool = false
		var removed_cmd string = ""
		for i, v := range commands.Cmd {
			if !Contains(removed_indexes, i) {
				new_commands.Cmd = append(new_commands.Cmd, v)
			} else {
				isfound = true
				removed_cmd = v
			}
		}
		if !isfound {
			fmt.Println("[*] Nothing removed")
			return
		}

		err = WriteCommands(new_commands)
		CheckErr(err)
		fmt.Printf("[+] Removed: %+v\n", removed_cmd)
	},
}

var copycmd = &cobra.Command{
	Use:   "cc <index>",
	Short: "copy command replacing <variable> with its <value>",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Too few arguments to call cc")
		} else if _, err := strconv.Atoi(args[0]); err != nil {
			return errors.New("argument must be integer")
		} else {
			return nil
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		commands, err := File2Commands(cmdFile_path)
		CheckErr(err)

		varmap, err := File2Map(varFile_path)
		CheckErr(err)

		//args[0] is filtered by Args func, so not produce err
		used_cmd_index, _ := strconv.Atoi(args[0])
		used_cmd := commands.Cmd[used_cmd_index]

		//Fills value to $VARNAME expression
		for varname := range varmap {
			used_cmd = FillVar(used_cmd, varname)
		}
		fmt.Println(used_cmd)
		fmt.Println(Copy(used_cmd))
	},
}
