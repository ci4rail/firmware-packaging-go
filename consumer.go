/*
Copyright © 2021 Ci4Rail GmbH <engineering@ci4rail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package fwpkg

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
)

// FirmwarePackageConsumer is a handle to consume firmware package archive files
type FirmwarePackageConsumer struct {
	fileName string
	file     *os.File
	manifest *FwManifest
}

// NewFirmwarePackageConsumerFromFile creates an object to work with the firmware package in fileName
// The file is opened and the manifest is parsed and checed for validity. If not valid, an error is returned
func NewFirmwarePackageConsumerFromFile(fileName string) (*FirmwarePackageConsumer, error) {
	p := new(FirmwarePackageConsumer)
	p.fileName = fileName

	f, err := os.Open(fileName)
	if err != nil {
		return nil, errors.New("can't open " + fileName + ": " + err.Error())
	}
	p.file = f

	m, err := p.getManifest()
	if err != nil {
		return nil, err
	}
	p.manifest = m
	return p, nil
}

// Manifest returns the parsed manifest as ago struct
func (p *FirmwarePackageConsumer) Manifest() (manifest *FwManifest) {
	return p.manifest
}

// File extracts the firmware binary file from the firmware package and writes it to w
func (p *FirmwarePackageConsumer) File(w io.Writer) error {

	_, err := p.file.Seek(0, 0)
	if err != nil {
		return err
	}
	err = untarFileContent(p.file, "./"+p.manifest.File, w)
	if err != nil {
		return errors.New("can't untar firmware binary " + p.manifest.File + ": " + err.Error())
	}
	return nil
}

func (p *FirmwarePackageConsumer) getManifest() (*FwManifest, error) {
	mJSON := new(bytes.Buffer)

	_, err := p.file.Seek(0, 0)
	if err != nil {
		return nil, err
	}
	err = untarFileContent(p.file, "./manifest.json", mJSON)
	if err != nil {
		return nil, errors.New("can't untar manifest: " + err.Error())
	}
	m, err := decodeManifest(mJSON.Bytes())
	if err != nil {
		return nil, errors.New("error in manifest: " + err.Error())
	}
	return m, nil
}

func decodeManifest(b []byte) (*FwManifest, error) {
	var m *FwManifest

	err := json.Unmarshal(b, &m)
	if err != nil {
		return nil, errors.New("can't decode manifest: " + err.Error())
	}
	if m.Name == "" {
		return nil, errors.New("missing \"name\" in manifest")
	}
	if m.Version == "" {
		return nil, errors.New("missing \"version\" in manifest")
	}
	if m.File == "" {
		return nil, errors.New("missing \"file\" in manifest")
	}
	if m.Compatibility.HW == "" {
		return nil, errors.New("missing \"compatibility.hw\" in manifest")
	}
	if len(m.Compatibility.MajorRevs) == 0 {
		return nil, errors.New("missing \"compatibility.major_revs\" in manifest")
	}
	return m, nil
}
