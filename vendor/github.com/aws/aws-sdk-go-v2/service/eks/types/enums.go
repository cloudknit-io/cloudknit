// Code generated by smithy-go-codegen DO NOT EDIT.

package types

type AddonIssueCode string

// Enum values for AddonIssueCode
const (
	AddonIssueCodeAccessDenied                 AddonIssueCode = "AccessDenied"
	AddonIssueCodeInternalFailure              AddonIssueCode = "InternalFailure"
	AddonIssueCodeClusterUnreachable           AddonIssueCode = "ClusterUnreachable"
	AddonIssueCodeInsufficientNumberOfReplicas AddonIssueCode = "InsufficientNumberOfReplicas"
	AddonIssueCodeConfigurationConflict        AddonIssueCode = "ConfigurationConflict"
	AddonIssueCodeAdmissionRequestDenied       AddonIssueCode = "AdmissionRequestDenied"
	AddonIssueCodeUnsupportedAddonModification AddonIssueCode = "UnsupportedAddonModification"
	AddonIssueCodeK8sResourceNotFound          AddonIssueCode = "K8sResourceNotFound"
)

// Values returns all known values for AddonIssueCode. Note that this can be
// expanded in the future, and so it is only as up to date as the client. The
// ordering of this slice is not guaranteed to be stable across updates.
func (AddonIssueCode) Values() []AddonIssueCode {
	return []AddonIssueCode{
		"AccessDenied",
		"InternalFailure",
		"ClusterUnreachable",
		"InsufficientNumberOfReplicas",
		"ConfigurationConflict",
		"AdmissionRequestDenied",
		"UnsupportedAddonModification",
		"K8sResourceNotFound",
	}
}

type AddonStatus string

// Enum values for AddonStatus
const (
	AddonStatusCreating     AddonStatus = "CREATING"
	AddonStatusActive       AddonStatus = "ACTIVE"
	AddonStatusCreateFailed AddonStatus = "CREATE_FAILED"
	AddonStatusUpdating     AddonStatus = "UPDATING"
	AddonStatusDeleting     AddonStatus = "DELETING"
	AddonStatusDeleteFailed AddonStatus = "DELETE_FAILED"
	AddonStatusDegraded     AddonStatus = "DEGRADED"
)

// Values returns all known values for AddonStatus. Note that this can be expanded
// in the future, and so it is only as up to date as the client. The ordering of
// this slice is not guaranteed to be stable across updates.
func (AddonStatus) Values() []AddonStatus {
	return []AddonStatus{
		"CREATING",
		"ACTIVE",
		"CREATE_FAILED",
		"UPDATING",
		"DELETING",
		"DELETE_FAILED",
		"DEGRADED",
	}
}

type AMITypes string

// Enum values for AMITypes
const (
	AMITypesAl2X8664          AMITypes = "AL2_x86_64"
	AMITypesAl2X8664Gpu       AMITypes = "AL2_x86_64_GPU"
	AMITypesAl2Arm64          AMITypes = "AL2_ARM_64"
	AMITypesCustom            AMITypes = "CUSTOM"
	AMITypesBottlerocketArm64 AMITypes = "BOTTLEROCKET_ARM_64"
	AMITypesBottlerocketX8664 AMITypes = "BOTTLEROCKET_x86_64"
)

// Values returns all known values for AMITypes. Note that this can be expanded in
// the future, and so it is only as up to date as the client. The ordering of this
// slice is not guaranteed to be stable across updates.
func (AMITypes) Values() []AMITypes {
	return []AMITypes{
		"AL2_x86_64",
		"AL2_x86_64_GPU",
		"AL2_ARM_64",
		"CUSTOM",
		"BOTTLEROCKET_ARM_64",
		"BOTTLEROCKET_x86_64",
	}
}

type CapacityTypes string

// Enum values for CapacityTypes
const (
	CapacityTypesOnDemand CapacityTypes = "ON_DEMAND"
	CapacityTypesSpot     CapacityTypes = "SPOT"
)

// Values returns all known values for CapacityTypes. Note that this can be
// expanded in the future, and so it is only as up to date as the client. The
// ordering of this slice is not guaranteed to be stable across updates.
func (CapacityTypes) Values() []CapacityTypes {
	return []CapacityTypes{
		"ON_DEMAND",
		"SPOT",
	}
}

type ClusterStatus string

// Enum values for ClusterStatus
const (
	ClusterStatusCreating ClusterStatus = "CREATING"
	ClusterStatusActive   ClusterStatus = "ACTIVE"
	ClusterStatusDeleting ClusterStatus = "DELETING"
	ClusterStatusFailed   ClusterStatus = "FAILED"
	ClusterStatusUpdating ClusterStatus = "UPDATING"
	ClusterStatusPending  ClusterStatus = "PENDING"
)

