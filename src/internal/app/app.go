package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/OliverSchlueter/goutils/ratelimit"
	"github.com/fancyinnovations/fancyspaces/internal/analytics"
	analyticsCache "github.com/fancyinnovations/fancyspaces/internal/analytics/cache"
	analyticsDatabase "github.com/fancyinnovations/fancyspaces/internal/analytics/database/clickhouse"
	"github.com/fancyinnovations/fancyspaces/internal/auth"
	"github.com/fancyinnovations/fancyspaces/internal/spaces"
	fakeSpacesDB "github.com/fancyinnovations/fancyspaces/internal/spaces/database/fake"
	spacesHandler "github.com/fancyinnovations/fancyspaces/internal/spaces/handler"
	"github.com/fancyinnovations/fancyspaces/internal/versions"
	mongoVersionsDB "github.com/fancyinnovations/fancyspaces/internal/versions/database/mongo"
	memoryVersionFileStorage "github.com/fancyinnovations/fancyspaces/internal/versions/filestorage/memory"
	minioVersionFileStorage "github.com/fancyinnovations/fancyspaces/internal/versions/filestorage/minio"
	versionsHandler "github.com/fancyinnovations/fancyspaces/internal/versions/handler"
	"github.com/minio/minio-go/v7"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const apiPrefix = "/api/v1"

type Configuration struct {
	Mux        *http.ServeMux
	Mongo      *mongo.Database
	ClickHouse driver.Conn
	MinIO      *minio.Client
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
	spacesStore := spaces.New(spaces.Configuration{
		DB: seedSpacesDB(),
	})
	sh := spacesHandler.New(spacesHandler.Configuration{
		Store:       spacesStore,
		UserFromCtx: auth.UserFromContext,
		Analytics:   as,
	})
	sh.Register(apiPrefix, cfg.Mux)

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
		UserFromCtx: auth.UserFromContext,
	})
	vh.Register(apiPrefix, cfg.Mux)
}

