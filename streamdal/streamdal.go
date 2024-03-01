package streamdal

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/streamdal/streamdal/libs/protos/build/go/protos"
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

	raw := make([]map[string]interface{}, 0)
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
