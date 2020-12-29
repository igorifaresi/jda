package interpath

import (
	"github.com/igorifaresi/jda"
)

type Node struct {
	Value interface{}
	Err   jda.LoggerErrorQueue
}

func isMapOfInterfaces(value interface{}) bool {
	switch value.(type) {
	case map[string]interface{}:
		return true
	}
	return false
}

func Init(value interface{}) Node {
	return Node{
		Value: value,
		Err: nil,
	}
}

func (node *Node) MapSlice(key string) Node {
	if node.Err != nil {
		return Node
	}
	if !isMapOfInterfaces(node.Value) {
		
	}
}

func (node *Node) At(index int) Node {
	if node.Err != nil {
		return Node
	}
}

func (node *Node) Int(key string) Node {
	if node.Err != nil {
		return Node
	}
}

func (node *Node) String(key string) Node {
	if node.Err != nil {
		return Node
	}
}

func (node *Node) End() (interface{}, error) {
	if node.Err != nil {
		return interface{}(nil), node.Err
	}
	return node.Value, nil
}
