package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"github.com/astaxie/beego"
	"log"
)
var key = []byte("sadsadsadsadsads")
var iv = []byte("sadsadsadsadsads")

func Encrypt(text []byte,typ int, token string) (string,error) {
	//登录时没有token 采用默认方式
	if typ==1 {
		md:=md5.New()
		md.Write([]byte(token))
		ciphertext:=hex.EncodeToString(md.Sum(nil))
		sha1:=sha1.New()
		sha1.Write([]byte(ciphertext))
		ciphertext=hex.EncodeToString(sha1.Sum(nil))
		md=md5.New()
		md.Write([]byte(ciphertext))
		ciphertext=hex.EncodeToString(md.Sum(nil))
		key=[]byte(ciphertext[:16])
		iv=[]byte(ciphertext[16:])
	}else if typ==2 {
		key = []byte("648215a86188855a")
		iv = []byte("d14eb2db3fce328c")
	}else if typ==0 {
		key = []byte("sadsadsadsadsads")
		iv = []byte("sadsadsadsadsads")
	}
	//生成cipher.Block 数据块
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Println("错误 -" +err.Error())
		return "",err
	}
	//填充内容，如果不足16位字符
	blockSize := block.BlockSize()
	originData := pad(text,blockSize)
	//加密方式
	blockMode := cipher.NewCBCEncrypter(block,iv)
	//加密，输出到[]byte数组
	crypted := make([]byte,len(originData))
	blockMode.CryptBlocks(crypted,originData)
	return base64.StdEncoding.EncodeToString(crypted) , nil
}

func pad(ciphertext []byte, blockSize int) []byte{
	padding := blockSize - len(ciphertext) % blockSize
	padtext := bytes.Repeat([]byte{byte(padding)},padding)
	return append(ciphertext,padtext...)
}

func Decrypt(text string,typ int ,token string) (string,error){
	//登录时没有token 采用默认方式
	if typ==1 {
		md:=md5.New()
		md.Write([]byte(token))
		ciphertext:=hex.EncodeToString(md.Sum(nil))
		sha1:=sha1.New()
		sha1.Write([]byte(ciphertext))
		ciphertext=hex.EncodeToString(sha1.Sum(nil))
		md=md5.New()
		md.Write([]byte(ciphertext))
		ciphertext=hex.EncodeToString(md.Sum(nil))
		key=[]byte(ciphertext[:16])
		iv=[]byte(ciphertext[16:])
	}else if typ==2 {
		key = []byte("648215a86188855a")
		iv = []byte("d14eb2db3fce328c")
	}else if typ==0 {
		key = []byte("sadsadsadsadsads")
		iv = []byte("sadsadsadsadsads")
	}
	decodeData,err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return "",nil
	}
	beego.Info(string(key),string(iv),typ)
	//生成密码数据块cipher.Block
	block,_ := aes.NewCipher(key)
	//解密模式
	blockMode := cipher.NewCBCDecrypter(block,iv)
	//输出到[]byte数组
	originData := make([]byte,len(decodeData))
	blockMode.CryptBlocks(originData, decodeData)
	//去除填充,并返回
	return string(unpad(originData)),nil
}

func unpad(ciphertext []byte) []byte{
	length := len(ciphertext)
	//去掉最后一次的padding
	unpadding := int(ciphertext[length - 1])
	return ciphertext[:(length - unpadding)]
}
