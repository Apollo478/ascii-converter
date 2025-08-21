package converter

type Ascii_t struct {
	AsciiChars [][]rune
	RgbColors [][]Rgb
}
type Rgb struct {
	R uint32
	G uint32
	B uint32
}
