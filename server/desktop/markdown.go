package main

import (
	"fmt"

	"github.com/gomarkdown/markdown/ast"
)

func walkFile(doc *ast.Node) {
	ast.WalkFunc(*doc, func(node ast.Node, entering bool) ast.WalkStatus {
		switch t := node.(type) {
		case *ast.Heading:
			if len(t.GetChildren()) > 0 {
				for _, child := range t.GetChildren() {
					if child.AsLeaf() != nil {
						leaf := child.AsLeaf()
						if len(leaf.Literal) > 0 {
							fmt.Printf("%s\n", string(leaf.Literal))
						}
					}
				}
			}
			// fmt.Printf("%s\n", string(t.Container.Content))
		}
		// if img, ok := node.(*ast.Image); ok && entering {
		// 	attr := img.Attribute
		// 	if attr == nil {
		// 		attr = &ast.Attribute{}
		// 	}
		// 	// TODO: might be duplicate
		// 	attr.Classes = append(attr.Classes, []byte("blog-img"))
		// 	img.Attribute = attr
		// }
		//
		// if link, ok := node.(*ast.Link); ok && entering {
		// 	isExternalURI := func(uri string) bool {
		// 		return (strings.HasPrefix(uri, "https://") || strings.HasPrefix(uri, "http://")) && !strings.Contains(uri, "blog.kowalczyk.info")
		// 	}
		// 	if isExternalURI(string(link.Destination)) {
		// 		link.AdditionalAttributes = append(link.AdditionalAttributes, `target="_blank"`)
		// 	}
		// }

		return ast.GoToNext
	})
}
