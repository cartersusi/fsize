package main

import (
	"flag"
	"fmt"
	"os/exec"

	sh "github.com/cartersusi/script_helper"
)

var version string

func GetFlag[T string | int | bool](long_flag, short_flag *T, value T, required bool, flag_name ...string) T {
	flagstr := "flags"
	if len(flag_name) > 0 {
		flagstr = flag_name[0]
	}

	if *long_flag == value && *short_flag == value {
		if required {
			sh.Error(fmt.Sprintf("Missing values for %s. Please provide one.", flagstr), true)
		}
	}

	if *long_flag != value && *short_flag != value {
		if required {
			sh.Error(fmt.Sprintf("Both %s are set. Please provide only one.", flagstr), true)
		}
	}

	if *long_flag != value {
		return *long_flag
	}

	if *short_flag != value {
		return *short_flag
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
		//find . -type f -exec du -sh {} + | sort -rh | head -n 5
		cmd = fmt.Sprintf("find %s -type f -exec du -sh {} + | sort -rh | head -n %d", directory, head_n)
	} else {
		//ls -lSha | head -n 10
		cmd = fmt.Sprintf("du -sh %s/* | sort -rh | head -n %d", directory, head_n)
	}

	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		sh.Error(err.Error(), true)
	}

	fmt.Println("\n" + string(out))
}
