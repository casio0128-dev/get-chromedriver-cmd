package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
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
	downloadChromeDriver(downloadURL, "driver"+string(os.PathSeparator)+"chromedriver.zip")
	if err := unzip("driver"+string(os.PathSeparator)+"chromedriver.zip", "driver"); err != nil {
		log.Fatalln(err)
	}
	const driverPATH = "driver" + string(os.PathSeparator) + "chromedriver"
	filepath.WalkDir(driverPATH, func(path string, d fs.DirEntry, err error) error {
		fmt.Println(path)
		dst, err := os.Create(path)
		if err != nil {
			return err
		}
		src, err := os.Open(driverPATH)
		if err != nil {
			return err
		}
		io.Copy(dst, src)

		return nil
	})

	os.Remove("driver/chromedriver.zip")
}

// unZip zipファイルを展開する
func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	ext := filepath.Ext(src)
	rep := regexp.MustCompile(ext + "$")
	dir := filepath.Base(rep.ReplaceAllString(src, ""))

	destDir := filepath.Join(dest, dir)
	// ファイル名のディレクトリを作成する
	if err := os.MkdirAll(destDir, os.ModeDir); err != nil {
		return err
	}

	for _, f := range r.File {
		if f.Mode().IsDir() {
			// ディレクトリは無視して構わない
			continue
		}
		if err := saveUnZipFile(destDir, *f); err != nil {
			return err
		}
	}
	return nil
}

// saveUnZipFile 展開したZipファイルをそのままローカルに保存する
func saveUnZipFile(destDir string, f zip.File) error {
	// 展開先のパスを設定する
	destPath := filepath.Join(destDir, f.Name)
	// 子孫ディレクトリがあれば作成する
	if err := os.MkdirAll(filepath.Dir(destPath), f.Mode()); err != nil {
		return err
	}
	// Zipファイルを開く
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()
	// 展開先ファイルを作成する
	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()
	// 展開先ファイルに書き込む
	if _, err := io.Copy(destFile, rc); err != nil {
		return err
	}

	return nil
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

		fmt.Println(downloads.ChromeDriver)
		for index, cd := range downloads.ChromeDriver {
			cp := cd.Platform
			fmt.Println(index, cd)
			if strings.HasPrefix(cp, getOSKey()) {
				return cd.URL, nil
			} else if strings.HasPrefix(cp, getOSKey()) {
				return cd.URL, nil
			} else if strings.HasPrefix(cp, getOSKey()) {
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
