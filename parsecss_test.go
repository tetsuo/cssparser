package cssparser_test

import (
	"bytes"
	"testing"

	"github.com/tetsuo/cssparser"
)

func TestParse(t *testing.T) {
	tests := []struct {
		desc     string
		input    string
		expected []cssparser.Node
	}{
		{
			desc:  "empty input",
			input: "",
		},
		{
			desc:  "whitespace only",
			input: "   \n\t\r ",
		},
		{
			desc:  "single property",
			input: "color: red;",
			expected: []cssparser.Node{
				{
					Property: []byte("color"),
					Value:    []byte("red"),
				},
			},
		},
		{
			desc:  "multiple declarations",
			input: "color: red; font-size: 16px;",
			expected: []cssparser.Node{
				{Property: []byte("color"), Value: []byte("red")},
				{Property: []byte("font-size"), Value: []byte("16px")},
			},
		},
		{
			desc:  "quoted string value",
			input: `font-family: "Open Sans";`,
			expected: []cssparser.Node{
				{Property: []byte("font-family"), Value: []byte(`"Open Sans"`)},
			},
		},
		{
			desc:  "with parentheses in value",
			input: `background-image: url(data:image/png;base64,iVBORw0KGgo...);`,
			expected: []cssparser.Node{
				{Property: []byte("background-image"), Value: []byte(`url(data:image/png;base64,iVBORw0KGgo...)`)},
			},
		},
		{
			desc:  "attribute selector property",
			input: "input[type=text]: border: 1px solid #ccc;",
			expected: []cssparser.Node{
				{Property: []byte("input[type=text]"), Value: []byte("border: 1px solid #ccc")},
			},
		},
		{
			desc:  "semicolon optional at end",
			input: "color: red",
			expected: []cssparser.Node{
				{Property: []byte("color"), Value: []byte("red")},
			},
		},
		{
			desc:  "skips invalid lines",
			input: "color: red; ???invalid;;; font-size: 12px;",
			expected: []cssparser.Node{
				{Property: []byte("color"), Value: []byte("red")},
				{Property: []byte("font-size"), Value: []byte("12px")},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			nodes, err := cssparser.Parse([]byte(tt.input))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(nodes) != len(tt.expected) {
				t.Fatalf("got %d nodes, want %d", len(nodes), len(tt.expected))
			}
			for i, n := range nodes {
				exp := tt.expected[i]
				if !bytes.Equal(n.Property, exp.Property) {
					t.Errorf("node %d: got property %q, want %q", i, n.Property, exp.Property)
				}
				if !bytes.Equal(n.Value, exp.Value) {
					t.Errorf("node %d: got value %q, want %q", i, n.Value, exp.Value)
				}
			}
		})
	}
}
