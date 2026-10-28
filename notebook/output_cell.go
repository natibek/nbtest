package notebook

import (
	"encoding/json"
	"fmt"
)

// Output cells
type outputType string

const (
	ErrorOutput       outputType = "error"
	StreamOutput      outputType = "stream"
	DisplayDataOutput outputType = "display_data"
	ExecResultOutput  outputType = "execute_result"
	ExecInputOutput   outputType = "execute_input"
)

type Outputs []Output

type Output interface {
	Type() outputType
}

func (o *Error) Type() outputType         { return ErrorOutput }
func (o *Stream) Type() outputType        { return StreamOutput }
func (o *DisplayData) Type() outputType   { return DisplayDataOutput }
func (o *ExecuteResult) Type() outputType { return ExecResultOutput }

func (o *Outputs) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("outputs: %w", err)
	}
	if raw == nil {
		*o = nil
		return nil
	}

	outputs := make(Outputs, 0, len(raw))
	for i, r := range raw {
		var head struct {
			OutputType outputType `json:"output_type"`
		}
		if err := json.Unmarshal(r, &head); err != nil {
			return fmt.Errorf("outputs[%d]: %w", i, err)
		}

		var output Output
		switch head.OutputType {
		case StreamOutput:
			output = new(Stream)
		case ErrorOutput:
			output = new(Error)
		case DisplayDataOutput:
			output = new(DisplayData)
		case ExecResultOutput, ExecInputOutput:
			output = new(ExecuteResult)
		default:
			return fmt.Errorf("outputs[%d]: unknown output_type %q", i, head.OutputType)
		}

		if err := json.Unmarshal(r, output); err != nil {
			return fmt.Errorf("outputs[%d]: %w", i, err)
		}
		outputs = append(outputs, output)
	}
	*o = outputs
	return nil
}

// {
//     'output_type': 'error',
//     'ename' : str,   # Exception name, as a string
//     'evalue' : str,  # Exception value, as a string

//	    # The traceback will contain a list of frames,
//	    # represented each as a string.
//	    'traceback' : list,
//	}
type Error struct {
	CellId     string     `json:"id"`
	OutputType outputType `json:"output_type"`
	Ename      string     `json:"ename"`
	Evalue     string     `json:"evalue"`
	Traceback  []string   `json:"traceback"`
}

//	{
//	  "output_type" : "stream",
//	  "name" : "stdout", # or stderr
//	  "text" : ["multiline stream text"],
//	}
type streamOutputName string

const (
	Stdout streamOutputName = "stdout"
	Stderr streamOutputName = "stderr"
)

type Stream struct {
	CellId     string           `json:"id"`
	OutputType outputType       `json:"output_type"`
	Name       streamOutputName `json:"name"`
	Text       multiLineString  `json:"text"`
}

//	{
//	    "output_type": "display_data",
//	    "data": {
//	        "text/plain": "[multiline text data]",
//	        "image/png": "[base64-encoded-multiline-png-data]",
//	        "application/json": {
//	            # JSON data is included as-is
//	            "key1": "data",
//	            "key2": ["some", "values"],
//	            "key3": {"more": "data"},
//	        },
//	        "application/vnd.exampleorg.type+json": {
//	            # JSON data, included as-is, when the mime-type key ends in +json
//	            "key1": "data",
//	            "key2": ["some", "values"],
//	            "key3": {"more": "data"},
//	        },
//	    },
//	    "metadata": {
//	        "image/png": {
//	            "width": 640,
//	            "height": 480,
//	        },
//	    },
//	}
type data struct {
	Text  multiLineString `json:"text/plain"`
	Image []byte          `json:"image/png"`
	Json  map[string]any  `json:"application/json"`
}

type imgMetadata struct {
	Width  uint `json:"width"`
	Height uint `json:"height"`
}

type metadata struct {
	ImgMetadata imgMetadata `json:"image/png"`
}

type DisplayData struct {
	CellId     string     `json:"id"`
	OutputType outputType `json:"output_type"`
	Data       data       `json:"data"`
	MetaData   metadata   `json:"metadata"`
}

//	{
//	    "output_type": "execute_result",
//	    "execution_count": 42,
//	    "data": {
//	        "text/plain": "[multiline text data]",
//	        "image/png": "[base64-encoded-multiline-png-data]",
//	        "application/json": {
//	            # JSON data is included as-is
//	            "json": "data",
//	        },
//	    },
//	    "metadata": {
//	        "image/png": {
//	            "width": 640,
//	            "height": 480,
//	        },
//	    },
//	}
type ExecuteResult struct {
	DisplayData
	ExecutionCount uint       `json:"execution_count"`
	OutputType     outputType `json:"output_type"`
}
