package util

import (
	"archive/zip"
	"bufio"
	"bytes"
	"context"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/jonas-p/go-shp"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

type Charset string

const (
	UTF8    = Charset("UTF-8")
	GB18030 = Charset("GB18030")
)

// BoxStringToArray 'BOX(111.964430789849 30.9798174714992,116.727509619973 34.8957377544938)' to [4]float
func BoxStringToArray(box2d string) (bounds [4]float64) {

	result := strings.FieldsFunc(box2d, func(r rune) bool {
		if r == 'B' || r == 'O' || r == 'X' || r == '(' || r == ' ' || r == ',' || r == ')' {
			return true
		} else {
			return false
		}
	})
	bounds[0], _ = strconv.ParseFloat(result[0], 64)
	bounds[1], _ = strconv.ParseFloat(result[1], 64)
	bounds[2], _ = strconv.ParseFloat(result[2], 64)
	bounds[3], _ = strconv.ParseFloat(result[3], 64)
	return bounds
}

func ExeCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	stdout, _ := cmd.StdoutPipe()
	stderr, err := cmd.StderrPipe()
	var wg sync.WaitGroup
	wg.Add(2)
	go read(context.Background(), &wg, stdout)
	go read(context.Background(), &wg, stderr)
	err = cmd.Start()
	if err != nil {
		return err
	}
	wg.Wait()
	if err != nil {
		return err
	}
	return nil
}

func read(ctx context.Context, wg *sync.WaitGroup, std io.ReadCloser) {
	reader := bufio.NewReader(std)
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			readString, err := reader.ReadString('\n')
			if err != nil || err == io.EOF {
				return
			}
			fmt.Print(ConvertByte2String([]byte(readString), "GB18030"))
		}
	}
}

func ConvertByte2String(byte []byte, charset Charset) string {
	var str string
	switch charset {
	case GB18030:
		var decodeBytes, _ = simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
		str = string(decodeBytes)
	case UTF8:
		fallthrough
	default:
		str = string(byte)
	}
	return str
}

func ShapefileToGeojson(input string) ([]byte, error) {
	shape, err := shp.Open(input)
	if err != nil {
		return nil, err
	}
	defer func(shape *shp.Reader) {
		err := shape.Close()
		if err != nil {
			return
		}
	}(shape)

	// fields from the attribute table (DBF)
	fields := shape.Fields()
	fc := geojson.NewFeatureCollection()
	// loop through all features in the shapefile
	for shape.Next() {
		n, s := shape.Shape()
		var feature *geojson.Feature

		switch s.(type) {
		case *shp.Point:
			pt, _ := s.(*shp.Point)
			feature = geojson.NewFeature(orb.Point{pt.X, pt.Y})
		case *shp.PolyLine:
			line, _ := s.(*shp.PolyLine)
			if line.NumParts != 1 {
				fmt.Println("Warning: more than 1 part in polyline!")
			}
			var lines orb.LineString
			for _, pt := range line.Points {
				point := orb.Point{pt.X, pt.Y}
				lines = append(lines, point)
			}
			feature = geojson.NewFeature(lines)
		case *shp.Polygon:
			polygon, _ := s.(*shp.Polygon)
			//coordinates := make([][][]float64, polygon.NumParts)
			var i int32
			var polygons orb.Polygon
			for i = 0; i < polygon.NumParts; i = i + 1 {
				var startIndex, endIndex int32
				startIndex = polygon.Parts[i]
				if i == polygon.NumParts-1 {
					endIndex = int32(len(polygon.Points))
				} else {
					endIndex = polygon.Parts[i+1]
				}

				var ring orb.Ring
				for j := startIndex; j < endIndex; j = j + 1 {
					ring = append(ring, orb.Point{polygon.Points[j].X, polygon.Points[j].Y})
				}
				polygons = append(polygons, ring)
			}
			feature = geojson.NewFeature(polygons)
		default:
			fmt.Println("Not support geometry type", reflect.TypeOf(s).Elem())
			continue
		}

		for k, f := range fields {
			val := shape.ReadAttribute(n, k)
			feature.Properties[f.String()] = val
		}
		fc.Append(feature)

	}

	rawJSON, err := fc.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return rawJSON, nil
}

// Decompress
func Decompress(zipFile, destDir string) (error, []string) {
	zipReader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err, nil
	}
	defer zipReader.Close()

	var decodeName string
	var fileNames []string
	for _, innerFile := range zipReader.File {
		if innerFile.Flags == 0 {
			//如果标致位是0  则是默认的本地编码   默认为gbk
			i := bytes.NewReader([]byte(innerFile.Name))
			decoder := transform.NewReader(i, simplifiedchinese.GB18030.NewDecoder())
			content, _ := io.ReadAll(decoder)
			decodeName = string(content)
		} else {
			//如果标志为是 1 << 11也就是 2048  则是utf-8编码
			decodeName = innerFile.Name
		}
		srcFile, err := innerFile.Open()
		if err != nil {
			return err, nil
		}

		defer srcFile.Close()

		if innerFile.FileInfo().IsDir() {
			continue
		}
		name := GetFileName(decodeName)
		fp := filepath.Join(destDir, name)
		newFile, err0 := os.OpenFile(fp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, innerFile.Mode())
		if err0 != nil {
			return err0, nil
		}
		defer newFile.Close()

		_, err = io.Copy(newFile, srcFile)
		if err != nil {
			return err, nil
		}
		fileNames = append(fileNames, name)
		defer newFile.Close()
	}
	return nil, fileNames
}

func GetDir(path string) string {
	path = filepath.Clean(path)
	dir, _ := filepath.Split(path)
	return filepath.Base(dir)
}

func GetFileName(path string) string {
	path = filepath.Clean(path)
	_, filename := filepath.Split(path)
	return filename
}

func SearchFile(fp string) (list []string, err error) {
	fp, err = filepath.Abs(fp)
	files, err := os.ReadDir(fp)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		list = append(list, filepath.Join(fp, file.Name()))
	}
	return list, nil
}

// VerifyShapeFile 验证shapefile 是否齐全（必备：shp,shx,dbf,prj）
func VerifyShapeFile(fp string) bool {
	files, err := ioutil.ReadDir(GetDir(fp))
	name := strings.Split(GetFileName(fp), ".")[0]
	shx := fmt.Sprintf("%s.shx", name)
	dbf := fmt.Sprintf("%s.dbf", name)
	prj := fmt.Sprintf("%s.prj", name)
	if err != nil {
		return false
	}
	flag := 0
	for _, file := range files {
		switch file.Name() {
		case shx:
			flag += 1
		case dbf:
			flag += 1
		case prj:
			flag += 1
		default:
			flag += 0
		}
	}
	if flag == 3 {
		return true
	} else {
		return false
	}
}

func MoveFile(fp string, newLoc string) error {
	err := os.Rename(fp, newLoc)
	if err != nil {
		return err
	}
	return nil
}

func MoveShapeFile(fp string, newLoc string) {
	shx := strings.ReplaceAll(fp, ".shp", ".shx")
	dbf := strings.ReplaceAll(fp, ".shp", ".dbf")
	prj := strings.ReplaceAll(fp, ".shp", ".prj")
	err := MoveFile(fp, newLoc)
	if err != nil {
		return
	}
	err = MoveFile(shx, newLoc)
	if err != nil {
		return
	}
	err = MoveFile(dbf, newLoc)
	if err != nil {
		return
	}
	err = MoveFile(prj, newLoc)
	if err != nil {
		return
	}

}
