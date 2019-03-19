package MapCutter

import (
	"bytes"
	"fmt"
	"github.com/jonas-p/go-shp"
	"github.com/opensatelliteproject/SatHelperApp/ImageProcessor/Projector"
	"golang.org/x/image/draw"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"image"
	"math"
	"strings"
	"unicode"
)

const rootIndexAlphabet = "abcdefghijklmnopqrstuvwxyz0123456789"
const defaultMargin = 5

type BorderSection struct {
	Name       string
	Bounds     shp.Box
	Properties map[string]string
}

type BSSearchTreeItem struct {
	CurrentName string
	Childs      []string
	Tree        map[string]*BSSearchTreeItem
}

func MakeBSSearchTreeItem(currentName string) *BSSearchTreeItem {
	return &BSSearchTreeItem{
		CurrentName: currentName,
		Tree:        make(map[string]*BSSearchTreeItem),
		Childs:      make([]string, 0),
	}
}

type MapCutter struct {
	sections     map[string]BorderSection
	searchTree   map[string]*BSSearchTreeItem
	marginPixels int
}

func MakeMapCutter(shapeFile string) (*MapCutter, error) {
	shape, err := shp.Open(shapeFile)
	if err != nil {
		return nil, err
	}
	defer shape.Close()

	return MakeMapDrawerFromShapes([]*shp.Reader{shape})
}

func MakeMapCutterFromFiles(shapeFiles []string) (*MapCutter, error) {
	readers := make([]*shp.Reader, 0)

	defer func() {
		for _, v := range readers {
			_ = v.Close()
		}
	}()

	for _, shapeFile := range shapeFiles {
		shape, err := shp.Open(shapeFile)
		if err != nil {
			return nil, err
		}
		readers = append(readers, shape)
	}

	return MakeMapDrawerFromShapes(readers)
}

func MakeMapDrawerFromShapes(shapes []*shp.Reader) (*MapCutter, error) {

	mc := &MapCutter{
		sections:     map[string]BorderSection{},
		searchTree:   map[string]*BSSearchTreeItem{},
		marginPixels: defaultMargin,
	}

	for _, shape := range shapes {
		// Cache all Border Sections
		fields := shape.Fields()
		for shape.Next() {
			n, p := shape.Shape()

			var poly *shp.Polygon

			switch v := p.(type) {
			case *shp.Polygon:
				poly = v
			}

			if poly == nil {
				continue
			}

			s := BorderSection{
				Bounds:     poly.BBox(),
				Properties: make(map[string]string),
			}

			for k, f := range fields {
				fieldBytes := f.Name[:]

				a := bytes.Split(fieldBytes, []byte{0})

				field := string(a[0])

				val := shape.ReadAttribute(n, k)
				s.Properties[field] = val
				if field == "name" {
					s.Name = val
				}
			}

			if s.Name != "" {
				mc.sections[nameToIndex(s.Name)] = s
			}
		}
	}

	mc.buildTree()

	return mc, nil
}

func MergeMapCutters(a, b *MapCutter) *MapCutter {
	mc := &MapCutter{
		sections:     map[string]BorderSection{},
		searchTree:   map[string]*BSSearchTreeItem{},
		marginPixels: defaultMargin,
	}

	// Copy sections
	for k, v := range a.sections {
		mc.sections[k] = v
	}
	for k, v := range b.sections {
		mc.sections[k] = v
	}

	// Build New Tree
	mc.buildTree()

	return mc
}

func (mc *MapCutter) buildTree() {
	// Initialize ROOT
	for _, v := range rootIndexAlphabet {
		c := string(v)
		mc.searchTree[c] = &BSSearchTreeItem{
			CurrentName: c,
			Tree:        make(map[string]*BSSearchTreeItem),
			Childs:      make([]string, 0),
		}
	}

	// Build Indexes
	for sectionName := range mc.sections {
		currentName := string(sectionName[0])
		remainingName := sectionName[1:]
		currentNode := mc.searchTree[string(currentName)]
		currentNode.Childs = append(currentNode.Childs, sectionName)

		// Start building the tree with remaining chars
		for len(remainingName) > 0 {
			currentName += string(remainingName[0])
			remainingName = remainingName[1:]

			if currentNode.Tree[currentName] == nil {
				currentNode.Tree[currentName] = MakeBSSearchTreeItem(currentName)
			}

			currentNode = currentNode.Tree[currentName]
			currentNode.Childs = append(currentNode.Childs, sectionName)
		}
	}
}

