package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

var IsLetter = regexp.MustCompile(`^[a-zA-Z]+$`).MatchString

type systemInfo struct {
	WorkedLoad int `json:"cpu"`
	InWorkCPU  int `json:"memory"`
	TotalCPU   int `json:"memory_total"`
}

func findRAM() int {
	app := "ps"
	arg0 := "-eo"
	arg1 := "%cpu"
	arg2 := "--sort"
	arg3 := "-%cpu"
	errorValue := -1

	cmd := exec.Command(app, arg0, arg1, arg2, arg3)
	stdout, err := cmd.Output()

	if err != nil {
		println(err.Error())
		return errorValue
	}

	stringArray := string(stdout)
	RamPercent := 0
	countOfProcess := 0
	stringTmp := ""
	codeOfSpace := 32
	flag := false

	for i := 1; i < len(stringArray); i++ {
		if (string(stringArray[i]) != "0.0") && stringArray[i] != uint8(codeOfSpace) && (!IsLetter(string(stringArray[i]))) && string(stringArray[i]) != "\n" {
			stringTmp += string(stringArray[i])
			flag = true
		}
		if string(stringArray[i]) == "\n" && flag {
			tmp := 0.0
			tmp, err = strconv.ParseFloat(stringTmp, 64)
			RamPercent += int(tmp)
			countOfProcess++
			stringTmp = ""
			flag = false
		}
	}

	// Calculated Percent of RAM
	return RamPercent
}

func findNumber(s string) int {
	re := regexp.MustCompile("[0-9]+")
	totalMemory, _ := strconv.Atoi(strings.Join(re.FindAllString(s, -1), " "))
	return totalMemory
}

func findRAMPercent() (int, int) {
	cmd := exec.Command("bash", "-c", "vmstat -s | grep 'total memory'")
	stdout, err := cmd.Output()
	errorValue := -1

	if err != nil {
		println(err.Error())
		return errorValue, errorValue
	}

	totalMemory := findNumber(string(stdout))

	cmd = exec.Command("bash", "-c", "vmstat -s | grep 'used memory'")
	stdout, err = cmd.Output()

	if err != nil {
		println(err.Error())
		return errorValue, errorValue
	}

	usedMemory := findNumber(string(stdout))

	println("Percentage of RAM in work: ", int(float32(usedMemory)/float32(totalMemory)*100), "%")

	println("RAM of System: ", int(float32(totalMemory)*0.000977), "MB")

	return int(float32(usedMemory) / float32(totalMemory) * 100), int(float32(totalMemory) * 0.000977)

}

func takeInfo(w http.ResponseWriter, r *http.Request) {

	cpuPercent := findRAM()
	print("Processor load percentage: ", cpuPercent, "%\n")

	ramInWork, totalRam := findRAMPercent()

	info := systemInfo{}
	info.WorkedLoad = cpuPercent
	info.InWorkCPU = ramInWork
	info.TotalCPU = totalRam

	jsonInfo, _ := json.Marshal(info)
	fmt.Println(string(jsonInfo))

	fmt.Fprintf(w, string(jsonInfo))
}

func main() {

	http.HandleFunc("/system/load", takeInfo)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
