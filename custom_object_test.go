package hubspot_test

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/WAS-Insurance/go-hubspot"
	"github.com/google/go-cmp/cmp"
)

type CustomObjectProperties struct {
	Name        *hubspot.HsStr `json:"name,omitempty"`
	CustomField *hubspot.HsStr `json:"custom_field,omitempty"`
}

func TestCustomObjectServiceOp_Create(t *testing.T) {
	object := &CustomObjectProperties{
		Name:        hubspot.NewString("Custom Object Name"),
		CustomField: hubspot.NewString("Custom Value"),
	}

	tests := []struct {
		name    string
		client  *hubspot.Client
		want    *hubspot.ResponseResource
		wantErr error
	}{
		{
			name: "Successfully create a custom object",
			client: hubspot.NewMockClient(&hubspot.MockConfig{
				Status: http.StatusCreated,
				Header: http.Header{},
				Body:   []byte(`{"id":"1001","archived":false,"properties":{"name":"Custom Object Name","custom_field":"Custom Value"},"createdAt":"2019-10-30T03:30:17.883Z","updatedAt":"2019-12-07T16:50:06.678Z"}`),
			}),
			want: &hubspot.ResponseResource{
				ID:       "1001",
				Archived: false,
				Properties: &CustomObjectProperties{
					Name:        hubspot.NewString("Custom Object Name"),
					CustomField: hubspot.NewString("Custom Value"),
				},
				CreatedAt: &createdAt,
				UpdatedAt: &updatedAt,
			},
			wantErr: nil,
		},
		{
			name: "Received invalid request",
			client: hubspot.NewMockClient(&hubspot.MockConfig{
				Status: http.StatusBadRequest,
				Header: http.Header{},
				Body:   []byte(`{"message":"Invalid input","correlationId":"c-id","category":"VALIDATION_ERROR"}`),
			}),
			want: nil,
			wantErr: &hubspot.APIError{
				HTTPStatusCode: http.StatusBadRequest,
				Message:        "Invalid input",
				CorrelationID:  "c-id",
				Category:       "VALIDATION_ERROR",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.client.CRM.CustomObject.Create("2-123456", object)
			if !reflect.DeepEqual(tt.wantErr, err) {
				t.Errorf("Create() error mismatch: want %v got %v", tt.wantErr, err)
				return
			}
			if diff := cmp.Diff(tt.want, got, cmpTimeOption); diff != "" {
				t.Errorf("Create() response mismatch (-want +got):%s", diff)
			}
		})
	}
}

func TestCustomObjectServiceOp_Get(t *testing.T) {
	type fields struct {
		client *hubspot.Client
	}
	type args struct {
		objectID string
		object   interface{}
		option   *hubspot.RequestQueryOption
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *hubspot.ResponseResource
		wantErr error
	}{
		{
			name: "Successfully get a custom object",
			fields: fields{
				client: hubspot.NewMockClient(&hubspot.MockConfig{
					Status: http.StatusOK,
					Header: http.Header{},
					Body:   []byte(`{"id":"1001","archived":false,"properties":{"name":"Custom Object Name","custom_field":"Custom Value"},"createdAt":"2019-10-30T03:30:17.883Z","updatedAt":"2019-12-07T16:50:06.678Z"}`),
				}),
			},
			args: args{
				objectID: "1001",
				object:   &CustomObjectProperties{},
				option: &hubspot.RequestQueryOption{
					CustomProperties: []string{"name", "custom_field"},
				},
			},
			want: &hubspot.ResponseResource{
				ID:       "1001",
				Archived: false,
				Properties: &CustomObjectProperties{
					Name:        hubspot.NewString("Custom Object Name"),
					CustomField: hubspot.NewString("Custom Value"),
				},
				CreatedAt: &createdAt,
				UpdatedAt: &updatedAt,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fields.client.CRM.CustomObject.Get("2-123456", tt.args.objectID, tt.args.object, tt.args.option)
			if !reflect.DeepEqual(tt.wantErr, err) {
				t.Errorf("Get() error mismatch: want %v got %v", tt.wantErr, err)
				return
			}
			if diff := cmp.Diff(tt.want, got, cmpTimeOption); diff != "" {
				t.Errorf("Get() response mismatch (-want +got):%s", diff)
			}
		})
	}
}

