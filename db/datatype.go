package db

type TableItemType struct {
	Content string
}

//=============================================================================
// Content
//=============================================================================

type ContentField struct {
	ID    string
	Value string
}

type Content struct {
	ID                 string
	ItemType           string
	ContractID         string `validate:"required"`
	CustomerID         string `validate:"required"`
	PublishedStartTime string
	PublishedEndTime   string
	PreviewKey         string
	Status             string
	CreatedAt          string
	UpdatedAt          string
	ContentFields      []ContentField
}
