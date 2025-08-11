package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +kubebuilder:storageversion
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:printcolumn:name="Account",type=string,JSONPath=`.status.account`,description="Account to retrieve from cyberark"
// +kubebuilder:printcolumn:name="Safe",type=string,JSONPath=`.status.safe`,description="Safe name containing account"
// +kubebuilder:printcolumn:name="Credential Updated",type=date,JSONPath=`.status.accountUpdate`,description="Date account/password was created/modified"
// +kubebuilder:printcolumn:name="Secret",type=string,JSONPath=`.status.secretName`,description="Kubernetes secret name(s)"
// +kubebuilder:printcolumn:name="Synced",type=date,JSONPath=`.status.lastSync`,description="When the account was last synced with secret"
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`,description="When the resource was created"

// +kubebuilder:subresource:status
type CyberArk struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec   CyberArkSpec   `json:"spec"`
	Status CyberArkStatus `json:"status"`
}

type CyberArkStatus struct {
	SecretHashes  map[string]string `json:"secretHashes"`
	LastSync      metav1.Time       `json:"lastSync"`
	AccountUpdate metav1.Time       `json:"accountUpdate"`
	SecretName    string            `json:"secretName"`
	Account       string            `json:"account"`
	Safe          string            `json:"safe"`
}

type CyberArkSpec struct {
	Target  CyberArkTarget  `json:"target"`
	Account CyberArkAccount `json:"account"`
}

type CyberArkTarget struct {
	// +optional
	Secret CyberArkTargetSecret `json:"secret,omitempty"`

	// Add future targets here
}

type CyberArkTargetSecret struct {
	Name              string            `json:"name"`
	UsernameKeys      []string          `json:"usernameKeys"`
	PasswordKeys      []string          `json:"passwordKeys"`
	AdditionalSecrets map[string]string `json:"additionalSecrets,omitempty"`
}

// +kubebuilder:validation:Enum=contains;startswith
type CyberArkSearchType string

const (
	SearchContains   CyberArkSearchType = "contains"
	SearchStartswith CyberArkSearchType = "startswith"
)

type CyberArkAccount struct {
	ID         string             `json:"id,omitempty"`
	UserName   string             `json:"userName,omitempty"`
	SafeName   string             `json:"safeName,omitempty"`
	Search     string             `json:"search,omitempty"`
	SearchType CyberArkSearchType `json:"searchType,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type CyberArkList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []CyberArk `json:"items"`
}
