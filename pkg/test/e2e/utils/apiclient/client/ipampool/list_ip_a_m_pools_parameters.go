// Code generated by go-swagger; DO NOT EDIT.

package ipampool

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

// NewListIPAMPoolsParams creates a new ListIPAMPoolsParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewListIPAMPoolsParams() *ListIPAMPoolsParams {
	return &ListIPAMPoolsParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewListIPAMPoolsParamsWithTimeout creates a new ListIPAMPoolsParams object
// with the ability to set a timeout on a request.
func NewListIPAMPoolsParamsWithTimeout(timeout time.Duration) *ListIPAMPoolsParams {
	return &ListIPAMPoolsParams{
		timeout: timeout,
	}
}

// NewListIPAMPoolsParamsWithContext creates a new ListIPAMPoolsParams object
// with the ability to set a context for a request.
func NewListIPAMPoolsParamsWithContext(ctx context.Context) *ListIPAMPoolsParams {
	return &ListIPAMPoolsParams{
		Context: ctx,
	}
}

// NewListIPAMPoolsParamsWithHTTPClient creates a new ListIPAMPoolsParams object
// with the ability to set a custom HTTPClient for a request.
func NewListIPAMPoolsParamsWithHTTPClient(client *http.Client) *ListIPAMPoolsParams {
	return &ListIPAMPoolsParams{
		HTTPClient: client,
	}
}

/* ListIPAMPoolsParams contains all the parameters to send to the API endpoint
   for the list IP a m pools operation.

   Typically these are written to a http.Request.
*/
type ListIPAMPoolsParams struct {

	// SeedName.
	SeedName string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the list IP a m pools params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *ListIPAMPoolsParams) WithDefaults() *ListIPAMPoolsParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the list IP a m pools params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *ListIPAMPoolsParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the list IP a m pools params
func (o *ListIPAMPoolsParams) WithTimeout(timeout time.Duration) *ListIPAMPoolsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the list IP a m pools params
func (o *ListIPAMPoolsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the list IP a m pools params
func (o *ListIPAMPoolsParams) WithContext(ctx context.Context) *ListIPAMPoolsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the list IP a m pools params
func (o *ListIPAMPoolsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the list IP a m pools params
func (o *ListIPAMPoolsParams) WithHTTPClient(client *http.Client) *ListIPAMPoolsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the list IP a m pools params
func (o *ListIPAMPoolsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithSeedName adds the seedName to the list IP a m pools params
func (o *ListIPAMPoolsParams) WithSeedName(seedName string) *ListIPAMPoolsParams {
	o.SetSeedName(seedName)
	return o
}

// SetSeedName adds the seedName to the list IP a m pools params
func (o *ListIPAMPoolsParams) SetSeedName(seedName string) {
	o.SeedName = seedName
}

// WriteToRequest writes these params to a swagger request
func (o *ListIPAMPoolsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param seed_name
	if err := r.SetPathParam("seed_name", o.SeedName); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}