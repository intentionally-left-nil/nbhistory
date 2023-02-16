package main

import (
	"encoding/json"
	"errors"
)

type Notebook struct {
	Metadata     map[string]interface{} `json:"metadata"`
	Version      int                    `json:"nbformat"`
	MinorVersion int                    `json:"nbformat_minor"`
	Cells        []NotebookCell         `json:"cells"`
}

type CellType string

const (
	RawCellType      CellType = "raw"
	MarkdownCellType CellType = "markdown"
	CodeCellType     CellType = "code"
)

type NotebookCell struct {
	Type     CellType
	Raw      *RawCell
	Markdown *MarkdownCell
	Code     *CodeCell
}

func (c *NotebookCell) UnmarshalJSON(data []byte) error {
	var underlying struct {
		Type CellType `json:"cell_type"`
	}
	err := json.Unmarshal(data, &underlying)
	if err != nil {
		return nil
	}
	c.Type = underlying.Type
	switch underlying.Type {
	case RawCellType:
		c.Raw = &RawCell{}
		err = json.Unmarshal(data, c.Raw)
	case MarkdownCellType:
		c.Markdown = &MarkdownCell{}
		err = json.Unmarshal(data, c.Markdown)
	case CodeCellType:
		c.Code = &CodeCell{}
		err = json.Unmarshal(data, c.Code)
	default:
		err = errors.New("Unexpected cell_type")
	}
	return err
}

func (c NotebookCell) MarshalJSON() ([]byte, error) {
	switch c.Type {
	case RawCellType:
		return json.Marshal(*c.Raw)
	case MarkdownCellType:
		return json.Marshal(*c.Markdown)
	case CodeCellType:
		return json.Marshal(*c.Code)
	default:
		return nil, errors.New("Unexpected cell_type")
	}
}

type Cell struct {
	Id       string                 `json:"cell_id"`
	Type     CellType               `json:"cell_type"`
	Metadata map[string]interface{} `json:"metadata"`
	Source   Source                 `json:"source"`
}

type RawCell struct {
	Cell
	Attachments *map[string]interface{} `json:"attachments,omitempty"`
}

type MarkdownCell struct {
	Cell
	Attachments *map[string]interface{} `json:"attachments,omitempty"`
}

type CodeCell struct {
	Cell
	// We don't actually want the output, ever. So just don't unmarshal/marshal that field
	// Outputs        interface{} `json:"outputs"`
	ExecutionCount *int `json:"execution_count"`
}

type Source struct {
	Data   []string
	isList bool
}

func (s *Source) UnmarshalJSON(data []byte) error {
	var raw interface{}
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}
	switch v := raw.(type) {
	case string:
		s.Data = []string{v}
		s.isList = false
	case []interface{}:
		s.Data = make([]string, len(v))
		for i := range v {
			val, ok := v[i].(string)
			if !ok {
				return errors.New("Expected array of strings for Source")
			}
			s.Data[i] = val
		}
		s.isList = true
	case []string:
		s.Data = v
		s.isList = true
	default:
		return errors.New("Expected string or array for Source")
	}
	return nil
}

func (s Source) MarshalJSON() ([]byte, error) {
	if s.isList {
		return json.Marshal(s.Data)
	} else if len(s.Data) == 1 {
		return json.Marshal(s.Data[0])
	} else {
		return nil, errors.New("Cannot marshal source into a single string")
	}
}