func TestCustomObjectServiceOp_Get_QueryOption(t *testing.T) {
	roundTripper := roundTripFunc(func(req *http.Request) *http.Response {
		if req.Method != http.MethodGet {
			t.Errorf("unexpected method: got %s want %s", req.Method, http.MethodGet)
		}
		if req.URL.Path != "/crm/v3/objects/2-123456/1001" {
			t.Errorf("unexpected path: got %s", req.URL.Path)
		}

		query := req.URL.Query()
		if got := query.Get("properties"); got != "name,custom_field" {
			t.Errorf("unexpected properties query: got %s want %s", got, "name,custom_field")
		}
		if got := query.Get("associations"); got != "contacts,companies" {
			t.Errorf("unexpected associations query: got %s want %s", got, "contacts,companies")
		}
		if got := query.Get("idProperty"); got != "external_id" {
			t.Errorf("unexpected idProperty query: got %s want %s", got, "external_id")
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body: io.NopCloser(bytes.NewBuffer([]byte(
				`{"id":"1001","archived":false,"properties":{"name":"Custom Object Name","custom_field":"Custom Value"}}`,
			))),
			Header: http.Header{},
		}
	})

	httpClient := &http.Client{Transport: roundTripper}
	client, err := hubspot.NewClient(hubspot.SetPrivateAppToken("token"), hubspot.WithHTTPClient(httpClient))
	if err != nil {
		t.Fatalf("NewClient() unexpected error: %v", err)
	}

	_, err = client.CRM.CustomObject.Get("2-123456", "1001", &CustomObjectProperties{}, &hubspot.RequestQueryOption{
		CustomProperties: []string{"name", "custom_field"},
		Associations:     []string{"contacts", "companies"},
		IDProperty:       "external_id",
	})
	if err != nil {
		t.Fatalf("Get() unexpected error: %v", err)
	}
}

func TestCustomObjectServiceOp_Update(t *testing.T) {
	object := &CustomObjectProperties{
		Name:        hubspot.NewString("Updated Object Name"),
		CustomField: hubspot.NewString("Updated Value"),
	}

	client := hubspot.NewMockClient(&hubspot.MockConfig{
		Status: http.StatusOK,
		Header: http.Header{},
		Body:   []byte(`{"id":"1001","archived":false,"properties":{"name":"Updated Object Name","custom_field":"Updated Value"},"createdAt":"2019-10-30T03:30:17.883Z","updatedAt":"2019-12-07T16:50:06.678Z"}`),
	})

	got, err := client.CRM.CustomObject.Update("2-123456", "1001", object)
	if err != nil {
		t.Fatalf("Update() unexpected error: %v", err)
	}

	want := &hubspot.ResponseResource{
		ID:       "1001",
		Archived: false,
		Properties: &CustomObjectProperties{
			Name:        hubspot.NewString("Updated Object Name"),
			CustomField: hubspot.NewString("Updated Value"),
		},
		CreatedAt: &createdAt,
		UpdatedAt: &updatedAt,
	}
	if diff := cmp.Diff(want, got, cmpTimeOption); diff != "" {
		t.Errorf("Update() response mismatch (-want +got):%s", diff)
	}
}

func TestCustomObjectServiceOp_Delete(t *testing.T) {
	client := hubspot.NewMockClient(&hubspot.MockConfig{
		Status: http.StatusNoContent,
		Header: http.Header{},
		Body:   []byte(``),
	})

	if err := client.CRM.CustomObject.Delete("2-123456", "1001"); err != nil {
		t.Fatalf("Delete() unexpected error: %v", err)
	}
}

