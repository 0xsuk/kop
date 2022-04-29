package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var (
	homeDir, _   = homedir.Dir()
	varFile_path = homeDir + "/.kopvar.json"
	cmdFile_path = homeDir + "/.kopcmd.json"
)

//TODO create Function to Return map, struct of file
//Read JsonFile and show
var rootCmd = &cobra.Command{
	Use:   "kop [command]",
	Short: "Visit https://github.com/0xsuk/kop to understand basic usage&concepts",
	Run: func(cmd *cobra.Command, args []string) {
		searchvar.Run(searchvar, []string{})
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	configInit()
	rootCmd.AddCommand(searchvar)
	rootCmd.AddCommand(copyvar)
	rootCmd.AddCommand(addvar)
	rootCmd.AddCommand(removevar)
	rootCmd.AddCommand(addcmd)
	rootCmd.AddCommand(searchcmd)
	rootCmd.AddCommand(removecmd)
	rootCmd.AddCommand(copycmd)

	//For usage generation
	//TODO: find a better way to generate organized help message
	//var varCmds []cobra.Command
	//var cmdCmds []cobra.Command
	//for _, c := range rootCmd.Commands() {
	//	cmdName := c.Name()
	//	if string(cmdName[len(cmdName)-1]) == "c" {
	//		fmt.Println("cmd")
	//		cmdCmds = append(cmdCmds, *c)
	//	}
	//	if string(cmdName[len(cmdName)-1]) == "v" {
	//		fmt.Println("var")
	//		varCmds = append(varCmds, *c)
	//	}
	//}

	////for variable related Commands
	//for _, c := range varCmds {
	//	fmt.Println(c.Use, c.Short)
	//}
	//for _, c := range cmdCmds {
	//	fmt.Println(c.Use, c.Short)
	//}

	rootCmd.SetUsageFunc(func(cmd *cobra.Command) error {
		fmt.Println("Available Commands:")
		fmt.Println("         help [<command>]        help about any command")
		fmt.Println("")
		fmt.Println("   [Variable related]")
		fmt.Println("         av <variable> <value>   add <variable> <value> pair")
		fmt.Println("         cv <variable>           copy <variable>'s <value> to system clipboard")
		fmt.Println("         rv <variable>...        remove variables. If no <variable> provided, remove all variables")
		fmt.Println("         sv [<search_query>]     search variable that contains <search_query>. If no query provided, print all variables")
		fmt.Println("")
		fmt.Println("   [Command related]")
		fmt.Println("         ac [<key>]              add command with specified key. Key is optional")
		fmt.Println("         cc <key>                copy a command replacing <variable> with its <value>")
		fmt.Println("         rc <key>...             remove commands specified by keys")
		fmt.Println("         sc [<search_query>]     search commands. If no query, get all commands")
		fmt.Println("")
		fmt.Println("Available Flags:")
		fmt.Println("         kop [<command>] -h     help about any command")
		return nil
	})
}

func configInit() {
	if _, err := os.Stat(varFile_path); os.IsNotExist(err) {
		f, err := os.Create(varFile_path)
		if err != nil {
			log.Fatalln(err)
		}
		defer f.Close()
		fmt.Println("[+] Initialized ~/.kopvar.json")
	}

	if _, err := os.Stat(cmdFile_path); os.IsNotExist(err) {
		f, err := os.Create(cmdFile_path)
		if err != nil {
			log.Fatalln(err)
		}
		defer f.Close()
		fmt.Println("[+] Initialized ~/.kopcmd.json")
	}
}
