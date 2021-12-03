package profanities

import "github.com/MikkelHJuul/profaneword"

type SwearStringer struct {}

func (_ SwearStringer) String() string {
	return GetSwear()
}

func GetSwear() string {
	randDevice := profaneword.CryptoRand{}
	rand := randDevice.RandMax(int64(len(swear)))
	return swear[rand]}


var swear = [...]string{
	"Good lord",
	"Jesus Christ",
	"arse",
	"arsehole",
	"ball",
	"boob",
	"bugger",
	"cake", //lol
	"choad",
	"cock",
	"crap",
	"cunt",
	"damn",
	"dorky",
	"eat",
	"fat",
	"frick",
	"frothing",
	"fuck",
	"fuckup",
	"fugly",
	"golly",
	"hail",
	"hard",
	"hell",
	"huge",
	"jerkass",
	"jizz",
	"kills",
	"lard",
	"lick",
	"limp",
	"master", // not a swear
	"my lord",
	"piss",
	"poop",
	"prick",
	"puta",
	"puto",
	"rapping",
	"rectal",
	"shit",
	"shite",
	"sod-off",
	"son of a",
	"son-of-a",
	"stupid",
	"suck",
	"tit", // the bird, right?
	"trojan",
	"wank",
	"x-rated",
	"xxx",

}