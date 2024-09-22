package main

import (
	"github.com/aonescu/zt/pkg/action"
	"github.com/aonescu/zt/pkg/bluetooth"
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/rds"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)



// mockLLMPrompt simulates an LLM prompt. In real use cases, this would come from a processed LLM response.
func mockLLMPrompt() map[string]interface{} {
	return map[string]interface{}{
		"action":          "connect_bluetooth",
		"bluetooth_device": "MySpeaker", // Device name to connect to
	}
}

// runBluetoothAction handles the Bluetooth connection action using the provided LLM input.
func runBluetoothAction() {
	// Initialize Bluetooth connector
	btConnector, err := bluetooth.NewBluetoothConnector()
	if err != nil {
		fmt.Printf("Error initializing Bluetooth connector: %v\n", err)
		return
	}

	// Define the Bluetooth connect handler
	handler := func(input map[string]interface{}) (map[string]interface{}, error) {
		deviceName, ok := input["bluetooth_device"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid input: missing bluetooth_device")
		}

		// Connect to Bluetooth device
		err := btConnector.ConnectToDevice(deviceName)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to device: %w", err)
		}

		return map[string]interface{}{
			"connection_status": "connected",
		}, nil
	}

	// Create a new action for connecting to Bluetooth
	connectBluetooth := action.NewAction(
		"Connect Bluetooth",
		"Connect to a Bluetooth speaker",
		"medium",
		handler,
	)

	// Simulate LLM prompt
	llmPrompt := mockLLMPrompt()

	// Update action input from the LLM prompt
	connectBluetooth.UpdateInput(llmPrompt)

	// Execute the action
	err = connectBluetooth.Execute()
	if err != nil {
		fmt.Printf("Error executing action: %v\n", err)
	} else {
		fmt.Printf("Action output: %v\n", connectBluetooth.Output)
	}

	// Optionally disconnect from the Bluetooth device after operation
	err = btConnector.Disconnect()
	if err != nil {
		fmt.Printf("Error disconnecting: %v\n", err)
	}
}

// runPulumiInfrastructure handles the Pulumi infrastructure setup (AWS RDS and security groups).
func runPulumiInfrastructure() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create a security group for the database
		sg, err := ec2.NewSecurityGroup(ctx, "db-sg", &ec2.SecurityGroupArgs{
			Description: pulumi.String("Allow database access"),
			Ingress: ec2.SecurityGroupIngressArray{
				&ec2.SecurityGroupIngressArgs{
					Protocol:   pulumi.String("tcp"),
					FromPort:   pulumi.Int(5432),
					ToPort:     pulumi.Int(5432),
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				},
			},
		})
		if err != nil {
			return err
		}

		// Create a PostgreSQL RDS instance
		db, err := rds.NewInstance(ctx, "corebit-db", &rds.InstanceArgs{
			Engine:              pulumi.String("postgres"),
			InstanceClass:       pulumi.String("db.t3.micro"),
			AllocatedStorage:    pulumi.Int(20),
			Name:                pulumi.String("corebit"),
			Username:            pulumi.String("admin"),
			Password:            pulumi.String("password123"),
			VpcSecurityGroupIds: pulumi.StringArray{sg.ID()},
		})
		if err != nil {
			return err
		}

		// Output database endpoint
		ctx.Export("dbEndpoint", db.Endpoint)

		return nil
	})
}

func main() {
	// Execute the Bluetooth action
	runBluetoothAction()

	// Execute the Pulumi infrastructure setup
	runPulumiInfrastructure()
}