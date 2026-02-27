package spaces_test

import (
	"errors"
	"testing"
	"time"

	spacesStore "github.com/fancyinnovations/fancyspaces/core/internal/spaces"
	"github.com/fancyinnovations/fancyspaces/core/internal/spaces/database/fake"
	"github.com/fancyinnovations/fancyspaces/integrations/idp-go-sdk/idp"
	"github.com/fancyinnovations/fancyspaces/integrations/spaces-go-sdk/spaces"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func TestStore_GetByID(t *testing.T) {
	now := time.Now()

	for _, tc := range []struct {
		Name    string
		Exiting []spaces.Space
		ID      string
		Exp     *spaces.Space
		ExpErr  error
	}{
		{
			Name: "Found space",
			Exiting: []spaces.Space{
				{
					ID:          "space-1",
					Slug:        "spaceOne",
					Title:       "Space One",
					Description: "This is the first space.",
					Categories:  []spaces.Category{spaces.CategoryMinecraftPlugin},
					IconURL:     "https://example.com/icon1.png",
					Status:      spaces.StatusApproved,
					CreatedAt:   now,
					Members: []spaces.Member{
						{UserID: "user-1", Role: spaces.RoleAdmin},
					},
				},
			},
			ID: "space-1",
			Exp: &spaces.Space{
				ID:          "space-1",
				Slug:        "spaceOne",
				Title:       "Space One",
				Description: "This is the first space.",
				Categories:  []spaces.Category{spaces.CategoryMinecraftPlugin},
				IconURL:     "https://example.com/icon1.png",
				Status:      spaces.StatusApproved,
				CreatedAt:   now,
				Members: []spaces.Member{
					{UserID: "user-1", Role: spaces.RoleAdmin},
				},
			},
			ExpErr: nil,
		},
		{
			Name: "Found space, multiple exists",
			Exiting: []spaces.Space{
				{
					ID:          "space-1",
					Slug:        "spaceOne",
					Title:       "Space One",
					Description: "This is the first space.",
					Categories:  []spaces.Category{spaces.CategoryMinecraftPlugin},
					IconURL:     "https://example.com/icon1.png",
					Status:      spaces.StatusApproved,
					CreatedAt:   now,
					Members: []spaces.Member{
						{UserID: "user-1", Role: spaces.RoleAdmin},
					},
				},
				{
					ID:          "space-2",
					Slug:        "spaceTwo",
					Title:       "Space Two",
					Description: "This is the second space.",
					Categories:  []spaces.Category{spaces.CategoryMinecraftMod},
					IconURL:     "https://example.com/icon2.png",
					Status:      spaces.StatusApproved,
					CreatedAt:   now,
					Members: []spaces.Member{
						{UserID: "user-2", Role: spaces.RoleAdmin},
					},
				},
				{
					ID:          "space-3",
					Slug:        "spaceThree",
					Title:       "Space Three",
					Description: "This is the third space.",
					Categories:  []spaces.Category{spaces.CategoryMinecraftServer},
					IconURL:     "https://example.com/icon3.png",
					Status:      spaces.StatusApproved,
					CreatedAt:   now,
					Members: []spaces.Member{
						{UserID: "user-3", Role: spaces.RoleAdmin},
					},
				},
			},
			ID: "space-2",
			Exp: &spaces.Space{
				ID:          "space-2",
				Slug:        "spaceTwo",
				Title:       "Space Two",
				Description: "This is the second space.",
				Categories:  []spaces.Category{spaces.CategoryMinecraftMod},
				IconURL:     "https://example.com/icon2.png",
				Status:      spaces.StatusApproved,
				CreatedAt:   now,
				Members: []spaces.Member{
					{UserID: "user-2", Role: spaces.RoleAdmin},
				},
			},
			ExpErr: nil,
		},
		{
			Name:    "No spaces exist",
			Exiting: []spaces.Space{},
			ID:      "space-1",
			Exp:     nil,
			ExpErr:  spaces.ErrSpaceNotFound,
		},
		{
			Name: "Space not found",
			Exiting: []spaces.Space{
				{
					ID:          "space-1",
					Slug:        "spaceOne",
					Title:       "Space One",
					Description: "This is the first space.",
					Categories:  []spaces.Category{spaces.CategoryMinecraftPlugin},
					IconURL:     "https://example.com/icon1.png",
					Status:      spaces.StatusApproved,
					CreatedAt:   now,
					Members: []spaces.Member{
						{UserID: "user-1", Role: spaces.RoleAdmin},
					},
				},
			},
			ID:     "space-2",
			Exp:    nil,
			ExpErr: spaces.ErrSpaceNotFound,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			db := fake.New()
			store := spacesStore.New(spacesStore.Configuration{
				DB: db,
			})

			for _, s := range tc.Exiting {
				if err := db.Create(&s); err != nil {
					t.Fatalf("Could not setup db: %v", err)
				}
			}

			got, err := store.GetByID(tc.ID)
			if !errors.Is(err, tc.ExpErr) {
				t.Fatalf("Expected error %v, got %v", tc.ExpErr, err)
			}

			if diff := cmp.Diff(tc.Exp, got); diff != "" {
				t.Errorf("Unexpected diff: %s", diff)
			}
		})
	}
}

