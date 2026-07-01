package handler

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fancyinnovations/fancyspaces/core/internal/maven"
	spacesStore "github.com/fancyinnovations/fancyspaces/core/internal/spaces"
	"github.com/fancyinnovations/fancyspaces/integrations/idp-go-sdk/idp"
	"github.com/fancyinnovations/fancyspaces/integrations/spaces-go-sdk/spaces"
)

// fakeMavenDB is a simple in-memory implementation of maven.DB for tests.
type fakeMavenDB struct {
	repos     map[string]maven.Repository
	artifacts map[string]maven.Artifact
}

func newFakeMavenDB() *fakeMavenDB {
	return &fakeMavenDB{repos: map[string]maven.Repository{}, artifacts: map[string]maven.Artifact{}}
}

func repoKey(spaceID, name string) string { return spaceID + ":" + name }
func artKey(spaceID, repoName, group, id string) string {
	return spaceID + ":" + repoName + ":" + group + ":" + id
}

func (f *fakeMavenDB) GetRepository(_ context.Context, spaceID, repoName string) (*maven.Repository, error) {
	r, ok := f.repos[repoKey(spaceID, repoName)]
	if !ok {
		return nil, maven.ErrRepositoryNotFound
	}
	return &r, nil
}
func (f *fakeMavenDB) GetRepositories(_ context.Context, spaceID string) ([]maven.Repository, error) {
	var out []maven.Repository
	for k, v := range f.repos {
		if k[:len(spaceID)+1] == spaceID+":" {
			out = append(out, v)
		}
	}
	return out, nil
}
func (f *fakeMavenDB) CreateRepository(_ context.Context, repo maven.Repository) error {
	f.repos[repoKey(repo.SpaceID, repo.Name)] = repo
	return nil
}
func (f *fakeMavenDB) UpdateRepository(_ context.Context, repo maven.Repository) error {
	f.repos[repoKey(repo.SpaceID, repo.Name)] = repo
	return nil
}
func (f *fakeMavenDB) DeleteRepository(_ context.Context, spaceID, repoName string) error {
	delete(f.repos, repoKey(spaceID, repoName))
	return nil
}

func (f *fakeMavenDB) GetArtifact(_ context.Context, spaceID, repoName, groupID, artifactID string) (*maven.Artifact, error) {
	k := artKey(spaceID, repoName, groupID, artifactID)
	a, ok := f.artifacts[k]
	if !ok {
		return nil, maven.ErrArtifactNotFound
	}
	return &a, nil
}
func (f *fakeMavenDB) GetArtifacts(_ context.Context, spaceID, repoName string) ([]maven.Artifact, error) {
	var out []maven.Artifact
	for k, v := range f.artifacts {
		prefix := spaceID + ":" + repoName + ":"
		if len(k) >= len(prefix) && k[:len(prefix)] == prefix {
			out = append(out, v)
		}
	}
	return out, nil
}
func (f *fakeMavenDB) CreateArtifact(_ context.Context, spaceID, repoName string, artifact maven.Artifact) error {
	k := artKey(spaceID, repoName, artifact.Group, artifact.ID)
	artifact.SpaceID = spaceID
	artifact.Repository = repoName
	f.artifacts[k] = artifact
	return nil
}
func (f *fakeMavenDB) UpdateArtifact(_ context.Context, spaceID, repoName string, artifact maven.Artifact) error {
	k := artKey(spaceID, repoName, artifact.Group, artifact.ID)
	f.artifacts[k] = artifact
	return nil
}
func (f *fakeMavenDB) DeleteArtifact(_ context.Context, spaceID, repoName, groupID, artifactID string) error {
	k := artKey(spaceID, repoName, groupID, artifactID)
	delete(f.artifacts, k)
	return nil
}

// fakeFileStore implements maven.FileStorage using a map keyed by a simple path
type fakeFileStore struct {
	files map[string][]byte
}

func newFakeFileStore() *fakeFileStore { return &fakeFileStore{files: map[string][]byte{}} }
func pathKey(spaceID, repoName, groupPath, artifactID, version, fileName string) string {
	return spaceID + "/" + repoName + "/" + groupPath + "/" + artifactID + "/" + version + "/" + fileName
}
func (f *fakeFileStore) UploadArtifactFile(_ context.Context, spaceID, repoName, groupPath, artifactID, version, fileName string, data []byte) error {
	f.files[pathKey(spaceID, repoName, groupPath, artifactID, version, fileName)] = append([]byte(nil), data...)
	return nil
}
func (f *fakeFileStore) DownloadArtifactFile(_ context.Context, spaceID, repoName, groupPath, artifactID, version, fileName string) ([]byte, error) {
	d, ok := f.files[pathKey(spaceID, repoName, groupPath, artifactID, version, fileName)]
	if !ok {
		return nil, fmt.Errorf("file not found")
	}
	return append([]byte(nil), d...), nil
}
func (f *fakeFileStore) DeleteArtifactFile(_ context.Context, spaceID, repoName, groupPath, artifactID, version, fileName string) error {
	delete(f.files, pathKey(spaceID, repoName, groupPath, artifactID, version, fileName))
	return nil
}

