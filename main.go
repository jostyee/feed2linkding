package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"net/url"
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

	c = http.Client{Timeout: time.Duration(5) * time.Second}

	excludedDomains     = make(map[string]struct{})
	createAPI, checkAPI = "", ""
)

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

	baseURL := strings.TrimSuffix(*givenLinkdingURL, "/")
	createAPI = baseURL + "/api/bookmarks/"
	checkAPI = baseURL + "/api/bookmarks/check/?url="

	if *d {
		log.Println("using", createAPI, "as URL to POST a link to")
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
		excluded, err := excludedDomain(ld.Link)
		if err != nil {
			log.Fatalln(err)
		}
		if excluded {
			if *d {
				log.Println("domain excluded: ", ld.Link)
			}
			continue
		}

		existed, err := checkLink(ld.Link)
		if err != nil {
			log.Fatalln(err)
		}
		if existed {
			if *d {
				log.Println("link exists in the database, skipping: ", ld.Link)
			}

			continue
		}

		ld.IsArchived = *givenIsArchived
		ld.IsUnread = *givenIsUnread
		ld.IsShared = *givenIsShared

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

		if err = postLink(ldJSON); err != nil {
			log.Println("error while posting link:", err)
		}
	}

	if *d {
		log.Printf("printElapsedTime: time elapsed %.2fs\n", time.Since(start).Seconds())
	}
}

func postLink(ldLinkJSON []byte) error {
	req, err := http.NewRequest("POST", createAPI, bytes.NewBuffer(ldLinkJSON))
	req.Header.Set("Authorization", "Token "+*givenLinkdingToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.Do(req)
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

type CheckResponse struct {
	Bookmark any `json:"bookmark"`
}

func excludedDomain(link string) (bool, error) {
	parsedUrl, err := url.Parse(link)
	if err != nil {
		return false, fmt.Errorf("Error parsing URL:", err)
	}

	domain := parsedUrl.Hostname()
	_, ok := excludedDomains[domain]
	return ok, nil
}

// checkLink checks if a link is already in the Linkding database
// returns true if the link is already in the database, false otherwise
func checkLink(link string) (bool, error) {
	req, err := http.NewRequest("GET", checkAPI+url.QueryEscape(link), nil)
	req.Header.Set("Authorization", "Token "+*givenLinkdingToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	if *dd {
		log.Println("response Status:", resp.Status)
		log.Println("response Headers:", resp.Header)
		log.Println("response Body:", string(body))
	}

	var response CheckResponse
	if err := json.Unmarshal([]byte(body), &response); err != nil {
		return false, err
	}

	if response.Bookmark != nil {
		return true, nil
	}

	return false, nil
}
