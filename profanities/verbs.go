package profanities

import "github.com/MikkelHJuul/profaneword"

type VerbStringer struct {}

func (_ VerbStringer) String() string {
	return GetVerb()
}

func GetVerb() string {
	randDevice := profaneword.CryptoRand{}
	rand := randDevice.RandMax(int64(len(verbs)))
	return verbs[rand]
}

var verbs = [...]string{
	"addictive",
	"anal",
	"ass-hat",
	"assbanging",
	"autoerotic",
	"bad-breathed",
	"badmouthing",
	"bastard",
	"benign",
	"bestial",
	"big",
	"bitchy",
	"bloody",
	"bollox",
	"brutal",
	"bugger",
	"cancerous",
	"cheese eating",
	"clownhugging",
	"crikey",
	"damn",
	"damning",
	"dark",
	"destroying",
	"dim",
	"dirty",
	"dominating",
	"dorking",
	"drunken",
	"dull",
	"eat",
	"elitist",
	"enlarged",
	"enormous",
	"envious",
	"enviously",
	"eradicating",
	"fanatic",
	"fisting",
	"fucking",
	"ginormous",
	"gluttonous",
	"gonorrheal",
	"grabbing",
	"hairy",
	"half-arsed",
	"hardcore",
	"heartbroken",
	"hentai-addict",
	"homicidal",
	"horny",
	"hot",
	"huge",
	"idiotic",
	"illegitimate",
	"incestuous",
	"injured",
	"injuring",
	"jerking",
	"jihadist",
	"lame",
	"lard eating",
	"lazy",
	"limp",
	"lustful",
	"lusty",
	"malformed",
	"masturbating",
	"milking",
	"misbehaved",
	"misbehaving",
	"molesting",
	"moron",
	"motherfucking",
	"muff-diving",
	"munching",
	"munging",
	"nazi",
	"necrophile",
	"no-brain",
	"no-good",
	"obnoxious",
	"particularly",
	"peeping",
	"pious",
	"pissed off",
	"pissing",
	"poopmongering",
	"proud",
	"raping",
	"rectal",
	"rim-job",
	"rubbish",
	"salty",
	"satanic",
	"shagging",
	"sick",
	"skanky",
	"slapping",
	"slimy",
	"smelly",
	"sodomizing",
	"splooge",
	"spooky",
	"suicidal",
	"swallowing",
	"tea bagging",
	"teabagging",
	"testicular",
	"testifying",
	"tittyfucking",
	"trashy",
	"ugly",
	"un-hung",
	"uninteresting",
	"unwed",
	"vomitting",
	"wrathful",
	"x-rated",
	"yeasty",
}