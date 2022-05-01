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

type Command struct {
	Key   string
	Value string
}

//addcmd adds strigns provided in arguments to cmdFile_path as a command
var addcmd = &cobra.Command{
	Use:   "ac [<key>]",
	Short: "add command with specified key. Key is optional",
	Run: func(cmd *cobra.Command, args []string) {
		commands, err := File2Commands(cmdFile_path)
		if err != nil {
			log.Fatalln(err)
		}

		newCommand := Command{}
		//ifnore case args > 1
		if len(args) > 0 {
			_, err := strconv.Atoi(args[0])
			//do not allow Number ID
			if err == nil {
				fmt.Println("[-] Cannot use an integer as a key")
				return
			}
			newCommand.Key = args[0]
		}

		fmt.Print("[+] Type command: ")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			if scanner.Text() != "" {
				newCommand.Value = scanner.Text()
			} else {
				fmt.Println("[*] Aborting")
				return
			}
		}
		//fmt.Printf("Adding Command: %+v\n", newcommand)
		commands = append(commands, newCommand)

		err = WriteCommands(commands)
		if err != nil {
			log.Fatalln(err)
		}

	},
}

//getcmd gets commands in cmdFile_path, and provide them with index number of Commands.Cmd array.
var searchcmd = &cobra.Command{
	Use:   "sc [<search_query>]",
	Short: "search commands. If no query, get all commands",
	Run: func(cmd *cobra.Command, args []string) {
		///unmarshaled json stored in
		commands, err := File2Commands(cmdFile_path)
		if err != nil {
			log.Fatalln(err)
		}
		//if no query specified
		if len(args) < 1 {
			for i, v := range commands {
				fmt.Printf("[%v] %v: %v\n", i, v.Key, v.Value)
			}
			return
		}

		//Accespt space-seperated args as a query
		query := strings.Join(args, " ")
		for i, v := range commands {
			if strings.Contains(v.Value, query) {
				fmt.Printf("[%v] %v: %v\n", i, v.Key, v.Value)
			}
		}
	},
}

//removecmd removes a command specified in argument.
var removecmd = &cobra.Command{
	Use:   "rc <key|index>...",
	Short: "remove commands specified by keys or ids",
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
		if err != nil {
			log.Fatalln(err)
		}

		//TODO: ?Do not initialize, because elements will be added later
		newCommands := []Command{}
		removingIndexes := []int{}
		for i := 0; i < len(args); i++ {
			id, err := strconv.Atoi(args[i])
			//args[i] is a key
			if err != nil {
				id = Key2Id(commands, args[i])
				if id == -1 {
					defer fmt.Println("[*] No command has provided key:", args[i])
					continue
				}
			}
			removingIndexes = append(removingIndexes, id)
		}
		var isfound bool = false
		var removedCmds []Command
		for i, v := range commands {
			if !Contains(removingIndexes, i) {
				newCommands = append(newCommands, v)
			} else {
				isfound = true
				removedCmds = append(removedCmds, v)
				fmt.Printf("[+] Removing [%v] %v: %v\n", i, v.Key, v.Value)
			}
		}
		if !isfound {
			fmt.Println("[*] Nothing removed")
			return
		}

		err = WriteCommands(newCommands)
		if err != nil {
			log.Fatalln(err)
		}
	},
}

var copycmd = &cobra.Command{
	Use:   "cc <key|index>",
	Short: "copy a command replacing <variable> with its <value>",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Too few arguments to call cc")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		commands, err := File2Commands(cmdFile_path)
		if err != nil {
			log.Fatalln(err)
		}

		varmap, err := File2Map(varFile_path)
		if err != nil {
			log.Fatalln(err)
		}

		//args[0] is filtered by Args func, so not produce err
		//TODO
		id, err := strconv.Atoi(args[0])
		if err != nil {
			id = Key2Id(commands, args[0])
			if id == -1 {
				fmt.Println("[-] No command has provided key:", args[0])
				return
			}
		}
		command := commands[id]

		//Fills value to $VARNAME expression
		out := command.Value
		for variable := range varmap {
			out = FillVar(varmap, command.Value, variable)
		}
		fmt.Println(out)
		Copy(out)
	},
}
