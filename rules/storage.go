package rules

type Storage interface {
	LoadRules() ([]Rule, error)
}
