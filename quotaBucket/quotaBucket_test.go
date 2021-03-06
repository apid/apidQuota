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

package quotaBucket_test


import (
	"github.com/apid/apidQuota/constants"
	. "github.com/apid/apidQuota/quotaBucket"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strings"
	"time"
)

var _ = Describe("Test AcceptedQuotaTimeUnitTypes", func() {
	It("testTimeUnit", func() {
		if !IsValidTimeUnit("second") {
			Fail("Expected true: second is a valid TimeUnit")
		}
		if !IsValidTimeUnit("minute") {
			Fail("Expected true: minute is a valid TimeUnit")
		}
		if !IsValidTimeUnit("hour") {
			Fail("Expected true: hour is a valid TimeUnit")
		}
		if !IsValidTimeUnit("day") {
			Fail("Expected true: day is a valid TimeUnit")
		}
		if !IsValidTimeUnit("week") {
			Fail("Expected true: week is a valid TimeUnit")
		}
		if !IsValidTimeUnit("month") {
			Fail("Expected true: month is a valid TimeUnit")
		}

		//invalid type
		if IsValidTimeUnit("invalidType") {
			Fail("Expected false: invalidType is not a valid TimeUnit")
		}
	})
})


var _ = Describe("Test AcceptedQuotaTypes", func() {
	It("testTimeUnit", func() {
		if !IsValidType("calendar") {
			Fail("Expected true: calendar is a valid quotaType")
		}
		if !IsValidType("rollingwindow") {
			Fail("Expected true: rollingwindow is a valid quotaType")
		}
		if IsValidType("invalidType") {
			Fail("Expected false: invalidType is not a valid quotaType")
		}
	})
})

