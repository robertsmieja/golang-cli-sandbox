package util

const requiredRegistry = "example.com/registry"
const insecureRegistries = "insecure-registries"

func ValidateDockerDaemonConfig(config *map[string]interface{}) bool {
	registries, ok := (*config)[insecureRegistries].([]string)

	if !ok {
		if registries == nil {
			registries = make([]string, 0)
		} else {
			panic(registries)
		}
	}

	for _, registry := range registries {
		if registry == requiredRegistry {
			return true
		}
	}

	return false
}

func AddRequiredRegistry(config *map[string]interface{}) {
	registries, ok := (*config)[insecureRegistries].([]string)

	if !ok {
		if registries == nil {
			registries = make([]string, 0)
		} else {
			panic(registries)
		}
	}

	registries = append(registries, requiredRegistry)
	(*config)[insecureRegistries] = registries
}
