package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var givenFeed = flag.String("feed", "", "RSS or Atom feed link")
var givenLinkdingURL = flag.String("ld-url", "", "URL to your linkding instance")
var givenLinkdingToken = flag.String("ld-token", "", "linkding token from your settings under the path $LinkdingURL/settings/integrations")
var givenIsArchived = flag.Bool("is-archived", true, "shall the URL be archived in linkding?")
var givenIsUnread = flag.Bool("unread", false, "shall the URL be unread in linkding?")
var givenIsShared = flag.Bool("shared", false, "shall the URL be shared in linkding?")
var givenTag = flag.String("tag", "", "give a tag name which gets added to that link")
var v = flag.Bool("v", false, "print the program's version")
var d = flag.Bool("d", false, "print debug output")
var dd = flag.Bool("dd", false, "like parameter -d, but print more debug output")
var ddd = flag.Bool("ddd", false, "like parameter -dd, but print even more debug output")

func handleFlags() {
	flag.Parse()
	if *v {
		fmt.Printf("version %v\n", version)
		os.Exit(0)
	}

	if *givenFeed == "" {
		log.Fatalln("providing a URL to a feed is a must, please use parameter -feed for this")
		os.Exit(1)
	}
	if *givenLinkdingURL == "" {
		log.Fatalln("providing the linkding URL is a must, please use parameter -ld-url for this")
		os.Exit(1)
	}
	if *givenLinkdingToken == "" {
		log.Fatalln("providing the linkding token is a must, please use parameter -ld-token for this")
		os.Exit(1)
	}

	if *ddd {
		*dd = true
		*d = true
		log.Println("Debug Level ddd is active")
	} else if *dd {
		*d = true
		log.Println("Debug Level dd is active")
	} else if *d {
		log.Println("Debug Level d is active")
	}
}
