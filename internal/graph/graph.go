package graph

import (
	"sort"

	"github.com/kernelstub/cognitor/internal/util"
	"github.com/kernelstub/cognitor/pkg/model"
)

func Build(snapshot model.Snapshot, findings []model.Finding) model.Graph {
	nodes := map[string]model.Node{}
	var edges []model.Edge
	for _, binary := range snapshot.Binaries {
		binaryID := "binary:" + binary.Path
		nodes[binaryID] = model.Node{ID: binaryID, Type: NodeBinary, Label: binary.Path, Attrs: map[string]string{"kind": binary.Kind}}
		for _, fn := range binary.Functions {
			fnID := "function:" + binary.Path + ":" + fn.Name
			nodes[fnID] = model.Node{ID: fnID, Type: NodeFunction, Label: fn.Name, Attrs: map[string]string{"binary": binary.Path}}
			edges = append(edges, model.Edge{From: binaryID, To: fnID, Type: EdgeContains, Attrs: map[string]string{}})
		}
		for _, imp := range binary.Imports {
			impID := "import:" + imp
			nodes[impID] = model.Node{ID: impID, Type: NodeImport, Label: imp, Attrs: map[string]string{}}
			edges = append(edges, model.Edge{From: binaryID, To: impID, Type: EdgeImports, Attrs: map[string]string{}})
		}
	}
	for _, finding := range findings {
		findingID := "finding:" + finding.ID
		nodes[findingID] = model.Node{ID: findingID, Type: NodeFinding, Label: finding.Title, Attrs: map[string]string{"severity": finding.Severity, "category": finding.Category}}
		targetID := "function:" + finding.AffectedBinary + ":" + finding.NewFunction
		if _, ok := nodes[targetID]; ok {
			edges = append(edges, model.Edge{From: findingID, To: targetID, Type: EdgeAffects, Attrs: map[string]string{}})
		}
	}
	var graph model.Graph
	for _, node := range nodes {
		graph.Nodes = append(graph.Nodes, node)
	}
	sort.Slice(graph.Nodes, func(i, j int) bool { return graph.Nodes[i].ID < graph.Nodes[j].ID })
	sort.Slice(edges, func(i, j int) bool {
		return util.StableID(edges[i].From, edges[i].To, edges[i].Type) < util.StableID(edges[j].From, edges[j].To, edges[j].Type)
	})
	graph.Edges = edges
	return graph
}
