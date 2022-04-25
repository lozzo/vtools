package tools

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func KillOcr() {
	c := exec.Command("taskkill.exe", "/f", "/im", "PaddleOCRServer.exe")
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err := c.Run(); err != nil {
		fmt.Println("killOcr Error: ", err)
	}
}

func RunOcr() {
	KillOcr()
	c := exec.Command("cmd", "/C", "PaddleOCRServer.exe")
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err := c.Run(); err != nil {
		fmt.Println("runOcr Error: ", err)
	}
	os.Exit(1)
}
func Ocr(lang string, image []byte) (string, error) {
	image_base64 := base64.StdEncoding.EncodeToString(image)
	req, _ := http.NewRequest("POST", "http://localhost:19941/ocr", strings.NewReader(image_base64))
	req.Header.Set("Content-Type", "text/plain")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Print(string(body))
	body_str := string(body)
	return body_str, nil
}
