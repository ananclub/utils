package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/goinggo/mapstructure"
	"github.com/google/uuid"
	"github.com/json-iterator/go"
	"github.com/orcaman/concurrent-map/v2"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

var Json = jsoniter.ConfigCompatibleWithStandardLibrary
var json = Json

func NewCMap[v any]() cmap.ConcurrentMap[v] {
	return cmap.New[v]()
}
func MapToStruct(v interface{}) (err error) {
	mapInstance := make(map[string]interface{})
	err = mapstructure.Decode(mapInstance, &v)
	return
}
func StructToMap(obj interface{}) map[string]interface{} {
	obj1 := reflect.TypeOf(obj)
	obj2 := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < obj1.NumField(); i++ {
		data[obj1.Field(i).Name] = obj2.Field(i).Interface()
	}
	return data
}
func BytesToUint16(b []byte) uint16 {
	buf := bytes.NewBuffer(b)
	var tmp uint16
	binary.Read(buf, binary.BigEndian, &tmp)
	return tmp
}
func BytesToUint32(b []byte) uint32 {
	buf := bytes.NewBuffer(b)
	var tmp uint32
	binary.Read(buf, binary.BigEndian, &tmp)
	return tmp
}
func BytesToUint64(b []byte) uint64 {
	buf := bytes.NewBuffer(b)
	var tmp uint64
	binary.Read(buf, binary.BigEndian, &tmp)
	return tmp
}
func Uint16ToBytes(n uint16) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(n))
	return b
}
func Uint32ToBytes(n uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(n))
	return b
}
func Uint64ToBytes(n uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, n)
	return b
}
func ArrayRemove(s []string, i int) []string {
	return append(s[:i], s[i+1:]...)
}
func GetCurrentTimeStr() (s string) {
	zone := "Asia/Shanghai"
	format := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation(zone)
	return time.Now().In(loc).Format(format)
}
func GetCurrentTimeStrWithZone(format, zone string) (s string, err error) {
	if zone == "" {
		zone = "Asia/Shanghai"
	}
	if format == "" {
		format = "2006-01-02 15:04:05"
	}
	loc, err := time.LoadLocation(zone)
	return time.Now().In(loc).Format(format), err
}
func GetCurrentTime(zone string) time.Time {
	if zone == "" {
		zone = "Asia/Shanghai"
	}
	loc, _ := time.LoadLocation(zone)
	return time.Now().In(loc)
}
func Truncate(s string, n int) string {
	runes := []rune(s)
	if len(runes) > n {
		return string(runes[:n])
	}
	return s
}

func UUID(hasSlice bool) string {
	if hasSlice {
		return uuid.New().String()

	} else {
		return strings.ReplaceAll(uuid.New().String(), "-", "")
	}
}

func AesDecrypt(ciphertext []byte, keystring string) ([]byte, error) {
	// Key
	key := []byte(keystring)

	// Create the AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Before even testing the decryption,
	// if the text is too small, then it is incorrect
	if len(ciphertext) < aes.BlockSize {
		err = errors.New("Text is too short")
		return nil, nil
	}

	// Get the 16 byte IV
	iv := ciphertext[:aes.BlockSize]

	// Remove the IV from the ciphertext
	ciphertext = ciphertext[aes.BlockSize:]

	// Return a decrypted stream
	stream := cipher.NewCFBDecrypter(block, iv)

	// Decrypt bytes from ciphertext
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, nil
}

func AesEncrypt(plaintext []byte, keystring string) ([]byte, error) {

	// Key
	key := []byte(keystring)

	// Create the AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Empty array of 16 + plaintext length
	// Include the IV at the beginning
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))

	// Slice of first 16 bytes
	iv := ciphertext[:aes.BlockSize]

	// Write 16 rand bytes to fill iv
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	// Return an encrypted stream
	stream := cipher.NewCFBEncrypter(block, iv)

	// Encrypt bytes from plaintext to ciphertext
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext, nil
}

func SplitDockerUrl(url string) (host string, port int, repo string, version string, isDocker bool, err error) {

	strs := strings.Split(url, "/")
	switch len(strs) {
	case 1:
		isDocker = true
		rs := strings.Split(url, ":")
		switch len(rs) {
		case 1:
			repo = rs[0]
			version = "latest"
		case 2:
			repo = rs[0]
			version = rs[1]
		}
		return
	case 2:
		if !strings.Contains(strs[0], ".") {
			isDocker = true
			host = strs[0]

		} else {
			hs := strings.Split(strs[0], ":")
			switch len(hs) {
			case 1:
				host = hs[0]
			case 2:
				host = hs[0]
				port, err = strconv.Atoi(hs[1])
				if err != nil {
					return
				}
			default:
				err = errors.New("not valid host:" + strs[0])
				return
			}
		}
		rs := strings.Split(strs[1], ":")
		switch len(rs) {
		case 1:
			repo = rs[0]
			version = "latest"
		case 2:
			repo = rs[0]
			version = rs[1]
		}
		return
	default:
		err = errors.New("not valid url:" + url)
	}
	return
}

