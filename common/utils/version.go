package utils

import "github.com/hashicorp/go-version"

// 使用注意：需要新增包依赖：github.com/hashicorp/go-version
// 用法:
// CompareVersion(cv," > 3.2.0")
// CompareVersion(cv," = 1.0.0")
func CompareVersion(v string, con string) (bool, error) {
	v1, err := version.NewVersion(v)
	if err != nil {
		return false, err
	}
	constraints, err := version.NewConstraint(con)
	if err != nil {
		return false, err
	}

	return constraints.Check(v1), nil
}
