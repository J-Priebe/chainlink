package resolver

import (
	"github.com/graph-gophers/graphql-go"

	"github.com/smartcontractkit/chainlink/core/bridges"
)

// BridgeResolver resolves the Bridge type.
type BridgeResolver struct {
	bridge bridges.BridgeType
}

func NewBridge(bridge bridges.BridgeType) *BridgeResolver {
	return &BridgeResolver{bridge: bridge}
}

func NewBridges(bridges []bridges.BridgeType) []*BridgeResolver {
	resolvers := []*BridgeResolver{}
	for _, b := range bridges {
		resolvers = append(resolvers, NewBridge(b))
	}

	return resolvers
}

// ID resolves the bridge's name.
func (r *BridgeResolver) Name() string {
	return string(r.bridge.Name)
}

// URL resolves the bridge's url.
func (r *BridgeResolver) URL() string {
	return string(r.bridge.URL.String())
}

// Confirmations resolves the bridge's url.
func (r *BridgeResolver) Confirmations() int32 {
	return int32(r.bridge.Confirmations)
}

// OutgoingToken resolves the bridge's outgoing token.
func (r *BridgeResolver) OutgoingToken() string {
	return r.bridge.OutgoingToken
}

// MinimumContractPayment resolves the bridge's minimum contract payment.
func (r *BridgeResolver) MinimumContractPayment() string {
	return r.bridge.MinimumContractPayment.String()
}

// CreatedAt resolves the bridge's created at field.
func (r *BridgeResolver) CreatedAt() graphql.Time {
	return graphql.Time{Time: r.bridge.CreatedAt}
}

// BridgePayloadResolver resolves a single bridge response
type BridgePayloadResolver struct {
	bridge bridges.BridgeType
	err    error
}

func NewBridgePayload(bridge bridges.BridgeType, err error) *BridgePayloadResolver {
	return &BridgePayloadResolver{
		bridge: bridge,
		err:    err,
	}
}

// ToBridge implements the Bridge union type of the payload
func (r *BridgePayloadResolver) ToBridge() (*BridgeResolver, bool) {
	if r.err == nil {
		return NewBridge(r.bridge), true
	}

	return nil, false
}

// ToBridge implements the NotFoundError union type of the payload
func (r *BridgePayloadResolver) ToNotFoundError() (*NotFoundErrorResolver, bool) {
	if r.err != nil {
		return NewNotFoundError("bridge not found"), true
	}

	return nil, false
}

// BridgesPayloadResolver resolves a page of bridges
type BridgesPayloadResolver struct {
	bridges []bridges.BridgeType
	total   int32
}

func NewBridgesPayload(bridges []bridges.BridgeType, total int32) *BridgesPayloadResolver {
	return &BridgesPayloadResolver{
		bridges: bridges,
		total:   total,
	}
}

// Results returns the bridges.
func (r *BridgesPayloadResolver) Results() []*BridgeResolver {
	return NewBridges(r.bridges)
}

// Metadata returns the pagination metadata.
func (r *BridgesPayloadResolver) Metadata() *PaginationMetadataResolver {
	return NewPaginationMetadata(r.total)
}

// CreateBridgePayloadResolver
type CreateBridgePayloadResolver struct {
	bridge        bridges.BridgeType
	incomingToken string
}

func NewCreateBridgePayload(bridge bridges.BridgeType, incomingToken string) *CreateBridgePayloadResolver {
	return &CreateBridgePayloadResolver{
		bridge:        bridge,
		incomingToken: incomingToken,
	}
}

func (r *CreateBridgePayloadResolver) ToCreateBridgeSuccess() (*CreateBridgeSuccessResolver, bool) {
	return NewCreateBridgeSuccessResolver(r.bridge, r.incomingToken), true
}

type CreateBridgeSuccessResolver struct {
	bridge        bridges.BridgeType
	incomingToken string
}

func NewCreateBridgeSuccessResolver(bridge bridges.BridgeType, incomingToken string) *CreateBridgeSuccessResolver {
	return &CreateBridgeSuccessResolver{
		bridge:        bridge,
		incomingToken: incomingToken,
	}
}

// Bridge resolves the bridge.
func (r *CreateBridgeSuccessResolver) Bridge() *BridgeResolver {
	return NewBridge(r.bridge)
}

// Token resolves the bridge's incoming token.
func (r *CreateBridgeSuccessResolver) IncomingToken() string {
	return r.incomingToken
}

type UpdateBridgePayloadResolver struct {
	bridge *bridges.BridgeType
	err    error
}

func NewUpdateBridgePayload(bridge *bridges.BridgeType, err error) *UpdateBridgePayloadResolver {
	return &UpdateBridgePayloadResolver{
		bridge: bridge,
		err:    err,
	}
}

func (r *UpdateBridgePayloadResolver) ToUpdateBridgeSuccess() (*UpdateBridgeSuccessResolver, bool) {
	if r.bridge != nil {
		return NewUpdateBridgeSuccess(*r.bridge), true
	}

	return nil, false
}

func (r *UpdateBridgePayloadResolver) ToNotFoundError() (*NotFoundErrorResolver, bool) {
	if r.err != nil {
		return NewNotFoundError("bridge not found"), true
	}

	return nil, false
}

// UpdateBridgePayloadResolver resolves
type UpdateBridgeSuccessResolver struct {
	bridge bridges.BridgeType
}

func NewUpdateBridgeSuccess(bridge bridges.BridgeType) *UpdateBridgeSuccessResolver {
	return &UpdateBridgeSuccessResolver{
		bridge: bridge,
	}
}

// Bridge resolves the success payload's bridge.
func (r *UpdateBridgeSuccessResolver) Bridge() *BridgeResolver {
	return NewBridge(r.bridge)
}
