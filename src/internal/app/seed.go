package app

import (
	"fmt"
	"time"

	"github.com/fancyinnovations/fancyspaces/internal/spaces"
	fakeSpacesDB "github.com/fancyinnovations/fancyspaces/internal/spaces/database/fake"
)

func seedSpacesDB() *fakeSpacesDB.DB {
	db := fakeSpacesDB.New()

	if err := seedMinecraftPlugins(db); err != nil {
		panic(fmt.Errorf("could not seed spaces db: %w", err))
	}

	if err := seedHytalePlugins(db); err != nil {
		panic(fmt.Errorf("could not seed spaces db: %w", err))
	}

	if err := seedOther(db); err != nil {
		panic(fmt.Errorf("could not seed spaces db: %w", err))
	}

	if err := seedByOtherCreators(db); err != nil {
		panic(fmt.Errorf("could not seed spaces db: %w", err))
	}

	return db
}

func seedMinecraftPlugins(db *fakeSpacesDB.DB) error {
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
		IssueSettings: spaces.IssueSettings{
			Enabled: true,
		},
		ReleaseSettings: spaces.ReleaseSettings{
			Enabled: true,
		},
		MavenRepositorySettings: spaces.MavenRepositorySettings{
			Enabled: true,
		},
	}
	if err := db.Create(fancynpcs); err != nil {
		return fmt.Errorf("could not seed spaces db: %w", err)
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
		IssueSettings: spaces.IssueSettings{
			Enabled: true,
		},
		ReleaseSettings: spaces.ReleaseSettings{
			Enabled: true,
		},
		MavenRepositorySettings: spaces.MavenRepositorySettings{
			Enabled: true,
		},
	}
	if err := db.Create(fancyholograms); err != nil {
		return fmt.Errorf("could not seed spaces db: %w", err)
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
		IssueSettings: spaces.IssueSettings{
			Enabled: true,
		},
		ReleaseSettings: spaces.ReleaseSettings{
			Enabled: true,
		},
		MavenRepositorySettings: spaces.MavenRepositorySettings{
			Enabled: true,
		},
	}
	if err := db.Create(fancydialogs); err != nil {
		return fmt.Errorf("could not seed spaces db: %w", err)
	}

	return nil
}

func seedHytalePlugins(db *fakeSpacesDB.DB) error {
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
		IssueSettings: spaces.IssueSettings{
			Enabled: true,
		},
		ReleaseSettings: spaces.ReleaseSettings{
			Enabled: true,
		},
		MavenRepositorySettings: spaces.MavenRepositorySettings{
			Enabled: true,
		},
	}
	if err := db.Create(fancycore); err != nil {
		return fmt.Errorf("could not seed spaces db: %w", err)
	}

	fancycorewebsite := &spaces.Space{
		ID:         "fancycorewebsite",
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
		IssueSettings: spaces.IssueSettings{
			Enabled: true,
		},
		ReleaseSettings: spaces.ReleaseSettings{
			Enabled: true,
		},
		MavenRepositorySettings: spaces.MavenRepositorySettings{
			Enabled: true,
		},
	}
	if err := db.Create(fancycorewebsite); err != nil {
		return fmt.Errorf("could not seed spaces db: %w", err)
	}

	fancyplots := &spaces.Space{
		ID:         "fancyplots",
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
		IssueSettings: spaces.IssueSettings{
			Enabled: true,
		},
		ReleaseSettings: spaces.ReleaseSettings{
			Enabled: true,
		},
		MavenRepositorySettings: spaces.MavenRepositorySettings{
			Enabled: true,
		},
	}
	if err := db.Create(fancyplots); err != nil {
		return fmt.Errorf("could not seed spaces db: %w", err)
	}

	fancyaudits := &spaces.Space{
		ID:         "fancyaudits",
		Slug:       "fancyaudits",
		Title:      "FancyAudits",
		Summary:    "Log various player actions + dupe detection for Hytale.",
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
		IssueSettings: spaces.IssueSettings{
			Enabled: true,
		},
		ReleaseSettings: spaces.ReleaseSettings{
			Enabled: true,
		},
		MavenRepositorySettings: spaces.MavenRepositorySettings{
			Enabled: true,
		},
	}
	if err := db.Create(fancyaudits); err != nil {
		return fmt.Errorf("could not seed spaces db: %w", err)
	}

	fancyconnect := &spaces.Space{
		ID:         "fancyconnect",
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
		IssueSettings: spaces.IssueSettings{
			Enabled: true,
		},
		ReleaseSettings: spaces.ReleaseSettings{
			Enabled: true,
		},
		MavenRepositorySettings: spaces.MavenRepositorySettings{
			Enabled: true,
		},
	}
	if err := db.Create(fancyconnect); err != nil {
		return fmt.Errorf("could not seed spaces db: %w", err)
	}

	fancyshops := &spaces.Space{
		ID:         "fancyshops",
		Slug:       "fancyshops",
		Title:      "FancyShops",
		Summary:    "Admin shops, chest shops, auctions, trade system and more for Hytale.",
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
		IssueSettings: spaces.IssueSettings{
			Enabled: true,
		},
		ReleaseSettings: spaces.ReleaseSettings{
			Enabled: true,
		},
		MavenRepositorySettings: spaces.MavenRepositorySettings{
			Enabled: true,
		},
	}
	if err := db.Create(fancyshops); err != nil {
		return fmt.Errorf("could not seed spaces db: %w", err)
	}

	citypass := &spaces.Space{
		ID:         "citypass",
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
		IssueSettings: spaces.IssueSettings{
			Enabled: true,
		},
		ReleaseSettings: spaces.ReleaseSettings{
			Enabled: true,
		},
		MavenRepositorySettings: spaces.MavenRepositorySettings{
			Enabled: true,
		},
	}
	if err := db.Create(citypass); err != nil {
		return fmt.Errorf("could not seed spaces db: %w", err)
	}

	cityquests := &spaces.Space{
		ID:         "cityquests",
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
		IssueSettings: spaces.IssueSettings{
			Enabled: true,
		},
		ReleaseSettings: spaces.ReleaseSettings{
			Enabled: true,
		},
		MavenRepositorySettings: spaces.MavenRepositorySettings{
			Enabled: true,
		},
	}
	if err := db.Create(cityquests); err != nil {
		return fmt.Errorf("could not seed spaces db: %w", err)
	}

	return nil
}

