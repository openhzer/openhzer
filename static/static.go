package static

import "embed"

var (
	//go:embed admin/css
	AdminCss embed.FS

	//go:embed admin/fonts
	AdminFonts embed.FS

	//go:embed admin/js
	AdminJs embed.FS
	/*
		//go:embed admin
		Admin embed.FS

		//go:embed proot
		Proot embed.FS

		//go:embed css
		Css embed.FS

		//go:embed fonts
		Fonts embed.FS

		//go:embed js
		Js embed.FS*/
)