func (mc *MapCutter) SearchSection(name string) []string {
	name = nameToIndex(name)
	var currentNode *BSSearchTreeItem
	l := ""
	for _, v := range name {
		c := string(v)
		l += c
		if currentNode == nil {
			currentNode = mc.searchTree[l]
		} else {
			currentNode = currentNode.Tree[l]
		}

		if currentNode == nil {
			break
		}
	}

	if currentNode != nil {
		return currentNode.Childs
	}

	return make([]string, 0)
}

func (mc *MapCutter) GetSection(name string) (BorderSection, error) {
	s := mc.SearchSection(name)
	if len(s) == 0 {
		return BorderSection{}, fmt.Errorf("no such section with name %s", name)
	}

	if len(s) > 1 {
		return BorderSection{}, fmt.Errorf("there is more than one section that matches name %s", name)
	}

	return mc.sections[s[0]], nil
}

func (mc *MapCutter) CutMap(section string, img image.Image, gc Projector.ProjectionConverter) (image.Image, error) {
	s, err := mc.GetSection(section)
	if err != nil {
		return nil, err
	}

	// Convert BBox to Pixels
	X0, Y0 := gc.LatLon2XYf(s.Bounds.MinY, s.Bounds.MinX)
	X1, Y1 := gc.LatLon2XYf(s.Bounds.MaxY, s.Bounds.MaxX)

	minX := int(math.Min(X0, X1)) - mc.marginPixels
	minY := int(math.Min(Y0, Y1)) - mc.marginPixels
	maxX := int(math.Max(X0, X1)) + mc.marginPixels
	maxY := int(math.Max(Y0, Y1)) + mc.marginPixels

	if minX < 0 {
		minX = 0
	}

	if minY < 0 {
		minY = 0
	}

	if maxX > img.Bounds().Dx() {
		maxX = img.Bounds().Dx()
	}

	if maxY > img.Bounds().Dy() {
		maxY = img.Bounds().Dy()
	}

	// Slice image

	out := image.NewRGBA(image.Rect(0, 0, maxX-minX, maxY-minY))
	draw.Draw(out, out.Bounds(), img, image.Point{X: minX, Y: minY}, draw.Src)
	return out, nil
}

func (mc *MapCutter) CutMapMany(sections []string, img image.Image, gc Projector.ProjectionConverter) (image.Image, error) {
	if len(sections) == 0 {
		return nil, fmt.Errorf("no sections specified")
	}

	if len(sections) == 1 {
		return mc.CutMap(sections[0], img, gc)
	}

	minX := int(math.MaxInt32)
	maxX := int(math.MinInt32)
	minY := minX
	maxY := maxX

	for _, section := range sections {
		s, err := mc.GetSection(section)
		if err != nil {
			return nil, err
		}

		// Convert BBox to Pixels
		X0, Y0 := gc.LatLon2XYf(s.Bounds.MinY, s.Bounds.MinX)
		X1, Y1 := gc.LatLon2XYf(s.Bounds.MaxY, s.Bounds.MaxX)

		_minX := int(math.Min(X0, X1)) - mc.marginPixels
		_minY := int(math.Min(Y0, Y1)) - mc.marginPixels
		_maxX := int(math.Max(X0, X1)) + mc.marginPixels
		_maxY := int(math.Max(Y0, Y1)) + mc.marginPixels

		if _minX < minX {
			minX = _minX
		}

		if _maxX > maxX {
			maxX = _maxX
		}

		if _minY < minY {
			minY = _minY
		}

		if _maxY > maxY {
			maxY = _maxY
		}

		if minX < 0 {
			minX = 0
		}

		if minY < 0 {
			minY = 0
		}

		if maxX > img.Bounds().Dx() {
			maxX = img.Bounds().Dx()
		}

		if maxY > img.Bounds().Dy() {
			maxY = img.Bounds().Dy()
		}
	}

	// Slice image

	out := image.NewRGBA(image.Rect(0, 0, maxX-minX, maxY-minY))
	draw.Draw(out, out.Bounds(), img, image.Point{X: minX, Y: minY}, draw.Src)
	return out, nil
}

func nameToIndex(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, s)
	return strings.ToLower(result)
}
