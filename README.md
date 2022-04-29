# kop
kop is a cli tool that I personally use for daily CTF. 
<br/><br/><br/>
**Concept: Do not type same commands**
<br/><br/><br/>
If you are CTF player, you are probably typing same shell commands everyday, like `nmap <ipaddress>`, or `gobuster dir --url　http://<ipaddress>`... and so on.  
The problem is that these commands slightly change every time, because of the "variable" factors, like IP address.

kop was made to address this problem. 

<br/><br/>
Here's a quick example usage.  
![image](https://github.com/0xsuk/kop/blob/main/.github/example1.png)    
  
brief explanation:  
- `kop ac`: `ac` stands for "add command".   
Adding nmap command with IP variable ${IP}.  ${} is a variable notation.    
- `kop av`: `av` stands for "add variable".    
Setting IP variable to 10.11.12.13.
- `kop sc`: `sc` stands for "search command".   
Searching commands that contains string "nmap", and kop shows match and its index which turns to be 0 (since this is the first command added)
- `kop cc`: `cc` stands for "copy command".   
Copying a command at index 0, to my **system clipboard**.  

The nmap command was copied to my clipboard, **replacing IP variable with its value.**  
You can paste the command, without typing loooong shell command again.

Commands and variables added by `ac` and `av` remain in json files, so you don't have to type them again.  

# Install
You must have `xclip` for Linux, `pbcopy` for OSX to make full use of kop.  
```
go install github.com/0xsuk/kop@latest
```

Or install from [releases](https://github.com/0xsuk/kop/releases)

# Usage

```
Available Commands:
         help [<command>]        help about any command

   [Variable related]
         av <variable> <value>   add <variable> <value> pair
         cv <variable>           copy <variable>'s <value> to system clipboard
         rv <variable>...        remove variables. If no <variable> provided, remove all variables
         sv [<search_query>]     search variable that contains <search_query>. If no query provided, print all variables

   [Command related]
         ac [<key>]              add command with specified key. Key is optional
         cc <key|index>          copy a command replacing <variable> with its <value>
         rc <key|index>...       remove commands specified by keys or ids
         sc [<search_query>]     search commands. If no query, get all commands

Available Flags:
         kop [<command>] -h     help about any command
```

<details>
<summary>Detail</summary>
  <li> &lt;key&gt;: Key is a string id that can be used as an argument to specify a command.  </li>
 <li> &lt;key|index&gt;: Key or Index. Index is an index of a command. All commands are assigned an index when added with `kop ac`, no matter &lt;key&gt; is provided or not. Check indexes of commands by `kop sc`.  
  </li>
</details>

CTF Life is so much better with kop!  
![image](https://github.com/0xsuk/kop/blob/main/.github/example3.png)  

# TODO
- [x] Adding command with string id
- [ ] Incremental Searching
- [ ] Executing command with kop


