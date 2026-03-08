package hubspot

import "fmt"

// CustomObjectService is an interface of custom object endpoints of the HubSpot API.
// It supports dynamic object types (e.g. "2-123456") and caller-defined property structures.
// Reference: https://developers.hubspot.com/docs/api/crm/crm-custom-objects
type CustomObjectService interface {
	Get(objectType string, objectID string, object interface{}, option *RequestQueryOption) (*ResponseResource, error)
	Create(objectType string, object interface{}) (*ResponseResource, error)
	Update(objectType string, objectID string, object interface{}) (*ResponseResource, error)
	Delete(objectType string, objectID string) error
	List(objectType string, option *CustomObjectListOption) (*CustomObjectListResponse, error)
	Search(objectType string, req *CustomObjectSearchRequest) (*CustomObjectSearchResponse, error)
}

// CustomObjectServiceOp handles communication with custom object endpoints of the HubSpot API.
type CustomObjectServiceOp struct {
	customObjectPath string
	client           *Client
}

var _ CustomObjectService = (*CustomObjectServiceOp)(nil)

type CustomObjectListOption struct {
	Limit                 int      `url:"limit,omitempty"`
	After                 string   `url:"after,omitempty"`
	Properties            []string `url:"properties,comma,omitempty"`
	PropertiesWithHistory []string `url:"propertiesWithHistory,comma,omitempty"`
	Associations          []string `url:"associations,comma,omitempty"`
	Archived              bool     `url:"archived,omitempty"`
	IDProperty            string   `url:"idProperty,omitempty"`
}

type CustomObjectPagingData struct {
	After string `json:"after,omitempty"`
	Link  string `json:"link,omitempty"`
}

type CustomObjectPaging struct {
	Next *CustomObjectPagingData `json:"next,omitempty"`
}

type CustomObjectRecord struct {
	ID         string                 `json:"id,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
	CreatedAt  *HsTime                `json:"createdAt,omitempty"`
	UpdatedAt  *HsTime                `json:"updatedAt,omitempty"`
	Archived   bool                   `json:"archived,omitempty"`
	ArchivedAt *HsTime                `json:"archivedAt,omitempty"`
}

type CustomObjectListResponse struct {
	Results []CustomObjectRecord `json:"results,omitempty"`
	Paging  *CustomObjectPaging  `json:"paging,omitempty"`
}

type CustomObjectSearchRequest struct {
	SearchOptions
}

type CustomObjectSearchResponse struct {
	Total   int64                `json:"total,omitempty"`
	Results []CustomObjectRecord `json:"results,omitempty"`
	Paging  *CustomObjectPaging  `json:"paging,omitempty"`
}

// Get gets a custom object.
// In order to bind the content, a structure must be specified as an argument.
// If you want to get custom fields, specify the field names in RequestQueryOption.CustomProperties.
func (s *CustomObjectServiceOp) Get(objectType string, objectID string, object interface{}, option *RequestQueryOption) (*ResponseResource, error) {
	resource := &ResponseResource{Properties: object}
	path := fmt.Sprintf("%s/%s/%s", s.customObjectPath, objectType, objectID)
	if err := s.client.Get(path, resource, option.setupProperties(nil)); err != nil {
		return nil, err
	}
	return resource, nil
}

// Create creates a custom object.
// In order to bind the created content, a structure must be specified as an argument.
func (s *CustomObjectServiceOp) Create(objectType string, object interface{}) (*ResponseResource, error) {
	req := &RequestPayload{Properties: object}
	resource := &ResponseResource{Properties: object}
	path := fmt.Sprintf("%s/%s", s.customObjectPath, objectType)
	if err := s.client.Post(path, req, resource); err != nil {
		return nil, err
	}
	return resource, nil
}

// Update updates a custom object.
// In order to bind the updated content, a structure must be specified as an argument.
func (s *CustomObjectServiceOp) Update(objectType string, objectID string, object interface{}) (*ResponseResource, error) {
	req := &RequestPayload{Properties: object}
	resource := &ResponseResource{Properties: object}
	path := fmt.Sprintf("%s/%s/%s", s.customObjectPath, objectType, objectID)
	if err := s.client.Patch(path, req, resource); err != nil {
		return nil, err
	}
	return resource, nil
}

// Delete deletes a custom object.
func (s *CustomObjectServiceOp) Delete(objectType string, objectID string) error {
	path := fmt.Sprintf("%s/%s/%s", s.customObjectPath, objectType, objectID)
	return s.client.Delete(path, nil)
}

// List gets custom objects by object type.
func (s *CustomObjectServiceOp) List(objectType string, option *CustomObjectListOption) (*CustomObjectListResponse, error) {
	resource := &CustomObjectListResponse{}
	path := fmt.Sprintf("%s/%s", s.customObjectPath, objectType)
	if err := s.client.Get(path, resource, option); err != nil {
		return nil, err
	}
	return resource, nil
}

// Search searches custom objects by object type.
func (s *CustomObjectServiceOp) Search(objectType string, req *CustomObjectSearchRequest) (*CustomObjectSearchResponse, error) {
	resource := &CustomObjectSearchResponse{}
	path := fmt.Sprintf("%s/%s/search", s.customObjectPath, objectType)
	if err := s.client.Post(path, req, resource); err != nil {
		return nil, err
	}
	return resource, nil
}
