package streamdal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/streamdal/streamdal/libs/protos/build/go/protos"
	"github.com/streamdal/terraform-provider-streamdal/util"
)

type Streamdal struct {
	Token    string
	Client   protos.ExternalClient
	grpcConn *grpc.ClientConn
}

type Config struct {
	Address string
	Token   string
	Timeout int
}

func New(cfg *Config) (*Streamdal, error) {
	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	timeout := time.Duration(cfg.Timeout) * time.Second
	dialContext, dialCancel := context.WithTimeout(context.Background(), timeout)
	defer dialCancel()

	conn, err := grpc.DialContext(dialContext, cfg.Address, opts...)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to grpc address '%s': %s", cfg.Address, err)
	}

	return &Streamdal{
		Token:    cfg.Token,
		Client:   protos.NewExternalClient(conn),
		grpcConn: conn,
	}, nil
}

func (s *Streamdal) Close() error {
	return s.grpcConn.Close()
}

func (s *Streamdal) CreatePipeline(ctx context.Context, req *protos.CreatePipelineRequest) (*protos.CreatePipelineResponse, error) {
	md := metadata.New(map[string]string{"auth-token": s.Token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	return s.Client.CreatePipeline(ctx, req)
}

func (s *Streamdal) UpdatePipeline(ctx context.Context, req *protos.UpdatePipelineRequest) (*protos.StandardResponse, error) {
	md := metadata.New(map[string]string{"auth-token": s.Token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	return s.Client.UpdatePipeline(ctx, req)
}

func (s *Streamdal) DeletePipeline(ctx context.Context, req *protos.DeletePipelineRequest) (*protos.StandardResponse, error) {
	md := metadata.New(map[string]string{"auth-token": s.Token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	return s.Client.DeletePipeline(ctx, req)
}

func (s *Streamdal) GetPipeline(ctx context.Context, req *protos.GetPipelineRequest) (*protos.GetPipelineResponse, error) {
	md := metadata.New(map[string]string{"auth-token": s.Token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	return s.Client.GetPipeline(ctx, req)
}

// GetPipelineFilter obtains a pipeline for a data source
func (s *Streamdal) GetPipelineFilter(filters []*Filter) (map[string]interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	md := metadata.New(map[string]string{"auth-token": s.Token})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	resp, err := s.Client.GetPipelines(ctx, &protos.GetPipelinesRequest{})
	if err != nil {
		return nil, diag.FromErr(err)
	}

	respBytes, err := json.Marshal(resp.GetPipelines())
	if err != nil {
		return nil, diag.FromErr(err)
	}

	raw := map[string]interface{}{}
	if err := json.Unmarshal(respBytes, &raw); err != nil {
		return nil, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to parse response",
			Detail:   err.Error(),
		})
	}

	pipelines, moreDiags := filterJSON(raw, filters)
	if moreDiags.HasError() {
		return nil, moreDiags
	}

	if len(pipelines) < 1 {
		// No connection found using filter
		return nil, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to find pipeline",
			Detail:   "Filters: " + filterString(filters),
		})
	} else if len(pipelines) > 1 {
		// Filter must find only one connection
		return nil, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Filter returned more than one pipeline",
		})
	}

	return pipelines[0], diags
}

func (s *Streamdal) CreateNotification(ctx context.Context, req *protos.CreateNotificationRequest) (*protos.CreateNotificationResponse, error) {
	md := metadata.New(map[string]string{"auth-token": s.Token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	return s.Client.CreateNotification(ctx, req)
}

func (s *Streamdal) UpdateNotification(ctx context.Context, req *protos.UpdateNotificationRequest) (*protos.StandardResponse, error) {
	md := metadata.New(map[string]string{"auth-token": s.Token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	return s.Client.UpdateNotification(ctx, req)
}

func (s *Streamdal) DeleteNotification(ctx context.Context, req *protos.DeleteNotificationRequest) (*protos.StandardResponse, error) {
	md := metadata.New(map[string]string{"auth-token": s.Token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	return s.Client.DeleteNotification(ctx, req)
}

func (s *Streamdal) GetNotification(ctx context.Context, req *protos.GetNotificationRequest) (*protos.GetNotificationResponse, error) {
	md := metadata.New(map[string]string{"auth-token": s.Token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	return s.Client.GetNotification(ctx, req)
}

func (s *Streamdal) GetNotificationConfigFilter(filters []*Filter) (map[string]interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	md := metadata.New(map[string]string{"auth-token": s.Token})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	resp, err := s.Client.GetNotifications(ctx, &protos.GetNotificationsRequest{})
	if err != nil {
		return nil, diag.FromErr(err)
	}

	respBytes, err := json.Marshal(resp.GetNotifications())
	if err != nil {
		return nil, diag.FromErr(err)
	}

	raw := map[string]interface{}{}
	if err := json.Unmarshal(respBytes, &raw); err != nil {
		return nil, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to parse response",
			Detail:   err.Error(),
		})
	}

	notificationCfgs, moreDiags := filterJSON(raw, filters)
	if moreDiags.HasError() {
		return nil, moreDiags
	}

	if len(notificationCfgs) < 1 {
		// No notification config found using filter
		return nil, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to find notification config",
			Detail:   "Filters: " + filterString(filters),
		})
	} else if len(notificationCfgs) > 1 {
		// Filter must find only one notification config
		return nil, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Filter returned more than one notification config",
		})
	}

	return notificationCfgs[0], diags
}

func (s *Streamdal) GetAudienceFilter(filters []*Filter) (map[string]interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	md := metadata.New(map[string]string{"auth-token": s.Token})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	resp, err := s.Client.GetAll(ctx, &protos.GetAllRequest{})
	if err != nil {
		return nil, diag.FromErr(err)
	}

	respBytes, err := json.Marshal(resp.GetAudiences())
	if err != nil {
		return nil, diag.FromErr(err)
	}

	raw := map[string]interface{}{}
	if err := json.Unmarshal(respBytes, &raw); err != nil {
		return nil, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to parse response",
			Detail:   err.Error(),
		})
	}

	audiences, moreDiags := filterJSON(raw, filters)
	if moreDiags.HasError() {
		return nil, moreDiags
	}

	if len(audiences) < 1 {
		// No audience found using filter
		return nil, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to find audience",
			Detail:   "Filters: " + filterString(filters),
		})
	} else if len(audiences) > 1 {
		// Filter must find only one audience
		return nil, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Filter returned more than one audience",
		})
	}

	return audiences[0], diags
}

