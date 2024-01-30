// Code generated by mdatagen. DO NOT EDIT.

package metadata

import (
	"go.opentelemetry.io/collector/pdata/pcommon"
)

// ResourceBuilder is a helper struct to build resources predefined in metadata.yaml.
// The ResourceBuilder is not thread-safe and must not to be used in multiple goroutines.
type ResourceBuilder struct {
	config ResourceAttributesConfig
	res    pcommon.Resource
}

// NewResourceBuilder creates a new ResourceBuilder. This method should be called on the start of the application.
func NewResourceBuilder(rac ResourceAttributesConfig) *ResourceBuilder {
	return &ResourceBuilder{
		config: rac,
		res:    pcommon.NewResource(),
	}
}

// SetCloudAccountID sets provided value as "cloud.account.id" attribute.
func (rb *ResourceBuilder) SetCloudAccountID(val string) {
	if rb.config.CloudAccountID.Enabled {
		rb.res.Attributes().PutStr("cloud.account.id", val)
	}
}

// SetCloudAvailabilityZone sets provided value as "cloud.availability_zone" attribute.
func (rb *ResourceBuilder) SetCloudAvailabilityZone(val string) {
	if rb.config.CloudAvailabilityZone.Enabled {
		rb.res.Attributes().PutStr("cloud.availability_zone", val)
	}
}

// SetCloudPlatform sets provided value as "cloud.platform" attribute.
func (rb *ResourceBuilder) SetCloudPlatform(val string) {
	if rb.config.CloudPlatform.Enabled {
		rb.res.Attributes().PutStr("cloud.platform", val)
	}
}

// SetCloudProvider sets provided value as "cloud.provider" attribute.
func (rb *ResourceBuilder) SetCloudProvider(val string) {
	if rb.config.CloudProvider.Enabled {
		rb.res.Attributes().PutStr("cloud.provider", val)
	}
}

// SetCloudRegion sets provided value as "cloud.region" attribute.
func (rb *ResourceBuilder) SetCloudRegion(val string) {
	if rb.config.CloudRegion.Enabled {
		rb.res.Attributes().PutStr("cloud.region", val)
	}
}

// SetHostID sets provided value as "host.id" attribute.
func (rb *ResourceBuilder) SetHostID(val string) {
	if rb.config.HostID.Enabled {
		rb.res.Attributes().PutStr("host.id", val)
	}
}

// SetHostImageID sets provided value as "host.image.id" attribute.
func (rb *ResourceBuilder) SetHostImageID(val string) {
	if rb.config.HostImageID.Enabled {
		rb.res.Attributes().PutStr("host.image.id", val)
	}
}

// SetHostName sets provided value as "host.name" attribute.
func (rb *ResourceBuilder) SetHostName(val string) {
	if rb.config.HostName.Enabled {
		rb.res.Attributes().PutStr("host.name", val)
	}
}

// SetHostType sets provided value as "host.type" attribute.
func (rb *ResourceBuilder) SetHostType(val string) {
	if rb.config.HostType.Enabled {
		rb.res.Attributes().PutStr("host.type", val)
	}
}

// Emit returns the built resource and resets the internal builder state.
func (rb *ResourceBuilder) Emit() pcommon.Resource {
	r := rb.res
	rb.res = pcommon.NewResource()
	return r
}