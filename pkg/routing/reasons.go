package routing

// Outcome tokens live in OutcomeKind (types.go); only spec-named reasons
// belong here, including the explicitly named no_eligible_candidates reason.
// ReasonCode is the shared vocabulary. Adding a code changes a frozen contract.
type ReasonCode string

const (
	ReasonNoEligibleCandidates              ReasonCode = "no_eligible_candidates"
	ReasonCandidateDependentSettingMismatch ReasonCode = "candidate_dependent_setting_mismatch"
	ReasonSelectedCandidateUnavailable      ReasonCode = "selected_candidate_unavailable"
	ReasonQuotaScopeUnmapped                ReasonCode = "quota_scope_unmapped"
	ReasonQuotaNoApplicableWindow           ReasonCode = "quota_no_applicable_window"
	ReasonQuotaStale                        ReasonCode = "quota_stale"
	ReasonQuotaExpired                      ReasonCode = "quota_expired"
	ReasonQuotaInvalid                      ReasonCode = "quota_invalid"
	ReasonQuotaPartial                      ReasonCode = "quota_partial"
	ReasonQuotaUnknown                      ReasonCode = "quota_unknown"
	ReasonQuotaOver                         ReasonCode = "quota_over"
	ReasonQuotaReserve                      ReasonCode = "quota_reserve"
	ReasonQuotaExpiring                     ReasonCode = "quota_expiring"
	ReasonQuotaOnPace                       ReasonCode = "quota_on_pace"
	ReasonQuotaConserve                     ReasonCode = "quota_conserve"
	ReasonInflightSpread                    ReasonCode = "inflight_spread"
	ReasonInflightUnknown                   ReasonCode = "inflight_unknown"
	ReasonReserveBreached                   ReasonCode = "reserve_breached"
	ReasonHeadroomGroupsOverlap             ReasonCode = "headroom_groups_overlap"
)