// Values returns all known values for ClusterStatus. Note that this can be
// expanded in the future, and so it is only as up to date as the client. The
// ordering of this slice is not guaranteed to be stable across updates.
func (ClusterStatus) Values() []ClusterStatus {
	return []ClusterStatus{
		"CREATING",
		"ACTIVE",
		"DELETING",
		"FAILED",
		"UPDATING",
		"PENDING",
	}
}

type ConfigStatus string

// Enum values for ConfigStatus
const (
	ConfigStatusCreating ConfigStatus = "CREATING"
	ConfigStatusDeleting ConfigStatus = "DELETING"
	ConfigStatusActive   ConfigStatus = "ACTIVE"
)

// Values returns all known values for ConfigStatus. Note that this can be expanded
// in the future, and so it is only as up to date as the client. The ordering of
// this slice is not guaranteed to be stable across updates.
func (ConfigStatus) Values() []ConfigStatus {
	return []ConfigStatus{
		"CREATING",
		"DELETING",
		"ACTIVE",
	}
}

type ConnectorConfigProvider string

// Enum values for ConnectorConfigProvider
const (
	ConnectorConfigProviderEksAnywhere ConnectorConfigProvider = "EKS_ANYWHERE"
	ConnectorConfigProviderAnthos      ConnectorConfigProvider = "ANTHOS"
	ConnectorConfigProviderGke         ConnectorConfigProvider = "GKE"
	ConnectorConfigProviderAks         ConnectorConfigProvider = "AKS"
	ConnectorConfigProviderOpenshift   ConnectorConfigProvider = "OPENSHIFT"
	ConnectorConfigProviderTanzu       ConnectorConfigProvider = "TANZU"
	ConnectorConfigProviderRancher     ConnectorConfigProvider = "RANCHER"
	ConnectorConfigProviderEc2         ConnectorConfigProvider = "EC2"
	ConnectorConfigProviderOther       ConnectorConfigProvider = "OTHER"
)

// Values returns all known values for ConnectorConfigProvider. Note that this can
// be expanded in the future, and so it is only as up to date as the client. The
// ordering of this slice is not guaranteed to be stable across updates.
func (ConnectorConfigProvider) Values() []ConnectorConfigProvider {
	return []ConnectorConfigProvider{
		"EKS_ANYWHERE",
		"ANTHOS",
		"GKE",
		"AKS",
		"OPENSHIFT",
		"TANZU",
		"RANCHER",
		"EC2",
		"OTHER",
	}
}

type ErrorCode string

// Enum values for ErrorCode
const (
	ErrorCodeSubnetNotFound               ErrorCode = "SubnetNotFound"
	ErrorCodeSecurityGroupNotFound        ErrorCode = "SecurityGroupNotFound"
	ErrorCodeEniLimitReached              ErrorCode = "EniLimitReached"
	ErrorCodeIpNotAvailable               ErrorCode = "IpNotAvailable"
	ErrorCodeAccessDenied                 ErrorCode = "AccessDenied"
	ErrorCodeOperationNotPermitted        ErrorCode = "OperationNotPermitted"
	ErrorCodeVpcIdNotFound                ErrorCode = "VpcIdNotFound"
	ErrorCodeUnknown                      ErrorCode = "Unknown"
	ErrorCodeNodeCreationFailure          ErrorCode = "NodeCreationFailure"
	ErrorCodePodEvictionFailure           ErrorCode = "PodEvictionFailure"
	ErrorCodeInsufficientFreeAddresses    ErrorCode = "InsufficientFreeAddresses"
	ErrorCodeClusterUnreachable           ErrorCode = "ClusterUnreachable"
	ErrorCodeInsufficientNumberOfReplicas ErrorCode = "InsufficientNumberOfReplicas"
	ErrorCodeConfigurationConflict        ErrorCode = "ConfigurationConflict"
	ErrorCodeAdmissionRequestDenied       ErrorCode = "AdmissionRequestDenied"
	ErrorCodeUnsupportedAddonModification ErrorCode = "UnsupportedAddonModification"
	ErrorCodeK8sResourceNotFound          ErrorCode = "K8sResourceNotFound"
)

