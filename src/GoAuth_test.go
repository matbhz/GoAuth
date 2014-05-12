package main

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey")

func TestOneShouldBeEqualOne(t *testing.T) {
	Convey("1 should equal 1", t, func() {
			So(1, ShouldEqual, 1)
		})
}
