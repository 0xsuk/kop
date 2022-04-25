package cmd

import (
	"fmt"
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
		varmap, err := File2Map(varFile_path)
		CheckErr(err)

		//iterate over jsonValue
		for name, value := range varmap {
			fmt.Printf("%v: %v\n", name, value)
		}
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
  rootCmd.SetUsageFunc(func(cmd *cobra.Command) error {
    fmt.Println(`Available Commands:
                 help [<command>]       help about any command
    [Variable]
                 av <variable> <value>  add <variable> <value> pair
                 cv <variable>          copy <variable>'s <value> to system clipboard
                 rv <variable>...       remove variables. If no <variable> provided, remove all variables
                 sv [<search_query>]    search variable that contains <search_query>. If no query provided, print all variables
    [Command]
                 ac                     add command
                 cc <index>             copy command replacing <variable> with its <value>
                 rc <index>...          remove commands specified by index
                 sc [<search_query>]    search commands. If no query provided, print all commands
                 
Available Flags:
                 kop [<command>] -h     help about any command
                `)
    return nil
  })
	rootCmd.AddCommand(searchvar)
  //searchvar.SetHelpFunc(func (*cobra.Command, []string) {})
  searchvar.SetUsageFunc(func(*cobra.Command) error {return nil})
	rootCmd.AddCommand(copyvar)
	rootCmd.AddCommand(addvar)
	rootCmd.AddCommand(removevar)
	rootCmd.AddCommand(addcmd)
	rootCmd.AddCommand(searchcmd)
	rootCmd.AddCommand(removecmd)
	rootCmd.AddCommand(copycmd)
}

func configInit() {
	if _, err := os.Stat(varFile_path); os.IsNotExist(err) {
		f, err := os.Create(varFile_path)
		CheckErr(err)
		defer f.Close()
		fmt.Println("[+] Initialized ~/.kopvar.json")
	}

	if _, err := os.Stat(cmdFile_path); os.IsNotExist(err) {
		f, err := os.Create(cmdFile_path)
		CheckErr(err)
		defer f.Close()
		fmt.Println("[+] Initialized ~/.kopcmd.json")
	}
}
