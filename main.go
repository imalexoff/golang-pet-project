package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	groups := *flag.String("groups", "3312,3315,3330,3475,3362,3367,3403,3404,3405", "Goods groups")
	flag.Parse()

	if groups == "" {
		log("Need to provide groups")

		fmt.Print("Press any key for exit...")
		fmt.Scan()

		os.Exit(0)
	}

	logf("Groups: %v", groups)

	err := getGoods(groups)
	if err != nil {
		log(err.Error())
	}

	fmt.Print("Press any key for exit...")
	fmt.Scan()

	os.Exit(0)
}