func TestStore_GetBySlug(t *testing.T) {
	now := time.Now()

	tests := []struct {
		Name    string
		Exiting []spaces.Space
		Slug    string
		Exp     *spaces.Space
		ExpErr  error
	}{
		{
			Name: "Found space by slug",
			Exiting: []spaces.Space{
				{
					ID:          "space-1",
					Slug:        "spaceOne",
					Title:       "Space One",
					Description: "First space.",
					Categories:  []spaces.Category{spaces.CategoryMinecraftPlugin},
					IconURL:     "https://example.com/icon1.png",
					Status:      spaces.StatusApproved,
					CreatedAt:   now,
					Members:     []spaces.Member{{UserID: "user-1", Role: spaces.RoleAdmin}},
				},
			},
			Slug: "spaceOne",
			Exp: &spaces.Space{
				ID:          "space-1",
				Slug:        "spaceOne",
				Title:       "Space One",
				Description: "First space.",
				Categories:  []spaces.Category{spaces.CategoryMinecraftPlugin},
				IconURL:     "https://example.com/icon1.png",
				Status:      spaces.StatusApproved,
				CreatedAt:   now,
				Members:     []spaces.Member{{UserID: "user-1", Role: spaces.RoleAdmin}},
			},
			ExpErr: nil,
		},
		{
			Name: "Multiple spaces, find by slug",
			Exiting: []spaces.Space{
				{
					ID:          "space-1",
					Slug:        "spaceOne",
					Title:       "Space One",
					Description: "First space.",
					Categories:  []spaces.Category{spaces.CategoryMinecraftPlugin},
					IconURL:     "https://example.com/icon1.png",
					Status:      spaces.StatusApproved,
					CreatedAt:   now,
					Members:     []spaces.Member{{UserID: "user-1", Role: spaces.RoleAdmin}},
				},
				{
					ID:          "space-2",
					Slug:        "spaceTwo",
					Title:       "Space Two",
					Description: "Second space.",
					Categories:  []spaces.Category{spaces.CategoryMinecraftMod},
					IconURL:     "https://example.com/icon2.png",
					Status:      spaces.StatusApproved,
					CreatedAt:   now,
					Members:     []spaces.Member{{UserID: "user-2", Role: spaces.RoleAdmin}},
				},
			},
			Slug: "spaceTwo",
			Exp: &spaces.Space{
				ID:          "space-2",
				Slug:        "spaceTwo",
				Title:       "Space Two",
				Description: "Second space.",
				Categories:  []spaces.Category{spaces.CategoryMinecraftMod},
				IconURL:     "https://example.com/icon2.png",
				Status:      spaces.StatusApproved,
				CreatedAt:   now,
				Members:     []spaces.Member{{UserID: "user-2", Role: spaces.RoleAdmin}},
			},
			ExpErr: nil,
		},
		{
			Name:    "No spaces exist",
			Exiting: []spaces.Space{},
			Slug:    "spaceOne",
			Exp:     nil,
			ExpErr:  spaces.ErrSpaceNotFound,
		},
		{
			Name: "Space not found by slug",
			Exiting: []spaces.Space{
				{
					ID:          "space-1",
					Slug:        "spaceOne",
					Title:       "Space One",
					Description: "First space.",
					Categories:  []spaces.Category{spaces.CategoryMinecraftPlugin},
					IconURL:     "https://example.com/icon1.png",
					Status:      spaces.StatusApproved,
					CreatedAt:   now,
					Members:     []spaces.Member{{UserID: "user-1", Role: spaces.RoleAdmin}},
				},
			},
			Slug:   "spaceTwo",
			Exp:    nil,
			ExpErr: spaces.ErrSpaceNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			db := fake.New()
			store := spacesStore.New(spacesStore.Configuration{DB: db})

			for _, s := range tc.Exiting {
				if err := db.Create(&s); err != nil {
					t.Fatalf("Could not setup db: %v", err)
				}
			}

			got, err := store.GetBySlug(tc.Slug)
			if !errors.Is(err, tc.ExpErr) {
				t.Fatalf("Expected error %v, got %v", tc.ExpErr, err)
			}

			if diff := cmp.Diff(tc.Exp, got); diff != "" {
				t.Errorf("Unexpected diff: %s", diff)
			}
		})
	}
}