// fakeSpacesDB implements spacesStore.DB minimal methods used by spaces.Store.Get
type fakeSpacesDB struct{}

func (f *fakeSpacesDB) GetByID(id string) (*spaces.Space, error) {
	return &spaces.Space{ID: id, Creator: "user1", Status: spaces.StatusApproved, MavenRepositorySettings: spaces.MavenRepositorySettings{Enabled: true}}, nil
}
func (f *fakeSpacesDB) GetBySlug(slug string) (*spaces.Space, error) {
	return &spaces.Space{ID: slug, Creator: "user1", Status: spaces.StatusApproved, MavenRepositorySettings: spaces.MavenRepositorySettings{Enabled: true}}, nil
}
func (f *fakeSpacesDB) GetForCreator(userID string) ([]spaces.Space, error)    { return nil, nil }
func (f *fakeSpacesDB) GetForCategory(category string) ([]spaces.Space, error) { return nil, nil }
func (f *fakeSpacesDB) GetAll() ([]spaces.Space, error)                        { return nil, nil }
func (f *fakeSpacesDB) Create(s *spaces.Space) error                           { return nil }
func (f *fakeSpacesDB) Update(id string, s *spaces.Space) error                { return nil }
func (f *fakeSpacesDB) Delete(id string) error                                 { return nil }

// helper to build handler + stores preloaded with an artifact
func setupHandlerWithArtifacts(t *testing.T) (*Handler, *fakeMavenDB, *fakeFileStore, *spacesStore.Store) {
	mdb := newFakeMavenDB()
	fs := newFakeFileStore()

	// create repo
	repo := maven.Repository{SpaceID: "space1", Name: "repo1", Public: true}
	mdb.repos[repoKey(repo.SpaceID, repo.Name)] = repo

	// create artifact
	artifact := maven.Artifact{
		SpaceID:    "space1",
		Repository: "repo1",
		// handler's metadata parsing currently expects the stored group to include the artifact token
		Group: "com.example.demo",
		ID:    "demo",
		Versions: []*maven.ArtifactVersion{
			{Version: "1.0-20230701.123456-1", PublishedAt: time.Date(2023, 7, 1, 12, 34, 56, 0, time.UTC), Files: []*maven.ArtifactVersionFile{{Name: "demo-1.0-20230701.123456-1.jar", Size: 10}}},
			{Version: "1.0-20230702.000001-2", PublishedAt: time.Date(2023, 7, 2, 0, 0, 1, 0, time.UTC), Files: []*maven.ArtifactVersionFile{{Name: "demo-1.0-20230702.000001-2.jar", Size: 11}, {Name: "demo-1.0-20230702.000001-2-sources.jar", Size: 12}}},
			{Version: "2.0", PublishedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), Files: []*maven.ArtifactVersionFile{{Name: "demo-2.0.jar", Size: 20}}},
		},
	}
	mdb.artifacts[artKey("space1", "repo1", "com.example.demo", "demo")] = artifact

	// upload files into fake file store (use group path with dots replaced by /)
	groupPath := "com/example/demo"
	_ = fs.UploadArtifactFile(context.Background(), "space1", "repo1", groupPath, "demo", "1.0-20230701.123456-1", "demo-1.0-20230701.123456-1.jar", []byte("olddata1"))
	_ = fs.UploadArtifactFile(context.Background(), "space1", "repo1", groupPath, "demo", "1.0-20230702.000001-2", "demo-1.0-20230702.000001-2.jar", []byte("snapshotjar"))
	_ = fs.UploadArtifactFile(context.Background(), "space1", "repo1", groupPath, "demo", "1.0-20230702.000001-2", "demo-1.0-20230702.000001-2-sources.jar", []byte("sources"))
	_ = fs.UploadArtifactFile(context.Background(), "space1", "repo1", groupPath, "demo", "2.0", "demo-2.0.jar", []byte("releasejar"))

	mstore := maven.New(maven.Configuration{Spaces: nil, DB: mdb, FileStore: fs, FileCache: fs, JavadocCache: nil, Analytics: nil})

	spacesS := spacesStore.New(spacesStore.Configuration{DB: &fakeSpacesDB{}})

	h := New(Configuration{Store: mstore, Spaces: spacesS, Analytics: nil, UserFromCtx: func(ctx context.Context) *idp.User { return &idp.User{ID: "user1", Verified: true, IsActive: true} }})
	// disable analytics/rate limit side-effects by ensuring analytics is nil
	return h, mdb, fs, spacesS
}

