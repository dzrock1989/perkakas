package authorization

type Claims struct {
	UserUUID          string    `json:"user_uuid"`
	UserName          string    `json:"user_name"`
	Name              string    `json:"name"`
	Nrp               string    `json:"nrp"`
	KepolisianUUID    string    `json:"kepolisian_uuid"`
	KepolisianLevel   string    `json:"kepolisian_level"`
	Created           int64     `json:"created"`
	IsSuperadmin      bool      `json:"is_superadmin,omitempty"`
	Roles             []Role    `json:"roles"`
	Wilayahs          []Wilayah `json:"wilayahs"`
	IsFaceRecogActive bool      `json:"is_face_recog_active"`
	IsFaceRecogSetup  bool      `json:"is_face_recog_setup"`
	DirektoratId      string    `json:"direktorat_id"`
	Direktorat        string    `json:"direktorat"`
	SubDirektoratId   string    `json:"sub_direktorat_id"`
	SubDirektorat     string    `json:"sub_direktorat"`
	Impersonate       bool      `json:"impersonate,omitempty"`
}

type Role struct {
	Uuid string `json:"uuid"`
	Name string `json:"name"`
}

type Wilayah struct {
	Uuid string `json:"uuid"`
	Name string `json:"name"`
}

type Permission map[string]map[string]string

var keys = []string{
	"uuid",
	"wilayah",
	"wilayah_id",
	"tahun",
	"jenis",
	"key",
}
