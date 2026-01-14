package docker

import "context"

// container 관련 API

func (c *Client) ListContainers(ctx context.Context) ([]Container, error) {
	list, err := c.cli.ContainerList(ctx, types.ContainerListOptions{
		All: true,
	})
	if err != nil {
		return nil, err
	}

	result := make([]Container, 0, len(list))
	for _, v := range list {
		result = append(result, Container{
			ID:     v.ID[:12],
			Name:   v.Name[0][1:],
			Image:  v.Image,
			State:  v.State,
			Status: v.Status,
		})
	}
	return result, nil
}