func TestMetadataEndpoints(t *testing.T) {
	h, _, _, _ := setupHandlerWithArtifacts(t)

	// artifact-level metadata
	req := httptest.NewRequest(http.MethodGet, "/space1/repo1/com/example/demo/maven-metadata.xml", nil)
	w := httptest.NewRecorder()
	h.handleFetchFile(w, req, &spaces.Space{ID: "space1", MavenRepositorySettings: spaces.MavenRepositorySettings{Enabled: true}, Status: spaces.StatusApproved}, &maven.Repository{SpaceID: "space1", Name: "repo1", Public: true})

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK for artifact metadata, got %d", res.StatusCode)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	body := buf.String()
	if !bytes.Contains([]byte(body), []byte("<latest>2.0</latest>")) {
		t.Fatalf("artifact metadata did not contain latest=2.0, body: %s", body)
	}

	// snapshot version-level metadata
	req2 := httptest.NewRequest(http.MethodGet, "/space1/repo1/com/example/demo/1.0-SNAPSHOT/maven-metadata.xml", nil)
	w2 := httptest.NewRecorder()
	h.handleFetchFile(w2, req2, &spaces.Space{ID: "space1", MavenRepositorySettings: spaces.MavenRepositorySettings{Enabled: true}, Status: spaces.StatusApproved}, &maven.Repository{SpaceID: "space1", Name: "repo1", Public: true})
	res2 := w2.Result()
	if res2.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK for snapshot metadata, got %d", res2.StatusCode)
	}
	buf2 := new(bytes.Buffer)
	buf2.ReadFrom(res2.Body)
	body2 := buf2.String()
	if !bytes.Contains([]byte(body2), []byte("<snapshot>")) || !bytes.Contains([]byte(body2), []byte("<snapshotVersions>")) {
		t.Fatalf("snapshot metadata missing expected elements, body: %s", body2)
	}
}

func TestStoreAndFetchFileThroughHandler(t *testing.T) {
	h, mdb, fs, _ := setupHandlerWithArtifacts(t)

	// prepare PUT request to upload a new version file
	putBody := []byte("newcontent")
	req := httptest.NewRequest(http.MethodPut, "/space1/repo1/com/example/demo/3.0/demo-3.0.jar", bytes.NewReader(putBody))
	w := httptest.NewRecorder()

	space := &spaces.Space{ID: "space1", Creator: "user1", MavenRepositorySettings: spaces.MavenRepositorySettings{Enabled: true}, Status: spaces.StatusApproved}
	repo := &maven.Repository{SpaceID: "space1", Name: "repo1", Public: true}

	h.handleStoreFile(w, req, space, repo)
	res := w.Result()
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d", res.StatusCode)
	}

	// verify artifact exists in DB and file is stored
	a, err := mdb.GetArtifact(context.Background(), "space1", "repo1", "com.example", "demo")
	if err != nil {
		t.Fatalf("artifact not found after upload: %v", err)
	}
	v := a.GetVersion("3.0")
	if v == nil {
		t.Fatalf("version 3.0 not present in artifact after upload")
	}
	f := v.GetFile("demo-3.0.jar")
	if f == nil {
		t.Fatalf("uploaded file metadata not present")
	}
	// check file content in fake file store
	data, err := fs.DownloadArtifactFile(context.Background(), "space1", "repo1", "com/example", "demo", "3.0", "demo-3.0.jar")
	if err != nil {
		t.Fatalf("failed to download stored file: %v", err)
	}
	if !bytes.Equal(data, putBody) {
		t.Fatalf("stored file contents mismatch: expected %s got %s", string(putBody), string(data))
	}

	// fetch the just uploaded file via handler
	reqGet := httptest.NewRequest(http.MethodGet, "/space1/repo1/com/example/demo/3.0/demo-3.0.jar", nil)
	wGet := httptest.NewRecorder()
	h.handleFetchFile(wGet, reqGet, space, repo)
	resGet := wGet.Result()
	if resGet.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK fetching uploaded file, got %d", resGet.StatusCode)
	}
	got := new(bytes.Buffer)
	got.ReadFrom(resGet.Body)
	if !bytes.Equal(got.Bytes(), putBody) {
		t.Fatalf("fetched file content mismatch: expected %s got %s", string(putBody), string(got.Bytes()))
	}
}
