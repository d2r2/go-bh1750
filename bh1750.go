package bh1750

import (
	"encoding/binary"
	"errors"
	"time"

	i2c "github.com/d2r2/go-i2c"
	"github.com/davecgh/go-spew/spew"
)

// Command bytes
const (
	// No active state.
	CMD_POWER_DOWN = 0x00

	// Waiting for measurement command.
	CMD_POWER_ON = 0x01

	// Reset Data register value.
	// Reset command is not acceptable in Power Down mode.
	CMD_RESET = 0x07

	// Start measurement at 1lx resolution.
	// Measurement Time is typically 120ms.
	CMD_CONTINUOUSLY_H_RES_MODE = 0x10

	// Start measurement at 0.5lx resolution.
	// Measurement Time is typically 120ms.
	CMD_CONTINUOUSLY_H_RES_MODE2 = 0x11

	// Start measurement at 4lx resolution.
	// Measurement Time is typically 16ms.
	CMD_CONTINUOUSLY_L_RES_MODE = 0x13

	// Start measurement at 1lx resolution.
	// Measurement Time is typically 120ms.
	// It is automatically set to Power Down mode after measurement
	CMD_ONE_TIME_H_RES_MODE = 0x20

	// Start measurement at 0.5lx resolution.
	// Measurement Time is typically 120ms.
	// It is automatically set to Power Down mode after measurement.
	CMD_ONE_TIME_H_RES_MODE2 = 0x21

	// Start measurement at 4lx resolution.
	// Measurement Time is typically 16ms.
	// It is automatically set to Power Down mode after measurement.
	CMD_ONE_TIME_L_RES_MODE = 0x23

	// Change measurement time. 01000_MT[7,6,5]
	CMD_CHANGE_MEAS_TIME_HIGH = 0x40

	// Change measurement time. 011_MT[4,3,2,1,0]
	CMD_CHANGE_MEAS_TIME_LOW = 0x60
)

// ResolutionMode define sensor sensitivity
// and measure time. Be aware, that improving
// sensitivity lead to increasing of measurement time.
type ResolutionMode int

const (
	// LowResolution precision 4 lx, 16 ms measurement time
	LowResolution ResolutionMode = iota + 1
	// HighResolution precision 1 lx, 120 ms measurement time
	HighResolution
	// HighestResolution precision 0.5 lx, 120 ms measurement time
	HighestResolution
)

// String define stringer interface.
func (v ResolutionMode) String() string {
	switch v {
	case LowResolution:
		return "Low Resolution"
	case HighResolution:
		return "High Resolution"
	case HighestResolution:
		return "Highest Resolution"
	default:
		return "<unknown>"
	}
}

// BH1750 it's a sensor itself.
type BH1750 struct {
	// Since sensor have no register
	// to report current state, we save
	// last issued command to fill this gap.
	lastCmd        byte
	lastResolution ResolutionMode
	factor         byte
}

// NewBH1750 return new sensor instance.
func NewBH1750() *BH1750 {
	v := &BH1750{}
	v.factor = v.GetDefaultSensivityFactor()
	return v
}

// Reset clear ambient light register value.
func (v *BH1750) Reset(i2c *i2c.I2C) error {
	lg.Debug("Reset sensor...")
	_, err := i2c.WriteBytes([]byte{CMD_RESET})
	if err != nil {
		return err
	}
	v.lastCmd = CMD_RESET
	time.Sleep(time.Microsecond * 3)
	return nil
}

// PowerDown return register to idle state.
func (v *BH1750) PowerDown(i2c *i2c.I2C) error {
	lg.Debug("Power down sensor...")
	_, err := i2c.WriteBytes([]byte{CMD_POWER_DOWN})
	if err != nil {
		return err
	}
	v.lastCmd = CMD_POWER_DOWN
	return nil
}

// PowerOn activate sensor.
func (v *BH1750) PowerOn(i2c *i2c.I2C) error {
	lg.Debug("Power on sensor...")
	_, err := i2c.WriteBytes([]byte{CMD_POWER_ON})
	if err != nil {
		return err
	}
	v.lastCmd = CMD_POWER_ON
	return nil
}

