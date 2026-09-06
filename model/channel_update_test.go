package model

import (
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// setupIsolatedDB returns an in-memory sqlite DB with only the Channel and
// Ability tables migrated, then assigns it to the package-level DB so that
// Channel.Update() exercises production code paths.
func setupIsolatedDB(t *testing.T) func() {
	originalDB := DB
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&Channel{}, &Ability{}); err != nil {
		t.Fatal(err)
	}
	DB = db
	return func() { DB = originalDB }
}

func boolPtr(v bool) *bool       { return &v }
func int64Ptr(v int64) *int64   { return &v }
func strPtr(v string) *string    { return &v }
func intPtr(v int) *int          { return &v }
func uintPtr(v uint) *uint       { return &v }
func stringPtrSliceEmpty() *string { s := ""; return &s }

// allFieldsPresent mirrors a full edit modal payload: every mutable JSON key
// is sent, so Update() writes all of them.
func allFieldsPresent() map[string]bool {
	return map[string]bool{
		"type":              true,
		"name":              true,
		"key":               true,
		"status":            true,
		"base_url":          true,
		"models":            true,
		"group":             true,
		"model_mapping":     true,
		"priority":          true,
		"max_concurrency":   true,
		"rpm":               true,
		"is_fallback":       true,
		"fallback_priority": true,
	}
}

func TestChannelUpdate_PersistsAllFallbackFields(t *testing.T) {
	defer setupIsolatedDB(t)()

	// Seed: a normal (non-fallback) channel with a known key + base_url.
	original := &Channel{
		Id:             1,
		Type:           1,
		Name:           "OriginalName",
		Key:            "sk-original-key",
		Status:         ChannelStatusEnabled,
		BaseURL:        strPtr("https://api.example.com"),
		Models:         "gpt-4",
		Group:          "default",
		ModelMapping:   stringPtrSliceEmpty(),
		Priority:       int64Ptr(0),
		MaxConcurrency: intPtr(0),
		RPM:            intPtr(0),
		IsFallback:     boolPtr(false),
		FallbackPriority: int64Ptr(0),
	}
	if err := DB.Create(original).Error; err != nil {
		t.Fatalf("seed insert: %v", err)
	}

	// User toggles the channel into a fallback one with priority 7.
	edit := &Channel{
		Id:               1,
		Type:             1,
		Name:             "OriginalName",
		Key:              "sk-original-key", // omitted by frontend, would be zero on a fresh bind
		Status:           ChannelStatusEnabled,
		BaseURL:          strPtr("https://api.example.com"),
		Models:           "gpt-4",
		Group:            "default",
		ModelMapping:     stringPtrSliceEmpty(),
		Priority:         int64Ptr(0),
		MaxConcurrency:   intPtr(0),
		RPM:              intPtr(0),
		IsFallback:       boolPtr(true),
		FallbackPriority: int64Ptr(7),
	}
	if err := edit.Update(allFieldsPresent()); err != nil {
		t.Fatalf("Update: %v", err)
	}

	var got Channel
	if err := DB.First(&got, 1).Error; err != nil {
		t.Fatalf("read after update: %v", err)
	}
	if !got.GetIsFallback() {
		t.Errorf("IsFallback=true not persisted (got false)")
	}
	if got.GetFallbackPriority() != 7 {
		t.Errorf("FallbackPriority=7 not persisted (got %d)", got.GetFallbackPriority())
	}
}

