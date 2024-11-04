package main

import (
	"bytes"
	"encoding/json"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

const ldPostLinkAPI = "api/bookmarks/"

func main() {
	start := time.Now()
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	handleFlags()

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(*givenFeed)
	if err != nil {
		log.Fatalln("error parsing feed URL", *givenFeed, "with the error:", err)
	}
	var ldLink string
	if strings.HasSuffix(*givenLinkdingURL, "/") {
		ldLink = *givenLinkdingURL + ldPostLinkAPI
	} else {
		ldLink = *givenLinkdingURL + "/" + ldPostLinkAPI
	}
	if *d {
		log.Println("using", ldLink, "as URL to POST a link to")
	}
	for i := 0; i < len(feed.Items); i++ {
		fi := feed.Items[i]
		var ld LinkdingURL
		title := strings.ReplaceAll(fi.Title, "\n", "")
		for strings.Contains(title, "  ") {
			title = strings.ReplaceAll(title, "  ", " ")
		}
		ld.Title = html.UnescapeString(title)
		ld.Link = fi.Link
		ld.IsArchived = *givenIsArchived
		ld.IsUnread = *givenIsUnread
		ld.IsShared = *givenIsShared
		var ld.TagNames []string
		if *givenTag != "" {
			ld.TagNames = append(ld.TagNames, *givenTag)
		}
		if *ddd {
			log.Println("Link object:", ld)
		}
		ldJSON, err := json.Marshal(ld)
		if err != nil {
			log.Fatalln(err)
		}
		if *dd {
			log.Println("Link as JSON", string(ldJSON))
		}
		err = postLink(ldLink, ldJSON)
		if err != nil {
			log.Println("error while posting link:", err)
		}
	}

	if *d {
		log.Printf("printElapsedTime: time elapsed %.2fs\n", time.Since(start).Seconds())
	}
}

func postLink(ldLink string, ldLinkJSON []byte) error {
	req, err := http.NewRequest("POST", ldLink, bytes.NewBuffer(ldLinkJSON))
	req.Header.Set("Authorization", "Token "+*givenLinkdingToken)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if *dd {
		log.Println("response Status:", resp.Status)
		log.Println("response Headers:", resp.Header)
		log.Println("response Body:", string(body))
	}
	return nil
}
