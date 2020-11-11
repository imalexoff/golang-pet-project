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
	maxThreadsCountStr := *flag.String("max-threads-count", "5", "Max threads count per group")
	flag.Parse()

	logf("Groups: %v", groupsStr)
	logf("Max threads: %v", maxThreadsCountStr)

	groups := strings.Split(groupsStr, ",")
	maxThreadsCount, _ := strconv.Atoi(maxThreadsCountStr)

	err := parse(groups, maxThreadsCount)
	if err != nil {
		log(err.Error())
	}

	fmt.Print("Press any key for exit...")
	fmt.Scan()

	os.Exit(0)
}
