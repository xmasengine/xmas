// Package xfig implements an XFig vector format parser.
package xfig

import "io"
import "image"
import "fmt"
import "strings"
import "errors"
import "unicode"
import "strconv"
import "bufio"

type Comment string

func ScanString[S interface{ ~string }](state fmt.ScanState) (S, error) {
	tok, err := state.Token(false, func(r rune) bool { return r != '\n' })
	if err != nil {
		return "", err
	}
	return S(tok), nil
}

func (c *Comment) Scan(state fmt.ScanState, verb rune) error {
	ch, _, _ := state.ReadRune()
	if ch != '#' {
		state.UnreadRune()
		return fmt.Errorf("not a comment")
	}
	tok, err := ScanString[Comment](state)
	if err != nil {
		return err
	}
	*c = tok
	return nil
}

func ScanOneOf[S interface{ ~string }](state fmt.ScanState, options ...S) (S, error) {
	tok, err := ScanString[S](state)
	if err != nil {
		return tok, err
	}
	for _, option := range options {
		if option == tok {
			return S(tok), nil
		}
	}
	return "", fmt.Errorf("For %T: got %s, expected one of %v", tok, tok, options)
}

func ScanOneOfTo[S interface{ ~string }](to *S, state fmt.ScanState, options ...S) error {
	tok, err := ScanOneOf[S](state, options...)
	if err != nil {
		return err
	}
	if to != nil {
		*to = tok
	}
	return nil
}

type Orientation string

const Landscape Orientation = "Landscape"
const Portrait Orientation = "Portrait"

func (o *Orientation) Scan(state fmt.ScanState, r rune) error {
	return ScanOneOfTo(o, state, Landscape, Portrait)
}

type Justification string

const Center Justification = "Center"
const FlushLeft Justification = "Flush Left"

func (j *Justification) Scan(state fmt.ScanState, r rune) error {
	return ScanOneOfTo(j, state, Center, FlushLeft)
}

type Units string

const MetricUnits = "Metric"
const USAUnits = "Inches"

func (u *Units) Scan(state fmt.ScanState, r rune) error {
	return ScanOneOfTo(u, state, MetricUnits, USAUnits)
}

type PaperSize string

const (
	LetterSize  = "Letter"
	LegalSize   = "Legal"
	LedgerSize  = "Ledger"
	TabloidSize = "Tabloid"
	ASize       = "A"
	BSize       = "B"
	CSize       = "C"
	DSize       = "D"
	ESize       = "E"
	A4Size      = "A4"
	A3Size      = "A3"
	A2Size      = "A2"
	A1Size      = "A1"
	A0Size      = "A0"
	B5Size      = "B5"
)

func (p *PaperSize) Scan(state fmt.ScanState, r rune) error {
	return ScanOneOfTo(p, state,
		LetterSize,
		LegalSize,
		LedgerSize,
		TabloidSize,
		ASize,
		BSize,
		CSize,
		DSize,
		ESize,
		A4Size,
		A3Size,
		A2Size,
		A1Size,
		A0Size,
		B5Size,
	)
}

type Pages string

const (
	SinglePage    Pages = "Single"
	MultiplePages Pages = "Multiple"
)

func (p *Pages) Scan(state fmt.ScanState, r rune) error {
	return ScanOneOfTo(p, state, SinglePage, MultiplePages)
}

type Color int

const (
	BackgroundColor Color = -3
	NoColor         Color = -2
	DefaultColor    Color = -1
)

func ScanInt[S interface{ ~int }](state fmt.ScanState) (S, error) {
	tok, err := state.Token(true,
		func(r rune) bool {
			if r == '-' {
				return true
			}
			return unicode.IsNumber(r)
		})
	if err != nil {
		return 0, err
	}
	i, err := strconv.Atoi(string(tok))
	if err != nil {
		return 0, err
	}
	return S(i), nil
}

func ScanFloat[S interface{ ~float64 }](state fmt.ScanState) (S, error) {
	tok, err := state.Token(true,
		func(r rune) bool {
			if r == '-' || r == '.' {
				return true
			}
			return unicode.IsNumber(r)
		})
	if err != nil {
		return 0, err
	}
	f, err := strconv.ParseFloat(string(tok), 0)
	if err != nil {
		return 0, err
	}
	return S(f), nil
}

