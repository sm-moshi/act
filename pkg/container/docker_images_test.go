package container

import (
	"context"
	"fmt"
	"io"
	"runtime"
	"testing"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func TestImageExistsLocally(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	ctx := context.Background()
	// to help make this test reliable and not flaky, we need to have
	// an image that will exist, and onew that won't exist

	// Test if image exists with specific tag
	invalidImageTag, err := ImageExistsLocally(ctx, "library/alpine:this-random-tag-will-never-exist", "linux/amd64")
	assert.Nil(t, err)
	assert.Equal(t, false, invalidImageTag)

	// Test if image exists with specific architecture (image platform)
	invalidImagePlatform, err := ImageExistsLocally(ctx, "alpine:latest", "windows/amd64")
	assert.Nil(t, err)
	assert.Equal(t, false, invalidImagePlatform)

	// pull an image
	cli, err := client.NewClientWithOpts(client.FromEnv)
	assert.Nil(t, err)
	cli.NegotiateAPIVersion(context.Background())

	// Chose alpine latest because it's so small
	// maybe we should build an image instead so that tests aren't reliable on dockerhub
	readerDefault, err := cli.ImagePull(ctx, "node:16-buster-slim", image.PullOptions{
		Platform: "linux/amd64",
	})
	assert.Nil(t, err)
	defer readerDefault.Close()
	_, err = io.ReadAll(readerDefault)
	assert.Nil(t, err)

	inspectImage, err := cli.ImageInspect(ctx, "node:16-buster-slim")
	assert.Nil(t, err)

	actualPlatform := fmt.Sprintf("%s/%s", inspectImage.Os, inspectImage.Architecture)
	imageDefaultArchExists, err := ImageExistsLocally(ctx, "node:16-buster-slim", actualPlatform)
	assert.Nil(t, err)
	assert.Equal(t, true, imageDefaultArchExists)

	mismatchedPlatform := "linux/amd64"
	if actualPlatform == mismatchedPlatform {
		mismatchedPlatform = "linux/arm64"
	}
	if runtime.GOOS == "windows" && actualPlatform == mismatchedPlatform {
		mismatchedPlatform = "windows/amd64"
	}

	imageArm64Exists, err := ImageExistsLocally(ctx, "node:16-buster-slim", mismatchedPlatform)
	assert.Nil(t, err)
	assert.Equal(t, false, imageArm64Exists)
}
