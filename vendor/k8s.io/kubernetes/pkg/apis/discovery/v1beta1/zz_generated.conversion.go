// +build !ignore_autogenerated

/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by conversion-gen. DO NOT EDIT.

package v1beta1

import (
	unsafe "unsafe"

	v1 "k8s.io/api/core/v1"
	v1beta1 "k8s.io/api/discovery/v1beta1"
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
	core "k8s.io/kubernetes/pkg/apis/core"
	discovery "k8s.io/kubernetes/pkg/apis/discovery"
)

func init() {
	localSchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*v1beta1.EndpointConditions)(nil), (*discovery.EndpointConditions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_EndpointConditions_To_discovery_EndpointConditions(a.(*v1beta1.EndpointConditions), b.(*discovery.EndpointConditions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*discovery.EndpointConditions)(nil), (*v1beta1.EndpointConditions)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_discovery_EndpointConditions_To_v1beta1_EndpointConditions(a.(*discovery.EndpointConditions), b.(*v1beta1.EndpointConditions), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.EndpointHints)(nil), (*discovery.EndpointHints)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_EndpointHints_To_discovery_EndpointHints(a.(*v1beta1.EndpointHints), b.(*discovery.EndpointHints), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*discovery.EndpointHints)(nil), (*v1beta1.EndpointHints)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_discovery_EndpointHints_To_v1beta1_EndpointHints(a.(*discovery.EndpointHints), b.(*v1beta1.EndpointHints), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.EndpointPort)(nil), (*discovery.EndpointPort)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_EndpointPort_To_discovery_EndpointPort(a.(*v1beta1.EndpointPort), b.(*discovery.EndpointPort), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*discovery.EndpointPort)(nil), (*v1beta1.EndpointPort)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_discovery_EndpointPort_To_v1beta1_EndpointPort(a.(*discovery.EndpointPort), b.(*v1beta1.EndpointPort), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.EndpointSlice)(nil), (*discovery.EndpointSlice)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_EndpointSlice_To_discovery_EndpointSlice(a.(*v1beta1.EndpointSlice), b.(*discovery.EndpointSlice), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*discovery.EndpointSlice)(nil), (*v1beta1.EndpointSlice)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_discovery_EndpointSlice_To_v1beta1_EndpointSlice(a.(*discovery.EndpointSlice), b.(*v1beta1.EndpointSlice), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.EndpointSliceList)(nil), (*discovery.EndpointSliceList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_EndpointSliceList_To_discovery_EndpointSliceList(a.(*v1beta1.EndpointSliceList), b.(*discovery.EndpointSliceList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*discovery.EndpointSliceList)(nil), (*v1beta1.EndpointSliceList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_discovery_EndpointSliceList_To_v1beta1_EndpointSliceList(a.(*discovery.EndpointSliceList), b.(*v1beta1.EndpointSliceList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta1.ForZone)(nil), (*discovery.ForZone)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_ForZone_To_discovery_ForZone(a.(*v1beta1.ForZone), b.(*discovery.ForZone), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*discovery.ForZone)(nil), (*v1beta1.ForZone)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_discovery_ForZone_To_v1beta1_ForZone(a.(*discovery.ForZone), b.(*v1beta1.ForZone), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*discovery.Endpoint)(nil), (*v1beta1.Endpoint)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_discovery_Endpoint_To_v1beta1_Endpoint(a.(*discovery.Endpoint), b.(*v1beta1.Endpoint), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1beta1.Endpoint)(nil), (*discovery.Endpoint)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_Endpoint_To_discovery_Endpoint(a.(*v1beta1.Endpoint), b.(*discovery.Endpoint), scope)
	}); err != nil {
		return err
	}
	return nil
}

func autoConvert_v1beta1_Endpoint_To_discovery_Endpoint(in *v1beta1.Endpoint, out *discovery.Endpoint, s conversion.Scope) error {
	out.Addresses = *(*[]string)(unsafe.Pointer(&in.Addresses))
	if err := Convert_v1beta1_EndpointConditions_To_discovery_EndpointConditions(&in.Conditions, &out.Conditions, s); err != nil {
		return err
	}
	out.Hostname = (*string)(unsafe.Pointer(in.Hostname))
	out.TargetRef = (*core.ObjectReference)(unsafe.Pointer(in.TargetRef))
	// WARNING: in.Topology requires manual conversion: does not exist in peer-type
	out.NodeName = (*string)(unsafe.Pointer(in.NodeName))
	out.Hints = (*discovery.EndpointHints)(unsafe.Pointer(in.Hints))
	return nil
}

func autoConvert_discovery_Endpoint_To_v1beta1_Endpoint(in *discovery.Endpoint, out *v1beta1.Endpoint, s conversion.Scope) error {
	out.Addresses = *(*[]string)(unsafe.Pointer(&in.Addresses))
	if err := Convert_discovery_EndpointConditions_To_v1beta1_EndpointConditions(&in.Conditions, &out.Conditions, s); err != nil {
		return err
	}
	out.Hostname = (*string)(unsafe.Pointer(in.Hostname))
	out.TargetRef = (*v1.ObjectReference)(unsafe.Pointer(in.TargetRef))
	// WARNING: in.DeprecatedTopology requires manual conversion: does not exist in peer-type
	out.NodeName = (*string)(unsafe.Pointer(in.NodeName))
	// WARNING: in.Zone requires manual conversion: does not exist in peer-type
	out.Hints = (*v1beta1.EndpointHints)(unsafe.Pointer(in.Hints))
	return nil
}

func autoConvert_v1beta1_EndpointConditions_To_discovery_EndpointConditions(in *v1beta1.EndpointConditions, out *discovery.EndpointConditions, s conversion.Scope) error {
	out.Ready = (*bool)(unsafe.Pointer(in.Ready))
	out.Serving = (*bool)(unsafe.Pointer(in.Serving))
	out.Terminating = (*bool)(unsafe.Pointer(in.Terminating))
	return nil
}

// Convert_v1beta1_EndpointConditions_To_discovery_EndpointConditions is an autogenerated conversion function.
func Convert_v1beta1_EndpointConditions_To_discovery_EndpointConditions(in *v1beta1.EndpointConditions, out *discovery.EndpointConditions, s conversion.Scope) error {
	return autoConvert_v1beta1_EndpointConditions_To_discovery_EndpointConditions(in, out, s)
}

func autoConvert_discovery_EndpointConditions_To_v1beta1_EndpointConditions(in *discovery.EndpointConditions, out *v1beta1.EndpointConditions, s conversion.Scope) error {
	out.Ready = (*bool)(unsafe.Pointer(in.Ready))
	out.Serving = (*bool)(unsafe.Pointer(in.Serving))
	out.Terminating = (*bool)(unsafe.Pointer(in.Terminating))
	return nil
}

// Convert_discovery_EndpointConditions_To_v1beta1_EndpointConditions is an autogenerated conversion function.
func Convert_discovery_EndpointConditions_To_v1beta1_EndpointConditions(in *discovery.EndpointConditions, out *v1beta1.EndpointConditions, s conversion.Scope) error {
	return autoConvert_discovery_EndpointConditions_To_v1beta1_EndpointConditions(in, out, s)
}

func autoConvert_v1beta1_EndpointHints_To_discovery_EndpointHints(in *v1beta1.EndpointHints, out *discovery.EndpointHints, s conversion.Scope) error {
	out.ForZones = *(*[]discovery.ForZone)(unsafe.Pointer(&in.ForZones))
	return nil
}

// Convert_v1beta1_EndpointHints_To_discovery_EndpointHints is an autogenerated conversion function.
func Convert_v1beta1_EndpointHints_To_discovery_EndpointHints(in *v1beta1.EndpointHints, out *discovery.EndpointHints, s conversion.Scope) error {
	return autoConvert_v1beta1_EndpointHints_To_discovery_EndpointHints(in, out, s)
}

func autoConvert_discovery_EndpointHints_To_v1beta1_EndpointHints(in *discovery.EndpointHints, out *v1beta1.EndpointHints, s conversion.Scope) error {
	out.ForZones = *(*[]v1beta1.ForZone)(unsafe.Pointer(&in.ForZones))
	return nil
}

// Convert_discovery_EndpointHints_To_v1beta1_EndpointHints is an autogenerated conversion function.
func Convert_discovery_EndpointHints_To_v1beta1_EndpointHints(in *discovery.EndpointHints, out *v1beta1.EndpointHints, s conversion.Scope) error {
	return autoConvert_discovery_EndpointHints_To_v1beta1_EndpointHints(in, out, s)
}

func autoConvert_v1beta1_EndpointPort_To_discovery_EndpointPort(in *v1beta1.EndpointPort, out *discovery.EndpointPort, s conversion.Scope) error {
	out.Name = (*string)(unsafe.Pointer(in.Name))
	out.Protocol = (*core.Protocol)(unsafe.Pointer(in.Protocol))
	out.Port = (*int32)(unsafe.Pointer(in.Port))
	out.AppProtocol = (*string)(unsafe.Pointer(in.AppProtocol))
	return nil
}

// Convert_v1beta1_EndpointPort_To_discovery_EndpointPort is an autogenerated conversion function.
func Convert_v1beta1_EndpointPort_To_discovery_EndpointPort(in *v1beta1.EndpointPort, out *discovery.EndpointPort, s conversion.Scope) error {
	return autoConvert_v1beta1_EndpointPort_To_discovery_EndpointPort(in, out, s)
}

func autoConvert_discovery_EndpointPort_To_v1beta1_EndpointPort(in *discovery.EndpointPort, out *v1beta1.EndpointPort, s conversion.Scope) error {
	out.Name = (*string)(unsafe.Pointer(in.Name))
	out.Protocol = (*v1.Protocol)(unsafe.Pointer(in.Protocol))
	out.Port = (*int32)(unsafe.Pointer(in.Port))
	out.AppProtocol = (*string)(unsafe.Pointer(in.AppProtocol))
	return nil
}

// Convert_discovery_EndpointPort_To_v1beta1_EndpointPort is an autogenerated conversion function.
func Convert_discovery_EndpointPort_To_v1beta1_EndpointPort(in *discovery.EndpointPort, out *v1beta1.EndpointPort, s conversion.Scope) error {
	return autoConvert_discovery_EndpointPort_To_v1beta1_EndpointPort(in, out, s)
}

func autoConvert_v1beta1_EndpointSlice_To_discovery_EndpointSlice(in *v1beta1.EndpointSlice, out *discovery.EndpointSlice, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	out.AddressType = discovery.AddressType(in.AddressType)
	if in.Endpoints != nil {
		in, out := &in.Endpoints, &out.Endpoints
		*out = make([]discovery.Endpoint, len(*in))
		for i := range *in {
			if err := Convert_v1beta1_Endpoint_To_discovery_Endpoint(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Endpoints = nil
	}
	out.Ports = *(*[]discovery.EndpointPort)(unsafe.Pointer(&in.Ports))
	return nil
}

// Convert_v1beta1_EndpointSlice_To_discovery_EndpointSlice is an autogenerated conversion function.
func Convert_v1beta1_EndpointSlice_To_discovery_EndpointSlice(in *v1beta1.EndpointSlice, out *discovery.EndpointSlice, s conversion.Scope) error {
	return autoConvert_v1beta1_EndpointSlice_To_discovery_EndpointSlice(in, out, s)
}

func autoConvert_discovery_EndpointSlice_To_v1beta1_EndpointSlice(in *discovery.EndpointSlice, out *v1beta1.EndpointSlice, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	out.AddressType = v1beta1.AddressType(in.AddressType)
	if in.Endpoints != nil {
		in, out := &in.Endpoints, &out.Endpoints
		*out = make([]v1beta1.Endpoint, len(*in))
		for i := range *in {
			if err := Convert_discovery_Endpoint_To_v1beta1_Endpoint(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Endpoints = nil
	}
	out.Ports = *(*[]v1beta1.EndpointPort)(unsafe.Pointer(&in.Ports))
	return nil
}

// Convert_discovery_EndpointSlice_To_v1beta1_EndpointSlice is an autogenerated conversion function.
func Convert_discovery_EndpointSlice_To_v1beta1_EndpointSlice(in *discovery.EndpointSlice, out *v1beta1.EndpointSlice, s conversion.Scope) error {
	return autoConvert_discovery_EndpointSlice_To_v1beta1_EndpointSlice(in, out, s)
}

func autoConvert_v1beta1_EndpointSliceList_To_discovery_EndpointSliceList(in *v1beta1.EndpointSliceList, out *discovery.EndpointSliceList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]discovery.EndpointSlice, len(*in))
		for i := range *in {
			if err := Convert_v1beta1_EndpointSlice_To_discovery_EndpointSlice(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

// Convert_v1beta1_EndpointSliceList_To_discovery_EndpointSliceList is an autogenerated conversion function.
func Convert_v1beta1_EndpointSliceList_To_discovery_EndpointSliceList(in *v1beta1.EndpointSliceList, out *discovery.EndpointSliceList, s conversion.Scope) error {
	return autoConvert_v1beta1_EndpointSliceList_To_discovery_EndpointSliceList(in, out, s)
}

func autoConvert_discovery_EndpointSliceList_To_v1beta1_EndpointSliceList(in *discovery.EndpointSliceList, out *v1beta1.EndpointSliceList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]v1beta1.EndpointSlice, len(*in))
		for i := range *in {
			if err := Convert_discovery_EndpointSlice_To_v1beta1_EndpointSlice(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

// Convert_discovery_EndpointSliceList_To_v1beta1_EndpointSliceList is an autogenerated conversion function.
func Convert_discovery_EndpointSliceList_To_v1beta1_EndpointSliceList(in *discovery.EndpointSliceList, out *v1beta1.EndpointSliceList, s conversion.Scope) error {
	return autoConvert_discovery_EndpointSliceList_To_v1beta1_EndpointSliceList(in, out, s)
}

func autoConvert_v1beta1_ForZone_To_discovery_ForZone(in *v1beta1.ForZone, out *discovery.ForZone, s conversion.Scope) error {
	out.Name = in.Name
	return nil
}

// Convert_v1beta1_ForZone_To_discovery_ForZone is an autogenerated conversion function.
func Convert_v1beta1_ForZone_To_discovery_ForZone(in *v1beta1.ForZone, out *discovery.ForZone, s conversion.Scope) error {
	return autoConvert_v1beta1_ForZone_To_discovery_ForZone(in, out, s)
}

func autoConvert_discovery_ForZone_To_v1beta1_ForZone(in *discovery.ForZone, out *v1beta1.ForZone, s conversion.Scope) error {
	out.Name = in.Name
	return nil
}

// Convert_discovery_ForZone_To_v1beta1_ForZone is an autogenerated conversion function.
func Convert_discovery_ForZone_To_v1beta1_ForZone(in *discovery.ForZone, out *v1beta1.ForZone, s conversion.Scope) error {
	return autoConvert_discovery_ForZone_To_v1beta1_ForZone(in, out, s)
}
