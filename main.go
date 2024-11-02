package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/core"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load .env file: %v", err)
	}
}

func main() {
	config := common.DefaultConfigProvider() // deafault ~/.oci/config 설정 사용

	computeClient, err := core.NewComputeClientWithConfigurationProvider(config)
	if err != nil {
		log.Fatalf("Failed to create compute client: %v", err)
	}

	instance, err := createInstance(computeClient)
	if err != nil {
		log.Fatalf("Failed to create instance: %v", err)
	}
	fmt.Printf("Instance created: %v\n", *instance.Id)
}

func createInstance(computeClient core.ComputeClient) (*core.Instance, error) {
	compartmentID := os.Getenv("OCI_COMPARTMENT_ID")
	subnetID := os.Getenv("OCI_SUBNET_ID")
	imageID := os.Getenv("OCI_IMAGE_ID")
	availabilityDomain := os.Getenv("OCI_AVAILABILITY_DOMAIN")
	shape := os.Getenv("OCI_SHAPE")
	pubKey, err := os.ReadFile("/home/ubuntu/.oci/ssh-key-2024-11-01.key.pub")
	if err != nil {
		log.Fatalf("Failed to read public key: %v", err)
	}
	bootVolumeSizeInGBs := os.Getenv("OCI_BOOT_VOLUME_SIZE_IN_GBS")
	bootVolumeSizeInGBsInt, err := strconv.ParseInt(bootVolumeSizeInGBs, 10, 64)
	if err != nil {
		log.Fatalf("Failed to parse boot volume size in GBs: %v", err)
	}

	ocupus := os.Getenv("OCI_OCPUS")
	ocpusFloat, err := strconv.ParseFloat(ocupus, 32)
	if err != nil {
		log.Fatalf("Failed to parse ocpus: %v", err)
	}
	memoryInGBs := os.Getenv("OCI_MEMORY_IN_GBS")
	memoryInGBsFloat, err := strconv.ParseFloat(memoryInGBs, 32)
	if err != nil {
		log.Fatalf("Failed to parse memory in GBs: %v", err)
	}

	launchDetails := core.LaunchInstanceDetails{
		CompartmentId:      &compartmentID,
		AvailabilityDomain: &availabilityDomain,
		DisplayName:        common.String(os.Getenv("OCI_DISPLAY_NAME")),
		Shape:              &shape,
		CreateVnicDetails: &core.CreateVnicDetails{
			AssignPublicIp:         common.Bool(true),
			SubnetId:               &subnetID,
			AssignPrivateDnsRecord: common.Bool(true),
			AssignIpv6Ip:           common.Bool(false),
		},
		SourceDetails: core.InstanceSourceViaImageDetails{
			ImageId:             &imageID,
			BootVolumeSizeInGBs: common.Int64(bootVolumeSizeInGBsInt),
			// BootVolumeVpusPerGB: common.Int64(10),
		},
		IsPvEncryptionInTransitEnabled: common.Bool(true),
		Metadata: map[string]string{
			"ssh_authorized_keys": string(pubKey),
		},
		AgentConfig: &core.LaunchInstanceAgentConfigDetails{
			PluginsConfig: []core.InstanceAgentPluginConfigDetails{
				{Name: common.String("Vulnerability Scanning"), DesiredState: core.InstanceAgentPluginConfigDetailsDesiredStateDisabled},
				{Name: common.String("Management Agent"), DesiredState: core.InstanceAgentPluginConfigDetailsDesiredStateDisabled},
				{Name: common.String("Custom Logs Monitoring"), DesiredState: core.InstanceAgentPluginConfigDetailsDesiredStateEnabled},
				{Name: common.String("Compute RDMA GPU Monitoring"), DesiredState: core.InstanceAgentPluginConfigDetailsDesiredStateDisabled},
				{Name: common.String("Compute Instance Monitoring"), DesiredState: core.InstanceAgentPluginConfigDetailsDesiredStateEnabled},
				{Name: common.String("Compute HPC RDMA Auto-Configuration"), DesiredState: core.InstanceAgentPluginConfigDetailsDesiredStateDisabled},
				{Name: common.String("Compute HPC RDMA Authentication"), DesiredState: core.InstanceAgentPluginConfigDetailsDesiredStateDisabled},
				{Name: common.String("Cloud Guard Workload Protection"), DesiredState: core.InstanceAgentPluginConfigDetailsDesiredStateEnabled},
				{Name: common.String("Block Volume Management"), DesiredState: core.InstanceAgentPluginConfigDetailsDesiredStateDisabled},
				{Name: common.String("Bastion"), DesiredState: core.InstanceAgentPluginConfigDetailsDesiredStateDisabled},
			},
			IsMonitoringDisabled: common.Bool(false),
			IsManagementDisabled: common.Bool(false),
		},
		// DefinedTags:  map[string]map[string]interface{}{},
		// FreeformTags: map[string]string{},
		InstanceOptions: &core.InstanceOptions{
			AreLegacyImdsEndpointsDisabled: common.Bool(false),
		},
		AvailabilityConfig: &core.LaunchInstanceAvailabilityConfigDetails{
			RecoveryAction: core.LaunchInstanceAvailabilityConfigDetailsRecoveryActionRestoreInstance,
		},
		ShapeConfig: &core.LaunchInstanceShapeConfigDetails{
			Ocpus:       common.Float32(float32(ocpusFloat)),
			MemoryInGBs: common.Float32(float32(memoryInGBsFloat)),
		},
	}

	launchInstanceRequest := core.LaunchInstanceRequest{
		LaunchInstanceDetails: launchDetails,
	}

	response, err := computeClient.LaunchInstance(context.Background(), launchInstanceRequest)
	if err != nil {
		return nil, err
	}

	return &response.Instance, nil
}
