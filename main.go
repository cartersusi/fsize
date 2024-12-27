package main

import (
	"flag"
	"fmt"
	"os/exec"

	sh "github.com/cartersusi/script_helper"
)

var version string

func GetFlag[T string | int | bool](longFlag, shortFlag *T, value T, required bool, flagName ...string) T {
	flagstr := "flags"
	if len(flagName) > 0 {
		flagstr = flagName[0]
	}

	if *longFlag == value && *shortFlag == value {
		if required {
			sh.Error(fmt.Sprintf("Missing values for %s. Please provide one.", flagstr), true)
		}
	}

	if *longFlag != value && *shortFlag != value {
		if required {
			sh.Error(fmt.Sprintf("Both %s are set. Please provide only one.", flagstr), true)
		}
	}

	if *longFlag != value {
		return *longFlag
	}

	if *shortFlag != value {
		return *shortFlag
	}

	return value
}

func main() {
	dir := flag.String("dir", "", "Directory to find file sizes for")
	d := flag.String("d", "", "Directory to find file sizes for")

	recursive := flag.Bool("recursive", false, "Recursively find file sizes")
	r := flag.Bool("r", false, "Recursively find file sizes")

	n := flag.Int("n", 25, "Number of files/directories to list")

	v := flag.Bool("v", false, "Print version")
	_version := flag.Bool("version", false, "Print version")

	flag.Parse()
	if *v || *_version {
		fmt.Println("fsize | Version:", version)
		return
	}

	directory := GetFlag(dir, d, "", true, "directory")
	recurse := GetFlag(recursive, r, false, false, "recursive")
	head_n := *n

	sh.Success(fmt.Sprintf("Directory: %s", directory))
	sh.Success(fmt.Sprintf("Recursive: %t", recurse))
	sh.Success(fmt.Sprintf("Number: %d", head_n))

	var cmd string
	if recurse {
		// Includes hidden files by default since 'find' doesn’t exclude them
		cmd = fmt.Sprintf("find %s -type f -exec du -sh {} + | sort -rh | head -n %d", directory, head_n)
	} else {
		// Include hidden files by expanding .* and ignoring errors from . and ..
		cmd = fmt.Sprintf("du -sh %s/* %s/.* 2>/dev/null | sort -rh | head -n %d", directory, directory, head_n)
	}

	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		sh.Error(err.Error(), true)
	}

	fmt.Println("\n" + string(out))
}
