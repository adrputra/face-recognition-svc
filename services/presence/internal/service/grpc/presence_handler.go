package grpcservice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/adrputra/face-recognition-svc/presence-svc/internal/controller"
	domain "github.com/adrputra/face-recognition-svc/presence-svc/internal/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
)

const (
	MethodCreatePresence = "/presence.PresenceService/CreatePresence"
	MethodGetPresence    = "/presence.PresenceService/GetPresence"
	MethodListPresence   = "/presence.PresenceService/ListPresence"
	MethodUpdatePresence = "/presence.PresenceService/UpdatePresence"
	MethodDeletePresence = "/presence.PresenceService/DeletePresence"

	presencePayloadTypeURL = "type.googleapis.com/presence.Payload"
)

type PresenceServiceServer interface {
	CreatePresence(context.Context, *anypb.Any) (*anypb.Any, error)
	GetPresence(context.Context, *anypb.Any) (*anypb.Any, error)
	ListPresence(context.Context, *anypb.Any) (*anypb.Any, error)
	UpdatePresence(context.Context, *anypb.Any) (*anypb.Any, error)
	DeletePresence(context.Context, *anypb.Any) (*anypb.Any, error)
}

type Handler struct {
	controller controller.InterfacePresenceController
}

func NewHandler(controller controller.InterfacePresenceController) *Handler {
	return &Handler{controller: controller}
}

func RegisterPresenceServiceServer(s grpc.ServiceRegistrar, srv PresenceServiceServer) {
	s.RegisterService(&PresenceServiceDesc, srv)
}

func (h *Handler) CreatePresence(ctx context.Context, req *anypb.Any) (*anypb.Any, error) {
	payload, err := parsePayload(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid payload")
	}

	item, err := h.controller.CreatePresence(ctx, &domain.Entity{
		UserID:        toString(payload["user_id"]),
		InstitutionID: toString(payload["institution_id"]),
		Status:        toString(payload["status"]),
		Note:          toString(payload["note"]),
	})
	if err != nil {
		return nil, mapError(err)
	}

	return newPayload(entityToMap(item))
}

func (h *Handler) GetPresence(ctx context.Context, req *anypb.Any) (*anypb.Any, error) {
	payload, err := parsePayload(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid payload")
	}

	item, err := h.controller.GetPresence(ctx, toString(payload["id"]))
	if err != nil {
		return nil, mapError(err)
	}

	return newPayload(entityToMap(item))
}

func (h *Handler) ListPresence(ctx context.Context, req *anypb.Any) (*anypb.Any, error) {
	payload, err := parsePayload(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid payload")
	}

	items, err := h.controller.ListPresence(ctx, domain.Filter{
		UserID:        toString(payload["user_id"]),
		InstitutionID: toString(payload["institution_id"]),
		Status:        toString(payload["status"]),
	})
	if err != nil {
		return nil, mapError(err)
	}

	outItems := make([]interface{}, 0, len(items))
	for _, item := range items {
		outItems = append(outItems, entityToMap(item))
	}

	return newPayload(map[string]interface{}{
		"items": outItems,
		"total": len(outItems),
	})
}

func (h *Handler) UpdatePresence(ctx context.Context, req *anypb.Any) (*anypb.Any, error) {
	payload, err := parsePayload(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid payload")
	}

	patch := domain.Patch{}
	if raw, ok := payload["user_id"]; ok {
		value := toString(raw)
		patch.UserID = &value
	}
	if raw, ok := payload["institution_id"]; ok {
		value := toString(raw)
		patch.InstitutionID = &value
	}
	if raw, ok := payload["status"]; ok {
		value := toString(raw)
		patch.Status = &value
	}
	if raw, ok := payload["note"]; ok {
		value := toString(raw)
		patch.Note = &value
	}

	item, err := h.controller.UpdatePresence(ctx, toString(payload["id"]), patch)
	if err != nil {
		return nil, mapError(err)
	}

	return newPayload(entityToMap(item))
}

