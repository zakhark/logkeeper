package logkeeper

import (
	"fmt"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const logKeeperDB = "logkeeper"

type Test struct {
	Id        bson.ObjectId          `bson:"_id"`
	BuildId   interface{}            `bson:"build_id"`
	BuildName string                 `bson:"build_name"`
	Name      string                 `bson:"name"`
	Command   string                 `bson:"command"`
	Started   time.Time              `bson:"started"`
	Ended     *time.Time             `bson:"ended"`
	Info      map[string]interface{} `bson:"info"`
	Failed    bool                   `bson:"failed"`
	Phase     string                 `bson:"phase"`
	Seq       int                    `bson:"seq",omitempty`
}

type LogKeeperBuild struct {
	Id       interface{}            `bson:"_id"`
	Builder  string                 `bson:"builder"`
	BuildNum int                    `bson:"buildnum"`
	Started  time.Time              `bson:"started"`
	Name     string                 `bson:"name"`
	Info     map[string]interface{} `bson:"info"`
	Phases   []string               `bson:"phases"`
	Seq      int                    `bson:"seq",omitempty`
}

// If "raw" is a bson.ObjectId, returns the string value of its .Hex() function.
// Otherwise, returns it's string representation if it implements Stringer, or
// string representation generated by fmt's %v formatter.
func stringifyId(raw interface{}) string {
	if buildObjId, ok := raw.(bson.ObjectId); ok {
		return buildObjId.Hex()
	}
	if asStr, ok := raw.(fmt.Stringer); ok {
		return asStr.String()
	}
	return fmt.Sprintf("%v", raw)
}

func idFromString(raw string) interface{} {
	if bson.IsObjectIdHex(raw) {
		return bson.ObjectIdHex(raw)
	}
	return raw
}

func findTest(db *mgo.Database, id string) (*Test, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, nil
	}
	test := &Test{}

	err := db.C("tests").Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(test)
	if err == mgo.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return test, nil
}

func findTestsForBuild(db *mgo.Database, buildId string) ([]Test, error) {
	queryBuildId := idFromString(buildId)
	tests := []Test{}

	err := db.C("tests").Find(bson.M{"build_id": queryBuildId}).Sort("started").All(&tests)
	if err != nil {
		return nil, err
	}
	return tests, nil
}

func findBuildById(db *mgo.Database, id string) (*LogKeeperBuild, error) {
	queryBuildId := idFromString(id)
	build := &LogKeeperBuild{}

	err := db.C("builds").Find(bson.M{"_id": queryBuildId}).One(build)
	if err == mgo.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return build, nil
}

func findBuildByBuilder(db *mgo.Database, builder string, buildnum int) (*LogKeeperBuild, error) {
	build := &LogKeeperBuild{}

	err := db.C("builds").Find(bson.M{"builder": builder, "buildnum": buildnum}).One(build)
	if err == mgo.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return build, nil
}
