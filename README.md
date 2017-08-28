Anonymize JSON docs

- Replace strings with sha1 hash of string
- Replace numbers with 0
- Anonymize both keys and values, or just values
- Supply a list of regex skip field filters so you can do things like skip keys that start with an underscore (_)
