package main

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/fzzy/radix/redis"

)

func TestOneShouldBeEqualOne(t *testing.T) {
	Convey("1 should equal 1", t, func() {
			So(1, ShouldEqual, 1)
		})
}

func Test_userPasswordShouldBeDifferent_whenSavingToDatastore(t *testing.T) {
	Convey("When I connect to the Datastore", t, func() {
			conn, err := redis.Dial("tcp", "127.0.0.1:6379")
			errorHandler(err)
			Convey("I save a user", func() {
					user := User{"Matheus", "12345"}
					saveUser(user, conn)
					Convey("Then I retrieve the saved user from the Datastore", func() {
							value, err := conn.Cmd("get", "Matheus").Str()
							Convey("There should be no errors", func() {
									So(err, ShouldEqual, nil)
								})
							Convey("The saved value for the user should be different since it is MD5 format", func() {
									So(value, ShouldNotEqual, user.Pass)
								})
						})
				})
		})
}
