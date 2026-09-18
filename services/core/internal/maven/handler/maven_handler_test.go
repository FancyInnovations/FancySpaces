package handler_test

import (
	"bytes"
	"context"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fancyinnovations/fancyspaces/core/internal/maven"
	fakeMavenDB "github.com/fancyinnovations/fancyspaces/core/internal/maven/database/fake"
	"github.com/fancyinnovations/fancyspaces/core/internal/maven/filestorage/memory"
	"github.com/fancyinnovations/fancyspaces/core/internal/maven/handler"
	"github.com/fancyinnovations/fancyspaces/core/internal/maven/javadoccache"
	spacesStore "github.com/fancyinnovations/fancyspaces/core/internal/spaces"
	fakeSpacesDB "github.com/fancyinnovations/fancyspaces/core/internal/spaces/database/fake"
	"github.com/fancyinnovations/fancyspaces/integrations/idp-go-sdk/idp"
	"github.com/fancyinnovations/fancyspaces/integrations/spaces-go-sdk/spaces"
)

func setupTestServer(t *testing.T) (*http.ServeMux, *maven.Store) {
	t.Helper()

	mdb := fakeMavenDB.New()
	sdb := fakeSpacesDB.New()
	memStorage := memory.NewStorage()
	jdc := javadoccache.NewService()
	spStore := spacesStore.New(spacesStore.Configuration{
		DB: sdb,
	})

	mStore := maven.New(maven.Configuration{
		Spaces:       spStore,
		DB:           mdb,
		FileStore:    memStorage,
		FileCache:    memStorage,
		JavadocCache: jdc,
	})

	testUser := &idp.User{
		ID:       "user-123",
		Verified: true,
		IsActive: true,
	}

	testSpace := &spaces.Space{
		ID:      "test-space",
		Slug:    "test-space",
		Title:   "Test Space",
		Status:  spaces.StatusApproved,
		Creator: testUser.ID,
		MavenRepositorySettings: spaces.MavenRepositorySettings{
			Enabled: true,
		},
		Members: []spaces.Member{
			{
				UserID: testUser.ID,
				Role:   spaces.RoleAdmin,
			},
		},
	}
	if err := sdb.Create(testSpace); err != nil {
		t.Fatalf("failed to create test space: %v", err)
	}

	if err := mStore.CreateRepository(context.Background(), testSpace.ID, maven.Repository{
		SpaceID:   testSpace.ID,
		Name:      "releases",
		Public:    true,
		CreatedAt: time.Now(),
	}); err != nil {
		t.Fatalf("failed to create test repo: %v", err)
	}

	h := handler.New(handler.Configuration{
		Store:  mStore,
		Spaces: spStore,
		UserFromCtx: func(ctx context.Context) *idp.User {
			return testUser
		},
	})

	mux := http.NewServeMux()
	h.RegisterMavenEndpoints(mux)
	h.RegisterAPIEndpoints("/api/v1", mux)

	return mux, mStore
}

