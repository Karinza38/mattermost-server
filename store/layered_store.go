// Copyright (c) 2016-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package store

import (
	"context"

	"github.com/mattermost/mattermost-server/einterfaces"
	"github.com/mattermost/mattermost-server/mlog"
	"github.com/mattermost/mattermost-server/model"
)

const (
	ENABLE_EXPERIMENTAL_REDIS = false
)

type LayeredStoreDatabaseLayer interface {
	LayeredStoreSupplier
	Store
}

type LayeredStore struct {
	TmpContext      context.Context
	RoleStore       RoleStore
	SchemeStore     SchemeStore
	DatabaseLayer   LayeredStoreDatabaseLayer
	LocalCacheLayer *LocalCacheSupplier
	RedisLayer      *RedisSupplier
	LayerChainHead  LayeredStoreSupplier
}

func NewLayeredStore(db LayeredStoreDatabaseLayer, metrics einterfaces.MetricsInterface, cluster einterfaces.ClusterInterface) Store {
	store := &LayeredStore{
		TmpContext:      context.TODO(),
		DatabaseLayer:   db,
		LocalCacheLayer: NewLocalCacheSupplier(metrics, cluster),
	}

	store.RoleStore = &LayeredRoleStore{store}
	store.SchemeStore = &LayeredSchemeStore{store}

	// Setup the chain
	if ENABLE_EXPERIMENTAL_REDIS {
		mlog.Debug("Experimental redis enabled.")
		store.RedisLayer = NewRedisSupplier()
		store.RedisLayer.SetChainNext(store.DatabaseLayer)
		store.LayerChainHead = store.RedisLayer
	} else {
		store.LocalCacheLayer.SetChainNext(store.DatabaseLayer)
		store.LayerChainHead = store.LocalCacheLayer
	}

	return store
}

type QueryFunction func(LayeredStoreSupplier) *LayeredStoreSupplierResult

func (s *LayeredStore) GetCurrentSchemaVersion() string {
	return s.DatabaseLayer.GetCurrentSchemaVersion()
}

func (s *LayeredStore) Team() TeamStore {
	return s.DatabaseLayer.Team()
}

func (s *LayeredStore) Channel() ChannelStore {
	return s.DatabaseLayer.Channel()
}

func (s *LayeredStore) Post() PostStore {
	return s.DatabaseLayer.Post()
}

func (s *LayeredStore) User() UserStore {
	return s.DatabaseLayer.User()
}

func (s *LayeredStore) Bot() BotStore {
	return s.DatabaseLayer.Bot()
}

func (s *LayeredStore) Audit() AuditStore {
	return s.DatabaseLayer.Audit()
}

func (s *LayeredStore) ClusterDiscovery() ClusterDiscoveryStore {
	return s.DatabaseLayer.ClusterDiscovery()
}

func (s *LayeredStore) Compliance() ComplianceStore {
	return s.DatabaseLayer.Compliance()
}

func (s *LayeredStore) Session() SessionStore {
	return s.DatabaseLayer.Session()
}

func (s *LayeredStore) OAuth() OAuthStore {
	return s.DatabaseLayer.OAuth()
}

func (s *LayeredStore) System() SystemStore {
	return s.DatabaseLayer.System()
}

func (s *LayeredStore) Webhook() WebhookStore {
	return s.DatabaseLayer.Webhook()
}

func (s *LayeredStore) Command() CommandStore {
	return s.DatabaseLayer.Command()
}

func (s *LayeredStore) CommandWebhook() CommandWebhookStore {
	return s.DatabaseLayer.CommandWebhook()
}

func (s *LayeredStore) Preference() PreferenceStore {
	return s.DatabaseLayer.Preference()
}

func (s *LayeredStore) License() LicenseStore {
	return s.DatabaseLayer.License()
}

func (s *LayeredStore) Token() TokenStore {
	return s.DatabaseLayer.Token()
}

func (s *LayeredStore) Emoji() EmojiStore {
	return s.DatabaseLayer.Emoji()
}

func (s *LayeredStore) Status() StatusStore {
	return s.DatabaseLayer.Status()
}

func (s *LayeredStore) FileInfo() FileInfoStore {
	return s.DatabaseLayer.FileInfo()
}

func (s *LayeredStore) Reaction() ReactionStore {
	return s.DatabaseLayer.Reaction()
}

func (s *LayeredStore) Job() JobStore {
	return s.DatabaseLayer.Job()
}

func (s *LayeredStore) UserAccessToken() UserAccessTokenStore {
	return s.DatabaseLayer.UserAccessToken()
}

