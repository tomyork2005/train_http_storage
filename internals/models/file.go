package models

type File struct {
	Id         int64  `json:"id" db:"id"`
	Alias      string `json:"alias" db:"alias"`
	PathToFile string `json:"path_to_file" db:"path_to_file"`
	UserId     int64  `json:"user_id" db:"user_id"`
}

func (f *File) Validate() error {
	return validate.Struct(f)
}
