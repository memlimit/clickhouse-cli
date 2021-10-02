package cli

const MultilineCLIPrefix = ":-] "

func (c *CLI) GetLivePrefixState() (string, bool) {
	return MultilineCLIPrefix, c.isMultilineInputStarted
}
