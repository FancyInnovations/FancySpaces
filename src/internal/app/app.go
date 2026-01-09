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
		Summary:     "Essential features every Hytale server needs. From permission management over world management to economy and more.",
		Description: "![](https://raw.githubusercontent.com/FancyInnovations/FancyDocs/refs/heads/main/public/logos-and-banners/fancycore-banner.png)\n\n> The all-in-one core plugin for Hytale servers. From powerful permission management and moderation tools to a flexible economy with multiple currencies and much more.\n\n## Why FancyCore?\n\n- All-in-one core plugin – fewer dependencies, fewer conflicts\n- Designed for both small servers and large networks\n- Highly configurable without sacrificing performance\n- Actively maintained with a clear development roadmap\n- Built with developers in mind (clean API & extensions)\n\n## Features\n\nWith **FancyCore**, you get a wide variety of features that are essential for running a modern Hytale server.\nIt includes **80+ commands** covering countless use cases for both small community servers and large server networks.\n\nFancyCore is designed with ease of use, high performance, and extensibility in mind.\n\n### Core Feature Categories\n\n- Flexible group and permission system\n- Robust economy system\n- Easy to use placeholders\n- Chat management\n- Powerful moderation tools\n- Teleportation features\n- World management\n- Player specific features\n- Inventory utilities\n- Server statistics\n- API for developers\n\nLearn more about each feature in the documentation: https://fancyinnovations.com/docs/hytale-plugins/fancycore/\n\n### Permissions\n\nFancyCore provides a powerful and flexible permission system suitable for any server size.\n\n- Create unlimited groups with inheritance\n- Per-group and per-player permissions\n- Temporary permissions and groups\n- Prefixes, suffixes, and priorities\n- Fully configurable via commands and files\n\nPerfect for managing staff hierarchies and player ranks.\n\n### Economy\n\nA feature-rich economy system built directly into the core.\n\n- Multiple currencies\n- Player balances with full command control\n- Admin and player economy commands\n- Optional integration with shops and other plugins\n- High-performance and safe data handling\n\nWhether you run a survival server or a complex RPG economy, FancyCore has you covered.\n\n### Placeholders\n\nFancyCore includes a built-in placeholder system for maximum compatibility.\n\n- Many placeholders from every feature category \n- Easy integration with chat messages, UIs and more\n- Lightweight and fast\n\n### Chat\n\nTake full control over your server chat.\n\n- Customizable chat formats\n- Group-based prefixes and suffixes\n- Chat channels (global, staff, per rank, etc.)\n- Chat cooldowns and filters\n- Player nicknames\n- Messaging system\n- Placeholder support in chat messages\n\nKeep your chat clean, organized, and immersive.\n\n### Moderation\n\nAll the moderation tools you need in one plugin.\n\n- Kick, mute, warn, and ban commands\n- Player reports\n- Chat and command logs\n- Staff-only chat channels\n- Silent punishments\n- Full permission control\n- Clear and consistent punishment messages\n\nDesigned to make moderation fast and effective.\n\n### Teleportation\n\nComprehensive teleportation features for players and staff.\n\n- Spawn\n- Server Warps\n- Homes and multiple home support\n- Teleport requests (TPA)\n- Cooldowns and permission-based limits\n\n### Worlds\n\nManage your worlds with ease.\n\n- Create worlds with different settings and environments\n- Teleport between worlds\n- Control player interactions per world\n- Ideal for hubs, minigames, and survival worlds\n\n### Player\n\nQuality-of-life features for everyday gameplay.\n\n- Player information commands\n- Inventory and gamemode utilities\n- AFK detection\n- Custom join and leave messages\n- Player-specific settings\n\n### Inventory\n\nAdvanced inventory utilities for players and staff.\n\n- View and manage player inventories\n- Create configurable kits\n- Virtual player backpacks\n\n### Server\n\n- View server health and statistics\n- Manage plugins\n\n### API\n\nFancyCore is built to be extended.\n\n- Clean and well-documented API\n- Access permissions, economy, placeholders, and more\n- Designed for developers and plugin integrations\n- Future-proof and actively maintained",
		Categories:  []spaces.Category{spaces.CategoryHytalePlugin},
		Links: []spaces.Link{
			{Name: "source_code", URL: "https://github.com/FancyInnovations/HytalePlugins"},
			{Name: "documentation", URL: "https://fancyinnovations.com/docs/hytale-plugins/fancycore"},
			{Name: "discord", URL: "https://discord.gg/ZUgYCEJUEx"},
		},
		IconURL:   "https://fancyinnovations.com/logos-and-banners/fancycore-logo.png",
		Status:    spaces.StatusApproved,
		CreatedAt: time.Date(2025, 11, 20, 20, 0, 0, 0, time.UTC),
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
		ID:         "fcw",
		Slug:       "fancycorewebsite",
		Title:      "FancyCore Website",
		Summary:    "Web frontend for the FancyCore Hytale plugin.",
		Categories: []spaces.Category{spaces.CategoryWebApp},
		IconURL:    "",
		Status:     spaces.StatusPrivate,
		CreatedAt:  time.Date(2025, 12, 5, 20, 0, 0, 0, time.UTC),
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
		ID:         "fp",
		Slug:       "fancyplots",
		Title:      "FancyPlots",
		Summary:    "Plot plugin for Hytale servers.",
		Categories: []spaces.Category{spaces.CategoryHytalePlugin},
		IconURL:    "",
		Status:     spaces.StatusPrivate,
		CreatedAt:  time.Date(2025, 12, 5, 20, 0, 0, 0, time.UTC),
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
		ID:         "fcon",
		Slug:       "fancyconnect",
		Title:      "FancyConnect",
		Summary:    "Proxy software for Hytale server networks.",
		Categories: []spaces.Category{spaces.CategoryHytalePlugin},
		IconURL:    "",
		Status:     spaces.StatusPrivate,
		CreatedAt:  time.Date(2025, 12, 5, 20, 0, 0, 0, time.UTC),
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
		ID:         "ca",
		Slug:       "cityauctions",
		Title:      "CityAuctions",
		Summary:    "Auction house plugin for Hytale.",
		Categories: []spaces.Category{spaces.CategoryHytalePlugin},
		IconURL:    "",
		Status:     spaces.StatusPrivate,
		CreatedAt:  time.Date(2025, 12, 5, 20, 0, 0, 0, time.UTC),
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
		ID:         "cp",
		Slug:       "citypass",
		Title:      "CityPass",
		Summary:    "Pass plugin for Hytale.",
		Categories: []spaces.Category{spaces.CategoryHytalePlugin},
		IconURL:    "",
		Status:     spaces.StatusPrivate,
		CreatedAt:  time.Date(2025, 12, 5, 20, 0, 0, 0, time.UTC),
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
		ID:         "cq",
		Slug:       "cityquests",
		Title:      "CityQuests",
		Summary:    "Quests plugin for Hytale.",
		Categories: []spaces.Category{spaces.CategoryHytalePlugin},
		IconURL:    "",
		Status:     spaces.StatusPrivate,
		CreatedAt:  time.Date(2025, 12, 5, 20, 0, 0, 0, time.UTC),
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
		Summary:     "Universal analytics platform especially made for the Minecraft and Hytale ecosystem. Track metrics, events and logs with ease.",
		Description: "Universal analytics platform especially made for the Minecraft and Hytale ecosystem. Track metrics, events and logs with ease.",
		Categories:  []spaces.Category{spaces.CategoryWebApp},
		Links: []spaces.Link{
			{Name: "website", URL: "https://fancyanalytics.net"},
			{Name: "source_code", URL: "https://github.com/FancyInnovations/FancyAnalytics"},
			{Name: "documentation", URL: "https://fancyinnovations.com/docs/web-services/fancyanalytics"},
			{Name: "discord", URL: "https://discord.gg/ZUgYCEJUEx"},
		},
		IconURL:   "https://fancyinnovations.com/logos-and-banners/fancyanalytics-logo.png",
		Status:    spaces.StatusApproved,
		CreatedAt: time.Date(2024, 1, 25, 20, 0, 0, 0, time.UTC),
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
		Summary:     "Tool to deploy Hytale and Minecraft plugins to multiple platforms via GitHub actions.",
		Description: "Tool to deploy Hytale and Minecraft plugins to multiple platforms via GitHub actions.",
		Categories:  []spaces.Category{spaces.CategoryOther},
		Links: []spaces.Link{
			{Name: "source_code", URL: "https://github.com/FancyInnovations/FancyVerteiler"},
			{Name: "discord", URL: "https://discord.gg/ZUgYCEJUEx"},
		},
		IconURL:   "https://fancyinnovations.com/logos-and-banners/fancyverteiler-logo.png",
		Status:    spaces.StatusApproved,
		CreatedAt: time.Date(2025, 12, 2, 20, 0, 0, 0, time.UTC),
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
