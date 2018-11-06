package main

import (
	"context"
	"time"

	bh1750 "github.com/d2r2/go-bh1750"
	i2c "github.com/d2r2/go-i2c"
	logger "github.com/d2r2/go-logger"
	shell "github.com/d2r2/go-shell"
)

var lg = logger.NewPackageLogger("main",
	logger.DebugLevel,
	// logger.InfoLevel,
)

func main() {
	defer logger.FinalizeLogger()
	// Create new connection to i2c-bus on 1 line with address 0x40.
	// Use i2cdetect utility to find device address over the i2c-bus
	i2c, err := i2c.NewI2C(0x23, 0)
	if err != nil {
		lg.Fatal(err)
	}
	defer i2c.Close()

	lg.Notify("**********************************************************************************************")
	lg.Notify("*** !!! READ THIS !!!")
	lg.Notify("*** You can change verbosity of output, by modifying logging level of modules \"i2c\", \"bh1750\".")
	lg.Notify("*** Uncomment/comment corresponding lines with call to ChangePackageLogLevel(...)")
	lg.Notify("*** !!! READ THIS !!!")
	lg.Notify("**********************************************************************************************")
	// Uncomment/comment next line to suppress/increase verbosity of output
	// logger.ChangePackageLogLevel("i2c", logger.InfoLevel)
	// logger.ChangePackageLogLevel("bh1750", logger.InfoLevel)

	sensor := bh1750.NewBH1750()
	// Reset sensor
	err = sensor.Reset(i2c)
	if err != nil {
		lg.Fatal(err)
	}
	// Reset sensitivity factor to default value
	err = sensor.ChangeSensivityFactor(i2c, sensor.GetDefaultSensivityFactor())
	if err != nil {
		lg.Fatal(err)
	}

	lg.Notify("**********************************************************************************************")
	lg.Notify("*** Measure ambient light one time")
	lg.Notify("**********************************************************************************************")
	// err = sensor.PowerOn(i2c)
	// if err != nil {
	// 	lg.Fatal(err)
	// }

	resolution := bh1750.LowResolution
	amb, err := sensor.MeasureAmbientLightOneTime(i2c, resolution)
	if err != nil {
		lg.Fatal(err)
	}
	lg.Infof("Ambient light (%s) = %v lx", resolution, amb)
	resolution = bh1750.HighResolution
	amb, err = sensor.MeasureAmbientLightOneTime(i2c, resolution)
	if err != nil {
		lg.Fatal(err)
	}
	lg.Infof("Ambient light (%s) = %v lx", resolution, amb)
	resolution = bh1750.HighestResolution
	amb, err = sensor.MeasureAmbientLightOneTime(i2c, resolution)
	if err != nil {
		lg.Fatal(err)
	}
	lg.Infof("Ambient light (%s) = %v lx", resolution, amb)

	lg.Notify("**********************************************************************************************")
	lg.Notify("*** Measure ambient light continuously")
	lg.Notify("**********************************************************************************************")
	resolution = bh1750.HighResolution
	wait, err := sensor.StartMeasureAmbientLightContinuously(i2c, resolution)
	if err != nil {
		lg.Fatal(err)
	}
	done := make(chan struct{})
	defer close(done)
	// Create context with cancellation possibility.
	ctx, cancel := context.WithCancel(context.Background())
	// Run goroutine waiting for OS termination events, including keyboard Ctrl+C.
	shell.CloseContextOnKillSignal(cancel, done)
	for i := 0; i < 10; i++ {
		// err = sensor.Reset(i2c)
		// if err != nil {
		// 	lg.Fatal(err)
		// }
		amb, err := sensor.FetchMeasuredAmbientLight(i2c)
		if err != nil {
			lg.Fatal(err)
		}
		lg.Infof("Ambient light (%s) = %v lx", resolution, amb)
		select {
		// Check for termination request.
		case <-ctx.Done():
			err = sensor.PowerDown(i2c)
			if err != nil {
				lg.Fatal(err)
			}
			lg.Fatal(ctx.Err())

			// Wait recommended duration.
			// You can increase delay - this
			// doesn't affect to measured value.
		case <-time.After(wait):
		}
	}
	err = sensor.PowerDown(i2c)
	if err != nil {
		lg.Fatal(err)
	}

	lg.Notify("**********************************************************************************************")
	lg.Notify("*** Increase light sensitivity factor in 2 times and repeat measures")
	lg.Notify("**********************************************************************************************")
	err = sensor.ChangeSensivityFactor(i2c, 138)
	if err != nil {
		lg.Fatal(err)
	}
	resolution = bh1750.LowResolution
	amb, err = sensor.MeasureAmbientLightOneTime(i2c, resolution)
	if err != nil {
		lg.Fatal(err)
	}
	lg.Infof("Ambient light (%s) = %v lx", resolution, amb)
	resolution = bh1750.HighResolution
	amb, err = sensor.MeasureAmbientLightOneTime(i2c, resolution)
	if err != nil {
		lg.Fatal(err)
	}
	lg.Infof("Ambient light (%s) = %v lx", resolution, amb)
	resolution = bh1750.HighestResolution
	amb, err = sensor.MeasureAmbientLightOneTime(i2c, resolution)
	if err != nil {
		lg.Fatal(err)
	}
	lg.Infof("Ambient light (%s) = %v lx", resolution, amb)

}
