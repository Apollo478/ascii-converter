package converter

type Ascii_t struct {
	AsciiChars [][]rune
	RgbColors [][]Rgb
}
type Rgb struct {
	R uint32
	G uint32
	B uint32
	A uint32
}

type Options struct {
	Width int
	Height int
	AspectRatio float64 
	Compression int 

	UseColor bool
	UseAlpha bool
	BlendPrev bool 

	CharSet string
	Invert bool

	FitTerminal bool
	ClearScreen bool

	ShowProgress bool
	Parallel bool
}