func seedOther(db *fakeSpacesDB.DB) error {
	fancyanalytics := &spaces.Space{
		ID:          "fa",
		Slug:        "fancyanalytics",
		Title:       "FancyAnalytics",
		Summary:     "Universal analytics platform especially made for the Minecraft and Hytale ecosystem. Track metrics, events and logs with ease.",
		Description: "![](https://fancyinnovations.com/logos-and-banners/fancyanalytics-banner.png)\n\nThe modern analytics platform for Minecraft server owners and plugin developers\n\n## Features\n\nWith FancyAnalytics, you will be able to understand the behavior of your users like never before. You will be able to:\n\n- Collect metrics like player count, server uptime, and more\n- Track events like player joins, deaths, and more\n- Visualize your data with beautiful charts and graphs on multiple dashboards\n- Analyze errors and exceptions that occur on your server or in your plugins\n\nFancyAnalytics is still in development, but we are working hard to bring you the best analytics platform for Minecraft. We plan to add many more features in the future, including:\n- Pre-defined metrics and dashboards for popular plugins\n- Community-driven metrics and dashboards\n- Alerts and notifications for important events\n- Real-time data streaming\n- Log viewing and analysis\n- Organization and team management\n- Survey and feedback collection\n- Support and issue tracking\n\n## Getting Started\n\nGetting started with FancyAnalytics is easy!\nJust create an account on our website ([https://fancyanalytics.net](https://fancyanalytics.net)) and follow the instructions to set up your server or plugin.\nYou can also check out our documentation for more information on how to use FancyAnalytics.",
		Categories:  []spaces.Category{spaces.CategoryWebApp, spaces.CategoryMinecraftPlugin, spaces.CategoryHytalePlugin},
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
		IssueSettings: spaces.IssueSettings{
			Enabled: true,
		},
		ReleaseSettings: spaces.ReleaseSettings{
			Enabled: true,
		},
		MavenRepositorySettings: spaces.MavenRepositorySettings{
			Enabled: true,
		},
	}
	if err := db.Create(fancyanalytics); err != nil {
		return fmt.Errorf("could not seed spaces db: %w", err)
	}

	fancyverteiler := &spaces.Space{
		ID:          "fv",
		Slug:        "fancyverteiler",
		Title:       "FancyVerteiler",
		Summary:     "Tool to deploy Hytale and Minecraft plugins to multiple platforms via GitHub actions or standalone app.",
		Description: "# FancyVerteiler\n\nWith FancyVerteiler, you can deploy Minecraft and Hytale plugins to multiple platforms at once via GitHub Actions or the standalone app.\n\n## Features\n\n- Configure multiple platforms in a single JSON configuration file.\n- Automatically read version and changelog from files.\n- Send notifications to a Discord channel via webhook.\n\nSupported Minecraft plugin platforms:\n- [FancySpaces](https://fancyspaces.net/)\n- [Modrinth](https://modrinth.com/)\n- [CurseForge](https://www.curseforge.com/)\n\nSupported Hytale plugin platforms:\n- [FancySpaces](https://fancyspaces.net/)\n- [Orbis](https://orbis.place/)\n- [Modtale](https://modtale.net/)\n- [UnifiedHytale](https://www.unifiedhytale.com)\n- [HytaHub](https://hytahub.com/)\n\n## Usage\n\n### GitHub Actions\n\nInclude the following in your GitHub Actions workflow:\n```yml\n- uses: fancyinnovations/fancyverteiler@main\n  with:\n    config_path: \"plugins/fancynpcs/release_deployment_config.json\"\n    commit_sha: \"see example for git integration below\"\n    commit_message: \"see example for git integration below\"\n    fancyspaces_api_key: ${{ secrets.FANCYSPACES_API_KEY }}\n    modrinth_api_key: ${{ secrets.MODRINTH_API_KEY }}\n    curseforge_api_key: ${{ secrets.CURSEFORGE_API_KEY }}\n    orbis_api_key: ${{ secrets.ORBIS_API_KEY }}\n    modtale_api_key: ${{ secrets.MODTALE_API_KEY }}\n    unifiedhytale_api_key: ${{ secrets.UNIFIEDHYTALE_API_KEY }}\n    hytahub_api_key: ${{ secrets.HYTAHUB_API_KEY }}\n    discord_webhook_url: ${{ secrets.DISCORD_WEBHOOK_URL }}\n```\n\nInputs:\n- `config_path` (required): Path to the JSON configuration file for FancyVerteiler.\n- `commit_sha` (optional): The commit SHA to replace in the changelog.\n- `commit_message` (optional): The commit message to replace in the changelog.\n- `<platform>_api_key` is only required if you want to publish to <platform>.\n\nExample json config:\n```json\n{\n  \"project_name\": \"FancyNpcs\",\n  \"plugin_jar_path\": \"./plugins/fancynpcs/build/libs/FancyNpcs-%VERSION%.jar\",\n  \"changelog_path\": \"./plugins/fancynpcs/CHANGELOG.md\",\n  \"version_path\": \"./plugins/fancynpcs/VERSION\",\n  \"fancyspaces\": {\n    \"space_id\": \"fn\",\n    \"platform\": \"paper\",\n    \"channel\": \"release\",\n    \"supported_versions\": [ \"1.21.10\", \"1.21.11\" ]\n  },\n  \"modrinth\": {\n    \"project_id\": \"EeyAn23L\",\n    \"supported_versions\": [ \"1.21.10\", \"1.21.11\" ],\n    \"channel\": \"release\",\n    \"loaders\": [ \"paper\", \"folia\" ],\n    \"featured\": true\n  },\n  \"curseforge\": {\n    \"project_id\": \"123456\",\n    \"type\": \"plugin\",\n    \"game_versions\": [ \"1.21.10\", \"1.21.11\" ],\n    \"release_type\": \"release\"\n  },\n  \"orbis\": {\n    \"resource_id\": \"1234\",\n    \"channel\": \"RELEASE\",\n    \"hytale_version_ids\": [ \"cmj1x42ef001k4qz9r03ojrpe\" ]\n  },\n  \"modtale\": {\n    \"project_id\": \"abcdef123456\",\n    \"channel\": \"RELEASE\",\n    \"game_versions\": [ \"2026.12.02.12312313\" ]\n  },\n  \"unifiedhytale\": {\n    \"project_id\": \"abcdef123456\",\n    \"game_versions\": [ \"2026.12.02.12312313\" ],\n    \"release_channel\": \"release\"\n  },\n  \"hytahub\": {\n    \"slug\": \"mymode\",\n    \"channel\": \"release\"\n  }\n}\n```\n\nTo automatically get the last commit SHA and message from git, you can add the following steps before the FancyVerteiler step:\n```yml\n      - name: Get last commit SHA and message\n        id: last_commit\n        run: |\n          {\n            echo \"commit_sha=$(git rev-parse --short HEAD)\"\n            echo \"commit_msg<<EOF\"\n            git log -1 --pretty=%B\n            echo \"EOF\"\n          } >> \"$GITHUB_OUTPUT\"\n\n      - name: Deploy\n        uses: fancyinnovations/fancyverteiler@main\n        with:\n          config_path: \"/plugins/fancynpcs/release_deployment_config.json\"\n          commit_sha: ${{ steps.last_commit.outputs.commit_sha }}\n          commit_message: ${{ steps.last_commit.outputs.commit_msg }}\n          modrinth_api_key: ${{ secrets.MODRINTH_API_KEY }}\n          discord_webhook_url: ${{ secrets.DISCORD_WEBHOOK_URL }}\n\n```\n\nThis will replace `%COMMIT_HASH%` and `%COMMIT_MESSAGE%` in the changelog with the actual commit hash and message.\n\n### Standalone\n\nYou can also run FancyVerteiler as a standalone app.\nEverything works the same way as in GitHub Actions, but you need to provide the inputs as environment variables.\n\nEnvironment variables:\n- `FV_CONFIG_PATH`\n- `FV_DISCORD_WEBHOOK_URL`\n- `FV_COMMIT_SHA`\n- `FV_MESSAGE_SHA`\n- `FV_{PLATFORM}_API_KEY` (example: `FV_FANCYSPACES_API_KEY`)\n\nYou can download the latest version of the standalone app from [FancySpaces](http://localhost:3001/spaces/fancyverteiler).",
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
		IssueSettings: spaces.IssueSettings{
			Enabled: true,
		},
		ReleaseSettings: spaces.ReleaseSettings{
			Enabled: true,
		},
		MavenRepositorySettings: spaces.MavenRepositorySettings{
			Enabled: true,
		},
	}
	if err := db.Create(fancyverteiler); err != nil {
		return fmt.Errorf("could not seed spaces db: %w", err)
	}

	return nil
}

