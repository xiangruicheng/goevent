package store

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"goevent/config"
	"io/ioutil"
	"os"
	"strings"
)

type FileStore struct {
	objType ObjType
	topic   string
	offset  int32
}

// Topic set topic
func (s *FileStore) Topic(topic string) Store {
	s.topic = topic
	return s
}

// ObjType set objType
func (s *FileStore) ObjType(objType ObjType) Store {
	s.objType = objType
	return s
}

// Offset set offset
func (s *FileStore) Offset(offset int32) Store {
	s.offset = offset
	return s
}

// Add append write file
func (s *FileStore) Add(data any) error {
	switch s.objType {
	case ObjTypeData:
		DataFileMutex.Lock()
		defer DataFileMutex.Unlock()
		break
	case ObjTypeSubscribe:
		SubscribeFileMutex.Lock()
		defer SubscribeFileMutex.Unlock()
		break
	case ObjTypeEvent:
		EventFileMutex.Lock()
		defer EventFileMutex.Unlock()
		break
	}

	filename := s.getFile()
	if !s.isFile(filename) {
		err := s.newFile(filename)
		if err != nil {
			return err
		}
	}

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	defer file.Close()
	if err != nil {
		return err
	}
	content, _ := json.Marshal(data)
	_, err = file.WriteString(string(content) + "\n")
	if err != nil {
		return err
	}
	return nil
}

// GetList read file from offset read length
func (s *FileStore) GetList(offest int32, length int32) ([]string, error) {
	filename := s.getFile()
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 创建 Scanner 对象
	scanner := bufio.NewScanner(file)
	rows := []string{}
	var num int32 = 0
	var i int32 = 1
	for i = 1; scanner.Scan(); i++ {
		if i >= offest && num < length {
			data := scanner.Text() // 获取当前行的文本
			rows = append(rows, data)
			num += 1
		}
	}
	return rows, nil
}

// Get read a row
func (s *FileStore) Get(offest int32) (string, error) {
	filename := s.getFile()
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 创建 Scanner 对象
	scanner := bufio.NewScanner(file)

	var i int32 = 1
	for i = 1; scanner.Scan(); i++ {
		if i == offest {
			data := scanner.Text() // 获取当前行的文本
			return data, nil
		}
	}
	return "", nil
}

// All read whole
func (s *FileStore) All() ([]string, error) {
	filename := s.getFile()
	contentTypes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	contentStr := strings.TrimRight(string(contentTypes), "\n")
	// 将文件内容按行分割
	lines := strings.Split(contentStr, "\n")
	return lines, nil
}

// Update update obj from offset
func (s *FileStore) Update(text any) error {
	switch s.objType {
	case ObjTypeData:
		return errors.New("data type not support update")
	case ObjTypeSubscribe:
		SubscribeFileMutex.Lock()
		defer SubscribeFileMutex.Unlock()
		break
	case ObjTypeEvent:
		EventFileMutex.Lock()
		defer EventFileMutex.Unlock()
		break
	}

	filename := s.getFile()
	line := s.offset

	// 读取文件的全部内容
	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	// 将文件内容按行分割
	lines := strings.Split(string(fileContent), "\n")
	// 检查行数是否足够
	if line < 1 || line > int32(len(lines)) {
		return errors.New("指定的行号超出范围" + filename)
	}

	// 更新指定行的内容
	newTextByte, err := json.Marshal(text)
	if err != nil {
		return err
	}
	lines[line-1] = string(newTextByte)
	// 将修改后的行重新组合成文件内容
	newContent := strings.Join(lines, "\n") // 确保文件最后有一个换行符
	err = os.WriteFile(filename, []byte(newContent), 0644)
	if err != nil {
		return err
	}
	return nil
}

// getFile get current obj filename
func (s *FileStore) getFile() string {
	switch s.objType {
	case ObjTypeEvent:
		filename := config.Data.Event.DataDir + "event.data"
		return filename
	case ObjTypeSubscribe:
		filename := config.Data.Event.DataDir + s.topic + "/subscribe.data"
		return filename
	case ObjTypeData:
		filenameIndex := (s.offset / config.Data.Event.DataFileMaxLines) * config.Data.Event.DataFileMaxLines
		filename := fmt.Sprintf("%016d.data", filenameIndex)
		filename = config.Data.Event.DataDir + s.topic + "/" + filename
		return filename
	}
	return ""
}

// newFile create a file and dir
func (s *FileStore) newFile(fullFilename string) error {
	//得到文件全路径
	pathArr := strings.Split(fullFilename, "/")
	//文件名
	filename := pathArr[len(pathArr)-1]
	//目录
	dir := strings.Replace(fullFilename, filename, "", 1)
	//创建目录
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}
	//创建文件
	file, err := os.Create(fullFilename)
	defer file.Close()
	return err
}

// isFile is a file
func (s *FileStore) isFile(filename string) bool {
	_, err := os.Stat(filename)
	if err == nil {
		return true
	}
	return false
}
