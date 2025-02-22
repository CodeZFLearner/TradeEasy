/*
 * @Author: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @Date: 2025-01-19 21:51:39
 * @LastEditors: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @LastEditTime: 2025-02-14 00:48:19
 * @FilePath: \HelloGolang\internal\oss\file.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package oss

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type FileHelper struct{}

// WriteToFile writes the given data to a JSON file.
func (fh *FileHelper) WriteToFile(filename string, data interface{}) error {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()

	// 创建 JSON 编码器
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // 设置缩进格式

	// 将 map 编码为 JSON 并写入文件
	if err := encoder.Encode(data); err != nil {
		fmt.Println(err)
		return fmt.Errorf("could not encode data to JSON: %v", err)
	}
	return nil
}

func (fh *FileHelper) WriteTxt(filename string, data string) error {
	byteData := []byte(data)

	// 使用 ioutil.WriteFile 将数据写入文件
	// 如果文件不存在，会创建文件；如果文件已存在，会覆盖文件内容
	err := os.WriteFile(filename, byteData, 0644)
	if err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}
	return nil
}
func (fh *FileHelper) ReadTxt(filename string) (string, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("读取文件失败: %v", err)
	}
	return string(bytes), nil
}

// ReadFromFile reads the JSON file and loads the data into the provided interface.
func (fh *FileHelper) ReadFromFile(filename string, data interface{}) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, data)
}
