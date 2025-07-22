package cyberark

type Account struct {
	CategoryModificationTime  int64             `json:"categoryModificationTime"`
	ID                        string            `json:"id"`
	Name                      string            `json:"name"`
	Address                   string            `json:"address"`
	PlatformID                string            `json:"platformId"`
	SafeName                  string            `json:"safeName"`
	UserName                  string            `json:"userName"`
	SecretType                SecretType        `json:"secretType"`
	PlatformAccountProperties map[string]string `json:"platformAccountProperties"`
	SecretManagement          SecretManagement  `json:"secretManagement"`
	CreatedTime               int64             `json:"createdTime"`
}

type SecretType string

const (
	SecretTypePassword SecretType = "password"
	SecretTypeKey      SecretType = "key"
)

type SecretManagement struct {
	AutomaticManagementEnabled bool                    `json:"automaticManagementEnabled"`
	ManualManagementReason     string                  `json:"manualManagementReason"`
	Status                     *SecretManagementStatus `json:"status,omitempty"`
	LastModifiedDateTime       uint64                  `json:"lastModifiedDateTime"`
	LastReconciledDateTime     *uint64                 `json:"lastReconciledDateTime"`
	LastVerifiedDateTime       *uint64                 `json:"lastVerifiedDateTime"`
}

type SecretManagementStatus string

const (
	SecretManagementStatusInProcess SecretManagementStatus = "inProcess"
	SecretManagementStatusSuccess   SecretManagementStatus = "success"
	SecretManagementStatusFailure   SecretManagementStatus = "failure"
)

type Password string
