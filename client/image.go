package client

import "encoding/json"

// ImageResource structure.
type ImageResource struct {
	Client *Client
}

// Image structure.
type Image struct {
	ID        string `json:"id,omitempty"`
	Image     string `json:"image,omitempty"`
	Type      string `json:"type,omitempty"`
	KOpts     string `json:"kOpts,omitempty"`
	BootTagID string `json:"bootTagId,omitempty"`
}

// JSON output for a image.
func (i *Image) JSON() []byte {
	b, _ := json.MarshalIndent(i, "", "  ")
	return b
}

// All images.
func (r *ImageResource) All() (*[]Image, error) {
	c := *r.Client
	j, err := c.All("/images")
	if err != nil {
		return nil, err
	}

	images := &[]Image{}
	if err := json.Unmarshal(j, images); err != nil {
		return nil, err
	}

	return images, nil
}

// Query for image.
func (r *ImageResource) Query(cond map[string]string) (*[]Image, error) {
	c := *r.Client
	j, err := c.Query("/images", cond)
	if err != nil {
		return nil, err
	}

	images := &[]Image{}
	if err := json.Unmarshal(j, images); err != nil {
		return nil, err
	}

	return images, nil
}

// Get image.
func (r *ImageResource) Get(id string) (*Image, error) {
	c := *r.Client
	j, err := c.Get("/images", id)
	if err != nil {
		return nil, err
	}

	image := &Image{}
	if err := json.Unmarshal(j, image); err != nil {
		return nil, err
	}

	return image, nil
}

// Create image.
func (r *ImageResource) Create(i *Image) (*Image, error) {
	c := *r.Client
	j, err := c.Create("/images", i)
	if err != nil {
		return nil, err
	}

	image := &Image{}
	if err := json.Unmarshal(j, image); err != nil {
		return nil, err
	}

	return image, nil
}

// Update image.
func (r *ImageResource) Update(id string, i *Image) (*Image, error) {
	c := *r.Client
	j, err := c.Update("/images/", id, i)
	if err != nil {
		return nil, err
	}

	image := &Image{}
	if err := json.Unmarshal(j, image); err != nil {
		return nil, err
	}

	return image, nil
}

// Delete image.
func (r *ImageResource) Delete(id string) (*Image, error) {
	c := *r.Client
	j, err := c.Delete("/images", id)
	if err != nil {
		return nil, err
	}

	image := &Image{}
	if err := json.Unmarshal(j, image); err != nil {
		return nil, err
	}

	return image, nil
}
