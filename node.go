package main

import (
    "fmt"
    "strings"
    "os"
    "bufio"
)

type Node struct {
    Id string
    Name string
    Category string
    Tags string
    Desc string
    Content string
    Executable bool
    ExecFile string
    Attachments []string // file path array.
    Links string // comma seperated node ids.
}

var CodePrefixSpace string = "    " // indent: 4

func (node Node) ShortString() string {
    res := ""
    res += fmt.Sprintf("        ID: %s\n", node.Id)
    res += fmt.Sprintf("      NAME: %s\n", node.Name)
    if node.Desc != "" {
        res += fmt.Sprintf("      DESC: %s\n", node.Desc)
    }
    return res
}

func (node Node) StringWithoutEmpty() string {
    res := ""
    res += fmt.Sprintf("        ID: %s\n", node.Id)
    res += fmt.Sprintf("      NAME: %s\n", node.Name)
    // res += fmt.Sprintf("  Category: %s\n", node.Category)
    if node.Tags != "" {
        res += fmt.Sprintf("      TAGS: %s\n", node.Tags)
    }
    res += fmt.Sprintf("EXECUTABLE: %t\n", node.Executable)

    if node.Executable {
        res += fmt.Sprintf("  EXECFILE: %s\n", node.ExecFile)
    }

    if node.Desc != "" {
        res += fmt.Sprintf("      DESC: %s\n", node.Desc)
    }
    codeLines := strings.Split(node.Content, "\n")
    for i, line := range codeLines {
        if i == 0 {
            res += fmt.Sprintln("   CONTENT:")
            res += fmt.Sprintln(CodePrefixSpace + line)
        } else {
            res += fmt.Sprintln(CodePrefixSpace + line)
        }
    }
    return res
}

func (node Node) String() string {
    res := ""
    res += fmt.Sprintf("        ID: %s\n", node.Id)
    res += fmt.Sprintf("      NAME: %s\n", node.Name)
    // res += fmt.Sprintf("  CATEGORY: %s\n", node.Category)
    res += fmt.Sprintf("      TAGS: %s\n", node.Tags)
    res += fmt.Sprintf("EXECUTABLE: %t\n", node.Executable)
    res += fmt.Sprintf("  EXECFILE: %s\n", node.ExecFile)
    res += fmt.Sprintf("      DESC: %s\n", node.Desc)
    codeLines := strings.Split(node.Content, "\n")
    for i, line := range codeLines {
        if i == 0 {
            res += fmt.Sprintln("   CONTENT:")
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
                } else {
                    resultParts = append(resultParts, part)
                }
            }
            return strings.Join(resultParts, sep)
        }
    }

    node.Id = normalizeStr(node.Id, "")
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
        rawLine := scanner.Text()
        line := strings.TrimSpace(rawLine)
        if strings.HasPrefix(line, "ID:") {
            node.Id = line[len("ID:"):]
        } else if strings.HasPrefix(line, "NAME:") {
            node.Name = line[len("NAME:"):]
        } else if strings.HasPrefix(line, "CATEGORY:") {
            node.Category = line[len("CATEGORY:"):]
        } else if strings.HasPrefix(line, "TAGS:") {
            node.Tags = line[len("TAGS:"):]
        } else if strings.HasPrefix(line, "EXECUTABLE:") {
            executableStr := line[len("EXECUTABLE:"):]
            if strings.EqualFold(executableStr, "true") {
                node.Executable = true
            } else {
                node.Executable = false
            }
        } else if strings.HasPrefix(line, "EXECFILE:") {
            node.ExecFile = line[len("EXECFILE:"):]
        } else if strings.HasPrefix(line, "DESC:") {
            node.Desc = line[len("DESC:"):]
            isDescLine = true
            isCodeLine = false
        } else if strings.HasPrefix(line, "CONTENT:") {
            node.Content = line[len("CONTENT:"):]
            isCodeLine = true
            isDescLine = false
        } else {
            if isDescLine {
                node.Desc = node.Desc + "\n" + line
            }
            if isCodeLine {
                codeLine := line
                if strings.HasPrefix(rawLine, CodePrefixSpace) {
                    codeLine = strings.TrimPrefix(rawLine, CodePrefixSpace)
                }
                node.Content = node.Content + "\n" + codeLine
            }
        }
    }

    node.Content = strings.TrimPrefix(node.Content, "\n")
    return node.Normalize(map[string]string{})
}
