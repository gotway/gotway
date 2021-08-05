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

type IngressHTTPSpec struct {
	Match   Match   `json:"match"`
	Backend Backend `json:"backend"`
}

type IngressHTTPStatus struct {
	Healthy bool `json:"healthy"`
}

// +genclient
// +k8s:deepcopy-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type IngressHTTP struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec   IngressHTTPSpec   `json:"spec"`
	Status IngressHTTPStatus `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type IngressHTTPList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []IngressHTTP `json:"items"`
}