func TestChannelUpdate_ToggleFallbackOff(t *testing.T) {
	defer setupIsolatedDB(t)()

	// Seed: a channel that *is* a fallback (mirroring the user's second
	// scenario: edit a fallback channel back to normal).
	isFb := true
	fp := int64(20)
	original := &Channel{
		Id:               2,
		Type:             1,
		Name:             "ExistingFallback",
		Key:              "sk-fb-key",
		Status:           ChannelStatusEnabled,
		BaseURL:          strPtr("https://api.example.com"),
		Models:           "gpt-4",
		Group:            "default",
		ModelMapping:     stringPtrSliceEmpty(),
		Priority:         int64Ptr(0),
		MaxConcurrency:   intPtr(0),
		RPM:              intPtr(0),
		IsFallback:       &isFb,
		FallbackPriority: &fp,
	}
	if err := DB.Create(original).Error; err != nil {
		t.Fatalf("seed insert: %v", err)
	}

	// User toggles it off and clears the priority.
	edit := &Channel{
		Id:               2,
		Type:             1,
		Name:             "ExistingFallback",
		Key:              "sk-fb-key",
		Status:           ChannelStatusEnabled,
		BaseURL:          strPtr("https://api.example.com"),
		Models:           "gpt-4",
		Group:            "default",
		ModelMapping:     stringPtrSliceEmpty(),
		Priority:         int64Ptr(0),
		MaxConcurrency:   intPtr(0),
		RPM:              intPtr(0),
		IsFallback:       boolPtr(false),
		FallbackPriority: int64Ptr(0),
	}
	if err := edit.Update(allFieldsPresent()); err != nil {
		t.Fatalf("Update: %v", err)
	}

	var got Channel
	if err := DB.First(&got, 2).Error; err != nil {
		t.Fatalf("read after update: %v", err)
	}
	if got.GetIsFallback() {
		t.Errorf("IsFallback=false not persisted (still true)")
	}
	if got.GetFallbackPriority() != 0 {
		t.Errorf("FallbackPriority=0 not persisted (still %d)", got.GetFallbackPriority())
	}
}

func TestChannelUpdate_PreservesUntouchedFields(t *testing.T) {
	defer setupIsolatedDB(t)()

	// Seed: existing channel with specific Group + Models.
	original := &Channel{
		Id:             3,
		Type:           1,
		Name:           "Sensitive",
		Key:            "sk-real-key",
		Status:         ChannelStatusEnabled,
		BaseURL:        strPtr("https://api.example.com"),
		Models:         "gpt-4,gpt-3.5-turbo",
		Group:          "vip,default",
		ModelMapping:   stringPtrSliceEmpty(),
		Priority:       int64Ptr(5),
		MaxConcurrency: intPtr(0),
		RPM:            intPtr(0),
		IsFallback:     boolPtr(false),
		FallbackPriority: int64Ptr(0),
	}
	if err := DB.Create(original).Error; err != nil {
		t.Fatalf("seed: %v", err)
	}

	// User only toggles fallback; everything else should stay as-is.
	edit := &Channel{
		Id:               3,
		Type:             1,
		Name:             "Sensitive",
		Key:              "sk-real-key",
		Status:           ChannelStatusEnabled,
		BaseURL:          strPtr("https://api.example.com"),
		Models:           "gpt-4,gpt-3.5-turbo",
		Group:            "vip,default",
		ModelMapping:     stringPtrSliceEmpty(),
		Priority:         int64Ptr(5),
		MaxConcurrency:   intPtr(0),
		RPM:              intPtr(0),
		IsFallback:       boolPtr(true),
		FallbackPriority: int64Ptr(99),
	}
	if err := edit.Update(allFieldsPresent()); err != nil {
		t.Fatalf("Update: %v", err)
	}

	var got Channel
	if err := DB.First(&got, 3).Error; err != nil {
		t.Fatalf("read: %v", err)
	}
	if got.Key != "sk-real-key" {
		t.Errorf("Key was clobbered: got %q", got.Key)
	}
	if got.Models != "gpt-4,gpt-3.5-turbo" {
		t.Errorf("Models clobbered: got %q", got.Models)
	}
	if got.Group != "vip,default" {
		t.Errorf("Group clobbered: got %q", got.Group)
	}
	if got.GetPriority() != 5 {
		t.Errorf("Priority clobbered: got %d", got.GetPriority())
	}
	if !got.GetIsFallback() {
		t.Errorf("IsFallback=true not persisted")
	}
	if got.GetFallbackPriority() != 99 {
		t.Errorf("FallbackPriority=99 not persisted, got %d", got.GetFallbackPriority())
	}
}

