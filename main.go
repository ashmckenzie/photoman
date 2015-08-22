package main

import (
  "os"
  log "github.com/Sirupsen/logrus"
  "fmt"
  "time"
  "path/filepath"

  "github.com/gosexy/exif"
  "github.com/codegangsta/cli"
  // "github.com/davecgh/go-spew/spew"
)

func procesImage(path string, f os.FileInfo, err error) error {
  if f.IsDir() { return nil }

  reader := exif.New()

  log.Debugf("Processing %s", path)
  err = reader.Open(path)
  if err != nil { log.Fatal(err) }

  str := fmt.Sprintf("%s", reader.Tags["Date and Time"])
  t := f.ModTime()

  if len(str) == 0 {
    log.Warnf("Date and Time EXIF tag missing for %s, falling back to mtime", path)
  } else {
    layout := "2006:01:02 15:04:05"
    t, err = time.Parse(layout, str)
    if err != nil { log.Fatal(err) }
  }

  newDir := fmt.Sprintf("%4d/%02d/%02d", t.Year(), t.Month(), t.Day())

  err = os.MkdirAll(newDir, 0777)
  if err != nil { log.Fatal(err) }

  newFile := fmt.Sprintf("%s/%s", newDir, f.Name())
  log.Debugf("Moving %s %s", path, newFile)
  err = os.Rename(path, newFile)
  if err != nil { log.Fatal(err) }

  return nil
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
  app.Action = func(c *cli.Context) {
    path := c.Args().First()
    if len(path) == 0 { log.Fatal("Please specify a path to process!") }
    err := filepath.Walk(path, procesImage)
    if err != nil { log.Fatal(err) }
  }

  app.Run(os.Args)


}
