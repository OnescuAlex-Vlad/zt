package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/rds"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
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
