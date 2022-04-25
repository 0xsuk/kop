# kop
kop is a cli tool that I personally use for daily CTF. 

**Concept: Do not type same commands**
  
If you are CTF player, you are probably typing same commands everyday, like `nmap <ipaddress>`, or `gobuster dir --url　http://<ipaddress>`... and so on.  
The problem is that these commands slightly change every time, because of the "variable" factors, like IP address.

kop was made to address this problem. 


Here's a quick example usage.  
![image](https://github.com/0xsuk/kop/blob/main/.github/example.png)  
brief explanation:  
- `kop ac`: `ac` stands for "add command".   
Adding nmap command with IP variable.    
- `kop av`: `av` stands for "add variable".    
Setting IP variable to 10.11.12.13.
- `kop sc`: `sc` stands for "search command".   
Searching commands that contains string "nmap", and kop shows hit and its index (which is 0).   
- `kop cc`: `cc` stands for "copy command".   
Copying a command at index 0, to my **system clipboard**.  

and then the nmap command was copied to my clipboard, **replacing IP variable with its value.**

Commands and variables added by `ac` and `av`, will remain in json files, so you don't have to type them again.  

# Install
You must have `xclip` for Linux, `pbcopy` for OSX to make full use of kop.  
```
go install github.com/0xsuk/kop@latest
```

# Usage
```
Available Commands:
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
```
