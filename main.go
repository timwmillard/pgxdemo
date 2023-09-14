package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
)

func main() {
	printJSON := flag.Bool("json", false, "print output as json")
	flag.Parse()

	ctx := context.Background()

	// Database setup
	var connString string
	if flag.NArg() > 0 {
		connString = fmt.Sprintf("database=%s", flag.Arg(0))
	}
	db, err := DBConnect(ctx, connString)
	if err != nil {
		return
	}

	// -----------

	people, err := GetPeople(ctx, db)
	if err != nil {
		slog.Error("GetPerson error", "error", err)
		return
	}
	if *printJSON {
		err = json.NewEncoder(os.Stdout).Encode(people)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		printPeople(people)
	}

	fmt.Println()

	companies, err := GetCompanies(ctx, db)
	if err != nil {
		slog.Error("GetCompanies error", "error", err)
		return
	}
	if *printJSON {
		err = json.NewEncoder(os.Stdout).Encode(companies)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		printCompanies(companies)
	}
}

func printPeople(people []Person) {
	fmt.Println("--- People ---")
	for _, person := range people {
		fmt.Println(person)
	}
}
func printCompanies(companies []Company) {
	fmt.Println("--- Companies ---")
	for _, company := range companies {
		fmt.Println(company)
	}
}
