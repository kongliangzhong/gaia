package main

import (
    "bufio"
    "os"
    "os/exec"
    "fmt"
    "strings"
    "strconv"
    "path/filepath"
    "io/ioutil"
)

var MetaFile string = "meta"
var Depencies string = "deps"

type ProjectType int

const (
    NODE ProjectType = iota
    TYPESCRIPT
    JAVA
    SCALA
    GO
    SHELL
)

type Executor struct {
    File string
    Type ProjectType
    TmpDir  string
    MainFile string
    Command string
}

func newExecutor(file string) *Executor {
    return &Executor{File: file}
}

func (executor *Executor) Execute() {
    executor.setType()
    filesMap := executor.parseFile()
    executor.generateTmpProject(filesMap)
    executor.generateBuildScript()
    executor.buildAndRun()
}

func (executor *Executor) setType() {
    ext := filepath.Ext(executor.File)
    switch ext {
    case ".js":
        executor.Type = NODE
    case ".ts":
        executor.Type = TYPESCRIPT
    case ".java":
        executor.Type = JAVA
    case ".scala":
        executor.Type = SCALA
    case ".go":
        executor.Type = GO

    default:
        fmt.Println("unsupport script file:" + executor.File)
        os.Exit(-1)
    }
}

/**
*
*/
func (executor *Executor) parseFile() map[string]string {
    f, err := os.Open(executor.File)
    if err != nil {
        panic(err)
    }
    defer f.Close()

    scanner := bufio.NewScanner(f)
    _, mainFile := filepath.Split(executor.File)
    currentFileName := mainFile
    executor.MainFile = mainFile
    filesMap := make(map[string]string)
    fileContent := ""
    for scanner.Scan() {
        line := scanner.Text()
        lineTrimed := strings.TrimSpace(line)

        if strings.HasPrefix(lineTrimed, "#!") {
            continue
        }

        if lineTrimed == "/***" {
            currentFileName = MetaFile
            // fmt.Println("parse meta info")
            fileContent = ""
        } else if lineTrimed == "*/" && currentFileName == MetaFile {
            filesMap[currentFileName] = fileContent
            currentFileName = mainFile
            fileContent = ""
        } else if strings.HasPrefix(lineTrimed, "/**#") && strings.HasSuffix(lineTrimed, "#*/") {
            filesMap[currentFileName] = fileContent
            len := len([]rune(lineTrimed))
            currentFileName = strings.TrimSpace(lineTrimed[4 : len-3])
            fileContent = ""

            // setup MainFile
            mainFileContent := filesMap[mainFile]
            if strings.TrimSpace(mainFileContent) == "" {
                executor.MainFile = currentFileName
            }
        } else{
            fileContent += line + "\n"
        }
    }

    filesMap[currentFileName] = fileContent
    return filesMap
}

func (executor *Executor) generateTmpProject(fileMap map[string]string) {
    projectDir := getTmpProjectDir()
    if projectDir == "" {
        fmt.Println("/tmp dir is full, please do some cleanup.")
        os.Exit(-1)
    }

    os.Mkdir(projectDir, os.ModePerm)
    fmt.Println("generate temp project in: " + projectDir)
    sourceFiles := []string{}
    for k, v := range fileMap {
        v = strings.TrimSpace(v)
        if v == "" {
            continue
        }
        file := projectDir + "/" + k
        sourceFiles = append(sourceFiles, k)
        f, err := os.Create(file)
        if err != nil {
            panic(err)
        }
        defer f.Close()

        f.WriteString(v)
    }

    executor.TmpDir = projectDir
}

func getTmpProjectDir() string {
    dir := "/tmp/gaia-tmp-"
    projectDir := ""
    for i := 0; i < 10000; i++ {
        projectDir = dir + strconv.Itoa(i)
        if _, err := os.Stat(projectDir); os.IsNotExist(err) {
            break
        }
    }
    return projectDir
}

func (executor *Executor) generateBuildScript() {
    metaFile := executor.TmpDir + "/" + MetaFile
    fileContent, err := ioutil.ReadFile(metaFile)
    if err != nil {
        return
    }

    fileLines := strings.Split(string(fileContent), "\n")
    switch executor.Type {
    case NODE:
        generateNodeFile(fileLines, executor.TmpDir)
        executor.Command = "npm install; node " + executor.MainFile
    case SCALA:
        generateSbtFile(fileLines, executor.TmpDir)
        appFile := refactorScalaMainFile(executor.TmpDir, executor.MainFile)
        executor.MainFile = appFile
        executor.Command = "sbt run"
    default:
        fmt.Println("executor type not implemented yet:" + string(executor.Type))
        os.Exit(-1)
    }
}