func (s *LayeredStore) ChannelMemberHistory() ChannelMemberHistoryStore {
	return s.DatabaseLayer.ChannelMemberHistory()
}

func (s *LayeredStore) Plugin() PluginStore {
	return s.DatabaseLayer.Plugin()
}

func (s *LayeredStore) Role() RoleStore {
	return s.RoleStore
}

func (s *LayeredStore) TermsOfService() TermsOfServiceStore {
	return s.DatabaseLayer.TermsOfService()
}

func (s *LayeredStore) UserTermsOfService() UserTermsOfServiceStore {
	return s.DatabaseLayer.UserTermsOfService()
}

func (s *LayeredStore) Scheme() SchemeStore {
	return s.SchemeStore
}

func (s *LayeredStore) Group() GroupStore {
	return s.DatabaseLayer.Group()
}

func (s *LayeredStore) LinkMetadata() LinkMetadataStore {
	return s.DatabaseLayer.LinkMetadata()
}

func (s *LayeredStore) MarkSystemRanUnitTests() {
	s.DatabaseLayer.MarkSystemRanUnitTests()
}

func (s *LayeredStore) Close() {
	s.DatabaseLayer.Close()
}

func (s *LayeredStore) LockToMaster() {
	s.DatabaseLayer.LockToMaster()
}

func (s *LayeredStore) UnlockFromMaster() {
	s.DatabaseLayer.UnlockFromMaster()
}

func (s *LayeredStore) DropAllTables() {
	defer s.LocalCacheLayer.Invalidate()
	s.DatabaseLayer.DropAllTables()
}

func (s *LayeredStore) TotalMasterDbConnections() int {
	return s.DatabaseLayer.TotalMasterDbConnections()
}

func (s *LayeredStore) TotalReadDbConnections() int {
	return s.DatabaseLayer.TotalReadDbConnections()
}

func (s *LayeredStore) TotalSearchDbConnections() int {
	return s.DatabaseLayer.TotalSearchDbConnections()
}

type LayeredRoleStore struct {
	*LayeredStore
}

func (s *LayeredRoleStore) Save(role *model.Role) (*model.Role, *model.AppError) {
	return s.LayerChainHead.RoleSave(s.TmpContext, role)
}

func (s *LayeredRoleStore) Get(roleId string) (*model.Role, *model.AppError) {
	return s.LayerChainHead.RoleGet(s.TmpContext, roleId)
}

func (s *LayeredRoleStore) GetAll() ([]*model.Role, *model.AppError) {
	return s.LayerChainHead.RoleGetAll(s.TmpContext)
}

func (s *LayeredRoleStore) GetByName(name string) (*model.Role, *model.AppError) {
	return s.LayerChainHead.RoleGetByName(s.TmpContext, name)
}

func (s *LayeredRoleStore) GetByNames(names []string) ([]*model.Role, *model.AppError) {
	return s.LayerChainHead.RoleGetByNames(s.TmpContext, names)
}

func (s *LayeredRoleStore) Delete(roldId string) (*model.Role, *model.AppError) {
	return s.LayerChainHead.RoleDelete(s.TmpContext, roldId)
}

func (s *LayeredRoleStore) PermanentDeleteAll() *model.AppError {
	return s.LayerChainHead.RolePermanentDeleteAll(s.TmpContext)
}

type LayeredSchemeStore struct {
	*LayeredStore
}

func (s *LayeredSchemeStore) Save(scheme *model.Scheme) (*model.Scheme, *model.AppError) {
	return s.LayerChainHead.SchemeSave(s.TmpContext, scheme)
}

func (s *LayeredSchemeStore) Get(schemeId string) (*model.Scheme, *model.AppError) {
	return s.LayerChainHead.SchemeGet(s.TmpContext, schemeId)
}

func (s *LayeredSchemeStore) GetByName(schemeName string) (*model.Scheme, *model.AppError) {
	return s.LayerChainHead.SchemeGetByName(s.TmpContext, schemeName)
}

func (s *LayeredSchemeStore) Delete(schemeId string) (*model.Scheme, *model.AppError) {
	return s.LayerChainHead.SchemeDelete(s.TmpContext, schemeId)
}

func (s *LayeredSchemeStore) GetAllPage(scope string, offset int, limit int) ([]*model.Scheme, *model.AppError) {
	return s.LayerChainHead.SchemeGetAllPage(s.TmpContext, scope, offset, limit)
}

func (s *LayeredSchemeStore) PermanentDeleteAll() *model.AppError {
	return s.LayerChainHead.SchemePermanentDeleteAll(s.TmpContext)
}
