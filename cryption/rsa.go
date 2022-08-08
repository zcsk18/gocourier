package cryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
)

type RSA struct {
	PublicKey *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
}

const keyLength = 2048
const chunkLength = 256
const chunkOffset = 11

func (this *RSA) Encode(plainText []byte) []byte {
	var cipherText []byte
	n := len(plainText)
	length:=chunkLength - chunkOffset
	for num:=0;num<n;num+=length {
		var p []byte
		if num+length > n {
			p = plainText[num:]
			length = n - num
		} else {
			p = plainText[num:num+length]
		}
		buf, err := rsa.EncryptPKCS1v15(rand.Reader, this.PublicKey, p[:length])
		if err!=nil{
			panic(err)
		}
		cipherText = append(cipherText, buf...)
	}
	return cipherText
}

func (this *RSA) Decode(cipherText []byte) []byte {
	var plainText []byte
	n := len(cipherText)
	length:=chunkLength

	for num:=0; num<n; num+=length {
		var p []byte
		if num+length > n {
			p = cipherText[num:]
			length = n - num
		} else {
			p = cipherText[num:num+length]
		}
		buf, err:=rsa.DecryptPKCS1v15(rand.Reader, this.PrivateKey, p)
		if err != nil {
			log.Printf("err msg[%d] %d\n", n, length)
			panic(err)
		}
		plainText = append(plainText, buf...)
	}
	return plainText
}

func (this *RSA) GenerateKey() {
	bits:= keyLength
	//GenerateKey函数使用随机数据生成器random生成一对具有指定字位数的RSA密钥
	//Reader是一个全局、共享的密码用强随机数生成器
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err!=nil{
		panic(err)
	}
	//保存私钥
	//通过x509标准将得到的ras私钥序列化为ASN.1 的 DER编码字符串
	X509PrivateKey := x509.MarshalPKCS1PrivateKey(privateKey)
	//使用pem格式对x509输出的内容进行编码
	//创建文件保存私钥
	privateFile, err := os.Create("private.pem")
	if err!=nil{
		panic(err)
	}
	defer privateFile.Close()
	//构建一个pem.Block结构体对象
	privateBlock:= pem.Block{Type: "RSA Private Key",Bytes:X509PrivateKey}
	//将数据保存到文件
	pem.Encode(privateFile,&privateBlock)
	//保存公钥
	//获取公钥的数据
	publicKey:=privateKey.PublicKey
	//X509对公钥编码
	X509PublicKey,err:=x509.MarshalPKIXPublicKey(&publicKey)
	if err!=nil{
		panic(err)
	}
	//pem格式编码
	//创建用于保存公钥的文件
	publicFile, err := os.Create("public.pem")
	if err!=nil{
		panic(err)
	}
	defer publicFile.Close()
	//创建一个pem.Block结构体对象
	publicBlock:= pem.Block{Type: "RSA Public Key",Bytes:X509PublicKey}
	//保存到文件
	pem.Encode(publicFile,&publicBlock)
}

func (this *RSA) SetPublicKey(path string) {
	//打开文件
	file,err:=os.Open(path)
	if err!=nil{
		panic(err)
	}
	defer file.Close()
	//读取文件的内容
	info, _ := file.Stat()
	buf:=make([]byte,info.Size())
	file.Read(buf)
	//pem解码
	block, _ := pem.Decode(buf)
	//x509解码
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err!=nil{
		panic(err)
	}
	//类型断言
	publicKey:=publicKeyInterface.(*rsa.PublicKey)
	this.PublicKey = publicKey
}

func (this *RSA) SetPrivateKey(path string) {
	file,err:=os.Open(path)
	if err!=nil{
		panic(err)
	}
	defer file.Close()
	//获取文件内容
	info, _ := file.Stat()
	buf:=make([]byte,info.Size())
	file.Read(buf)
	//pem解码
	block, _ := pem.Decode(buf)
	//X509解码
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err!=nil{
		panic(err)
	}
	this.PrivateKey = privateKey
}