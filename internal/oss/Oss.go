package oss

import (
	"compress/gzip"
	"context"
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"sync"
)

var (
	clientInstance *cos.Client
	optInstance    *cos.MultiDownloadOptions
	once           sync.Once
)

const (
	SecretId  = "AKIDOe5SzpjYM8EKZn9UzSA6viT2NVRErSII"
	SecretKey = "c6RZT8YrVPNvdu8WC9NoZHxZERgzBDNV"
	DownPath  = "./data/"
)

var fileCh = make(chan string, 3)
var wg sync.WaitGroup

func RunCommandBash(command string) (string, error) {
	//cmd := exec.Command("bash", "-c", command)
	cmd := exec.Command("cmd", "/c", command)

	output, err := cmd.CombinedOutput() // 捕获 stdout 和 stderr
	return string(output), err
}

func download(key string) {
	defer wg.Done()
	if clientInstance == nil {
		GetClient()
	}
	if optInstance == nil {
		GetDownloadOptions()
	}

	file := fmt.Sprintf("%s%s", DownPath, key)

	_, err := clientInstance.Object.Download(
		context.Background(), key, file, optInstance,
	)
	if err != nil {
		panic(err)
	}
	fileCh <- key
}

// GetClient ensures clientInstance is initialized only once.
func GetClient() *cos.Client {
	once.Do(func() {
		// Initialize the COS Client
		u, err := url.Parse("https://pankou-1311832543.cos.ap-beijing.myqcloud.com")
		if err != nil {
			panic("Failed to parse bucket URL: " + err.Error()) // Handle initialization error
		}
		b := &cos.BaseURL{BucketURL: u}
		clientInstance = cos.NewClient(b, &http.Client{
			Transport: &cos.AuthorizationTransport{
				SecretID:  SecretId,
				SecretKey: SecretKey,
			},
		})
	})
	return clientInstance
}

// GetDownloadOptions ensures optInstance is initialized only once.
func GetDownloadOptions() *cos.MultiDownloadOptions {
	once.Do(func() {
		optInstance = &cos.MultiDownloadOptions{
			ThreadPoolSize: 3,
		}
	})
	return optInstance
}

func gzExtract(file string) {
	file = fmt.Sprintf("%s%s", DownPath, file)
	toFile := file[:len(file)-3]
	gzFile, err := os.Open(file)
	if err != nil {
		fmt.Println("Failed to open gzip file:", err)
		return
	}
	defer gzFile.Close()

	gzipReader, err := gzip.NewReader(gzFile)
	if err != nil {
		fmt.Println("Failed to create gzip reader:", err)
		return
	}
	defer gzipReader.Close()

	sqlFile, err := os.Create(toFile)
	if err != nil {
		fmt.Println("Failed to create SQL file:", err)
		return
	}
	defer sqlFile.Close()

	_, err = io.Copy(sqlFile, gzipReader)
	if err != nil {
		fmt.Println("Failed to write SQL file:", err)
		return
	}

}

func sql(key string) {
	//cmd := fmt.Sprintf("gunzip -c %s | mysql -uroot -p%s dust", key, "zff666...")
	cmd := fmt.Sprintf("mysql -uroot -p%s dust < %s%s", "root", DownPath, key[:len(key)-3])
	fmt.Println(cmd)

	std, err := RunCommandBash(cmd)
	//
	if err != nil {
		panic(err)
	}
	fmt.Println("Cmd Output:", std)
}
