package entities

type Video struct {
	ID       int64  `json:"id"`
	FileName string `json:"file_name"`
	FilePath string `json:"file_path"`
}