func (h *Handler) DeletePresence(ctx context.Context, req *anypb.Any) (*anypb.Any, error) {
	payload, err := parsePayload(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid payload")
	}

	id := toString(payload["id"])
	if err := h.controller.DeletePresence(ctx, id); err != nil {
		return nil, mapError(err)
	}

	return newPayload(map[string]interface{}{
		"id":      strings.TrimSpace(id),
		"deleted": true,
	})
}

var PresenceServiceDesc = grpc.ServiceDesc{
	ServiceName: "presence.PresenceService",
	HandlerType: (*PresenceServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreatePresence",
			Handler:    _PresenceService_CreatePresence_Handler,
		},
		{
			MethodName: "GetPresence",
			Handler:    _PresenceService_GetPresence_Handler,
		},
		{
			MethodName: "ListPresence",
			Handler:    _PresenceService_ListPresence_Handler,
		},
		{
			MethodName: "UpdatePresence",
			Handler:    _PresenceService_UpdatePresence_Handler,
		},
		{
			MethodName: "DeletePresence",
			Handler:    _PresenceService_DeletePresence_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "presence.clean-architecture",
}

func _PresenceService_CreatePresence_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(anypb.Any)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PresenceServiceServer).CreatePresence(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MethodCreatePresence,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PresenceServiceServer).CreatePresence(ctx, req.(*anypb.Any))
	}
	return interceptor(ctx, in, info, handler)
}

func _PresenceService_GetPresence_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(anypb.Any)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PresenceServiceServer).GetPresence(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MethodGetPresence,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PresenceServiceServer).GetPresence(ctx, req.(*anypb.Any))
	}
	return interceptor(ctx, in, info, handler)
}

func _PresenceService_ListPresence_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(anypb.Any)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PresenceServiceServer).ListPresence(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MethodListPresence,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PresenceServiceServer).ListPresence(ctx, req.(*anypb.Any))
	}
	return interceptor(ctx, in, info, handler)
}

func _PresenceService_UpdatePresence_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(anypb.Any)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PresenceServiceServer).UpdatePresence(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MethodUpdatePresence,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PresenceServiceServer).UpdatePresence(ctx, req.(*anypb.Any))
	}
	return interceptor(ctx, in, info, handler)
}

func _PresenceService_DeletePresence_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(anypb.Any)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PresenceServiceServer).DeletePresence(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MethodDeletePresence,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PresenceServiceServer).DeletePresence(ctx, req.(*anypb.Any))
	}
	return interceptor(ctx, in, info, handler)
}

func mapError(err error) error {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, controller.ErrRequiredUserID),
		errors.Is(err, controller.ErrRequiredID),
		errors.Is(err, controller.ErrRequiredUpdateField),
		errors.Is(err, controller.ErrEmptyUserID),
		errors.Is(err, controller.ErrEmptyStatus):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}

func newPayload(in map[string]interface{}) (*anypb.Any, error) {
	raw, err := json.Marshal(in)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "marshal response: %v", err)
	}
	return &anypb.Any{
		TypeUrl: presencePayloadTypeURL,
		Value:   raw,
	}, nil
}

func parsePayload(in *anypb.Any) (map[string]interface{}, error) {
	if in == nil || len(in.Value) == 0 {
		return map[string]interface{}{}, nil
	}

	var out map[string]interface{}
	if err := json.Unmarshal(in.Value, &out); err != nil {
		return nil, err
	}
	if out == nil {
		return map[string]interface{}{}, nil
	}
	return out, nil
}

func entityToMap(item *domain.Entity) map[string]interface{} {
	if item == nil {
		return map[string]interface{}{}
	}
	return map[string]interface{}{
		"id":             item.ID,
		"user_id":        item.UserID,
		"institution_id": item.InstitutionID,
		"status":         item.Status,
		"note":           item.Note,
		"created_at":     item.CreatedAt.Format(time.RFC3339Nano),
		"updated_at":     item.UpdatedAt.Format(time.RFC3339Nano),
	}
}

func toString(value interface{}) string {
	if value == nil {
		return ""
	}

	switch typed := value.(type) {
	case string:
		return typed
	case float64:
		return strconv.FormatFloat(typed, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(typed)
	default:
		return fmt.Sprintf("%v", value)
	}
}