func seedByOtherCreators(db *fakeSpacesDB.DB) error {
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
		IssueSettings: spaces.IssueSettings{
			Enabled: true,
		},
		ReleaseSettings: spaces.ReleaseSettings{
			Enabled: true,
		},
		MavenRepositorySettings: spaces.MavenRepositorySettings{
			Enabled: false,
		},
	}
	if err := db.Create(clovepluralprojects); err != nil {
		return fmt.Errorf("could not seed spaces db: %w", err)
	}

	orbisguard := &spaces.Space{
		ID:          "orbisguard",
		Slug:        "orbisguard",
		Title:       "OrbisGuard",
		Summary:     "Region protection plugin for Hytale servers. Define areas, set flags (block-break, block-place, use), and control who can build where. Based on WorldGuard for Minecraft.",
		Description: "![OrbisGuard](https://docs.wiflow.dev/orbisguard-banner.png)\n\n[![Discord](https://img.shields.io/badge/Discord-Join%20Server-5865F2?style=for-the-badge&logo=discord&logoColor=white)](https://discord.gg/Dfy8eDGFnK) [![Documentation](https://img.shields.io/badge/Docs-Full%20Documentation-2ea44f?style=for-the-badge&logo=gitbook&logoColor=white)](https://docs.wiflow.dev/hytale-plugins)\n\n**Region protection for Hytale.** Define areas, set rules, control who can do what.\n\nThink **WorldGuard** but for Hytale.\n\n---\n\n## Features\n\n- **Mob Spawning Control** - Block natural mob and NPC spawns in regions\n- **Block Protection** - Prevent breaking, placing, and hammer cycling\n- **Container Locks** - Protect chests, barrels, and storage\n- **PvP Zones** - Enable or disable player combat per region\n- **Entry Control** - Restrict who can enter or leave areas\n- **Custom Messages** - Welcome titles and farewell messages\n- **Minimap Integration** - Regions render on the world map with labels\n- **Polygonal Regions** - Create non-rectangular areas\n- **LuckPerms Groups** - Add entire groups as members (`g:vip`)\n- **Public API** - Other plugins can query and modify regions\n\n---\n\n## Quick Start\n\n```\n/rg wand                     - Enable selection mode\nLeft-click: set pos1         - Break a block to mark corner 1\nRight-click: set pos2        - Place a block to mark corner 2\n/rg define myregion          - Create the region\n/rg flag myregion build deny - Block building\n/rg addmember myregion Steve - Add a member who can build\n```\n\n---\n\n## Extended Protection (Hyxin)\n\nSome Hytale mechanics aren't exposed through the plugin API. Without [Hyxin](https://www.curseforge.com/hytale/mods/hyxin), players can bypass protection using:\n\n| Exploit | What happens |\n|---------|--------------|\n| Auto-pickup | Items fly into inventory automatically |\n| F-key harvesting | Pressing F to harvest crops/ores |\n| Hammer cycling | Changing block variants with hammer |\n| Fluid placement | Placing water/lava with buckets |\n| Explosions | TNT/bombs destroying protected blocks |\n| Command bypass | Using /home to escape regions |\n\n**OrbisGuard is the only protection plugin that blocks ALL of these.**\n\n**Setup with Hyxin:**\n\n```\nHytale/\n├── earlyplugins/\n│   ├── Hyxin-x.x.x.jar             <- Mixin loader\n│   └── OrbisGuard-Mixins-0.3.0.jar <- Extended protection\n└── mods/\n    └── OrbisGuard-0.3.0.jar        <- Main plugin\n```\n\n---\n\n## All Commands\n\n| Command | Description |\n|---------|-------------|\n| `/rg wand` | Toggle selection wand mode |\n| `/rg selmode cuboid\\|poly` | Switch between cuboid and polygon selection |\n| `/rg define <name>` | Create region from selection |\n| `/rg definepoly <name>` | Create polygonal region |\n| `/rg remove <name>` | Delete a region |\n| `/rg flag <region> <flag> <value>` | Set a flag |\n| `/rg addmember <region> <player>` | Add member |\n| `/rg addowner <region> <player>` | Add owner |\n| `/rg info [region]` | Show region details |\n| `/rg list` | List all regions |\n| `/rg bypass` | Toggle admin bypass |\n\n[See all 20+ commands in the full docs →](https://docs.wiflow.dev/hytale-plugins/orbisguard/commands)\n\n---\n\n## Flags\n\n### Protection\n\n| Flag | Default | Description |\n|------|---------|-------------|\n| `build` | deny | Master toggle for all block changes |\n| `block-break` | deny | Breaking blocks |\n| `block-place` | deny | Placing blocks |\n| `hammer` | deny | Hammer cycling (requires Hyxin) |\n| `explosions` | deny | TNT/bomb damage (requires Hyxin) |\n| `chest-access` | deny | Opening containers |\n\n### Combat & Movement\n\n| Flag | Default | Description |\n|------|---------|-------------|\n| `pvp` | allow | Player vs player combat |\n| `damage` | allow | All damage (master toggle) |\n| `fall-damage` | allow | Fall damage only |\n| `entry` | allow | Entering the region |\n| `exit` | allow | Leaving the region |\n\n### Items & Interaction\n\n| Flag | Default | Description |\n|------|---------|-------------|\n| `use` | allow | Doors, buttons, levers |\n| `item-pickup` | allow | F-key item pickup |\n| `item-pickup-auto` | allow | Auto-pickup (requires Hyxin) |\n| `item-drop` | allow | Dropping items |\n| `crafting` | allow | Using crafting stations |\n\n### Messages & Minimap\n\n| Flag | Description |\n|------|-------------|\n| `greet-title` | Title shown when entering |\n| `farewell-title` | Title shown when leaving |\n| `deny-message` | Custom message when action is blocked |\n| `blocked-cmds` | Commands to block (requires Hyxin) |\n| `allowed-cmds` | Only these commands work (requires Hyxin) |\n| `minimap-color` | Region color on world map (hex) |\n| `minimap-label` | Text label on minimap |\n\n---\n\n## Plugin Integrations\n\n- **LuckPerms** - Add permission groups as members: `/rg addmember spawn g:vip`\n- **SimpleClaims** - Party claims render on world map alongside regions\n- **EasyClaims** - Player claims render on world map\n\n> **Patched versions** of SimpleClaims and EasyClaims are available on our [Discord](https://discord.gg/n6tTAVUSgN). These prevent players from claiming chunks that overlap with protected regions.\n\n---\n\n## For Developers\n\nOrbisGuard has a public API for querying and modifying regions:\n\n```java\nOrbisGuardAPI api = OrbisGuardAPI.getInstance();\nIRegionManager manager = api.getRegionContainer().getRegionManager(\"world\");\n\n// Query regions at a location\nSet<IRegion> regions = manager.getRegionsAt(100, 64, 200);\n\n// Create a region\nmanager.createRegion(\"shop\", BlockVector3.at(0, 0, 0), BlockVector3.at(100, 256, 100));\n\n// Listen for events\napi.getEventBus().register(RegionCreatedEvent.class, event -> {\n    System.out.println(\"Region created: \" + event.getRegion().getId());\n});\n```\n\n[Full API documentation →](https://docs.wiflow.dev/hytale-plugins)\n\n---\n\n## Support\n\nQuestions or issues? Join the [Discord server](https://discord.gg/n6tTAVUSgN) or contact **w1fl0w**.\n",
		Categories:  []spaces.Category{spaces.CategoryHytalePlugin},
		Links: []spaces.Link{
			{Name: "documentation", URL: "https://docs.wiflow.dev/hytale-plugins"},
			{Name: "discord", URL: "https://discord.gg/Dfy8eDGFnK"},
		},
		IconURL:   "https://docs.wiflow.dev/orbisguard-icon.png",
		Status:    spaces.StatusApproved,
		CreatedAt: time.Date(2026, 1, 17, 20, 0, 0, 0, time.UTC),
		Members: []spaces.Member{
			{
				UserID: "wiflow",
				Role:   spaces.RoleOwner,
			},
		},
		IssueSettings: spaces.IssueSettings{
			Enabled: true,
		},
		ReleaseSettings: spaces.ReleaseSettings{
			Enabled: true,
		},
		MavenRepositorySettings: spaces.MavenRepositorySettings{
			Enabled: false,
		},
	}
	if err := db.Create(orbisguard); err != nil {
		return fmt.Errorf("could not seed spaces db: %w", err)
	}

	orbismines := &spaces.Space{
		ID:          "orbismines",
		Slug:        "orbismines",
		Title:       "OrbisMines",
		Summary:     "A Hytale server plugin for creating and managing auto-resetting mines with customizable block compositions, timer-based or percentage-based resets, and full GUI editing support.",
		Description: "![OrbisMines](https://docs.wiflow.dev/orbismines-banner.png)\n\n[![Discord](https://img.shields.io/badge/Discord-Join%20Server-5865F2?style=for-the-badge&logo=discord&logoColor=white)](https://discord.gg/Dfy8eDGFnK) [![Documentation](https://img.shields.io/badge/Docs-Full%20Documentation-2ea44f?style=for-the-badge&logo=gitbook&logoColor=white)](https://docs.wiflow.dev/hytale-plugins)\n\n**Auto-resetting mines for Hytale.** Create mines that reset on a timer or when depleted.\n\nThink **MineResetLite** but for Hytale.\n\n---\n\n## Requirements\n\n**[OrbisGuard](https://www.curseforge.com/hytale/mods/orbisguard)** is required. Mines are linked to OrbisGuard regions for boundaries and protection.\n\n```\nHytale/\n└── mods/\n    ├── OrbisGuard-x.x.x.jar  <- Required dependency\n    └── OrbisMines-1.0.0.jar  <- This plugin\n```\n\n---\n\n## Features\n\n- **Timer Resets** - Automatically reset mines every X minutes\n- **Percentage Resets** - Reset when X% of blocks are mined\n- **Block Composition** - Define what blocks spawn and their percentages\n- **Surface Layer** - Optional different block for the top layer\n- **Fill Mode** - Only replace air/mined blocks (preserve structures)\n- **Teleport Safety** - Move players out before resetting\n- **Reset Warnings** - Configurable minute and second warnings\n- **Local Broadcast** - Only notify players inside the mine\n- **GUI Editor** - Visual editor for mine composition and settings\n- **Async Resets** - Non-blocking block placement\n\n---\n\n## Quick Start\n\n```\n1. Create an OrbisGuard region first:\n   /rg wand\n   (select two corners)\n   /rg define mymine\n\n2. Create a mine from the region:\n   /mine create mymine              - Uses region with same name\n   /mine create mymine otherregion  - Uses a different region\n\n3. Add blocks to the composition:\n   /mine set mymine Stone 50\n   /mine set mymine Coal_Ore 30\n   /mine set mymine Iron_Ore 20\n\n4. Set a reset timer (minutes):\n   /mine flag mymine delay 15\n\n5. Reset the mine:\n   /mine reset mymine\n```\n\n---\n\n## GUI Editor\n\nUse `/mine edit` to open a visual editor with:\n\n- **Blocks Tab** - Add/remove blocks, adjust percentages with +/- buttons\n- **Settings Tab** - Configure reset timer, percentage trigger, and options\n- **Region Editing** - Adjust mine coordinates with \"Set Here\" buttons\n- **Block Picker** - Searchable dropdown of all block types\n\n---\n\n## All Commands\n\n| Command | Description |\n|---------|-------------|\n| `/mine create <name> [region]` | Create mine (region defaults to name if not specified) |\n| `/mine delete <name>` | Delete a mine |\n| `/mine list` | List all mines |\n| `/mine info <name>` | Show mine details |\n| `/mine reset <name>` | Manually reset a mine |\n| `/mine set <name> <block> <percent>` | Add/update block in composition |\n| `/mine unset <name> <block>` | Remove block from composition |\n| `/mine flag <name> <flag> <value>` | Set mine flag |\n| `/mine tp <name>` | Teleport to mine |\n| `/mine settp <name>` | Set mine teleport point |\n| `/mine edit [name]` | Open GUI editor (or mine list) |\n\nAliases: `/mines`, `/om`\n\n---\n\n## Flags\n\n| Flag | Default | Description |\n|------|---------|-------------|\n| `delay` | 0 | Reset timer in minutes (0 = disabled) |\n| `percent` | -1 | Reset when X% mined (-1 = disabled) |\n| `fillmode` | false | Only replace air/mined blocks |\n| `silent` | false | Disable all reset messages |\n| `localbroadcast` | false | Only message players in the mine |\n| `teleport` | true | Move players out during reset |\n| `surface` | none | Block type for top layer (e.g., Grass) |\n\n---\n\n## Configuration\n\nEdit `mods/OrbisMines/config.json` to customize:\n\n- Reset broadcast messages\n- Warning messages (minutes and seconds)\n- Default values for new mines\n- Blocks per tick (performance tuning)\n\n---\n\n## Example Setup\n\n**Prison-style mine with mixed ores:**\n\n```\n/rg define mine_a\n/mine create mine_a\n/mine set mine_a Stone 60\n/mine set mine_a Coal_Ore 20\n/mine set mine_a Iron_Ore 15\n/mine set mine_a Gold_Ore 5\n/mine flag mine_a delay 10\n/mine flag mine_a teleport true\n/mine settp mine_a\n```\n\n**Quick-reset PvP mine:**\n\n```\n/mine create pvpmine\n/mine set pvpmine Rock_Crystal_Cyan_Block 100\n/mine flag pvpmine percent 80\n/mine flag pvpmine silent true\n```\n\n---\n\n## Permissions\n\n| Permission | Description |\n|------------|-------------|\n| `orbismines.admin` | Access to all mine commands |\n\n---\n\n## Support\n\nQuestions or issues? Join the [Discord server](https://discord.gg/n6tTAVUSgN) or contact **w1fl0w**.\n",
		Categories:  []spaces.Category{spaces.CategoryHytalePlugin},
		Links: []spaces.Link{
			{Name: "documentation", URL: "https://docs.wiflow.dev/hytale-plugins"},
			{Name: "discord", URL: "https://discord.gg/Dfy8eDGFnK"},
		},
		IconURL:   "https://docs.wiflow.dev/orbismines-icon.png",
		Status:    spaces.StatusApproved,
		CreatedAt: time.Date(2026, 1, 19, 20, 0, 0, 0, time.UTC),
		Members: []spaces.Member{
			{
				UserID: "wiflow",
				Role:   spaces.RoleOwner,
			},
		},
		IssueSettings: spaces.IssueSettings{
			Enabled: true,
		},
		ReleaseSettings: spaces.ReleaseSettings{
			Enabled: true,
		},
		MavenRepositorySettings: spaces.MavenRepositorySettings{
			Enabled: false,
		},
	}
	if err := db.Create(orbismines); err != nil {
		return fmt.Errorf("could not seed spaces db: %w", err)
	}

	wiflowsScoreboard := &spaces.Space{
		ID:          "wiflowscoreboard",
		Slug:        "wiflowscoreboard",
		Title:       "WiFlow's Scoreboard",
		Summary:     "A fully customizable scoreboard HUD with world-specific profiles, animated text, playtime tracking, and LuckPerms integration.",
		Description: "![WiFlowScoreboard](https://docs.wiflow.dev/wiflows-scoreboard-banner.png)\n\n[![Discord](https://img.shields.io/badge/Discord-Join%20Server-5865F2?style=for-the-badge&logo=discord&logoColor=white)](https://discord.gg/Dfy8eDGFnK) [![Documentation](https://img.shields.io/badge/Docs-Full%20Documentation-2ea44f?style=for-the-badge&logo=gitbook&logoColor=white)](https://docs.wiflow.dev/hytale-plugins)\n\n**Customizable scoreboard HUD for Hytale.** Different scoreboards per world, animated text, playtime tracking, and full placeholder support.\n\n---\n\n## Features\n\n- **World-Specific Profiles** - Different scoreboards for each world\n- **Permission Variants** - Show different content based on permissions\n- **Animated Text** - Rainbow, color cycling, and custom animations\n- **Playtime Tracking** - Session and total playtime (persisted across restarts)\n- **LuckPerms Integration** - Display prefix, suffix, and group\n- **Custom Placeholders** - API for other plugins to register placeholders\n- **MultipleHUD Compatible** - Works alongside other HUD plugins\n- **Hot Reload** - Change configs without restarting\n- **Custom Logo Support** - Add your own server logo (easy installer tool on Discord)\n\n---\n\n## Placeholders\n\n| Placeholder | Description |\n|-------------|-------------|\n| `{player}` | Player display name |\n| `{world}` | Current world name |\n| `{online}` | Online player count |\n| `{max_players}` | Server max players |\n| `{server}` | Server name |\n| `{tps}` | Server TPS |\n| `{playtime}` | Session playtime |\n| `{playtime_current}` | Session playtime |\n| `{playtime_total}` | Total playtime (all sessions) |\n| `{prefix}` | LuckPerms prefix |\n| `{suffix}` | LuckPerms suffix |\n| `{group}` | LuckPerms primary group |\n\n---\n\n## Animations\n\nAdd animations to any line using tags:\n\n```\n[anim:rainbow]Welcome to the server!\n[anim:colors:#ff0000:#00ff00:#0000ff]Cycling colors\n[anim:fire]Hot text\n```\n\n### Built-in Animations\n\n| Animation | Description |\n|-----------|-------------|\n| `rainbow` | Full spectrum color cycle |\n| `colors` | Custom color sequence |\n| `flash` | White/red alternating |\n| `fire` | Red/orange/yellow gradient |\n| `ocean` | Blue wave |\n| `forest` | Green shades |\n| `ice` | Cold blues and white |\n| `neon` | Vibrant cycling |\n| `gold` | Shimmering gold |\n\nCreate custom animations by adding JSON files to `mods/WiFlowScoreboard/animations/`\n\n---\n\n## Quick Start\n\n```\nmods/WiFlowScoreboard/\n├── config.json              # Global settings\n├── playtime.json            # Persistent playtime data\n├── scoreboards/\n│   ├── default.json         # Default profile (all worlds)\n│   └── world_survival.json  # World-specific profile\n└── animations/\n    └── custom.json          # Custom animations\n```\n\n### Example Profile (default.json)\n\n```json\n{\n  \"shared\": {\n    \"width\": 280,\n    \"height\": 450,\n    \"positionX\": 10,\n    \"positionY\": 300,\n    \"positionSide\": \"right\",\n    \"showLogo\": true,\n    \"showTitle\": true,\n    \"panelColor\": \"#00000055\"\n  },\n  \"profiles\": [\n    {\n      \"name\": \"default\",\n      \"default\": true,\n      \"title\": { \"text\": \"{server}\" },\n      \"lines\": [\n        \"\",\n        \"[#ffaa00]Player: [#ffffff]{player}\",\n        \"[#ffaa00]World: [#ffffff]{world}\",\n        \"\",\n        \"[#88ffff]Playtime: {playtime_total}\",\n        \"[#88ffff]Online: {online}/{max_players}\",\n        \"\",\n        \"[anim:rainbow]Have fun!\"\n      ]\n    }\n  ]\n}\n```\n\n---\n\n## World-Specific Profiles\n\nCreate `world_<name>.json` files for different worlds:\n\n```\nscoreboards/\n├── default.json           # Fallback for all worlds\n├── world_survival.json    # Survival world\n├── world_pvp.json         # PvP arena\n└── world_hub.json         # Hub/lobby\n```\n\n---\n\n## Permission-Based Variants\n\nShow different content based on player permissions:\n\n```json\n{\n  \"profiles\": [\n    {\n      \"name\": \"vip\",\n      \"permission\": \"scoreboard.vip\",\n      \"title\": { \"text\": \"[VIP] {player}\" },\n      \"lines\": [\"VIP exclusive content...\"]\n    },\n    {\n      \"name\": \"default\",\n      \"default\": true,\n      \"title\": { \"text\": \"{player}\" },\n      \"lines\": [\"Regular content...\"]\n    }\n  ]\n}\n```\n\n---\n\n## Commands\n\n| Command | Description |\n|---------|-------------|\n| `/sb reload` | Reload all configs |\n| `/sb hide` | Hide your scoreboard |\n| `/sb show` | Show your scoreboard |\n\nAliases: `/scoreboard`, `/wsb`\n\n---\n\n## For Developers\n\nRegister custom placeholders from your plugin:\n\n```java\nHyScoreboardAPI.registerPlaceholder(\"kills\", context -> {\n    int kills = getPlayerKills(context.getPlayerUuid());\n    return String.valueOf(kills);\n});\n\nHyScoreboardAPI.registerPlaceholder(\"balance\", context -> {\n    double bal = getBalance(context.getPlayerUuid());\n    return String.format(\"$%.2f\", bal);\n});\n```\n\n---\n\n## Custom Logo\n\nWant your server logo on the scoreboard? Download the **Logo Installer Tool** from our [Discord server](https://discord.gg/n6tTAVUSgN) - just drag and drop your logo.png and run the script!\n\n---\n\n## Support\n\nQuestions or issues? Join the [Discord server](https://discord.gg/n6tTAVUSgN) or contact **w1fl0w**.\n",
		Categories:  []spaces.Category{spaces.CategoryHytalePlugin},
		Links: []spaces.Link{
			{Name: "documentation", URL: "https://docs.wiflow.dev/hytale-plugins"},
			{Name: "discord", URL: "https://discord.gg/Dfy8eDGFnK"},
		},
		IconURL:   "https://docs.wiflow.dev/wiflows-scoreboard-icon.png",
		Status:    spaces.StatusApproved,
		CreatedAt: time.Date(2026, 1, 20, 20, 0, 0, 0, time.UTC),
		Members: []spaces.Member{
			{
				UserID: "wiflow",
				Role:   spaces.RoleOwner,
			},
		},
		IssueSettings: spaces.IssueSettings{
			Enabled: true,
		},
		ReleaseSettings: spaces.ReleaseSettings{
			Enabled: true,
		},
		MavenRepositorySettings: spaces.MavenRepositorySettings{
			Enabled: false,
		},
	}
	if err := db.Create(wiflowsScoreboard); err != nil {
		return fmt.Errorf("could not seed spaces db: %w", err)
	}

	return nil
}
