package main

import (
	"flag"
)

func main() {
	var folder string
	var email string

	flag.StringVar(&folder, "add", "", "add a new folder to scan for Git repositries")
	flag.StringVar(&email, "email", "agoyal3_be21@thapar.edu", "the email to scan for")

	flag.Parse()

	if folder != "" {
		scan(folder)
		return
	}

	stats(email)

}
