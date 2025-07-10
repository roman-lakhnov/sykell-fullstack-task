package analyzer

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

// StartAnalyzerWorker runs a background worker that continuously processes links
func StartAnalyzerWorker() {
	go func() {
		for {
			// Find one record with "created" status
			link, id, found, err := getNextLinkToProcess()
			if err != nil {
				fmt.Printf("Error retrieving link: %v\n", err)
				time.Sleep(1 * time.Second)
				continue
			}

			if !found {
				// No links to process, wait and check again
				time.Sleep(1 * time.Second)
				continue
			}

			// Process the link
			err = processLink(link, id)
			if err != nil {
				fmt.Printf("Error processing link %s: %v\n", link, err)
				time.Sleep(1 * time.Second)
			}
		}
	}()
}

// getNextLinkToProcess retrieves the next link with "created" status
func getNextLinkToProcess() (string, int, bool, error) {
	row := DB.QueryRow("SELECT id, url FROM results WHERE check_status = 'created' LIMIT 1")
	var id int
	var url string
	err := row.Scan(&id, &url)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return "", 0, false, nil // No links found, but no error
		}
		return "", 0, false, err // Actual error occurred
	}

	return url, id, true, nil
}

// processLink analyzes a single link and updates its record in the database
func processLink(link string, id int) error {
	updateInDB(
		id,
		link,
		"pending",
		time.Now(),
		"",
		"",
		map[string]int{},
		0, 0, 0,
		false,
		nil, // No inaccessible links at this point
	)

	resp, err := http.Get(link)
	if err != nil || resp.StatusCode >= 400 {
		// Update status to error in the database
		return updateInDB(
			id,
			link,
			"error",
			time.Now(),
			"",
			"",
			map[string]int{},
			0, 0, 0,
			false,
			[]LinkIssue{{URL: link, StatusCode: resp.StatusCode}},
		)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return updateInDB(
			id,
			link,
			"error",
			time.Now(),
			"",
			"",
			map[string]int{},
			0, 0, 0,
			false,
			[]LinkIssue{{URL: link, StatusCode: -1}}, // Connection error
		)
	}

	parsedURL, _ := url.Parse(link)

	// HTML version
	htmlVersion := "HTML5" // Default to HTML5

	// Check DOCTYPE to determine HTML version
	docTypeNode := doc.Find("html").Nodes
	if len(docTypeNode) > 0 && docTypeNode[0].PrevSibling != nil && docTypeNode[0].PrevSibling.Type == html.DoctypeNode {
		docType := docTypeNode[0].PrevSibling.Data

		switch {
		case strings.Contains(docType, "HTML 4.01"):
			htmlVersion = "HTML 4.01"
		case strings.Contains(docType, "XHTML 1.0"):
			htmlVersion = "XHTML 1.0"
		case strings.Contains(docType, "XHTML 1.1"):
			htmlVersion = "XHTML 1.1"
		case strings.Contains(docType, "HTML 3.2"):
			htmlVersion = "HTML 3.2"
		case strings.Contains(docType, "HTML 2.0"):
			htmlVersion = "HTML 2.0"
		}
	} else {
		// No doctype, look for HTML5 indicators
		if doc.Find("[data-*]").Length() > 0 ||
			doc.Find("article, section, nav, header, footer, aside").Length() > 0 {
			htmlVersion = "HTML5"
		} else {
			htmlVersion = "Unknown"
		}
	}

	// Title
	title := strings.TrimSpace(doc.Find("title").Text())

	// Heading counts
	headings := map[string]int{}
	for i := 1; i <= 6; i++ {
		tag := fmt.Sprintf("h%d", i)
		headings[strings.ToUpper(tag)] = doc.Find(tag).Length()
	}

	// Links
	internalLinks := 0
	externalLinks := 0
	inaccessibleLinks := 0
	var inaccessibleLinksDetails []LinkIssue

	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		linkParsed, err := url.Parse(href)
		if err != nil {
			return
		}
		finalURL := parsedURL.ResolveReference(linkParsed).String()

		if linkParsed.Host == "" || linkParsed.Host == parsedURL.Host {
			internalLinks++
		} else {
			externalLinks++
		}

		// Skip checking javascript: links, mailto: links, tel: links and fragments only
		if strings.HasPrefix(linkParsed.Scheme, "javascript") ||
			strings.HasPrefix(linkParsed.Scheme, "mailto") ||
			strings.HasPrefix(linkParsed.Scheme, "tel") ||
			(linkParsed.Path == "" && linkParsed.Fragment != "") {
			return
		}

		client := http.Client{
			Timeout: 5 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				// Allow up to 5 redirects
				if len(via) >= 5 {
					return http.ErrUseLastResponse
				}
				return nil
			},
		}

		req, err := http.NewRequest("HEAD", finalURL, nil)
		if err != nil {
			inaccessibleLinks++
			inaccessibleLinksDetails = append(inaccessibleLinksDetails, LinkIssue{
				URL:        finalURL,
				StatusCode: -1, // Connection error
			})
			return
		}

		// Add user agent to avoid getting blocked
		req.Header.Set("User-Agent", "Mozilla/5.0 Webpage Analyzer")

		r, err := client.Do(req)
		if err != nil {
			inaccessibleLinks++
			inaccessibleLinksDetails = append(inaccessibleLinksDetails, LinkIssue{
				URL:        finalURL,
				StatusCode: -1, // Connection error
			})
			return
		}
		defer r.Body.Close()

		// Mark as inaccessible if status code is 4xx or 5xx
		if r.StatusCode >= 400 {
			inaccessibleLinks++
			inaccessibleLinksDetails = append(inaccessibleLinksDetails, LinkIssue{
				URL:        finalURL,
				StatusCode: r.StatusCode,
			})
		}
	})

	// Login form detection
	loginForm := false
	doc.Find("form").EachWithBreak(func(i int, form *goquery.Selection) bool {
		if form.Find("input[type='password']").Length() > 0 {
			loginForm = true
			return false
		}
		return true
	})

	// Update the record in DB
	return updateInDB(
		id,
		link,
		"checked",
		time.Now(),
		title,
		htmlVersion,
		headings,
		internalLinks,
		externalLinks,
		inaccessibleLinks,
		loginForm,
		inaccessibleLinksDetails,
	)
}

