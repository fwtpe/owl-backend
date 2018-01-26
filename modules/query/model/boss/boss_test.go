package boss

import (
	rd "github.com/Pallinder/go-randomdata"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("PlatformDetail", func() {
	Context("ShortenDescription()", func() {
		longText := rd.RandStringRunes(201)

		DescribeTable("The result should be as expected one",
			func(source string, expectedResult string) {
				testedResult := (&PlatformDetail{Description: source}).ShortenDescription()

				Expect(testedResult).To(Equal(expectedResult))
			},
			Entry("Simple text", "This is a car", "This is a car"),
			Entry("Simple text(trimmed)", " True    color ", "True color"),
			Entry("Text contains linefeed and tab", "This\tis\ra\ncat", "This is a cat"),
			Entry("Text is more the 200 characters", longText, string([]rune(longText)[0:100])),
		)
	})
})

var _ = Describe("ConvertsPlatformIpsToHosts", func() {
	var getTestedHost = func(hosts []*Host, targetIp string) *Host {
		/**
		 * Finds the tested host
		 */
		var targetHost *Host
		for _, host := range hosts {
			if host.Ip == targetIp {
				targetHost = host
				break
			}
		}
		// :~)

		return targetHost
	}

	Context("Normal conversion", func() {
		It("The host[34.56.0.33] should match certain properties", func() {
			testedHosts := ConvertsPlatformIpsToHosts(
				[]*PlatformIps{
					{
						Name: "k01.c01",
						IpList: []*PlatformIp{
							{
								Ip: "10.7.0.33", Hostname: "ene-ks-101-007-000-033",
								Status: "1", PopId: "91", Type: "KEEP",
							},
							{
								Ip: "10.7.0.34", Hostname: "ene-ks-101-007-000-034",
								Status: "1", PopId: "92", Type: "KEEP",
							},
						},
					},
					{
						Name: "k01.c02",
						IpList: []*PlatformIp{
							{
								Ip: "34.56.0.33", Hostname: "ene-ks-034-056-000-033",
								Status: "1", PopId: "101", Type: "KEEP",
							},
							{
								Ip: "34.56.0.34", Hostname: "ene-ks-034-056-000-034",
								Status: "1", PopId: "102", Type: "KEEP",
							},
							{ // Skipped
								Ip: "39.20.0.34", Hostname: "",
								Status: "1", PopId: "119", Type: "KEEP",
							},
						},
					},
				},
			)

			Expect(testedHosts).To(HaveLen(4))

			testedHost := getTestedHost(testedHosts, "34.56.0.33")
			Expect(testedHost).To(PointTo(MatchAllFields(
				Fields{
					"Ip":        Equal("34.56.0.33"),
					"Hostname":  Equal("ene-ks-034-056-000-033"),
					"Platform":  Equal("k01.c02"),
					"Platforms": And(HaveLen(1), ConsistOf("k01.c02")),
					"Isp":       Equal("ene"),
					"Activate":  Equal("1"),
					"IdcId":     Equal("101"),
				},
			)))
		})
	})

	Context("Duplicated(multiple platforms) conversion", func() {
		It("The host[121.34.0.33](duplicated) should match certain properties", func() {
			testedHosts := ConvertsPlatformIpsToHosts(
				[]*PlatformIps{
					{
						Name: "z01.a01",
						IpList: []*PlatformIp{
							{
								Ip: "121.34.0.33", Hostname: "zza-kc-121-034-000-033",
								Status: "0", PopId: "151", Type: "KEEP",
							},
							{
								Ip: "121.34.0.34", Hostname: "zza-kc-121-034-000-034",
								Status: "1", PopId: "152", Type: "KEEP",
							},
						},
					},
					{
						Name: "z01.a02",
						IpList: []*PlatformIp{
							{
								Ip: "121.34.0.33", Hostname: "zza-kc-121-034-000-033",
								Status: "0", PopId: "151", Type: "KEEP",
							},
						},
					},
					{ // Non-effective platform, but set status to "1"
						Name: "z01.a03",
						IpList: []*PlatformIp{
							{
								Ip: "121.34.198.27", Hostname: "zza-kc-121-034-000-033",
								Status: "1", PopId: "151", Type: "KEEP",
							},
						},
					},
				},
			)

			Expect(testedHosts).To(HaveLen(2))

			testedHost := getTestedHost(testedHosts, "121.34.0.33")
			Expect(testedHost).To(PointTo(MatchFields(
				IgnoreExtras,
				Fields{
					"Ip":        Equal("121.34.0.33"),
					"Platform":  Equal("z01.a02"),
					"Platforms": And(HaveLen(2), ConsistOf([]string{"z01.a01", "z01.a02"})),
					"Activate":  Equal("1"),
				},
			)))
		})
	})
})