// TODO: add default props:
func generateNodeFile(fileLines []string, projectDir string) {
    generateJsonLine := func(k, v string) string {
        k = strings.Replace(k, "\"", "", -1)
        k = strings.Replace(k, "'", "", -1)
        k = strings.TrimSpace(k)
        v = strings.Replace(v, "\"", "", -1)
        v = strings.Replace(v, "'", "", -1)
        v = strings.TrimSpace(v)

        result := "\"" + k + "\": \"" + v + "\", \n"
        return result
    }

    removeTrailingComma := func(str string) string {
        res := strings.TrimSpace(str)
        res = strings.TrimRight(res, ",")
        return res + "\n"
    }

    generateJsonObject := func(name string, fields []string) string {
        if len(fields) == 0 {
            return ""
        }

        result := "\"" + name + "\": { \n"
        for _, item := range fields {
            parts := strings.Split(item, ":")
            result += generateJsonLine(parts[0], parts[1])
        }

        result = removeTrailingComma(result)
        result += "}, \n"
        return result
    }

    content := "{ \n"
    objectMap := make(map[string][]string)

    for _, line := range fileLines {
        line = strings.TrimSpace(line)
        if strings.Index(line, ":=") > 0 {
            parts := strings.Split(line, ":=")
            jsonLine := generateJsonLine(parts[0], parts[1])
            content += jsonLine
        } else if strings.Index(line, "+=") > 0 {
            parts := strings.Split(line, "+=")
            k := strings.TrimSpace(parts[0])
            if Depencies == k || "dependencies" == k || "devDependencies" == k {
                depsList := objectMap[Depencies]
                depsList = append(depsList, parts[1])
                objectMap["dependencies"] = depsList
            } else {
                fieldList := objectMap[k]
                fieldList = append(fieldList, parts[1])
                objectMap[k] = fieldList
            }
        } else {
            // fmt.Println("invalid meta info line: " + line)
            continue
        }
    }

    for name, fields := range objectMap {
        objStr := generateJsonObject(name, fields)
        content += objStr
    }

    content = removeTrailingComma(content)
    content += "} \n"

    pkgFile, _ := os.Create(projectDir + "/package.json")
    defer pkgFile.Close()

    pkgFile.WriteString(content)
}

// TODO: setup default fields.
func generateSbtFile(props []string, projectDir string) {
    fmt.Println("props", props)

    replaceDepsName := func(line string) string {
        res := line
        allSpaceTrimed := strings.Replace(line, " ", "", -1)
        if strings.HasPrefix(allSpaceTrimed, "deps+=") {
            res = "libraryDependencies" + strings.TrimLeft(res, "deps")
        }

        return res
    }

    lines := ""
    for _, line := range props {
        line = strings.TrimSpace(line)
        line = replaceDepsName(line)
        lines += line + "\n"
    }

    sbtFile, _ := os.Create(projectDir + "/build.sbt")
    defer sbtFile.Close()
    sbtFile.WriteString(lines)
}

func refactorScalaMainFile(tmpDir, mainFile string) string {
    parseMainClassName := func() string {
        className := strings.TrimRight(mainFile, ".scala")
        if strings.Index(className, "-") > 0 {
            parts := strings.Split(className, "-")
            className = ""
            for _, part := range parts {
                className += strings.Title(part)
            }
        } else {
            className = strings.Title(className)
        }
        return className
    }

    f := tmpDir + "/" + mainFile
    fileContent, err := ioutil.ReadFile(f)
    if err != nil {
        fmt.Println(err)
        os.Exit(-1)
    }

    fileLines := strings.Split(string(fileContent), "\n")
    newFileLines := ""
    appFileName := parseMainClassName()
    inClassHead := true
    for _, line := range fileLines {
        line = strings.TrimSpace(line)
        if strings.HasPrefix(line, "package ") || strings.HasPrefix(line, "import ") {
            newFileLines += line + "\n"
        } else {
            if inClassHead {
                inClassHead = false
                classDeclaration := "object " + appFileName + " extends App {"
                newFileLines += classDeclaration + "\n"
            } else {
                newFileLines += line + "\n"
            }
        }
    }

    newFileLines += "}\n"

    os.Remove(f)
    appFile, _ := os.Create(tmpDir + "/" + appFileName + ".scala")
    defer appFile.Close()
    appFile.WriteString(newFileLines)
    return appFileName + ".scala"
}

func (executor *Executor) buildAndRun() {
    fmt.Printf("executor: %+v\n", executor)
    commands := strings.Split(executor.Command, ";")
    for _, cmdStr := range commands {
        cmdStr = strings.TrimSpace(cmdStr)
        if cmdStr == "" {
            continue
        }
        fmt.Println("run command:", cmdStr)
        parts := strings.Split(cmdStr, " ")
        length := len(parts)
        cmd := exec.Command(parts[0], parts[1:length]...)
        cmd.Dir = executor.TmpDir
        out, err := cmd.CombinedOutput()
        if err != nil {
            fmt.Printf("error: %s\n", err)
            os.Exit(-1)
        }
        fmt.Print(string(out))
    }
}
