package main

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNoCells(t *testing.T) {
	file := `{
		"metadata": {
			"hello": "world"
		},
		"nbformat": 4,
		"nbformat_minor": 1,
		"cells": []
	}`
	var notebook Notebook
	require.NoError(t, json.Unmarshal([]byte(file), &notebook))
	require.Equal(t, Notebook{
		Metadata:     map[string]interface{}{"hello": "world"},
		Version:      4,
		MinorVersion: 1,
		Cells:        []NotebookCell{},
	}, notebook)
	validateRoundTrip(t, file)
}

func TestRawCell(t *testing.T) {
	file := `{
		"metadata": {
			"hello": "world"
		},
		"nbformat": 4,
		"nbformat_minor": 1,
		"cells": [
			{
				"cell_id": "123",
				"cell_type": "raw",
				"metadata": {"format": "mime/type"},
				"source": "[some nbformat output text]"
			}
		]
	}`
	var notebook Notebook
	require.NoError(t, json.Unmarshal([]byte(file), &notebook))
	require.Equal(t, Notebook{
		Metadata:     map[string]interface{}{"hello": "world"},
		Version:      4,
		MinorVersion: 1,
		Cells: []NotebookCell{{
			Type: RawCellType,
			Raw: &RawCell{
				Cell: Cell{Type: RawCellType, Id: "123", Metadata: map[string]interface{}{"format": "mime/type"}, Source: Source{isList: false, Data: []string{"[some nbformat output text]"}}},
			},
		}},
	}, notebook)
	validateRoundTrip(t, file)
}

func TestMarkdownCell(t *testing.T) {
	file := `{
		"metadata": {
			"hello": "world"
		},
		"nbformat": 4,
		"nbformat_minor": 1,
		"cells": [
			{
				"cell_id": "123",
				"cell_type": "markdown",
				"metadata": {},
				"source": ["# Header", "body"]
			}
		]
	}`
	var notebook Notebook
	require.NoError(t, json.Unmarshal([]byte(file), &notebook))
	require.Equal(t, Notebook{
		Metadata:     map[string]interface{}{"hello": "world"},
		Version:      4,
		MinorVersion: 1,
		Cells: []NotebookCell{{
			Type: MarkdownCellType,
			Markdown: &MarkdownCell{
				Cell: Cell{
					Type:     MarkdownCellType,
					Id:       "123",
					Metadata: map[string]interface{}{},
					Source:   Source{isList: true, Data: []string{"# Header", "body"}}},
			},
		}},
	}, notebook)
	validateRoundTrip(t, file)
}

func TestCodeCell(t *testing.T) {
	file := `{
		"metadata": {
			"hello": "world"
		},
		"nbformat": 4,
		"nbformat_minor": 1,
		"cells": [
			{
				"cell_id": "123",
				"cell_type": "code",
				"metadata": {"collapsed": true},
				"execution_count": 1,
				"source": ["import intake", "import requests"]
			}
		]
	}`
	one := 1
	var notebook Notebook
	require.NoError(t, json.Unmarshal([]byte(file), &notebook))
	require.Equal(t, Notebook{
		Metadata:     map[string]interface{}{"hello": "world"},
		Version:      4,
		MinorVersion: 1,
		Cells: []NotebookCell{{
			Type: CodeCellType,
			Code: &CodeCell{
				Cell: Cell{
					Type:     CodeCellType,
					Id:       "123",
					Metadata: map[string]interface{}{"collapsed": true},
					Source:   Source{isList: true, Data: []string{"import intake", "import requests"}}},
				ExecutionCount: &one,
			},
		}},
	}, notebook)
	validateRoundTrip(t, file)
}

func validateRoundTrip(t *testing.T, file string) {
	var notebook Notebook
	require.NoError(t, json.Unmarshal([]byte(file), &notebook))
	out, err := json.Marshal(notebook)
	require.NoError(t, err)
	require.JSONEq(t, file, string(out))

}
