package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type Match struct {
	Method     string `json:"method"`
	Host       string `json:"host"`
	Path       string `json:"path"`
	PathPrefix string `json:"pathPrefix"`
}

type Backend struct {
	URL        string `json:"url"`
	HealthPath string `json:"healthPath"`
}

type HTTPServiceSpec struct {
	Match   Match   `json:"match"`
	Backend Backend `json:"backend"`
}

type HTTPServiceStatus struct {
	Healthy bool `json:"healthy"`
}

// +genclient
// +k8s:deepcopy-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type HTTPService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec   HTTPServiceSpec   `json:"spec"`
	Status HTTPServiceStatus `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type HTTPServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []HTTPService `json:"items"`
}
