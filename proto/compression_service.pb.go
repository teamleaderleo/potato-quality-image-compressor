// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v6.30.2
// source: proto/compression_service.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// CompressImageRequest contains the data and parameters for image compression
type CompressImageRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The image data to be compressed
	ImageData []byte `protobuf:"bytes,1,opt,name=image_data,json=imageData,proto3" json:"image_data,omitempty"`
	// The requested quality level (1-100)
	Quality int32 `protobuf:"varint,2,opt,name=quality,proto3" json:"quality,omitempty"`
	// The requested output format
	Format string `protobuf:"bytes,3,opt,name=format,proto3" json:"format,omitempty"`
	// The requested compression strategy
	Strategy string `protobuf:"bytes,4,opt,name=strategy,proto3" json:"strategy,omitempty"`
	// Optional original filename
	Filename string `protobuf:"bytes,5,opt,name=filename,proto3" json:"filename,omitempty"`
}

func (x *CompressImageRequest) Reset() {
	*x = CompressImageRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_compression_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CompressImageRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CompressImageRequest) ProtoMessage() {}

func (x *CompressImageRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_compression_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CompressImageRequest.ProtoReflect.Descriptor instead.
func (*CompressImageRequest) Descriptor() ([]byte, []int) {
	return file_proto_compression_service_proto_rawDescGZIP(), []int{0}
}

func (x *CompressImageRequest) GetImageData() []byte {
	if x != nil {
		return x.ImageData
	}
	return nil
}

func (x *CompressImageRequest) GetQuality() int32 {
	if x != nil {
		return x.Quality
	}
	return 0
}

func (x *CompressImageRequest) GetFormat() string {
	if x != nil {
		return x.Format
	}
	return ""
}

func (x *CompressImageRequest) GetStrategy() string {
	if x != nil {
		return x.Strategy
	}
	return ""
}

func (x *CompressImageRequest) GetFilename() string {
	if x != nil {
		return x.Filename
	}
	return ""
}

// CompressImageResponse contains the compressed image and metadata
type CompressImageResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The compressed image data
	ImageData []byte `protobuf:"bytes,1,opt,name=image_data,json=imageData,proto3" json:"image_data,omitempty"`
	// The format of the compressed image
	Format string `protobuf:"bytes,2,opt,name=format,proto3" json:"format,omitempty"`
	// The size of the original image in bytes
	OriginalSize int64 `protobuf:"varint,3,opt,name=original_size,json=originalSize,proto3" json:"original_size,omitempty"`
	// The size of the compressed image in bytes
	CompressedSize int64 `protobuf:"varint,4,opt,name=compressed_size,json=compressedSize,proto3" json:"compressed_size,omitempty"`
	// The compression ratio (compressed_size / original_size)
	CompressionRatio float64 `protobuf:"fixed64,5,opt,name=compression_ratio,json=compressionRatio,proto3" json:"compression_ratio,omitempty"`
	// Time taken to compress the image in milliseconds
	ProcessingTimeMs int64 `protobuf:"varint,6,opt,name=processing_time_ms,json=processingTimeMs,proto3" json:"processing_time_ms,omitempty"`
	// Error message (if any)
	Error string `protobuf:"bytes,7,opt,name=error,proto3" json:"error,omitempty"`
	// The filename (if provided in the request)
	Filename string `protobuf:"bytes,8,opt,name=filename,proto3" json:"filename,omitempty"`
}

func (x *CompressImageResponse) Reset() {
	*x = CompressImageResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_compression_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CompressImageResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CompressImageResponse) ProtoMessage() {}

