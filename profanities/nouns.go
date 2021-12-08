package profanities

import "github.com/MikkelHJuul/profaneword"

type NounStringer struct{}

func (_ NounStringer) String() string {
	return GetNoun()
}

func GetNoun() string {
	randDevice := profaneword.CryptoRand{}
	rand := randDevice.RandMax(int64(len(nouns)))
	return nouns[rand]
}

var nouns = [...]string{
	"Sir Wank-a-Lot",
	"addict",
	"adulterer",
	"anal-secretion",
	"anus",
	"armpit",
	"arsehole",
	"bait",
	"bandit",
	"banshee",
	"barf",
	"beater",
	"bellend",
	"boner",
	"breasthugger",
	"butt",
	"cancer",
	"cat-killer",
	"chipmunk",
	"clown",
	"cracker",
	"damned",
	"denier",
	"dickhead",
	"dildo",
	"dimwit",
	"dipshit",
	"dog",
	"dong",
	"donkey",
	"dork",
	"douchebag",
	"dwarf",
	"eel",
	"ejaculate",
	"fanatic",
	"fatty",
	"feet",
	"fetish",
	"flaming",
	"flan",
	"fly",
	"funeral",
	"fungi",
	"fungus",
	"glutton",
	"goblin",
	"hag",
	"hard on",
	"harem",
	"hell",
	"hellhole",
	"ho",
	"hoe",
	"honey-pot",
	"horse",
	"idiot",
	"ignoramus",
	"imbecile",
	"imp",
	"indifferent",
	"injury",
	"jail bait",
	"jerk",
	"jockey",
	"john",
	"killer",
	"land lord",
	"liquidator",
	"lowlife",
	"master baiter",
	"molester",
	"monkey",
	"moronic",
	"mouse",
	"muff",
	"murderer",
	"naked",
	"necrophilic",
	"no-one",
	"nobody",
	"nude",
	"nudist",
	"nutsack",
	"offender",
	"ogre",
	"orgasmic",
	"pecker",
	"peepee",
	"penis",
	"peter",
	"phallus",
	"pig",
	"pirate",
	"pretender",
	"punch",
	"punk",
	"pussy",
	"rapist",
	"rat",
	"rectum",
	"retard",
	"rodent",
	"rubber duck",
	"satan",
	"scrotum",
	"scumbag",
	"shitter",
	"sissy",
	"skank",
	"skunk",
	"slapper",
	"slut",
	"snail",
	"snake",
	"sodomite",
	"son-of-a-bitch",
	"son-of-a-whore",
	"spunk",
	"stoner",
	"succubus",
	"sucker",
	"swinger",
	"tard",
	"temptress",
	"terminator",
	"thundercunt",
	"timesink",
	"towelhead",
	"twat",
	"vampire",
	"vomit",
	"waffle",
	"wail",
	"wannabe",
	"weasel",
	"weed",
	"weirdo",
	"whale",
	"willy",
	"witch",
	"x",
	"zoophile",
}