func (c *Color) Scan(state fmt.ScanState, r rune) error {
	newc, err := ScanInt[Color](state)
	if err != nil {
		return err
	}
	if newc < -3 || newc > 512 {
		return fmt.Errorf("Color out of range: %d", newc)
	}
	*c = newc
	return nil
}

type Header struct {
	Orientation
	Justification
	Units
	PaperSize
	Magnification float64
	Pages
	Transparent Color
	Comment
	Resolution int // units /inch
	Origin     int
}

type Parser struct {
	Errors  []error
	Scanner *bufio.Scanner
	Keep    bool
}

func NewParser(rd io.Reader) *Parser {
	return &Parser{Scanner: bufio.NewScanner(rd)}
}

func (p Parser) AllErrors() error {
	if len(p.Errors) < 1 {
		return nil
	}
	return errors.Join(p.Errors...)
}

func (p *Parser) Next() (string, bool) {
	if p.Keep {
		p.Keep = false
		return p.Scanner.Text(), true
	} else {
		ok := p.Scanner.Scan()
		if !ok {
			return "", false
		} else {
			return p.Scanner.Text(), true
		}
	}
}

func (s *Parser) Scanf(form string, args ...any) int {
	line, ok := s.Next()
	if !ok {
		s.Errors = append(s.Errors, io.EOF)
		return 0
	}
	found, err := fmt.Sscanf(line, form, args...)
	if err != nil {
		s.Errors = append(s.Errors, err)
	}
	return found
}

func (s *Parser) ScanComment(c *Comment) bool {
	line, ok := s.Next()
	if !ok {
		return false
	}
	if len(line) > 0 && line[0] == '#' {
		*c = Comment(line)
		return true
	} else {
		s.Keep = true
		return false
	}
}

func (s *Parser) Scanln(args ...any) int {
	line, ok := s.Next()
	if !ok {
		s.Errors = append(s.Errors, io.EOF)
		return 0
	}
	found, err := fmt.Sscanln(line, args...)
	if err != nil {
		s.Errors = append(s.Errors, err)
	}
	return found
}

func (s *Parser) Scan(args ...any) int {
	line, ok := s.Next()
	if !ok {
		s.Errors = append(s.Errors, io.EOF)
		return 0
	}
	found, err := fmt.Sscan(line, args...)
	if err != nil {
		s.Errors = append(s.Errors, err)
	}
	return found
}

func (h *Header) ParseWith(s *Parser) error {
	var comment Comment

	s.ScanComment(&comment)
	if !strings.HasPrefix(string(comment), "FIG 3.2") {
		return fmt.Errorf("Incorrect header: %s", comment)
	} else {
		println("comment>", comment, "<")
	}

	s.Scanln(&h.Orientation)
	s.Scanln(&h.Justification)
	s.Scanln(&h.Units)
	s.Scanln(&h.PaperSize)
	s.Scanln(&h.Magnification)
	s.Scanln(&h.Pages)
	s.Scanln(&h.Transparent)
	comment = ""
	for s.ScanComment(&comment) {
		h.Comment += "\n" + comment
	}

	s.Scanln(&h.Resolution, &h.Origin)
	return s.AllErrors()
}

type HexColor string

func (h *HexColor) Scan(state fmt.ScanState, verb rune) error {
	ch, _, _ := state.ReadRune()
	if ch != '#' {
		state.UnreadRune()
		return fmt.Errorf("not a color hex")
	}
	tok, err := ScanString[HexColor](state)
	if err != nil {
		return err
	}
	*h = tok
	return nil
}

type UserDefinedColor struct {
	ColorNumber Color
	HexColor    HexColor
}

func (u *UserDefinedColor) ParseWith(p *Parser) error {
	l := p.Scanln(&u.ColorNumber, &u.HexColor)
	if l != 2 {
		return fmt.Errorf("Could not parse user defined color")
	}
	return nil
}

type Common struct {
	LineStyle     int
	LineThickness int
	PenColor      Color
	FillColor     Color
	Depth         int
	PenStyle      int
	AreaFill      int
	StyleVal      float64
}

func (c *Common) ParseWith(p *Parser) error {
	l := p.Scan(&c.LineStyle, &c.LineThickness,
		&c.PenColor, &c.FillColor, &c.Depth, &c.PenStyle,
		&c.AreaFill, &c.StyleVal,
	)
	if l != 9 {
		return fmt.Errorf("Could not parse common block")
	}
	return nil
}

type Point image.Point