func (x *CompressImageResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_compression_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CompressImageResponse.ProtoReflect.Descriptor instead.
func (*CompressImageResponse) Descriptor() ([]byte, []int) {
	return file_proto_compression_service_proto_rawDescGZIP(), []int{1}
}

func (x *CompressImageResponse) GetImageData() []byte {
	if x != nil {
		return x.ImageData
	}
	return nil
}

func (x *CompressImageResponse) GetFormat() string {
	if x != nil {
		return x.Format
	}
	return ""
}

func (x *CompressImageResponse) GetOriginalSize() int64 {
	if x != nil {
		return x.OriginalSize
	}
	return 0
}

func (x *CompressImageResponse) GetCompressedSize() int64 {
	if x != nil {
		return x.CompressedSize
	}
	return 0
}

func (x *CompressImageResponse) GetCompressionRatio() float64 {
	if x != nil {
		return x.CompressionRatio
	}
	return 0
}

func (x *CompressImageResponse) GetProcessingTimeMs() int64 {
	if x != nil {
		return x.ProcessingTimeMs
	}
	return 0
}

func (x *CompressImageResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

func (x *CompressImageResponse) GetFilename() string {
	if x != nil {
		return x.Filename
	}
	return ""
}

// BatchCompressRequest contains multiple images to compress
type BatchCompressRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// List of compression requests
	Requests []*CompressImageRequest `protobuf:"bytes,1,rep,name=requests,proto3" json:"requests,omitempty"`
}

func (x *BatchCompressRequest) Reset() {
	*x = BatchCompressRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_compression_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BatchCompressRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BatchCompressRequest) ProtoMessage() {}

func (x *BatchCompressRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_compression_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BatchCompressRequest.ProtoReflect.Descriptor instead.
func (*BatchCompressRequest) Descriptor() ([]byte, []int) {
	return file_proto_compression_service_proto_rawDescGZIP(), []int{2}
}

func (x *BatchCompressRequest) GetRequests() []*CompressImageRequest {
	if x != nil {
		return x.Requests
	}
	return nil
}

// BatchCompressResponse contains multiple compressed images
type BatchCompressResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// List of compression responses
	Responses []*CompressImageResponse `protobuf:"bytes,1,rep,name=responses,proto3" json:"responses,omitempty"`
	// Total time taken to process the batch in milliseconds
	TotalProcessingTimeMs int64 `protobuf:"varint,2,opt,name=total_processing_time_ms,json=totalProcessingTimeMs,proto3" json:"total_processing_time_ms,omitempty"`
}

func (x *BatchCompressResponse) Reset() {
	*x = BatchCompressResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_compression_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BatchCompressResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BatchCompressResponse) ProtoMessage() {}

func (x *BatchCompressResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_compression_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BatchCompressResponse.ProtoReflect.Descriptor instead.
func (*BatchCompressResponse) Descriptor() ([]byte, []int) {
	return file_proto_compression_service_proto_rawDescGZIP(), []int{3}
}

func (x *BatchCompressResponse) GetResponses() []*CompressImageResponse {
	if x != nil {
		return x.Responses
	}
	return nil
}

func (x *BatchCompressResponse) GetTotalProcessingTimeMs() int64 {
	if x != nil {
		return x.TotalProcessingTimeMs
	}
	return 0
}

// ServiceStatsRequest is used to request service statistics
type ServiceStatsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Optional time period in seconds (0 = all time)
	TimePeriodSeconds int64 `protobuf:"varint,1,opt,name=time_period_seconds,json=timePeriodSeconds,proto3" json:"time_period_seconds,omitempty"`
}

func (x *ServiceStatsRequest) Reset() {
	*x = ServiceStatsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_compression_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServiceStatsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServiceStatsRequest) ProtoMessage() {}

func (x *ServiceStatsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_compression_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServiceStatsRequest.ProtoReflect.Descriptor instead.
func (*ServiceStatsRequest) Descriptor() ([]byte, []int) {
	return file_proto_compression_service_proto_rawDescGZIP(), []int{4}
}

func (x *ServiceStatsRequest) GetTimePeriodSeconds() int64 {
	if x != nil {
		return x.TimePeriodSeconds
	}
	return 0
}

// ServiceStatsResponse contains service statistics
type ServiceStatsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Total number of requests processed
	TotalRequests int64 `protobuf:"varint,1,opt,name=total_requests,json=totalRequests,proto3" json:"total_requests,omitempty"`
	// Total number of images processed
	TotalImages int64 `protobuf:"varint,2,opt,name=total_images,json=totalImages,proto3" json:"total_images,omitempty"`
	// Average processing time in milliseconds
	AvgProcessingTimeMs float64 `protobuf:"fixed64,3,opt,name=avg_processing_time_ms,json=avgProcessingTimeMs,proto3" json:"avg_processing_time_ms,omitempty"`
	// Average compression ratio
	AvgCompressionRatio float64 `protobuf:"fixed64,4,opt,name=avg_compression_ratio,json=avgCompressionRatio,proto3" json:"avg_compression_ratio,omitempty"`
	// Number of worker threads
	WorkerCount int32 `protobuf:"varint,5,opt,name=worker_count,json=workerCount,proto3" json:"worker_count,omitempty"`
	// Number of busy workers
	BusyWorkers int32 `protobuf:"varint,6,opt,name=busy_workers,json=busyWorkers,proto3" json:"busy_workers,omitempty"`
	// Current memory usage in bytes
	MemoryUsageBytes int64 `protobuf:"varint,7,opt,name=memory_usage_bytes,json=memoryUsageBytes,proto3" json:"memory_usage_bytes,omitempty"`
}

func (x *ServiceStatsResponse) Reset() {
	*x = ServiceStatsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_compression_service_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServiceStatsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServiceStatsResponse) ProtoMessage() {}

func (x *ServiceStatsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_compression_service_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServiceStatsResponse.ProtoReflect.Descriptor instead.
func (*ServiceStatsResponse) Descriptor() ([]byte, []int) {
	return file_proto_compression_service_proto_rawDescGZIP(), []int{5}
}

func (x *ServiceStatsResponse) GetTotalRequests() int64 {
	if x != nil {
		return x.TotalRequests
	}
	return 0
}

func (x *ServiceStatsResponse) GetTotalImages() int64 {
	if x != nil {
		return x.TotalImages
	}
	return 0
}

func (x *ServiceStatsResponse) GetAvgProcessingTimeMs() float64 {
	if x != nil {
		return x.AvgProcessingTimeMs
	}
	return 0
}

func (x *ServiceStatsResponse) GetAvgCompressionRatio() float64 {
	if x != nil {
		return x.AvgCompressionRatio
	}
	return 0
}

func (x *ServiceStatsResponse) GetWorkerCount() int32 {
	if x != nil {
		return x.WorkerCount
	}
	return 0
}

func (x *ServiceStatsResponse) GetBusyWorkers() int32 {
	if x != nil {
		return x.BusyWorkers
	}
	return 0
}

func (x *ServiceStatsResponse) GetMemoryUsageBytes() int64 {
	if x != nil {
		return x.MemoryUsageBytes
	}
	return 0
}

var File_proto_compression_service_proto protoreflect.FileDescriptor

var file_proto_compression_service_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x6d, 0x70, 0x72, 0x65, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x0b, 0x63, 0x6f, 0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x9f,
	0x01, 0x0a, 0x14, 0x43, 0x6f, 0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x49, 0x6d, 0x61, 0x67, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x69, 0x6d, 0x61, 0x67, 0x65,
	0x5f, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x69, 0x6d, 0x61,
	0x67, 0x65, 0x44, 0x61, 0x74, 0x61, 0x12, 0x18, 0x0a, 0x07, 0x71, 0x75, 0x61, 0x6c, 0x69, 0x74,
	0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x71, 0x75, 0x61, 0x6c, 0x69, 0x74, 0x79,
	0x12, 0x16, 0x0a, 0x06, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x74, 0x72, 0x61,
	0x74, 0x65, 0x67, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x74, 0x72, 0x61,
	0x74, 0x65, 0x67, 0x79, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65,
	0x22, 0xa9, 0x02, 0x0a, 0x15, 0x43, 0x6f, 0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x49, 0x6d, 0x61,
	0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x69, 0x6d,
	0x61, 0x67, 0x65, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09,
	0x69, 0x6d, 0x61, 0x67, 0x65, 0x44, 0x61, 0x74, 0x61, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x6f, 0x72,
	0x6d, 0x61, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x66, 0x6f, 0x72, 0x6d, 0x61,
	0x74, 0x12, 0x23, 0x0a, 0x0d, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x6c, 0x5f, 0x73, 0x69,
	0x7a, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e,
	0x61, 0x6c, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x27, 0x0a, 0x0f, 0x63, 0x6f, 0x6d, 0x70, 0x72, 0x65,
	0x73, 0x73, 0x65, 0x64, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x0e, 0x63, 0x6f, 0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x65, 0x64, 0x53, 0x69, 0x7a, 0x65, 0x12,
	0x2b, 0x0a, 0x11, 0x63, 0x6f, 0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x18, 0x05, 0x20, 0x01, 0x28, 0x01, 0x52, 0x10, 0x63, 0x6f, 0x6d, 0x70,
	0x72, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x61, 0x74, 0x69, 0x6f, 0x12, 0x2c, 0x0a, 0x12,
	0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x69, 0x6e, 0x67, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x5f,
	0x6d, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x10, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73,
	0x73, 0x69, 0x6e, 0x67, 0x54, 0x69, 0x6d, 0x65, 0x4d, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x72,
	0x72, 0x6f, 0x72, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72,
	0x12, 0x1a, 0x0a, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x08, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x55, 0x0a, 0x14,
	0x42, 0x61, 0x74, 0x63, 0x68, 0x43, 0x6f, 0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x3d, 0x0a, 0x08, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x72, 0x65, 0x73,
	0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x43, 0x6f, 0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x49, 0x6d, 0x61,
	0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x08, 0x72, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x73, 0x22, 0x92, 0x01, 0x0a, 0x15, 0x42, 0x61, 0x74, 0x63, 0x68, 0x43, 0x6f, 0x6d,
	0x70, 0x72, 0x65, 0x73, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x40, 0x0a,
	0x09, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x22, 0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x43,
	0x6f, 0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x52, 0x09, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x73, 0x12,
	0x37, 0x0a, 0x18, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73,
	0x69, 0x6e, 0x67, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x5f, 0x6d, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x15, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x69,
	0x6e, 0x67, 0x54, 0x69, 0x6d, 0x65, 0x4d, 0x73, 0x22, 0x45, 0x0a, 0x13, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x53, 0x74, 0x61, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x2e, 0x0a, 0x13, 0x74, 0x69, 0x6d, 0x65, 0x5f, 0x70, 0x65, 0x72, 0x69, 0x6f, 0x64, 0x5f, 0x73,
	0x65, 0x63, 0x6f, 0x6e, 0x64, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x11, 0x74, 0x69,
	0x6d, 0x65, 0x50, 0x65, 0x72, 0x69, 0x6f, 0x64, 0x53, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x73, 0x22,
	0xbd, 0x02, 0x0a, 0x14, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x53, 0x74, 0x61, 0x74, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x25, 0x0a, 0x0e, 0x74, 0x6f, 0x74, 0x61,
	0x6c, 0x5f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x0d, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x73, 0x12,
	0x21, 0x0a, 0x0c, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x73, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x49, 0x6d, 0x61, 0x67,
	0x65, 0x73, 0x12, 0x33, 0x0a, 0x16, 0x61, 0x76, 0x67, 0x5f, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73,
	0x73, 0x69, 0x6e, 0x67, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x5f, 0x6d, 0x73, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x01, 0x52, 0x13, 0x61, 0x76, 0x67, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x69, 0x6e,
	0x67, 0x54, 0x69, 0x6d, 0x65, 0x4d, 0x73, 0x12, 0x32, 0x0a, 0x15, 0x61, 0x76, 0x67, 0x5f, 0x63,
	0x6f, 0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x01, 0x52, 0x13, 0x61, 0x76, 0x67, 0x43, 0x6f, 0x6d, 0x70, 0x72,
	0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x61, 0x74, 0x69, 0x6f, 0x12, 0x21, 0x0a, 0x0c, 0x77,
	0x6f, 0x72, 0x6b, 0x65, 0x72, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x0b, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x21,
	0x0a, 0x0c, 0x62, 0x75, 0x73, 0x79, 0x5f, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x73, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x62, 0x75, 0x73, 0x79, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72,
	0x73, 0x12, 0x2c, 0x0a, 0x12, 0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x5f, 0x75, 0x73, 0x61, 0x67,
	0x65, 0x5f, 0x62, 0x79, 0x74, 0x65, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x03, 0x52, 0x10, 0x6d,
	0x65, 0x6d, 0x6f, 0x72, 0x79, 0x55, 0x73, 0x61, 0x67, 0x65, 0x42, 0x79, 0x74, 0x65, 0x73, 0x32,
	0x8a, 0x03, 0x0a, 0x17, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x43, 0x6f, 0x6d, 0x70, 0x72, 0x65, 0x73,
	0x73, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x56, 0x0a, 0x0d, 0x43,
	0x6f, 0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x12, 0x21, 0x2e, 0x63,
	0x6f, 0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x43, 0x6f, 0x6d, 0x70, 0x72,
	0x65, 0x73, 0x73, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x22, 0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x43, 0x6f,
	0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x5c, 0x0a, 0x13, 0x42, 0x61, 0x74, 0x63, 0x68, 0x43, 0x6f, 0x6d, 0x70,
	0x72, 0x65, 0x73, 0x73, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x73, 0x12, 0x21, 0x2e, 0x63, 0x6f, 0x6d,
	0x70, 0x72, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x42, 0x61, 0x74, 0x63, 0x68, 0x43, 0x6f,
	0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e,
	0x63, 0x6f, 0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x42, 0x61, 0x74, 0x63,
	0x68, 0x43, 0x6f, 0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x61, 0x0a, 0x14, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x43, 0x6f, 0x6d, 0x70, 0x72,
	0x65, 0x73, 0x73, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x73, 0x12, 0x21, 0x2e, 0x63, 0x6f, 0x6d, 0x70,
	0x72, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x43, 0x6f, 0x6d, 0x70, 0x72, 0x65, 0x73, 0x73,
	0x49, 0x6d, 0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e, 0x63,
	0x6f, 0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x43, 0x6f, 0x6d, 0x70, 0x72,
	0x65, 0x73, 0x73, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x28, 0x01, 0x30, 0x01, 0x12, 0x56, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x53, 0x74, 0x61, 0x74, 0x73, 0x12, 0x20, 0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x72, 0x65,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x53, 0x74, 0x61,
	0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x63, 0x6f, 0x6d, 0x70,
	0x72, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x53,
	0x74, 0x61, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x40, 0x5a, 0x3e,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74, 0x65, 0x61, 0x6d, 0x6c,
	0x65, 0x61, 0x64, 0x65, 0x72, 0x6c, 0x65, 0x6f, 0x2f, 0x70, 0x6f, 0x74, 0x61, 0x74, 0x6f, 0x2d,
	0x71, 0x75, 0x61, 0x6c, 0x69, 0x74, 0x79, 0x2d, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x2d, 0x63, 0x6f,
	0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x6f, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_compression_service_proto_rawDescOnce sync.Once
	file_proto_compression_service_proto_rawDescData = file_proto_compression_service_proto_rawDesc
)

func file_proto_compression_service_proto_rawDescGZIP() []byte {
	file_proto_compression_service_proto_rawDescOnce.Do(func() {
		file_proto_compression_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_compression_service_proto_rawDescData)
	})
	return file_proto_compression_service_proto_rawDescData
}

