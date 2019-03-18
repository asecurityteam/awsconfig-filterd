package v1

import (
	"context"
	"fmt"
	"io/ioutil"
	"testing"

	domain "github.com/asecurityteam/awsconfig-filterd/pkg/domain"
	"github.com/asecurityteam/logevent"
	"github.com/asecurityteam/runhttp"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

var validEvent = "{\"configurationItemDiff\":{\"changedProperties\":{},\"changeType\":\"CREATE\"},\"configurationItem\":{\"relatedEvents\":[],\"relationships\":[{\"resourceId\":\"eni-05721fa8354d07b8c\",\"resourceName\":null,\"resourceType\":\"AWS::EC2::NetworkInterface\",\"name\":\"Contains NetworkInterface\"},{\"resourceId\":\"sg-05d1b37a375dcca8e\",\"resourceName\":null,\"resourceType\":\"AWS::EC2::SecurityGroup\",\"name\":\"Is associated with SecurityGroup\"},{\"resourceId\":\"subnet-3d0b8c5a\",\"resourceName\":null,\"resourceType\":\"AWS::EC2::Subnet\",\"name\":\"Is contained in Subnet\"},{\"resourceId\":\"vol-0da7faa5400c54c4c\",\"resourceName\":null,\"resourceType\":\"AWS::EC2::Volume\",\"name\":\"Is attached to Volume\"},{\"resourceId\":\"vpc-b290fcd5\",\"resourceName\":null,\"resourceType\":\"AWS::EC2::VPC\",\"name\":\"Is contained in Vpc\"}],\"configuration\":{\"amiLaunchIndex\":0,\"imageId\":\"ami-0bbe6b35405ecebdb\",\"instanceId\":\"i-0a763ac3ee37d8d2b\",\"instanceType\":\"t2.micro\",\"kernelId\":null,\"keyName\":\"zactest2\",\"launchTime\":\"2019-02-22T20:30:10.000Z\",\"monitoring\":{\"state\":\"disabled\"},\"placement\":{\"availabilityZone\":\"us-west-2a\",\"affinity\":null,\"groupName\":\"\",\"partitionNumber\":null,\"hostId\":null,\"tenancy\":\"default\",\"spreadDomain\":null},\"platform\":null,\"privateDnsName\":\"ip-172-31-30-79.us-west-2.compute.internal\",\"privateIpAddress\":\"172.31.30.79\",\"productCodes\":[],\"publicDnsName\":\"ec2-34-222-120-66.us-west-2.compute.amazonaws.com\",\"publicIpAddress\":\"34.222.120.66\",\"ramdiskId\":null,\"state\":{\"code\":16,\"name\":\"running\"},\"stateTransitionReason\":\"\",\"subnetId\":\"subnet-3d0b8c5a\",\"vpcId\":\"vpc-b290fcd5\",\"architecture\":\"x86_64\",\"blockDeviceMappings\":[{\"deviceName\":\"/dev/sda1\",\"ebs\":{\"attachTime\":\"2019-02-22T20:30:11.000Z\",\"deleteOnTermination\":true,\"status\":\"attached\",\"volumeId\":\"vol-0da7faa5400c54c4c\"}}],\"clientToken\":\"\",\"ebsOptimized\":false,\"enaSupport\":true,\"hypervisor\":\"xen\",\"iamInstanceProfile\":null,\"instanceLifecycle\":null,\"elasticGpuAssociations\":[],\"elasticInferenceAcceleratorAssociations\":[],\"networkInterfaces\":[{\"association\":{\"ipOwnerId\":\"amazon\",\"publicDnsName\":\"ec2-34-222-120-66.us-west-2.compute.amazonaws.com\",\"publicIp\":\"34.222.120.66\"},\"attachment\":{\"attachTime\":\"2019-02-22T20:30:10.000Z\",\"attachmentId\":\"eni-attach-0195433bad822bc2f\",\"deleteOnTermination\":true,\"deviceIndex\":0,\"status\":\"attached\"},\"description\":\"\",\"groups\":[{\"groupName\":\"launch-wizard-2\",\"groupId\":\"sg-05d1b37a375dcca8e\"}],\"ipv6Addresses\":[],\"macAddress\":\"02:17:59:ed:8b:0a\",\"networkInterfaceId\":\"eni-05721fa8354d07b8c\",\"ownerId\":\"515665915980\",\"privateDnsName\":\"ip-172-31-30-79.us-west-2.compute.internal\",\"privateIpAddress\":\"172.31.30.79\",\"privateIpAddresses\":[{\"association\":{\"ipOwnerId\":\"amazon\",\"publicDnsName\":\"ec2-34-222-120-66.us-west-2.compute.amazonaws.com\",\"publicIp\":\"34.222.120.66\"},\"primary\":true,\"privateDnsName\":\"ip-172-31-30-79.us-west-2.compute.internal\",\"privateIpAddress\":\"172.31.30.79\"}],\"sourceDestCheck\":true,\"status\":\"in-use\",\"subnetId\":\"subnet-3d0b8c5a\",\"vpcId\":\"vpc-b290fcd5\"}],\"rootDeviceName\":\"/dev/sda1\",\"rootDeviceType\":\"ebs\",\"securityGroups\":[{\"groupName\":\"launch-wizard-2\",\"groupId\":\"sg-05d1b37a375dcca8e\"}],\"sourceDestCheck\":true,\"spotInstanceRequestId\":null,\"sriovNetSupport\":null,\"stateReason\":null,\"tags\":[{\"key\":\"business_unit\",\"value\":\"CISO-Security\"},{\"key\":\"service_name\",\"value\":\"foo-bar\"}],\"virtualizationType\":\"hvm\",\"cpuOptions\":{\"coreCount\":1,\"threadsPerCore\":1},\"capacityReservationId\":null,\"capacityReservationSpecification\":null,\"hibernationOptions\":{\"configured\":false},\"licenses\":[]},\"supplementaryConfiguration\":{},\"tags\":{\"service_name\":\"foo-bar\",\"business_unit\":\"CISO-Security\"},\"configurationItemVersion\":\"1.3\",\"configurationItemCaptureTime\":\"2019-02-22T20:43:10.208Z\",\"configurationStateId\":1550868190208,\"awsAccountId\":\"515665915980\",\"configurationItemStatus\":\"ResourceDiscovered\",\"resourceType\":\"AWS::EC2::Instance\",\"resourceId\":\"i-0a763ac3ee37d8d2b\",\"resourceName\":null,\"ARN\":\"arn:aws:ec2:us-west-2:515665915980:instance/i-0a763ac3ee37d8d2b\",\"awsRegion\":\"us-west-2\",\"availabilityZone\":\"us-west-2a\",\"configurationStateMd5Hash\":\"\",\"resourceCreationTime\":\"2019-02-22T20:30:10.000Z\"},\"notificationCreationTime\":\"2019-02-22T20:43:11.256Z\",\"messageType\":\"ConfigurationItemChangeNotification\",\"recordVersion\":\"1.3\"}"
var noValidResourceType = "{\"configurationItemDiff\":{\"changedProperties\":{},\"changeType\":\"CREATE\"},\"configurationItem\":{\"relatedEvents\":[],\"relationships\":[{\"resourceId\":\"eni-05721fa8354d07b8c\",\"resourceName\":null,\"resourceType\":\"AWS::EC2::NetworkInterface\",\"name\":\"Contains NetworkInterface\"},{\"resourceId\":\"sg-05d1b37a375dcca8e\",\"resourceName\":null,\"resourceType\":\"AWS::EC2::SecurityGroup\",\"name\":\"Is associated with SecurityGroup\"},{\"resourceId\":\"subnet-3d0b8c5a\",\"resourceName\":null,\"resourceType\":\"AWS::EC2::Subnet\",\"name\":\"Is contained in Subnet\"},{\"resourceId\":\"vol-0da7faa5400c54c4c\",\"resourceName\":null,\"resourceType\":\"AWS::EC2::Volume\",\"name\":\"Is attached to Volume\"},{\"resourceId\":\"vpc-b290fcd5\",\"resourceName\":null,\"resourceType\":\"AWS::EC2::VPC\",\"name\":\"Is contained in Vpc\"}],\"configuration\":{\"amiLaunchIndex\":0,\"imageId\":\"ami-0bbe6b35405ecebdb\",\"instanceId\":\"i-0a763ac3ee37d8d2b\",\"instanceType\":\"t2.micro\",\"kernelId\":null,\"keyName\":\"zactest2\",\"launchTime\":\"2019-02-22T20:30:10.000Z\",\"monitoring\":{\"state\":\"disabled\"},\"placement\":{\"availabilityZone\":\"us-west-2a\",\"affinity\":null,\"groupName\":\"\",\"partitionNumber\":null,\"hostId\":null,\"tenancy\":\"default\",\"spreadDomain\":null},\"platform\":null,\"privateDnsName\":\"ip-172-31-30-79.us-west-2.compute.internal\",\"privateIpAddress\":\"172.31.30.79\",\"productCodes\":[],\"publicDnsName\":\"ec2-34-222-120-66.us-west-2.compute.amazonaws.com\",\"publicIpAddress\":\"34.222.120.66\",\"ramdiskId\":null,\"state\":{\"code\":16,\"name\":\"running\"},\"stateTransitionReason\":\"\",\"subnetId\":\"subnet-3d0b8c5a\",\"vpcId\":\"vpc-b290fcd5\",\"architecture\":\"x86_64\",\"blockDeviceMappings\":[{\"deviceName\":\"/dev/sda1\",\"ebs\":{\"attachTime\":\"2019-02-22T20:30:11.000Z\",\"deleteOnTermination\":true,\"status\":\"attached\",\"volumeId\":\"vol-0da7faa5400c54c4c\"}}],\"clientToken\":\"\",\"ebsOptimized\":false,\"enaSupport\":true,\"hypervisor\":\"xen\",\"iamInstanceProfile\":null,\"instanceLifecycle\":null,\"elasticGpuAssociations\":[],\"elasticInferenceAcceleratorAssociations\":[],\"networkInterfaces\":[{\"association\":{\"ipOwnerId\":\"amazon\",\"publicDnsName\":\"ec2-34-222-120-66.us-west-2.compute.amazonaws.com\",\"publicIp\":\"34.222.120.66\"},\"attachment\":{\"attachTime\":\"2019-02-22T20:30:10.000Z\",\"attachmentId\":\"eni-attach-0195433bad822bc2f\",\"deleteOnTermination\":true,\"deviceIndex\":0,\"status\":\"attached\"},\"description\":\"\",\"groups\":[{\"groupName\":\"launch-wizard-2\",\"groupId\":\"sg-05d1b37a375dcca8e\"}],\"ipv6Addresses\":[],\"macAddress\":\"02:17:59:ed:8b:0a\",\"networkInterfaceId\":\"eni-05721fa8354d07b8c\",\"ownerId\":\"515665915980\",\"privateDnsName\":\"ip-172-31-30-79.us-west-2.compute.internal\",\"privateIpAddress\":\"172.31.30.79\",\"privateIpAddresses\":[{\"association\":{\"ipOwnerId\":\"amazon\",\"publicDnsName\":\"ec2-34-222-120-66.us-west-2.compute.amazonaws.com\",\"publicIp\":\"34.222.120.66\"},\"primary\":true,\"privateDnsName\":\"ip-172-31-30-79.us-west-2.compute.internal\",\"privateIpAddress\":\"172.31.30.79\"}],\"sourceDestCheck\":true,\"status\":\"in-use\",\"subnetId\":\"subnet-3d0b8c5a\",\"vpcId\":\"vpc-b290fcd5\"}],\"rootDeviceName\":\"/dev/sda1\",\"rootDeviceType\":\"ebs\",\"securityGroups\":[{\"groupName\":\"launch-wizard-2\",\"groupId\":\"sg-05d1b37a375dcca8e\"}],\"sourceDestCheck\":true,\"spotInstanceRequestId\":null,\"sriovNetSupport\":null,\"stateReason\":null,\"tags\":[{\"key\":\"business_unit\",\"value\":\"CISO-Security\"},{\"key\":\"service_name\",\"value\":\"foo-bar\"}],\"virtualizationType\":\"hvm\",\"cpuOptions\":{\"coreCount\":1,\"threadsPerCore\":1},\"capacityReservationId\":null,\"capacityReservationSpecification\":null,\"hibernationOptions\":{\"configured\":false},\"licenses\":[]},\"supplementaryConfiguration\":{},\"tags\":{\"service_name\":\"foo-bar\",\"business_unit\":\"CISO-Security\"},\"configurationItemVersion\":\"1.3\",\"configurationItemCaptureTime\":\"2019-02-22T20:43:10.208Z\",\"configurationStateId\":1550868190208,\"awsAccountId\":\"515665915980\",\"configurationItemStatus\":\"ResourceDiscovered\",\"resourceType\":\"AWS::EC2::NotAnInstance\",\"resourceId\":\"i-0a763ac3ee37d8d2b\",\"resourceName\":null,\"ARN\":\"arn:aws:ec2:us-west-2:515665915980:instance/i-0a763ac3ee37d8d2b\",\"awsRegion\":\"us-west-2\",\"availabilityZone\":\"us-west-2a\",\"configurationStateMd5Hash\":\"\",\"resourceCreationTime\":\"2019-02-22T20:30:10.000Z\"},\"notificationCreationTime\":\"2019-02-22T20:43:11.256Z\",\"messageType\":\"ConfigurationItemChangeNotification\",\"recordVersion\":\"1.3\"}"
var noResourceType = "{\"configurationItemDiff\":{\"changedProperties\":{},\"changeType\":\"CREATE\"},\"configurationItem\":{\"relatedEvents\":[],\"relationships\":[{\"resourceId\":\"eni-05721fa8354d07b8c\",\"resourceName\":null,\"resourceType\":\"AWS::EC2::NetworkInterface\",\"name\":\"Contains NetworkInterface\"},{\"resourceId\":\"sg-05d1b37a375dcca8e\",\"resourceName\":null,\"resourceType\":\"AWS::EC2::SecurityGroup\",\"name\":\"Is associated with SecurityGroup\"},{\"resourceId\":\"subnet-3d0b8c5a\",\"resourceName\":null,\"resourceType\":\"AWS::EC2::Subnet\",\"name\":\"Is contained in Subnet\"},{\"resourceId\":\"vol-0da7faa5400c54c4c\",\"resourceName\":null,\"resourceType\":\"AWS::EC2::Volume\",\"name\":\"Is attached to Volume\"},{\"resourceId\":\"vpc-b290fcd5\",\"resourceName\":null,\"resourceType\":\"AWS::EC2::VPC\",\"name\":\"Is contained in Vpc\"}],\"configuration\":{\"amiLaunchIndex\":0,\"imageId\":\"ami-0bbe6b35405ecebdb\",\"instanceId\":\"i-0a763ac3ee37d8d2b\",\"instanceType\":\"t2.micro\",\"kernelId\":null,\"keyName\":\"zactest2\",\"launchTime\":\"2019-02-22T20:30:10.000Z\",\"monitoring\":{\"state\":\"disabled\"},\"placement\":{\"availabilityZone\":\"us-west-2a\",\"affinity\":null,\"groupName\":\"\",\"partitionNumber\":null,\"hostId\":null,\"tenancy\":\"default\",\"spreadDomain\":null},\"platform\":null,\"privateDnsName\":\"ip-172-31-30-79.us-west-2.compute.internal\",\"privateIpAddress\":\"172.31.30.79\",\"productCodes\":[],\"publicDnsName\":\"ec2-34-222-120-66.us-west-2.compute.amazonaws.com\",\"publicIpAddress\":\"34.222.120.66\",\"ramdiskId\":null,\"state\":{\"code\":16,\"name\":\"running\"},\"stateTransitionReason\":\"\",\"subnetId\":\"subnet-3d0b8c5a\",\"vpcId\":\"vpc-b290fcd5\",\"architecture\":\"x86_64\",\"blockDeviceMappings\":[{\"deviceName\":\"/dev/sda1\",\"ebs\":{\"attachTime\":\"2019-02-22T20:30:11.000Z\",\"deleteOnTermination\":true,\"status\":\"attached\",\"volumeId\":\"vol-0da7faa5400c54c4c\"}}],\"clientToken\":\"\",\"ebsOptimized\":false,\"enaSupport\":true,\"hypervisor\":\"xen\",\"iamInstanceProfile\":null,\"instanceLifecycle\":null,\"elasticGpuAssociations\":[],\"elasticInferenceAcceleratorAssociations\":[],\"networkInterfaces\":[{\"association\":{\"ipOwnerId\":\"amazon\",\"publicDnsName\":\"ec2-34-222-120-66.us-west-2.compute.amazonaws.com\",\"publicIp\":\"34.222.120.66\"},\"attachment\":{\"attachTime\":\"2019-02-22T20:30:10.000Z\",\"attachmentId\":\"eni-attach-0195433bad822bc2f\",\"deleteOnTermination\":true,\"deviceIndex\":0,\"status\":\"attached\"},\"description\":\"\",\"groups\":[{\"groupName\":\"launch-wizard-2\",\"groupId\":\"sg-05d1b37a375dcca8e\"}],\"ipv6Addresses\":[],\"macAddress\":\"02:17:59:ed:8b:0a\",\"networkInterfaceId\":\"eni-05721fa8354d07b8c\",\"ownerId\":\"515665915980\",\"privateDnsName\":\"ip-172-31-30-79.us-west-2.compute.internal\",\"privateIpAddress\":\"172.31.30.79\",\"privateIpAddresses\":[{\"association\":{\"ipOwnerId\":\"amazon\",\"publicDnsName\":\"ec2-34-222-120-66.us-west-2.compute.amazonaws.com\",\"publicIp\":\"34.222.120.66\"},\"primary\":true,\"privateDnsName\":\"ip-172-31-30-79.us-west-2.compute.internal\",\"privateIpAddress\":\"172.31.30.79\"}],\"sourceDestCheck\":true,\"status\":\"in-use\",\"subnetId\":\"subnet-3d0b8c5a\",\"vpcId\":\"vpc-b290fcd5\"}],\"rootDeviceName\":\"/dev/sda1\",\"rootDeviceType\":\"ebs\",\"securityGroups\":[{\"groupName\":\"launch-wizard-2\",\"groupId\":\"sg-05d1b37a375dcca8e\"}],\"sourceDestCheck\":true,\"spotInstanceRequestId\":null,\"sriovNetSupport\":null,\"stateReason\":null,\"tags\":[{\"key\":\"business_unit\",\"value\":\"CISO-Security\"},{\"key\":\"service_name\",\"value\":\"foo-bar\"}],\"virtualizationType\":\"hvm\",\"cpuOptions\":{\"coreCount\":1,\"threadsPerCore\":1},\"capacityReservationId\":null,\"capacityReservationSpecification\":null,\"hibernationOptions\":{\"configured\":false},\"licenses\":[]},\"supplementaryConfiguration\":{},\"tags\":{\"service_name\":\"foo-bar\",\"business_unit\":\"CISO-Security\"},\"configurationItemVersion\":\"1.3\",\"configurationItemCaptureTime\":\"2019-02-22T20:43:10.208Z\",\"configurationStateId\":1550868190208,\"awsAccountId\":\"515665915980\",\"configurationItemStatus\":\"ResourceDiscovered\",\"resourceType\":\"\",\"resourceId\":\"i-0a763ac3ee37d8d2b\",\"resourceName\":null,\"ARN\":\"arn:aws:ec2:us-west-2:515665915980:instance/i-0a763ac3ee37d8d2b\",\"awsRegion\":\"us-west-2\",\"availabilityZone\":\"us-west-2a\",\"configurationStateMd5Hash\":\"\",\"resourceCreationTime\":\"2019-02-22T20:30:10.000Z\"},\"notificationCreationTime\":\"2019-02-22T20:43:11.256Z\",\"messageType\":\"ConfigurationItemChangeNotification\",\"recordVersion\":\"1.3\"}"
var timestamp = "2019-02-22T20:43:11.479Z"

func TestHandle(t *testing.T) {
	tc := []struct {
		name         string
		event        string
		err          error
		filterCalled bool
		filterOK     bool
		filterErr    error
	}{
		{
			name:         "success",
			event:        validEvent,
			err:          nil,
			filterCalled: true,
			filterOK:     true,
			filterErr:    nil,
		},
		{
			name:         "no valid resource type",
			event:        noValidResourceType,
			err:          domain.ErrEventDiscarded{Reason: "no valid resource type found"},
			filterCalled: true,
			filterOK:     false,
			filterErr:    fmt.Errorf("no valid resource type found"),
		},
		{
			name:         "no valid resource type",
			event:        noResourceType,
			err:          domain.ErrEventDiscarded{Reason: "empty resource type"},
			filterCalled: false,
			filterOK:     false,
			filterErr:    nil,
		},
		{
			name:         "cannot unmarshal ConfigEvent",
			event:        "0",
			err:          domain.ErrEventDiscarded{Reason: "json: cannot unmarshal number into Go value of type v1.ConfigEvent"},
			filterCalled: false,
			filterOK:     false,
			filterErr:    nil,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockFilter := NewMockConfigFilter(ctrl)
			if tt.filterCalled {
				mockFilter.EXPECT().Filter(gomock.Any()).Return(tt.filterOK, tt.filterErr)
			}

			input := ConfigNotification{
				Message:   tt.event,
				Timestamp: timestamp,
			}

			configFilterHandler := &ConfigFilterHandler{
				LogFn:  runhttp.LoggerFromContext,
				StatFn: runhttp.StatFromContext,
				Filter: mockFilter,
			}

			ctx := logevent.NewContext(context.Background(), logevent.New(logevent.Config{Output: ioutil.Discard}))
			_, err := configFilterHandler.Handle(ctx, input)
			require.Equal(t, tt.err, err)
		})
	}
}
