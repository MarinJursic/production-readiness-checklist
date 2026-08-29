package controlcheck

// RegistryVersion and RegistrySHA256 bind the reviewed family registry used by
// catalog/control-check-bindings.json without making this package decode that
// catalog artifact.
const (
	RegistryVersion = "0.1.0"
	RegistrySHA256  = "7448c14fd9e3c3ce7a1fa9ea28a2ce98f8d0a1a3fe55239ba1d4715c4870bc9b"
)

// The registry is deliberately closed. A catalog binding may select only a
// reviewed, versioned implementation contract. Runtime registration supplies
// code for a descriptor; it cannot add or weaken a contract.
var implementationRegistry = []Descriptor{
	{
		Family: FamilyAnalysisAdapter, ImplementationID: "prc.check.analysis-adapter@0.1",
		ImplementationDigest: "58874f628ccf411d678417c387f45cde2457575f8c9d6d3f504627a9b96cf8d3",
		Capabilities:         []Capability{{Target: TargetExecution, Authority: AuthorityExecuted}},
	},
	{
		Family: FamilyArtifactIntegrity, ImplementationID: "prc.check.artifact-integrity@0.1",
		ImplementationDigest: "b80869c071be1f64a4d131cc3a4738f5a45d6d65b63529c8031932d79d1ac9b5",
		Capabilities:         []Capability{{Target: TargetArtifact, Authority: AuthorityArtifact}},
	},
	{
		Family: FamilyCIPolicy, ImplementationID: "prc.check.ci-policy@0.1",
		ImplementationDigest: "3ba1e384d5ab180ccb2d806c13f664f0e87576d889cc821ebd96004037ea0899",
		Capabilities: []Capability{
			{Target: TargetEnvironment, Authority: AuthorityEnvironment},
			{Target: TargetRepository, Authority: AuthorityRepository},
		},
	},
	{
		Family: FamilyContainerIaC, ImplementationID: "prc.check.container-iac@0.1",
		ImplementationDigest: "d32efaf4030b652d15a08b45f2a21fd2410dbcccc9e15a53608325facfaf4338",
		Capabilities: []Capability{
			{Target: TargetEnvironment, Authority: AuthorityEnvironment},
			{Target: TargetRepository, Authority: AuthorityRepository},
		},
	},
	{
		Family: FamilyEnvironmentEvidence, ImplementationID: "prc.check.environment-evidence@0.1",
		ImplementationDigest: "a6e267998e29bb26034c8d7d75ee887df6f806118e4f89cb3d6123a51dd137fc",
		Capabilities: []Capability{
			{Target: TargetEnvironment, Authority: AuthorityEnvironment},
			{Target: TargetExternalRegistry, Authority: AuthorityExternalRegistry},
		},
	},
	{
		Family: FamilyExecutionEvidence, ImplementationID: "prc.check.execution-evidence@0.1",
		ImplementationDigest: "a75fe01ec7fb60727324ee533ec183f1b55cd8a66b39f12b6401959caa37dfca",
		Capabilities:         []Capability{{Target: TargetExecution, Authority: AuthorityExecuted}},
	},
	{
		Family: FamilyInventoryFact, ImplementationID: "prc.check.inventory-fact@0.1",
		ImplementationDigest: "e109ca581c021b1c2439e3e2dba40e470938cace7f79d03cd1f2a28327d5dc56",
		Capabilities:         []Capability{{Target: TargetRepository, Authority: AuthorityRepository}},
	},
	{
		Family: FamilyPackageMetadata, ImplementationID: "prc.check.package-metadata@0.1",
		ImplementationDigest: "d7a26047868ab897e9ae0e35a758c0ab787971d05309ad5e4a7012426b4695ff",
		Capabilities: []Capability{
			{Target: TargetArtifact, Authority: AuthorityArtifact},
			{Target: TargetExternalRegistry, Authority: AuthorityExternalRegistry},
			{Target: TargetRepository, Authority: AuthorityRepository},
		},
	},
	{
		Family: FamilySourceAST, ImplementationID: "prc.check.source-ast@0.1",
		ImplementationDigest: "6f6b32e965f38c924b7fd9600128237c6ff4c9b1d0668da56e3553b85982ebff",
		Capabilities:         []Capability{{Target: TargetRepository, Authority: AuthorityRepository}},
	},
	{
		Family: FamilyStructuredDocument, ImplementationID: "prc.check.structured-document@0.1",
		ImplementationDigest: "cd05c2c2c6cd7981722127c94b9e93e967fa72df62a99398fb0aacd77ad2d354",
		Capabilities:         []Capability{{Target: TargetRepository, Authority: AuthorityRepository}},
	},
	{
		Family: FamilyStructuredRecord, ImplementationID: "prc.check.structured-record@0.1",
		ImplementationDigest: "44d9d08b14d2aa7ae1b57c511ce9c9fa16cf232fdc127cd198580be90fbd6039",
		Capabilities:         []Capability{{Target: TargetStructuredRecord, Authority: AuthorityStructuredRecord}},
	},
}

// Implementations returns a defensive copy in stable family order.
func Implementations() []Descriptor {
	result := make([]Descriptor, len(implementationRegistry))
	for i, descriptor := range implementationRegistry {
		result[i] = cloneDescriptor(descriptor)
	}
	return result
}

// LookupImplementation finds an immutable reviewed implementation by its
// versioned identifier.
func LookupImplementation(id string) (Descriptor, bool) {
	for _, descriptor := range implementationRegistry {
		if descriptor.ImplementationID == id {
			return cloneDescriptor(descriptor), true
		}
	}
	return Descriptor{}, false
}

func cloneDescriptor(descriptor Descriptor) Descriptor {
	descriptor.Capabilities = append([]Capability(nil), descriptor.Capabilities...)
	return descriptor
}

func descriptorEqual(left, right Descriptor) bool {
	if left.Family != right.Family || left.ImplementationID != right.ImplementationID ||
		left.ImplementationDigest != right.ImplementationDigest || len(left.Capabilities) != len(right.Capabilities) {
		return false
	}
	for i := range left.Capabilities {
		if left.Capabilities[i] != right.Capabilities[i] {
			return false
		}
	}
	return true
}

func supports(descriptor Descriptor, target Target, authority Authority) bool {
	for _, capability := range descriptor.Capabilities {
		if capability.Target == target && capability.Authority == authority {
			return true
		}
	}
	return false
}
