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
    prefix = strings.TrimSpace(prefix)
    // fmt.Println("nodes len: ", len(nodes), "prefix:|", prefix, "|")
    rootNode := &TreeNode{}
    rootChildren := []*TreeNode{}

    child := &TreeNode{}
    for _, node := range nodes {
        if node.Name == prefix {
            if strings.Index(prefix, "-") < 0 {
                rootNode.Name = prefix
            } else {
                nameParts := strings.Split(prefix, "-")
                rootNode.Name = nameParts[len(nameParts) - 1]
            }
            rootNode.Id = node.Id
            continue
        }

        if prefix != "" {
            tail := strings.TrimPrefix(node.Name, prefix + "-")
            if strings.HasPrefix(node.Name, prefix + "-") {
                if strings.Index(tail, "-") < 0 {
                    child =  nodesToTree(nodes, node.Name)
                } else {
                    newPrefix := prefix + "-" + strings.Split(tail, "-")[0]
                    child = nodesToTree(nodes, newPrefix)
                }
                rootChildren = append(rootChildren, child)
            }
        } else {
            if strings.Index(node.Name, "-") < 0 {
                child = nodesToTree(nodes, node.Name)
            } else {
                child = nodesToTree(nodes, strings.Split(node.Name, "-")[0])
            }
            rootChildren = append(rootChildren, child)
        }

    }

    if rootNode.Name == "" {
        nameParts := strings.Split(prefix, "-")
        rootNode.Name = nameParts[len(nameParts) - 1]
    }

    rootNode.Children = rootChildren
    rootNode.removeDuplicatedChild()
    return rootNode
}

func (treeNode *TreeNode) removeDuplicatedChild() {
    newChildList := []*TreeNode{}
    nameMap := make(map[string]bool)

    for _, tn := range treeNode.Children {
        if !nameMap[tn.Name] {
            newChildList = append(newChildList, tn)
            nameMap[tn.Name] = true
        }
    }
    treeNode.Children = newChildList
}

/**
 root
 ├── a
 │   └── a1
 └── b
*/
func (treeNode *TreeNode) PrintToScreen(indent int) {
    treeLines := generateNodeLines(treeNode, "")
    prefix := strings.Repeat(" ", indent)
    var prefixedLines []string
    for _, line := range treeLines {
        prefixedLines = append(prefixedLines, prefix + line)
    }
    outputLines := strings.Join(prefixedLines, "\n")
    fmt.Println(outputLines)
}

func generateNodeLines(tn *TreeNode, prefix string) []string {
    label := tn.Name
    if tn.Name == "" {
        label = "ROOT"
    }
    if tn.Id != "" {
        label = tn.Name + "(" + tn.Id + ")"
    }

    nodeLines := []string{ prefix + label }
    var childLines []string
    var childLine string
    newPrefix := prefix
    if strings.HasSuffix(newPrefix, "└── ") {
        newPrefix = strings.TrimSuffix(prefix, "└── ") + "    "
    } else if strings.HasSuffix(newPrefix, "├── ") {
        newPrefix = strings.TrimSuffix(prefix, "├── ") + "│   "
    }

    for i, childNode := range tn.Children {
        var childPrefix string
        if i == len(tn.Children) - 1 {
            childLine = newPrefix + "└── " + childNode.Name
            childPrefix = newPrefix + "    "
        } else {
            childLine = newPrefix + "├── " + childNode.Name
            childPrefix = newPrefix  + "│   "
        }
        if strings.TrimSpace(childNode.Id) != "" {
            childLine += "(" + childNode.Id + ")"
        }
        childLines = append(childLines, childLine)

        for j, grandChild := range childNode.Children {
            var grandChildPrefix string
            if j == len(childNode.Children) - 1 {
                grandChildPrefix = childPrefix + "└── "
            } else {
                grandChildPrefix = childPrefix + "├── "
            }
            grandChildLines := generateNodeLines(grandChild, grandChildPrefix)
            childLines = append(childLines, grandChildLines...)
        }

    }
    nodeLines = append(nodeLines, childLines...)
    return nodeLines
}


func (treeNode *TreeNode) flattenForPrint(depth int, maxDepth int, isLast bool) []NodeForPrint {
    if depth >= maxDepth {
        return []NodeForPrint{}
    }

    label := treeNode.Name
    if treeNode.Id != "" {
        label = label + "(" + treeNode.Id + ")"
    }

    np := NodeForPrint{label: label, depth: depth, isLast: isLast}
    npList := []NodeForPrint{np}
    for i, tn := range treeNode.Children {
        subList := []NodeForPrint{}
        if i == len(treeNode.Children) - 1 {
            subList = tn.flattenForPrint(depth + 1, maxDepth, true)
        } else {
            subList = tn.flattenForPrint(depth + 1, maxDepth, false)
        }

        // fmt.Println("subList:", subList)
        for _, subNp := range subList {
            // fmt.Println("subNp:", subNp)
            npList = append(npList, subNp)
        }
    }

    return npList
}
