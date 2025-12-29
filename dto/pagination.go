package dto

type Pagination struct {
	Page       int `json:"page"`        // halaman saat ini
	Limit      int `json:"limit"`       // jumlah item per halaman
	TotalRows  int `json:"total_rows"`  // total item di database
	TotalPages int `json:"total_pages"` // total halaman
}
