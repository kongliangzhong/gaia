package main

import (
    "fmt"
    "strings"
)

type TreeNode struct {
    Name, Id string
    Children []*TreeNode
}

type NodeForPrint struct {
    label string
    depth int
    isLast bool
}

func mapListToTree(mapList map[string][]string, root string) *TreeNode {
    rootNode := &TreeNode{ Name: root }
    rootChildren := []*TreeNode{}

    for cate, nameArray := range mapList {
        cateNode := &TreeNode{ Name: cate }
        cateChildren := []*TreeNode{}
        for _, name := range nameArray {
            itemNode := &TreeNode{ Name: name }
            cateChildren = append(cateChildren, itemNode)
        }
        cateNode.Children = cateChildren
        rootChildren = append(rootChildren, cateNode)
    }
    rootNode.Children = rootChildren
    return rootNode
}

func nodesToTree(nodes []Node, prefix string) *TreeNode {
    rootNode := &TreeNode{Name: prefix}
    rootChildren := []*TreeNode{}

    for _, node := range nodes {
        if node.Name == prefix {
            rootNode.Id = node.Id
            continue
        }

        if strings.HasPrefix(node.Name, prefix + "-") {
            tail := strings.TrimPrefix(node.Name, prefix + "-")
            if strings.Index(tail, "-") < 0 {
                child := nodesToTree(nodes, node.Name)
                rootChildren = append(rootChildren, child)
            }
        }
    }
    rootNode.Children = rootChildren

    return rootNode
}

/**
 root
 ├── a
 │   └── a1
 └── b
*/
func (treeNode *TreeNode) PrintToScreen(maxDepth int) {
    if maxDepth < 2 {
        maxDepth = 2
    }

    npList := []*NodeForPrint{}
    treeNode.flattenForPrint(npList, 0, maxDepth, false)
    lines := ""
    isInLast := false
    for _, np := range npList {
        line := np.label
        if np.depth > 0 {
            if np.isLast {
                line = "└── " + line

                if np.depth == 1 {
                    isInLast = true
                }
            } else {
                line = "├── " + line
            }

            if !isInLast {
                if np.depth >= 2 {
                    prefix := "│   "
                    prefix += strings.Repeat(" ", 4 * (np.depth - 2))
                    line = prefix + line
                }
            }
        }
        lines += line + "\n"
    }

    fmt.Print(lines)
}

func (treeNode *TreeNode) flattenForPrint(npList []*NodeForPrint, depth int, maxDepth int, isLast bool) {
    if depth >= maxDepth {
        return
    }

    label := treeNode.Name
    if treeNode.Id != "" {
        label = label + "(" + treeNode.Id + ")"
    }

    np := &NodeForPrint{label: label, depth: depth, isLast: isLast}
    npList = append(npList, np)

    for i, tn := range treeNode.Children {
        if i == len(treeNode.Children) - 1 {
            tn.flattenForPrint(npList, depth + 1, maxDepth, true)
        } else {
            tn.flattenForPrint(npList, depth + 1, maxDepth, false)
        }
    }
}