// Values returns all known values for ErrorCode. Note that this can be expanded in
// the future, and so it is only as up to date as the client. The ordering of this
// slice is not guaranteed to be stable across updates.
func (ErrorCode) Values() []ErrorCode {
	return []ErrorCode{
		"SubnetNotFound",
		"SecurityGroupNotFound",
		"EniLimitReached",
		"IpNotAvailable",
		"AccessDenied",
		"OperationNotPermitted",
		"VpcIdNotFound",
		"Unknown",
		"NodeCreationFailure",
		"PodEvictionFailure",
		"InsufficientFreeAddresses",
		"ClusterUnreachable",
		"InsufficientNumberOfReplicas",
		"ConfigurationConflict",
		"AdmissionRequestDenied",
		"UnsupportedAddonModification",
		"K8sResourceNotFound",
	}
}

type FargateProfileStatus string

// Enum values for FargateProfileStatus
const (
	FargateProfileStatusCreating     FargateProfileStatus = "CREATING"
	FargateProfileStatusActive       FargateProfileStatus = "ACTIVE"
	FargateProfileStatusDeleting     FargateProfileStatus = "DELETING"
	FargateProfileStatusCreateFailed FargateProfileStatus = "CREATE_FAILED"
	FargateProfileStatusDeleteFailed FargateProfileStatus = "DELETE_FAILED"
)

// Values returns all known values for FargateProfileStatus. Note that this can be
// expanded in the future, and so it is only as up to date as the client. The
// ordering of this slice is not guaranteed to be stable across updates.
func (FargateProfileStatus) Values() []FargateProfileStatus {
	return []FargateProfileStatus{
		"CREATING",
		"ACTIVE",
		"DELETING",
		"CREATE_FAILED",
		"DELETE_FAILED",
	}
}

type IpFamily string

// Enum values for IpFamily
const (
	IpFamilyIpv4 IpFamily = "ipv4"
	IpFamilyIpv6 IpFamily = "ipv6"
)

// Values returns all known values for IpFamily. Note that this can be expanded in
// the future, and so it is only as up to date as the client. The ordering of this
// slice is not guaranteed to be stable across updates.
func (IpFamily) Values() []IpFamily {
	return []IpFamily{
		"ipv4",
		"ipv6",
	}
}

type LogType string

// Enum values for LogType
const (
	LogTypeApi               LogType = "api"
	LogTypeAudit             LogType = "audit"
	LogTypeAuthenticator     LogType = "authenticator"
	LogTypeControllerManager LogType = "controllerManager"
	LogTypeScheduler         LogType = "scheduler"
)

// Values returns all known values for LogType. Note that this can be expanded in
// the future, and so it is only as up to date as the client. The ordering of this
// slice is not guaranteed to be stable across updates.
func (LogType) Values() []LogType {
	return []LogType{
		"api",
		"audit",
		"authenticator",
		"controllerManager",
		"scheduler",
	}
}

type NodegroupIssueCode string

// Enum values for NodegroupIssueCode
const (
	NodegroupIssueCodeAutoScalingGroupNotFound             NodegroupIssueCode = "AutoScalingGroupNotFound"
	NodegroupIssueCodeAutoScalingGroupInvalidConfiguration NodegroupIssueCode = "AutoScalingGroupInvalidConfiguration"
	NodegroupIssueCodeEc2SecurityGroupNotFound             NodegroupIssueCode = "Ec2SecurityGroupNotFound"
	NodegroupIssueCodeEc2SecurityGroupDeletionFailure      NodegroupIssueCode = "Ec2SecurityGroupDeletionFailure"
	NodegroupIssueCodeEc2LaunchTemplateNotFound            NodegroupIssueCode = "Ec2LaunchTemplateNotFound"
	NodegroupIssueCodeEc2LaunchTemplateVersionMismatch     NodegroupIssueCode = "Ec2LaunchTemplateVersionMismatch"
	NodegroupIssueCodeEc2SubnetNotFound                    NodegroupIssueCode = "Ec2SubnetNotFound"
	NodegroupIssueCodeEc2SubnetInvalidConfiguration        NodegroupIssueCode = "Ec2SubnetInvalidConfiguration"
	NodegroupIssueCodeIamInstanceProfileNotFound           NodegroupIssueCode = "IamInstanceProfileNotFound"
	NodegroupIssueCodeIamLimitExceeded                     NodegroupIssueCode = "IamLimitExceeded"
	NodegroupIssueCodeIamNodeRoleNotFound                  NodegroupIssueCode = "IamNodeRoleNotFound"
	NodegroupIssueCodeNodeCreationFailure                  NodegroupIssueCode = "NodeCreationFailure"
	NodegroupIssueCodeAsgInstanceLaunchFailures            NodegroupIssueCode = "AsgInstanceLaunchFailures"
	NodegroupIssueCodeInstanceLimitExceeded                NodegroupIssueCode = "InstanceLimitExceeded"
	NodegroupIssueCodeInsufficientFreeAddresses            NodegroupIssueCode = "InsufficientFreeAddresses"
	NodegroupIssueCodeAccessDenied                         NodegroupIssueCode = "AccessDenied"
	NodegroupIssueCodeInternalFailure                      NodegroupIssueCode = "InternalFailure"
	NodegroupIssueCodeClusterUnreachable                   NodegroupIssueCode = "ClusterUnreachable"
)