//Tests for QuotaBucket
var _ = Describe("QuotaBucket", func() {

	//validate all fields set as expected.
	//validate period set as expected.
	//validate async QuotaBucket is empty.
	It("Create with NewQuotaBucket with all valid fields", func() {
		edgeOrgID := "sampleOrg"
		id := "sampleID"
		interval := 1
		timeUnit := "hour"
		quotaType := "calendar"
		preciseAtSecondsLevel := true
		maxCount := int64(10)
		weight := int64(1)
		distributed := true
		synchronous := true
		syncTimeInSec := int64(-1)
		syncMessageCount := int64(-1)

		//start time before now()
		startTime := time.Now().UTC().AddDate(0, -1, 0).Unix()

		quotaBucket, err := NewQuotaBucket(edgeOrgID, id, interval, timeUnit,
			quotaType, preciseAtSecondsLevel, startTime, maxCount,
			weight, distributed, synchronous, syncTimeInSec, syncMessageCount)
		Expect(err).NotTo(HaveOccurred())
		now := time.Now().UTC()
		currentHour := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, time.UTC)
		err = quotaBucket.Validate()
		Expect(err).NotTo(HaveOccurred())

		//also check if all the fields are set as expected
		getPeriod, err := quotaBucket.GetPeriod()
		Expect(err).NotTo(HaveOccurred())
		Expect(getPeriod.GetPeriodInputStartTime().String()).ShouldNot(BeEmpty())
		Expect(getPeriod.GetPeriodStartTime().String()).ShouldNot(BeEmpty())
		Expect(getPeriod.GetPeriodEndTime().String()).ShouldNot(BeEmpty())
		Expect(getPeriod.GetPeriodStartTime().String()).Should(Equal(currentHour.String()))
		Expect(getPeriod.GetPeriodEndTime().String()).Should(Equal(currentHour.Add(time.Hour).String()))

		//check if isDistributed and isSynchronous are true.
		if !quotaBucket.IsDistrubuted(){
			Fail("expected true, returned true.")
		}
		if !quotaBucket.IsSynchronous(){
			Fail("expected true, returned true.")
		}

		asyncBucket := quotaBucket.GetAsyncQuotaBucket()
		if asyncBucket != nil{
			Fail("asyncBucket should be nil for synchronous request.")
		}

		//testCase2: synchronous := true  and (syncTimeInSec and syncMessageCount) > -1 -> aSyncQuotaBucket is nil.
		syncTimeInSec = int64(10)
		syncMessageCount = int64(10)
		quotaBucket, err = NewQuotaBucket(edgeOrgID, id, interval, timeUnit,
			quotaType, preciseAtSecondsLevel, startTime, maxCount,
			weight, distributed, synchronous, syncTimeInSec, syncMessageCount)
		Expect(err).NotTo(HaveOccurred())
		now = time.Now().UTC()
		currentHour = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, time.UTC)
		err = quotaBucket.Validate()
		Expect(err).NotTo(HaveOccurred())
		//check if isDistributed and isSynchronous are true.
		if !quotaBucket.IsDistrubuted(){
			Fail("expected true, returned true.")
		}
		if !quotaBucket.IsSynchronous(){
			Fail("expected true, returned true.")
		}

		asyncBucket = quotaBucket.GetAsyncQuotaBucket()
		if asyncBucket != nil{
			Fail("asyncBucket should be nil for synchronous request.")
		}


	})

	//startTime for quotaBucket after time.Now()
	It("Create with NewQuotaBucket with start time after now()", func() {
		edgeOrgID := "sampleOrg"
		id := "sampleID"
		interval := 1
		timeUnit := "hour"
		quotaType := "calendar"
		preciseAtSecondsLevel := true
		maxCount := int64(10)
		weight := int64(1)
		distributed := true
		synchronous := true
		syncTimeInSec := int64(-1)
		syncMessageCount := int64(-1)
		//start time is after now() -> should still set period.
		startTime := time.Now().UTC().AddDate(0, 1, 0).Unix()
		quotaBucket, err := NewQuotaBucket(edgeOrgID, id, interval, timeUnit,
			quotaType, preciseAtSecondsLevel, startTime, maxCount,
			weight, distributed, synchronous, syncTimeInSec, syncMessageCount)
		Expect(err).NotTo(HaveOccurred())
		now := time.Now().UTC()
		currentHour := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, time.UTC)
		err = quotaBucket.Validate()
		Expect(err).NotTo(HaveOccurred())

		getPeriod, err := quotaBucket.GetPeriod()
		Expect(err).NotTo(HaveOccurred())
		Expect(getPeriod.GetPeriodInputStartTime().String()).ShouldNot(BeEmpty())
		Expect(getPeriod.GetPeriodStartTime().String()).ShouldNot(BeEmpty())
		Expect(getPeriod.GetPeriodEndTime().String()).ShouldNot(BeEmpty())
		Expect(getPeriod.GetPeriodStartTime().String()).Should(Equal(currentHour.String()))
		Expect(getPeriod.GetPeriodEndTime().String()).Should(Equal(currentHour.Add(time.Hour).String()))

		//check if isDistributed and isSynchronous are true.
		if !quotaBucket.IsDistrubuted(){
			Fail("expected true, returned false.")
		}
		if !quotaBucket.IsSynchronous(){
			Fail("expected true, returned false.")
		}

		asyncBucket := quotaBucket.GetAsyncQuotaBucket()
		if asyncBucket != nil{
			Fail("asyncBucket should be nil for synchronous request.")
		}
	})

	//Testcase1 : with syncTimeInSec
	//Testcase2 : with syncMessageCount
	//Testcase3 : InvalidTestCase - with syncTimeInSec and syncMessageCount
	It("Create with NewQuotaBucket for aSyncQuotBucket", func() {
		edgeOrgID := "sampleOrg"
		id := "sampleID"
		interval := 1
		timeUnit := "hour"
		quotaType := "calendar"
		preciseAtSecondsLevel := true
		maxCount := int64(10)
		weight := int64(1)
		distributed := true
		synchronous := false

		//Testcase1 : with syncTimeInSec
		syncTimeInSec := int64(10)
		syncMessageCount := int64(-1)
		//start time is after now() -> should still set period.
		startTime := time.Now().UTC().AddDate(0, 1, 0).Unix()
		quotaBucket, err := NewQuotaBucket(edgeOrgID, id, interval, timeUnit,
			quotaType, preciseAtSecondsLevel, startTime, maxCount,
			weight, distributed, synchronous, syncTimeInSec, syncMessageCount)
		Expect(err).NotTo(HaveOccurred())
		now := time.Now().UTC()
		currentHour := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, time.UTC)
		err = quotaBucket.Validate()
		Expect(err).NotTo(HaveOccurred())

		getPeriod, err := quotaBucket.GetPeriod()
		Expect(err).NotTo(HaveOccurred())
		Expect(getPeriod.GetPeriodInputStartTime().String()).ShouldNot(BeEmpty())
		Expect(getPeriod.GetPeriodStartTime().String()).ShouldNot(BeEmpty())
		Expect(getPeriod.GetPeriodEndTime().String()).ShouldNot(BeEmpty())
		Expect(getPeriod.GetPeriodStartTime().String()).Should(Equal(currentHour.String()))
		Expect(getPeriod.GetPeriodEndTime().String()).Should(Equal(currentHour.Add(time.Hour).String()))

		//check if isDistributed is true and isSynchronous is false.
		if !quotaBucket.IsDistrubuted(){
			Fail("expected true, returned false.")
		}
		if quotaBucket.IsSynchronous(){
			Fail("expected false, returned true.")
		}

		asyncBucket := quotaBucket.GetAsyncQuotaBucket()
		if asyncBucket == nil{
			Fail("asyncBucket should be nil for synchronous request.")
		}

		//Testcase2 : with syncMessageCount
		syncTimeInSec = int64(-1)
		syncMessageCount = int64(10)
		//start time is after now() -> should still set period.
		startTime = time.Now().UTC().AddDate(0, 1, 0).Unix()
		quotaBucket, err = NewQuotaBucket(edgeOrgID, id, interval, timeUnit,
			quotaType, preciseAtSecondsLevel, startTime, maxCount,
			weight, distributed, synchronous, syncTimeInSec, syncMessageCount)
		Expect(err).NotTo(HaveOccurred())
		now = time.Now().UTC()
		currentHour = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, time.UTC)
		err = quotaBucket.Validate()
		Expect(err).NotTo(HaveOccurred())

		getPeriod, err = quotaBucket.GetPeriod()
		Expect(err).NotTo(HaveOccurred())
		Expect(getPeriod.GetPeriodInputStartTime().String()).ShouldNot(BeEmpty())
		Expect(getPeriod.GetPeriodStartTime().String()).ShouldNot(BeEmpty())
		Expect(getPeriod.GetPeriodEndTime().String()).ShouldNot(BeEmpty())
		Expect(getPeriod.GetPeriodStartTime().String()).Should(Equal(currentHour.String()))
		Expect(getPeriod.GetPeriodEndTime().String()).Should(Equal(currentHour.Add(time.Hour).String()))

		//check if isDistributed is true and isSynchronous is false.
		if !quotaBucket.IsDistrubuted(){
			Fail("expected true, returned false.")
		}
		if quotaBucket.IsSynchronous(){
			Fail("expected false, returned true.")
		}

		asyncBucket = quotaBucket.GetAsyncQuotaBucket()
		if asyncBucket == nil{
			Fail("asyncBucket should be nil for synchronous request.")
		}

		//Testcase3 : InvalidTestCase - with syncTimeInSec and syncMessageCount
		syncTimeInSec = int64(10)
		syncMessageCount = int64(10)
		//start time is after now() -> should still set period.
		startTime = time.Now().UTC().AddDate(0, 1, 0).Unix()
		quotaBucket, err = NewQuotaBucket(edgeOrgID, id, interval, timeUnit,
			quotaType, preciseAtSecondsLevel, startTime, maxCount,
			weight, distributed, synchronous, syncTimeInSec, syncMessageCount)
		Expect(err).To(HaveOccurred())

	})

	It("Create with NonDistributed QuotaBucket", func() {
		edgeOrgID := "sampleOrg"
		id := "sampleID"
		interval := 1
		timeUnit := "hour"
		quotaType := "calendar"
		preciseAtSecondsLevel := true
		maxCount := int64(10)
		weight := int64(1)
		distributed := false
		synchronous := false
		syncTimeInSec := int64(-1)
		syncMessageCount := int64(-1)

		//start time before now()
		startTime := time.Now().UTC().AddDate(0, -1, 0).Unix()

		quotaBucket, err := NewQuotaBucket(edgeOrgID, id, interval, timeUnit,
			quotaType, preciseAtSecondsLevel, startTime, maxCount,
			weight, distributed, synchronous, syncTimeInSec, syncMessageCount)
		Expect(err).NotTo(HaveOccurred())
		now := time.Now().UTC()
		currentHour := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, time.UTC)
		err = quotaBucket.Validate()
		Expect(err).NotTo(HaveOccurred())
		Expect(quotaBucket.IsSynchronous()).Should(BeFalse())

		//also check if all the fields are set as expected
		getPeriod, err := quotaBucket.GetPeriod()
		Expect(err).NotTo(HaveOccurred())
		Expect(getPeriod.GetPeriodInputStartTime().String()).ShouldNot(BeEmpty())
		Expect(getPeriod.GetPeriodStartTime().String()).ShouldNot(BeEmpty())
		Expect(getPeriod.GetPeriodEndTime().String()).ShouldNot(BeEmpty())
		Expect(getPeriod.GetPeriodStartTime().String()).Should(Equal(currentHour.String()))
		Expect(getPeriod.GetPeriodEndTime().String()).Should(Equal(currentHour.Add(time.Hour).String()))

		//check if isDistributed is false and isSynchronous is false.
		if quotaBucket.IsDistrubuted(){
			Fail("expected false, returned true.")
		}
		if quotaBucket.IsSynchronous(){
			Fail("expected false, returned true.")
		}

		asyncBucket := quotaBucket.GetAsyncQuotaBucket()
		if asyncBucket != nil {
			Fail("asyncBucket should be nil for synchronous request.")
		}

		synchronous = true
		quotaBucket, err = NewQuotaBucket(edgeOrgID, id, interval, timeUnit,
			quotaType, preciseAtSecondsLevel, startTime, maxCount,
			weight, distributed, synchronous, syncTimeInSec, syncMessageCount)
		Expect(err).To(HaveOccurred())
	})



	It("Test invalid timeUnitType", func() {
		edgeOrgID := "sampleOrg"
		id := "sampleID"
		timeUnit := "invalidTimeUnit"
		quotaType := "calendar"
		interval := 1
		maxCount := int64(10)
		weight := int64(1)
		preciseAtSecondsLevel := true
		startTime := time.Now().UTC().AddDate(0, -1, 0).Unix()
		distributed := true
		synchronous := true
		syncTimeInSec := int64(-1)
		syncMessageCount := int64(-1)

		quotaBucket, err := NewQuotaBucket(edgeOrgID, id, interval, timeUnit,
			quotaType, preciseAtSecondsLevel, startTime, maxCount,
			weight, distributed, synchronous, syncTimeInSec, syncMessageCount)
		err = quotaBucket.Validate()
		Expect(err).To(HaveOccurred())

		if !strings.Contains(err.Error(), constants.InvalidQuotaTimeUnitType) {
			Fail("expected: " + constants.InvalidQuotaTimeUnitType + "but got: " + err.Error())
		}

	})

	It("Test invalid quotaType", func() {
		edgeOrgID := "sampleOrg"
		id := "sampleID"
		timeUnit := "hour"
		quotaType := "invalidTimeUnit"
		interval := 1
		maxCount := int64(10)
		weight := int64(1)
		preciseAtSecondsLevel := true
		startTime := time.Now().UTC().AddDate(0, -1, 0).Unix()
		distributed := true
		synchronous := true
		syncTimeInSec := int64(-1)
		syncMessageCount := int64(-1)

		quotaBucket, err := NewQuotaBucket(edgeOrgID, id, interval, timeUnit,
			quotaType, preciseAtSecondsLevel, startTime, maxCount,
			weight, distributed, synchronous, syncTimeInSec, syncMessageCount)
		err = quotaBucket.Validate()
		Expect(err).To(HaveOccurred())

		if !strings.Contains(err.Error(), constants.InvalidQuotaType) {
			Fail("expected: " + constants.InvalidQuotaType + "but got: " + err.Error())
		}
	})
})

