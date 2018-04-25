package main

import (
	"github.com/TheGrum/booktools/booktools/cmd"
)

func main() {
	cmd.Execute()

	//fmt.Println("Identifying Characters")
	//PrintNameFrequency(CharacterFrequencies(root, 3, 1))

	/*
		fmt.Printf("%v", root.Structure())
		fmt.Printf("%v", root.Children[1].Children[0])

		//section := root.Children[0].Children[0]
		fmt.Println("Prose:")
		m := make(map[string]int)
		for _, chapter := range root.Children {
			fmt.Println("Chapter:")
			for _, section := range chapter.Children {
				sec_text := section.String()
				NameFrequencyMapMerge(m, GetNameFrequencyViaProse(sec_text))
				fmt.Print(".")
			}
		}
		fmt.Println("NameFrequency")
		PrintNameFrequency(m)
	*/
}