// Get internal parameters used to program sensor.
func (v *BH1750) getResolutionData(resolution ResolutionMode) (cmd byte,
	wait time.Duration, divider uint32) {

	switch resolution {
	case LowResolution:
		cmd = CMD_ONE_TIME_L_RES_MODE
		divider = 1
		// typical measure time is 16 ms,
		// but as it was found 24 ms max time
		// gives better results
		wait = time.Millisecond * 24
	case HighResolution:
		cmd = CMD_ONE_TIME_H_RES_MODE
		divider = 1
		// typical measure time
		wait = time.Millisecond * 120
	case HighestResolution:
		cmd = CMD_ONE_TIME_H_RES_MODE2
		divider = 2
		// typical measure time
		wait = time.Millisecond * 120
	}
	wait = wait * time.Duration(v.factor) /
		time.Duration(v.GetDefaultSensivityFactor())

	return cmd, wait, divider
}

// MeasureAmbientLightOneTime measure and return ambient light once in lux.
func (v *BH1750) MeasureAmbientLightOneTime(i2c *i2c.I2C,
	resolution ResolutionMode) (uint16, error) {

	lg.Debug("Run one time measure...")

	cmd, wait, divider := v.getResolutionData(resolution)

	v.lastCmd = cmd
	v.lastResolution = resolution

	_, err := i2c.WriteBytes([]byte{cmd})
	if err != nil {
		return 0, err
	}

	time.Sleep(wait)

	var data struct {
		Data [2]byte
	}
	err = readDataToStruct(i2c, 2, binary.BigEndian, &data)
	if err != nil {
		return 0, err
	}

	amb := uint16(uint32(uint16(data.Data[0])<<8|uint16(data.Data[1])) *
		5 / 6 / divider)

	return amb, nil
}

// StartMeasureAmbientLightContinuously start continuous
// measurement process. Use FetchMeasuredAmbientLight to get
// average ambient light amount collected and calculated over a time.
// Use PowerDown to stop measurements and return sensor to idle state.
func (v *BH1750) StartMeasureAmbientLightContinuously(i2c *i2c.I2C,
	resolution ResolutionMode) (wait time.Duration, err error) {

	lg.Debug("Start measures continuously...")

	cmd, wait, _ := v.getResolutionData(resolution)

	v.lastCmd = cmd
	v.lastResolution = resolution

	_, err = i2c.WriteBytes([]byte{cmd})
	if err != nil {
		return 0, err
	}

	// Wait first time to collect necessary
	// amount of light for correct results.
	// It's not necessary to wait next time
	// same amount of time, because
	// sensor accumulate average lux amount
	// without any overwrite old value.
	time.Sleep(wait)

	// In any case we are returning same
	// recommended amount of time to wait
	// between measures.
	return wait, nil
}

// FetchMeasuredAmbientLight return current average ambient light in lux.
// Previous command should be any continuous measurement initiation,
// otherwise error will be reported.
func (v *BH1750) FetchMeasuredAmbientLight(i2c *i2c.I2C) (uint16, error) {

	lg.Debug("Fetch measured data...")

	cmd, _, divider := v.getResolutionData(v.lastResolution)

	if v.lastCmd != cmd {
		return 0, errors.New(
			"can't fetch measured ambient light, since last command doesn't match")
	}

	var data struct {
		Data [2]byte
	}
	err := readDataToStruct(i2c, 2, binary.BigEndian, &data)
	if err != nil {
		return 0, err
	}

	amb := uint16(uint32(uint16(data.Data[0])<<8|uint16(data.Data[1])) *
		5 / 6 / divider)

	return amb, nil
}

// GetDefaultSensivityFactor return factor value
// used when your sensor have no any protection cover.
// This is default setting according to specification.
func (v *BH1750) GetDefaultSensivityFactor() byte {
	return 69
}

// ChangeSensivityFactor used when you close sensor
// with protection cover, which change (ordinary decrease)
// expected amount of light falling on the sensor.
// In this case you should calibrate you sensor and find
// appropriate factor to get in output correct ambient light value.
// Be aware, that improving sensitivity will increase
// measurement time.
func (v *BH1750) ChangeSensivityFactor(i2c *i2c.I2C, factor byte) error {

	lg.Debug("Change sensitivity factor...")

	// minimum limit
	const minValue = 31
	// maximum limit
	const maxValue = 254

	if factor < minValue || factor > maxValue {
		return errors.New(spew.Sprintf("sensitivity factor value exceed range [%d..%d]",
			minValue, maxValue))
	}

	high := (factor & 0xE0) >> 5
	_, err := i2c.WriteBytes([]byte{CMD_CHANGE_MEAS_TIME_HIGH | high})
	if err != nil {
		return err
	}

	low := (factor & 0x1F)
	_, err = i2c.WriteBytes([]byte{CMD_CHANGE_MEAS_TIME_LOW | low})
	if err != nil {
		return err
	}

	v.factor = factor

	return nil
}
