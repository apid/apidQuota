// Copyright 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package quotaBucket

import (
	"github.com/apid/apidQuota/constants"
	"sync"
	"time"
	"errors"
)

var quotaCachelock = sync.RWMutex{}

type quotaBucketCache struct {
	qBucket    *QuotaBucket
	expiryTime int64
}

var quotaCache map[string]quotaBucketCache

func init() {
	quotaCache = make(map[string]quotaBucketCache)
}

func getFromCache(cacheKey string, weight int64) (*QuotaBucket, bool) {
	quotaCachelock.Lock()
	qBucketCache, ok := quotaCache[cacheKey]
	quotaCachelock.Unlock()

	if !ok {
		return nil, false
	}

	isExpired := time.Unix(qBucketCache.expiryTime, 0).Before(time.Now().UTC())
	if isExpired {

		removeFromCache(cacheKey, qBucketCache)
		return nil, false
	}

	// update expiry time every time you access.
	qBucketCache.qBucket.Weight = weight
	ttl := time.Now().UTC().Add(constants.CacheTTL).Unix()
	qBucketCache.expiryTime = ttl

	quotaCachelock.Lock()
	quotaCache[cacheKey] = qBucketCache
	quotaCachelock.Unlock()

	return qBucketCache.qBucket, true

}

func removeFromCache(cacheKey string, qBucketCache quotaBucketCache) error {
	//for async Stop the scheduler.

	if qBucketCache.qBucket.Distributed && !qBucketCache.qBucket.IsSynchronous() {
		aSyncBucket := qBucketCache.qBucket.GetAsyncQuotaBucket()
		if aSyncBucket == nil {
			return errors.New(constants.AsyncQuotaBucketEmpty + " : aSyncQuotaBucket to increment cannot be empty.")
		}
		qticker, err := aSyncBucket.getAsyncQTicker()
		if err != nil {
			return err
		}
		qticker.Stop()
	}

	quotaCachelock.Lock()
	delete(quotaCache, cacheKey)
	quotaCachelock.Unlock()
	return nil
}

func addToCache(qBucketToAdd *QuotaBucket) {

	cacheKey := qBucketToAdd.GetEdgeOrgID() + constants.CacheKeyDelimiter + qBucketToAdd.GetID()
	ttl := time.Now().UTC().Add(constants.CacheTTL).Unix()
	qCacheData := quotaBucketCache{
		qBucket:    qBucketToAdd,
		expiryTime: ttl,
	}

	quotaCachelock.Lock()
	quotaCache[cacheKey] = qCacheData
	quotaCachelock.Unlock()
}