func (p *Point) Scan(state fmt.ScanState, verb rune) error {
	x, err := ScanInt[int](state)
	if err != nil {
		return err
	}
	y, err := ScanInt[int](state)
	if err != nil {
		return err
	}
	p.X = x
	p.Y = y
	return nil
}

type FloatPoint image.Point

func (p *FloatPoint) Scan(state fmt.ScanState, verb rune) error {
	x, err := ScanFloat[float64](state)
	if err != nil {
		return err
	}
	y, err := ScanFloat[float64](state)
	if err != nil {
		return err
	}
	p.X = int(x)
	p.Y = int(y)
	return nil
}

type Rectangle image.Rectangle

func (r *Rectangle) Scan(state fmt.ScanState, verb rune) error {
	x1, err := ScanInt[int](state)
	if err != nil {
		return err
	}
	y1, err := ScanInt[int](state)
	if err != nil {
		return err
	}
	x2, err := ScanInt[int](state)
	if err != nil {
		return err
	}
	y2, err := ScanInt[int](state)
	if err != nil {
		return err
	}
	r.Min.X = int(x1)
	r.Min.Y = int(y1)
	r.Max.X = int(x2)
	r.Max.Y = int(y2)
	return nil
}

type Arrow struct {
	Type      int
	Style     int
	Thickness float64
	Width     float64
	Height    float64
}

type Arc struct {
	Common

	CapStyle      int
	Direction     int
	ForwardArrow  int
	BackwardArrow int
	Center        FloatPoint
	Points        [3]Point
	Arrows        [2]Arrow
}

func (a *Arc) ParseWith(p *Parser) error {
	err := a.Common.ParseWith(p)
	if err != nil {
		return err
	}

	l := p.Scan(&a.CapStyle, &a.Direction,
		&a.ForwardArrow, &a.BackwardArrow,
		&a.Center,
	)
	if l != 5 {
		return fmt.Errorf("Could not parse Arc")
	}
	for i := 0; i < len(a.Points); i++ {
		p.Scan(&a.Points[i])
	}
	if a.ForwardArrow != 0 {
		p.Scan(&a.Arrows[0])
	}
	if a.BackwardArrow != 0 {
		p.Scan(&a.Arrows[1])
	}
	return nil
}

type EllipseType int

const (
	EllipseRadii     EllipseType = 1
	EllipseDiameters EllipseType = 2
	CircleRadius     EllipseType = 3
	CircleDialeter   EllipseType = 4
)

type Ellipse struct {
	SubType EllipseType
	Common
	Direction int
	Angle     float64
	Center    Point
	Radius    Point
	Start     Point
	End       Point
}

func (p *EllipseType) Scan(state fmt.ScanState, verb rune) error {
	psf, err := ScanInt[EllipseType](state)
	if err != nil {
		return err
	}
	*p = psf
	return nil
}

func (e *Ellipse) ParseWith(p *Parser) error {
	p.Scanln(
		&e.SubType,
		&e.Common,
		&e.Direction,
		&e.Angle,
		&e.Center,
		&e.Radius,
		&e.Start,
		&e.End,
	)
	return p.AllErrors()
}

type PolylineType int

const (
	PolylinePolyline PolylineType = 1
	BoxPolyline      PolylineType = 2
	PolygonPolyline  PolylineType = 3
	ArcBoxPolyline   PolylineType = 4
	PicturePolyline  PolylineType = 5
)

func (p *PolylineType) Scan(state fmt.ScanState, verb rune) error {
	psf, err := ScanInt[PolylineType](state)
	if err != nil {
		return err
	}
	*p = psf
	return nil
}

type Picture struct {
	Flipped bool
	File    string
}

func (p *Picture) Scan(state fmt.ScanState, verb rune) error {
	flip, err := ScanInt[int](state)
	if err != nil {
		return err
	}
	p.Flipped = (flip == 1)
	file, err := ScanString[string](state)
	if err != nil {
		return err
	}
	p.File = file
	return nil
}

type Polyline struct {
	SubType PolylineType
	Common

	CapStyle      int
	Direction     int
	ForwardArrow  int
	BackwardArrow int
	NPoints       int
	Points        []Point
	Arrows        [2]Arrow
	Picture       Picture
}

