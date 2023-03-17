## Vertex
- metaData
- project
- document
- range
- resultSet
- definitionResult
- referenceResult
- hoverResult

metaData:
{"id":0,"type":"vertex","label":"metaData","version":"0.4.3","projectRoot":"file:///z/vim-fzf","positionEncoding":"utf-16","toolInfo":{"name":"lsif-clang","version":"0.1.0"}}

project:
{"id":1,"type":"vertex","label":"project","kind":"cpp"}

document:
{"id":9,"type":"vertex","label":"document","uri":"file:///z/vim-fzf/cmd/colorize_c_cpp.h","languageId":"c"}

range:
{"id":513,"type":"vertex","label":"range","start":{"line":28,"character":7},"end":{"line":28,"character":11}}

resultSet:
{"id":2,"type":"vertex","label":"resultSet"}

referenceResult:
{"id":3,"type":"vertex","label":"referenceResult"}

definitionResult:
{"id":20,"type":"vertex","label":"definitionResult"}

hoverResult:
{"id":4,"type":"vertex","label":"hoverResult","result":{"contents":[{"language":"c","value":"char quote"}]}}

## Edge
- item
- next
- contains
- textDocument/definition
- textDocument/references
- textDocument/hover

item:
{"id":517,"type":"edge","label":"item","outV":510,"inVs":[513],"document":512}

next:
{"id":12,"type":"edge","label":"next","outV":10,"inV":2}

contains:
{"id":11,"type":"edge","label":"contains","outV":9,"inVs":[10]}

textDocument/hover
{"id":5,"type":"edge","label":"textDocument/hover","outV":2,"inV":4}

textDocument/references
{"id":6,"type":"edge","label":"textDocument/references","outV":2,"inV":3}

textDocument/definition
{"id":8,"type":"edge","label":"textDocument/definition","outV":2,"inV":7}

## Drawing Graph
- gnuplot: math formular
- Mermaid: Generate diagrams from markdown-like text
- diagon: ascii flow chart
- go-callvis: call graph of a Go
- Graphviz: generate lsif relation chart

## Lsif Spec
- Document
  * Range line:10 character:15
  * Range line:39 character:8
	* ...
- ResultSet
  * Definition
	* Reference
	* Hover

Each `range` is a pos of symbol in that doc

    doc ─┬─ range ── resultSet ─┬─ define ── range
         ├─ range               ├─ refer  ── range
         ├─ range               └─ hover  ── range

Core understanding:
- `doc` contains multi `range``
- each `range` link to `resultSet`, defined symbol and refered symbol link to
  same `resultSet`
- `resultSet` link to 3 nodes:
  * `defineResult`
	* `referResult`
	* `hoverResult`
- `defineResult`/`referResult` link with a `range`
- `range` doesn't include file path, file path exists in out 'edge'

``` ascii
+-------------------+  contains  +--------------------------------------+     Doc:1
|  1[Doc]./main.c   | ---------> |         2[Range]line 0:5             | <-----------+
+-------------------+            +--------------------------------------+             |
          |                           ^           ^          |      |                 |
          | contains             Doc:1 refer Doc:1 define   next   next               |
          v                           |           |          v      v                 |
+-------------------+  Doc:1 refer  +----------------+   +--------------+  define  +-----------------+
| 7[Range]line 7:13 | <------------ | 5[ReferResult] |   | 3[ResultSet] | -------> | 4[DefineResutl] |
+-------------------+               +----------------+   +--------------+          +-----------------+
          |                            ^      refer        |    ^    |
          |                            +-------------------+    |    | hover
          |                                                     |    v
          |                        next                         |  +-----------------------------+
          +-----------------------------------------------------+  | 6[HoverResult]void foo(int) |
                                                                   +-----------------------------+
```

## cmd References
lsif:
``` bash
> lsif-clang compile_commands.json > dump.lsif # c/c++
```

## Reference
https://unix.stackexchange.com/questions/126630/creating-diagrams-in-ascii
https://stackoverflow.com/questions/3211801/graphviz-and-ascii-output
