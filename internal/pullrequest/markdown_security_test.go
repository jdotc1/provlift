package pullrequest

import (
    "strings"
    "testing"

    "github.com/provlift/provlift/internal/upgrade"
)

func TestSafeMarkdownNeutralizesUntrustedLinkAndImageSyntax(t *testing.T) {
    input := "[phishing](https://example.invalid) ![badge](https://example.invalid/image)"
    got := safeMarkdown(input)
    if strings.Contains(got, "](https://") {
        t.Fatalf("untrusted Markdown link survived: %q", got)
    }
    for _, escaped := range []string{
        `\[phishing\]\(https://example.invalid\)`,
        `\!\[badge\]\(https://example.invalid/image\)`,
    } {
        if !strings.Contains(got, escaped) {
            t.Errorf("missing escaped syntax %q in %q", escaped, got)
        }
    }
}

func TestPRBodyDoesNotRenderVersionMetadataAsALink(t *testing.T) {
    plan := &upgrade.Plan{Actions: []upgrade.Action{
        action(t, "envs/prod", "hashicorp/azurerm",
            "[review here](https://example.invalid)", "4.19.0", upgrade.LevelMinor, false),
    }}
    pr := Build(plan, nil, nil)
    if strings.Contains(pr.Body, "[review here](https://example.invalid)") {
        t.Fatalf("version metadata produced a clickable link: %s", pr.Body)
    }
    if !strings.Contains(pr.Body, `\[review here\]\(https://example.invalid\)`) {
        t.Errorf("escaped version metadata missing: %s", pr.Body)
    }
}