func TestStore_GetAll(t *testing.T) {
	now := time.Now()

	for _, tc := range []struct {
		Name    string
		Exiting []spaces.Space
	}{
		{
			Name:    "No spaces exist",
			Exiting: []spaces.Space{},
		},
		{
			Name: "One space exists",
			Exiting: []spaces.Space{
				{
					ID:          "space-1",
					Slug:        "spaceOne",
					Title:       "Space One",
					Description: "This is the first space.",
					Categories:  []spaces.Category{spaces.CategoryMinecraftPlugin},
					IconURL:     "https://example.com/icon1.png",
					Status:      spaces.StatusDraft,
					CreatedAt:   now,
					Members: []spaces.Member{
						{UserID: "user-1", Role: spaces.RoleAdmin},
					},
				},
			},
		},
		{
			Name: "Multiple spaces exists",
			Exiting: []spaces.Space{
				{
					ID:          "space-1",
					Slug:        "spaceOne",
					Title:       "Space One",
					Description: "This is the first space.",
					Categories:  []spaces.Category{spaces.CategoryMinecraftPlugin},
					IconURL:     "https://example.com/icon1.png",
					Status:      spaces.StatusDraft,
					CreatedAt:   now,
					Members: []spaces.Member{
						{UserID: "user-1", Role: spaces.RoleAdmin},
					},
				},
				{
					ID:          "space-2",
					Slug:        "spaceTwo",
					Title:       "Space Two",
					Description: "This is the second space.",
					Categories:  []spaces.Category{spaces.CategoryMinecraftMod},
					IconURL:     "https://example.com/icon2.png",
					Status:      spaces.StatusApproved,
					CreatedAt:   now,
					Members: []spaces.Member{
						{UserID: "user-2", Role: spaces.RoleAdmin},
					},
				},
				{
					ID:          "space-3",
					Slug:        "spaceThree",
					Title:       "Space Three",
					Description: "This is the third space.",
					Categories:  []spaces.Category{spaces.CategoryOther},
					IconURL:     "https://example.com/icon3.png",
					Status:      spaces.StatusArchived,
					CreatedAt:   now,
					Members: []spaces.Member{
						{UserID: "user-3", Role: spaces.RoleAdmin},
					},
				},
			},
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			db := fake.New()
			store := spacesStore.New(spacesStore.Configuration{
				DB: db,
			})

			for _, s := range tc.Exiting {
				if err := db.Create(&s); err != nil {
					t.Fatalf("Could not setup db: %v", err)
				}
			}

			got, err := store.GetAll()
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if diff := cmp.Diff(tc.Exiting, got); diff != "" {
				t.Errorf("Unexpected diff: %s", diff)
			}
		})
	}
}

