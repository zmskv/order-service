package dto

type ItemModel struct {
	ID          int64  `db:"id"`
	OrderUID    string `db:"order_uid"`
	ChrtID      int    `db:"chrt_id"`
	TrackNumber string `db:"track_number"`
	Price       int    `db:"price"`
	RID         string `db:"rid"`
	Name        string `db:"name"`
	Sale        int    `db:"sale"`
	Size        string `db:"size"`
	TotalPrice  int    `db:"total_price"`
	NmID        int    `db:"nm_id"`
	Brand       string `db:"brand"`
	Status      int    `db:"status"`
}
