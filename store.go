package main

type Store interface {
    Add(lf Leaf) error
    AddAlias(from, to string) error
    Update(lf Leaf) error
    Append(id string, extraContent string) error
    Search(category string, tagStr string) []Leaf
    Remove(id string) error
    GetById(id string) (Leaf, error)
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
