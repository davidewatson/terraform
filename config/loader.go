package config

// Load loads the Terraform configuration from a given file.
//
// This file can be any format that Terraform recognizes, and import any
// other format that Terraform recognizes.
func Load(path string) (*Config, error) {
	importTree, err := loadTree(path)
	if err != nil {
		return nil, err
	}

	configTree, err := importTree.ConfigTree()
	if err != nil {
		return nil, err
	}

	return configTree.Flatten()
}
