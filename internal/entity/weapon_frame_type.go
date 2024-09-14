package entity

type WeaponFrameType struct {
	Name  string
	Value int
}

var WeaponFrameTypes = map[string]int{
	"Pistol": 0,
	"SMG":    1,
	"AR":     2,
	"SR":     3,
}

var WeaponFrameTypeNameMap = getWeaponFrameTypeNameMap()

func getWeaponFrameTypeNameMap() map[string]WeaponFrameType {
	m := map[string]WeaponFrameType{}
	for k, v := range WeaponFrameTypes {
		m[k] = WeaponFrameType{
			Name:  k,
			Value: v,
		}
	}
	return m
}