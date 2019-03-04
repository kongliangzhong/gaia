package main

type Store interface {
    Add(node Node) error
    AddAlias(from, to string) error
    Update(node Node) error
    Append(id string, extraContent string) error
    Search(category string, tagStr string) []Node
    Remove(id string) error
    GetById(id string) (Node, error)
    GetStats() Stats
    GetAlias() map[string]string
}

type Stats struct {
    TotalSize int
    AllCates     []string
    AllTags      []string
    CateTagsMap  map[string][]string
    CateNumMap   map[string]int
    TagCatesMap  map[string][]string
    TagNumMap    map[string]int
}