// Values returns all known values for NodegroupIssueCode. Note that this can be
// expanded in the future, and so it is only as up to date as the client. The
// ordering of this slice is not guaranteed to be stable across updates.
func (NodegroupIssueCode) Values() []NodegroupIssueCode {
	return []NodegroupIssueCode{
		"AutoScalingGroupNotFound",
		"AutoScalingGroupInvalidConfiguration",
		"Ec2SecurityGroupNotFound",
		"Ec2SecurityGroupDeletionFailure",
		"Ec2LaunchTemplateNotFound",
		"Ec2LaunchTemplateVersionMismatch",
		"Ec2SubnetNotFound",
		"Ec2SubnetInvalidConfiguration",
		"IamInstanceProfileNotFound",
		"IamLimitExceeded",
		"IamNodeRoleNotFound",
		"NodeCreationFailure",
		"AsgInstanceLaunchFailures",
		"InstanceLimitExceeded",
		"InsufficientFreeAddresses",
		"AccessDenied",
		"InternalFailure",
		"ClusterUnreachable",
	}
}

type NodegroupStatus string

// Enum values for NodegroupStatus
const (
	NodegroupStatusCreating     NodegroupStatus = "CREATING"
	NodegroupStatusActive       NodegroupStatus = "ACTIVE"
	NodegroupStatusUpdating     NodegroupStatus = "UPDATING"
	NodegroupStatusDeleting     NodegroupStatus = "DELETING"
	NodegroupStatusCreateFailed NodegroupStatus = "CREATE_FAILED"
	NodegroupStatusDeleteFailed NodegroupStatus = "DELETE_FAILED"
	NodegroupStatusDegraded     NodegroupStatus = "DEGRADED"
)

// Values returns all known values for NodegroupStatus. Note that this can be
// expanded in the future, and so it is only as up to date as the client. The
// ordering of this slice is not guaranteed to be stable across updates.
func (NodegroupStatus) Values() []NodegroupStatus {
	return []NodegroupStatus{
		"CREATING",
		"ACTIVE",
		"UPDATING",
		"DELETING",
		"CREATE_FAILED",
		"DELETE_FAILED",
		"DEGRADED",
	}
}

type ResolveConflicts string

// Enum values for ResolveConflicts
const (
	ResolveConflictsOverwrite ResolveConflicts = "OVERWRITE"
	ResolveConflictsNone      ResolveConflicts = "NONE"
)

// Values returns all known values for ResolveConflicts. Note that this can be
// expanded in the future, and so it is only as up to date as the client. The
// ordering of this slice is not guaranteed to be stable across updates.
func (ResolveConflicts) Values() []ResolveConflicts {
	return []ResolveConflicts{
		"OVERWRITE",
		"NONE",
	}
}

type TaintEffect string

// Enum values for TaintEffect
const (
	TaintEffectNoSchedule       TaintEffect = "NO_SCHEDULE"
	TaintEffectNoExecute        TaintEffect = "NO_EXECUTE"
	TaintEffectPreferNoSchedule TaintEffect = "PREFER_NO_SCHEDULE"
)

// Values returns all known values for TaintEffect. Note that this can be expanded
// in the future, and so it is only as up to date as the client. The ordering of
// this slice is not guaranteed to be stable across updates.
func (TaintEffect) Values() []TaintEffect {
	return []TaintEffect{
		"NO_SCHEDULE",
		"NO_EXECUTE",
		"PREFER_NO_SCHEDULE",
	}
}

type UpdateParamType string

