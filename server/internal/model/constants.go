package model

const (
	RoleJadeEmperor     = "jade_emperor"
	RoleYanwang         = "yanwang"
	RoleJudge           = "judge"
	RoleImpermanence    = "impermanence"
	RoleMengpo          = "mengpo"
	RoleHeavenlyGeneral = "heavenly_general"
	RoleCityGod         = "city_god"
	RoleGhostClerk      = "ghost_clerk"
	RoleAuditor         = "auditor"

	SoulAlive                = "alive"
	SoulPendingCapture       = "pending_capture"
	SoulJudging              = "judging"
	SoulHell                 = "hell"
	SoulReincarnationWaiting = "reincarnation_waiting"
	SoulReincarnated         = "reincarnated"

	CapturePending   = "pending"
	CaptureRunning   = "running"
	CaptureCompleted = "completed"
	CaptureException = "exception"
	CaptureCanceled  = "canceled"

	ReincarnationQueued      = "queued"
	ReincarnationPendingSoup = "pending_soup"
	ReincarnationApproved    = "approved"
	ReincarnationRejected    = "rejected"
	ReincarnationReviewing   = "reviewing"
	ReincarnationCompleted   = "completed"
	ReincarnationException   = "exception"

	SoupPending     = "pending"
	SoupIssued      = "issued"
	SoupFailed      = "failed"
	SoupNotRequired = "not_required"

	ApprovalPending     = "pending"
	ApprovalApproved    = "approved"
	ApprovalRejected    = "rejected"
	ApprovalTransferred = "transferred"

	WishPending  = "pending"
	WishRouted   = "routed"
	WishResolved = "resolved"
	WishRejected = "rejected"
)
