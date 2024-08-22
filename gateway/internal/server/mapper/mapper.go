package mapper

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/plaja-app/plaja-api/pkg/logger"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/plaja-app/plaja-api/gateway/internal/server/response"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"
)

func New(log *logger.Logger) *ResponseMapper {
	return &ResponseMapper{log}
}

type ResponseMapper struct {
	log *logger.Logger
}

func (rm *ResponseMapper) MapGRPCErr(
	_ context.Context,
	_ *runtime.ServeMux,
	_ runtime.Marshaler,
	w http.ResponseWriter,
	_ *http.Request,
	err error,
) {
	s := status.Convert(err)
	res, err := json.Marshal(response.NewErrorResponse(s.Message()))
	if err != nil {
		rm.log.Error("Failed to convert marshall err response", zap.Error(err))
	}

	w.WriteHeader(runtime.HTTPStatusFromCode(s.Code()))
	if _, err = w.Write(res); err != nil {
		rm.log.Error("Failed to write error", zap.Error(err))
	}
}
