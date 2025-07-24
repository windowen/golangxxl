package passwordhelper

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"queueJob/pkg/zlogger"
)

// CBC能掩盖明文结构信息，保证相同密文可得不同明文，所以不容易主动攻击，安全性好于ECB，适合传输长度长的报文，是SSL和IPSec的标准。
// =================== CBC ======================
func AesEncryptCBC(origData []byte, key []byte) (encrypted []byte) {
	// 分组秘钥
	// NewCipher该函数限制了输入k的长度必须为16, 24或者32
	block, _ := aes.NewCipher(key)
	blockSize := block.BlockSize()                              // 获取秘钥块的长度
	origData = pkcs5Padding(origData, blockSize)                // 补全码
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize]) // 加密模式
	encrypted = make([]byte, len(origData))                     // 创建数组
	blockMode.CryptBlocks(encrypted, origData)                  // 加密
	return encrypted
}
func AesDecryptCBC(encrypted []byte, key []byte) (decrypted []byte, err error) {
	block, err := aes.NewCipher(key) // 分组秘钥
	if err != nil {
		zlogger.Errorf("AesDecryptCBC err :%v", err)
		return nil, err
	}
	blockSize := block.BlockSize()                              // 获取秘钥块的长度
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize]) // 加密模式
	decrypted = make([]byte, len(encrypted))                    // 创建数组

	if len(encrypted)%blockSize != 0 {
		zlogger.Errorf(fmt.Sprintf("encrypted:%v,blockSize:%v", len(encrypted), blockSize))
		return nil, err
	}
	blockMode.CryptBlocks(decrypted, encrypted) // 解密
	decrypted = pkcs5UnPadding(decrypted)       // 去除补全码
	return decrypted, nil
}
func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// EncryptByAes Aes加密 后 base64 再加 urlencoded
func EncryptByAesCBC(data []byte, key string) (string, error) {
	res := AesEncryptCBC(data, []byte(key))
	return url.QueryEscape(base64.StdEncoding.EncodeToString(res)), nil
}

func EncryptByAesCBC2(data []byte, key string) (string, error) {
	res := AesEncryptCBC(data, []byte(key))
	return url.QueryEscape(string(res)), nil
}

func DecryptByAesCBC(data string, key string) ([]byte, error) {
	data, _ = url.QueryUnescape(data)
	dataByte, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	rs, err := AesDecryptCBC(dataByte, []byte(key))
	return rs, nil
}

// =================== ECB ======================
func AesEncryptECB(origData []byte, key []byte) (encrypted []byte) {
	cipher, _ := aes.NewCipher(generateKey(key))
	length := (len(origData) + aes.BlockSize) / aes.BlockSize
	plain := make([]byte, length*aes.BlockSize)
	copy(plain, origData)
	pad := byte(len(plain) - len(origData))
	for i := len(origData); i < len(plain); i++ {
		plain[i] = pad
	}
	encrypted = make([]byte, len(plain))
	// 分组分块加密
	for bs, be := 0, cipher.BlockSize(); bs <= len(origData); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Encrypt(encrypted[bs:be], plain[bs:be])
	}

	return encrypted
}
func AesDecryptECB(encrypted []byte, key []byte) (decrypted []byte) {
	cipher, _ := aes.NewCipher(generateKey(key))
	decrypted = make([]byte, len(encrypted))
	//
	for bs, be := 0, cipher.BlockSize(); bs < len(encrypted); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Decrypt(decrypted[bs:be], encrypted[bs:be])
	}

	trim := 0
	if len(decrypted) > 0 {
		trim = len(decrypted) - int(decrypted[len(decrypted)-1])
	}

	return decrypted[:trim]
}
func generateKey(key []byte) (genKey []byte) {
	genKey = make([]byte, 16)
	copy(genKey, key)
	for i := 16; i < len(key); {
		for j := 0; j < 16 && i < len(key); j, i = j+1, i+1 {
			genKey[j] ^= key[i]
		}
	}
	return genKey
}

// MD5加密
func Md5Encrypt(str string) string {
	md := md5.New()
	md.Write([]byte(str))                // 需要加密的字符串为 str
	cipherStr := md.Sum(nil)             // 不需要拼接额外的数据，如果约定了额外加密数据，可以在这里传递
	return hex.EncodeToString(cipherStr) // 输出加密结果
}

// 管理员密码生成
func CreatManagerPwd(name, pwd string) string {
	var salt = "woshiyigesalt123456*&^%$#"
	var res = Md5Encrypt(name + pwd + salt)
	return res
}

/*
 加密规则：
  所有参数按照ASCII码升序排列 最后拼接密钥 例如：name=oscar&pwd=123&md5
  然后md5加密 转大写（加密后的串约定统一大写）
  时间戳采用格林威治时间 毫秒
*/
// md5验签
func CheckMD5(m map[string]interface{}, md5Key string) bool {
	strs := SortingMD5(m)
	strs = strs + md5Key
	md5str := Md5Encrypt(strs)
	// data := []byte(strs)
	// has := md5.Sum(data)
	// md5str :=fmt.Sprintf("%x",has)
	// 验签  sign统一大写
	return m["Sign"] == strings.ToUpper(md5str)
}

// 将参数按ASCII码升序排列 name=oscar&pwd=123&md5
func SortingMD5(m map[string]interface{}) string {
	var strs []string
	var ss string
	for k := range m {
		// 去掉sign拼接
		if k != "Sign" {
			strs = append(strs, fmt.Sprint(k))
		}
	}
	sort.Strings(strs)
	for _, v := range strs {
		ss = ss + v + "=" + fmt.Sprint(m[v]) + "&"
	}
	return ss
}
