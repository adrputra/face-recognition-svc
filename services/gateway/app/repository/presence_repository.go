package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/adrputra/face-recognition-svc/gateway/app/config"
	"github.com/adrputra/face-recognition-svc/gateway/app/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
)

const (
	createPresenceMethod = "/presence.PresenceService/CreatePresence"
	getPresenceMethod    = "/presence.PresenceService/GetPresence"
	listPresenceMethod   = "/presence.PresenceService/ListPresence"
	updatePresenceMethod = "/presence.PresenceService/UpdatePresence"
	deletePresenceMethod = "/presence.PresenceService/DeletePresence"

	presencePayloadTypeURL = "type.googleapis.com/presence.Payload"
)

type PresenceRepository struct {
	addr string

	mu   sync.Mutex
	conn *grpc.ClientConn
}

func NewPresenceRepository(cfg *config.Config) *PresenceRepository {
	host := strings.TrimSpace(cfg.API.PresenceSVC.Host)
	if host == "" {
		host = "localhost"
	}

	port := cfg.API.PresenceSVC.Port
	if port == 0 {
		port = 50051
	}

	return &PresenceRepository{
		addr: fmt.Sprintf("%s:%d", host, port),
	}
}

func (r *PresenceRepository) Create(ctx context.Context, req *domain.Presence) (*domain.Presence, error) {
	payload := map[string]interface{}{
		"user_id": req.UserID,
		"status":  req.Status,
	}
	if req.InstitutionID != "" {
		payload["institution_id"] = req.InstitutionID
	}
	if req.Note != "" {
		payload["note"] = req.Note
	}

	out, err := r.invoke(ctx, createPresenceMethod, payload)
	if err != nil {
		return nil, err
	}

	return presenceFromMap(out), nil
}

func (r *PresenceRepository) Get(ctx context.Context, id string) (*domain.Presence, error) {
	out, err := r.invoke(ctx, getPresenceMethod, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return nil, err
	}
	return presenceFromMap(out), nil
}

func (r *PresenceRepository) List(ctx context.Context, filter domain.PresenceFilter) ([]*domain.Presence, error) {
	payload := map[string]interface{}{}
	if filter.UserID != "" {
		payload["user_id"] = filter.UserID
	}
	if filter.InstitutionID != "" {
		payload["institution_id"] = filter.InstitutionID
	}
	if filter.Status != "" {
		payload["status"] = filter.Status
	}

	out, err := r.invoke(ctx, listPresenceMethod, payload)
	if err != nil {
		return nil, err
	}

	items := out["items"]
	list, ok := items.([]interface{})
	if !ok {
		return []*domain.Presence{}, nil
	}

	result := make([]*domain.Presence, 0, len(list))
	for _, item := range list {
		record, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		result = append(result, presenceFromMap(record))
	}
	return result, nil
}

func (r *PresenceRepository) Update(ctx context.Context, id string, req *domain.Presence) (*domain.Presence, error) {
	payload := map[string]interface{}{
		"id": id,
	}
	if req.UserID != "" {
		payload["user_id"] = req.UserID
	}
	if req.InstitutionID != "" {
		payload["institution_id"] = req.InstitutionID
	}
	if req.Status != "" {
		payload["status"] = req.Status
	}
	if req.Note != "" {
		payload["note"] = req.Note
	}

	out, err := r.invoke(ctx, updatePresenceMethod, payload)
	if err != nil {
		return nil, err
	}

	return presenceFromMap(out), nil
}

func (r *PresenceRepository) Delete(ctx context.Context, id string) error {
	_, err := r.invoke(ctx, deletePresenceMethod, map[string]interface{}{
		"id": id,
	})
	return err
}

func (r *PresenceRepository) invoke(ctx context.Context, method string, payload map[string]interface{}) (map[string]interface{}, error) {
	conn, err := r.getConn()
	if err != nil {
		return nil, domain.ThrowError(http.StatusServiceUnavailable, err)
	}

	req, err := newPayload(payload)
	if err != nil {
		return nil, domain.ThrowError(http.StatusBadRequest, err)
	}

	resp := &anypb.Any{}
	callCtx := forwardMetadata(ctx)
	if err := conn.Invoke(callCtx, method, req, resp); err != nil {
		return nil, grpcErrorToHTTP(err)
	}

	out, err := parsePayload(resp)
	if err != nil {
		return nil, domain.ThrowError(http.StatusInternalServerError, err)
	}
	return out, nil
}

func (r *PresenceRepository) getConn() (*grpc.ClientConn, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.conn != nil {
		return r.conn, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		r.addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("dial presence service %s: %w", r.addr, err)
	}

	r.conn = conn
	return r.conn, nil
}

func forwardMetadata(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		return metadata.NewOutgoingContext(ctx, md)
	}
	if md, ok := metadata.FromOutgoingContext(ctx); ok {
		return metadata.NewOutgoingContext(ctx, md)
	}
	return ctx
}

func grpcErrorToHTTP(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return domain.ThrowError(http.StatusInternalServerError, err)
	}

	var statusCode int
	switch st.Code() {
	case codes.InvalidArgument:
		statusCode = http.StatusBadRequest
	case codes.NotFound:
		statusCode = http.StatusNotFound
	case codes.Unauthenticated:
		statusCode = http.StatusUnauthorized
	case codes.PermissionDenied:
		statusCode = http.StatusForbidden
	case codes.Unavailable:
		statusCode = http.StatusServiceUnavailable
	default:
		statusCode = http.StatusInternalServerError
	}

	return domain.ThrowError(statusCode, errors.New(st.Message()))
}

func newPayload(in map[string]interface{}) (*anypb.Any, error) {
	raw, err := json.Marshal(in)
	if err != nil {
		return nil, err
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

func presenceFromMap(data map[string]interface{}) *domain.Presence {
	if data == nil {
		return &domain.Presence{}
	}
	return &domain.Presence{
		ID:            toString(data["id"]),
		UserID:        toString(data["user_id"]),
		InstitutionID: toString(data["institution_id"]),
		Status:        toString(data["status"]),
		Note:          toString(data["note"]),
		CreatedAt:     toString(data["created_at"]),
		UpdatedAt:     toString(data["updated_at"]),
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
