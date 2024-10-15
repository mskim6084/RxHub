package drugproduct

type DrugProduct struct {
	DrugCode                 uint32 `json:"drug_code"`
	ClassName                string `json:"class_name"`
	DrugIdentificationNumber string `json:"drug_identification_number"`
	BrandName                string `json:"brand_name"`
	Descriptor               string `json:"descriptor"`
	CompanyName              string `json:"company_name"`
	LastUpdateDate           string `json:"last_update_date"`
	ActiveIngredients        []ActiveIngredient
}

type ActiveIngredient struct {
	DrugCode       uint32 `json:"drug_code"`
	IngredientName string `json:"ingredient_name"`
	Strength       string `json:"strength"`
	StrengthUnit   string `json:"strength_unit"`
}
