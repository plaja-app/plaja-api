package marshaler

import (
	"errors"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/plaja-app/plaja-api/gateway/internal/server/response"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

const successMessage = "OK"

func New() *Marshaler {
	return &Marshaler{
		Marshaler: &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				EmitUnpopulated: true,
				UseProtoNames:   true,
			},
		},
	}
}

type Marshaler struct {
	runtime.Marshaler
}

func (m *Marshaler) Marshal(v any) ([]byte, error) {
	msg, ok := v.(proto.Message)
	if !ok {
		return nil, errors.New("failed to cast value to proto msg interface")
	}

	return m.Marshaler.Marshal(response.NewSuccessResponse(successMessage, msg))
}
