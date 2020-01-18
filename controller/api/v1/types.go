package v1

type NetworkPhase string
type NetworkResourceType string

const (
	PhasePending   NetworkPhase = "Pending"
	PhaseAvailable NetworkPhase = "Available"
	PhaseBound     NetworkPhase = "Bound"
	PhaseRelease   NetworkPhase = "Release"
	PhaseFailed    NetworkPhase = "Failed"
	PhaseLost      NetworkPhase = "Lost"

	ResourceIPBlock  NetworkResourceType = "IPBlock"
	ResourceEndpoint NetworkResourceType = "Endpoint"

	// AnnBindCompleted Annotation applies to EndpointClaim. It indicates that the lifecycle
	// of the EC has passed through the initial setup. This information changes how
	// we interpret some observations of the state of the objects. Value of this
	// Annotation does not matter.
	AnnBindCompleted     = "network.crd.firemiles.top/bind-completed"
	AnnBoundByController = "network.crd.firemiles.top/bound-by-controller"
)

// SubnetCIDR define subnet name and cidr block allocated
type SubnetCIDR struct {
	Name string `json:"subnet"`
	CIDR string `json:"cidr"`
}

// SubnetSlice defines multiple subnets for one endpoint
type SubnetSlice []SubnetCIDR

// PodInfo
type PodInfo struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}
