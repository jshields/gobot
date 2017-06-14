package curie

import (
	"errors"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/firmata"
)

const (
	CURIE_IMU              = 0x11
	CURIE_IMU_READ_ACCEL   = 0x00
	CURIE_IMU_READ_GYRO    = 0x01
	CURIE_IMU_READ_TEMP    = 0x02
	CURIE_IMU_SHOCK_DETECT = 0x03
	CURIE_IMU_STEP_COUNTER = 0x04
	CURIE_IMU_TAP_DETECT   = 0x05
	CURIE_IMU_READ_MOTION  = 0x06
)

type AccelerometerData struct {
	X int16
	Y int16
	Z int16
}

type GyroscopeData struct {
	X int16
	Y int16
	Z int16
}

// IMUDriver represents the IMU that is built-in to the Curie
type IMUDriver struct {
	name       string
	connection *firmata.Adaptor
	gobot.Eventer
}

// NewIMUDriver returns a new IMUDriver
func NewIMUDriver(a *firmata.Adaptor) *IMUDriver {
	imu := &IMUDriver{
		name:       gobot.DefaultName("CurieIMU"),
		connection: a,
		Eventer:    gobot.NewEventer(),
	}

	return imu
}

// Start starts up the IMUDriver
func (imu *IMUDriver) Start() (err error) {
	imu.connection.On("SysexResponse", func(res interface{}) {
		data := res.([]byte)
		if data[1] == CURIE_IMU {
			switch data[2] {
			case CURIE_IMU_READ_ACCEL:
				val, _ := parseAccelerometerData(data)
				imu.Publish("Accelerometer", val)

			case CURIE_IMU_READ_GYRO:
				val, _ := parseGyroscopeData(data)
				imu.Publish("Gyroscope", val)

			case CURIE_IMU_READ_TEMP:
				val, _ := parseTemperatureData(data)
				imu.Publish("Temperature", val)

			}
		}
	})
	return
}

// Halt stops the IMUDriver
func (imu *IMUDriver) Halt() (err error) {
	return
}

// Name returns the IMUDriver's name
func (imu *IMUDriver) Name() string { return imu.name }

// SetName sets the IMUDriver'ss name
func (imu *IMUDriver) SetName(n string) { imu.name = n }

// Connection returns the IMUDriver's Connection
func (imu *IMUDriver) Connection() gobot.Connection { return imu.connection }

// ReadAccelerometer calls the Curie's built-in accelerometer. The result will
// be returned by the Sysex response message
func (imu *IMUDriver) ReadAccelerometer() error {
	return imu.connection.WriteSysex([]byte{CURIE_IMU, CURIE_IMU_READ_ACCEL})
}

// ReadGyroscope calls the Curie's built-in gyroscope. The result will
// be returned by the Sysex response message
func (imu *IMUDriver) ReadGyroscope() error {
	return imu.connection.WriteSysex([]byte{CURIE_IMU, CURIE_IMU_READ_GYRO})
}

// ReadTemperature calls the Curie's built-in temperature sensor.
// The result will be returned by the Sysex response message
func (imu *IMUDriver) ReadTemperature() error {
	return imu.connection.WriteSysex([]byte{CURIE_IMU, CURIE_IMU_READ_TEMP})
}

func parseAccelerometerData(data []byte) (*AccelerometerData, error) {
	if len(data) < 9 {
		return nil, errors.New("Invalid data")
	}
	x := int16(uint16(data[3]) | uint16(data[4])<<7)
	y := int16(uint16(data[5]) | uint16(data[6])<<7)
	z := int16(uint16(data[7]) | uint16(data[8])<<7)

	res := &AccelerometerData{X: x, Y: y, Z: z}
	return res, nil
}

func parseGyroscopeData(data []byte) (*GyroscopeData, error) {
	if len(data) < 9 {
		return nil, errors.New("Invalid data")
	}
	x := int16(uint16(data[3]) | uint16(data[4])<<7)
	y := int16(uint16(data[5]) | uint16(data[6])<<7)
	z := int16(uint16(data[7]) | uint16(data[8])<<7)

	res := &GyroscopeData{X: x, Y: y, Z: z}
	return res, nil
}

func parseTemperatureData(data []byte) (float32, error) {
	if len(data) < 8 {
		return 0, errors.New("Invalid data")
	}
	t1 := int16(uint16(data[3]) | uint16(data[4])<<7)
	t2 := int16(uint16(data[5]) | uint16(data[6])<<7)

	res := (float32(t1+(t2*8)) / 512.0) + 23.0
	return res, nil
}
