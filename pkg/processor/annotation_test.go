package processor

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAnnotation(t *testing.T) {
	Convey("ParseAnnotations", t, func() {
		Convey("should parse annotations with argument", func() {
			So(
				ParseAnnotations([]string{
					"   @api   ",
				}),
				ShouldResemble,
				[]Annotation{
					Annotation{Type: AnnotationTypeAPI, Argument: "", Body: nil},
				},
			)
			So(
				ParseAnnotations([]string{
					"   @API  test arguments ",
				}),
				ShouldResemble,
				[]Annotation{
					Annotation{Type: AnnotationTypeAPI, Argument: "test arguments", Body: nil},
				},
			)
		})

		Convey("should parse annotations with body", func() {
			So(
				ParseAnnotations([]string{
					"   @API  argument ",
					" test  ",
					"   example",
				}),
				ShouldResemble,
				[]Annotation{
					Annotation{Type: AnnotationTypeAPI, Argument: "argument", Body: []string{
						"test",
						"example",
					}},
				},
			)
		})

		Convey("should parse multiple annotations", func() {
			So(
				ParseAnnotations([]string{
					"   @API  argument ",
					"@Version v0.1",
				}),
				ShouldResemble,
				[]Annotation{
					Annotation{Type: AnnotationTypeAPI, Argument: "argument", Body: nil},
					Annotation{Type: AnnotationTypeVersion, Argument: "v0.1", Body: nil},
				},
			)
			So(
				ParseAnnotations([]string{
					"   @API  argument ",
					"test",
					"@tag example",
					"  some example",
				}),
				ShouldResemble,
				[]Annotation{
					Annotation{Type: AnnotationTypeAPI, Argument: "argument", Body: []string{"test"}},
					Annotation{Type: AnnotationTypeTag, Argument: "example", Body: []string{"some example"}},
				},
			)
		})

		Convey("should ignore unknown text", func() {
			So(ParseAnnotations([]string{}), ShouldBeEmpty)
			So(ParseAnnotations([]string{
				"test",
			}), ShouldBeEmpty)
			So(ParseAnnotations([]string{
				"@unknown",
			}), ShouldBeEmpty)
			So(
				ParseAnnotations([]string{
					"  test",
					"   @API  argument ",
				}),
				ShouldResemble,
				[]Annotation{
					Annotation{Type: AnnotationTypeAPI, Argument: "argument", Body: nil},
				},
			)
		})

	})
}
