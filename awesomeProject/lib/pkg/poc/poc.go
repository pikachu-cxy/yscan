package poc

import (
	"io/ioutil"
)

const pocPath = "C:\\Users\\simple.chen\\Downloads\\yscan5\\awesomeProject\\pocs"

func Scan(targetURL string, templatePath string) {
	// 加载 poc 模板

}

// 加载模板
func loadTemplates() {

	_, err := ioutil.ReadDir(pocPath)
	if err != nil {
		println(err.Error())
	}

	// 获取文件夹内所有以 .yaml 结尾的文件
	// 指定后缀名
	/*
		extension := ".yaml"
		//filePattern := filepath.Join(pocPath, "*.yaml")
		err = filepath.Walk(pocPath, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				fmt.Printf("访问文件或文件夹 %s 时出错：%v\n", path, err)
				return nil
			}
			if !info.IsDir() && filepath.Ext(path) == extension {
				fmt.Println(path)
			}
			return nil
		})

		if err != nil {
			fmt.Println("Failed to read files:", err)
			return
		}

	*/

	// 遍历每个文件进行读取和解析
	/*
		for _, file := range files {
			// 读取文件内容
			data, err := ioutil.ReadFile(file)
			if err != nil {
				fmt.Printf("Failed to read file %s: %s\n", file, err)
				continue
			}
			println(data)
			// 解析 YAML 数据
			var config Poc
			err = yaml.Unmarshal(data, &config)
			if err != nil {
				fmt.Printf("Failed to parse YAML file %s: %s\n", file, err)
				continue
			} else {
				println(err.Error())
			}
			fmt.Println("id: ", config.Id)

		}

	*/

}

// 根据keyword 过滤 返回pocs
func filterPoc() {

}