func seedSpacesDB() *fakeSpacesDB.DB {
	db := fakeSpacesDB.New()

	fancycore := &spaces.Space{
		ID:          "fc",
		Slug:        "fancycore",
		Title:       "FancyCore",
		Description: "Essential features every Hytale server needs. From permission management over world management to economy and more.",
		Categories:  []spaces.Category{spaces.CategoryHytalePlugin},
		IconURL:     "https://fancyinnovations.com/logos-and-banners/fancycore-logo.png",
		Status:      spaces.StatusApproved,
		CreatedAt:   time.Date(2025, 11, 20, 20, 0, 0, 0, time.UTC),
		Members: []spaces.Member{
			{
				UserID: "admin-1",
				Role:   spaces.RoleOwner,
			},
		},
	}
	if err := db.Create(fancycore); err != nil {
		panic(fmt.Errorf("could not seed spaces db: %w", err))
	}

	fancycorewebsite := &spaces.Space{
		ID:          "fcw",
		Slug:        "fancycorewebsite",
		Title:       "FancyCore Website",
		Description: "Web frontend for the FancyCore Hytale plugin.",
		Categories:  []spaces.Category{spaces.CategoryWebApp},
		IconURL:     "",
		Status:      spaces.StatusPrivate,
		CreatedAt:   time.Date(2025, 12, 5, 20, 0, 0, 0, time.UTC),
		Members: []spaces.Member{
			{
				UserID: "admin-1",
				Role:   spaces.RoleOwner,
			},
		},
	}
	if err := db.Create(fancycorewebsite); err != nil {
		panic(fmt.Errorf("could not seed spaces db: %w", err))
	}

	fancyplots := &spaces.Space{
		ID:          "fp",
		Slug:        "fancyplots",
		Title:       "FancyPlots",
		Description: "Plot plugin for Hytale servers.",
		Categories:  []spaces.Category{spaces.CategoryHytalePlugin},
		IconURL:     "",
		Status:      spaces.StatusPrivate,
		CreatedAt:   time.Date(2025, 12, 5, 20, 0, 0, 0, time.UTC),
		Members: []spaces.Member{
			{
				UserID: "admin-1",
				Role:   spaces.RoleOwner,
			},
		},
	}
	if err := db.Create(fancyplots); err != nil {
		panic(fmt.Errorf("could not seed spaces db: %w", err))
	}

	fancyconnect := &spaces.Space{
		ID:          "fcon",
		Slug:        "fancyconnect",
		Title:       "FancyConnect",
		Description: "Proxy software for Hytale server networks.",
		Categories:  []spaces.Category{spaces.CategoryHytalePlugin},
		IconURL:     "",
		Status:      spaces.StatusPrivate,
		CreatedAt:   time.Date(2025, 12, 5, 20, 0, 0, 0, time.UTC),
		Members: []spaces.Member{
			{
				UserID: "admin-1",
				Role:   spaces.RoleOwner,
			},
		},
	}
	if err := db.Create(fancyconnect); err != nil {
		panic(fmt.Errorf("could not seed spaces db: %w", err))
	}

	cityauctions := &spaces.Space{
		ID:          "ca",
		Slug:        "cityauctions",
		Title:       "CityAuctions",
		Description: "Auction house plugin for Hytale.",
		Categories:  []spaces.Category{spaces.CategoryHytalePlugin},
		IconURL:     "",
		Status:      spaces.StatusPrivate,
		CreatedAt:   time.Date(2025, 12, 5, 20, 0, 0, 0, time.UTC),
		Members: []spaces.Member{
			{
				UserID: "admin-1",
				Role:   spaces.RoleOwner,
			},
		},
	}
	if err := db.Create(cityauctions); err != nil {
		panic(fmt.Errorf("could not seed spaces db: %w", err))
	}

	citypass := &spaces.Space{
		ID:          "cp",
		Slug:        "citypass",
		Title:       "CityPass",
		Description: "Pass plugin for Hytale.",
		Categories:  []spaces.Category{spaces.CategoryHytalePlugin},
		IconURL:     "",
		Status:      spaces.StatusPrivate,
		CreatedAt:   time.Date(2025, 12, 5, 20, 0, 0, 0, time.UTC),
		Members: []spaces.Member{
			{
				UserID: "admin-1",
				Role:   spaces.RoleOwner,
			},
		},
	}
	if err := db.Create(citypass); err != nil {
		panic(fmt.Errorf("could not seed spaces db: %w", err))
	}

	cityquests := &spaces.Space{
		ID:          "cq",
		Slug:        "cityquests",
		Title:       "CityQuests",
		Description: "Quests plugin for Hytale.",
		Categories:  []spaces.Category{spaces.CategoryHytalePlugin},
		IconURL:     "",
		Status:      spaces.StatusPrivate,
		CreatedAt:   time.Date(2025, 12, 5, 20, 0, 0, 0, time.UTC),
		Members: []spaces.Member{
			{
				UserID: "admin-1",
				Role:   spaces.RoleOwner,
			},
		},
	}
	if err := db.Create(cityquests); err != nil {
		panic(fmt.Errorf("could not seed spaces db: %w", err))
	}

	fancyanalytics := &spaces.Space{
		ID:          "fa",
		Slug:        "fancyanalytics",
		Title:       "FancyAnalytics",
		Description: "Universal analytics platform especially made for the Minecraft and Hytale ecosystem. Track metrics, events and logs with ease.",
		Categories:  []spaces.Category{spaces.CategoryWebApp},
		IconURL:     "https://fancyinnovations.com/logos-and-banners/fancyanalytics-logo.png",
		Status:      spaces.StatusApproved,
		CreatedAt:   time.Date(2024, 1, 25, 20, 0, 0, 0, time.UTC),
		Members: []spaces.Member{
			{
				UserID: "admin-1",
				Role:   spaces.RoleOwner,
			},
		},
	}
	if err := db.Create(fancyanalytics); err != nil {
		panic(fmt.Errorf("could not seed spaces db: %w", err))
	}

	fancyverteiler := &spaces.Space{
		ID:          "fv",
		Slug:        "fancyverteiler",
		Title:       "FancyVerteiler",
		Description: "Tool to deploy Hytale and Minecraft plugins to multiple platforms via GitHub actions.",
		Categories:  []spaces.Category{spaces.CategoryOther},
		IconURL:     "https://fancyinnovations.com/logos-and-banners/fancyverteiler-logo.png",
		Status:      spaces.StatusApproved,
		CreatedAt:   time.Date(2025, 12, 2, 20, 0, 0, 0, time.UTC),
		Members: []spaces.Member{
			{
				UserID: "admin-1",
				Role:   spaces.RoleOwner,
			},
		},
	}
	if err := db.Create(fancyverteiler); err != nil {
		panic(fmt.Errorf("could not seed spaces db: %w", err))
	}

	return db
}
