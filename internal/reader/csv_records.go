package reader

import "fmt"

type CSVFile struct {
	Name string
	Records [][]string
}

func (c CSVFile) Print() {
	fmt.Println(c.Name)
	for row, line := range c.Records {
		fmt.Println(row, line)
	}
}