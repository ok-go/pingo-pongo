package pingo_pongo

type Config struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

var (
	cfg            = Config{400, 300}
	positionFirst  = Point{50, cfg.Height / 2}
	positionSecond = Point{cfg.Width - 50, cfg.Height / 2}
)
