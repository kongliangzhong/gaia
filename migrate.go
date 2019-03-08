package main

import (
    "bufio"
    "encoding/base64"
    "errors"
    // "fmt"
    "os"
    "strings"
)

type CodeSegment struct {
    Id, Category, Tags, Desc, Code string
}

func getAll(file string) (csArray []CodeSegment, err error) {
    f, err := os.Open(file)
    if err != nil {
        return
    }
    defer f.Close()

    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        line := scanner.Text()
        csg, err := strToCodeSegment(line)
        if err == nil {
            csArray = append(csArray, csg)
        }
    }
    return
}

func strToCodeSegment(str string) (cs CodeSegment, err error) {
    flds := strings.Split(str, "|")
    if len(flds) != 5 {
        err = errors.New("parse segemnt str failed: " + str)
        return
    }
    id := flds[0]
    cate := flds[1]
    tags := flds[2]
    desc := flds[3]
    code := flds[4]

    descBs, err := base64.StdEncoding.DecodeString(desc)
    if err != nil {
        return
    }

    codeBs, err := base64.StdEncoding.DecodeString(code)
    if err != nil {
        return
    }

    desc = string(descBs)
    code = string(codeBs)

    cs = CodeSegment{id, cate, tags, desc, code}
    return
}

// func main() {
//     oldFile := "/Users/kongliang/.gaia/data/code-snippets.txt"
//     allOldCs, err := getAll(oldFile)
//     if err != nil {
//         fmt.Println("error:", err)
//         os.Exit(1)
//     }

//     fmt.Printf("allOldCs size: %d \n", len(allOldCs))

//     jsonStore := newJsonFileStore("/Users/kongliang/.gaia/data/data.json")
//     for _, cs := range allOldCs {
//         node := Node{
//             Name: cs.Id,
//             Category: cs.Category,
//             Tags: cs.Tags,
//             Desc: cs.Desc,
//             Content: cs.Code,
//         }

//         err = jsonStore.Add(node)
//         if err != nil {
//             fmt.Println("save node failed:", err)
//         }
//     }

// }
