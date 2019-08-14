package winenum

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)
//List of previously entered commands
func PSHistory(){
	commands := "Get-Content (Get-PSReadlineOption).HistorySavePath"
	out, err := exec.Command("powershell", commands).Output()
	if err != nil {
		fmt.Println(err)
	}
	result := string(out)
	fmt.Println(result)
}
// Display current $PATH and env information
func EnvInfo(){
	out, err := exec.Command("cmd", "/C", "SET | more").Output()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(out))

}
//Information about the name and version of the operating system, list of disks
func SystemInfo() {
	commands := []string {"systeminfo | findstr /c:'Host Name' /c:'OS Version'",
		"[System.Environment]::OSVersion.Version",
		"wmic logicaldisk get caption,description,providername"}
	for _,command := range commands{
		fmt.Println(command)
		out, err := exec.Command("powershell", command).Output()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(out))
	}
}
//Information about the current user, a list of other users, information about them and their rights in the system
func UserInfo(){
	fmt.Println("Current User Details and Other Users")
	commands := []string {"whoami /all", "query user", "wmic useraccount get name"}
	for i, command := range commands{
		fmt.Println("Command: ", command)
		out, err := exec.Command("powershell", command).Output()
		if err != nil {
			fmt.Println(err)
		}
		result := string(out)
		if i == 1{
			names := strings.Split(string(out), "\r\n")
			//fmt.Println(names, len(names))
			for _, user := range names{
				fmt.Println("net user " + user)
				UserDetail, err := exec.Command("powershell", "net user " + user).Output()
				if err != nil {
					fmt.Println(err)
				}
				//fmt.Println(string(UserDetail))
				result += string(UserDetail)
			}
		}
		fmt.Println(result)

	}
}

//obtaining a list of processes that are currently running on the computer, as well as a list of installed software
func Services(){
	fmt.Println("Services/Programs:")
	DefProcesses := []string {"audiodg.exe", "conhost.exe",	"csrss.exe", "lsass.exe",
		"lsm.exe", "MSCamS64.exe","naPrdMgr.exe","OSPPSVC.EXE",	"PresentationFontCache.exe",
		"SearchIndexer.exe","services.exe",	"smss.exe",	"spoolsv.exe",	"svchost.exe",
		"svchost.exe","svchost.exe","System","System Idle Process",	"UNS.exe","wininit.exe",
		"WmiApSrv.exe",	"WmiPrvSE.exe",	"wmpnetwk.exe",	"WUDFHost.exe"}
	commands := []string {"tasklist /v", "Get-ChildItem 'C:" + `\` +
		"Program Files', 'C:" + `\` + "Program Files (x86)' | ft Parent,Name,LastWriteTime"}
	for _, command := range commands{
		fmt.Println("Command: ", command)
		out, err := exec.Command("powershell", command).Output()
		if err != nil {
			fmt.Println(err)
		}
		result := strings.Split(string(out), "\r\n")
		fmt.Println(result, len(result))
		for i, subresult := range result {
			for _, proc := range DefProcesses{
				if strings.Contains(subresult, proc){
					result[i] = result[len(result)-1]
					result[len(result)-1] = ""
					result = result[:len(result)-1]
				}
			}
		}
		for _, subresult :=range result{
			fmt.Println(subresult)
		}

	}
}
//Scheduled Tasks List
func JobsInfo(){
	fmt.Println("Jobs/Tasks:")
	commands := []string {"SCHTASKS /Query /fo LIST", "Get-ScheduledTask"}
	for _, command := range commands{
		fmt.Println("Command: ", command)
		out, err := exec.Command("powershell", command).Output()
		if err != nil {
			fmt.Println(err)
		}
		result := string(out)
		fmt.Println(result)
	}
}

type ProcessName struct {
	Port int
	State int
	ProcName string
}
//List of network interface parameters, ARP table, DNS settings, processes and their ports
func NetworkInfo(){
	fmt.Println("Network:")
	commands := []string {"ipconfig /all", "arp -a", "Get-DnsClientServerAddress",
		"Get-NetTCPConnection | Select-Object @{n='Port';e={$_.LocalPort}},@{n='State';e={$_.State}},@{n='ProcName';e={(Get-Process -Id $_.OwningProcess).ProcessName}} | ConvertTo-Json -Compress",
	}
	for _, command := range commands{
		fmt.Println("Command: ", command)
		out, err := exec.Command("powershell", command).Output()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(out))
		var ProcessNames []ProcessName
		var ParseJsonStruct ProcessName
		if command == "Get-NetTCPConnection | Select-Object @{n='Port';e={$_.LocalPort}},@{n='State';e={$_.State}},@{n='ProcName';e={(Get-Process -Id $_.OwningProcess).ProcessName}} | ConvertTo-Json -Compress"{
			result := bytes.SplitAfter(out, []byte("}"))
			for _, subresult := range result{
				err := json.Unmarshal(subresult[1:], &ParseJsonStruct)
				if err != nil{
					fmt.Println(err)
				}
				ProcessNames = append(ProcessNames, ParseJsonStruct)

			}
			fmt.Println(ProcessNames)
		}

	}

}
/*
func Searches(){
	fmt.Println("Search any reg file with pass:")
	commands := []string {"REG QUERY HKCU /F " + "password " + "/t REG_SZ /S /K"}
	for _, command := range commands{
		fmt.Println("Command: ", command)
		out, err := exec.Command("powershell", command).Output()
		if err != nil {
			fmt.Println(err)
		}
		result := string(out)
		fmt.Println(result)
	}
}
*/
func WinEnumRun(){
	SystemInfo()
	UserInfo()
	EnvInfo()
	PSHistory()
	Services()
	JobsInfo()
	NetworkInfo()
}
