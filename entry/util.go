package entry

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

func getXSRFToken(doc *goquery.Document) (string, error) {

	csrf := doc.Find(`meta[name="_csrf_token"]`)

	token, e := csrf.First().Attr("content")

	if e {
		return token, nil
	} else {
		return "", fmt.Errorf("error finding xsrf token")
	}

}
