
BH1750 ambient light sensor
=====================

[![Build Status](https://travis-ci.org/d2r2/go-bh1750.svg?branch=master)](https://travis-ci.org/d2r2/go-bh1750)
[![Go Report Card](https://goreportcard.com/badge/github.com/d2r2/go-bh1750)](https://goreportcard.com/report/github.com/d2r2/go-bh1750)
[![GoDoc](https://godoc.org/github.com/d2r2/go-bh1750?status.svg)](https://godoc.org/github.com/d2r2/go-bh1750)
[![MIT License](http://img.shields.io/badge/License-MIT-yellow.svg)](./LICENSE)

BH1750 ([general specification](https://raw.github.com/d2r2/go-bh1750/master/docs/bh1750fvi-e-186247.pdf)) is a power effective ambient light sensor with spectral response close to human eye. Sensor returns measured ambient light value in lux units. Easily integrated with Arduino and Raspberry PI via i2c communication interface:
![image](https://raw.github.com/d2r2/go-bh1750/master/docs/bh1750.jpg)

Here is a library written in [Go programming language](https://golang.org/) for Raspberry PI and counterparts, which gives you in the output ambient light value (making all necessary i2c-bus interacting and values computing).

Golang usage
------------


```go
func main() {
	// Create new connection to i2c-bus on 0 line with address 0x23.
	// Use i2cdetect utility to find device address over the i2c-bus
	i2c, err := i2c.NewI2C(0x23, 0)
	if err != nil {
		log.Fatal(err)
	}
	defer i2c.Close()

	sensor := bh1750.NewBH1750()

	resolution := bh1750.HighResolution
	amb, err := sensor.MeasureAmbientLight(i2c, resolution)
	if err != nil {
		lg.Fatal(err)
	}
	log.Printf("Ambient light (%s) = %v lx", resolution, amb)
```


Getting help
------------

GoDoc [documentation](http://godoc.org/github.com/d2r2/go-bh1750)

Installation
------------

```bash
$ go get -u github.com/d2r2/go-bh1750
```

Troubleshooting
--------------

- *How to obtain fresh Golang installation to RPi device (either any RPi clone):*
If your RaspberryPI golang installation taken by default from repository is outdated, you may consider
to install actual golang manually from official Golang [site](https://golang.org/dl/). Download
tar.gz file containing armv6l in the name. Follow installation instructions.

- *How to enable I2C bus on RPi device:*
If you employ RaspberryPI, use raspi-config utility to activate i2c-bus on the OS level.
Go to "Interfacing Options" menu, to active I2C bus.
Probably you will need to reboot to load i2c kernel module.
Finally you should have device like /dev/i2c-1 present in the system.

- *How to find I2C bus allocation and device address:*
Use i2cdetect utility in format "i2cdetect -y X", where X may vary from 0 to 5 or more,
to discover address occupied by peripheral device. To install utility you should run
`apt install i2c-tools` on debian-kind system. `i2cdetect -y 1` sample output:
	```
	     0  1  2  3  4  5  6  7  8  9  a  b  c  d  e  f
	00:          -- -- -- -- -- -- -- -- -- -- -- -- --
	10: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
	20: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
	30: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
	40: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
	50: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
	60: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
	70: -- -- -- -- -- -- 76 --    
	```

Contact
-------

Please use [Github issue tracker](https://github.com/d2r2/go-bh1750/issues) for filing bugs or feature requests.


License
-------

Go-bh1750 is licensed under MIT License.
