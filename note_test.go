package hubspot_test

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/belong-inc/go-hubspot"
)

func TestNoteServiceOp_DeleteAssociation(t *testing.T) {
	type fields struct {
		notePath string
		client   *hubspot.Client
	}
	type args struct {
		noteID string
		conf   *hubspot.AssociationConfig
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "Successfully deleted association",
			fields: fields{
				notePath: hubspot.ExportNoteBasePath,
				client: hubspot.NewMockClient(&hubspot.MockConfig{
					Status: http.StatusNoContent,
					Header: http.Header{},
				}),
			},
			args: args{
				noteID: "note001",
				conf: &hubspot.AssociationConfig{
					ToObject:   hubspot.ObjectTypeDeal,
					ToObjectID: "deal001",
					Type:       hubspot.AssociationTypeNoteToDeal,
				},
			},
			wantErr: nil,
		},
		{
			name: "Received invalid association type error",
			fields: fields{
				notePath: hubspot.ExportNoteBasePath,
				client: hubspot.NewMockClient(&hubspot.MockConfig{
					Status: http.StatusBadRequest,
					Header: http.Header{},
					Body:   []byte(`{"status":"error","message":"test is not a valid association type between notes and deals","correlationId":"correlation_id","context":{"type":["test"],"fromObjectType":["notes"],"toObjectType":["deals"]},"category":"VALIDATION_ERROR","subCategory":"crm.associations.INVALID_ASSOCIATION_TYPE"}`),
				}),
			},
			args: args{
				noteID: "note001",
				conf: &hubspot.AssociationConfig{
					ToObject:   hubspot.ObjectTypeDeal,
					ToObjectID: "deal001",
					Type:       "test",
				},
			},
			wantErr: &hubspot.APIError{
				HTTPStatusCode: http.StatusBadRequest,
				Status:         "error",
				Message:        "test is not a valid association type between notes and deals",
				CorrelationID:  "correlation_id",
				Context: hubspot.ErrContext{
					Type:           []string{"test"},
					FromObjectType: []string{"notes"},
					ToObjectType:   []string{string(hubspot.ObjectTypeDeal)},
				},
				Category:    "VALIDATION_ERROR",
				SubCategory: "crm.associations.INVALID_ASSOCIATION_TYPE",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.client.CRM.Note.DeleteAssociation(tt.args.noteID, tt.args.conf)
			if !reflect.DeepEqual(tt.wantErr, err) {
				t.Errorf("DeleteAssociation() error mismatch: want %s got %s", tt.wantErr, err)
			}
		})
	}
}