var _ = Describe("IsCurrentPeriod", func() {
	It("Test IsCurrentPeriod for RollingType Window  - Valid TestCase", func() {

		edgeOrgID := "sampleOrg"
		id := "sampleID"
		timeUnit := "hour"
		quotaType := "rollingwindow"
		interval := 1
		maxCount := int64(10)
		weight := int64(1)
		preciseAtSecondsLevel := true
		distributed := true
		synchronous := true
		syncTimeInSec := int64(-1)
		syncMessageCount := int64(-1)

		//InputStart time is before now
		startTime := time.Now().UTC().AddDate(0, -1, 0).Unix()

		quotaBucket, err := NewQuotaBucket(edgeOrgID, id, interval, timeUnit,
			quotaType, preciseAtSecondsLevel, startTime, maxCount,
			weight, distributed, synchronous, syncTimeInSec, syncMessageCount)
		Expect(err).NotTo(HaveOccurred())
		err = quotaBucket.Validate()
		Expect(err).NotTo(HaveOccurred())

		period, err := quotaBucket.GetPeriod()
		if err != nil {
			Fail("no error expected")
		}
		if ok := period.IsCurrentPeriod(quotaBucket); !ok {
			Fail("Exprected true, returned: false")
		}

		//InputStart time is now
		startTime = time.Now().UTC().Unix()
		quotaBucket, err = NewQuotaBucket(edgeOrgID, id, interval, timeUnit,
			quotaType, preciseAtSecondsLevel, startTime, maxCount,
			weight, distributed, synchronous, syncTimeInSec, syncMessageCount)
		Expect(err).NotTo(HaveOccurred())
		err = quotaBucket.Validate()
		Expect(err).NotTo(HaveOccurred())

		period, err = quotaBucket.GetPeriod()
		if err != nil {
			Fail("no error expected")
		}
		period.IsCurrentPeriod(quotaBucket)
		if ok := period.IsCurrentPeriod(quotaBucket); !ok {
			Fail("Exprected true, returned: false")
		}
	})

	It("Test IsCurrentPeriod for RollingType Window - InValid TestCase", func() {

		edgeOrgID := "sampleOrg"
		id := "sampleID"
		timeUnit := "hour"
		quotaType := "rollingwindow"
		interval := 1
		maxCount := int64(10)
		weight := int64(1)
		preciseAtSecondsLevel := true
		//InputStart time is after now.
		startTime := time.Now().UTC().AddDate(0, 1, 0).Unix()
		distributed := true
		synchronous := true
		syncTimeInSec := int64(-1)
		syncMessageCount := int64(-1)

		quotaBucket, err := NewQuotaBucket(edgeOrgID, id, interval, timeUnit,
			quotaType, preciseAtSecondsLevel, startTime, maxCount,
			weight, distributed, synchronous, syncTimeInSec, syncMessageCount)
		Expect(err).NotTo(HaveOccurred())

		err = quotaBucket.Validate()
		Expect(err).NotTo(HaveOccurred())

		period, err := quotaBucket.GetPeriod()
		Expect(err).NotTo(HaveOccurred())
		if ok := period.IsCurrentPeriod(quotaBucket); ok {
			Fail("Exprected false, returned: true")
		}

	})

	It("Test IsCurrentPeriod for calendarType Window - Valid TestCases", func() {

		edgeOrgID := "sampleOrg"
		id := "sampleID"
		timeUnit := "hour"
		quotaType := "calendar"
		interval := 1
		maxCount := int64(10)
		weight := int64(1)
		preciseAtSecondsLevel := true
		distributed := true
		synchronous := true
		syncTimeInSec := int64(-1)
		syncMessageCount := int64(-1)

		//InputStart time is before now
		startTime := time.Now().UTC().UTC().AddDate(-1, -1, 0).Unix()

		quotaBucket, err := NewQuotaBucket(edgeOrgID, id, interval, timeUnit,
			quotaType, preciseAtSecondsLevel, startTime, maxCount,
			weight, distributed, synchronous, syncTimeInSec, syncMessageCount)
		Expect(err).NotTo(HaveOccurred())

		err = quotaBucket.Validate()
		Expect(err).NotTo(HaveOccurred())

		period, err := quotaBucket.GetPeriod()
		if err != nil {
			Fail("no error expected but returned " + err.Error())
		}
		if ok := period.IsCurrentPeriod(quotaBucket); !ok {
			Fail("Exprected true, returned: false")
		}

		//InputStart time is now
		startTime = time.Now().UTC().Unix()
		quotaBucket, err = NewQuotaBucket(edgeOrgID, id, interval, timeUnit,
			quotaType, preciseAtSecondsLevel, startTime, maxCount,
			weight, distributed, synchronous, syncTimeInSec, syncMessageCount)
		Expect(err).NotTo(HaveOccurred())

		err = quotaBucket.Validate()
		Expect(err).NotTo(HaveOccurred())

		period, err = quotaBucket.GetPeriod()
		if err != nil {
			Fail("no error expected but returned " + err.Error())
		}
		period.IsCurrentPeriod(quotaBucket)
		if ok := period.IsCurrentPeriod(quotaBucket); !ok {
			Fail("Exprected true, returned: false")
		}

	})

	It("Test IsCurrentPeriod for calendarType Window InValid TestCase", func() {

		edgeOrgID := "sampleOrg"
		id := "sampleID"
		timeUnit := "hour"
		quotaType := "calendar"
		interval := 1
		maxCount := int64(10)
		weight := int64(1)
		preciseAtSecondsLevel := true
		distributed := true
		synchronous := true
		syncTimeInSec := int64(-1)
		syncMessageCount := int64(-1)

		//InputStart time is after now.
		startTime := time.Now().UTC().AddDate(0, 1, 0).Unix()

		quotaBucket, err := NewQuotaBucket(edgeOrgID, id, interval, timeUnit,
			quotaType, preciseAtSecondsLevel, startTime, maxCount,
			weight, distributed, synchronous, syncTimeInSec, syncMessageCount)
		Expect(err).NotTo(HaveOccurred())

		err = quotaBucket.Validate()
		Expect(err).NotTo(HaveOccurred())

		period, err := quotaBucket.GetPeriod()
		Expect(err).NotTo(HaveOccurred())

		if ok := period.IsCurrentPeriod(quotaBucket); ok {
			Fail("Exprected false, returned: true")
		}
	})
})

