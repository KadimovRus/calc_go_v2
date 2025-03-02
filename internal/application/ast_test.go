package application_test

import (
	"github.com/KadimovRus/calc_go/internal/application"
	"strconv"
	"testing"
)

func TestParseAST(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		wantErr  bool
	}{
		{"3+5", "(+ 3 5)", false},
		{"10-2", "(- 10 2)", false},
		{"2*3", "(* 2 3)", false},
		{"8/4", "(/ 8 4)", false},
		{"(1+2)*3", "(* (+ 1 2) 3)", false},
		{"", "", true},
		{"5+", "", true},
		{"(3+4", "", true},
	}

	for _, tt := range tests {
		node, err := application.ParseAST(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("ParseAST(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
		}
		if err == nil {
			got := astToString(node)
			if got != tt.expected {
				t.Errorf("ParseAST(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		}
	}
}

func astToString(node *application.ASTNode) string {
	if node.IsLeaf {
		return strconv.FormatFloat(node.Value, 'f', -1, 64)
	}
	return "(" + node.Operator + " " + astToString(node.Left) + " " + astToString(node.Right) + ")"
}
