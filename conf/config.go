package conf

import (
	"gopkg.in/ini.v1"
	"log"
)

var iniFile *ini.File

var BufLen int
var KeySize int

func init() {
	BufLen = 40960
	KeySize = 256
}

func SetIni(name string) {
	file, err := ini.Load(name)
	if err != nil {
		log.Printf("Fail to read file: %v", err)
		return
	}
	iniFile = file
}

func GetIniValue(section string, key string) string {
	return iniFile.Section(section).Key(key).String()
}
