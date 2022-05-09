# firmware-packaging-go
A go library to read firmware packages. Initially created for io4edge firmware packages.

A firmware package is a tar file containing a manifest and the firmware binary:

```bash 
$ tar tvf testdata/t1.fwpkg
drwx------ root/root         0 2021-09-07 14:24 ./
-rw-r--r-- root/root    155688 2021-09-05 13:01 ./fw_varA_1_1_0.json
-rw-r--r-- root/root       179 2021-09-07 14:24 ./manifest.json
```

where `manifest.json` describes the firmware binary:
```json
{
  "name": "cpu01-tty_accdl",
  "version": "1.1.0",
  "file": "fw_varA_1_1_0.json",
  "compatibility": {
    "hw": "s101-cpu01",
    "major_revs": [
      1,
      2
    ]
  }
}
```

Example usage of the library:
```go
    import fwpkg "github.com/ci4rail/firmware-packaging-go"

	pkg, err := fwpkg.NewFirmwarePackageConsumerFromFile("t1.fwpkg")
	if err != nil {
		...
	}
	manifest := pkg.Manifest()
    fmt.Printf("manifest: %v\n", manifest)

    // load fw file into memory
    fwFile := new(bytes.Buffer)
	err = pkg.File(fwFile)
	if err != nil {
		...
	}
```
