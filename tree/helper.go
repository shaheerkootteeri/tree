package tree

import (
	"fmt"
	"sort"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	paddingStart = 3
	Unknown      = "<unknown>"
)

func TreeFromMapValues(mapValue map[string]any) []TreeNode {
	if mapValue == nil {
		return nil
	}
	keys := make([]string, 0, len(mapValue))
	for k := range mapValue {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	nodes := []TreeNode{}
	for _, k := range keys {
		v := mapValue[k]
		n := NewNode(k)
		switch t := v.(type) {
		case map[string]any:
			n.AddChildren(TreeFromMapValues(t))
		case []any:
			for index, item := range t {
				if itemMap, ok := item.(map[string]any); ok {
					indexNode := NewNode(fmt.Sprintf("[%d]", index+1))
					indexNode.AddChildren(TreeFromMapValues(itemMap))
					n.AddChild(indexNode)
				} else {
					n.AddChild(NewNode(fmt.Sprintf("[%d] %v", index+1, item)))
				}
			}
		default:
			n.SetValue(fmt.Sprintf("%s: %v", n.Value(), t))
		}
		nodes = append(nodes, n)
	}
	return nodes
}

func TreeFromMapStringValues(mapValue map[string]string) []TreeNode {
	mapStrings := make(map[string]interface{}, len(mapValue))
	for k, v := range mapValue {
		mapStrings[k] = v
	}
	return TreeFromMapValues(mapStrings)
}

func TreeFromMatchExpressions(expressions []metav1.LabelSelectorRequirement) []TreeNode {
	nodes := []TreeNode{}
	for index, value := range expressions {
		indexNode := NewNode(fmt.Sprintf("[%d]", index+1))
		indexNode.AddChild(NewNode(fmt.Sprintf("%s: %v", "key", value.Key)))
		indexNode.AddChild(NewNode(fmt.Sprintf("%s: %v", "operator", value.Operator)))
		valueNode := NewNode("values")
		for index, item := range value.Values {
			valueNode.AddChild(NewNode(fmt.Sprintf("[%d] %v", index+1, item)))
		}
		indexNode.AddChild(valueNode)
		nodes = append(nodes, indexNode)
	}
	return nodes
}
