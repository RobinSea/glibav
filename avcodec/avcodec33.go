package avcodec

//#include <libavutil/avutil.h>
//#include <libavcodec/avcodec.h>
//////////////////#cgo CFLAGS: -D_REENTRANT
////////////////#cgo LDFLAGS:-static -D_REENTRANT
// #cgo pkg-config: libavcodec libavutil
import "C"

import (
	"unsafe"

	"github.com/imkira/go-libav/avutil"
)

type CodecParameters struct {
	CAVCodecParameters *C.AVCodecParameters
}

func NewCodecParameters() (*CodecParameters, error) {
	cPkt := (*C.AVCodecParameters)(C.avcodec_parameters_alloc())
	if cPkt == nil {
		return nil, ErrAllocationError
	}
	return NewCodecParametersFromC(unsafe.Pointer(cPkt)), nil
}

func NewCodecParametersFromC(cPSD unsafe.Pointer) *CodecParameters {
	return &CodecParameters{CAVCodecParameters: (*C.AVCodecParameters)(cPSD)}
}

func (cParams *CodecParameters) Free() {
	C.avcodec_parameters_free(&cParams.CAVCodecParameters)
}

func (ctx *Context) CopyFromBackup(dst *Context) error {
	// added in lavc 57.33.100
	parameters, err := NewCodecParameters()
	if err != nil {
		return err
	}
	defer parameters.Free()
	cParams := (*C.AVCodecParameters)(unsafe.Pointer(parameters.CAVCodecParameters))
	code := C.avcodec_parameters_from_context(cParams, ctx.CAVCodecContext)
	if code < 0 {
		return avutil.NewErrorFromCode(avutil.ErrorCode(code))
	}
	code = C.avcodec_parameters_to_context(dst.CAVCodecContext, cParams)
	if code < 0 {
		return avutil.NewErrorFromCode(avutil.ErrorCode(code))
	}
	return nil
}

func (ctx *Context) CopyFrom(dst *Context, Codecpar *Parameters) error {

	code := C.avcodec_parameters_from_context(Codecpar.CAVCodecParameters, dst.CAVCodecContext)
	if code < 0 {
		return avutil.NewErrorFromCode(avutil.ErrorCode(code))
	}
	return nil
}

func (ctx *Context) CopyTo(dst *Context, Codecpar *Parameters) error {

	code := C.avcodec_parameters_to_context(dst.CAVCodecContext, Codecpar.CAVCodecParameters)
	if code < 0 {
		return avutil.NewErrorFromCode(avutil.ErrorCode(code))
	}
	return nil
}

func (ctx *Context) DecodeVideo(pkt *Packet, frame *avutil.Frame) (bool, error) {
	cFrame := (*C.AVFrame)(unsafe.Pointer(frame.CAVFrame))
	cPkt := (*C.AVPacket)(unsafe.Pointer(pkt.CAVPacket))
	C.avcodec_send_packet(ctx.CAVCodecContext, cPkt)
	code := C.avcodec_receive_frame(ctx.CAVCodecContext, cFrame)
	var err error
	if code < 0 {
		err = avutil.NewErrorFromCode(avutil.ErrorCode(code))
		//code = 0
	}
	return code >= 0, err
}

func (ctx *Context) DecodeVideoTwo(frame *avutil.Frame) (bool, error) {
	cFrame := (*C.AVFrame)(unsafe.Pointer(frame.CAVFrame))
	//C.avcodec_send_packet(ctx.CAVCodecContext, nil)
	code := C.avcodec_receive_frame(ctx.CAVCodecContext, cFrame)
	var err error
	if code < 0 {
		err = avutil.NewErrorFromCode(avutil.ErrorCode(code))
		//code = 0
	}
	return code >= 0, err
}

func (ctx *Context) DecodeAudio(pkt *Packet, frame *avutil.Frame) (bool, error) {
	cFrame := (*C.AVFrame)(unsafe.Pointer(frame.CAVFrame))
	cPkt := (*C.AVPacket)(unsafe.Pointer(pkt.CAVPacket))
	C.avcodec_send_packet(ctx.CAVCodecContext, cPkt)
	code := C.avcodec_receive_frame(ctx.CAVCodecContext, cFrame)
	var err error
	if code < 0 {
		err = avutil.NewErrorFromCode(avutil.ErrorCode(code))
		//code = 0
	}
	return code >= 0, err
}

func (ctx *Context) DecodeAudioTwo(frame *avutil.Frame) (bool, error) {
	cFrame := (*C.AVFrame)(unsafe.Pointer(frame.CAVFrame))
	//C.avcodec_send_packet(ctx.CAVCodecContext, nil)
	code := C.avcodec_receive_frame(ctx.CAVCodecContext, cFrame)
	var err error
	if code < 0 {
		err = avutil.NewErrorFromCode(avutil.ErrorCode(code))
		//code = 0
	}
	return code >= 0, err
}

func (ctx *Context) EncodeVideoSend(pkt *Packet, frame *avutil.Frame) (bool, error) {
	var cFrame *C.AVFrame
	if frame != nil {
		cFrame = (*C.AVFrame)(unsafe.Pointer(frame.CAVFrame))
	}
	//cPkt := (*C.AVPacket)(unsafe.Pointer(pkt.CAVPacket))
	code := C.avcodec_send_frame(ctx.CAVCodecContext, cFrame)

	var err error
	if code < 0 {
		err = avutil.NewErrorFromCode(avutil.ErrorCode(code))
		return code < 0, err
	}

	// code = C.avcodec_receive_packet(ctx.CAVCodecContext, cPkt)
	// if code < 0 {
	// 	err = avutil.NewErrorFromCode(avutil.ErrorCode(code))
	// }
	return code >= 0, err
}

func (ctx *Context) EncodeVideoReceive(pkt *Packet) (bool, error) {

	cPkt := (*C.AVPacket)(unsafe.Pointer(pkt.CAVPacket))
	var err error
	code := C.avcodec_receive_packet(ctx.CAVCodecContext, cPkt)
	if code < 0 {
		err = avutil.NewErrorFromCode(avutil.ErrorCode(code))
	}
	return code >= 0, err
}

func (ctx *Context) EncodeAudioSend(pkt *Packet, frame *avutil.Frame) (bool, error) {
	var cFrame *C.AVFrame
	if frame != nil {
		cFrame = (*C.AVFrame)(unsafe.Pointer(frame.CAVFrame))
	}
	//	cPkt := (*C.AVPacket)(unsafe.Pointer(pkt.CAVPacket))

	var err error
	code := C.avcodec_send_frame(ctx.CAVCodecContext, cFrame)
	if code < 0 {
		err = avutil.NewErrorFromCode(avutil.ErrorCode(code))
		return code < 0, err
	}

	// code = C.avcodec_receive_packet(ctx.CAVCodecContext, cPkt)
	// if code < 0 {
	// 	err = avutil.NewErrorFromCode(avutil.ErrorCode(code))
	// }
	return code >= 0, err
}

func (ctx *Context) EncodeAudioReceive(pkt *Packet) (bool, error) {

	cPkt := (*C.AVPacket)(unsafe.Pointer(pkt.CAVPacket))
	var err error
	code := C.avcodec_receive_packet(ctx.CAVCodecContext, cPkt)
	if code < 0 {
		err = avutil.NewErrorFromCode(avutil.ErrorCode(code))
	}
	return code >= 0, err
}
