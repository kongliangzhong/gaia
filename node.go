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
    res += fmt.Sprintf("Executable: %s\n", node.Executable)
    res += fmt.Sprintf("  ExecFile: %s\n", node.ExecFile)
    res += fmt.Sprintf("      Desc: %s\n", node.Desc)
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

// func (node Node) PrintToFile(fpath string) error {
//     f, err := os.OpenFile(fpath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0660)
//     if err != nil {
//         return err
//     }
//     defer f.Close()

//     f.WriteString("Id:       " + node.Id + "\n")
//     f.WriteString("Category: " + node.Category + "\n")
//     f.WriteString("Tags:     " + node.Tags + "\n")
//     f.WriteString("Desc:     " + node.Desc + "\n")

//     codeLines := strings.Split(node.Content, "\n")
//     for i, line := range codeLines {
//         if i == 0 {
//             f.WriteString("Content:  " + line + "\n")
//         } else {
//             f.WriteString(CodePrefixSpace + line + "\n")
//         }
//     }
//     return nil
// }

func (node *Node) Normalize() error {
    node.Name = strings.ToLower(node.Name)
    node.Category = strings.ToLower(node.Category)
    node.Tags = strings.ToLower(node.Tags)
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
