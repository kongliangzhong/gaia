package main

type Store interface {
    Add(node Node) error
    AddAlias(from, to string) error
    Update(node Node) error
    Append(id string, extraContent string) error
    Search(category string, keywords []string) []Node
    Remove(id string) error
    GetById(id string) (Node, error)
    GetStats() Stats
    GetAlias() map[string]string
    ListCategories() map[string][]string
    ListNodes(names []string) []Node
}

type Stats struct {
    CategorySize int
    NodeSize     int
    TagSize      int
    Name0Size    int
}
