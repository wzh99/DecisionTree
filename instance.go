package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type Instance struct {
	decision int
	attrib [6]int
}

func DataFromFile(path string) (data []Instance, err error) {
	data = make([]Instance, 0)
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		line, _, err  := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				return data, nil
			} else {
				return data, err
			}
		}
		var inst Instance
		_, _ = fmt.Sscanf(string(line), "%d %d %d %d %d %d %d", &inst.decision, &inst.attrib[0],
			&inst.attrib[1], &inst.attrib[2], &inst.attrib[3], &inst.attrib[4], &inst.attrib[5])
		data = append(data, inst)
	}
}
