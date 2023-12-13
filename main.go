package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	chromeVersion, err := getInstalledChromeVersion()
	if err != nil {
		panic(err)
	}

	fmt.Println("Detected Chrome Version:", chromeVersion)

	chromeDriverVersion, err := getChromeDriverVersion(chromeVersion)
	if err != nil {
		panic(err)
	}

	fmt.Println("Compatible ChromeDriver Version:", chromeDriverVersion)

	downloadURL := chromeDriverVersion
	downloadChromeDriver(downloadURL, "driver/chromedriver.zip")
	unzip("driver/chromedriver.zip", "driver")
}

func unzip(src, dst string) (fp *os.File, err error) {
	zipFp, err := zip.OpenReader(src)
	if err != nil {
		return nil, err
	}
	fp, err = os.Create(dst)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	for _, zfp := range zipFp.File {
		zp, err := zfp.Open()
		if err != nil {
			return nil, err
		}
		zipByte, err := io.ReadAll(zp)
		zp.Close()
		if err != nil {
			return nil, nil
		}
		fp.Write(zipByte)
	}
	return nil, nil
}

func getInstalledChromeVersion() (string, error) {
	// Chromeのバージョンを取得するコマンドを実行
	out, err := exec.Command("reg", "query", "HKEY_CURRENT_USER\\Software\\Google\\Chrome\\BLBeacon", "/v", "version").Output()
	if err != nil {
		return "", err
	}

	output := string(out)
	start := strings.Index(output, "REG_SZ")
	if start == -1 {
		return "", fmt.Errorf("cannot find Chrome version in registry")
	}

	version := strings.TrimSpace(output[start+6:])
	return version, nil
}

func getChromeDriverVersion(chromeVersion string) (string, error) {
	// Chromeのメジャーバージョンに基づいてChromeDriverのバージョンを取得
	majorVersion := strings.Split(chromeVersion, ".")[0]
	majorVersionInt, err := strconv.Atoi(majorVersion)
	if err != nil {
		return "", err
	}
	var resp *http.Response
	if majorVersionInt <= 114 {
		resp, err = http.Get(fmt.Sprintf("https://chromedriver.storage.googleapis.com/LATEST_RELEASE_%s", majorVersion))
	} else {
		resp, err = http.Get(("https://googlechromelabs.github.io/chrome-for-testing/latest-versions-per-milestone-with-downloads.json"))
	}
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var version string
	if majorVersionInt > 114 {
		chromeDriverData := &ChromeDriverData{}
		if err := json.Unmarshal(respBytes, chromeDriverData); err != nil {
			return "", err
		}
		downloads := chromeDriverData.Milestones[majorVersion].Downloads
		for index, cd := range downloads.ChromeDriver {
			cp := cd.Platform
			if strings.HasPrefix(cp, "win") || strings.HasPrefix(cp, "mac") || strings.HasPrefix(cp, "linux") {
				fmt.Println(index, cd)
				return cd.URL, nil
			}
		}
	}

	_, err = fmt.Fscanf(resp.Body, "%s", &version)
	return version, err
}

func downloadChromeDriver(url string, filepath string) {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println("ChromeDriver downloaded successfully:", filepath)
}
