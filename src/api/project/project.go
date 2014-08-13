package project

type Project struct {
	ProjectKey ProjectKey
	ProjectName ProjectName
}

type ProjectKey string
type ProjectName string

const (
	FAVOR_PROJECT ProjectKey = "favor"
	COSECHA_PROJECT ProjectKey = "cosecha"
	OPERADORA_PROJECT ProjectKey = "operadora"
	CARTERO_PROJECT ProjectKey = "cartero"
	TEST_PROJECT ProjectKey = "gittest"

	FAVOR_PROJECT_NAME ProjectName = "Favor"
	COSECHA_PROJECT_NAME ProjectName = "Cosecha"
	OPERADORA_PROJECT_NAME ProjectName = "Operadora"
	CARTERO_PROJECT_NAME ProjectName = "Cartero"
	TEST_PROJECT_NAME ProjectName = "Gittest"
)

func GetProjects() []Project {
	return []Project{
		{ FAVOR_PROJECT, FAVOR_PROJECT_NAME },
		{ COSECHA_PROJECT, COSECHA_PROJECT_NAME },
		{ OPERADORA_PROJECT, OPERADORA_PROJECT_NAME },
		{ CARTERO_PROJECT, CARTERO_PROJECT_NAME },
		{ TEST_PROJECT, TEST_PROJECT_NAME },
	}
}
