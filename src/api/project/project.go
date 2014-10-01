package project

import "api/command"
import "api/repository"
import "api/config"
import "fmt"
import "strings"

type Project struct {
	ProjectKey ProjectKey
	ProjectName ProjectName
}

type ProjectKey string
type ProjectName string

type ReleaseType struct {
	ReleaseTypeKey ReleaseTypeKey
	ReleaseTypeName ReleaseTypeName
}

type ReleaseTypeKey string
type ReleaseTypeName string

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

	MAJOR_RELEASE_TYPE ReleaseTypeKey = "major"
	MINOR_RELEASE_TYPE ReleaseTypeKey = "minor"
	HOTFIX_RELEASE_TYPE ReleaseTypeKey = "hotfix"

	MAJOR_RELEASE_TYPE_NAME ReleaseTypeName = "Major"
	MINOR_RELEASE_TYPE_NAME ReleaseTypeName = "Minor"
	HOTFIX_RELEASE_TYPE_NAME ReleaseTypeName = "Hotfix"
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

func GetReleaseTypes() []ReleaseType {
	return []ReleaseType{
		{ MAJOR_RELEASE_TYPE, MAJOR_RELEASE_TYPE_NAME },
		{ MINOR_RELEASE_TYPE, MINOR_RELEASE_TYPE_NAME },
		{ HOTFIX_RELEASE_TYPE, HOTFIX_RELEASE_TYPE_NAME },
	}
}

func IncrementVersionByReleaseType(repositoryName repository.Repository, releaseType ReleaseTypeKey) {
	if releaseType == MAJOR_RELEASE_TYPE {
		IncrementMajorVersion(repositoryName)
	} else if releaseType == MINOR_RELEASE_TYPE {
		IncrementMinorVersion(repositoryName)
	} else if releaseType == HOTFIX_RELEASE_TYPE {
		IncrementHotfixVersion(repositoryName)
	}
}

func IncrementMajorVersion(repositoryName repository.Repository) {
	executeGradleTask(repositoryName, "incrementMajorVersion")
}

func IncrementMinorVersion(repositoryName repository.Repository) {
	executeGradleTask(repositoryName, "incrementMinorVersion")
}

func IncrementHotfixVersion(repositoryName repository.Repository) {
	executeGradleTask(repositoryName, "incrementHotfixVersion")
}

func SetReleaseVersion(repositoryName repository.Repository) {
	executeGradleTask(repositoryName, "setReleaseVersion")
}

func SetDevelopmentVersion(repositoryName repository.Repository) {
	executeGradleTask(repositoryName, "setDevelopmentVersion")
}

func PrintVersion(repositoryName repository.Repository) string {
	versionRaw, err := executeGradleTask(repositoryName, "printVersion")
	if (versionRaw != "" && err == nil) {
		return strings.TrimSpace(strings.Split(versionRaw, "Version:")[1])
	} else {
		return "Unknown";
	}
}

func executeGradleTask(repositoryName repository.Repository, gradleTask string) (string, error) {
	gradlePath := config.GetProperty("gradle.path")
	projectsContainerPath := config.GetProperty("repository.local.path")
	projectPath := fmt.Sprintf(projectsContainerPath, repositoryName)
	output, err := command.ExecuteCommand(fmt.Sprintf("%v -Dorg.gradle.daemon=true -p %v -q %v",
		gradlePath, projectPath, gradleTask))

	return output, err
}