func (l *Polyline) ParseWith(p *Parser) error {
	p.Scanln(
		&l.SubType,
		&l.Common,
		&l.CapStyle,
		&l.Direction,
		&l.ForwardArrow,
		&l.BackwardArrow,
		&l.NPoints,
	)

	if l.ForwardArrow != 0 {
		p.Scanln(&l.Arrows[0])
	}
	if l.BackwardArrow != 0 {
		p.Scanln(&l.Arrows[1])
	}
	if l.SubType == PicturePolyline {
		p.Scanln(&l.Picture)
	}

	le := l.NPoints
	fmt.Printf("Polyline %#v\n", l)
	l.Points = make([]Point, le)
	points := make([]any, le)
	for i := 0; i < le; i++ {
		points[i] = &l.Points[i]
	}
	p.Scanln(points...)
	return p.AllErrors()
}

type SplineType int

const (
	OpenApproximatedSpline   SplineType = 0
	ClosedApproximatedSpline SplineType = 1
	OpenInterpretedSpline    SplineType = 2
	ClosedInterpretedSpline  SplineType = 3
	OpenXSpline              SplineType = 4
	ClosedXSpline            SplineType = 5
)

func (p *SplineType) Scan(state fmt.ScanState, verb rune) error {
	psf, err := ScanInt[SplineType](state)
	if err != nil {
		return err
	}
	*p = psf
	return nil
}

type Spline struct {
	SubType SplineType
	Common

	CapStyle      int
	Direction     int
	ForwardArrow  int
	BackwardArrow int
	NPoints       int
	Points        []Point
	Arrows        [2]Arrow
	Picture       Picture
	Controls      []float64
}

func (l *Spline) ParseWith(p *Parser) error {
	p.Scanln(
		&l.SubType,
		&l.Common,
		&l.CapStyle,
		&l.Direction,
		&l.ForwardArrow,
		&l.BackwardArrow,
		&l.NPoints,
	)

	if l.ForwardArrow != 0 {
		p.Scanln(&l.Arrows[0])
	}
	if l.BackwardArrow != 0 {
		p.Scanln(&l.Arrows[1])
	}

	le := l.NPoints
	l.Points = make([]Point, le)
	points := make([]any, le)
	for i := 0; i < le; i++ {
		points[i] = &l.Points[i]
	}
	l.Controls = make([]float64, le)
	controls := make([]any, le)
	for i := 0; i < le; i++ {
		controls[i] = &l.Controls[i]
	}
	p.Scanln(controls...)
	return p.AllErrors()
}

type TextType int

const (
	LeftJustifiedText   TextType = 0
	CenterJustifiedText TextType = 1
	RightJustifiedText  TextType = 2
)

func (p *TextType) Scan(state fmt.ScanState, verb rune) error {
	psf, err := ScanInt[TextType](state)
	if err != nil {
		return err
	}
	*p = psf
	return nil
}

type Text struct {
	SubType   TextType
	Color     Color
	Depth     int
	PenStyle  int
	Font      int
	FontSize  float64
	Angle     float64
	FontFlags FontFlag
	Height    int
	Length    int
	Origin    Point
	Text      EscapedText
}

func (t *Text) ParseWith(p *Parser) error {
	p.Scanln(
		&t.SubType,
		&t.Color,
		&t.Depth,
		&t.PenStyle,
		&t.Font,
		&t.FontSize,
		&t.Angle,
		&t.FontFlags,
		&t.Height,
		&t.Length,
		&t.Origin,
		&t.Text,
	)
	return p.AllErrors()
}

type EscapedText string

func (e *EscapedText) Scan(state fmt.ScanState, verb rune) error {
	escaped, err := ScanString[string](state)
	if err != nil {
		return err
	}
	un, err := strconv.Unquote("\"" + escaped + "\"")
	if err != nil {
		return err
	}
	*e = EscapedText(un)
	return nil
}

type Compound struct {
	Bounds  Rectangle
	Objects []Object
}

func (sub *Object) ParseByCode(p *Parser) error {
	switch sub.ObjectCode {
	case ColorCode:
		return sub.UserDefinedColor.ParseWith(p)
	case EllipseCode:
		return sub.Ellipse.ParseWith(p)
	case PolylineCode:
		return sub.Polyline.ParseWith(p)
	case SplineCode:
		return sub.Spline.ParseWith(p)
	case TextCode:
		return sub.Text.ParseWith(p)
	case ArcCode:
		return sub.Arc.ParseWith(p)
	case CompoundCode:
		return sub.Compound.ParseWith(p)
	default:
		return fmt.Errorf("Unknown object code")
	}
}

