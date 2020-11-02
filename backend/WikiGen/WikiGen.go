package WikiGen

import (
	"fmt"
	"os"
	"strings"
	"time"

	"backend/IPSheets/Parsing"
)

func GetSubPage(parent string) string {
	mpl := Parsing.GetMpl()
	var sb strings.Builder

	sb.WriteString(fmt.Sprintln("=", parent, "-", mpl[parent].Desc, "=\n"))

	sb.WriteString("==Bill of Materials==\n")
	sb.WriteString(fmt.Sprintln("Generated on:", time.Now()))

	sb.WriteString(GenBOM(parent))

	sb.WriteString("\n----\nEnd-of-autogenerated")
	return sb.String()
}

func GenBOM(parent string) string {
	mpl := Parsing.GetMpl()
	Children := Parsing.GetSubs()[parent]

	data := make([][]interface{}, len(Children)+1)
	data[0] = []interface{}{"SKU", "Description", "Qty", "Location/Notes"}
	for i, child := range Children {
		line := []interface{}{fmt.Sprint("[[", child.Child, "]]"), mpl[child.Child].Desc, child.Qty, child.Location}
		data[i+1] = line
	}

	table := genTable(data)
	return table
}

func genTable(data [][]interface{}) string {
	var sb strings.Builder
	sb.WriteString("{| class=\"wikitable\"\n")

	for i, line := range data {
		if i == 0 {
			sb.WriteString("|+\n")
			for _, val := range line {
				sb.WriteString("!")
				sb.WriteString(fmt.Sprint(val))
				sb.WriteString("\n")
			}
		} else {
			sb.WriteString("|-\n")
			for _, val := range line {
				sb.WriteString("|")
				sb.WriteString(fmt.Sprint(val))
				sb.WriteString("\n")
			}
		}
	}
	sb.WriteString("|}\n")

	return sb.String()
}

func CreateFile(content string, title string) {
	f, err := os.Create(fmt.Sprint("generatedTextFiles\\", title))
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = f.WriteString(content)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	err = f.Close()
}
