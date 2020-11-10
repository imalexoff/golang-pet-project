package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	groupsStr := *flag.String("groups", "3312,3315,3330,3475,3362,3367,3403,3404,3405", "Goods groups")
	threadsCountStr := *flag.String("threads-count", "5", "Threads count for each group")
	flag.Parse()

	logf("Groups: %v", groupsStr)
	logf("Threads: %v", threadsCountStr)

	groups := strings.Split(groupsStr, ",")
	threadsCount, _ := strconv.Atoi(threadsCountStr)

	err := parse(groups, threadsCount)
	if err != nil {
		log(err.Error())
	}

	fmt.Print("Press any key for exit...")
	fmt.Scan()

	os.Exit(0)
}
