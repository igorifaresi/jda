package interpath

import (
	"fmt"
)

type Node struct {
	Value interface{}
	Err   error
	//TODO: Add client-level verbose
}

func isMapOfInterfaces(value interface{}) bool {
	switch value.(type) {
	case map[string]interface{}:
		return true
	}
	return false
}

func isMapOfInterfacesSlice(value interface{}) bool {
	switch value.(type) {
	case []map[string]interface{}:
		return true
	}
	return false
}

func Init(value interface{}) *Node {
	return &Node{
		Value: value,
		Err: nil,
	}
}

func (node *Node) MapSlice(key string) *Node {
	if node.Err != nil {
		return node
	}
	if !isMapOfInterfaces(node.Value) {
		node.Err = fmt.Errorf(`Cannot get map slice "`+key+
			`", root is not a map[string]interface{}`)
		return node
	}
	v, ok := node.Value.(map[string]interface{})[key]
	if !ok {
		node.Err = fmt.Errorf(`Cannot get map slice "`+key+
			`", key dont exists in root map[string]interface{}`)
		return node
	}
	switch v.(type) {
	case []map[string]interface{}:
		node.Value = v
		return node
	}
	node.Err = fmt.Errorf(`Cannot get map slice "`+key+
		`", value is not []map[string]interface{}`)
	return node
}

func (node *Node) MapAt(index uint) *Node {
	if node.Err != nil {
		return node
	}
	if !isMapOfInterfacesSlice(node.Value) {
		node.Err = fmt.Errorf(`Cannot get map[string]interface{} `+
			`at %d, root is not a []map[string]interface{}`, index)
		return node
	}
	if int(index) >= len(node.Value.([]map[string]interface{})) {
		node.Err = fmt.Errorf(`Cannot get map[string]interface{} `+
			`at %d, invalid index`, index)
		return node
	}
	node.Value = node.Value.([]map[string]interface{})[index]
	return node
}

func (node *Node) Int(key string) *Node {
	if node.Err != nil {
		return Node
	}
	if !isMapOfInterfaces(node.Value) { //Can be removed, can have a flag or a diferent compilation file
		node.Err = fmt.Errorf(`Cannot get integer "`+key+
			`", root is not a map[string]interface{}`)
		return node
	}
	v, ok := node.Value.(map[string]interface{})[key]
	if !ok {
		node.Err = fmt.Errorf(`Cannot get integer "`+key+
			`", key dont exists in root map[string]interface{}`)
		return node
	}
	var integer int
	switch v.(type) {
	case int:
		integer = v.(int)
	case int32:
		integer = int(v.(int32))
	case int64:
		integer = int(v.(int64))
	case uint:
		integer = int(v.(uint))
	case uint32:
		integer = int(v.(uint32))
	case uint64:
		integer = int(v.(uint64))
	case float32:
		integer = int(v.(float32))
	case float64:
		integer = int(v.(float64))
	case byte:
		integer = int(v.(byte))	
	default:
		node.Err = fmt.Errorf(`Cannot get integer "`+key+
			`", value is not int, int32, int64, uint, uint32, uint64, float32, float64 or byte`)
		return node
	}
	node.Value = integer
	return node
}

func (node *Node) String(key string) *Node {
	if node.Err != nil {
		return Node
	}
	if !isMapOfInterfaces(node.Value) { //Can be removed, can have a flag or a diferent compilation file
		node.Err = fmt.Errorf(`Cannot get string "`+key+
			`", root is not a map[string]interface{}`)
		return node
	}
	v, ok := node.Value.(map[string]interface{})[key]
	if !ok {
		node.Err = fmt.Errorf(`Cannot get string "`+key+
			`", key dont exists in root map[string]interface{}`)
		return node
	}
	switch v.(type) {
	case string:
		node.Value = v.(string)
		return node
	}
	node.Err = fmt.Errorf(`Cannot get string "`+key+
		`", value is not string`)
	return node
}

func (node *Node) End() (interface{}, error) {
	if node.Err != nil {
		return interface{}(nil), node.Err
	}
	return node.Value, nil
}
