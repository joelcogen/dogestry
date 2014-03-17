package remote

import (
  docker "github.com/blake-education/go-dockerclient"

  "encoding/json"
  "fmt"
  "io/ioutil"
  "log"
  "os"
  "os/exec"
  "path/filepath"
  "strings"
)

type RsyncRemote struct {
  config RemoteConfig
  Url string
}

func NewRsyncRemote(config RemoteConfig) (*RsyncRemote, error) {
  // TODO validate

  return &RsyncRemote{
    config: config,
    Url:  config.Url.String(),
  }, nil
}

func (remote *RsyncRemote) Validate() error {
  return nil
}

func (remote *RsyncRemote) Desc() string {
  return fmt.Sprintf("rsync(%s)", remote.Url)
}

// push all of imageRoot to the remote
func (remote *RsyncRemote) Push(image, imageRoot string) error {
  log.Println("pushing rsync", remote.Url)

  return remote.rsyncTo(imageRoot, "")
}

// pull image with id into dst
func (remote *RsyncRemote) PullImageId(id ID, dst string) error {
  log.Println("pulling rsync", "images/"+id, "->", dst)

  return remote.rsyncFrom("images/"+string(id), dst)
}

func (remote *RsyncRemote) ImageFullId(id ID) (ID, error) {
  // look for an image
  imagesRoot := filepath.Join(filepath.Clean(remote.Url), "images")
  file, err := os.Open(imagesRoot)
  if err != nil {
    return "", err
  }

  names, err := file.Readdirnames(-1)
  if err != nil {
    return "", err
  }

  for _, name := range names {
    if strings.HasPrefix(name, string(id)) {
      return ID(name), nil
    }
  }

  return "", ErrNoSuchImage
}

func (remote *RsyncRemote) WalkImages(id ID, walker ImageWalkFn) error {
  return WalkImages(remote, id, walker)
}

func (remote *RsyncRemote) ResolveImageNameToId(image string) (ID, error) {
  return ResolveImageNameToId(remote, image)
}

func (remote *RsyncRemote) ParseTag(repo, tag string) (ID, error) {
  repoPath := filepath.Join(filepath.Clean(remote.Url), "repositories", repo, tag)

  if id, err := ioutil.ReadFile(repoPath); err == nil {
    return ID(id), nil
  } else if os.IsNotExist(err) {
    return "", nil
  } else {
    return "", err
  }
}

func (remote *RsyncRemote) ImageMetadata(id ID) (docker.Image, error) {
  image := docker.Image{}

  imageJson, err := ioutil.ReadFile(filepath.Join(remote.imagePath(id), "json"))
  if os.IsNotExist(err) {
    return image, ErrNoSuchImage
  } else if err != nil {
    return image, err
  }

  if err := json.Unmarshal(imageJson, &image); err != nil {
    return image, err
  }

  return image, nil
}

func (remote *RsyncRemote) rsyncTo(src, dst string) error {
  return remote.rsync(src+"/", remote.RemotePath(dst)+"/")
}

func (remote *RsyncRemote) rsyncFrom(src, dst string) error {
  return remote.rsync(remote.RemotePath(src)+"/", dst+"/")
}

func (remote *RsyncRemote) rsync(src, dst string) error {
  out, err := exec.Command("rsync", "-av", src, dst).CombinedOutput()
  if err != nil {
    return fmt.Errorf("rsync failed: %s\noutput: %s", err, string(out))
  }
  log.Println(string(out))

  return nil
}

func (remote *RsyncRemote) imagePath(id ID) string {
  return remote.RemotePath("images", string(id))
}

func (remote *RsyncRemote) RemotePath(part ...string) string {
  return strings.TrimRight(remote.Url, "/") + "/" + strings.TrimLeft(filepath.Join(part...), "/")
}

