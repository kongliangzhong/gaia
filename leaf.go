package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

type Node struct {
    Id, Name, Category, Tags, Desc, Content string
    IsExecutable bool
    MainFileName string
}

func (node Node) PrintToScreen() {
    fmt.Printf("     ID: %s\n   NAME: %s\n   CATE: %s\n   TAGS: %s\n",
        node.Id, node.Name, node.Category, node.Tags)
    fmt.Printf("   DESC: %s\n", node.Desc)
    codeLines := strings.Split(node.Content, "\n")
    for i, line := range codeLines {
        if i == 0 {
            fmt.Println("CONTENT:")
            fmt.Println("    " + line)
        } else {
            fmt.Println("    " + line)
        }
    }
}

func (node Node) PrintToFile(fpath string) error {
    f, err := os.OpenFile(fpath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0660)
    if err != nil {
        return err
    }
    defer f.Close()

    f.WriteString("Id:       " + node.Id + "\n")
    f.WriteString("Category: " + node.Category + "\n")
    f.WriteString("Tags:     " + node.Tags + "\n")
    f.WriteString("Desc:     " + node.Desc + "\n")

    codeLines := strings.Split(node.Content, "\n")
    for i, line := range codeLines {
        if i == 0 {
            f.WriteString("Content:  " + line + "\n")
        } else {
            f.WriteString(CodePrefixSpace + line + "\n")
        }
    }
    return nil
}

func (node *Node) ReadFromFile(fpath string) error {
    f, err := os.Open(fpath)
    if err != nil {
        return err
    }
    defer f.Close()

    isCodeLine := false
    isDescLine := false
    scanner := bufio.NewScanner(f)
    //var codePrefixSpace string
    for scanner.Scan() {
        line := scanner.Text()
        if strings.HasPrefix(line, "Id:") {
            node.Id = strings.TrimSpace(line[len("Id:"):])
        } else if strings.HasPrefix(line, "Category:") {
            node.Category = strings.TrimSpace(line[len("Category:"):])
        } else if strings.HasPrefix(line, "Tags:") {
            node.Tags = strings.TrimSpace(line[len("Tags:"):])
        } else if strings.HasPrefix(line, "Desc:") {
            node.Desc = strings.TrimSpace(line[len("Desc:"):])
            isDescLine = true
            isCodeLine = false
        } else if strings.HasPrefix(line, "Content:") {
            node.Content = strings.TrimSpace(line[len("Content:"):])
            isCodeLine = true
            isDescLine = false
        } else {
            if isDescLine {
                var descLine string
                if strings.HasPrefix(line, CodePrefixSpace) {
                    descLine = line[len(CodePrefixSpace):]
                } else {
                    descLine = strings.TrimSpace(line)
                }
                node.Desc = node.Desc + "\n" + descLine
                //fmt.Println("node.Desc:", node.Desc)
            }

            if isCodeLine {
                var codeLine string
                if strings.HasPrefix(line, CodePrefixSpace) {
                    codeLine = line[len(CodePrefixSpace):]
                } else {
                    codeLine = strings.TrimSpace(line)
                }
                node.Content = node.Content + "\n" + codeLine
            }
        }
    }
    return nil
}