func (c *Compound) ParseWith(p *Parser) error {
	p.Scanln(&c.Bounds)
	for {
		sub := Object{}
		p.Scan(&sub.ObjectCode)
		println("compound", sub.ObjectCode)
		if sub.ObjectCode == -6 {
			return nil // Compound is done
		}
		if len(p.Errors) > 0 {
			return errors.Join(p.Errors...)
		}
		var comment Comment
		for p.ScanComment(&comment) {
			sub.Comment += "\n" + comment
		}
		err := sub.ParseByCode(p)
		if err != nil {
			return err
		}
		c.Objects = append(c.Objects, sub)
	}
}

// ObjectCode is the type of object
type ObjectCode int

const (
	ColorCode    ObjectCode = 0
	EllipseCode  ObjectCode = 1
	PolylineCode ObjectCode = 2
	SplineCode   ObjectCode = 3
	TextCode     ObjectCode = 4
	ArcCode      ObjectCode = 5
	CompoundCode ObjectCode = 6
	FirstCode    ObjectCode = ColorCode
	LastCode     ObjectCode = CompoundCode
)

func (c *ObjectCode) Scan(state fmt.ScanState, r rune) error {
	newc, err := ScanInt[ObjectCode](state)
	if err != nil {
		return err
	}
	if newc < FirstCode || newc > LastCode {
		return fmt.Errorf("Object Code out of range: %d", newc)
	}
	*c = newc
	return nil
}

type Object struct {
	Comment Comment
	ObjectCode
	UserDefinedColor
	Ellipse
	Polyline
	Spline
	Text
	Arc
	Compound
}

func (o *Object) ParseWith(p *Parser) error {
	var comment Comment
	for p.ScanComment(&comment) {
		o.Comment += "\n" + comment
	}
	p.Scan(&o.ObjectCode)

	fmt.Printf("Object %#v\n", o)
	if len(p.Errors) > 0 {
		return errors.Join(p.Errors...)
	}
	err := o.ParseByCode(p)
	if err != nil {
		return err
	}
	return nil
}

func (f *Fig) ParseWith(p *Parser) error {
	err := f.Header.ParseWith(p)
	if err != nil {
		return err
	}

	for {
		var o Object
		err = o.ParseWith(p)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		f.Objects = append(f.Objects, o)
	}
}

type Fig struct {
	Header
	Objects []Object
}

type FontFlag int

const (
	RigidFontFlag FontFlag = 1 << iota
	SpecialFontFlag
	PostscriptFontFlag
	HiddenFontFlag
)

func (p *FontFlag) Scan(state fmt.ScanState, verb rune) error {
	psf, err := ScanInt[FontFlag](state)
	if err != nil {
		return err
	}
	*p = psf
	return nil
}

type LatexFont int

func (p *LatexFont) Scan(state fmt.ScanState, verb rune) error {
	psf, err := ScanInt[LatexFont](state)
	if err != nil {
		return err
	}
	*p = psf
	return nil
}

const (
	DefaultLatexFont LatexFont = iota
	RomanLatexFont
	BoldLatexFont
	ItalicLatexFont
	SansSerifLatexFont
	TypewriterLatexFont
)

type PostscriptFont int

func (p *PostscriptFont) Scan(state fmt.ScanState, verb rune) error {
	psf, err := ScanInt[PostscriptFont](state)
	if err != nil {
		return err
	}
	*p = psf
	return nil
}

const (
	DefaultPostscriptFont PostscriptFont = iota - 1
	TimesRoman
	TimesItalic
	TimesBold
	TimesBoldItalic
	AvantGardeBook
	AvantGardeBookOblique
	AvantGardeDemi
	AvantGardeDemiOblique
	BookmanLight
	BookmanLightItalic
	BookmanDemi
	BookmanDemiItalic
	Courier
	CourierOblique
	CourierBold
	CourierBoldOblique
	Helvetica
	HelveticaOblique
	HelveticaBold
	HelveticaBoldOblique
	HelveticaNarrow
	HelveticaNarrowOblique
	HelveticaNarrowBold
	HelveticaNarrowBoldOblique
	NewCenturySchoolbookRoman
	NewCenturySchoolbookItalic
	NewCenturySchoolbookBold
	NewCenturySchoolbookBoldItalic
	PalatinoRoman
	PalatinoItalic
	PalatinoBold
	PalatinoBoldItalic
	Symbol
	ZapfChanceryMediumItalic
	ZapfDingbats
)
