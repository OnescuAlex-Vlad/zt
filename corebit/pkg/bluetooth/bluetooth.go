package bluetooth

import (
	"errors"
	"fmt"
	"tinygo.org/x/bluetooth"
)

// BluetoothConnector represents a struct that handles Bluetooth connections.
type BluetoothConnector struct {
	adapter bluetooth.Adapter
	device  *bluetooth.Device
}

// NewBluetoothConnector initializes a new BluetoothConnector.
func NewBluetoothConnector() (*BluetoothConnector, error) {
	adapter := bluetooth.DefaultAdapter
	if err := adapter.Enable(); err != nil {
		return nil, fmt.Errorf("failed to enable Bluetooth adapter: %w", err)
	}

	return &BluetoothConnector{
		adapter: *adapter,
	}, nil
}

// ConnectToDevice connects to a Bluetooth device by its name.
func (bc *BluetoothConnector) ConnectToDevice(deviceName string) error {
	fmt.Printf("Starting scan for Bluetooth devices...\n")

	var device bluetooth.Device
	err := bc.adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
		if result.LocalName() == deviceName {
			fmt.Printf("Found device: %s\n", result.LocalName())
			device, _ = adapter.Connect(result.Address, bluetooth.ConnectionParams{})
			bc.device = &device
			adapter.StopScan()
		}
	})

	if err != nil {
		return fmt.Errorf("failed to scan for devices: %w", err)
	}

	if device.Address.String() == "" {
		return errors.New("device not found")
	}

	fmt.Printf("Connected to Bluetooth device: %s\n", deviceName)
	return nil
}

// Disconnect disconnects from the Bluetooth device.
func (bc *BluetoothConnector) Disconnect() error {
	if bc.device == nil {
		return errors.New("no device to disconnect")
	}

	if err := bc.device.Disconnect(); err != nil {
		return fmt.Errorf("failed to disconnect from device: %w", err)
	}

	fmt.Printf("Disconnected from Bluetooth device.\n")
	return nil
}