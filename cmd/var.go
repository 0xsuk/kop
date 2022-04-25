package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var searchvar = &cobra.Command{
	Use:   "sv [<search_query>]",
	Short: "search variable that contains <search_query>. If no query provided, print all variables",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("Too many arguments to call gv.")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		varmap, err := File2Map(varFile_path)
		CheckErr(err)

		var searchq string = ""
		if len(args) == 1 {
			searchq = args[0]
		}

		var isfound bool = false
		for varname := range varmap {
			if strings.Contains(varname, searchq) {
				fmt.Println(varname + ": " + varmap[varname])
				isfound = true
			}
		}
		if !isfound {
			fmt.Printf("Variable not found for search_query=%q\n", searchq)
		}
	},
}

var copyvar = &cobra.Command{
	Use:   "cv <variable>",
	Short: "copy <variable>'s <value> to system clipboard",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("Please specify variable name")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		searchvar := args[0]
		varmap, err := File2Map(varFile_path)
		CheckErr(err)

		var isfound bool = false
		for varname := range varmap {
			if varname == searchvar {
				isfound = true
				fmt.Println(varname + ": " + varmap[varname])
				fmt.Println(Copy(varmap[varname]))
			}
		}
		if !isfound {
			fmt.Printf("No variable found for variable=%q\n", searchvar)
		}
	},
}

var addvar = &cobra.Command{
	Use:   "av <variable> <value>",
	Short: "add <variable> <value> pair",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("Too few argument to call add")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		varmap, err := File2Map(varFile_path)
		CheckErr(err)

		varmap[args[0]] = args[1]

		err = WriteMap(varmap, varFile_path)
		CheckErr(err)

		fmt.Println("[+] Successfully Added")
	},
}

//resetvar resests varFile_path
var removevar = &cobra.Command{
	Use:   "rv <variable>...",
	Short: "remove variables. If no <variable> provided, remove all variables",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Printf("Are you sure you want to reset all variables?[Y/n]")
			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				if scanner.Text() == "n" {
					fmt.Println("[-] Stopped rv")
					return
				}
				if scanner.Text() != "Y" && scanner.Text() != "" {
					fmt.Println("[*] Type Y or nothing to reset all, n to stop.")
					cmd.Run(cmd, args)
				}
			}
			varFile, err := os.Create(varFile_path)
			CheckErr(err)
			defer varFile.Close()
			fmt.Println("[+] Successfully reset")
			return
		}

		varmap, err := File2Map(varFile_path)
		CheckErr(err)
		new_varmap := map[string]string{}
		var isfound bool = false
		var removed_var string = ""
		for varname := range varmap {
			if !Contains(args, varname) {
				new_varmap[varname] = varmap[varname]
			} else {
				isfound = true
				removed_var = varname
			}
		}
		if !isfound {
			fmt.Println("[*] Nothing removed")
			return
		}
		err = WriteMap(new_varmap, varFile_path)
		CheckErr(err)
		fmt.Printf("[+] Successfully removed: %+v\n", removed_var)
	},
}
