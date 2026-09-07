package qr

import (
	"bytes"
	"encoding/xml"
	"io"
	"os"
	"path/filepath"
	"testing"
	"webtyp.com/fmt"
)

func TestSVGQuietZoneClamping(t *testing.T) {
	m, err := Encode("https://webtyp.com", LevelM)
	if err != nil {
		t.Fatalf("Encode error: %v", err)
	}

	svg1 := m.SVG(SVGOptions{QuietZone: 1})
	svg4 := m.SVG(SVGOptions{QuietZone: 4})

	if svg1 != svg4 {
		t.Fatalf("SVG with QuietZone=1 differed from QuietZone=4")
	}

	svg0 := m.SVG(SVGOptions{})
	if svg0 != svg4 {
		t.Fatalf("SVG with QuietZone=0 differed from QuietZone=4")
	}

	if m.Dark(-1, 0) {
		t.Errorf("Dark(-1, 0) = true; want false")
	}
	if m.Dark(m.Size(), 0) {
		t.Errorf("Dark(Size, 0) = true; want false")
	}
}

func TestSVGXMLStructureAndPathCount(t *testing.T) {
	m, err := Encode("hello world", LevelM)
	if err != nil {
		t.Fatalf("Encode error: %v", err)
	}

	svgStr := m.SVG(SVGOptions{})

	decoder := xml.NewDecoder(bytes.NewReader([]byte(svgStr)))
	pathCount := 0
	rectCount := 0

	for {
		tok, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Fatalf("XML parse error: %v", err)
		}
		switch se := tok.(type) {
		case xml.StartElement:
			if se.Name.Local == "path" {
				pathCount++
			}
			if se.Name.Local == "rect" {
				rectCount++
			}
		}
	}

	if rectCount != 1 {
		t.Errorf("got %d <rect> elements; want 1", rectCount)
	}
	if pathCount != 1 {
		t.Errorf("got %d <path> elements; want exactly 1", pathCount)
	}
}

func TestSVGOutputVisualArtifacts(t *testing.T) {
	outDir := filepath.Join("testdata", "out")
	if err := os.MkdirAll(outDir, 0755); err != nil {
		t.Fatalf("failed to create out dir: %v", err)
	}

	url := "https://192.168.1.5:8080/__webtyp/ca"
	levels := []struct {
		name  string
		level Level
	}{
		{"level_m.svg", LevelM},
		{"level_l.svg", LevelL},
		{"level_q.svg", LevelQ},
		{"level_h.svg", LevelH},
	}

	for _, l := range levels {
		m, err := Encode(url, l.level)
		if err != nil {
			t.Fatalf("Encode level %s error: %v", l.name, err)
		}
		svgData := m.SVG(SVGOptions{})
		outPath := filepath.Join(outDir, l.name)
		if err := os.WriteFile(outPath, []byte(svgData), 0644); err != nil {
			t.Fatalf("failed to write %s: %v", outPath, err)
		}
	}

	longData := fmt.Repeat("A", 300)
	mHigh, err := Encode(longData, LevelM)
	if err != nil {
		t.Fatalf("Encode long data error: %v", err)
	}
	highPath := filepath.Join(outDir, "high_version.svg")
	if err := os.WriteFile(highPath, []byte(mHigh.SVG(SVGOptions{})), 0644); err != nil {
		t.Fatalf("failed to write high_version.svg: %v", err)
	}
}
