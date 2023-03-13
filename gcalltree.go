package main

import (
	"bufio"
	"encoding/json"
	"os"
	"strings"

	"github.com/cjun714/glog/log"
)

type Vertex struct {
	ID    int
	Label string
}

type Edge struct {
	ID    int
	OutV  int
	InV   int
	InVs  []int
	Label string
}

func init() {
	log.SetNullFormat()
}

func main() {
	file, e := os.Open("/z/dump.lsif")
	if e != nil {
		log.F(e)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "edge") {
			//log.I(line)
			var edge Edge
			e := json.Unmarshal([]byte(line), &edge)
			if e != nil {
				log.F(e)
			}

			if edge.InV != 0 {
				log.I(edge.ID, edge.Label, edge.OutV, "->", edge.InV)
			} else {
				log.I(edge.ID, edge.Label, edge.OutV, "->", edge.InVs[0])
			}
		} else if strings.Contains(line, "vertex") {
			var vertex Vertex
			e := json.Unmarshal([]byte(line), &vertex)
			if e != nil {
				log.F(e)
			}

			log.I(vertex.Label)
		}

		if e := scanner.Err(); e != nil {
			log.F(e)
		}
	}
}