// TestChannelUpdate_PreservesNilPointerFields guards against a regression
// where the Update() map-based patch would set untouched *T pointer fields
// to NULL instead of leaving them alone. The classic repro: user opens an
// edit modal, doesn't touch Weight / Priority / MaxConcurrency / RPM /
// BaseURL / ModelMapping, toggles only is_fallback — those untouched
// columns must keep their pre-existing values.
func TestChannelUpdate_PreservesNilPointerFields(t *testing.T) {
	defer setupIsolatedDB(t)()

	original := &Channel{
		Id:               4,
		Type:             1,
		Name:             "NilPointerTest",
		Key:              "sk-nil-test",
		Status:           ChannelStatusEnabled,
		BaseURL:          strPtr("https://orig.example.com"),
		Models:           "gpt-4",
		Group:            "vip",
		ModelMapping:     strPtr("{\"gpt-4\":\"gpt-4-turbo\"}"),
		Priority:         int64Ptr(7),
		MaxConcurrency:   intPtr(3),
		RPM:              intPtr(60),
		IsFallback:       boolPtr(false),
		FallbackPriority: int64Ptr(0),
	}
	if err := DB.Create(original).Error; err != nil {
		t.Fatalf("seed: %v", err)
	}

	// Simulate an edit where ONLY the fallback fields are populated; every
	// other pointer field is left nil (i.e. the JSON omitted them).
	edit := &Channel{
		Id:               4,
		Type:             1,
		Name:             "NilPointerTest",
		Key:              "sk-nil-test",
		Status:           ChannelStatusEnabled,
		Models:           "gpt-4",
		Group:            "vip",
		IsFallback:       boolPtr(true),
		FallbackPriority: int64Ptr(8),
	}
	provided := map[string]bool{
		"id": true, "type": true, "name": true, "key": true, "status": true,
		"models": true, "group": true, "is_fallback": true, "fallback_priority": true,
	}
	if err := edit.Update(provided); err != nil {
		t.Fatalf("Update: %v", err)
	}

	var got Channel
	if err := DB.First(&got, 4).Error; err != nil {
		t.Fatalf("read: %v", err)
	}
	if got.GetBaseURL() != "https://orig.example.com" {
		t.Errorf("BaseURL clobbered to NULL, got %q", got.GetBaseURL())
	}
	if got.GetModelMapping() == nil {
		t.Errorf("ModelMapping clobbered to NULL")
	}
	if got.GetPriority() != 7 {
		t.Errorf("Priority clobbered, got %d", got.GetPriority())
	}
	if got.GetMaxConcurrency() != 3 {
		t.Errorf("MaxConcurrency clobbered, got %d", got.GetMaxConcurrency())
	}
	if got.GetRPM() != 60 {
		t.Errorf("RPM clobbered, got %d", got.GetRPM())
	}
	if !got.GetIsFallback() {
		t.Errorf("IsFallback toggle not applied")
	}
	if got.GetFallbackPriority() != 8 {
		t.Errorf("FallbackPriority not applied, got %d", got.GetFallbackPriority())
	}
}

