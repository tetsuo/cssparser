# cssparser

Parses CSS declarations from raw input.

## Installation

```sh
go get github.com/tetsuo/cssparser
```

## Usage

```go
css := []byte(`
    color: blue;
    font-size: 12px;
`)
nodes, _ := cssparser.Parse(css)
for _, n := range nodes {
  fmt.Printf("Decl @ %+v: %s = %s\n", n.Position, n.Property, n.Value)
}
```

Output:

```
Decl @ {Start:{Line:2 Column:7} End:{Line:2 Column:19}}: color = blue
Decl @ {Start:{Line:3 Column:7} End:{Line:3 Column:23}}: font-size = 12px
```

## License

MIT
