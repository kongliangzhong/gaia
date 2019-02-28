package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

type Leaf struct {
    Id, Name, Category, Tags, Desc, Content string
    IsExecutable bool
    MainFileName string
}

func (lf Leaf) PrintToScreen() {
    fmt.Printf("     ID: %s\n   NAME: %s\n   CATE: %s\n   TAGS: %s\n",
        lf.Id, lf.Name, lf.Category, lf.Tags)
    fmt.Printf("   DESC: %s\n", lf.Desc)
    codeLines := strings.Split(lf.Content, "\n")
    for i, line := range codeLines {
        if i == 0 {
            fmt.Println("CONTENT:")
            fmt.Println("    " + line)
        } else {
            fmt.Println("    " + line)
        }
    }
}

func (lf Leaf) PrintToFile(fpath string) error {
    f, err := os.OpenFile(fpath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0660)
    if err != nil {
        return err
    }
    defer f.Close()

    f.WriteString("Id:       " + lf.Id + "\n")
    f.WriteString("Category: " + lf.Category + "\n")
    f.WriteString("Tags:     " + lf.Tags + "\n")
    f.WriteString("Desc:     " + lf.Desc + "\n")

    codeLines := strings.Split(lf.Content, "\n")
    for i, line := range codeLines {
        if i == 0 {
            f.WriteString("Content:  " + line + "\n")
        } else {
            f.WriteString(CodePrefixSpace + line + "\n")
        }
    }
    return nil
}

func (lf *Leaf) ReadFromFile(fpath string) error {
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
            lf.Id = strings.TrimSpace(line[len("Id:"):])
        } else if strings.HasPrefix(line, "Category:") {
            lf.Category = strings.TrimSpace(line[len("Category:"):])
        } else if strings.HasPrefix(line, "Tags:") {
            lf.Tags = strings.TrimSpace(line[len("Tags:"):])
        } else if strings.HasPrefix(line, "Desc:") {
            lf.Desc = strings.TrimSpace(line[len("Desc:"):])
            isDescLine = true
            isCodeLine = false
        } else if strings.HasPrefix(line, "Content:") {
            lf.Content = strings.TrimSpace(line[len("Content:"):])
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
                lf.Desc = lf.Desc + "\n" + descLine
                //fmt.Println("lf.Desc:", lf.Desc)
            }

            if isCodeLine {
                var codeLine string
                if strings.HasPrefix(line, CodePrefixSpace) {
                    codeLine = line[len(CodePrefixSpace):]
                } else {
                    codeLine = strings.TrimSpace(line)
                }
                lf.Content = lf.Content + "\n" + codeLine
            }
        }
    }
    return nil
}
