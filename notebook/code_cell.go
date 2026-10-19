package notebook

type Cell interface {
	Match(other Cell) (bool, error)
}

type CodeCell struct {
	Source    string
	Outputs   []*Output
	ExecCount int
	Metadata  map[string]any
	CellId    uint
}

// have different types of errors

func (c CodeCell) Match(other Cell) (bool, error) {
	target, ok := other.(CodeCell)
	if !ok {
		return false, nil
	}

	if c.Source != target.Source {
		return false, nil
	}

	if len(c.Outputs) != len(target.Outputs) {
		return false, nil
	}

	// match outputs
	for i, output := range c.Outputs {
		match, err := (*output).Match(*target.Outputs[i])
		if err != nil {
			return false, err
		}
		if !match {
			return false, nil
		}
	}
	return true, nil
}

type Output interface {
	Match(other Output) (bool, error)
}

// 3 types of output
// error
//
//	{
//	  "output_type" : "error",
//	  'ename' : str,   # Exception name, as a string
//	  'evalue' : str,  # Exception value, as a string
//	  'traceback' : list,
//	}
type ErrorOutput struct {
	Ename     string
	Evalue    string
	Traceback []string
}

// stream
//
//	{
//	  "output_type" : "stream",
//	  "name" : "stdout", # or stderr
//	  "text" : ["multiline stream text"],
//	}
type streamOutputName string

const Stdout streamOutputName = "stdout"
const Stderr streamOutputName = "stderr"

type multilineText *[]string

type StreamOutput struct {
	Name streamOutputName
	Text multilineText
}

// execute_result / display_data
//
//	{
//	  "output_type" : "execute_result" | "display_data",
//	  "execution_count": 42, # if "execute_result"
//	  "data" : {
//	    "text/plain" : ["multiline text data"],
//	    "image/png": ["base64-encoded-png-data"],
//	    "application/json": {
//	      # JSON data is included as-is
//	      "json": "data",
//	    },
//	  },
//	  "metadata" : {
//	    "image/png": {
//	      "width": 640,
//	      "height": 480,
//	    },
//	  },
type imgData *[]byte

type data struct {
	Text  multilineText
	Image imgData
	Json  map[string]any
}

type imgMetadata struct {
	Width  uint
	Height uint
}

type metadata struct {
	ImgMetadata imgMetadata
}

type DisplayData struct {
	Data     data
	MetaData metadata
}

type ExecuteResult struct {
	ExecutionCount uint
	DisplayData
}
