package config

type cfgKeysBuilder map[string]interface{}

func (c *cfgKeysBuilder) extractVariables(prefix string, in map[string]interface{}) (out []string, err error) {
	for k, v := range in {
		if newMap, ok := v.(cfgKeysBuilder); ok {
			values, err := c.extractVariables(prefix+"_"+k, newMap)
			if err != nil {
				return nil, err
			}
			out = append(out, values...)
		} else {
			k = prefix + "_" + k

			out = append(out, k[1:])
		}
	}
	return out, nil
}
