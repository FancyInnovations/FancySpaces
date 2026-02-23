package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/OliverSchlueter/goutils/broker"
	"github.com/OliverSchlueter/goutils/ratelimit"
	"github.com/OliverSchlueter/goutils/sitemapgen"
	"github.com/fancyinnovations/fancyspaces/core/internal/analytics"
	analyticsCache "github.com/fancyinnovations/fancyspaces/core/internal/analytics/cache"
	analyticsDatabase "github.com/fancyinnovations/fancyspaces/core/internal/analytics/database/clickhouse"
	"github.com/fancyinnovations/fancyspaces/core/internal/badges"
	"github.com/fancyinnovations/fancyspaces/core/internal/fflags"
	"github.com/fancyinnovations/fancyspaces/core/internal/frontend"
	"github.com/fancyinnovations/fancyspaces/core/internal/issues"
	mongoIssuesDB "github.com/fancyinnovations/fancyspaces/core/internal/issues/database/mongo"
	issuesHandler "github.com/fancyinnovations/fancyspaces/core/internal/issues/handler"
	"github.com/fancyinnovations/fancyspaces/core/internal/issues/issuesync"
	"github.com/fancyinnovations/fancyspaces/core/internal/maven"
	mongoMavenDB "github.com/fancyinnovations/fancyspaces/core/internal/maven/database/mongo"
	memoryMavenFileStorage "github.com/fancyinnovations/fancyspaces/core/internal/maven/filestorage/memory"
	minioMavenFileStorage "github.com/fancyinnovations/fancyspaces/core/internal/maven/filestorage/minio"
	mavenHandler "github.com/fancyinnovations/fancyspaces/core/internal/maven/handler"
	"github.com/fancyinnovations/fancyspaces/core/internal/maven/javadoccache"
	"github.com/fancyinnovations/fancyspaces/core/internal/secrets"
	mongoSecretsDB "github.com/fancyinnovations/fancyspaces/core/internal/secrets/database/mongo"
	"github.com/fancyinnovations/fancyspaces/core/internal/secrets/encrypter/aes"
	secretsHandler "github.com/fancyinnovations/fancyspaces/core/internal/secrets/handler"
	"github.com/fancyinnovations/fancyspaces/core/internal/sitemapprovider"
	"github.com/fancyinnovations/fancyspaces/core/internal/spaces"
	mongoSpacesDB "github.com/fancyinnovations/fancyspaces/core/internal/spaces/database/mongo"
	spacesHandler "github.com/fancyinnovations/fancyspaces/core/internal/spaces/handler"
	"github.com/fancyinnovations/fancyspaces/core/internal/versions"
	mongoVersionsDB "github.com/fancyinnovations/fancyspaces/core/internal/versions/database/mongo"
	memoryVersionFileStorage "github.com/fancyinnovations/fancyspaces/core/internal/versions/filestorage/memory"
	minioVersionFileStorage "github.com/fancyinnovations/fancyspaces/core/internal/versions/filestorage/minio"
	versionsHandler "github.com/fancyinnovations/fancyspaces/core/internal/versions/handler"
	"github.com/fancyinnovations/fancyspaces/integrations/idp-go-sdk/idp"
	"github.com/minio/minio-go/v7"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const apiPrefix = "/api/v1"

type Configuration struct {
	Mux      *http.ServeMux
	MavenMux *http.ServeMux

	Broker     broker.Broker
	Mongo      *mongo.Database
	ClickHouse driver.Conn
	MinIO      *minio.Client

	SecretsMasterKey []byte
}