func TestSnapshotFilesUploadAndFetch(t *testing.T) {
	mux, mStore := setupTestServer(t)

	// 1. Upload POM for 1.0.0-SNAPSHOT
	pomPath := "/test-space/releases/com/example/myapp/1.0.0-SNAPSHOT/myapp-1.0.0-SNAPSHOT.pom"
	pomBody := []byte("<project><modelVersion>4.0.0</modelVersion><groupId>com.example</groupId><artifactId>myapp</artifactId><version>1.0.0-SNAPSHOT</version></project>")
	req := httptest.NewRequest(http.MethodPut, pomPath, bytes.NewReader(pomBody))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("PUT pom expected 201 Created, got %d: %s", w.Code, w.Body.String())
	}

	// 2. Upload JAR for 1.0.0-SNAPSHOT
	jarPath := "/test-space/releases/com/example/myapp/1.0.0-SNAPSHOT/myapp-1.0.0-SNAPSHOT.jar"
	jarBody := []byte("fake-jar-data-content")
	req = httptest.NewRequest(http.MethodPut, jarPath, bytes.NewReader(jarBody))
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("PUT jar expected 201 Created, got %d: %s", w.Code, w.Body.String())
	}

	// 3. Upload sources JAR for 1.0.0-SNAPSHOT
	sourcesPath := "/test-space/releases/com/example/myapp/1.0.0-SNAPSHOT/myapp-1.0.0-SNAPSHOT-sources.jar"
	sourcesBody := []byte("fake-sources-jar-content")
	req = httptest.NewRequest(http.MethodPut, sourcesPath, bytes.NewReader(sourcesBody))
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("PUT sources expected 201 Created, got %d: %s", w.Code, w.Body.String())
	}

	// Verify all 3 files exist in the artifact model
	artifact, err := mStore.GetArtifact(context.Background(), "test-space", "releases", "com.example", "myapp")
	if err != nil {
		t.Fatalf("failed to get artifact: %v", err)
	}
	v := artifact.GetVersion("1.0.0-SNAPSHOT")
	if v == nil {
		t.Fatalf("expected version 1.0.0-SNAPSHOT to exist")
	}
	if len(v.Files) != 3 {
		t.Fatalf("expected 3 files in 1.0.0-SNAPSHOT, got %d", len(v.Files))
	}

	// Fetch POM back
	req = httptest.NewRequest(http.MethodGet, pomPath, nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("GET pom expected 200 OK, got %d", w.Code)
	}
	if !bytes.Equal(w.Body.Bytes(), pomBody) {
		t.Fatalf("GET pom content mismatch")
	}

	// Fetch JAR back
	req = httptest.NewRequest(http.MethodGet, jarPath, nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("GET jar expected 200 OK, got %d", w.Code)
	}
	if !bytes.Equal(w.Body.Bytes(), jarBody) {
		t.Fatalf("GET jar content mismatch")
	}

	// Fetch Sources back
	req = httptest.NewRequest(http.MethodGet, sourcesPath, nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("GET sources expected 200 OK, got %d", w.Code)
	}
	if !bytes.Equal(w.Body.Bytes(), sourcesBody) {
		t.Fatalf("GET sources content mismatch")
	}
}