// Enum values for UpdateParamType
const (
	UpdateParamTypeVersion                  UpdateParamType = "Version"
	UpdateParamTypePlatformVersion          UpdateParamType = "PlatformVersion"
	UpdateParamTypeEndpointPrivateAccess    UpdateParamType = "EndpointPrivateAccess"
	UpdateParamTypeEndpointPublicAccess     UpdateParamType = "EndpointPublicAccess"
	UpdateParamTypeClusterLogging           UpdateParamType = "ClusterLogging"
	UpdateParamTypeDesiredSize              UpdateParamType = "DesiredSize"
	UpdateParamTypeLabelsToAdd              UpdateParamType = "LabelsToAdd"
	UpdateParamTypeLabelsToRemove           UpdateParamType = "LabelsToRemove"
	UpdateParamTypeTaintsToAdd              UpdateParamType = "TaintsToAdd"
	UpdateParamTypeTaintsToRemove           UpdateParamType = "TaintsToRemove"
	UpdateParamTypeMaxSize                  UpdateParamType = "MaxSize"
	UpdateParamTypeMinSize                  UpdateParamType = "MinSize"
	UpdateParamTypeReleaseVersion           UpdateParamType = "ReleaseVersion"
	UpdateParamTypePublicAccessCidrs        UpdateParamType = "PublicAccessCidrs"
	UpdateParamTypeLaunchTemplateName       UpdateParamType = "LaunchTemplateName"
	UpdateParamTypeLaunchTemplateVersion    UpdateParamType = "LaunchTemplateVersion"
	UpdateParamTypeIdentityProviderConfig   UpdateParamType = "IdentityProviderConfig"
	UpdateParamTypeEncryptionConfig         UpdateParamType = "EncryptionConfig"
	UpdateParamTypeAddonVersion             UpdateParamType = "AddonVersion"
	UpdateParamTypeServiceAccountRoleArn    UpdateParamType = "ServiceAccountRoleArn"
	UpdateParamTypeResolveConflicts         UpdateParamType = "ResolveConflicts"
	UpdateParamTypeMaxUnavailable           UpdateParamType = "MaxUnavailable"
	UpdateParamTypeMaxUnavailablePercentage UpdateParamType = "MaxUnavailablePercentage"
)

// Values returns all known values for UpdateParamType. Note that this can be
// expanded in the future, and so it is only as up to date as the client. The
// ordering of this slice is not guaranteed to be stable across updates.
func (UpdateParamType) Values() []UpdateParamType {
	return []UpdateParamType{
		"Version",
		"PlatformVersion",
		"EndpointPrivateAccess",
		"EndpointPublicAccess",
		"ClusterLogging",
		"DesiredSize",
		"LabelsToAdd",
		"LabelsToRemove",
		"TaintsToAdd",
		"TaintsToRemove",
		"MaxSize",
		"MinSize",
		"ReleaseVersion",
		"PublicAccessCidrs",
		"LaunchTemplateName",
		"LaunchTemplateVersion",
		"IdentityProviderConfig",
		"EncryptionConfig",
		"AddonVersion",
		"ServiceAccountRoleArn",
		"ResolveConflicts",
		"MaxUnavailable",
		"MaxUnavailablePercentage",
	}
}

type UpdateStatus string

// Enum values for UpdateStatus
const (
	UpdateStatusInProgress UpdateStatus = "InProgress"
	UpdateStatusFailed     UpdateStatus = "Failed"
	UpdateStatusCancelled  UpdateStatus = "Cancelled"
	UpdateStatusSuccessful UpdateStatus = "Successful"
)

// Values returns all known values for UpdateStatus. Note that this can be expanded
// in the future, and so it is only as up to date as the client. The ordering of
// this slice is not guaranteed to be stable across updates.
func (UpdateStatus) Values() []UpdateStatus {
	return []UpdateStatus{
		"InProgress",
		"Failed",
		"Cancelled",
		"Successful",
	}
}

type UpdateType string

// Enum values for UpdateType
const (
	UpdateTypeVersionUpdate                      UpdateType = "VersionUpdate"
	UpdateTypeEndpointAccessUpdate               UpdateType = "EndpointAccessUpdate"
	UpdateTypeLoggingUpdate                      UpdateType = "LoggingUpdate"
	UpdateTypeConfigUpdate                       UpdateType = "ConfigUpdate"
	UpdateTypeAssociateIdentityProviderConfig    UpdateType = "AssociateIdentityProviderConfig"
	UpdateTypeDisassociateIdentityProviderConfig UpdateType = "DisassociateIdentityProviderConfig"
	UpdateTypeAssociateEncryptionConfig          UpdateType = "AssociateEncryptionConfig"
	UpdateTypeAddonUpdate                        UpdateType = "AddonUpdate"
)

// Values returns all known values for UpdateType. Note that this can be expanded
// in the future, and so it is only as up to date as the client. The ordering of
// this slice is not guaranteed to be stable across updates.
func (UpdateType) Values() []UpdateType {
	return []UpdateType{
		"VersionUpdate",
		"EndpointAccessUpdate",
		"LoggingUpdate",
		"ConfigUpdate",
		"AssociateIdentityProviderConfig",
		"DisassociateIdentityProviderConfig",
		"AssociateEncryptionConfig",
		"AddonUpdate",
	}
}
