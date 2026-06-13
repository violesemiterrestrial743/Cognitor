package store

import (
	"context"

	"github.com/kernelstub/cognitor/internal/util"
	"github.com/kernelstub/cognitor/pkg/model"
)

func (s *Store) SaveGraph(ctx context.Context, graph model.Graph) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, `delete from graph_edges`); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `delete from graph_nodes`); err != nil {
		return err
	}
	for _, node := range graph.Nodes {
		attrs, err := encode(node.Attrs)
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `insert into graph_nodes(id,type,label,attrs_json) values(?,?,?,?)`, node.ID, node.Type, node.Label, attrs); err != nil {
			return err
		}
	}
	for _, edge := range graph.Edges {
		attrs, err := encode(edge.Attrs)
		if err != nil {
			return err
		}
		id := util.StableID(edge.From, edge.To, edge.Type)
		if _, err := tx.ExecContext(ctx, `insert into graph_edges(id,from_id,to_id,type,attrs_json) values(?,?,?,?,?)`, id, edge.From, edge.To, edge.Type, attrs); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *Store) LoadGraph(ctx context.Context) (model.Graph, error) {
	var graph model.Graph
	nodeRows, err := s.db.QueryContext(ctx, `select id,type,label,attrs_json from graph_nodes order by type,label,id`)
	if err != nil {
		return graph, err
	}
	defer nodeRows.Close()
	for nodeRows.Next() {
		var node model.Node
		var attrs string
		if err := nodeRows.Scan(&node.ID, &node.Type, &node.Label, &attrs); err != nil {
			return graph, err
		}
		node.Attrs, err = decode[map[string]string](attrs)
		if err != nil {
			return graph, err
		}
		graph.Nodes = append(graph.Nodes, node)
	}
	edgeRows, err := s.db.QueryContext(ctx, `select from_id,to_id,type,attrs_json from graph_edges order by type,from_id,to_id`)
	if err != nil {
		return graph, err
	}
	defer edgeRows.Close()
	for edgeRows.Next() {
		var edge model.Edge
		var attrs string
		if err := edgeRows.Scan(&edge.From, &edge.To, &edge.Type, &attrs); err != nil {
			return graph, err
		}
		edge.Attrs, err = decode[map[string]string](attrs)
		if err != nil {
			return graph, err
		}
		graph.Edges = append(graph.Edges, edge)
	}
	return graph, edgeRows.Err()
}
