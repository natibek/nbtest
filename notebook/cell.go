package notebook

import (
	"encoding/json"
	"fmt"
	"strings"
)

// https://nbformat.readthedocs.io/en/latest/format_description.html
// Cell ids

// Since the 4.5 schema release, all cells have an id field which must be a string
// of length 1-64 with alphanumeric, -, and _ as legal characters to use.
// These ids must be unique to any given Notebook following the nbformat spec.

//	{
//	    "cell_type": "type",
//	    "metadata": {},
//	    "source": "single string or [list, of, strings]",
//	}
type Cell interface {
	Type() cellType
}

func (o *Code) Type() cellType     { return CodeCell }
func (o *Markdown) Type() cellType { return MarkdownCell }
func (o *Raw) Type() cellType      { return RawCell }

//	{
//	    "cell_type": "code",
//	    "execution_count": 1,  # integer or null
//	    "metadata": {
//	        "collapsed": True,  # whether the output of the cell is collapsed
//	        "scrolled": False,  # any of true, false or "auto"
//	    },
//	    "source": "[some multi-line code]",
//	    "outputs": [
//	        {
//	            # list of output dicts (described below)
//	            "output_type": "stream",
//	            # ...
//	        }
//	    ],
//	}

type cellType string

const (
	CodeCell     cellType = "code"
	MarkdownCell cellType = "markdown"
	RawCell      cellType = "raw"
)

type multiLineString string

func (mls *multiLineString) UnmarshalJSON(data []byte) error {
	var parts []string
	if err := json.Unmarshal(data, &parts); err == nil {
		// nbformat keeps the newline on each element, so join with ""
		*mls = multiLineString(strings.Join(parts, ""))
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("source must be a string or list of strings: %w", err)
	}
	*mls = multiLineString(s)
	return nil
}

type Code struct {
	CellId    string          `json:"id"`
	CellType  cellType        `json:"cell_type"`
	ExecCount *uint           `json:"execution_count"`
	Metadata  map[string]any  `json:"metadata"`
	Source    multiLineString `json:"source"`
	Outputs   Outputs         `json:"outputs"`
}

//	{
//	    "cell_type": "raw",
//	    "metadata": {
//	        # the mime-type of the target nbconvert format.
//	        # nbconvert to formats other than this will exclude this cell.
//	        "format": "mime/type"
//	    },
//	    "source": ["some nbformat output text"],
//	}
type Raw struct {
	CellId   string          `json:"id"`
	CellType cellType        `json:"cell_type"`
	Metadata map[string]any  `json:"metadata"`
	Source   multiLineString `json:"source"`
}

//	{
//	    "cell_type": "markdown",
//	    "metadata": {},
//	    "source": "[multi-line *markdown*]",
//	}
type Markdown struct {
	CellId   string          `json:"id"`
	CellType cellType        `json:"cell_type"`
	Source   multiLineString `json:"source"`
}

type Cells []Cell

func (c *Cells) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("cells: %w", err)
	}
	if raw == nil {
		*c = nil
		return nil
	}

	cells := make(Cells, 0, len(raw))
	for i, r := range raw {
		var head struct {
			CellType cellType `json:"cell_type"`
		}
		if err := json.Unmarshal(r, &head); err != nil {
			return fmt.Errorf("cell[%d]: %w", i, err)
		}

		var cell Cell
		switch head.CellType {
		case CodeCell:
			cell = new(Code)
		case RawCell:
			cell = new(Raw)
		case MarkdownCell:
			cell = new(Markdown)
		default:
			return fmt.Errorf("cells[%d]: unknown cell type: %q", i, head.CellType)
		}

		if err := json.Unmarshal(r, cell); err != nil {
			return fmt.Errorf("cell[%d]: %w", i, err)
		}
		cells = append(cells, cell)
	}
	*c = cells
	return nil
}