func TestStore_Create(t *testing.T) {
	now := time.Now()

	normalUser := idp.User{
		ID:        "user-1",
		Provider:  idp.ProviderBasic,
		Name:      "User",
		Email:     "user@fancyspaces.net",
		Verified:  true,
		Password:  "pwd",
		Roles:     []string{"user"},
		CreatedAt: now,
		IsActive:  true,
		Metadata:  map[string]string{},
	}

	for _, tc := range []struct {
		Name    string
		Exiting []spaces.Space
		Creator *idp.User
		Req     *spacesStore.CreateOrUpdateSpaceReq
		Exp     *spaces.Space
		ExpErr  error
	}{
		{
			Name:    "Successful creation",
			Exiting: []spaces.Space{},
			Creator: &normalUser,
			Req: &spacesStore.CreateOrUpdateSpaceReq{
				Slug:        "space-one",
				Title:       "Space One",
				Description: "This is the first space.",
				Categories:  []spaces.Category{spaces.CategoryMinecraftPlugin},
				IconURL:     "https://example.com/icon1.png",
			},
			Exp: &spaces.Space{
				ID:          "",
				Slug:        "space-one",
				Title:       "Space One",
				Description: "This is the first space.",
				Categories:  []spaces.Category{spaces.CategoryMinecraftPlugin},
				IconURL:     "https://example.com/icon1.png",
				Status:      spaces.StatusDraft,
				CreatedAt:   now,
				Creator:     "user-1",
				Members:     []spaces.Member{},
				Links:       []spaces.Link{},
				AnalyticsSettings: spaces.AnalyticsSettings{
					Enabled: true,
				},
			},
			ExpErr: nil,
		},
		{
			Name:    "User is not active",
			Exiting: []spaces.Space{},
			Creator: &idp.User{
				ID:       "user-2",
				Verified: true,
				IsActive: false,
				Roles:    []string{"user"},
			},
			Req: &spacesStore.CreateOrUpdateSpaceReq{
				Slug:        "space-one",
				Title:       "Space One",
				Description: "This is the first space.",
				Categories:  []spaces.Category{spaces.CategoryMinecraftPlugin},
				IconURL:     "https://example.com/icon1.png",
			},
			Exp:    nil,
			ExpErr: spacesStore.ErrUserNotActive,
		},
		{
			Name:    "User is not verified",
			Exiting: []spaces.Space{},
			Creator: &idp.User{
				ID:       "user-2",
				Verified: false,
				IsActive: true,
				Roles:    []string{"user"},
			},
			Req: &spacesStore.CreateOrUpdateSpaceReq{
				Slug:        "space-one",
				Title:       "Space One",
				Description: "This is the first space.",
				Categories:  []spaces.Category{spaces.CategoryMinecraftPlugin},
				IconURL:     "https://example.com/icon1.png",
			},
			Exp:    nil,
			ExpErr: spacesStore.ErrUserNotVerified,
		},
		{
			Name: "Slug is already taken",
			Exiting: []spaces.Space{
				{
					ID:          uuid.New().String(),
					Slug:        "space-one",
					Title:       "Space One",
					Description: "This is the first space.",
					Categories:  []spaces.Category{spaces.CategoryMinecraftPlugin},
					IconURL:     "https://example.com/icon1.png",
					Status:      spaces.StatusDraft,
					CreatedAt:   now,
					Creator:     "user-1",
					Members:     []spaces.Member{},
				},
			},
			Creator: &normalUser,
			Req: &spacesStore.CreateOrUpdateSpaceReq{
				Slug:        "space-one",
				Title:       "Space One",
				Description: "This is the first space.",
				Categories:  []spaces.Category{spaces.CategoryMinecraftPlugin},
				IconURL:     "https://example.com/icon1.png",
			},
			Exp:    nil,
			ExpErr: spacesStore.ErrSpaceAlreadyExists,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			db := fake.New()
			store := spacesStore.New(spacesStore.Configuration{
				DB: db,
			})

			for _, s := range tc.Exiting {
				if err := db.Create(&s); err != nil {
					t.Fatalf("Could not setup db: %v", err)
				}
			}

			got, err := store.Create(tc.Creator, tc.Req)
			if !errors.Is(err, tc.ExpErr) {
				t.Fatalf("Expected error %v, got %v", tc.ExpErr, err)
			}

			if tc.Exp != nil {
				if len(got.ID) == 0 {
					t.Errorf("Expected non-empty ID, got empty")
				}
				if got.CreatedAt.IsZero() {
					t.Errorf("Expected non-zero CreatedAt, got zero")
				}

				// Normalize dynamic fields for comparison
				got.ID = ""
				got.CreatedAt = now
			}

			if diff := cmp.Diff(tc.Exp, got); diff != "" {
				t.Errorf("Unexpected diff: %s", diff)
			}
		})
	}
}