func TestMultipleSnapshotVersions(t *testing.T) {
	mux, mStore := setupTestServer(t)

	// Upload 1.0.0-SNAPSHOT
	req := httptest.NewRequest(http.MethodPut, "/test-space/releases/com/example/myapp/1.0.0-SNAPSHOT/myapp-1.0.0-SNAPSHOT.jar", bytes.NewReader([]byte("v1")))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("PUT 1.0.0-SNAPSHOT failed: %d", w.Code)
	}

	// Upload 1.1.0-SNAPSHOT
	req = httptest.NewRequest(http.MethodPut, "/test-space/releases/com/example/myapp/1.1.0-SNAPSHOT/myapp-1.1.0-SNAPSHOT.jar", bytes.NewReader([]byte("v2")))
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("PUT 1.1.0-SNAPSHOT failed: %d", w.Code)
	}

	// Check that BOTH snapshot versions exist
	artifact, err := mStore.GetArtifact(context.Background(), "test-space", "releases", "com.example", "myapp")
	if err != nil {
		t.Fatalf("failed to get artifact: %v", err)
	}
	if len(artifact.Versions) != 2 {
		t.Fatalf("expected 2 versions, got %d", len(artifact.Versions))
	}
	if artifact.GetVersion("1.0.0-SNAPSHOT") == nil {
		t.Fatalf("1.0.0-SNAPSHOT was overwritten or deleted")
	}
	if artifact.GetVersion("1.1.0-SNAPSHOT") == nil {
		t.Fatalf("1.1.0-SNAPSHOT missing")
	}

	// Fetch artifact maven-metadata.xml
	req = httptest.NewRequest(http.MethodGet, "/test-space/releases/com/example/myapp/maven-metadata.xml", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("GET maven-metadata.xml expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var meta maven.MetadataXML
	if err := xml.Unmarshal(w.Body.Bytes(), &meta); err != nil {
		t.Fatalf("failed to unmarshal metadata: %v", err)
	}
	if meta.Versioning == nil || meta.Versioning.Latest != "1.1.0-SNAPSHOT" {
		t.Fatalf("expected latest 1.1.0-SNAPSHOT, got %v", meta.Versioning)
	}
	if meta.Versioning.Release != "" {
		t.Fatalf("expected empty release for snapshot-only, got %s", meta.Versioning.Release)
	}

	// Fetch version-level maven-metadata.xml for 1.0.0-SNAPSHOT
	req = httptest.NewRequest(http.MethodGet, "/test-space/releases/com/example/myapp/1.0.0-SNAPSHOT/maven-metadata.xml", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("GET version maven-metadata.xml expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var versionMeta maven.MetadataXML
	if err := xml.Unmarshal(w.Body.Bytes(), &versionMeta); err != nil {
		t.Fatalf("failed to unmarshal version metadata: %v", err)
	}
	if versionMeta.Version != "1.0.0-SNAPSHOT" {
		t.Fatalf("expected version 1.0.0-SNAPSHOT, got %s", versionMeta.Version)
	}
	if versionMeta.Versioning == nil || versionMeta.Versioning.Snapshot == nil {
		t.Fatalf("expected snapshot info in version metadata")
	}

	// Fetch version-level metadata checksums
	req = httptest.NewRequest(http.MethodGet, "/test-space/releases/com/example/myapp/1.0.0-SNAPSHOT/maven-metadata.xml.sha1", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("GET maven-metadata.xml.sha1 expected 200, got %d", w.Code)
	}
	if len(w.Body.String()) != 40 {
		t.Fatalf("expected 40-character sha1 hash, got %q", w.Body.String())
	}

	req = httptest.NewRequest(http.MethodGet, "/test-space/releases/com/example/myapp/1.0.0-SNAPSHOT/maven-metadata.xml.md5", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("GET maven-metadata.xml.md5 expected 200, got %d", w.Code)
	}
	if len(w.Body.String()) != 32 {
		t.Fatalf("expected 32-character md5 hash, got %q", w.Body.String())
	}
}

func TestTimestampedSnapshotDeploymentAndResolution(t *testing.T) {
	mux, _ := setupTestServer(t)

	// Deploy timestamped POM and JAR
	pomPath := "/test-space/releases/com/example/demo/1.0.0-SNAPSHOT/demo-1.0.0-20240101.120000-1.pom"
	pomBody := []byte("<project><modelVersion>4.0.0</modelVersion><groupId>com.example</groupId><artifactId>demo</artifactId><version>1.0.0-SNAPSHOT</version></project>")
	req := httptest.NewRequest(http.MethodPut, pomPath, bytes.NewReader(pomBody))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("PUT timestamped POM failed: %d", w.Code)
	}

	jarPath := "/test-space/releases/com/example/demo/1.0.0-SNAPSHOT/demo-1.0.0-20240101.120000-1.jar"
	jarBody := []byte("demo-jar-content")
	req = httptest.NewRequest(http.MethodPut, jarPath, bytes.NewReader(jarBody))
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("PUT timestamped JAR failed: %d", w.Code)
	}

	// Fetch version metadata
	metaPath := "/test-space/releases/com/example/demo/1.0.0-SNAPSHOT/maven-metadata.xml"
	req = httptest.NewRequest(http.MethodGet, metaPath, nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("GET version metadata failed: %d", w.Code)
	}

	var meta maven.MetadataXML
	if err := xml.Unmarshal(w.Body.Bytes(), &meta); err != nil {
		t.Fatalf("failed to unmarshal metadata: %v", err)
	}
	if meta.Versioning == nil || meta.Versioning.Snapshot == nil {
		t.Fatalf("expected snapshot info")
	}
	if meta.Versioning.Snapshot.Timestamp != "20240101.120000" || meta.Versioning.Snapshot.BuildNumber != 1 {
		t.Fatalf("unexpected snapshot info: %+v", meta.Versioning.Snapshot)
	}

	// Fetch via generic SNAPSHOT URL
	genericJarPath := "/test-space/releases/com/example/demo/1.0.0-SNAPSHOT/demo-1.0.0-SNAPSHOT.jar"
	req = httptest.NewRequest(http.MethodGet, genericJarPath, nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("GET generic SNAPSHOT jar failed: %d", w.Code)
	}
	if !bytes.Equal(w.Body.Bytes(), jarBody) {
		t.Fatalf("GET generic SNAPSHOT jar content mismatch")
	}

	// Fetch via exact timestamped URL
	req = httptest.NewRequest(http.MethodGet, jarPath, nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("GET timestamped jar failed: %d", w.Code)
	}
	if !bytes.Equal(w.Body.Bytes(), jarBody) {
		t.Fatalf("GET timestamped jar content mismatch")
	}
}
