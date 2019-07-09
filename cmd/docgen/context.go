package main

import (
	"fmt"
	"io/ioutil"
	"path"
	"sort"
	"strings"

	"github.com/nyaruka/gocommon/urns"
	"github.com/nyaruka/goflow/cmd/docgen/context"
	"github.com/nyaruka/goflow/utils"

	"github.com/pkg/errors"
)

func init() {
	registerGenerator("context map", generateContextMap)
}

func generateContextMap(baseDir string, outputDir string, items map[string][]*TaggedItem) error {
	ctx := context.NewContext()

	// the dynamic types in the context aren't described in the code so we add them manually here
	ctx.AddType(context.NewDynamicType("fields", "field-keys", context.NewProperty("{key}", "{key} for the contact", "any")))
	ctx.AddType(context.NewDynamicType("results", "result-keys", context.NewProperty("{key}", "{key} value for the run", "result")))

	// the urns type also added here as it's "dynamic" in sense that keys are known at build time
	ctx.AddType(createURNsType())

	// now add the types from tagged docstrings
	for _, item := range items["context"] {
		// examples are actually property descriptors for context items
		properties := make([]*context.Property, len(item.examples))
		for i, propDesc := range item.examples {
			prop := context.ParseProperty(propDesc)
			if prop == nil {
				return errors.Errorf("invalid format for property description \"%s\"", propDesc)
			}
			properties[i] = prop
		}

		if item.tagValue == "root" {
			ctx.SetRoot(properties)
		} else {
			ctx.AddType(context.NewStaticType(item.tagValue, properties))
		}
	}

	if err := ctx.Validate(); err != nil {
		return err
	}

	mapPath := path.Join(outputDir, "context.json")
	marshaled, _ := utils.JSONMarshalPretty(ctx)
	ioutil.WriteFile(mapPath, marshaled, 0755)

	fmt.Printf(" > %d context types written to %s\n", len(items["context"]), mapPath)

	nodes := ctx.EnumerateNodes(map[string][]string{
		"field-keys":  []string{"age", "gender"},
		"result-keys": []string{"response_1"},
	})
	nodeOutput := &strings.Builder{}
	for _, n := range nodes {
		nodeOutput.WriteString(fmt.Sprintf("%s -> %s\n", n.Path, n.Help))
	}

	listPath := path.Join(outputDir, "context.txt")
	ioutil.WriteFile(listPath, []byte(nodeOutput.String()), 0755)

	return nil
}

func createURNsType() context.Type {
	properties := make([]*context.Property, 0, len(urns.ValidSchemes))
	for k := range urns.ValidSchemes {
		name := strings.Title(k)
		properties = append(properties, context.NewProperty(k, name+" URN for the contact", "text"))
	}
	sort.SliceStable(properties, func(i, j int) bool { return properties[i].Key < properties[j].Key })

	return context.NewStaticType("urns", properties)
}
