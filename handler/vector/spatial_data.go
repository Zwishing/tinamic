package vector

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"io/ioutil"
	"strings"
)

type SupportType string

const (
	geojson SupportType = ".geojson"
	zipShp  SupportType = ".zip"
	shp     SupportType = ".shp"
)

const (
	//上传数据的存储路径
	uploadPath = "./uplaod"

	//数据存储的目录
	storePath = "./data"

	//上传的字段名称
	uplaodField = "file"
)

// Upload 上传数据
//func Upload(ctx *fiber.Ctx) error {
//	file, err := ctx.FormFile(uplaodField)
//	if err != nil {
//		return err
//	}
//	if !IsSupport(file.Filename, string(geojson), string(zipShp)) {
//		err = response.Fail(ctx, "", "数据格式不支持！")
//		if err != nil {
//			return err
//		}
//	}
//	fp := filepath.Join(uploadPath, file.Filename)
//	err = ctx.SaveFile(file, fp)
//	if err != nil {
//		return err
//	}
//	err = RecordFile(fp)
//	if err != nil {
//		return err
//	}
//	return nil
//}

//func save2database(fp string) {
//	fileList,err:= utils.SearchFile(fp)
//	if err!=nil{
//		return
//	}
//	for _,file:=range fileList {
//		if filepath.Ext(file)==".shp"{
//			shapeFile, err := geos.NewShapeFile(file)
//			if err != nil {
//				return
//			}
//			err = shapeFile.MoveFile(storePath)
//			if err != nil {
//				return
//			}
//			uid, _ :=uuid.NewV4()
//			data:=models.SpatialData{
//				Uid:       uid,
//				Name:      shapeFile.Name,
//				IsPublish: false,
//				FileType:  "shapefile",
//				Size:      shapeFile.TotalSize(),
//				FilePath:  shapeFile.ShpPath[0],
//				CreateAt: shapeFile.CreatedAt,
//				UpdateAt: time.Now(),
//			}
//			_, err = queries.InsertSpatialData(data)
//			if err != nil {
//				return
//			}
//		}else if filepath.Ext(file)==".geojson"{
//			_,name:=filepath.Split(file)
//			fileInfo,err:=os.Stat(file)
//			newPath:=filepath.Join(storePath,name)
//			err=os.Rename(file,newPath)
//			if err!=nil {
//				return
//			}
//			uid, _ :=uuid.NewV4()
//			data:=models.SpatialData{
//				Uid:       uid,
//				Name:      strings.Split(name,".")[0],
//				IsPublish: false,
//				FileType:  "geojson",
//				Size:      fileInfo.Size(),
//				FilePath:  newPath,
//				CreateAt: fileInfo.ModTime(),
//				UpdateAt: time.Now(),
//			}
//			_, err = queries.InsertSpatialData(data)
//			if err != nil {
//				return
//			}
//		}
//	}
//}

func Publish(ctx *fiber.Ctx) error {
	type pub struct {
		Uid string `json:"uid" xml:"uid" form:"uid"`
	}
	a := string(ctx.Body())
	fmt.Println(a)
	p := new(pub)
	//解析路径下数据，
	ctx.BodyParser(p)
	// filepath<-uid
	//调用入库接口
	//fp:="D:\\Code\\go-web\\tinamic\\uplaod\\river.shp"
	//queries.Shp2PgSql("layers",p.Uid,fp)

	//读取数据获取基本信息，投影，范围，

	//写入数据库基本信息

	return nil
}

//func QuerySpatialData(ctx *fiber.Ctx) error {
//	data, err := model.QuerySpatialData()
//	if err != nil {
//		return err
//	}
//	err = response.Success(ctx, &data, "数据查询成功！")
//	if err != nil {
//		err = response.Fail(ctx, "", "数据查询失败！")
//		return err
//	}
//	return nil
//}

func Delete() {

}

// IsSupport 判断是否支持特定格式的文件
func IsSupport(file string, suffixes ...string) bool {
	for _, suffix := range suffixes {
		if strings.HasSuffix(file, suffix) {
			return true
		}
	}
	return false
}

//func RecordFile(fp string) error {
//	sd := model.NewSpatialData(fp)
//	if strings.HasSuffix(fp, string(geojson)) {
//		sd.FileType = string(geojson)
//	} else if strings.HasSuffix(fp, string(zipShp)) {
//		err, names := util.Decompress(fp, uploadPath)
//		if err != nil {
//			return err
//		}
//		// 删除压缩包
//		err = os.Remove(fp)
//		if err != nil {
//			return err
//		}
//		for _, name := range names {
//			if strings.HasSuffix(name, string(shp)) {
//				sd.Name = strings.Split(name, ".")[0]
//				sd.FileType = string(shp)
//				break
//			}
//		}
//		if sd.FileType == "" {
//			return errors.New("no shapefile in zip")
//		}
//	}
//	_, err := model.InsertSpatialData(sd)
//	if err != nil {
//		return err
//	}
//	return nil
//}

func SearchFile(fp string) {
	ioutil.ReadDir(fp)
}
