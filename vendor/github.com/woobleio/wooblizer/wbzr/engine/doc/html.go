package doc

import (
	"regexp"

	h "golang.org/x/net/html"
	"strings"
)

type html struct {
	doc     *h.Node
	curNode *h.Node
}

var exclNodes []interface{}

func NewHTML(doc string) (*html, error) {
	r := strings.NewReader(doc)
	node, err := h.Parse(r)
	if err != nil {
		return nil, err
	}

	addExcludedNodes("body", "html", "head", h.DoctypeNode, h.ErrorNode, h.DocumentNode, h.CommentNode)

	return &html{
		node,
		node,
	}, nil
}

// pIndex is parent node index in the tree (excluded and invalid nodes are not indexed)
func (html *html) ReadAndExecute(fn func(*h.Node, int) int, pIndex int) {
	n := html.curNode

	// Fixes html string format, avoid " " text nodes, for insecable space use &nbsp;
	isInvalid, _ := regexp.MatchString("^\\s+$", n.Data)
	if isInvalid {
		return
	}

	if !isExcludedNode(n) {
		pIndex = fn(n, pIndex)
	}
	for c := html.curNode.FirstChild; c != nil; c = c.NextSibling {
		html.curNode = c
		html.ReadAndExecute(fn, pIndex)
	}
}

func addExcludedNodes(nodes ...interface{}) {
	for _, n := range nodes {
		exclNodes = append(exclNodes, n)
	}
}

func isExcludedNode(node *h.Node) bool {
	for _, n := range exclNodes {
		switch n.(type) {
		case string:
			if node.Data == n {
				return true
			}
		case h.NodeType:
			if node.Type == n {
				return true
			}
		}
	}
	return false
}
