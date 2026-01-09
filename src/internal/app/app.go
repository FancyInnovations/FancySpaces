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
	"github.com/fancyinnovations/fancyspaces/internal/frontend"
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

	// Frontend
	fh := frontend.NewHandler(frontend.Configuration{
		Files: frontend.Files,
	})
	fh.Register(cfg.Mux)
}

func seedSpacesDB() *fakeSpacesDB.DB {
	db := fakeSpacesDB.New()

	fancynpcs := &spaces.Space{
		ID:          "fn",
		Slug:        "fancynpcs",
		Title:       "FancyNpcs",
		Summary:     "Simple, lightweight and fast NPC plugin using packets.",
		Description: "<div align=\"center\">\n\n![FancyNpcs Banner](https://github.com/FancyInnovations/FancyPlugins/blob/main/docs/src/static/logos-and-banners/fancynpcs-banner.png?raw=true)\n\n<br />\n\nSimple, lightweight and feature-rich NPC plugin for **[Paper](https://papermc.io/software/paper)** (and **[Folia](https://papermc.io/software/folia)**) servers using packets.\n\n</div>\n\n## Features\n\nWith this plugin you can create NPCs with customizable properties like:\n\n- **Type** (Cow, Pig, Player, etc.)\n- **Skin** (from username or texture URL)\n- **Glowing** (in all colors)\n- **Attributes** (pose, visibility, variant, etc.)\n- **Equipment** (eg. holding a diamond sword and wearing leather armor)\n- **Interactions** (player commands, console commands, messages)\n- ...and much more!\n\nCheck out **[images section](#images)** down below.\n\n<br />\n\n## Installation\n\nPaper **1.20** or newer with **Java 21** (or higher) is required. Plugin should also work on **Paper** forks.\n\n**Spigot** is **not** supported.\n\n<br />\n\n## Documentation\n\nOfficial documentation is hosted **[here](https://fancyinnovations.com/docs/minecraft-plugins/fancynpcs)**. Quick reference:\n\n- **[Getting started](https://fancyinnovations.com/docs/minecraft-plugins/fancynpcs/getting-started)**\n- **[Command Reference](https://fancyinnovations.com/docs/minecraft-plugins/fancynpcs/commands/npc)**\n- **[Using API](https://fancyinnovations.com/docs/minecraft-plugins/fancynpcs/api/getting-started)**\n\n**Have more questions?** Feel free to ask them on our **[Discord](https://discord.gg/ZUgYCEJUEx)** server.\n\n<br />\n\n## Developer API\n\nMore information can be found in **[Documentation](https://fancyinnovations.com/docs/minecraft-plugins/fancynpcs/api/getting-started)** and **[Javadocs](https://repo.fancyinnovations.com/javadoc/releases/de/oliver/FancyNpcs/latest)**.\n\n### Maven\n\n```xml\n<repository>\n    <id>fancyinnovations-releases</id>\n    <name>FancyInnovations Repository</name>\n    <url>https://repo.fancyinnovations.com/releases</url>\n</repository>\n```\n\n```xml\n<dependency>\n    <groupId>de.oliver</groupId>\n    <artifactId>FancyNpcs</artifactId>\n    <version>[VERSION]</version>\n    <scope>provided</scope>\n</dependency>\n```\n\n### Gradle\n\n```groovy\nrepositories {\n    maven(\"https://repo.fancyinnovations.com/releases\")\n}\n\ndependencies {\n    compileOnly(\"de.oliver:FancyNpcs:[VERSION]\")\n}\n```\n\n<br />\n\n## Images\n\nImages showcasing the plugin, sent to us by our community.\n\n![Screenshot 1](https://github.com/FancyMcPlugins/FancyNpcs/blob/main/images/screenshots/niceron1.jpeg?raw=true)  \n<sup>Provided by [Explorer's Eden](https://explorerseden.eu/)</sup>\n\n![Screenshot 2](https://github.com/FancyMcPlugins/FancyNpcs/blob/main/images/screenshots/niceron2.jpeg?raw=true)  \n<sup>Provided by [Explorer's Eden](https://explorerseden.eu/)</sup>\n\n![Screenshot 3](https://github.com/FancyMcPlugins/FancyNpcs/blob/main/images/screenshots/niceron3.jpeg?raw=true)  \n<sup>Provided by [Explorer's Eden](https://explorerseden.eu/)</sup>\n\n![Screenshot 4](https://github.com/FancyMcPlugins/FancyNpcs/blob/main/images/screenshots/dave1.jpeg?raw=true)  \n<sup>Provided by [Beacon's Quest](https://www.beaconsquest.net/)</sup>\n\n![Screenshot 5](https://github.com/FancyMcPlugins/FancyNpcs/blob/main/images/screenshots/oliver1.jpeg?raw=true)  \n<sup>Provided by [@OliverSchlueter](https://github.com/OliverSchlueter)</sup>\n\n![Screenshot 6](https://github.com/FancyMcPlugins/FancyNpcs/blob/main/images/screenshots/oliver2.jpeg?raw=true)  \n<sup>Provided by [@OliverSchlueter](https://github.com/OliverSchlueter)</sup>\n\n![Screenshot 7](https://github.com/FancyMcPlugins/FancyNpcs/blob/main/images/screenshots/grabsky1.jpeg?raw=true)  \n<sup>Provided by [@Grabsky](https://github.com/Grabsky)</sup>\n",
		Categories:  []spaces.Category{spaces.CategoryMinecraftPlugin},
		Links: []spaces.Link{
			{Name: "source_code", URL: "https://github.com/FancyInnovations/FancyPlugins"},
			{Name: "documentation", URL: "https://fancyinnovations.com/docs/minecraft-plugins/fancynpcs"},
			{Name: "discord", URL: "https://discord.gg/ZUgYCEJUEx"},
		},
		IconURL:   "https://fancyinnovations.com/logos-and-banners/fancynpcs-logo.png",
		Status:    spaces.StatusApproved,
		CreatedAt: time.Date(2022, 12, 19, 20, 0, 0, 0, time.UTC),
		Members: []spaces.Member{
			{
				UserID: "oliver",
				Role:   spaces.RoleOwner,
			},
		},
	}
	if err := db.Create(fancynpcs); err != nil {
		panic(fmt.Errorf("could not seed spaces db: %w", err))
	}

	fancyholograms := &spaces.Space{
		ID:          "fh",
		Slug:        "fancyholograms",
		Title:       "FancyHolograms",
		Summary:     "Simple, lightweight and fast NPC plugin using packets.",
		Description: "<div align=\"center\">\n\n![FancyHolograms Banner](https://github.com/FancyInnovations/FancyPlugins/blob/main/docs/src/static/logos-and-banners/fancyholograms-banner.png?raw=true)\n\n<br />\n\nSimple, lightweight and feature-rich hologram plugin for **[Paper](https://papermc.io/software/paper)** (and **[Folia](https://papermc.io/software/folia)**) servers using **[display entities](https://minecraft.wiki/w/Display)**\nand packets.\n\n</div>\n\n## Features\n\nWith this plugin you can create holograms with customizable properties like:\n\n- **Hologram Type** (text, item or block)\n- **Position**, **Rotation** and **Scale**\n- **Text Alignment**, **Background Color** and **Shadow**.\n- **Billboard** (fixed, center, horizontal, vertical)\n- **MiniMessage** formatting.\n- Placeholders support through **[PlaceholderAPI](https://github.com/PlaceholderAPI/PlaceholderAPI)** and **[MiniPlaceholders](https://github.com/MiniPlaceholders/MiniPlaceholders)** integration.\n- **[FancyNpcs](https://modrinth.com/plugin/fancynpcs)** integration.\n- ...and much more!\n\nCheck out **[images section](#images)** down below.\n\n<br />\n\n## Installation\n\nPaper **1.20** or newer with **Java 21** (or higher) is required. Plugin should also work on **Paper** forks.\n\n**Spigot** is **not** supported.\n\n<br />\n\n## Documentation\n\nOfficial documentation is hosted **[here](https://fancyinnovations.com/docs/minecraft-plugins/fancyholograms)**. Quick reference:\n\n- **[Getting Started](https://fancyinnovations.com/docs/minecraft-plugins/fancyholograms/getting-started)**\n- **[Command Reference](https://fancyinnovations.com/docs/minecraft-plugins/fancyholograms/commands/hologram)**\n- **[Using API](https://fancyinnovations.com/docs/minecraft-plugins/fancyholograms/api/getting-started)**\n\n**Have more questions?** Feel free to ask them on our **[Discord](https://discord.gg/ZUgYCEJUEx)** server.\n\n<br />\n\n## Developer API\n\nMore information can be found in **[Documentation](https://fancyinnovations.com/docs/minecraft-plugins/fancyholograms/api/getting-started)** and **[Javadocs](https://repo.fancyinnovations.com/javadoc/releases/de/oliver/FancyHolograms/latest)**.\n\n### Maven\n\n```xml\n<repository>\n    <id>fancyinnovations-releases</id>\n    <name>FancyInnovations Repository</name>\n    <url>https://repo.fancyinnovations.com/releases</url>\n</repository>\n```\n\n```xml\n<dependency>\n    <groupId>de.oliver</groupId>\n    <artifactId>FancyHolograms</artifactId>\n    <version>[VERSION]</version>\n    <scope>provided</scope>\n</dependency>\n```\n\n### Gradle\n\n```groovy\nrepositories {\n    maven(\"https://repo.fancyinnovations.com/releases\")\n}\n\ndependencies {\n    compileOnly(\"de.oliver:FancyHolograms:[VERSION]\")\n}\n```\n\n<br />\n\n## Images\n\nImages showcasing the plugin, sent to us by our community.\n\n![Screenshot 1](https://github.com/FancyMcPlugins/FancyHolograms/blob/main/images/screenshots/example1.jpeg?raw=true)  \n<sup>Provided by [@OliverSchlueter](https://github.com/OliverSchlueter)</sup>\n\n![Screenshot 2](https://github.com/FancyMcPlugins/FancyHolograms/blob/main/images/screenshots/example2.jpeg?raw=true)  \n<sup>Provided by [@OliverSchlueter](https://github.com/OliverSchlueter)</sup>\n\n![Screenshot 3](https://github.com/FancyMcPlugins/FancyHolograms/blob/main/images/screenshots/example3.jpeg?raw=true)  \n<sup>Provided by [@OliverSchlueter](https://github.com/OliverSchlueter)</sup>\n\n![Screenshot 4](https://github.com/FancyMcPlugins/FancyHolograms/blob/main/images/screenshots/example4.jpeg?raw=true)  \n<sup>Provided by [@OliverSchlueter](https://github.com/OliverSchlueter)</sup>\n\n![Screenshot 5](https://github.com/FancyMcPlugins/FancyHolograms/blob/main/images/screenshots/example5.jpeg?raw=true)  \n<sup>Provided by [@OliverSchlueter](https://github.com/OliverSchlueter)</sup>\n",
		Categories:  []spaces.Category{spaces.CategoryMinecraftPlugin},
		Links: []spaces.Link{
			{Name: "source_code", URL: "https://github.com/FancyInnovations/FancyPlugins"},
			{Name: "documentation", URL: "https://fancyinnovations.com/docs/minecraft-plugins/fancyholograms"},
			{Name: "discord", URL: "https://discord.gg/ZUgYCEJUEx"},
		},
		IconURL:   "https://fancyinnovations.com/logos-and-banners/fancyholograms-logo.png",
		Status:    spaces.StatusApproved,
		CreatedAt: time.Date(2023, 2, 18, 20, 0, 0, 0, time.UTC),
		Members: []spaces.Member{
			{
				UserID: "oliver",
				Role:   spaces.RoleOwner,
			},
		},
	}
	if err := db.Create(fancyholograms); err != nil {
		panic(fmt.Errorf("could not seed spaces db: %w", err))
	}

	fancydialogs := &spaces.Space{
		ID:          "fd",
		Slug:        "fancydialogs",
		Title:       "FancyDialogs",
		Summary:     "Simple, easy to use and lightweight plugin to show dialogs. You can show a welcome-screen. Other plugins can use FancyDialogs to integrate fancy dialogs into their plugins.",
		Description: "![FancyNpcs Banner](https://github.com/FancyInnovations/FancyPlugins/blob/main/docs/src/static/logos-and-banners/fancydialogs-banner.png?raw=true)\n\n<br />\n\nSimple and lightweight plugin to create and manage the new dialogs.\nBecause the dialogs were added in 1.21.6, only players on 1.21.6 or newer can view them.\n\nYou can create create dialogs (in JSON files) and then assign them to many cool features.\n\n## Features\n- Welcome dialog: shows when a player joins for the first time\n- Tutorials: explain how your amazing server works in multiple chapters [WIP]\n- MiniMessage formatting is supported\n- You can use placeholders by PlaceholderAPI and MiniPlaceholders\n- Awesome API for other plugin developers to use\n\n**Core advantages of FancyDialogs:**\n- Simple custom dialog creation (in JSON format or in code)\n- A lot of different dialog components (text, buttons, input fields, etc.)\n- Custom button actions (e.g. run commands, give items, open other dialogs)\n- MiniMessages and PlaceholderAPI support\n\n**For servers**\n\nIf you use FancyDialogs as a plugin, you can take advantage of the following features:\n- Dialog as welcome-screen for players joining the first time\n- Custom tutorial dialogs\n- FancyNpcs 'open_dialog' action for NPCs\n\n**For plugin developers**\n\nIf you are a plugin developer and want to spice up your plugin with dialogs, you can use FancyDialogs as a library.\n\nYou can define default dialogs, which will be persisted in the `plugins/FancyDialogs/data/dialogs` folder.\nThis allows server administrators to customize the dialogs for their server.\nYou can then use the dialogs in your plugin and show them to the players whenever you want.\n\n**Common use cases are:**\n- Help dialogs (e.g. for commands or features)\n- Confirmation dialogs for critical actions\n- Shop UIs (replacing inventories UIs)\n- Dialogs for quests\n",
		Categories:  []spaces.Category{spaces.CategoryMinecraftPlugin},
		Links: []spaces.Link{
			{Name: "source_code", URL: "https://github.com/FancyInnovations/FancyPlugins"},
			{Name: "documentation", URL: "https://fancyinnovations.com/docs/minecraft-plugins/fancydialogs"},
			{Name: "discord", URL: "https://discord.gg/ZUgYCEJUEx"},
		},
		IconURL:   "https://fancyinnovations.com/logos-and-banners/fancydialogs-logo.png",
		Status:    spaces.StatusApproved,
		CreatedAt: time.Date(2023, 2, 18, 20, 0, 0, 0, time.UTC),
		Members: []spaces.Member{
			{
				UserID: "oliver",
				Role:   spaces.RoleOwner,
			},
		},
	}
	if err := db.Create(fancydialogs); err != nil {
		panic(fmt.Errorf("could not seed spaces db: %w", err))
	}

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
				UserID: "oliver",
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
				UserID: "oliver",
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
				UserID: "oliver",
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
				UserID: "oliver",
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
				UserID: "oliver",
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
				UserID: "oliver",
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
				UserID: "oliver",
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
				UserID: "oliver",
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
				UserID: "oliver",
				Role:   spaces.RoleOwner,
			},
		},
	}
	if err := db.Create(fancyverteiler); err != nil {
		panic(fmt.Errorf("could not seed spaces db: %w", err))
	}

	clovepluralprojects := &spaces.Space{
		ID:          "cpp",
		Slug:        "cpp",
		Title:       "ClovePluralProjects",
		Summary:     "A plural accessibility tool",
		Description: "An accessibility tool for plural communities everywhere, over Minecraft and Hytale, with a web dash",
		Categories:  []spaces.Category{spaces.CategoryMinecraftPlugin, spaces.CategoryHytalePlugin, spaces.CategoryWebApp, spaces.CategoryMinecraftMod},
		Links: []spaces.Link{
			{Name: "website", URL: "https://clovelib.win"},
			{Name: "source_code", URL: "https://github.com/CloveLib/"},
			{Name: "documentation", URL: "https://clovelib.win/docs"},
			{Name: "discord", URL: "https://discord.gg/k8HrBvDaQn"},
		},
		IconURL:   "https://clovelib.win/icons/cpt.png",
		Status:    spaces.StatusApproved,
		CreatedAt: time.Date(2026, 1, 9, 17, 17, 0, 0, time.UTC),
		Members: []spaces.Member{
			{
				UserID: "clovelib",
				Role:   spaces.RoleOwner,
			},
		},
	}
	if err := db.Create(clovepluralprojects); err != nil {
		panic(fmt.Errorf("could not seed spaces db: %w", err))
	}

	return db
}
