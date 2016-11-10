package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
)

const (
	NORMAL  = "normal"
	STANDBY = "standby"
	UNKNOWN = "unknown"
)

type appXmlCfg struct {
	CMPP2Connect struct {
		Enable bool `xml:"enable"`
	}
}

func getDb2Mode() (string, error) {
	f, err := os.Open(cfg.Mas.Appxml.Run)
	if err != nil {
		return "", err
	}
	defer f.Close()

	d := xml.NewDecoder(f)
	cmpp2 := appXmlCfg{}
	err = d.Decode(&cmpp2)
	if err != nil {
		return "", err
	}

	if cmpp2.CMPP2Connect.Enable == true {
		return STANDBY, nil
	}
	return NORMAL, nil
}

func switchDb2Mode(wf flushWriter) error {
	wf.Send("begin\n\nreplace app.xml ...\n")

	cur, err := getDb2Mode()
	if err != nil {
		return err
	}
	if cur == NORMAL {
		wf.Send(fmt.Sprintf("replace app.xml with '%s'\n", cfg.Mas.Appxml.Standby))
		err = copyFile(cfg.Mas.Appxml.Standby, cfg.Mas.Appxml.Run)
	} else {
		wf.Send(fmt.Sprintf("replace app.xml with '%s'\n", cfg.Mas.Appxml.Normal))
		err = copyFile(cfg.Mas.Appxml.Normal, cfg.Mas.Appxml.Run)
	}

	rst := "success"
	if err != nil {
		rst = "fail" + err.Error()
	}
	wf.Send(fmt.Sprintf("replace app.xml %s\n", rst))
	return err
}

type flushWriter struct {
	f http.Flusher
	w io.Writer
}

func (fw flushWriter) Send(msg string) {
	fw.w.Write([]byte(msg))
	if fw.f != nil {
		fw.f.Flush()
	}
}
func restartDb2Service(wf flushWriter) error {
	var err error = nil
	if cfg.Mas.Workdir != "" {
		cwd, _ := os.Getwd()
		err = os.Chdir(cfg.Mas.Workdir)
		if err != nil {
			return err
		}
		defer os.Chdir(cwd)
	}

	cmd := exec.Command("sh", cfg.Mas.Restart)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}

	go func() {
		buff := make([]byte, 1024)
		for {
			n, err := stdout.Read(buff)
			if err != nil {
				return
			}
			wf.Send(string(buff[:n]))
		}
	}()
	if err := cmd.Wait(); err != nil {
		return err
	}
	wf.Send("\nover\n")
	return nil
}

func copyFile(src string, dst string) error {
	data, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(dst, data, 0777)
}