func TestStore_Update(t *testing.T) {
	now := time.Now()

	for _, tc := range []struct {
		Name    string
		Exiting []spaces.Space
		ID      string
		Req     *spacesStore.CreateOrUpdateSpaceReq
		Exp     *spaces.Space
		ExpErr  error
	}{
		{
			Name: "Successful update",
			Exiting: []spaces.Space{
				{
					ID:          "space-1",
					Slug:        "my-space",
					Title:       "Space One",
					Description: "This is the first space.",
					Categories:  []spaces.Category{spaces.CategoryMinecraftPlugin},
					IconURL:     "https://example.com/icon1.png",
					Status:      spaces.StatusDraft,
					CreatedAt:   now,
					Creator:     "user-1",
					Members:     []spaces.Member{},
				},
			},
			ID: "space-1",
			Req: &spacesStore.CreateOrUpdateSpaceReq{
				Slug:        "my-space",
				Title:       "Space Two",
				Description: "This is the second space.",
				Categories:  []spaces.Category{spaces.CategoryMinecraftPlugin, spaces.CategoryWebApp},
				IconURL:     "https://example.com/icon2.png",
			},
			Exp: &spaces.Space{
				ID:          "space-1",
				Slug:        "my-space",
				Title:       "Space Two",
				Description: "This is the second space.",
				Categories:  []spaces.Category{spaces.CategoryMinecraftPlugin, spaces.CategoryWebApp},
				IconURL:     "https://example.com/icon2.png",
				Status:      spaces.StatusDraft,
				CreatedAt:   now,
				Creator:     "user-1",
				Members:     []spaces.Member{},
			},
			ExpErr: nil,
		},
		{
			Name: "Slug is already taken",
			Exiting: []spaces.Space{
				{
					ID:          "space-1",
					Slug:        "space-one",
					Title:       "Space One",
					Description: "This is the first space.",
					Categories:  []spaces.Category{spaces.CategoryMinecraftPlugin},
					IconURL:     "https://example.com/icon1.png",
					Status:      spaces.StatusDraft,
					CreatedAt:   now,
					Creator:     "user-1",
					Members:     []spaces.Member{},
				},
				{
					ID:          "space-2",
					Slug:        "space-two",
					Title:       "Space Two",
					Description: "This is the second space.",
					Categories:  []spaces.Category{spaces.CategoryMinecraftMod},
					IconURL:     "https://example.com/icon2.png",
					Status:      spaces.StatusApproved,
					CreatedAt:   now,
					Creator:     "user-2",
					Members:     []spaces.Member{},
				},
			},
			ID: "space-1",
			Req: &spacesStore.CreateOrUpdateSpaceReq{
				Slug:        "space-two",
				Title:       "Space One",
				Description: "This is the first space.",
				Categories:  []spaces.Category{spaces.CategoryMinecraftPlugin},
				IconURL:     "https://example.com/icon1.png",
			},
			Exp:    nil,
			ExpErr: spacesStore.ErrSpaceAlreadyExists,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			db := fake.New()
			store := spacesStore.New(spacesStore.Configuration{
				DB: db,
			})

			for _, s := range tc.Exiting {
				if err := db.Create(&s); err != nil {
					t.Fatalf("Could not setup db: %v", err)
				}
			}

			err := store.Update(tc.ID, tc.Req)
			if !errors.Is(err, tc.ExpErr) {
				t.Fatalf("Expected error %v, got %v", tc.ExpErr, err)
			}

			if tc.Exp != nil {
				got, err := db.GetByID(tc.ID)
				if err != nil {
					t.Fatalf("Could not fetch updated space: %v", err)
				}
				if diff := cmp.Diff(tc.Exp, got); diff != "" {
					t.Errorf("Unexpected diff: %s", diff)
				}
			}
		})
	}
}