func TestCustomObjectServiceOp_List(t *testing.T) {
	client := hubspot.NewMockClient(&hubspot.MockConfig{
		Status: http.StatusOK,
		Header: http.Header{},
		Body:   []byte(`{"results":[{"id":"1001","properties":{"name":"Object A","custom_field":"A"},"createdAt":"2019-10-30T03:30:17.883Z","updatedAt":"2019-12-07T16:50:06.678Z","archived":false}],"paging":{"next":{"after":"2","link":"https://example.com"}}}`),
	})

	got, err := client.CRM.CustomObject.List("2-123456", &hubspot.CustomObjectListOption{
		Limit:      10,
		After:      "1",
		Properties: []string{"name", "custom_field"},
	})
	if err != nil {
		t.Fatalf("List() unexpected error: %v", err)
	}

	if len(got.Results) != 1 {
		t.Fatalf("List() results length mismatch: got %d want %d", len(got.Results), 1)
	}
	if got.Results[0].ID != "1001" {
		t.Errorf("List() result id mismatch: got %s want %s", got.Results[0].ID, "1001")
	}
	if got.Results[0].Properties["name"] != "Object A" {
		t.Errorf("List() property mismatch: got %v want %v", got.Results[0].Properties["name"], "Object A")
	}
	if got.Paging == nil || got.Paging.Next == nil || got.Paging.Next.After != "2" {
		t.Errorf("List() paging mismatch: got %+v", got.Paging)
	}
}

func TestCustomObjectServiceOp_Search(t *testing.T) {
	client := hubspot.NewMockClient(&hubspot.MockConfig{
		Status: http.StatusOK,
		Header: http.Header{},
		Body:   []byte(`{"total":1,"results":[{"id":"1001","properties":{"name":"Object A","custom_field":"A"},"createdAt":"2019-10-30T03:30:17.883Z","updatedAt":"2019-12-07T16:50:06.678Z","archived":false}]}`),
	})

	got, err := client.CRM.CustomObject.Search("2-123456", &hubspot.CustomObjectSearchRequest{
		SearchOptions: hubspot.SearchOptions{
			FilterGroups: []hubspot.FilterGroup{
				{
					Filters: []hubspot.Filter{
						{
							PropertyName: "name",
							Operator:     hubspot.EQ,
							Value:        hubspot.NewString("Object A"),
						},
					},
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("Search() unexpected error: %v", err)
	}
	if got.Total != 1 {
		t.Errorf("Search() total mismatch: got %d want %d", got.Total, 1)
	}
	if len(got.Results) != 1 {
		t.Fatalf("Search() results length mismatch: got %d want %d", len(got.Results), 1)
	}
	if got.Results[0].Properties["custom_field"] != "A" {
		t.Errorf("Search() property mismatch: got %v want %v", got.Results[0].Properties["custom_field"], "A")
	}
}

type roundTripFunc func(req *http.Request) *http.Response

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func TestCustomObjectServiceOp_List_QueryOption(t *testing.T) {
	roundTripper := roundTripFunc(func(req *http.Request) *http.Response {
		if req.Method != http.MethodGet {
			t.Errorf("unexpected method: got %s want %s", req.Method, http.MethodGet)
		}
		if req.URL.Path != "/crm/v3/objects/2-123456" {
			t.Errorf("unexpected path: got %s", req.URL.Path)
		}

		query := req.URL.Query()
		assertQueryValue(t, query, "limit", "10")
		assertQueryValue(t, query, "after", "20")
		assertQueryValue(t, query, "properties", "name,custom_field")
		assertQueryValue(t, query, "propertiesWithHistory", "name")
		assertQueryValue(t, query, "associations", "contacts,companies")
		assertQueryValue(t, query, "archived", "true")
		assertQueryValue(t, query, "idProperty", "external_id")

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBuffer([]byte(`{"results":[]}`))),
			Header:     http.Header{},
		}
	})

	httpClient := &http.Client{Transport: roundTripper}
	client, err := hubspot.NewClient(hubspot.SetPrivateAppToken("token"), hubspot.WithHTTPClient(httpClient))
	if err != nil {
		t.Fatalf("NewClient() unexpected error: %v", err)
	}

	_, err = client.CRM.CustomObject.List("2-123456", &hubspot.CustomObjectListOption{
		Limit:                 10,
		After:                 "20",
		Properties:            []string{"name", "custom_field"},
		PropertiesWithHistory: []string{"name"},
		Associations:          []string{"contacts", "companies"},
		Archived:              true,
		IDProperty:            "external_id",
	})
	if err != nil {
		t.Fatalf("List() unexpected error: %v", err)
	}
}

func assertQueryValue(t *testing.T, query url.Values, key string, want string) {
	t.Helper()
	if got := query.Get(key); got != want {
		t.Errorf("unexpected %s query: got %s want %s", key, got, want)
	}
}
