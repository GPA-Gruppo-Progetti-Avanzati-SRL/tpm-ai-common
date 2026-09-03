package agentexecution

import (
	"fmt"
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

// @tpm-schematics:start-region("top-file-section")
// @tpm-schematics:end-region("top-file-section")

func UpdateMethodsGoInfo() string {
	i := fmt.Sprintf("tpm_morphia query filter support generated for %s package on %s", "author", time.Now().String())
	return i
}

type UnsetMode int64

const (
	UnSpecified     UnsetMode = 0
	KeepCurrent               = 1
	UnsetData                 = 2
	SetData2Default           = 3
)

type UnsetOption func(uopt *UnsetOptions)

type UnsetOptions struct {
	DefaultMode UnsetMode
	OId         UnsetMode
	Domain      UnsetMode
	Site        UnsetMode
	Bid         UnsetMode
	Et          UnsetMode
	Status      UnsetMode
	BatchId     UnsetMode
	Weight      UnsetMode
	BidRef      UnsetMode
	Params      UnsetMode
	Group       UnsetMode
}

func (uo *UnsetOptions) ResolveUnsetMode(um UnsetMode) UnsetMode {
	if um == UnSpecified {
		um = uo.DefaultMode
	}

	return um
}

func WithDefaultUnsetMode(m UnsetMode) UnsetOption {
	return func(uopt *UnsetOptions) {
		uopt.DefaultMode = m
	}
}
func WithOIdUnsetMode(m UnsetMode) UnsetOption {
	return func(uopt *UnsetOptions) {
		uopt.OId = m
	}
}
func WithDomainUnsetMode(m UnsetMode) UnsetOption {
	return func(uopt *UnsetOptions) {
		uopt.Domain = m
	}
}
func WithSiteUnsetMode(m UnsetMode) UnsetOption {
	return func(uopt *UnsetOptions) {
		uopt.Site = m
	}
}
func WithBidUnsetMode(m UnsetMode) UnsetOption {
	return func(uopt *UnsetOptions) {
		uopt.Bid = m
	}
}
func WithEtUnsetMode(m UnsetMode) UnsetOption {
	return func(uopt *UnsetOptions) {
		uopt.Et = m
	}
}
func WithStatusUnsetMode(m UnsetMode) UnsetOption {
	return func(uopt *UnsetOptions) {
		uopt.Status = m
	}
}
func WithBatchIdUnsetMode(m UnsetMode) UnsetOption {
	return func(uopt *UnsetOptions) {
		uopt.BatchId = m
	}
}
func WithWeightUnsetMode(m UnsetMode) UnsetOption {
	return func(uopt *UnsetOptions) {
		uopt.Weight = m
	}
}
func WithBidRefUnsetMode(m UnsetMode) UnsetOption {
	return func(uopt *UnsetOptions) {
		uopt.BidRef = m
	}
}
func WithParamsUnsetMode(m UnsetMode) UnsetOption {
	return func(uopt *UnsetOptions) {
		uopt.Params = m
	}
}
func WithGroupUnsetMode(m UnsetMode) UnsetOption {
	return func(uopt *UnsetOptions) {
		uopt.Group = m
	}
}

type UpdateOption func(ud *UpdateDocument)
type UpdateOptions []UpdateOption

// GetUpdateDocumentFromOptions convenience method to create an update document from single updates instead of a whole object
func GetUpdateDocumentFromOptions(opts ...UpdateOption) UpdateDocument {
	ud := UpdateDocument{}
	for _, o := range opts {
		o(&ud)
	}

	return ud
}

// GetUpdateDocument
// Convenience method to create an Update Document from the values of the top fields of the object. The convenience is in the handling
// the unset because if I pass an empty struct to the update it generates an empty object anyway in the db. Handling the unset eliminates
// the issue and delete an existing value without creating an empty struct.
func GetUpdateDocument(obj *AgentExecution, opts ...UnsetOption) UpdateDocument {

	uo := &UnsetOptions{DefaultMode: KeepCurrent}
	for _, o := range opts {
		o(uo)
	}

	ud := UpdateDocument{}
	ud.setOrUnsetDomain(obj.Domain, uo.ResolveUnsetMode(uo.Domain))
	ud.setOrUnsetSite(obj.Site, uo.ResolveUnsetMode(uo.Site))
	ud.setOrUnset_bid(obj.Bid, uo.ResolveUnsetMode(uo.Bid))
	ud.setOrUnset_et(obj.Et, uo.ResolveUnsetMode(uo.Et))
	ud.setOrUnsetStatus(obj.Status, uo.ResolveUnsetMode(uo.Status))
	ud.setOrUnsetBatch_id(obj.BatchId, uo.ResolveUnsetMode(uo.BatchId))
	ud.setOrUnsetWeight(obj.Weight, uo.ResolveUnsetMode(uo.Weight))
	ud.setOrUnsetBid_ref(&obj.BidRef, uo.ResolveUnsetMode(uo.BidRef))
	ud.setOrUnsetParams(obj.Params, uo.ResolveUnsetMode(uo.Params))
	ud.setOrUnsetGroup(obj.Group, uo.ResolveUnsetMode(uo.Group))

	return ud
}

// SetOId No Remarks
func (ud *UpdateDocument) SetOId(p bson.ObjectID) *UpdateDocument {
	mName := fmt.Sprintf(OIdFieldName)
	ud.Set().Add(func() bson.E {
		return bson.E{Key: mName, Value: p}
	})
	return ud
}

// UnsetOId No Remarks
func (ud *UpdateDocument) UnsetOId() *UpdateDocument {
	mName := fmt.Sprintf(OIdFieldName)
	ud.Unset().Add(func() bson.E {
		return bson.E{Key: mName, Value: ""}
	})
	return ud
}

// setOrUnsetOId No Remarks
func (ud *UpdateDocument) setOrUnsetOId(p bson.ObjectID, um UnsetMode) {
	if !p.IsZero() {
		ud.SetOId(p)
	} else {
		switch um {
		case KeepCurrent:
		case UnsetData:
			ud.UnsetOId()
		case SetData2Default:
			ud.UnsetOId()
		}
	}
}

// @tpm-schematics:start-region("o-id-field-update-section")
// @tpm-schematics:end-region("o-id-field-update-section")

// SetDomain No Remarks
func (ud *UpdateDocument) SetDomain(p string) *UpdateDocument {
	mName := fmt.Sprintf(DomainFieldName)
	ud.Set().Add(func() bson.E {
		return bson.E{Key: mName, Value: p}
	})
	return ud
}

// UnsetDomain No Remarks
func (ud *UpdateDocument) UnsetDomain() *UpdateDocument {
	mName := fmt.Sprintf(DomainFieldName)
	ud.Unset().Add(func() bson.E {
		return bson.E{Key: mName, Value: ""}
	})
	return ud
}

// setOrUnsetDomain No Remarks
func (ud *UpdateDocument) setOrUnsetDomain(p string, um UnsetMode) {
	if p != "" {
		ud.SetDomain(p)
	} else {
		switch um {
		case KeepCurrent:
		case UnsetData:
			ud.UnsetDomain()
		case SetData2Default:
			ud.UnsetDomain()
		}
	}
}

func UpdateWithDomain(p string) UpdateOption {
	return func(ud *UpdateDocument) {
		if p != "" {
			ud.SetDomain(p)
		} else {
			ud.UnsetDomain()
		}
	}
}

// @tpm-schematics:start-region("domain-field-update-section")
// @tpm-schematics:end-region("domain-field-update-section")

// SetSite No Remarks
func (ud *UpdateDocument) SetSite(p string) *UpdateDocument {
	mName := fmt.Sprintf(SiteFieldName)
	ud.Set().Add(func() bson.E {
		return bson.E{Key: mName, Value: p}
	})
	return ud
}

// UnsetSite No Remarks
func (ud *UpdateDocument) UnsetSite() *UpdateDocument {
	mName := fmt.Sprintf(SiteFieldName)
	ud.Unset().Add(func() bson.E {
		return bson.E{Key: mName, Value: ""}
	})
	return ud
}

// setOrUnsetSite No Remarks
func (ud *UpdateDocument) setOrUnsetSite(p string, um UnsetMode) {
	if p != "" {
		ud.SetSite(p)
	} else {
		switch um {
		case KeepCurrent:
		case UnsetData:
			ud.UnsetSite()
		case SetData2Default:
			ud.UnsetSite()
		}
	}
}

func UpdateWithSite(p string) UpdateOption {
	return func(ud *UpdateDocument) {
		if p != "" {
			ud.SetSite(p)
		} else {
			ud.UnsetSite()
		}
	}
}

// @tpm-schematics:start-region("site-field-update-section")
// @tpm-schematics:end-region("site-field-update-section")

// Set_bid No Remarks
func (ud *UpdateDocument) Set_bid(p string) *UpdateDocument {
	mName := fmt.Sprintf(BidFieldName)
	ud.Set().Add(func() bson.E {
		return bson.E{Key: mName, Value: p}
	})
	return ud
}

// Unset_bid No Remarks
func (ud *UpdateDocument) Unset_bid() *UpdateDocument {
	mName := fmt.Sprintf(BidFieldName)
	ud.Unset().Add(func() bson.E {
		return bson.E{Key: mName, Value: ""}
	})
	return ud
}

// setOrUnset_bid No Remarks
func (ud *UpdateDocument) setOrUnset_bid(p string, um UnsetMode) {
	if p != "" {
		ud.Set_bid(p)
	} else {
		switch um {
		case KeepCurrent:
		case UnsetData:
			ud.Unset_bid()
		case SetData2Default:
			ud.Unset_bid()
		}
	}
}

func UpdateWith_bid(p string) UpdateOption {
	return func(ud *UpdateDocument) {
		if p != "" {
			ud.Set_bid(p)
		} else {
			ud.Unset_bid()
		}
	}
}

// @tpm-schematics:start-region("-bid-field-update-section")
// @tpm-schematics:end-region("-bid-field-update-section")

// Set_et No Remarks
func (ud *UpdateDocument) Set_et(p string) *UpdateDocument {
	mName := fmt.Sprintf(EtFieldName)
	ud.Set().Add(func() bson.E {
		return bson.E{Key: mName, Value: p}
	})
	return ud
}

// Unset_et No Remarks
func (ud *UpdateDocument) Unset_et() *UpdateDocument {
	mName := fmt.Sprintf(EtFieldName)
	ud.Unset().Add(func() bson.E {
		return bson.E{Key: mName, Value: ""}
	})
	return ud
}

// setOrUnset_et No Remarks
func (ud *UpdateDocument) setOrUnset_et(p string, um UnsetMode) {
	if p != "" {
		ud.Set_et(p)
	} else {
		switch um {
		case KeepCurrent:
		case UnsetData:
			ud.Unset_et()
		case SetData2Default:
			ud.Unset_et()
		}
	}
}

func UpdateWith_et(p string) UpdateOption {
	return func(ud *UpdateDocument) {
		if p != "" {
			ud.Set_et(p)
		} else {
			ud.Unset_et()
		}
	}
}

// @tpm-schematics:start-region("-et-field-update-section")
// @tpm-schematics:end-region("-et-field-update-section")

// SetStatus No Remarks
func (ud *UpdateDocument) SetStatus(p string) *UpdateDocument {
	mName := fmt.Sprintf(StatusFieldName)
	ud.Set().Add(func() bson.E {
		return bson.E{Key: mName, Value: p}
	})
	return ud
}

// UnsetStatus No Remarks
func (ud *UpdateDocument) UnsetStatus() *UpdateDocument {
	mName := fmt.Sprintf(StatusFieldName)
	ud.Unset().Add(func() bson.E {
		return bson.E{Key: mName, Value: ""}
	})
	return ud
}

// setOrUnsetStatus No Remarks
func (ud *UpdateDocument) setOrUnsetStatus(p string, um UnsetMode) {
	if p != "" {
		ud.SetStatus(p)
	} else {
		switch um {
		case KeepCurrent:
		case UnsetData:
			ud.UnsetStatus()
		case SetData2Default:
			ud.UnsetStatus()
		}
	}
}

func UpdateWithStatus(p string) UpdateOption {
	return func(ud *UpdateDocument) {
		if p != "" {
			ud.SetStatus(p)
		} else {
			ud.UnsetStatus()
		}
	}
}

// @tpm-schematics:start-region("status-field-update-section")
// @tpm-schematics:end-region("status-field-update-section")

// SetBatch_id No Remarks
func (ud *UpdateDocument) SetBatch_id(p string) *UpdateDocument {
	mName := fmt.Sprintf(BatchIdFieldName)
	ud.Set().Add(func() bson.E {
		return bson.E{Key: mName, Value: p}
	})
	return ud
}

// UnsetBatch_id No Remarks
func (ud *UpdateDocument) UnsetBatch_id() *UpdateDocument {
	mName := fmt.Sprintf(BatchIdFieldName)
	ud.Unset().Add(func() bson.E {
		return bson.E{Key: mName, Value: ""}
	})
	return ud
}

// setOrUnsetBatch_id No Remarks
func (ud *UpdateDocument) setOrUnsetBatch_id(p string, um UnsetMode) {
	if p != "" {
		ud.SetBatch_id(p)
	} else {
		switch um {
		case KeepCurrent:
		case UnsetData:
			ud.UnsetBatch_id()
		case SetData2Default:
			ud.UnsetBatch_id()
		}
	}
}

func UpdateWithBatch_id(p string) UpdateOption {
	return func(ud *UpdateDocument) {
		if p != "" {
			ud.SetBatch_id(p)
		} else {
			ud.UnsetBatch_id()
		}
	}
}

// @tpm-schematics:start-region("batch-id-field-update-section")
// @tpm-schematics:end-region("batch-id-field-update-section")

// SetWeight No Remarks
func (ud *UpdateDocument) SetWeight(p int32) *UpdateDocument {
	mName := fmt.Sprintf(WeightFieldName)
	ud.Set().Add(func() bson.E {
		return bson.E{Key: mName, Value: p}
	})
	return ud
}

// UnsetWeight No Remarks
func (ud *UpdateDocument) UnsetWeight() *UpdateDocument {
	mName := fmt.Sprintf(WeightFieldName)
	ud.Unset().Add(func() bson.E {
		return bson.E{Key: mName, Value: ""}
	})
	return ud
}

// setOrUnsetWeight No Remarks
func (ud *UpdateDocument) setOrUnsetWeight(p int32, um UnsetMode) {
	if p != 0 {
		ud.SetWeight(p)
	} else {
		switch um {
		case KeepCurrent:
		case UnsetData:
			ud.UnsetWeight()
		case SetData2Default:
			ud.UnsetWeight()
		}
	}
}

func UpdateWithWeight(p int32) UpdateOption {
	return func(ud *UpdateDocument) {
		if p != 0 {
			ud.SetWeight(p)
		} else {
			ud.UnsetWeight()
		}
	}
}

// @tpm-schematics:start-region("weight-field-update-section")
// @tpm-schematics:end-region("weight-field-update-section")

// SetBid_ref No Remarks
func (ud *UpdateDocument) SetBid_ref(p *BidEtPair) *UpdateDocument {
	mName := fmt.Sprintf(BidRefFieldName)
	ud.Set().Add(func() bson.E {
		return bson.E{Key: mName, Value: p}
	})
	return ud
}

// UnsetBid_ref No Remarks
func (ud *UpdateDocument) UnsetBid_ref() *UpdateDocument {
	mName := fmt.Sprintf(BidRefFieldName)
	ud.Unset().Add(func() bson.E {
		return bson.E{Key: mName, Value: ""}
	})
	return ud
}

// setOrUnsetBid_ref No Remarks - here2
func (ud *UpdateDocument) setOrUnsetBid_ref(p *BidEtPair, um UnsetMode) {
	if p != nil && !p.IsZero() {
		ud.SetBid_ref(p)
	} else {
		switch um {
		case KeepCurrent:
		case UnsetData:
			ud.UnsetBid_ref()
		case SetData2Default:
			ud.UnsetBid_ref()
		}
	}
}

func UpdateWithBid_ref(p *BidEtPair) UpdateOption {
	return func(ud *UpdateDocument) {
		if p != nil && !p.IsZero() {
			ud.SetBid_ref(p)
		} else {
			ud.UnsetBid_ref()
		}
	}
}

// @tpm-schematics:start-region("bid-ref-field-update-section")
// @tpm-schematics:end-region("bid-ref-field-update-section")

// SetParams No Remarks
func (ud *UpdateDocument) SetParams(p bson.M) *UpdateDocument {
	mName := fmt.Sprintf(ParamsFieldName)
	ud.Set().Add(func() bson.E {
		return bson.E{Key: mName, Value: p}
	})
	return ud
}

// UnsetParams No Remarks
func (ud *UpdateDocument) UnsetParams() *UpdateDocument {
	mName := fmt.Sprintf(ParamsFieldName)
	ud.Unset().Add(func() bson.E {
		return bson.E{Key: mName, Value: ""}
	})
	return ud
}

// setOrUnsetParams No Remarks
func (ud *UpdateDocument) setOrUnsetParams(p bson.M, um UnsetMode) {
	if len(p) != 0 {
		ud.SetParams(p)
	} else {
		switch um {
		case KeepCurrent:
		case UnsetData:
			ud.UnsetParams()
		case SetData2Default:
			ud.UnsetParams()
		}
	}
}

func UpdateWithParams(p bson.M) UpdateOption {
	return func(ud *UpdateDocument) {
		if len(p) != 0 {
			ud.SetParams(p)
		} else {
			ud.UnsetParams()
		}
	}
}

// @tpm-schematics:start-region("params-field-update-section")
// @tpm-schematics:end-region("params-field-update-section")

// SetGroup No Remarks
func (ud *UpdateDocument) SetGroup(p string) *UpdateDocument {
	mName := fmt.Sprintf(GroupFieldName)
	ud.Set().Add(func() bson.E {
		return bson.E{Key: mName, Value: p}
	})
	return ud
}

// UnsetGroup No Remarks
func (ud *UpdateDocument) UnsetGroup() *UpdateDocument {
	mName := fmt.Sprintf(GroupFieldName)
	ud.Unset().Add(func() bson.E {
		return bson.E{Key: mName, Value: ""}
	})
	return ud
}

// setOrUnsetGroup No Remarks
func (ud *UpdateDocument) setOrUnsetGroup(p string, um UnsetMode) {
	if p != "" {
		ud.SetGroup(p)
	} else {
		switch um {
		case KeepCurrent:
		case UnsetData:
			ud.UnsetGroup()
		case SetData2Default:
			ud.UnsetGroup()
		}
	}
}

func UpdateWithGroup(p string) UpdateOption {
	return func(ud *UpdateDocument) {
		if p != "" {
			ud.SetGroup(p)
		} else {
			ud.UnsetGroup()
		}
	}
}

// @tpm-schematics:start-region("group-field-update-section")
// @tpm-schematics:end-region("group-field-update-section")

// @tpm-schematics:start-region("bottom-file-section")
// @tpm-schematics:end-region("bottom-file-section")
