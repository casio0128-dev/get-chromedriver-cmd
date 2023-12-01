package main

import (
	"os"
)

type ChromeDriver struct {
	majorVersion string
	version      string
	zipFile      os.File
}

func (c *ChromeDriver) MajorVersion() string {
	return c.majorVersion
}

func (c *ChromeDriver) SetMajorVersion(majorVersion string) {
	c.majorVersion = majorVersion
}

func (c *ChromeDriver) Version() string {
	return c.version
}

func (c *ChromeDriver) SetVersion(version string) {
	c.version = version
}

func (c *ChromeDriver) ZipFile() os.File {
	return c.zipFile
}

func (c *ChromeDriver) SetZipFile(zipFile os.File) {
	c.zipFile = zipFile
}
