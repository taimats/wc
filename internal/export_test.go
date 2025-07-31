package internal

var (
	CountWords = countWords
	ExecuteWC  = executeWC
	NewRecord  = newRecord
)

func NewRecordSlice(records ...record) []record {
	if len(records) == 0 {
		return []record{}
	}
	slice := make([]record, len(records))
	copy(slice, records)
	return slice
}
