package main

type gallery struct {
	Title   string
	Date    string
	Grid    string
	Zoom    bool
	GridCSS string
	Images  []string
}

var galleryIndex = map[string]gallery{
	"untitled-no1": {
		Title: "untitled #1", Date: "2025-01-18", Grid: "3x3", Zoom: true,
		Images: []string{
			"x-25n11_sharp.png", "x-dante_menu.png", "x-security_guards_xmas.png",
			"x-nypd_pizza.png", "x-sehaj_punjabi_deli.png", "x-w4_station.png",
			"x-w4_basketball.png", "x-stg_23n11_highline_overpass.png", "x-birds_7th_ave_xmas.png",
		},
	},
	"untitled-no2": {
		Title: "untitled #2", Date: "2025-05-31", Grid: "3x3", Zoom: true,
		Images: []string{
			"stealth.png", "22n10_hang_prints.png", "neal.png",
			"blossom_1.png", "rebecca.png", "blossom_2.png",
			"fof_guy.png", "xstreet.png", "reena.png",
		},
	},
	"untitled-no3": {
		Title: "untitled #3", Date: "2025-10-19", Grid: "3x3",
		Images: []string{
			"l16.jpg", "l5.jpg", "l11.jpg",
			"l3.jpg", "l12.jpg", "l10.jpg",
			"l17.jpg", "l8.jpg", "l9.jpg",
		},
	},
	"architecture-no1": {
		Title: "architecture #1", Date: "2025-10-19", Grid: "3x1",
		Images: []string{"p10.jpg", "p11.jpg", "p1.jpg"},
	},
	"architecture-no2": {
		Title: "architecture #2", Date: "2025-10-19", Grid: "1x1",
		GridCSS: "repeat(auto-fit, minmax(calc(var(--page-width) / 4), 0.5fr))",
		Images:  []string{"p3.jpg"},
	},
	"human-no1": {
		Title: "human #1", Date: "2025-10-19", Grid: "3x3",
		Images: []string{
			"l4.jpg", "l1.jpg", "l0.jpg",
			"l6.jpg", "l15.jpg", "l14.jpg",
			"l13.jpg", "l7.jpg", "l2.jpg",
		},
	},
	"human-no2": {
		Title: "human #2", Date: "2025-10-19", Grid: "3x1",
		Images: []string{"p21.jpg", "p20.jpg", "p12.jpg"},
	},
	"countryside-no1": {
		Title: "countryside #1", Date: "2025-10-19", Grid: "3x3",
		Images: []string{
			"l1.jpg", "l2.jpg", "l15.jpg",
			"l6.jpg", "l14.jpg", "l8.jpg",
			"l17.jpg", "l9.jpg", "l0.jpg",
		},
	},
}
