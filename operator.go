package main

import (
    "errors"
    "fmt"
    "io/ioutil"
    "os"
    "os/exec"
    "strings"
    "github.com/satori/go.uuid"
    "strconv"
)

const resultDelimiter = "--------------------------------------------------------"

type Operator struct {
    err   error
    store Store
}

func newOperator(store Store) *Operator {
    return &Operator{nil, store}
}

func (op *Operator) Add(node Node) {
    // var validate = func() {
    //     node.Content = strings.TrimSpace(node.Content)
    //     if node.Content == "" {
    //         op.err = errors.New("content can not be empty.")
    //     }

    //     if node.Category == "" && node.Tags == "" {
    //         op.err = errors.New("category and tags can not be both empty.")
    //     }

    //     if strings.Contains(node.Category, "|") || strings.Contains(node.Tags, "|") {
    //         op.err = errors.New("category and tagStr can not contains '|' charactor.")
    //     }
    //     return
    // }

    // if op.err != nil {
    //     return
    // }

    // validate()
    if op.err != nil {
        return
    }

    op.err = op.store.Add(node)
}

func (op *Operator) AddAlias(from, to string) {
    op.err = op.store.AddAlias(from, to)
}

func (op *Operator) Update(node Node) {
    if node.Id == "" {
        op.err = errors.New("id is empty")
        return
    }

    op.err = op.store.Update(node)
}

func (op *Operator) Append(id string, extraContent string) {
    if id == "" || extraContent == "" {
        op.err = errors.New("id or content is nil")
        return
    }

    op.err = op.store.Append(id, extraContent)
}

func (op *Operator) Search(category string, tags string) {
    matchedNode := op.store.Search(category, tags)
    size := len(matchedNode)
    if size > 10 {
        fmt.Println("Found", size, "matched content segments, print first 10 as below:")
    } else {
        fmt.Println("Found", size, "matched content segments, print as below:")
    }
    for i, node := range matchedNode {
        if i < 10 {
            fmt.Println(resultDelimiter)
            node.PrintToScreen()
        } else {
            break
        }
    }
    fmt.Println(resultDelimiter)
}

func (op *Operator) Remove(id string) {
    if op.err != nil {
        return
    }
    op.err = op.store.Remove(id)
}

func (op *Operator) Merge(ids ...string) {
    var arrContains = func(arr []string, str string) bool {
        for _, s := range arr {
            if str == s {
                return true
            }
        }
        return false
    }

    var cate string
    var allTags []string
    var desc string
    var content string
    for i, id := range ids {
        node, err := op.store.GetById(id)
        if err != nil {
            op.err = err
            return
        }

        if i == 0 {
            cate = node.Category
        }

        desc = desc + "\n" + node.Desc
        content = content + "\n" + node.Content

        if node.Category != cate {
            op.err = errors.New("categorys are not equal, can not merge.")
            return
        }

        tags := strings.Split(node.Tags, ",")
        for _, t := range tags {
            if !arrContains(allTags, t) {
                allTags = append(allTags, t)
            }
        }
    }

    desc = strings.TrimSpace(desc)
    content = strings.TrimSpace(content)
    allTagsStr := strings.Join(allTags, ",")
    mergedNode := Node{Id: "", Category: cate, Tags: allTagsStr, Desc: desc, Content: content}
    for _, id := range ids {
        op.Remove(id)
    }
    op.Add(mergedNode)
}

func (op *Operator) Edit(id string) {
    node, err := op.store.GetById(id)
    if err != nil {
        op.err = err
        return
    }

    tmpDir := os.TempDir()
    uuidObj, err := uuid.NewV4()
    tmpFileName := uuidObj.String()
    tmpFile, err := ioutil.TempFile(tmpDir, tmpFileName)
    if err != nil {
        op.err = err
        return
    }
    defer tmpFile.Close()

    node.PrintToFile(tmpFile.Name())

    path, err := exec.LookPath("vi")
    if err != nil {
        op.err = errors.New("Error while looking for vi: " + err.Error())
        return
    }

    cmd := exec.Command(path, tmpFile.Name())
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    err = cmd.Start()
    if err != nil {
        op.err = err
        return
    }

    err = cmd.Wait()

    if err != nil {
        op.err = err
        return
    }

    //fmt.Println("tmpFile: ", tmpFile.Name())
    err = (&node).ReadFromFile(tmpFile.Name())
    if err != nil {
        op.err = err
        return
    }
    //node.PrintToScreen()

    oldId := node.Id
    node.Id = ""
    op.Remove(oldId)
    op.Add(node)
}

func (op *Operator) ListAlias() {
    aliasMap := op.store.GetAlias()

    if len(aliasMap) == 0 {
        fmt.Println("No Alias Mapping")
        return
    } else {
        for k, v := range aliasMap {
            fmt.Printf("%s -> %s\n", k, v)
        }
    }
}

func (op *Operator) ListCates() {
    stats := op.store.GetStats()
    head := []string{"INDEX   ", "CATEGORY        ", "NODE-NUM     ", "TAGS"}
    index := 0
    format := fmt.Sprintf("%%-%ds%%-%ds%%-%ds%%-%ds\n", len(head[0]), len(head[1]), len(head[2]), len(head[3]))
    //fmt.Printf("%s%s%s%s\n", head[0], head[1], head[2], head[3])

    lineMax := 50
    for cate, tags := range stats.CateTagsMap {
        index ++
        num := stats.CateNumMap[cate]
        tagLines := []string{}
        line := ""
        for i, tag := range tags {
            if line == "" {
                line = tag
            } else {
                line = line + "," + tag
            }

            if len(line) > lineMax {
                if i != len(tags) - 1 {
                    line = line + ","
                }
                tagLines = append(tagLines, line)
                line = ""
            } else {
                if i == len(tags) - 1 {
                    if strings.HasSuffix(line, ",") {
                        line = line[:len(line)-1]
                    }
                    tagLines = append(tagLines, line)
                }
            }
        }

        for i, tagLine := range tagLines {
            if i == 0 {
                fmt.Printf(format, strconv.Itoa(index), cate, strconv.Itoa(num), tagLine)
            } else {
                formatNewLine := fmt.Sprintf("%%%ds\n", len(head[0]) + len(head[1]) + len(head[2]) + len(tagLine))
                fmt.Printf(formatNewLine, tagLine)
            }
        }

    }
}

func (op *Operator) ListTags() {
    stats := op.store.GetStats()
    head := []string{"INDEX    ", "TAG                    ", "NODE-NUM ", "CATEGORIES    "}
    index := 0
    format := fmt.Sprintf("%%-%ds%%-%ds%%-%ds%%-%ds\n", len(head[0]), len(head[1]), len(head[2]), len(head[3]))
    fmt.Printf("%s%s%s%s\n", head[0], head[1], head[2], head[3])
    for tag, cates := range stats.TagCatesMap {
        index ++
        num := stats.TagNumMap[tag]
        fmt.Printf(format, strconv.Itoa(index), tag, strconv.Itoa(num), strings.Join(cates, ","))
    }
}

func (op *Operator) Exec(file string) {
    executor := newExecutor(file)
    executor.Execute()
}