func (s *Streamdal) GetAudience(ctx context.Context, id string) (*protos.Audience, error) {
	md := metadata.New(map[string]string{"auth-token": s.Token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	aud := util.AudienceFromStr(id)
	if aud == nil {
		return nil, errors.New("invalid audience id")
	}

	resp, err := s.Client.GetAll(ctx, &protos.GetAllRequest{})
	if err != nil {
		return nil, err
	}

	if util.AudienceInList(aud, resp.GetAudiences()) {
		return aud, nil
	}

	return nil, errors.New("audience not found")
}

// GetPipelinesForAudience returns a list of pipeline IDs that are associated with the given audience.
// Used for obtaining pipeline assignments for the audience resource
func (s *Streamdal) GetPipelinesForAudience(ctx context.Context, aud *protos.Audience) ([]string, error) {
	pipelineIDs := make([]string, 0)

	resp, err := s.Client.GetAll(ctx, &protos.GetAllRequest{})
	if err != nil {
		return nil, err
	}

	for pipelineID, pipeline := range resp.GetPipelines() {
		for _, audience := range pipeline.Audiences {
			if util.AudienceEquals(aud, audience) {
				pipelineIDs = append(pipelineIDs, pipelineID)
			}
		}
	}

	return pipelineIDs, nil
}

func (s *Streamdal) SetPipelines(ctx context.Context, aud *protos.Audience, pipelineIDs []string) (*protos.StandardResponse, error) {
	md := metadata.New(map[string]string{"auth-token": s.Token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	return s.Client.SetPipelines(ctx, &protos.SetPipelinesRequest{
		Audience:    aud,
		PipelineIds: pipelineIDs,
	})
}

func (s *Streamdal) CreateAudience(ctx context.Context, req *protos.CreateAudienceRequest) (*protos.StandardResponse, error) {
	md := metadata.New(map[string]string{"auth-token": s.Token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	return s.Client.CreateAudience(ctx, req)
}

func (s *Streamdal) DeleteAudience(ctx context.Context, req *protos.DeleteAudienceRequest) (*protos.StandardResponse, error) {
	md := metadata.New(map[string]string{"auth-token": s.Token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	return s.Client.DeleteAudience(ctx, req)
}
