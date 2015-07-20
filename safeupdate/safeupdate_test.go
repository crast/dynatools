package safeupdate

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/underarmour/dynago"
)

func TestVariableGen(t *testing.T) {
	v := variableGen{[]rune{'#'}}
	seq := []vgtest{
		vgTest(0, "#a", "#b", "#c", "#d", "#e", "#f"),
		vgTest(19, "#z", "#aa", "#ab", "#ac"),
		vgTest(21, "#ay", "#az", "#ba", "#bb", "#bc"),
		vgTest(21, "#by", "#bz", "#ca", "#cb"),
		vgTest(621, "#zz", "#aaa", "#aab"),
	}
	i := 0
	for _, inner := range seq {
		// skip some so we don't need a complete sequence
		for s := 0; s < inner.skip; s++ {
			v.next()
			i++
		}
		for _, expected := range inner.items {
			i++
			value := v.next()
			if value != expected {
				t.Errorf("Iteration %d: Expected %s, got %s", i, expected, value)
			}
		}
	}
}

func vgTest(skip int, items ...string) vgtest {
	return vgtest{skip, items}
}

type vgtest struct {
	skip  int
	items []string
}

func TestPrepareBuild(t *testing.T) {
	Convey("With a known document", t, func() {
		doc := dynago.Document{
			"Id": "45", "Range": 78, // Key props should be excluded
			"FirstName": "Bob", "Age": 72, "Something": 88, // Non-reserved props get no transform
			"Count": "78", "Domain": "foo.com", // reserved props should be mapped
		}

		Convey("prepare should give known sues", func() {
			s, derivedKey := prepare([2]string{"Id", "Range"}, doc)
			So(s.key, ShouldResemble, [2]string{"Id", "Range"})
			So(len(s.mapped), ShouldEqual, 5)
			So(s.mapped[0], ShouldResemble, attributeDesc{"Age", "", ":a"})
			So(s.mapped[1], ShouldResemble, attributeDesc{"Count", "#b", ":b"})
			So(s.mapped[2], ShouldResemble, attributeDesc{"Domain", "#c", ":c"})
			So(s.mapped[3], ShouldResemble, attributeDesc{"FirstName", "", ":d"})

			So(len(s.eaNames), ShouldEqual, 2)
			So(s.eaNames[0], ShouldResemble, dynago.Param{"#b", "Count"})
			So(derivedKey, ShouldResemble, dynago.Document{"Id": "45", "Range": 78})

		})

		Convey("Build should map the fields to parameters", func() {
			update := Build([]string{"Id", "Range"}, doc)
			So(update.Expression, ShouldEqual, "SET Age=:a, #b=:b, #c=:c, FirstName=:d, Something=:e")
			So(update.Key, ShouldResemble, dynago.Document{"Id": "45", "Range": 78})
			So(len(update.EAValues), ShouldEqual, 5)
			So(update.EAValues[4], ShouldResemble, dynago.Param{":e", 88})
		})
	})

	Convey("With a document containing empty values", t, func() {
		doc := dynago.Document{
			"Id":          "45",
			"Foo":         "Hello",
			"EmptyString": "",
			"Domain":      dynago.StringSet{},
			"Index":       1,
		}
		Convey("prepare should have delete expressions", func() {
			s, derivedKey := prepare([2]string{"Id", ""}, doc)
			So(derivedKey, ShouldResemble, dynago.Document{"Id": "45"})
			So(len(s.mapped), ShouldEqual, 2)
			So(s.mapped[1], ShouldResemble, attributeDesc{"Index", "#b", ":b"})
			So(s.expression, ShouldEqual, "SET Foo=:a, #b=:b REMOVE #c, EmptyString")
		})
		Convey("Build should map only relevant stuff to parameters", func() {
			update := Build([]string{"Id"}, doc)
			So(update.EANames, ShouldResemble, Params{{"#b", "Index"}, {"#c", "Domain"}})
			So(update.EAValues, ShouldResemble, Params{{":a", "Hello"}, {":b", 1}})
		})
	})
}
