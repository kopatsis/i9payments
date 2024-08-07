package db

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty"`
	Name              string             `bson:"name"`
	Username          string             `bson:"username"`
	Paying            bool               `bson:"paying"`
	Provider          string             `bson:"provider"`
	Level             float32            `bson:"level"`
	BannedExercises   []string           `bson:"bannedExer"`
	BannedStretches   []string           `bson:"bannedStr"`
	BannedParts       []int              `bson:"bannedParts"`
	PlyoTolerance     int                `bson:"plyoToler"`
	ExerFavoriteRates map[string]float32 `bson:"exerfavs"`
	ExerModifications map[string]float32 `bson:"exermods"`
	TypeModifications map[string]float32 `bson:"typemods"`
	RoundEndurance    map[int]float32    `bson:"roundendur"`
	TimeEndurance     map[int]float32    `bson:"timeendur"`
	PushupSetting     string             `bson:"pushsetting"`
	LastMinutes       float32            `bson:"lastmins"`
	LastDifficulty    int                `bson:"lastdiff"`
	Assessed          bool               `bson:"assessed"`
	Badges            []string           `bson:"badges"`       //New
	StrWOStartedCt    int                `bson:"strwostartct"` //New
	StrWOGeneratedCt  int                `bson:"strwogenct"`   //New
	WOGeneratedCt     int                `bson:"wogenct"`      //New
	WOStartedCt       int                `bson:"wostartct"`    //New
	WORatedCt         int                `bson:"completed"`    //New
	DisplayLevel      int                `bson:"displevel"`    //New
	ResetDate         primitive.DateTime `bson:"reset"`        //New
}

type UserPayment struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	UserMongoID    string             `bson:"userid"`
	Username       string             `bson:"username"`
	Provider       string             `bson:"provider"`
	SubscriptionID string             `bson:"subid"`
	SubLength      string             `bson:"length"`
	EndDate        primitive.DateTime `bson:"end"`
	Expires        primitive.DateTime `bson:"expires"`
	Processing     bool               `bson:"processing"`
	Ending         bool               `bson:"ending"`
}
