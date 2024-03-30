package controllers

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

var ip string

func init(){
	readConfig(ip)
}

// Mvt 返回矢量瓦片
func Mvt(mvt []byte, x int, y int, z int, params string)  {
	res, _ := http.Get(fmt.Sprintf("http://localhost:7800/dwh_gis.g_administrative_boundary/%d/%d/%d.pbf?)",z,x,y))
	defer res.Body.Close() // 在回复后必须关闭回复的主体
	mvt, _ = ioutil.ReadAll(res.Body)
	//return body
}




func readConfig(ip string){

}