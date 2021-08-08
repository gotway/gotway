package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type Match struct {
	Method     string `json:"method"`
	Host       string `json:"host"`
	Port       string `json:"port"`
	Path       string `json:"path"`
	PathPrefix string `json:"pathPrefix"`
}

type Service struct {
	Name       string `json:"name"`
	URL        string `json:"url"`
	HealthPath string `json:"healthPath"`
}

type Cache struct {
	TTL      int64    `json:"ttl"`
	Statuses []int    `json:"statuses"`
	Tags     []string `json:"tags"`
}

type IngressHTTPSpec struct {
	Match   Match   `json:"match"`
	Service Service `json:"service"`
	Cache   Cache   `json:"cache"`
}

type IngressHTTPStatus struct {
	IsServiceHealthy bool `json:"isServiceHealthy"`
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
