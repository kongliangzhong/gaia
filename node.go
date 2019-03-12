package main

import (
    "fmt"
    "strings"
    "os"
    "bufio"
)

type Node struct {
    Id, Name, Category, Tags, Desc, Content string
    Executable bool
    ExecFile string
    Attachments []string // file path array.
}

var CodePrefixSpace string = "          " // len:10

func (node Node) String() string {
    res := ""
    res += fmt.Sprintf("        Id: %s\n", node.Id)
    res += fmt.Sprintf("      Name: %s\n", node.Name)
    res += fmt.Sprintf("  Category: %s\n", node.Category)
    res += fmt.Sprintf("      Tags: %s\n", node.Tags)
    res += fmt.Sprintf("Executable: %t\n", node.Executable)
    res += fmt.Sprintf("  ExecFile: %s\n", node.ExecFile)
    res += fmt.Sprintf("      Desc: %s\n", node.Desc)
    codeLines := strings.Split(node.Content, "\n")
    for i, line := range codeLines {
        if i == 0 {
            res += fmt.Sprintln("   Content:")
            res += fmt.Sprintln("    " + line)
        } else {
            res += fmt.Sprintln("    " + line)
        }
    }
    return res
}

func (node Node) PrintToScreen() {
    fmt.Println(node.String())
}

func (node *Node) Normalize(aliasMap map[string]string) error {
    normalizeStr := func(s string, sep string) string {
        result := strings.ToLower(strings.TrimSpace(s))
        if sep == "" {
            return result
        } else {
            parts := strings.Split(result, sep)
            resultParts := []string{}
            for _, part := range parts {
                if aliasMap[part] != "" {
                    resultParts = append(resultParts, aliasMap[part])
                }
            }

            return strings.Join(resultParts, sep)
        }
    }

    node.Name = normalizeStr(node.Name, "-")
    node.Category = normalizeStr(node.Category, "")
    node.Tags = normalizeStr(node.Tags, ",")
    node.Desc = strings.TrimSpace(node.Desc)
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
    for scanner.Scan() {
        line := scanner.Text()
        line = strings.TrimSpace(line)
        if strings.HasPrefix(line, "Id:") {
            node.Id = line[len("Id:"):]
        } else if strings.HasPrefix(line, "Name:") {
            node.Name = line[len("Name:"):]
        } else if strings.HasPrefix(line, "Category:") {
            node.Category = line[len("Category:"):]
        } else if strings.HasPrefix(line, "Tags:") {
            node.Tags = line[len("Tags:"):]
        } else if strings.HasPrefix(line, "Executable:") {
            executableStr := line[len("Executable:"):]
            if strings.EqualFold(executableStr, "true") {
                node.Executable = true
            } else {
                node.Executable = false
            }
        } else if strings.HasPrefix(line, "ExecFile:") {
            node.ExecFile = line[len("ExecFile:"):]
        } else if strings.HasPrefix(line, "Desc:") {
            node.Desc = line[len("Desc:"):]
            isDescLine = true
            isCodeLine = false
        } else if strings.HasPrefix(line, "Content:") {
            node.Content = line[len("Content:"):]
            isCodeLine = true
            isDescLine = false
        } else {
            if isDescLine {
                node.Desc = node.Desc + "\n" + line
            }
            if isCodeLine {
                node.Content = node.Content + "\n" + line
            }
        }
    }
    return nil
}