var file_proto_compression_service_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_proto_compression_service_proto_goTypes = []interface{}{
	(*CompressImageRequest)(nil),  // 0: compression.CompressImageRequest
	(*CompressImageResponse)(nil), // 1: compression.CompressImageResponse
	(*BatchCompressRequest)(nil),  // 2: compression.BatchCompressRequest
	(*BatchCompressResponse)(nil), // 3: compression.BatchCompressResponse
	(*ServiceStatsRequest)(nil),   // 4: compression.ServiceStatsRequest
	(*ServiceStatsResponse)(nil),  // 5: compression.ServiceStatsResponse
}
var file_proto_compression_service_proto_depIdxs = []int32{
	0, // 0: compression.BatchCompressRequest.requests:type_name -> compression.CompressImageRequest
	1, // 1: compression.BatchCompressResponse.responses:type_name -> compression.CompressImageResponse
	0, // 2: compression.ImageCompressionService.CompressImage:input_type -> compression.CompressImageRequest
	2, // 3: compression.ImageCompressionService.BatchCompressImages:input_type -> compression.BatchCompressRequest
	0, // 4: compression.ImageCompressionService.StreamCompressImages:input_type -> compression.CompressImageRequest
	4, // 5: compression.ImageCompressionService.GetServiceStats:input_type -> compression.ServiceStatsRequest
	1, // 6: compression.ImageCompressionService.CompressImage:output_type -> compression.CompressImageResponse
	3, // 7: compression.ImageCompressionService.BatchCompressImages:output_type -> compression.BatchCompressResponse
	1, // 8: compression.ImageCompressionService.StreamCompressImages:output_type -> compression.CompressImageResponse
	5, // 9: compression.ImageCompressionService.GetServiceStats:output_type -> compression.ServiceStatsResponse
	6, // [6:10] is the sub-list for method output_type
	2, // [2:6] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_proto_compression_service_proto_init() }
func file_proto_compression_service_proto_init() {
	if File_proto_compression_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_compression_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CompressImageRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_compression_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CompressImageResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_compression_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BatchCompressRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_compression_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BatchCompressResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_compression_service_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServiceStatsRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_compression_service_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServiceStatsResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_compression_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_compression_service_proto_goTypes,
		DependencyIndexes: file_proto_compression_service_proto_depIdxs,
		MessageInfos:      file_proto_compression_service_proto_msgTypes,
	}.Build()
	File_proto_compression_service_proto = out.File
	file_proto_compression_service_proto_rawDesc = nil
	file_proto_compression_service_proto_goTypes = nil
	file_proto_compression_service_proto_depIdxs = nil
}
