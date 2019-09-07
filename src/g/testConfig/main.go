package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"smartping/src/g"
	"strings"
)

var region = map[string]string {
	"杭州":"浙江",
	"合肥":"安徽",
	"福州":"厦门",
	"南昌":"江西",
	"济南":"山东",
	"郑州":"河南",
	"武汉":"湖北",
	"长沙":"湖南",
	"广州":"广东",
	"海口":"海南",
	"太原":"陕西",
	"西宁":"青海",
	"南京":"江苏",
	"沈阳":"辽宁",
	"长春":"吉林",
	"石家庄":"河北",
	"贵阳":"贵州",
	"成都":"四川",
	"昆明":"云南",
	"西安":"陕西",
	"兰州":"甘肃",
	"南宁":"广西",
	"银川":"宁夏",
	"乌鲁木齐":"新疆",
	"拉萨":"西藏",
	"哈尔滨":"黑龙江",
	"呼和浩特":"内蒙古",
	"北京":"北京",
	"上海":"上海",
	"重庆":"重庆",
	"天津":"天津",
}

var ispAcronymMap = map[string]string {
	"联通":"cucc",
	"移动":"cmcc",
	"电信":"ctcc",
}

func getRegion(city string) string {
	return region[city]
}

func getISPAcronym(isp string) string {
	return ispAcronymMap[isp]
}

func main() {

	ipmap := ReadIpList("/Users/zhangyuwei/GitLab/CNPing/iptoconfig.dat")
	for k, v := range ipmap {
		fmt.Println("--- 地点:", k)
		for kk, vv := range v {
			fmt.Println("运营商:", kk)
			for _, vvv := range vv {
				fmt.Println("\tip:", vvv)
			}
		}
	}

	Cfg := g.ReadConfig("/Users/zhangyuwei/GitLab/smartping/conf/config.json.test")

	//for k, v := range Cfg.Chinamap {
	//	for kk, _ := range v {
	//		Cfg.Chinamap[k][kk] = append(Cfg.Chinamap[k][kk], ipmap[k][kk]...)
	//	}
	//}

	for _, k := range region {
		fmt.Println("k:", k)
		for _, kk := range ispAcronymMap {
			fmt.Println("kk:", kk)

			//Cfg.Chinamap[k][kk] = append(Cfg.Chinamap[k][kk], ipmap[k][kk]...)
		}
	}

	g.AppendIPToConfigFile(Cfg, "", "/Users/zhangyuwei/GitLab/smartping/conf/config.json.auto.0906")

	fmt.Println("------")
	for k, v := range Cfg.Chinamap {
		fmt.Println("--- k:", k)
		for kk, vv := range v {
			fmt.Println("-- kk:", kk)
			for _, vvv := range vv {
				fmt.Println("vvv:", vvv)
			}
		}
	}

}

func ReadIpList(fileName string) map[string]map[string][]string {
	if !g.IsExist(fileName) {
		fmt.Println("ip list file dose not exist")
		return nil
	}
	// map[位置]map[运营商][]ip
	ipCouldPing := make(map[string]map[string][]string)
	// map[运营商][]ip
	ispMap := make(map[string][]string)
	file, err := os.Open(fileName)
	defer file.Close()
	if err != nil {
		log.Fatal("Config Not Found!")
	} else {
		inputReader := bufio.NewReader(file)
		for {
			inputTmp, readerError := inputReader.ReadString('\n')
			inputString := strings.TrimRight(inputTmp, "\n")
			sp := make([]string, 3)
			if strings.Contains(inputString, "begin") {
				sp = strings.Split(inputString, " ")	//sp[1] 位置;  sp[2] 运营商
				ipSlice := make([]string, 0)
				for {
					ipTmp, err := inputReader.ReadString('\n')
					ip := strings.TrimRight(ipTmp, "\n")
					if strings.Contains(ip, "end") {
						// ip 遍历结束
						ispMap[getISPAcronym(sp[2])] = ipSlice
						break
					}
					if err == io.EOF {
						// 文件结尾
						break
					}
					ipSlice = append(ipSlice, ip)
				}
			}

			if len(ispMap) == 3  {
				ipCouldPing[getRegion(sp[1])] = ispMap
				ispMap = make(map[string][]string)
			}

			//cleanCIDR := strings.TrimRight(inputString, "\n")
			//fmt.Printf("The input was: %s", inputString)
			if readerError == io.EOF {
				break
			}
		}
	}
	fmt.Println("ReadConfig success")
	return ipCouldPing
}