package analyzer

import (
    "github.com/PuerkitoBio/goquery"
    "golang.org/x/net/html"
    "strings"
    "io"
)

func ParseHTML(r io.Reader) (*goquery.Document, *html.Node, error) {
    root, err := html.Parse(r)
    if err != nil {
        return nil, nil, err
    }
    doc := goquery.NewDocumentFromNode(root)
    return doc, root, nil
}

// inspects the root node for a <!DOCTYPE ...> and return the HTML version.
func DetectHTMLVersion(root *html.Node) string {
    for n := root.FirstChild; n != nil; n = n.NextSibling {
        if n.Type == html.DoctypeNode {
            doctype := strings.ToLower(n.Data)

            if doctype == "html" && len(n.Attr) == 0 {
                return "HTML5"
            }

            for _, attr := range n.Attr {
                pub := strings.ToLower(attr.Val)

                switch {
                case strings.Contains(pub, "-//w3c//dtd html 4.01 transitional"):
                    return "HTML 4.01 Transitional"
                case strings.Contains(pub, "-//w3c//dtd html 4.01 strict"):
                    return "HTML 4.01 Strict"
                case strings.Contains(pub, "-//w3c//dtd html 4.01 frameset"):
                    return "HTML 4.01 Frameset"
                case strings.Contains(pub, "-//w3c//dtd xhtml 1.0 transitional"):
                    return "XHTML 1.0 Transitional"
                case strings.Contains(pub, "-//w3c//dtd xhtml 1.0 strict"):
                    return "XHTML 1.0 Strict"
                case strings.Contains(pub, "-//w3c//dtd xhtml 1.0 frameset"):
                    return "XHTML 1.0 Frameset"
                case strings.Contains(pub, "-//w3c//dtd xhtml 1.1"):
                    return "XHTML 1.1"
                }
            }

            return "Unknown Doctype: " + doctype
        }
    }

    return "Unknown (no doctype)"
}
