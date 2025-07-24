package component

var Countries = make(map[string]*Country)

func init() {

	cambodia := &Country{"Cambodia", "柬埔寨", "KH", "855", -1}

	Countries["+855"] = cambodia

}

type Country struct {
	englishName string
	chineseName string
	shortName   string
	phoneCode   string
	timeDiff    int
}
