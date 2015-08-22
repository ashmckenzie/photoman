package main

import (
  "fmt"
  "io"
  log "github.com/Sirupsen/logrus"
  "os"
  "path/filepath"
  "strings"
  "time"

  "github.com/gosexy/exif"
  "github.com/codegangsta/cli"
  // "github.com/davecgh/go-spew/spew"
)

var mode string;
var supportedPhotoTypes = map[string]bool{
  ".jpg": true,
}

func isSupportPhotoType(extension string) bool {
  if _, ok := supportedPhotoTypes[extension]; ok {
    return true
  } else {
    return false
  }
}

func procesImage(path string, f os.FileInfo, err error) error {
  if f.IsDir() { return nil }

  log.Debugf("Processing %s", path)

  extension := filepath.Ext(f.Name())
  if !isSupportPhotoType(strings.ToLower(extension)) {
    log.Warnf("%s's file type %s is unsupported", path, extension)
    return nil
  }

  reader := exif.New()

  err = reader.Open(path)
  if err != nil { log.Fatal(err) }

  str := fmt.Sprintf("%s", reader.Tags["Date and Time"])
  t := f.ModTime()

  if len(str) == 0 {
    log.Warnf("Date and Time EXIF tag missing for %s", path)
  } else {
    layout := "2006:01:02 15:04:05"
    t, err = time.Parse(layout, str)
    if err != nil { log.Fatal(err) }
  }

  newDir := fmt.Sprintf("%4d/%02d/%02d", t.Year(), t.Month(), t.Day())

  err = os.MkdirAll(newDir, 0777)
  if err != nil { log.Fatal(err) }

  newFile := fmt.Sprintf("%s/%s", newDir, f.Name())

  if mode == "move" {
    log.Debugf("Moving %s %s", path, newFile)
    err = os.Rename(path, newFile)
  } else {
    if _, err := os.Stat(newFile); err == nil {
      log.Warnf("Photo %s already exists", newFile)
    } else {
      log.Debugf("Copying %s %s", path, newFile)
      err = copyFile(path, newFile)
    }
  }

  if err != nil { log.Fatal(err) }

  return nil
}

func copyFile(src, dst string) error {
	s, err := os.Open(src)
	if err != nil { return err }

	// no need to check errors on read only file, we already got everything
	// we need from the filesystem, so nothing can go wrong now.
	defer s.Close()
	d, err := os.Create(dst)
	if err != nil { return err }

	if _, err := io.Copy(d, s); err != nil {
		d.Close()
		return err
	}

	return d.Close()
}

func setupLogging() {
  logLevel := log.InfoLevel
  if os.Getenv("DEBUG") == "true" { logLevel = log.DebugLevel }
  log.SetLevel(logLevel)
}

func main() {
  setupLogging()

  app := cli.NewApp()
  app.EnableBashCompletion = true
  app.Name = "photoman"
  app.Version = "0.0.1"
  app.Usage = "Manage your photos into nice a simple YYYY/MM/DD structure"

  app.Flags = []cli.Flag {
    cli.StringFlag{
      Name: "mode",
      Value: "move",
      Usage: "Move or copy photos",
    },
  }

  app.Action = func(c *cli.Context) {
    path := c.Args().First()
    if len(path) == 0 { log.Fatal("Please specify a path to process!") }

    if c.String("mode") == "move" || c.String("mode") == "copy" {
      mode = c.String("mode")
    } else {
      log.Fatal("Support --mode's are move or copy")
    }

    err := filepath.Walk(path, procesImage)
    if err != nil { log.Fatal(err) }
  }

  app.Run(os.Args)
}