// TestChannelUpdate_FrontendPayload mimics the exact JSON body the
// Channel.vue edit modal sends via PUT /api/channel/. The list endpoint
// Omit("key") so the frontend always posts key=""; we verify the real key
// is preserved AND the fallback toggle persists.
func TestChannelUpdate_FrontendPayload(t *testing.T) {
	defer setupIsolatedDB(t)()

	original := &Channel{
		Id:               5,
		Type:             1,
		Name:             "FrontendPayload",
		Key:              "sk-frontend-test",
		Status:           ChannelStatusEnabled,
		BaseURL:          strPtr("https://api.example.com"),
		Models:           "gpt-4",
		Group:            "default",
		ModelMapping:     stringPtrSliceEmpty(),
		Priority:         int64Ptr(0),
		MaxConcurrency:   intPtr(0),
		RPM:              intPtr(0),
		IsFallback:       boolPtr(false),
		FallbackPriority: int64Ptr(0),
	}
	if err := DB.Create(original).Error; err != nil {
		t.Fatalf("seed: %v", err)
	}

	// Body as Channel.vue would send it (note: key="" because frontend
	// never sees the real value).
	body := &Channel{
		Id:               5,
		Type:             1,
		Name:             "FrontendPayload",
		Key:              "", // <-- frontend always sends empty
		Status:           ChannelStatusEnabled,
		BaseURL:          strPtr("https://api.example.com"),
		Models:           "gpt-4",
		Group:            "default",
		ModelMapping:     stringPtrSliceEmpty(),
		IsFallback:       boolPtr(true),
		FallbackPriority: int64Ptr(11),
	}
	provided := map[string]bool{
		"id": true, "type": true, "name": true, "key": true, "status": true,
		"base_url": true, "models": true, "group": true, "model_mapping": true,
		"is_fallback": true, "fallback_priority": true,
	}
	if err := body.Update(provided); err != nil {
		t.Fatalf("Update: %v", err)
	}

	var got Channel
	if err := DB.First(&got, 5).Error; err != nil {
		t.Fatalf("read: %v", err)
	}
	if got.Key != "sk-frontend-test" {
		t.Errorf("Key was clobbered by empty key from frontend, got %q", got.Key)
	}
	if !got.GetIsFallback() {
		t.Errorf("IsFallback=true not persisted via frontend-shaped payload")
	}
	if got.GetFallbackPriority() != 11 {
		t.Errorf("FallbackPriority=11 not persisted, got %d", got.GetFallbackPriority())
	}
}

// TestChannelUpdate_PartialStatusToggle reproduces the reported bug:
// Channel.vue toggles enable/disable by PUT /api/channel/ with body
// {id, status}. The update must only change the `status` column and must NOT
// wipe the channel's other fields (name/models/group/config/...) to empty.
func TestChannelUpdate_PartialStatusToggle(t *testing.T) {
	defer setupIsolatedDB(t)()

	original := &Channel{
		Id:               6,
		Type:             1,
		Name:             "KeptName",
		Key:              "sk-keep",
		Status:           ChannelStatusEnabled,
		BaseURL:          strPtr("https://api.example.com"),
		Models:           "gpt-4",
		Group:            "vip,default",
		Config:           `{"region":"us"}`,
		CooldownSeconds:  30,
		Balance:          12.5,
		ModelMapping:     stringPtrSliceEmpty(),
		Priority:         int64Ptr(2),
		MaxConcurrency:   intPtr(4),
		RPM:              intPtr(100),
		IsFallback:       boolPtr(false),
		FallbackPriority: int64Ptr(0),
	}
	if err := DB.Create(original).Error; err != nil {
		t.Fatalf("seed: %v", err)
	}

	// Only {id, status} is sent by the toggle button.
	body := &Channel{Id: 6, Status: ChannelStatusManuallyDisabled}
	provided := map[string]bool{"id": true, "status": true}
	if err := body.Update(provided); err != nil {
		t.Fatalf("Update: %v", err)
	}

	var got Channel
	if err := DB.First(&got, 6).Error; err != nil {
		t.Fatalf("read: %v", err)
	}
	if got.Status != ChannelStatusManuallyDisabled {
		t.Errorf("Status not updated to %d, got %d", ChannelStatusManuallyDisabled, got.Status)
	}
	if got.Name != "KeptName" {
		t.Errorf("Name clobbered to %q", got.Name)
	}
	if got.Models != "gpt-4" {
		t.Errorf("Models clobbered to %q", got.Models)
	}
	if got.Group != "vip,default" {
		t.Errorf("Group clobbered to %q", got.Group)
	}
	if got.Config != `{"region":"us"}` {
		t.Errorf("Config clobbered to %q", got.Config)
	}
	if got.CooldownSeconds != 30 {
		t.Errorf("CooldownSeconds clobbered to %d", got.CooldownSeconds)
	}
	if got.Balance != 12.5 {
		t.Errorf("Balance clobbered to %v", got.Balance)
	}
	if got.Key != "sk-keep" {
		t.Errorf("Key clobbered to %q", got.Key)
	}
	if got.GetBaseURL() != "https://api.example.com" {
		t.Errorf("BaseURL clobbered to %q", got.GetBaseURL())
	}
	if got.GetPriority() != 2 {
		t.Errorf("Priority clobbered to %d", got.GetPriority())
	}
}