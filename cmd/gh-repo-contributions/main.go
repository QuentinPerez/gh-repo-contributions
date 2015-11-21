package main

import (
	"fmt"
	"os"

	"github.com/QuentinPerez/gh-repo-contributions"
	"github.com/Sirupsen/logrus"
)

func main() {
	if len(os.Args) != 2 {
		logrus.Fatalf("Usage: %v <GithubName>", os.Args[0])
	}
	result, err := github.GetAllContibutions(os.Args[1])
	if err != nil {
		logrus.Fatal(err)
	}
	maxlen := 0
	for _, c := range result {
		if len(c.Repo)+18 > maxlen {
			maxlen = len(c.Repo) + 18
		}
	}
	fmt.Printf("%-*s | %s\n", maxlen, "Repositories", "Commits")
	for i := 0; i <= maxlen; i++ {
		fmt.Printf("-")
	}
	fmt.Printf("|")
	for i := 0; i < 10; i++ {
		fmt.Printf("-")
	}
	fmt.Println()
	total := 0
	for _, c := range result {
		// FIXME Don't use + operator to concat strings
		fmt.Printf("%-*s | %d\n", maxlen, "http://github.com/"+c.Repo, c.Commits)
		total += c.Commits
	}
	fmt.Printf("%-*s | %d\n", maxlen, "", total)
	fmt.Println()
}
