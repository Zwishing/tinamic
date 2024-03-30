package controllers

//func TestPublish(t *testing.T) {
//	//解析路径下数据，
//	//如果是压缩文件需要先解压，
//	err, files := utils.DeCompress("../../temp/spatialdata/spatialdata.zip","../../temp/spatialdata")
//	if err != nil {
//		return
//	}
//	fmt.Println(files)
//	//调用入库接口
//	for _,file:=range files{
//		if !strings.HasSuffix(file,".shp"){
//			continue
//		}
//		queries.Shp2PgSql(filepath.Join("../../temp/spatialdata",file),db)
//	}
//}