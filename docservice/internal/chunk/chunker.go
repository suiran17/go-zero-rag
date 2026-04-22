package chunk

import "strings"

// maxChunkChars is the max characters per chunk (~500 tokens for Chinese text).
const maxChunkChars = 800

// Split splits text into chunks by paragraph boundaries, keeping each chunk
// under maxChunkChars characters. Adjacent short paragraphs are merged.
func Split(text string) []string {
	paras := splitParagraphs(text)
	var chunks []string
	cur := ""

	for _, p := range paras {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		// If a single paragraph exceeds the limit, hard-split it.
		if len(p) > maxChunkChars {
			if cur != "" {
				chunks = append(chunks, cur)
				cur = ""
			}
			chunks = append(chunks, hardSplit(p)...)
			continue
		}
		if cur == "" {
			cur = p
		} else if len(cur)+len(p)+1 <= maxChunkChars {
			cur += "\n" + p
		} else {
			chunks = append(chunks, cur)
			cur = p
		}
	}
	if cur != "" {
		chunks = append(chunks, cur)
	}
	return chunks
}

func splitParagraphs(text string) []string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	return strings.Split(text, "\n\n")
}

func hardSplit(text string) []string {
	var parts []string
	runes := []rune(text)
	for len(runes) > 0 {
		end := maxChunkChars
		if end > len(runes) {
			end = len(runes)
		}
		parts = append(parts, string(runes[:end]))
		runes = runes[end:]
	}
	return parts
}