func Start(cfg Configuration) {
	// Analytics
	aDB := analyticsDatabase.NewDB(&analyticsDatabase.Configuration{
		CH: cfg.ClickHouse,
	})
	if err := aDB.Setup(context.Background()); err != nil {
		panic(fmt.Errorf("could not setup analytics database: %w", err))
	}
	ac := analyticsCache.NewCache()
	as := analytics.New(analytics.Configuration{
		DB:    aDB,
		Cache: ac,
		GetIP: ratelimit.GetIP,
	})

	// Spaces
	spacesDB := mongoSpacesDB.NewDB(&mongoSpacesDB.Configuration{
		Mongo: cfg.Mongo,
	})
	spacesStore := spaces.New(spaces.Configuration{
		DB: spacesDB,
	})
	seedSpaces(spacesStore)
	sh := spacesHandler.New(spacesHandler.Configuration{
		Store:       spacesStore,
		UserFromCtx: idp.UserFromCtx,
		Analytics:   as,
	})
	sh.Register(apiPrefix, cfg.Mux)
	snh := spacesHandler.NewNatsHandler(spacesHandler.NatsConfiguration{
		Broker: cfg.Broker,
		Store:  spacesStore,
	})
	if err := snh.Register(); err != nil {
		panic(fmt.Errorf("could not register spaces nats handler: %w", err))
	}

	// Versions
	versionsDB := mongoVersionsDB.NewDB(&mongoVersionsDB.Configuration{
		Mongo: cfg.Mongo,
	})
	versionFileStorage := minioVersionFileStorage.NewStorage(cfg.MinIO)
	if err := versionFileStorage.Setup(context.Background()); err != nil {
		panic(fmt.Errorf("could not setup version file storage: %w", err))
	}
	versionFileCache := memoryVersionFileStorage.NewStorage()
	versionsStore := versions.New(versions.Configuration{
		DB:          versionsDB,
		FileStorage: versionFileStorage,
		FileCache:   versionFileCache,
		Analytics:   as,
	})
	vh := versionsHandler.New(versionsHandler.Configuration{
		Store:       versionsStore,
		Spaces:      spacesStore,
		Analytics:   as,
		UserFromCtx: idp.UserFromCtx,
	})
	vh.Register(apiPrefix, cfg.Mux)

	// Maven repository
	mavenDB := mongoMavenDB.NewDB(&mongoMavenDB.Configuration{
		Mongo: cfg.Mongo,
	})
	mavenFileStorage := minioMavenFileStorage.NewStorage(cfg.MinIO)
	if err := mavenFileStorage.Setup(context.Background()); err != nil {
		panic(fmt.Errorf("could not setup maven file storage: %w", err))
	}
	mavenFileCache := memoryMavenFileStorage.NewStorage()
	mavenStore := maven.New(maven.Configuration{
		Spaces:       spacesStore,
		DB:           mavenDB,
		FileStore:    mavenFileStorage,
		FileCache:    mavenFileCache,
		JavadocCache: javadoccache.NewService(),
		Analytics:    as,
	})
	seedMavenRepos(mavenStore)
	mh := mavenHandler.New(mavenHandler.Configuration{
		Store:       mavenStore,
		Spaces:      spacesStore,
		Analytics:   as,
		UserFromCtx: idp.UserFromCtx,
	})
	mh.RegisterAPIEndpoints(apiPrefix, cfg.Mux)
	mh.RegisterMavenEndpoints(cfg.MavenMux)

	// Issues
	issuesDB := mongoIssuesDB.NewDB(&mongoIssuesDB.Configuration{
		Mongo: cfg.Mongo,
	})
	issuesStore := issues.New(issues.Configuration{
		DB: issuesDB,
	})
	ih := issuesHandler.New(issuesHandler.Configuration{
		Store:       issuesStore,
		Spaces:      spacesStore,
		UserFromCtx: idp.UserFromCtx,
	})
	ih.Register(apiPrefix, cfg.Mux)

	// Secrets
	secretsDB := mongoSecretsDB.NewDB(&mongoSecretsDB.Configuration{
		Mongo: cfg.Mongo,
	})
	secretsStore := secrets.New(secrets.Configuration{
		Database:  secretsDB,
		Encrypter: &aes.AES{},
		MasterKey: cfg.SecretsMasterKey,
	})
	secH := secretsHandler.New(secretsHandler.Configuration{
		Store:       secretsStore,
		Spaces:      spacesStore,
		UserFromCtx: idp.UserFromCtx,
	})
	secH.Register(apiPrefix, cfg.Mux)
	secNatsH := secretsHandler.NewNatsHandler(secretsHandler.NatsConfiguration{
		Broker: cfg.Broker,
		Store:  secretsStore,
	})
	if err := secNatsH.Register(); err != nil {
		panic(fmt.Errorf("could not register secrets nats handler: %w", err))
	}

	// Frontend
	fh := frontend.NewHandler(frontend.Configuration{
		Files: frontend.Files,
	})
	fh.Register(cfg.Mux)

	// Sitemap
	smp := sitemapprovider.NewService(&sitemapprovider.Configuration{
		Spaces: spacesStore,
	})
	smg := sitemapgen.NewHandler(sitemapgen.Configuration{
		Provider:      smp.GenerateUrls,
		Ratelimit:     nil, // leave to default
		CacheDuration: nil, // leave to default
	})
	smg.Register(cfg.Mux)

	// Badges
	bh := badges.NewHandler(badges.Configuration{
		Spaces:    spacesStore,
		Versions:  versionsStore,
		Analytics: as,
	})
	bh.Register(apiPrefix, cfg.Mux)

	// Issue Syncer
	if !fflags.DisableIssueSyncer.IsEnabled() {
		issueSyncer := issuesync.NewService(&issuesync.Configuration{
			SpacesStore: spacesStore,
			IssuesStore: issuesStore,
		})
		issueSyncer.StartScheduler()
		go func() {
			issueSyncer.SyncIssuesForAllSpaces()
		}()
	}
}
