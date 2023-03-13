package main

import (
	"bufio"
	"encoding/json"
	"os"
	"strconv"
	"strings"

	"github.com/cjun714/glog/log"
)

const TYPE_NODE = `"type":"vertex"`
const TYPE_EDGE = `"type":"edge"`

const NODE_LABEL_METADATA = `"label":"metaData"`
const NODE_LABEL_PROJECT = `"label":"project"`
const NODE_LABEL_DOCUMENT = `"label":"document"`
const NODE_LABEL_RANGE = `"label":"range"`
const NODE_LABEL_RESULTSET = `"label":"resultSet"`
const NODE_LABEL_REFERENCE_RESULT = `"label":"referenceResult"`
const NODE_LABEL_DEFINITION_RESULT = `"label":"definitionResult"`
const NODE_LABEL_HOVER_RESULT = `"label":"hoverResult"`

const EDGE_LABEL_ITEM = `"label":"item"`
const EDGE_LABEL_NEXT = `"label":"next"`
const EDGE_LABEL_CONTAINS = `"label":"contains"`
const EDGE_LABEL_HOVER = `"label":"textDocument/hover"`
const EDGE_LABEL_REFERENCES = `"label":"textDocument/references"`
const EDGE_LABEL_DEFINITION = `"label":"textDocument/definition"`

type Vertex struct {
	ID    int
	Label string
}

func init() {
	log.SetNullFormat()
}

type Node struct {
	ID    int
	Label string
}

func (n Node) String() string {
	id := strconv.Itoa(n.ID)
	return id + " [label = \"" + id + " " + n.Label + "\"]"
}

type MetaDataNode struct {
	Node
	Version          string
	ProjectRoot      string
	PositionEncoding string
	ToolInfo         struct {
		Name    string
		Version string
	}
}

func (n MetaDataNode) String() string {
	id := strconv.Itoa(n.ID)
	return id + "[label = \"" + id + " MetaData\"]"
}

type ProjectNode struct {
	Node
	Kind string
}

func (n ProjectNode) String() string {
	id := strconv.Itoa(n.ID)
	return id + "[label = \"" + id + " Project\"]"
}

type DocumentNode struct {
	Node
	Uri        string
	LanguageId string
}

func (n DocumentNode) String() string {
	id := strconv.Itoa(n.ID)
	return id + " [label = \"" + id + "[" + n.Label + "]" + n.Uri + "\"]"
}

type Pos struct {
	Line      int
	Character int
}

func (p Pos) String() string {
	return "pos: " + strconv.Itoa(p.Line) + ":" + strconv.Itoa(p.Character)
}

type RangeNode struct {
	Node
	Start Pos
	End   Pos
}

func (n RangeNode) String() string {
	id := strconv.Itoa(n.ID)
	return id + " [label = \"" + id + "[" + n.Label + "]" + n.Start.String() + "\"]"
}

type HoverResultNode struct {
	Node
	Result struct {
		Contents []struct {
			Language string
			Value    string
		}
	}
}

func (n HoverResultNode) String() string {
	id := strconv.Itoa(n.ID)
	return id + " [label = \"" + id + "[" + n.Label + "]" + n.Result.Contents[0].Value + " \"]"
}

type Edge struct {
	ID    int
	Label string
	OutV  int
	InV   int
	InVs  []int
}

func (e Edge) String() string {
	if e.InV != 0 {
		return strconv.Itoa(e.OutV) + " -> " + strconv.Itoa(e.InV) + " [label =\"" + e.Label + "\"]"
	} else {
		return strconv.Itoa(e.OutV) + " -> " + strconv.Itoa(e.InVs[0]) + "  [label =\"" + e.Label + "\"]"
	}
}

type ItemEdge struct {
	Edge
	Document int
	Property string
}

func (e ItemEdge) String() string {
	return strconv.Itoa(e.OutV) + " -> " + strconv.Itoa(e.InVs[0]) + " [label =\"document: " + strconv.Itoa(e.Document) + "\"]"
}

type Graph struct {
	Nodes map[int]interface{}
	Edges map[int]interface{}
}

func main() {
	g, e := LoadLsif("/z/dump.lsif")
	if e != nil {
		log.F(e)
	}

	for _, edge := range g.Nodes {
		log.I(edge)
	}
}

func LoadLsif(path string) (Graph, error) {
	var g Graph
	g.Nodes = make(map[int]interface{}, 100) // TODO estimate size by .lsif size
	g.Edges = make(map[int]interface{}, 100) // TODO estimate size by .lsif size

	file, e := os.Open("/z/dump.lsif")
	if e != nil {
		return g, e
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if e := scanner.Err(); e != nil {
			return g, e
		}
		line := scanner.Text()

		// if edge
		if strings.Contains(line, TYPE_EDGE) {
			if strings.Contains(line, EDGE_LABEL_ITEM) {
				var edge ItemEdge
				e := json.Unmarshal([]byte(line), &edge)
				if e != nil {
					return g, e
				}
				g.Edges[edge.ID] = edge
			} else {
				var edge Edge
				e := json.Unmarshal([]byte(line), &edge)
				if e != nil {
					return g, e
				}
				g.Edges[edge.ID] = edge

			}
		} else {
			if strings.Contains(line, NODE_LABEL_METADATA) {
				var node MetaDataNode
				e := json.Unmarshal([]byte(line), &node)
				if e != nil {
					return g, e
				}
				g.Nodes[node.ID] = node
			} else if strings.Contains(line, NODE_LABEL_PROJECT) {
				var node ProjectNode
				e := json.Unmarshal([]byte(line), &node)
				if e != nil {
					return g, e
				}
				g.Nodes[node.ID] = node
			} else if strings.Contains(line, NODE_LABEL_DOCUMENT) {
				var node DocumentNode
				e := json.Unmarshal([]byte(line), &node)
				if e != nil {
					return g, e
				}
				g.Nodes[node.ID] = node
			} else if strings.Contains(line, NODE_LABEL_RANGE) {
				var node RangeNode
				e := json.Unmarshal([]byte(line), &node)
				if e != nil {
					return g, e
				}
				g.Nodes[node.ID] = node
			} else if strings.Contains(line, NODE_LABEL_HOVER_RESULT) {
				var node HoverResultNode
				e := json.Unmarshal([]byte(line), &node)
				if e != nil {
					return g, e
				}
				g.Nodes[node.ID] = node
			} else {
				var node Node
				e := json.Unmarshal([]byte(line), &node)
				if e != nil {
					return g, e
				}
				g.Nodes[node.ID] = node
			}
		}
	}

	return g, nil
}

func test2() {
	var node HoverResultNode

	str := `{"id":4,"type":"vertex","label":"hoverResult","result":{"contents":[{"language":"c","value":"char quote"}]}}`

	e := json.Unmarshal([]byte(str), &node)
	if e != nil {
		log.F(e)
	}
	log.I(node.Result.Contents[0].Value)
}

func main2() {
	log.I("digraph g{")
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
				log.I("  ", edge.OutV, "->", edge.InV, "[label = \"", edge.Label, "\"]")
			} else {
				log.I("  ", edge.OutV, "->", edge.InVs[0], "[label = \"", edge.Label, "\"]")
			}
		} else if strings.Contains(line, "vertex") {
			var vertex Vertex
			e := json.Unmarshal([]byte(line), &vertex)
			if e != nil {
				log.F(e)
			}

			log.I("  ", vertex.ID, "[label = \"", vertex.Label, "\"]")
			//log.I(vertex.Label)
		}

		if e := scanner.Err(); e != nil {
			log.F(e)
		}
	}
	log.I("}")
}

func test() {
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