func TestStore_Delete(t *testing.T) {
	now := time.Now()

	for _, tc := range []struct {
		Name    string
		Exiting []spaces.Space
		ID      string
		ExpErr  error
	}{
		{
			Name: "Successful deletion",
			Exiting: []spaces.Space{
				{
					ID:          "space-1",
					Slug:        "spaceOne",
					Title:       "Space One",
					Description: "This is the first space.",
					Categories:  []spaces.Category{spaces.CategoryMinecraftPlugin},
					IconURL:     "https://example.com/icon1.png",
					Status:      spaces.StatusDraft,
					CreatedAt:   now,
					Members: []spaces.Member{
						{UserID: "user-1", Role: spaces.RoleAdmin},
					},
				},
			},
			ID:     "space-1",
			ExpErr: nil,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			db := fake.New()
			store := spacesStore.New(spacesStore.Configuration{
				DB: db,
			})

			for _, s := range tc.Exiting {
				if err := db.Create(&s); err != nil {
					t.Fatalf("Could not setup db: %v", err)
				}
			}

			err := store.Delete(tc.ID)
			if !errors.Is(err, tc.ExpErr) {
				t.Fatalf("Expected error %v, got %v", tc.ExpErr, err)
			}

			if tc.ExpErr == nil {
				_, err := db.GetByID(tc.ID)
				if !errors.Is(err, spaces.ErrSpaceNotFound) {
					t.Errorf("Expected space to be deleted, but it still exists")
				}
			}
		})
	}
}

func TestStore_ChangeStatus(t *testing.T) {
	now := time.Now()

	for _, tc := range []struct {
		Name   string
		Space  spaces.Space
		To     spaces.Status
		Exp    spaces.Status
		ExpErr error
	}{
		{
			Name: "Successful change status from draft to review",
			Space: spaces.Space{
				ID:          "space-1",
				Slug:        "space-one",
				Title:       "Space One",
				Description: "This is the first space.",
				Categories:  []spaces.Category{spaces.CategoryMinecraftPlugin},
				IconURL:     "https://example.com/icon1.png",
				Status:      spaces.StatusDraft,
				CreatedAt:   now,
				Members: []spaces.Member{
					{UserID: "user-1", Role: spaces.RoleAdmin},
				},
			},
			To:     spaces.StatusReview,
			Exp:    spaces.StatusReview,
			ExpErr: nil,
		},
		{
			Name: "Successful change status from review to approved",
			Space: spaces.Space{
				ID:          "space-1",
				Slug:        "space-one",
				Title:       "Space One",
				Description: "This is the first space.",
				Categories:  []spaces.Category{spaces.CategoryMinecraftPlugin},
				IconURL:     "https://example.com/icon1.png",
				Status:      spaces.StatusReview,
				CreatedAt:   now,
				Members: []spaces.Member{
					{UserID: "user-1", Role: spaces.RoleAdmin},
				},
			},
			To:     spaces.StatusApproved,
			Exp:    spaces.StatusApproved,
			ExpErr: nil,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			db := fake.New()
			store := spacesStore.New(spacesStore.Configuration{
				DB: db,
			})

			if err := db.Create(&tc.Space); err != nil {
				t.Fatalf("Could not setup db: %v", err)
			}

			err := store.ChangeStatus(&tc.Space, tc.To)
			if !errors.Is(err, tc.ExpErr) {
				t.Fatalf("Expected error %v, got %v", tc.ExpErr, err)
			}

			if tc.ExpErr == nil {
				if tc.Space.Status != tc.Exp {
					t.Errorf("Expected status %v, got %v", tc.Exp, tc.Space.Status)
				}
			}
		})
	}
}
