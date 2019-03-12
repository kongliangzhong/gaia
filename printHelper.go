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
func (treeNode *TreeNode) PrintToScreen(maxDepth int) {
    if maxDepth < 2 {
        maxDepth = 2
    }

    npList := treeNode.flattenForPrint(0, maxDepth, false)
    lines := ""
    isInLast := false
    for _, np := range npList {
        // fmt.Println("np:", np)

        line := np.label
        if np.depth > 0 {
            if np.isLast {
                line = "└── " + line
                isInLast = true
            } else {
                line = "├── " + line
            }

            if !isInLast {
                if np.depth >= 2 {
                    prefix := "│   "
                    prefix += strings.Repeat(prefix, np.depth - 1)
                    line = prefix + line
                }
            } else {
                line = strings.Repeat(" ", 4 * (np.depth - 1)) + line
            }
        }
        lines += "  " + line + "\n"
    }

    fmt.Print(lines)
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