var _ = Describe("Test GetPeriod and the timeInterval in period set as expected", func() {
	It("Valid GetPeriod", func() {
		edgeOrgID := "sampleOrg"
		id := "sampleID"
		timeUnit := "hour"
		quotaType := "rollingwindow"
		interval := 1
		maxCount := int64(10)
		weight := int64(1)
		preciseAtSecondsLevel := true
		distributed := true
		synchronous := true
		syncTimeInSec := int64(-1)
		syncMessageCount := int64(-1)

		//InputStart time is before now
		startTime := time.Now().UTC().AddDate(0, -1, 0).Unix()
		quotaBucket, err := NewQuotaBucket(edgeOrgID, id, interval, timeUnit,
			quotaType, preciseAtSecondsLevel, startTime, maxCount,
			weight, distributed, synchronous, syncTimeInSec, syncMessageCount)
		Expect(err).NotTo(HaveOccurred())

		err = quotaBucket.Validate()
		Expect(err).NotTo(HaveOccurred())
		qPeriod, err := quotaBucket.GetPeriod()
		Expect(err).NotTo(HaveOccurred())

		// check if the rolling window was set properly
		Expect(qPeriod.GetPeriodInputStartTime()).Should(Equal(quotaBucket.GetStartTime()))
		if !qPeriod.GetPeriodEndTime().After(qPeriod.GetPeriodStartTime()) {
			Fail("Rolling Window was not set as expected")
		}
		intervalDuration := qPeriod.GetPeriodEndTime().Sub(qPeriod.GetPeriodStartTime())
		expectedDuration, err := GetIntervalDurtation(quotaBucket)
		Expect(intervalDuration).Should(Equal(expectedDuration))

		//for non rolling Type window setCurrentPeriod as endTime is < time.now.
		quotaType = "calendar"
		quotaBucket, err = NewQuotaBucket(edgeOrgID, id, interval, timeUnit,
			quotaType, preciseAtSecondsLevel, startTime, maxCount,
			weight, distributed, synchronous, syncTimeInSec, syncMessageCount)
		Expect(err).NotTo(HaveOccurred())

		qPeriod, err = quotaBucket.GetPeriod()
		Expect(err).NotTo(HaveOccurred())
		// check if the calendar window was set properly
		Expect(qPeriod.GetPeriodInputStartTime()).Should(Equal(quotaBucket.GetStartTime()))
		if !qPeriod.GetPeriodEndTime().After(qPeriod.GetPeriodStartTime()) {
			Fail("period for Non Rolling Window Type was not set as expected")
		}
		intervalDuration = qPeriod.GetPeriodEndTime().Sub(qPeriod.GetPeriodStartTime())
		expectedDuration, err = GetIntervalDurtation(quotaBucket)
		Expect(intervalDuration).Should(Equal(expectedDuration))
	})
})
