package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
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

type ProjectNode struct {
	Node
	Kind string
}

type DocumentNode struct {
	Node
	Uri        string
	LanguageId string
}

type Pos struct {
	Line      int
	Character int
}

type RangeNode struct {
	Node
	Start Pos
	End   Pos
}

type ResultSetNode Node

type DefinitionNode Node

type ReferenceNode Node

type HoverNode struct {
	Node
	Result struct {
		Contents []struct {
			Language string
			Value    string
		}
	}
}

type Edge struct {
	ID    int
	Label string
	OutV  int
	InV   int
	InVs  []int
}

type ItemEdge struct {
	Edge
	Document int
	Property string
}

type Graph struct {
	Nodes map[int]interface{}
	Edges map[int]interface{}
}

func main() {
	path := "/z/dump.lsif"
	ToDot(path)

	//g, e := LoadLsif(path)
	//if e != nil {
	//log.F(e)
	//}

	//genCall(g)
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
			} else if strings.Contains(line, NODE_LABEL_RESULTSET) {
				var node ResultSetNode
				e := json.Unmarshal([]byte(line), &node)
				if e != nil {
					return g, e
				}
				g.Nodes[node.ID] = node
			} else if strings.Contains(line, NODE_LABEL_DEFINITION_RESULT) {
				var node DefinitionNode
				e := json.Unmarshal([]byte(line), &node)
				if e != nil {
					return g, e
				}
				g.Nodes[node.ID] = node
			} else if strings.Contains(line, NODE_LABEL_REFERENCE_RESULT) {
				var node ReferenceNode
				e := json.Unmarshal([]byte(line), &node)
				if e != nil {
					return g, e
				}
				g.Nodes[node.ID] = node
			} else if strings.Contains(line, NODE_LABEL_HOVER_RESULT) {
				var node HoverNode
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
		if nd, ok := node.(MetaDataNode); ok {
			log.I(fmt.Sprintf(`%d [label = "%d[%s]"]`, nd.ID, nd.ID, nd.Label))
		} else if nd, ok := node.(ProjectNode); ok {
			log.I(fmt.Sprintf(`%d [label = "%d[%s]"]`, nd.ID, nd.ID, nd.Label))
		} else if nd, ok := node.(DocumentNode); ok {
			log.I(fmt.Sprintf(`%d [label = "%d[%s]%s"]`, nd.ID, nd.ID, nd.Label, nd.Uri))
		} else if nd, ok := node.(RangeNode); ok {
			log.I(fmt.Sprintf(`%d [label = "%d[%s]line:%d:%d"]`,
				nd.ID, nd.ID, nd.Label, nd.Start.Line, nd.Start.Character))
		} else if nd, ok := node.(ResultSetNode); ok {
			log.I(fmt.Sprintf(`%d [label = "%d[%s]"]`, nd.ID, nd.ID, nd.Label))
		} else if nd, ok := node.(DefinitionNode); ok {
			log.I(fmt.Sprintf(`%d [label = "%d[%s]"]`, nd.ID, nd.ID, nd.Label))
		} else if nd, ok := node.(ReferenceNode); ok {
			log.I(fmt.Sprintf(`%d [label = "%d[%s]"]`, nd.ID, nd.ID, nd.Label))
		} else if nd, ok := node.(HoverNode); ok {
			log.I(fmt.Sprintf(`%d [label = "%d[%s]%s"]`,
				nd.ID, nd.ID, nd.Label, nd.Result.Contents[0].Value))
		}
	}

	for _, edge := range g.Edges {
		if eg, ok := edge.(ItemEdge); ok { // ItemEdge
			log.I(fmt.Sprintf(`%d -> %d [label = "%d|doc:%d|%s"]`,
				eg.OutV, eg.InVs[0], eg.ID, eg.Document, eg.Property))
		} else if eg, ok := edge.(Edge); ok { // Edge
			label := strings.Replace(eg.Label, "textDocument/", "", -1)
			if eg.InV != 0 {
				log.I(fmt.Sprintf(`%d -> %d [label = "%s"]`, eg.OutV, eg.InV, label))
			} else {
				log.I(fmt.Sprintf(`%d -> %d [label = "%s"]`, eg.OutV, eg.InVs[0], label))
			}
		}
	}

	log.I("}")

	return nil
}

func genCall(g Graph) {
	nodes := g.Nodes
	//edges := g.Edges

	for _, node := range nodes {
		if nd, ok := node.(Node); ok &&
			strings.Contains(NODE_LABEL_RANGE, nd.Label) {
			log.I(nd.Label)

		}
	}

}
