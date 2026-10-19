package notebook

type notebookMetadata struct {
}

type Notebook struct {
	Meatadata notebookMetadata
	Cells     *[]Cell
}
