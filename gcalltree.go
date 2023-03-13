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
	return id + " [label = \"" + id + "[" + n.Label + "]\"]"
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
	return id + " [label = \"" + id + " MetaData\"]"
}

type ProjectNode struct {
	Node
	Kind string
}

func (n ProjectNode) String() string {
	id := strconv.Itoa(n.ID)
	return id + " [label = \"" + id + "[" + n.Label + "]\"]"
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
	return id + " [label = \"" + id + "[" + n.Label + "]" + n.Result.Contents[0].Value + "\"]"
}

type Edge struct {
	ID    int
	Label string
	OutV  int
	InV   int
	InVs  []int
}

func (e Edge) String() string {
	id := strconv.Itoa(e.ID)
	label := strings.Replace(e.Label, "textDocument", "", -1)

	if e.InV != 0 {
		return strconv.Itoa(e.OutV) + " -> " + strconv.Itoa(e.InV) +
			" [label = \"" + id + " " + label + "\"]"
	} else {
		return strconv.Itoa(e.OutV) + " -> " + strconv.Itoa(e.InVs[0]) +
			" [label = \"" + id + " " + label + "\"]"
	}
}

type ItemEdge struct {
	Edge
	Document int
	Property string
}

func (e ItemEdge) String() string {
	id := strconv.Itoa(e.ID)
	return strconv.Itoa(e.OutV) + " -> " + strconv.Itoa(e.InVs[0]) +
		" [label =\"" + id + " doc:" + strconv.Itoa(e.Document) + " " + e.Property + "\"]"
}

type Graph struct {
	Nodes map[int]interface{}
	Edges map[int]interface{}
}

func main() {
	ToDot("/z/dump.lsif")
}

func LoadLsif(path string) (Graph, error) {
	var g Graph
	g.Nodes = make(map[int]interface{}, 100) // TODO estimate size by .lsif size
	g.Edges = make(map[int]interface{}, 100) // TODO estimate size by .lsif size

	file, e := os.Open(path)
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

func ToDot(path string) error {
	g, e := LoadLsif(path)
	if e != nil {
		return e
	}

	log.I("digraph g{")
	for _, node := range g.Nodes {
		log.I(node)
	}
	for _, edge := range g.Edges {
		log.I(edge)
	}
	log.I("}")

	return nil
}