func CheckDockerRepo(host string, port int, repo string, v string, authcode string) (tags []interface{}, err error) {

	url := "https://" + host + ":" + strconv.Itoa(port) + "/v2/" + repo + "/tags/list"
	//master1.meleclass.com registry authcode

	headers := make(map[string]string)
	headers["Authorization"] = "Basic " + authcode

	_, body, err := HttpDo(url, "GET", nil, nil, headers, nil, nil)
	if err != nil {
		return
	}

	var s map[string]interface{}
	if err = json.Unmarshal(body, &s); err != nil {
		return
	}
	var ok bool
	if tags, ok = s["tags"].([]interface{}); ok {

		for _, d := range tags {
			if ds, ok := d.(string); ok && ds == v {
				return
			}
		}

		err = errors.New("no this version:" + v)
		return

	} else {
		err = errors.New("json tags decode err:" + string(body))
		return
	}

	return
}
func DockerImageChangeVer(image string, newVer string) (newImage string, err error) {
	if image == "" {
		err = errors.New("image param can not be empty string")
		return
	}
	strs := strings.Split(image, "/")

	switch len(strs) {
	case 1:
		rs := strings.Split(strs[0], ":")

		newImage = rs[0]
		if newVer != "" {
			newImage += ":" + newVer
		}

		return
	case 2:

		newImage = strs[0]
		rs := strings.Split(strs[1], ":")
		newImage += "/" + rs[0]
		if newVer != "" {
			newImage += ":" + newVer
		}
		return
	default:
		err = errors.New("not valid image url:" + image)
	}
	return
}

func PathExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func GBKToUTF8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func UTF8ToGBK(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

var UrlEncode func(string) string = url.QueryEscape
var UrlDecode func(string) (string, error) = url.QueryUnescape

func CheckParamStrict(params map[string]interface{}, key string, in interface{}, canNull bool) (isnull bool, err error) {
	param, ok := params[key]
	if !ok {
		isnull = true
		if !canNull {
			err = fmt.Errorf("param %s is null", key)
		}
		return
	}
	tp := reflect.TypeOf(param)
	p := reflect.ValueOf(in)
	if p.Kind() != reflect.Ptr {
		err = fmt.Errorf("dest param must be ptr")
		return
	}

	v := p.Elem()
	t := v.Type()

	if tp != t {
		err = fmt.Errorf("want type %s but %T", t.Name(), param)
		return
	}
	*v.Addr().Interface().(*interface{}) = param
	return
}

func CheckParam(params map[string]interface{}, key string, in interface{}, canNull bool) (isNull bool, err error) {
	param, ok := params[key]
	if !ok {
		isNull = true
		if !canNull {
			err = fmt.Errorf("param %s is null", key)
		}
		return
	}
	tp := reflect.TypeOf(param)
	ptr := reflect.ValueOf(in)
	if ptr.Kind() != reflect.Ptr {
		err = fmt.Errorf("dest param must be ptr type,please contact coder to fix it")
		return
	}
	v := ptr.Elem()
	t := v.Type()

	switch tp.Kind() {

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch t.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v.SetInt(param.(int64))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			v.SetUint(uint64(param.(int64)))
		case reflect.Float32, reflect.Float64:
			v.SetFloat(float64(param.(int64)))
		case reflect.String:
			v.SetString(strconv.FormatInt(param.(int64), 10))
		}
		return
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch t.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v.SetInt(int64(param.(uint64)))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			v.SetUint(param.(uint64))
		case reflect.Float32, reflect.Float64:
			v.SetFloat(float64(param.(uint64)))
		case reflect.String:
			v.SetString(strconv.FormatUint(param.(uint64), 10))
		}
		return
	case reflect.Float32, reflect.Float64:
		switch t.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v.SetInt(int64(param.(float64)))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			v.SetUint(uint64(param.(float64)))
		case reflect.Float32, reflect.Float64:
			v.SetFloat(param.(float64))
		case reflect.String:
			v.SetString(strconv.FormatFloat(param.(float64), 'f', -1, 64))
		}
		return
	case reflect.String:
		switch t.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			var i int64
			i, err = strconv.ParseInt(param.(string), 10, 64)
			if err == nil {
				v.SetInt(i)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			var i uint64
			i, err = strconv.ParseUint(param.(string), 10, 64)
			if err == nil {
				v.SetUint(i)
			}
		case reflect.Float32, reflect.Float64:
			var f float64
			f, err = strconv.ParseFloat(param.(string), 64)
			if err == nil {
				v.SetFloat(f)
			}
		case reflect.String:
			v.SetString(param.(string))

		}
		return
	}
	err = fmt.Errorf("dest param want type %s but %T", t.Name(), param)
	return
}
func VerifyEmailFormat(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}
