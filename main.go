package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/v33/github"
	"io/ioutil"
	"strings"
	"time"
)

// Fetch all the public organizations' membership of a user.
//
func fetchOrganizations(username string) ([]*github.Organization, error) {
	client := github.NewClient(nil)
	orgs, _, err := client.Organizations.List(context.Background(), username, nil)
	return orgs, err
}

func writeData() {

}

func main() {

	now := time.Now()
	since := now.AddDate(0, 0, -90)
	// var username string
	// fmt.Print("Enter GitHub username: ")
	// fmt.Scanf("%s", &username)

	// organizations, err := fetchOrganizations(username)
	// if err != nil {
	// 	fmt.Printf("Error: %v\n", err)
	// 	return
	// }

	// for i, organization := range organizations {
	// 	fmt.Printf("%v. %v\n", i+1, organization.GetLogin())
	// }

	client := github.NewClient(nil)
	// opts := &github.SearchOptions{Sort: "stars", Order: "desc"}
	opts := &github.SearchOptions{ListOptions: github.ListOptions{PerPage: 100}, Sort: "stars", Order: "desc"}
	//query := "language:Go stars:10..1000"

	query := fmt.Sprintf("language:Go created:>%s", since.Format("2006-01-02"))
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

	ioutil.WriteFile("test.csv", csvData, 0622)
	// for _, topic := range topics.Topics {
	// 	fmt.Println(*topic.Name)
	// }
}
