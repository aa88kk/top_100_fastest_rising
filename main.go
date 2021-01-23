package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/v33/github"
	"io/ioutil"
	"strings"
	"time"
)

var (
	client *github.Client
)

func output(lang string, cnt int, pastDays int) {

	now := time.Now()
	since := now.AddDate(0, 0, -pastDays)
	opts := &github.SearchOptions{ListOptions: github.ListOptions{PerPage: 100}, Sort: "stars", Order: "desc"}
	//query := "language:Go stars:10..1000"

	query := fmt.Sprintf("language:%s created:>%s", lang, since.Format("2006-01-02"))
	fmt.Println("query:", query)
	results, _, err := client.Search.Repositories(context.Background(), query, opts)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	var csvData []byte
	var line string

	for _, repo := range results.Repositories {
		if repo.Description != nil {
			line = fmt.Sprintf("%d, %s, %s, %q, %q, %v\n", *repo.StargazersCount, *repo.Name, *repo.HTMLURL, *repo.Description, strings.Join(repo.Topics, ","), repo.GetCreatedAt())
		} else {
			line = fmt.Sprintf("%d, %s, %s, \"\", %q, %v\n", *repo.StargazersCount, *repo.Name, *repo.HTMLURL, strings.Join(repo.Topics, ","), repo.GetCreatedAt())
		}
		csvData = append(csvData, []byte(line)...)
	}

	filename := fmt.Sprintf("top_%d_fastest_rising_since_%s", cnt, since.Format("2006-01-02"))
	ioutil.WriteFile("./data/"+filename+".csv", csvData, 0622)
	// for _, topic := range topics.Topics {
	// 	fmt.Println(*topic.Name)
	// }

	var mdData []byte
	var desc string
	mdData = []byte("| Stars | Name | Desc | Created | \n")
	mdData = append(mdData, []byte("| ----- | ------- | ------------- | --------- |\n")...)

	for _, repo := range results.Repositories {
		if repo.Description != nil {
			desc = strings.Replace(*repo.Description, "|", "\\|", -1)
			line = fmt.Sprintf("| %d | [%s](%s) | %q | %v |\n", *repo.StargazersCount, *repo.Name, *repo.HTMLURL, desc, repo.GetCreatedAt())
		} else {
			line = fmt.Sprintf("| %d | [%s](%s) |  | %v |\n", *repo.StargazersCount, *repo.Name, *repo.HTMLURL, repo.GetCreatedAt())
		}
		mdData = append(mdData, []byte(line)...)
	}

	mdData = append(mdData, []byte("\n")...)

	mdFilename := fmt.Sprintf("%s_top_%d_past_%d_days.md", lang, cnt, pastDays)
	ioutil.WriteFile("./"+mdFilename, mdData, 0622)

}

func main() {

	client = github.NewClient(nil)
	// opts := &github.SearchOptions{Sort: "stars", Order: "desc"}

	output("Go", 100, 30)
}
