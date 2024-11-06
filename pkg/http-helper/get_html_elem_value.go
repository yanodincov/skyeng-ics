package httphelper

import (
	"bytes"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
)

func GetHTMLElemValue(body []byte, selector string) (string, error) {
	// Load HTML document
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return "", errors.Wrap(err, "failed to load HTML document")
	}

	// Find the input element by the selector and retrieve its value
	value, exists := doc.Find(selector).Attr("value")
	if !exists {
		return "", errors.Wrapf(err, "failed to find the value of the selector: %s", selector)
	}

	return value, nil
}
