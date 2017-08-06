package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/Sirupsen/logrus"
	"github.com/docker/go-plugins-helpers/volume"
)

const driverName = "samplefs"

type driver struct {
	mountpoint string
}

func (d driver) Create(request volume.Request) volume.Response {
	logrus.WithField("method", "create").Debugf("%#v", request)
	name := request.Name
	err := os.MkdirAll(path.Join(d.mountpoint, name), 777)
	var response volume.Response
	if err != nil {
		response = volume.Response{Err: err.Error()}
	} else {
		response = volume.Response{}
	}
	logrus.WithField("method", "create").Debugf("%#v", response)
	return response
}

func (d driver) List(request volume.Request) volume.Response {
	logrus.WithField("method", "list").Debugf("%#v", request)
	var vols []*volume.Volume
	dirs, err := ioutil.ReadDir(d.mountpoint)
	if err != nil {
		return volume.Response{Err: err.Error()}
	}
	for _, dir := range dirs {
		name := dir.Name()
		vols = append(vols, &volume.Volume{
			Name:       name,
			Mountpoint: path.Join(d.mountpoint, name),
		})
	}
	response := volume.Response{Volumes: vols}
	logrus.WithField("method", "list").Debugf("%#v", response)
	return response
}

func (d driver) Get(request volume.Request) volume.Response {
	logrus.WithField("method", "get").Debugf("%#v", request)
	name := request.Name
	mountpoint := path.Join(d.mountpoint, name)
	_, err := os.Stat(mountpoint)
	var response volume.Response
	if os.IsNotExist(err) {
		response = volume.Response{Err: "not found"}
	} else {
		response = volume.Response{Volume: &volume.Volume{
			Name:       name,
			Mountpoint: path.Join(d.mountpoint, name),
		}}
	}
	logrus.WithField("method", "get").Debugf("%#v", response)
	return response
}

func (d driver) Remove(request volume.Request) volume.Response {
	logrus.WithField("method", "remove").Debugf("%#v", request)
	err := os.RemoveAll(path.Join(d.mountpoint, request.Name))
	var response volume.Response
	if err != nil {
		response = volume.Response{Err: err.Error()}
	} else {
		response = volume.Response{}
	}
	logrus.WithField("method", "remove").Debugf("%#v", response)
	return response
}

func (d driver) Path(request volume.Request) volume.Response {
	logrus.WithField("method", "path").Debugf("%#v", request)
	name := request.Name
	mountpoint := path.Join(d.mountpoint, name)
	_, err := os.Stat(mountpoint)
	var response volume.Response
	if os.IsNotExist(err) {
		response = volume.Response{Err: "not found"}
	} else {
		response = volume.Response{Mountpoint: path.Join(d.mountpoint, name)}
	}
	logrus.WithField("method", "path").Debugf("%#v", response)
	return response
}

func (d driver) Mount(request volume.MountRequest) volume.Response {
	logrus.WithField("method", "mount").Debugf("%#v", request)
	name := request.Name
	mountpoint := path.Join(d.mountpoint, name)
	_, err := os.Stat(mountpoint)
	var response volume.Response
	if os.IsNotExist(err) {
		response = volume.Response{Err: "not found"}
	} else {
		response = volume.Response{Mountpoint: path.Join(d.mountpoint, name)}
	}
	logrus.WithField("method", "mount").Debugf("%#v", response)
	return response
}

func (d driver) Unmount(request volume.UnmountRequest) volume.Response {
	logrus.WithField("method", "unmount").Debugf("%#v", request)
	response := volume.Response{}
	logrus.WithField("method", "unmount").Debugf("%#v", response)
	return response
}

func (d driver) Capabilities(request volume.Request) volume.Response {
	logrus.WithField("method", "capabilities").Debugf("%#v", request)
	response := volume.Response{Capabilities: volume.Capability{Scope: "local"}}
	logrus.WithField("method", "capabilities").Debugf("%#v", response)
	return response
}

func newDriver(mountpoint string) volume.Driver {
	return driver{mountpoint}
}

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	d := newDriver("/mnt")
	h := volume.NewHandler(d)
	fmt.Println(h.ServeUnix(driverName, 0))
}
