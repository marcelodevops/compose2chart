package convert

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateChart_Minimal(t *testing.T) {
	compose := `version: "3"
services:
  web:
    image: nginx:alpine
    ports:
      - "8080:80"
`
	tmp := os.TempDir()
	inPath := filepath.Join(tmp, "compose-test.yml")
	ioutil.WriteFile(inPath, []byte(compose), 0644)
	outDir := filepath.Join(tmp, "chart-out")
	os.RemoveAll(outDir)

	opts := Options{ComposePath: inPath, OutDir: outDir, ChartName: "testchart", AppVersion: "0.1", Version: "0.1"}
	if err := GenerateChart(opts); err != nil {
		t.Fatal(err)
	}

	// verify Chart.yaml exists
	if _, err := os.Stat(filepath.Join(outDir, "Chart.yaml")); os.IsNotExist(err) {
		t.Fatalf("Chart.yaml not created")
	}
}