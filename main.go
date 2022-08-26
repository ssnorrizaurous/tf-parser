package main

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"io/ioutil"
	"log"
	"strings"
)

type TfModule struct {
	Type         string   `hcl:"source,label"`
	PluginConfig hcl.Body `hcl:",remain"`
}

func main() {
	//var tfModule TfModule
	//hclwrite.File{}
	//err := hclsimple.DecodeFile("example.tf", nil, &tfModule)
	//if err != nil {
	//	log.Fatalf("Failed to load configuration: %s", err)
	//}
	//log.Printf("Configuration is %#v", tfModule)

	//parser := hclparse.NewParser()
	//f, diags := parser.ParseHCLFile("example.tf")
	//
	//hclwrite.File{}
	//log.Printf("Configuration is %#v", f.)
	//log.Printf("Configuration is %#v", diags)
	fn := "example.tf"
	src, err := ioutil.ReadFile(fn)
	if err != nil {
		log.Printf("Failed to read file %q: %s", fn, err)
		return
	}

	//log.Printf("%v", string(src))

	f, diags := hclwrite.ParseConfig(src, fn, hcl.Pos{Line: 1, Column: 1})
	log.Printf("%v", diags)
	log.Printf("%v", f.Body().Blocks()[0].Labels())

	parseBody(f.Body(), nil)
	//if diags.HasErrors() {
	//	for _, diag := range diags {
	//		if diag.Subject != nil {
	//			log.Printf("[%s:%d] %s: %s", diag.Subject.Filename, diag.Subject.Start.Line, diag.Summary, diag.Detail)
	//		} else {
	//			log.Printf("%s: %s", diag.Summary, diag.Detail)
	//		}
	//	}
	//	return
	//}
}

func parseTokens(tokens hclwrite.Tokens) []string {
	lst := make([]string, 0)
	for _, t := range tokens {
		lst = append(lst, string(t.Bytes))
	}

	return lst
}

func parseBody(body *hclwrite.Body, inBlocks []string) {
	attrs := body.Attributes()
	for name, attr := range attrs {
		//var cleanedExprTokens hclwrite.Tokens
		tokens := attr.Expr().BuildTokens(nil)
		//if len(inBlocks) == 1 {
		//inBlock := inBlocks[0]
		fmt.Printf("%s%s = %s\n", strings.Repeat("\t", len(inBlocks)), name, strings.Join(parseTokens(tokens), ""))
		if inBlocks[0] == "filter" {
			fmt.Print("found filter!\n")
		}
		//if inBlock == "variable" && name == "type" {
		//	cleanedExprTokens = cleanTypeExpr(tokens)
		//	body.SetAttributeRaw(name, cleanedExprTokens)
		//	continue
		//} else if (inBlock == "resource" || inBlock == "data") && name == "provider" {
		//	cleanedExprTokens = cleanProviderExpr(tokens)
		//	body.SetAttributeRaw(name, cleanedExprTokens)
		//	continue
		//}
		//} else {
		//	// TODO
		//	//fmt.Printf("%s: %vn", inBlocks[0], attr.Expr().Variables())
		//	//fmt.Printf("%s = %s\n", name, strings.Join(parseTokens(tokens), ""))
		//}
		//cleanedExprTokens = cleanValueExpr(tokens)
		//body.SetAttributeRaw(name, cleanedExprTokens)
	}

	blocks := body.Blocks()
	// TODO why is filter not being parsed
	for offset, block := range blocks {

		//fmt.Printf("%s.%v\n", block.Type(), block.Labels())

		// TODO fix
		if len(block.Labels()) > 0 {
			fmt.Printf("%s%s.%s\n", strings.Repeat("\t", offset), block.Type(), strings.Join(block.Labels(), "."))
		} else {
			fmt.Printf("%s%s\n", strings.Repeat("\t", offset), block.Type())
		}
		inBlocks := append(inBlocks, block.Type())
		//fmt.Printf("inblocks: %v, body: %v\n", inBlocks, block.Labels())

		parseBody(block.Body(), inBlocks)
	}
}
